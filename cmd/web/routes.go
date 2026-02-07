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

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	mux.Handle("GET /", dynamic.ThenFunc(handlers.Home(app)))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(handlers.ShowSnippets(app)))
	mux.Handle("GET /snippets/create", dynamic.ThenFunc(handlers.CreateSnippet(app)))
	mux.Handle("POST /snippets/create", dynamic.ThenFunc(handlers.PostCreateSnippet(app)))

	chain := alice.New(middleware.PanicRecover, middleware.SlogRequest, middleware.SecureHeaders)

	return chain.Then(mux)
}
