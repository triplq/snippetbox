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

	staticChain := alice.New(
		middleware.PanicRecover,
		middleware.SecureHeaders)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", staticChain.Then(http.StripPrefix("/static", fileServer)))

	dynamicChain := alice.New(
		middleware.PanicRecover,
		middleware.SecureHeaders,
		middleware.SlogRequest,
		app.SessionManager.LoadAndSave)

	mux.Handle("GET /", dynamicChain.ThenFunc(handlers.Home(app)))
	mux.Handle("GET /snippets/view/{id}", dynamicChain.ThenFunc(handlers.ShowSnippets(app)))
	mux.Handle("GET /snippets/create", dynamicChain.ThenFunc(handlers.CreateSnippet(app)))
	mux.Handle("POST /snippets/create", dynamicChain.ThenFunc(handlers.PostCreateSnippet(app)))
	mux.Handle("GET /snippets/login", dynamicChain.ThenFunc(handlers.LogIn(app)))
	mux.Handle("POST /snippets/login", dynamicChain.ThenFunc(handlers.LogInPost(app)))
	mux.Handle("GET /snippets/signup", dynamicChain.ThenFunc(handlers.SignUp(app)))
	mux.Handle("POST /snippets/signup", dynamicChain.ThenFunc(handlers.SignUp(app)))
	mux.Handle("POST /snippets/logout", dynamicChain.ThenFunc(handlers.LogOut(app)))

	return mux
}
