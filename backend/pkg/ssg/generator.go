package ssg

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	ContentDir  string
	TemplateDir string
	StaticDir   string
	OutputDir   string
	BasePath    string
	BlogTitle   string
	BlogURL     string
}

type Generator struct {
	cfg Config
}

func NewGenerator(cfg Config) *Generator {
	// Normalize base path: e.g. "/go-blog" or ""
	cfg.BasePath = strings.TrimSuffix(cfg.BasePath, "/")
	if cfg.BasePath != "" && !strings.HasPrefix(cfg.BasePath, "/") {
		cfg.BasePath = "/" + cfg.BasePath
	}
	if cfg.BlogTitle == "" {
		cfg.BlogTitle = "My Blog"
	}
	return &Generator{cfg: cfg}
}

func (g *Generator) Build() error {
	// 1. Ensure output directory exists and is clean
	if err := os.MkdirAll(g.cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// 2. Load all posts
	posts, err := LoadPosts(g.cfg.ContentDir)
	if err != nil {
		return fmt.Errorf("loading posts: %w", err)
	}

	// Assign BasePath to all posts
	for _, p := range posts {
		p.BasePath = g.cfg.BasePath
	}

	// 3. Render Index page
	indexTmplPath := filepath.Join(g.cfg.TemplateDir, "static_index.html")
	indexTmpl, err := template.ParseFiles(indexTmplPath)
	if err != nil {
		return fmt.Errorf("parsing index template: %w", err)
	}

	indexPath := filepath.Join(g.cfg.OutputDir, "index.html")
	indexFile, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("creating index.html: %w", err)
	}
	defer indexFile.Close()

	indexData := struct {
		Posts    []*Post
		BasePath string
		Title    string
	}{
		Posts:    posts,
		BasePath: g.cfg.BasePath,
		Title:    g.cfg.BlogTitle,
	}

	if err := indexTmpl.Execute(indexFile, indexData); err != nil {
		return fmt.Errorf("rendering index template: %w", err)
	}

	// 4. Render Post pages
	postTmplPath := filepath.Join(g.cfg.TemplateDir, "post.html")
	postTmpl, err := template.ParseFiles(postTmplPath)
	if err != nil {
		return fmt.Errorf("parsing post template: %w", err)
	}

	for _, p := range posts {
		postDir := filepath.Join(g.cfg.OutputDir, "posts", p.Slug)
		if err := os.MkdirAll(postDir, 0755); err != nil {
			return fmt.Errorf("creating post dir %s: %w", postDir, err)
		}

		postFilePath := filepath.Join(postDir, "index.html")
		pf, err := os.Create(postFilePath)
		if err != nil {
			return fmt.Errorf("creating post file %s: %w", postFilePath, err)
		}

		if err := postTmpl.Execute(pf, p); err != nil {
			pf.Close()
			return fmt.Errorf("rendering post %s: %w", p.Slug, err)
		}
		pf.Close()
	}

	// 5. Render 404 page
	notFoundTmplPath := filepath.Join(g.cfg.TemplateDir, "404.html")
	if _, err := os.Stat(notFoundTmplPath); err == nil {
		notFoundTmpl, err := template.ParseFiles(notFoundTmplPath)
		if err == nil {
			nfFile, err := os.Create(filepath.Join(g.cfg.OutputDir, "404.html"))
			if err == nil {
				_ = notFoundTmpl.Execute(nfFile, struct{ BasePath string }{BasePath: g.cfg.BasePath})
				nfFile.Close()
			}
		}
	}

	// 6. Copy static assets
	staticDest := filepath.Join(g.cfg.OutputDir, "static")
	if err := copyDir(g.cfg.StaticDir, staticDest); err != nil {
		return fmt.Errorf("copying static assets: %w", err)
	}

	// 7. Write RSS Feed (feed.xml)
	if err := g.generateRSS(posts); err != nil {
		// Log but don't fail build
		fmt.Printf("warning: generating RSS feed: %v\n", err)
	}

	// 8. Create .nojekyll for GitHub Pages (prevents Jekyll processing)
	noJekyllPath := filepath.Join(g.cfg.OutputDir, ".nojekyll")
	_ = os.WriteFile(noJekyllPath, []byte(""), 0644)

	return nil
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author,omitempty"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"item"`
}

type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

func (g *Generator) generateRSS(posts []*Post) error {
	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       g.cfg.BlogTitle,
			Link:        g.cfg.BlogURL + g.cfg.BasePath,
			Description: "Latest blog posts",
		},
	}

	for _, p := range posts {
		link := fmt.Sprintf("%s%s/posts/%s/", g.cfg.BlogURL, g.cfg.BasePath, p.Slug)
		rss.Channel.Items = append(rss.Channel.Items, RSSItem{
			Title:       p.Title,
			Link:        link,
			Description: p.Summary,
			PubDate:     p.Date.Format(time.RFC1123Z),
			Author:      p.Author,
		})
	}

	data, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}

	fullXML := []byte(xml.Header + string(data) + "\n")
	return os.WriteFile(filepath.Join(g.cfg.OutputDir, "feed.xml"), fullXML, 0644)
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
