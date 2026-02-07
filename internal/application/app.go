package application

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/triplq/snippetbox/internal/models"
	"github.com/triplq/snippetbox/internal/templates"
)

type Application struct {
	Snippets       *models.SnippetModel
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *templates.TemplateData) error {
	err := templates.Canvas(w, status, data, app.TemplateCache, page)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoder *form.InvalidDecoderError
		if errors.As(err, &invalidDecoder) {
			panic(err)
		}

		return err
	}

	return nil
}
