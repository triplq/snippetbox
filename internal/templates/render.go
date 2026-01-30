package templates

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func Canvas(w http.ResponseWriter, status int, data *TemplateData, cache map[string]*template.Template, page string) error {
	ts, ok := cache[page]
	if !ok {
		err_string := fmt.Sprintf("this template: %s is not exist", page)
		err := errors.New(err_string)
		return err
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	buf.WriteTo(w)

	return nil
}
