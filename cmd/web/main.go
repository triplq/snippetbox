package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", showSnippets)
	mux.HandleFunc("/snippets/create", createSnippet)

	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
