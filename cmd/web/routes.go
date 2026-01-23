package main

import (
	"net/http"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/handlers"
)

func routes(app *application.Application) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handlers.Home(app))
	mux.HandleFunc("/snippets/view", handlers.ShowSnippets(app))
	mux.HandleFunc("/snippets/create", handlers.CreateSnippet(app))

	return mux
}
