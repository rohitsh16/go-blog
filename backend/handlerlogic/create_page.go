package handlerlogic

import (
	"log"
	"net/http"
	"text/template"
)

// CreatePostHandler serves the page to create a new blog post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("../frontend/templates/create.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Handle form submission here
		title := r.FormValue("title")
		content := r.FormValue("content")
		log.Printf("New Post: %s - %s", title, content)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
