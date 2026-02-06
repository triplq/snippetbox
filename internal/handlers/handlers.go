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

type SnippetForm struct {
	title   string
	content string
	expires int
	errors  map[string]string
}

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

		templateForm := new(SnippetForm)

		templateForm.title = r.PostForm.Get("title")
		templateForm.content = r.PostForm.Get("content")
		templateForm.expires, err = strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(templateForm.title) == "" {
			templateForm.errors["title"] = "This field can not be empty"
		} else if utf8.RuneCountInString(templateForm.title) > 100 {
			templateForm.errors["title"] = "This filed can not be large then 100"
		}

		if utf8.RuneCountInString(templateForm.content) > 100 {
			templateForm.errors["content"] = "This filed can not be large then 100"

		}

		if templateForm.expires != 1 && templateForm.expires != 7 && templateForm.expires != 365 {
			templateForm.errors["expires"] = "This field must equals 1, 7, 365"
		}

		if len(templateForm.errors) > 0 {
			data := templates.NewTemplateData(r)
			data.Form = templateForm
			err := app.Render(w, http.StatusUnprocessableEntity, "create.html", data)
			if err != nil {
				helpers.ServerError(w, err)
			}
			return
		}

		id, err := app.Snippets.Insert(templateForm.title, templateForm.content, templateForm.expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	}
}
