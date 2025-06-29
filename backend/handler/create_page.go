package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler serves the page to create a new blog post
func (s *Service) CreatePostHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("frontend/templates/create.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Template error: %v", err)
			return
		}
		tmpl.Execute(c.Writer, nil)
	} else if c.Request.Method == http.MethodPost {
		title := c.PostForm("title")
		content := c.PostForm("content")
		log.Printf("New Post: %s - %s", title, content)
		c.Redirect(http.StatusSeeOther, "/")
	}
}
