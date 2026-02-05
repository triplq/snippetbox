package main

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/handlers"
	"github.com/triplq/snippetbox/internal/middleware"
)

func routes(app *application.Application) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", handlers.Home(app))
	mux.HandleFunc("GET /snippets/view/{id}", handlers.ShowSnippets(app))
	mux.HandleFunc("GET /snippets/create", handlers.CreateSnippet(app))
	mux.HandleFunc("POST /snippets/create", handlers.PostCreateSnippet(app))

	chain := alice.New(middleware.PanicRecover, middleware.SlogRequest, middleware.SecureHeaders)

	return chain.Then(mux)
}
