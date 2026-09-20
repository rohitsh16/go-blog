.PHONY: all bin/ssg build-static build-pages serve-static new-post build-backend run-backend clean

# Default target builds the static site
all: build-static

# Compile SSG CLI binary if missing or modified
bin/ssg: backend/cmd/ssg/main.go backend/pkg/ssg/*.go
	@mkdir -p bin
	cd backend && go build -o ../bin/ssg ./cmd/ssg

# Generate static HTML site in dist/
build-static: bin/ssg
	./bin/ssg build -out dist

# Generate static HTML site with base path for GitHub Pages
build-pages: bin/ssg
	./bin/ssg build -base "/go-blog" -out dist

# Serve static site locally on http://localhost:3000
serve-static: bin/ssg
	./bin/ssg serve -dir dist -port 3000

# Create a new post in content/posts/ (Usage: make new-post TITLE="Post Title")
new-post: bin/ssg
	@if [ -z "$(TITLE)" ]; then echo "Usage: make new-post TITLE=\"Your Post Title\""; exit 1; fi
	./bin/ssg new -title "$(TITLE)"

# Build dynamic Go backend binary
build-backend:
	@mkdir -p bin
	cd backend && go build -o ../bin/go-blog .

# Run dynamic Go backend
run-backend:
	cd backend && go run main.go

# Clean generated static files
clean:
	rm -rf dist bin/ssg
