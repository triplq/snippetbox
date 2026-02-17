package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/helpers"
	"github.com/triplq/snippetbox/internal/types"
)

func SecureHeaders(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
			w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")

			next.ServeHTTP(w, r)
		})
	}
}

func SlogRequest(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("received request",
				"ip", r.RemoteAddr,
				"proto", r.Proto,
				"method", r.Method,
				"uri", r.URL.RequestURI())

			next.ServeHTTP(w, r)
		})
	}
}

func PanicRecover(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					helpers.ServerError(w, fmt.Errorf("%s", err), app)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func AuthIsRequired(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !app.IsAuthenticated(r) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			w.Header().Set("Cache-Control", "no-store")

			next.ServeHTTP(w, r)
		})
	}
}

func NoSurf(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
		})

		return csrfHandler
	}
}

func Authenticate(app *application.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}

			exists, err := app.Users.Exists(id)
			if err != nil {
				helpers.ServerError(w, err, app)
				return
			}

			if exists {
				ctx := context.WithValue(r.Context(), types.IsAuthenticatedContextKey, true)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
