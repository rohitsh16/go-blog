package handler

import (
	"database/sql"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rohitsh16/go-blog/backend/dto"
)

// HomeHandler serves the homepage with a list of blog posts fetched from MySQL.
// It expects to be used as a method on Service so it can access the DB connection.
func (s *Service) HomeHandler(c *gin.Context) {
	// only display published posts to users
	rows, err := s.Mysql.Query("SELECT id, title, content, slug, author_id, author_name, published, created_at, updated_at FROM posts WHERE published = 1 ORDER BY created_at DESC")
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

	// Use the user-facing index template which expects a slice of PostResponse
	tmpl, err := template.ParseFiles("../frontend/templates/index_user.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "template error: %v", err)
		return
	}

	// Execute the template writing to the gin context's writer
	if err := tmpl.Execute(c.Writer, posts); err != nil {
		c.String(http.StatusInternalServerError, "template execution error: %v", err)
		return
	}
}

// AdminHandler renders an admin view of posts (management features).
func (s *Service) AdminHandler(c *gin.Context) {
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

	tmpl, err := template.ParseFiles("../frontend/templates/index_admin.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "template error: %v", err)
		return
	}

	if err := tmpl.Execute(c.Writer, posts); err != nil {
		c.String(http.StatusInternalServerError, "template execution error: %v", err)
		return
	}
}
