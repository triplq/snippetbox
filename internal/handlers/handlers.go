package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/helpers"
	"github.com/triplq/snippetbox/internal/models"
	"github.com/triplq/snippetbox/internal/templates"
	"github.com/triplq/snippetbox/internal/validator"
)

type SnippetForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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

		flash := app.SessionManager.PopString(r.Context(), "flash")

		data := templates.NewTemplateData(r)
		data.Snippet = snippet
		data.Flash = flash

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
		var templateForm SnippetForm

		err := app.DecodePostForm(r, &templateForm)
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		err = app.FormDecoder.Decode(&templateForm, r.PostForm)
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		templateForm.CheckError(validator.NotBlank(templateForm.Title), "title", "Field can not be empty")
		templateForm.CheckError(validator.MaxChars(templateForm.Title, 100), "title", "Field can not be large then 100")
		templateForm.CheckError(validator.NotBlank(templateForm.Content), "content", "Field can not be empty")
		templateForm.CheckError(validator.MaxChars(templateForm.Content, 100), "title", "Field can not be large then 100")
		templateForm.CheckError(validator.PermittedInt(templateForm.Expires, 1, 7, 365), "expires", "This field must be 1, 7 or 365")

		if !templateForm.Valid() {
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

		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
	}
}

func LogIn(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func LogInPost(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func SignUp(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func SignUpPost(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func LogOut(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
