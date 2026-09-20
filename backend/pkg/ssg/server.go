package ssg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	dir  string
	port int
	base string
}

func NewServer(dir string, port int, base string) *Server {
	if port <= 0 {
		port = 3000
	}
	base = strings.TrimSuffix(base, "/")
	if base != "" && !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	return &Server{dir: dir, port: port, base: base}
}

func (s *Server) Start(ctx context.Context) error {
	absDir, err := filepath.Abs(s.dir)
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir(absDir))

	mux := http.NewServeMux()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path

		// Strip base prefix if configured
		if s.base != "" && strings.HasPrefix(reqPath, s.base) {
			reqPath = strings.TrimPrefix(reqPath, s.base)
			if reqPath == "" {
				reqPath = "/"
			}
			r.URL.Path = reqPath
		}

		// Clean URL support: if asking for /posts/slug, check for /posts/slug/index.html
		localPath := filepath.Join(absDir, filepath.Clean(reqPath))
		fi, err := os.Stat(localPath)
		if os.IsNotExist(err) {
			// Try .html
			if _, err := os.Stat(localPath + ".html"); err == nil {
				r.URL.Path = reqPath + ".html"
			} else {
				// Serve 404.html if exists
				custom404 := filepath.Join(absDir, "404.html")
				if _, err404 := os.Stat(custom404); err404 == nil {
					w.WriteHeader(http.StatusNotFound)
					http.ServeFile(w, r, custom404)
					return
				}
			}
		} else if err == nil && fi.IsDir() {
			indexFile := filepath.Join(localPath, "index.html")
			if _, errIndex := os.Stat(indexFile); os.IsNotExist(errIndex) {
				custom404 := filepath.Join(absDir, "404.html")
				if _, err404 := os.Stat(custom404); err404 == nil {
					w.WriteHeader(http.StatusNotFound)
					http.ServeFile(w, r, custom404)
					return
				}
			}
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("/", handler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	previewURL := fmt.Sprintf("http://localhost:%d%s/", s.port, s.base)
	log.Printf("Serving static blog from %s", absDir)
	log.Printf("Preview URL: %s", previewURL)
	log.Printf("Press Ctrl+C to stop")

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutCtx)
	case err := <-errCh:
		return err
	}
}
