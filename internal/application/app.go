package application

import (
	"html/template"
	"net/http"

	"github.com/triplq/snippetbox/internal/models"
	"github.com/triplq/snippetbox/internal/templates"
)

type Application struct {
	Snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *templates.TemplateData) error {
	err := templates.Canvas(w, status, data, app.TemplateCache, page)
	if err != nil {
		return err
	}
	return nil
}
