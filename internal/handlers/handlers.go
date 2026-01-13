package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/triplq/snippetbox/internal/helpers"
)

func Home(errLog *log.Logger) http.HandlerFunc {
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
			helpers.ServerError(w, err, errLog)
			return
		}

		if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
			helpers.ServerError(w, err, errLog)
		}
	}
}

func ShowSnippets(errLog *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			helpers.NotFound(w)
			return
		}

		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func CreateSnippet(errLog *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			helpers.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		fmt.Fprint(w, "Creating a snippet...")
	}
}
