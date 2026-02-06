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
	Title   string
	Content string
	Expires int
	Errors  map[string]string
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
		data.Form = SnippetForm{
			Expires: 365,
		}
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

		exp, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		templateForm := SnippetForm{
			Title:   r.PostForm.Get("title"),
			Content: r.PostForm.Get("content"),
			Expires: exp,
			Errors:  make(map[string]string),
		}

		if strings.TrimSpace(templateForm.Title) == "" {
			templateForm.Errors["title"] = "This field can not be empty"
		} else if utf8.RuneCountInString(templateForm.Title) > 100 {
			templateForm.Errors["title"] = "This filed can not be large then 100"
		}

		if strings.TrimSpace(templateForm.Content) == "" {
			templateForm.Errors["content"] = "This field can not be empty"
		} else if utf8.RuneCountInString(templateForm.Content) > 100 {
			templateForm.Errors["content"] = "This filed can not be large then 100"
		}

		if templateForm.Expires != 1 && templateForm.Expires != 7 && templateForm.Expires != 365 {
			templateForm.Errors["expires"] = "This field must equals 1, 7, 365"
		}

		if len(templateForm.Errors) > 0 {
			data := templates.NewTemplateData(r)
			data.Form = templateForm
			err := app.Render(w, http.StatusUnprocessableEntity, "create.html", data)
			if err != nil {
				helpers.ServerError(w, err)
			}
			return
		}

		id, err := app.Snippets.Insert(templateForm.Title, templateForm.Content, templateForm.Expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	}
}
