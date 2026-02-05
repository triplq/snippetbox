package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/helpers"
	"github.com/triplq/snippetbox/internal/models"
	"github.com/triplq/snippetbox/internal/templates"
)

func Home(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			helpers.NotFound(w)
			return
		}

		snippets, err := app.Snippets.Latest()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := templates.NewTemplateData(r)
		data.Snippets = snippets

		err = app.Render(w, http.StatusOK, "home.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func ShowSnippets(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
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

		data := templates.NewTemplateData(r)
		data.Snippet = snippet

		err = app.Render(w, http.StatusOK, "view.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func CreateSnippet(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := templates.NewTemplateData(r)
		err := app.Render(w, http.StatusOK, "create.html", data)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}
}

func PostCreateSnippet(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		errors := make(map[string]string)

		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(title) == "" {
			errors["title"] = "This field can not be empty"
		} else if utf8.RuneCountInString(title) > 100 {
			errors["title"] = "This filed can not be large then 100"
		}

		if utf8.RuneCountInString(content) > 100 {
			errors["content"] = "This filed can not be large then 100"

		}

		if expires != 1 && expires != 7 && expires != 365 {
			errors["expires"] = "This field must equals 1, 7, 365"
		}

		if len(errors) > 0 {
			fmt.Fprint(w, errors)
			return
		}

		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	}
}
