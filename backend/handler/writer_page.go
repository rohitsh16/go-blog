package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rohitsh16/go-blog/backend/dto"
)

// WriterDashboardHandler shows a writer-facing dashboard with links to create/edit posts.
// Currently this lists all posts (no auth implemented).
func (s *Service) WriterDashboardHandler(c *gin.Context) {
	rows, err := s.Mysql.Query("SELECT id, title, content, slug, author_id, author_name, published, created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to query posts: %v", err)
		return
	}
	defer rows.Close()

	posts := make([]dto.PostResponse, 0)
	for rows.Next() {
		var id int64
		var title string
		var content string
		var slug sql.NullString
		var authorID sql.NullInt64
		var authorName sql.NullString
		var published bool
		var createdAt sql.NullTime
		var updatedAt sql.NullTime

		if err := rows.Scan(&id, &title, &content, &slug, &authorID, &authorName, &published, &createdAt, &updatedAt); err != nil {
			c.String(http.StatusInternalServerError, "failed to scan post: %v", err)
			return
		}

		var slugVal string
		if slug.Valid {
			slugVal = slug.String
		}

		var aID *int64
		if authorID.Valid {
			v := authorID.Int64
			aID = &v
		}

		var aName *string
		if authorName.Valid {
			v := authorName.String
			aName = &v
		}

		var uAt *time.Time
		if updatedAt.Valid {
			t := updatedAt.Time
			uAt = &t
		}

		posts = append(posts, dto.NewPostResponse(id, title, content, slugVal, aID, aName, published, createdAt.Time, uAt))
	}

	tmpl, err := template.ParseFiles("../frontend/templates/writer_dashboard.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "template error: %v", err)
		return
	}
	if err := tmpl.Execute(c.Writer, posts); err != nil {
		c.String(http.StatusInternalServerError, "template execution error: %v", err)
		return
	}
}

// WriterEditHandler supports GET to render an edit form and POST to update a post.
func (s *Service) WriterEditHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid id")
		return
	}

	if c.Request.Method == http.MethodGet {
		var title, content string
		var slug sql.NullString
		var authorID sql.NullInt64
		var authorName sql.NullString
		var published bool
		var createdAt sql.NullTime
		var updatedAt sql.NullTime

		err := s.Mysql.QueryRow("SELECT title, content, slug, author_id, author_name, published, created_at, updated_at FROM posts WHERE id = ?", id).
			Scan(&title, &content, &slug, &authorID, &authorName, &published, &createdAt, &updatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				c.String(http.StatusNotFound, "post not found")
				return
			}
			c.String(http.StatusInternalServerError, "failed to fetch post: %v", err)
			return
		}

		var slugVal string
		if slug.Valid {
			slugVal = slug.String
		}

		var aID *int64
		if authorID.Valid {
			v := authorID.Int64
			aID = &v
		}

		var aName *string
		if authorName.Valid {
			v := authorName.String
			aName = &v
		}

		var uAt *time.Time
		if updatedAt.Valid {
			t := updatedAt.Time
			uAt = &t
		}

		data := dto.NewPostResponse(id, title, content, slugVal, aID, aName, published, createdAt.Time, uAt)
		tmpl, err := template.ParseFiles("../frontend/templates/writer_edit.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "template error: %v", err)
			return
		}
		if err := tmpl.Execute(c.Writer, data); err != nil {
			c.String(http.StatusInternalServerError, "template execution error: %v", err)
			return
		}
		return
	}

	// POST - update
	title := c.PostForm("title")
	content := c.PostForm("content")
	authorName := c.PostForm("author_name")
	slugVal := c.PostForm("slug")
	publishedVal := c.PostForm("published")
	published := false
	if publishedVal == "1" || publishedVal == "on" {
		published = true
	}
	if title == "" || content == "" {
		c.String(http.StatusBadRequest, "title and content are required")
		return
	}

	_, err = s.Mysql.Exec("UPDATE posts SET title = ?, content = ?, slug = ?, author_name = ?, published = ? WHERE id = ?", title, content, slugVal, authorName, published, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to update post: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/writer")
}
