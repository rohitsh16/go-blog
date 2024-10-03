package main

import (
	"go-blog/handlerlogic"
	"go-blog/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	storage.InitDB()
	r := mux.NewRouter()

	handlerlogic.HandlerRoutes(r)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
