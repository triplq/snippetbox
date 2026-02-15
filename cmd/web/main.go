package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/models"
	"github.com/triplq/snippetbox/internal/server"
	"github.com/triplq/snippetbox/internal/templates"
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

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		slog.Error("critical error", "err", err)
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application.Application{
		Snippets:       &models.SnippetModel{DB: db},
		Users:          &models.UserModel{DB: db},
		TemplateCache:  templateCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:         *address,
		Handler:      server.Routes(app),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Listening on adress", "addr", *address)
	if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
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
