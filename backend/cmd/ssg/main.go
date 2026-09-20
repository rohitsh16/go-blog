package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/rohitsh16/go-blog/backend/db"
	"github.com/rohitsh16/go-blog/backend/pkg/ssg"
	"github.com/rohitsh16/go-blog/backend/server/config"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "build":
		runBuild(os.Args[2:])
	case "serve":
		runServe(os.Args[2:])
	case "new":
		runNew(os.Args[2:])
	case "export-db":
		runExportDB(os.Args[2:])
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Static Site Generator (SSG) for go-blog

Usage:
  ssg <command> [flags]

Commands:
  build       Compile markdown posts and templates into static HTML (dist/)
  serve       Preview the static site locally using a lightweight HTTP server
  new         Scaffold a new markdown blog post with frontmatter
  export-db   Export existing MySQL blog posts to markdown files in content/posts/

Flags for 'build':
  -content    Path to markdown content directory (default: "content/posts")
  -templates  Path to HTML templates directory (default: "frontend/templates")
  -static     Path to static assets directory (default: "frontend/static")
  -out        Output directory for static site (default: "dist")
  -base       Base URL path prefix (e.g., "/go-blog" for GitHub Pages, default: "")
  -title      Blog title (default: "My Blog")
  -url        Base website URL (default: "https://rohitsh16.github.io")

Flags for 'serve':
  -dir        Directory to serve (default: "dist")
  -port       HTTP port (default: 3000)
  -base       Base URL path prefix to match build (default: "")

Flags for 'new':
  -title      Title of the new blog post (required)
  -author     Author name (default: "Author")
  -slug       Custom URL slug (optional, auto-derived from title)

Examples:
  # Build for local preview (root base path)
  ssg build

  # Build for GitHub Pages (under repo name subpath)
  ssg build -base "/go-blog"

  # Preview locally
  ssg serve -port 3000

  # Create a new post
  ssg new -title "Building Distributed Systems in Go"`)
}

func findProjectRoot() string {
	// Check current directory
	if _, err := os.Stat("frontend"); err == nil {
		return "."
	}
	// Check parent directory
	if _, err := os.Stat("../frontend"); err == nil {
		return ".."
	}
	return "."
}

func runBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	root := findProjectRoot()

	contentDir := fs.String("content", filepath.Join(root, "content", "posts"), "Content directory")
	templateDir := fs.String("templates", filepath.Join(root, "frontend", "templates"), "Template directory")
	staticDir := fs.String("static", filepath.Join(root, "frontend", "static"), "Static assets directory")
	outputDir := fs.String("out", filepath.Join(root, "dist"), "Output directory")
	basePath := fs.String("base", "", "Base path for URLs")
	blogTitle := fs.String("title", "My Blog", "Blog title")
	blogURL := fs.String("url", "https://rohitsh16.github.io", "Blog canonical URL")

	_ = fs.Parse(args)

	cfg := ssg.Config{
		ContentDir:  *contentDir,
		TemplateDir: *templateDir,
		StaticDir:   *staticDir,
		OutputDir:   *outputDir,
		BasePath:    *basePath,
		BlogTitle:   *blogTitle,
		BlogURL:     *blogURL,
	}

	log.Printf("Starting static site build...")
	log.Printf("Content:   %s", cfg.ContentDir)
	log.Printf("Templates: %s", cfg.TemplateDir)
	log.Printf("Static:    %s", cfg.StaticDir)
	log.Printf("Output:    %s", cfg.OutputDir)
	if cfg.BasePath != "" {
		log.Printf("Base Path: %s", cfg.BasePath)
	}

	gen := ssg.NewGenerator(cfg)
	if err := gen.Build(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	log.Printf("Build successful! Static site generated in: %s", cfg.OutputDir)
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	root := findProjectRoot()

	dir := fs.String("dir", filepath.Join(root, "dist"), "Directory to serve")
	port := fs.Int("port", 3000, "Port to listen on")
	base := fs.String("base", "", "Base path prefix")

	_ = fs.Parse(args)

	// Check if dir exists, if not inform user to run build first
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Printf("Output directory '%s' does not exist. Running static build first...", *dir)
		runBuild(nil)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := ssg.NewServer(*dir, *port, *base)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	title := fs.String("title", "", "Post title (required)")
	author := fs.String("author", "Rohit Shukla", "Author name")
	slug := fs.String("slug", "", "Post slug (optional)")

	_ = fs.Parse(args)

	if *title == "" {
		fmt.Fprintln(os.Stderr, "Error: -title is required")
		fs.Usage()
		os.Exit(1)
	}

	root := findProjectRoot()
	contentDir := filepath.Join(root, "content", "posts")
	_ = os.MkdirAll(contentDir, 0755)

	if *slug == "" {
		*slug = strings.ToLower(strings.ReplaceAll(*title, " ", "-"))
		*slug = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				return r
			}
			return -1
		}, *slug)
	}

	filename := filepath.Join(contentDir, fmt.Sprintf("%s.md", *slug))
	if _, err := os.Stat(filename); err == nil {
		log.Fatalf("Post file already exists: %s", filename)
	}

	today := time.Now().Format("2006-01-02")
	content := fmt.Sprintf(`---
title: %q
slug: %q
author: %q
date: %q
published: true
summary: "Write a brief summary of the article here."
---

Write your article in Markdown here...
`, *title, *slug, *author, today)

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to create post file: %v", err)
	}

	log.Printf("Created new post: %s", filename)
}

func runExportDB(args []string) {
	root := findProjectRoot()
	cfgPath := filepath.Join(root, "backend", "config_files", "config.yaml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		cfgPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Could not load config file (%s): %v", cfgPath, err)
	}

	mysql, err := db.GetInstance("mysql", cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Could not connect to MySQL: %v", err)
	}

	rows, err := mysql.Query("SELECT title, content, slug, author_name, published, created_at FROM posts")
	if err != nil {
		log.Fatalf("Failed to query posts: %v", err)
	}
	defer rows.Close()

	contentDir := filepath.Join(root, "content", "posts")
	_ = os.MkdirAll(contentDir, 0755)

	count := 0
	for rows.Next() {
		var title, body, authorName string
		var slugVal *string
		var published bool
		var createdAt time.Time

		if err := rows.Scan(&title, &body, &slugVal, &authorName, &published, &createdAt); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		slug := ""
		if slugVal != nil && *slugVal != "" {
			slug = *slugVal
		} else {
			slug = strings.ToLower(strings.ReplaceAll(title, " ", "-"))
		}

		filePath := filepath.Join(contentDir, fmt.Sprintf("%s.md", slug))
		postContent := fmt.Sprintf(`---
title: %q
slug: %q
author: %q
date: %q
published: %t
summary: %q
---

%s
`, title, slug, authorName, createdAt.Format("2006-01-02"), published, "", body)

		if err := os.WriteFile(filePath, []byte(postContent), 0644); err != nil {
			log.Printf("Failed to write %s: %v", filePath, err)
			continue
		}
		count++
		log.Printf("Exported: %s -> %s", title, filePath)
	}

	log.Printf("Export completed: %d post(s) exported to %s", count, contentDir)
}
