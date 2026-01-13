package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Home(errLog *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		}

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/pages/home.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			errLog.Fatal(err)
			http.Error(w, "Interanl Server Error", http.StatusInternalServerError)
			return
		}

		if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
			errLog.Fatal(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func ShowSnippets(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "No method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Creating a snippet...")
}
