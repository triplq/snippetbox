package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snpbox"))
}

func snippetShow(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Showing snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating snippet..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", snippetShow)
	mux.HandleFunc("/snippets/create", snippetCreate)

	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
