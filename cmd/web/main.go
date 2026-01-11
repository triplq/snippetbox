package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	address := flag.String("addr", ":4000", "Address to use")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippets/view", showSnippets)
	mux.HandleFunc("/snippets/create", createSnippet)

	srv := &http.Server{
		Addr:     *address,
		ErrorLog: errLog,
		Handler:  mux,
	}

	infoLog.Println("Listening on adress", *address)
	if err := srv.ListenAndServe(); err != nil {
		errLog.Fatal(err)
	}
}
