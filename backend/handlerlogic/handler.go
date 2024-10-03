package handlerlogic

import "github.com/gorilla/mux"

func HandlerRoutes(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/create", CreatePostHandler).Methods("GET", "POST")
}
