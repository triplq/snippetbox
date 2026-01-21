package main

import (
	"net/http"

	"github.com/triplq/snippetbox/internal/handlers"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handlers.Home())
	mux.HandleFunc("/snippets/view", handlers.ShowSnippets())
	mux.HandleFunc("/snippets/create", handlers.CreateSnippet())

	return mux
}
