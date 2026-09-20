package handler

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostViewData struct {
	Title         string
	Summary       string
	Author        string
	FormattedDate string
	HTMLContent   template.HTML
	BasePath      string
}

// PostDetailHandler renders an individual blog post by slug or numeric ID.
// Routes: GET /post/:identifier
func (s *Service) PostDetailHandler(c *gin.Context) {
	identifier := strings.TrimSpace(c.Param("identifier"))
	if identifier == "" {
		c.String(http.StatusBadRequest, "post identifier is required")
		return
	}

	var (
		id         int64
		title      string
		content    string
		slug       sql.NullString
		authorName sql.NullString
		published  bool
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
		queryErr   error
	)

	// Check if identifier is numeric ID or slug
	if numID, err := strconv.ParseInt(identifier, 10, 64); err == nil {
		queryErr = s.Mysql.QueryRow(
			"SELECT id, title, content, slug, author_name, published, created_at, updated_at FROM posts WHERE (id = ? OR slug = ?) AND published = 1",
			numID, identifier,
		).Scan(&id, &title, &content, &slug, &authorName, &published, &createdAt, &updatedAt)
	} else {
		queryErr = s.Mysql.QueryRow(
			"SELECT id, title, content, slug, author_name, published, created_at, updated_at FROM posts WHERE slug = ? AND published = 1",
			identifier,
		).Scan(&id, &title, &content, &slug, &authorName, &published, &createdAt, &updatedAt)
	}

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			// Render 404 template if available
			if tmpl404, err404 := ParseTemplate("404.html"); err404 == nil {
				c.Status(http.StatusNotFound)
				_ = tmpl404.Execute(c.Writer, struct{ BasePath string }{BasePath: ""})
				return
			}
			c.String(http.StatusNotFound, "Post not found")
			return
		}
		c.String(http.StatusInternalServerError, "failed to query post: %v", queryErr)
		return
	}

	author := ""
	if authorName.Valid {
		author = authorName.String
	}

	postDate := time.Now()
	if createdAt.Valid {
		postDate = createdAt.Time
	}

	// Format body content (convert newlines to paragraphs if raw text)
	htmlBody := formatContentHTML(content)

	viewData := PostViewData{
		Title:         title,
		Summary:       "",
		Author:        author,
		FormattedDate: postDate.Format("January 2, 2006"),
		HTMLContent:   template.HTML(htmlBody),
		BasePath:      "",
	}

	tmpl, err := ParseTemplate("post.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "template error: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, viewData); err != nil {
		c.String(http.StatusInternalServerError, "template execution error: %v", err)
		return
	}
}

// formatContentHTML safely converts plaintext or basic markdown to HTML
func formatContentHTML(content string) string {
	if strings.Contains(content, "<p>") || strings.Contains(content, "<div>") {
		return content
	}
	paragraphs := strings.Split(content, "\n\n")
	var sb strings.Builder
	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		escaped := template.HTMLEscapeString(trimmed)
		escaped = strings.ReplaceAll(escaped, "\n", "<br>")
		sb.WriteString("<p>" + escaped + "</p>\n")
	}
	return sb.String()
}
