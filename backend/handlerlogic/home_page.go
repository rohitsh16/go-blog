package handlerlogic

import (
	"go-blog/dto"
	"net/http"
	"text/template"
)

// HomeHandler serves the homepage with a list of blog posts
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts := []dto.BlogPost{
		{Title: "My First Blog Post", Content: "Welcome to my blog!"},
		{Title: "Learning Golang", Content: "Today, I learned how to build a simple web application with Go."},
	}

	tmpl, err := template.ParseFiles("../Frontend/templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, posts)
}
