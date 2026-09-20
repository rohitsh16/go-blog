package handler

import (
	"html/template"
	"os"
	"path/filepath"
)

// ResolveTemplatePath returns the correct path to a template file
// whether the working directory is the repository root or the backend/ directory.
func ResolveTemplatePath(filename string) string {
	candidates := []string{
		filepath.Join("frontend", "templates", filename),
		filepath.Join("..", "frontend", "templates", filename),
		filepath.Join("..", "..", "frontend", "templates", filename),
		filepath.Join("templates", filename),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	// Fallback to relative parent path
	return filepath.Join("..", "frontend", "templates", filename)
}

// ParseTemplate parses a template by filename using directory-agnostic resolution.
func ParseTemplate(filename string) (*template.Template, error) {
	path := ResolveTemplatePath(filename)
	return template.ParseFiles(path)
}
