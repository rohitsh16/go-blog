package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rohitsh16/go-blog/backend/dto"
)

// CreatePostHandler serves the create page on GET and accepts both
// form POSTs (from the HTML form) and JSON POSTs (API clients).
// Routes:
//
//	GET  /create       -> renders create.html
//	POST /create       -> accepts form submission and redirects to '/'
//	POST /v1/post      -> accepts JSON and returns created post as JSON (registered separately)
func (s *Service) CreatePostHandler(c *gin.Context) {
	// If GET, render the create HTML page
	if c.Request.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("../frontend/templates/create.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Template error: %v", err)
			return
		}
		tmpl.Execute(c.Writer, nil)
		return
	}

	// For POST, decide whether it's JSON or form based on Content-Type
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var req dto.CreatePost
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// prepare values for optional fields
		slugVal := req.Slug
		var authorID interface{}
		if req.AuthorID != nil {
			authorID = *req.AuthorID
		} else {
			authorID = nil
		}
		var authorName interface{}
		if req.AuthorName != nil {
			authorName = *req.AuthorName
		} else {
			authorName = nil
		}
		published := true
		if req.Published != nil {
			published = *req.Published
		}

		result, err := s.Mysql.Exec("INSERT INTO posts (title, content, slug, author_id, author_name, published) VALUES (?, ?, ?, ?, ?, ?)", req.Title, req.Content, slugVal, authorID, authorName, published)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert post"})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve inserted id"})
			return
		}

		var createdAt time.Time
		var updatedAt sql.NullTime
		var dbSlug sql.NullString
		var dbAuthorID sql.NullInt64
		var dbAuthorName sql.NullString
		var dbPublished bool
		err = s.Mysql.QueryRow("SELECT slug, author_id, author_name, published, created_at, updated_at FROM posts WHERE id = ?", id).
			Scan(&dbSlug, &dbAuthorID, &dbAuthorName, &dbPublished, &createdAt, &updatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch created post metadata"})
			return
		}

		// map DB nullable fields to DTO pointers
		var aID *int64
		if dbAuthorID.Valid {
			v := dbAuthorID.Int64
			aID = &v
		}
		var aName *string
		if dbAuthorName.Valid {
			v := dbAuthorName.String
			aName = &v
		}
		var slugOut string
		if dbSlug.Valid {
			slugOut = dbSlug.String
		}
		var uAt *time.Time
		if updatedAt.Valid {
			t := updatedAt.Time
			uAt = &t
		}

		resp := dto.NewPostResponse(id, req.Title, req.Content, slugOut, aID, aName, dbPublished, createdAt, uAt)
		c.JSON(http.StatusCreated, resp)
		return
	}

	// Otherwise treat as form submission
	title := c.PostForm("title")
	content := c.PostForm("content")
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		c.String(http.StatusBadRequest, "title and content are required")
		return
	}

	authorName := c.PostForm("author_name")
	slugVal := c.PostForm("slug")
	publishedVal := c.PostForm("published")
	published := false
	if publishedVal == "1" || strings.ToLower(publishedVal) == "on" {
		published = true
	}

	_, err := s.Mysql.Exec("INSERT INTO posts (title, content, slug, author_name, published) VALUES (?, ?, ?, ?, ?)", title, content, slugVal, authorName, published)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to insert post")
		return
	}

	// Redirect back to home page after creating via form
	c.Redirect(http.StatusSeeOther, "/")
}
