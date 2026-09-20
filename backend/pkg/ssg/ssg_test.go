package ssg_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rohitsh16/go-blog/backend/pkg/ssg"
)

func TestParsePost_WithFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	postFile := filepath.Join(tmpDir, "test-post.md")

	content := `---
title: "Testing SSG Architecture"
slug: "testing-ssg-arch"
author: "Test Author"
date: "2026-09-20"
published: true
summary: "A test summary for SSG."
---

# Heading 1

This is a **bold** paragraph with *italics* and ` + "`inline code`" + `.

- Item 1
- Item 2

> A great quote

` + "```go\nfunc main() {}\n```"

	if err := os.WriteFile(postFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test post: %v", err)
	}

	post, err := ssg.ParsePost(postFile)
	if err != nil {
		t.Fatalf("ParsePost failed: %v", err)
	}

	if post.Title != "Testing SSG Architecture" {
		t.Errorf("expected title 'Testing SSG Architecture', got %s", post.Title)
	}
	if post.Slug != "testing-ssg-arch" {
		t.Errorf("expected slug 'testing-ssg-arch', got %s", post.Slug)
	}
	if post.Author != "Test Author" {
		t.Errorf("expected author 'Test Author', got %s", post.Author)
	}
	if !post.Published {
		t.Errorf("expected Published true")
	}
	if post.Summary != "A test summary for SSG." {
		t.Errorf("expected summary, got %s", post.Summary)
	}

	expectedYear := 2026
	if post.Date.Year() != expectedYear {
		t.Errorf("expected year %d, got %d", expectedYear, post.Date.Year())
	}

	htmlStr := string(post.HTMLContent)
	if !strings.Contains(htmlStr, "<h1>Heading 1</h1>") {
		t.Errorf("expected h1 tag in rendered HTML, got %s", htmlStr)
	}
	if !strings.Contains(htmlStr, "<strong>bold</strong>") {
		t.Errorf("expected strong tag in rendered HTML, got %s", htmlStr)
	}
	if !strings.Contains(htmlStr, "<em>italics</em>") {
		t.Errorf("expected em tag in rendered HTML, got %s", htmlStr)
	}
	if !strings.Contains(htmlStr, "<code>inline code</code>") {
		t.Errorf("expected code tag in rendered HTML, got %s", htmlStr)
	}
	if !strings.Contains(htmlStr, "<ul>") || !strings.Contains(htmlStr, "<li>Item 1</li>") {
		t.Errorf("expected ul list in rendered HTML, got %s", htmlStr)
	}
	if !strings.Contains(htmlStr, "<blockquote><p>A great quote</p></blockquote>") {
		t.Errorf("expected blockquote in rendered HTML, got %s", htmlStr)
	}
	if !strings.Contains(htmlStr, "<pre><code class=\"language-go\">func main() {}\n</code></pre>") {
		t.Errorf("expected code block in rendered HTML, got %s", htmlStr)
	}
}

func TestParsePost_FallbackMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	postFile := filepath.Join(tmpDir, "simple-note.md")

	content := "Just some plain markdown without frontmatter."
	if err := os.WriteFile(postFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test post: %v", err)
	}

	post, err := ssg.ParsePost(postFile)
	if err != nil {
		t.Fatalf("ParsePost failed: %v", err)
	}

	if post.Title != "simple-note" {
		t.Errorf("expected filename fallback title 'simple-note', got %s", post.Title)
	}
	if post.Slug != "simple-note" {
		t.Errorf("expected slug 'simple-note', got %s", post.Slug)
	}
	if !post.Published {
		t.Errorf("expected default published true")
	}
	if post.Date.IsZero() {
		t.Errorf("expected non-zero mod time date")
	}
}

func TestLoadPosts(t *testing.T) {
	tmpDir := t.TempDir()

	post1 := filepath.Join(tmpDir, "post1.md")
	_ = os.WriteFile(post1, []byte("---\ntitle: Post 1\ndate: 2026-01-01\npublished: true\n---\nBody 1"), 0644)

	post2 := filepath.Join(tmpDir, "post2.md")
	_ = os.WriteFile(post2, []byte("---\ntitle: Post 2\ndate: 2026-06-01\npublished: true\n---\nBody 2"), 0644)

	postDraft := filepath.Join(tmpDir, "draft.md")
	_ = os.WriteFile(postDraft, []byte("---\ntitle: Draft\ndate: 2026-03-01\npublished: false\n---\nDraft body"), 0644)

	posts, err := ssg.LoadPosts(tmpDir)
	if err != nil {
		t.Fatalf("LoadPosts failed: %v", err)
	}

	// Draft should be excluded
	if len(posts) != 2 {
		t.Fatalf("expected 2 published posts, got %d", len(posts))
	}

	// Should be sorted by date descending (Post 2 first, then Post 1)
	if posts[0].Title != "Post 2" || posts[1].Title != "Post 1" {
		t.Errorf("expected descending sort order, got: %s, %s", posts[0].Title, posts[1].Title)
	}
}

func TestGenerator_BuildEndToEnd(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	templateDir := filepath.Join(tmpDir, "templates")
	staticDir := filepath.Join(tmpDir, "static")
	outDir := filepath.Join(tmpDir, "dist")

	_ = os.MkdirAll(contentDir, 0755)
	_ = os.MkdirAll(templateDir, 0755)
	_ = os.MkdirAll(staticDir, 0755)

	// Sample post
	postFile := filepath.Join(contentDir, "hello.md")
	_ = os.WriteFile(postFile, []byte("---\ntitle: Hello World\nslug: hello-world\ndate: 2026-09-20\npublished: true\n---\nWelcome!"), 0644)

	// Sample templates
	indexTmpl := filepath.Join(templateDir, "static_index.html")
	_ = os.WriteFile(indexTmpl, []byte("<html><body>Index: {{len .Posts}} posts</body></html>"), 0644)

	postTmpl := filepath.Join(templateDir, "post.html")
	_ = os.WriteFile(postTmpl, []byte("<html><body>Post: {{.Title}} - {{.HTMLContent}}</body></html>"), 0644)

	notFoundTmpl := filepath.Join(templateDir, "404.html")
	_ = os.WriteFile(notFoundTmpl, []byte("<html><body>Not Found</body></html>"), 0644)

	// Sample static asset
	cssFile := filepath.Join(staticDir, "style.css")
	_ = os.WriteFile(cssFile, []byte("body { background: #fff; }"), 0644)

	cfg := ssg.Config{
		ContentDir:  contentDir,
		TemplateDir: templateDir,
		StaticDir:   staticDir,
		OutputDir:   outDir,
		BasePath:    "/test-blog",
		BlogTitle:   "Test Blog",
		BlogURL:     "https://example.com",
	}

	gen := ssg.NewGenerator(cfg)
	if err := gen.Build(); err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify generated files
	expectedFiles := []string{
		filepath.Join(outDir, "index.html"),
		filepath.Join(outDir, "404.html"),
		filepath.Join(outDir, ".nojekyll"),
		filepath.Join(outDir, "feed.xml"),
		filepath.Join(outDir, "static", "style.css"),
		filepath.Join(outDir, "posts", "hello-world", "index.html"),
	}

	for _, f := range expectedFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Errorf("expected generated file %s does not exist", f)
		}
	}
}
