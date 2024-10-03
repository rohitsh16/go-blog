package handlerlogic

import (
	"go-blog/dto"
	"net/http"
	"text/template"
)

// HomeHandler serves the homepage with a list of blog posts
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	posts := []dto.BlogPost{
		{Title: "blog post", Content: "contents of the blog!"},
	}

	tmpl, err := template.ParseFiles("../frontend/templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, posts)
}
