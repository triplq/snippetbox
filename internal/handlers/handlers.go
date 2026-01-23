package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/helpers"
	"github.com/triplq/snippetbox/internal/models"
)

func Home(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			helpers.NotFound(w)
			return
		}

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/pages/home.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
			helpers.ServerError(w, err)
		}
	}
}

func ShowSnippets(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			helpers.NotFound(w)
			return
		}

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				helpers.NotFound(w)
			} else {
				helpers.ServerError(w, err)
			}
			return
		}

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/pages/view.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		if err = ts.ExecuteTemplate(w, "base", snippet); err != nil {
			helpers.ServerError(w, err)
		}
	}
}

func CreateSnippet(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			helpers.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		title := "snail"
		content := "snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			helpers.ServerError(w, err)
		}

		fmt.Fprint(w, "Creating a snippet...", id)
	}
}
