package server

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/handlers"
	"github.com/triplq/snippetbox/internal/middleware"
	"github.com/triplq/snippetbox/ui"
)

func Routes(app *application.Application) http.Handler {
	mux := http.NewServeMux()

	staticChain := alice.New(
		middleware.PanicRecover(app),
		middleware.SecureHeaders(app))

	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("GET /static/", staticChain.Then(fileServer))

	dynamicChain := alice.New(
		middleware.PanicRecover(app),
		middleware.SecureHeaders(app),
		middleware.SlogRequest(app),
		app.SessionManager.LoadAndSave,
		middleware.NoSurf(app),
		middleware.Authenticate(app))

	mux.Handle("GET /", dynamicChain.ThenFunc(handlers.Home(app)))
	mux.Handle("GET /snippets/view/{id}", dynamicChain.ThenFunc(handlers.ShowSnippet(app)))
	mux.Handle("GET /user/login", dynamicChain.ThenFunc(handlers.LogIn(app)))
	mux.Handle("POST /user/login", dynamicChain.ThenFunc(handlers.LogInPost(app)))
	mux.Handle("GET /user/signup", dynamicChain.ThenFunc(handlers.SignUp(app)))
	mux.Handle("POST /user/signup", dynamicChain.ThenFunc(handlers.SignUpPost(app)))

	protectedChain := dynamicChain.Append(middleware.AuthIsRequired(app))

	mux.Handle("GET /snippets/create", protectedChain.ThenFunc(handlers.CreateSnippet(app)))
	mux.Handle("POST /snippets/create", protectedChain.ThenFunc(handlers.CreateSnippetPost(app)))
	mux.Handle("POST /user/logout", protectedChain.ThenFunc(handlers.LogOutPost(app)))

	return mux
}
