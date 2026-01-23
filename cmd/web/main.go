package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/models"
)

func main() {
	address := flag.String("addr", ":4000", "Address to use")
	dsn := flag.String("dsn", "web_app:@tcp(127.0.0.1:3306)/snippetbox?parseTime=true", "MySQL connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := openDB(*dsn)
	if err != nil {
		slog.Error("critical error", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	app := &application.Application{
		Snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:    *address,
		Handler: routes(app),
	}

	slog.Info("Listening on adress", "addr", *address)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("critical error", "err", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
