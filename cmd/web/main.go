package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", showSnippets)
	mux.HandleFunc("/snippets/create", createSnippet)

	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
