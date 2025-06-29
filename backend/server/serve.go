package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rohitsh16/go-blog/backend/db"
	"github.com/rohitsh16/go-blog/backend/handler"
	"github.com/rohitsh16/go-blog/backend/server/config"
)

type Server struct {
	httpServer *http.Server
	Config     *config.Config
}

// creates new server
func NewServer() (*Server, error) {
	cfg, cfgErr := config.LoadConfig("config.yaml")
	if cfgErr != nil {
		log.Fatalf("Unable to load config, %v", cfgErr)
	}

	mysql, dbErr := db.GetInstance("mysql", cfg.DatabaseConfig)
	if dbErr != nil {
		log.Fatalf("Unable to init Mysql DB, %v", dbErr)
	}

	router := gin.Default() // previously wasn't initiaised it in single object, hence error

	service := handler.Service{
		Config:       cfg,
		Mysql:        mysql,
		GinFramework: router,
	}

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	handler.RegisterRoutes(service)

	return &Server{
		httpServer: httpServer,
		Config:     cfg,
	}, nil
}

// Start runs the HTTP server
func (s *Server) Start() error {
	log.Println("Starting server on", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
