package helpers

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/triplq/snippetbox/internal/application"
)

func ServerError(w http.ResponseWriter, err error, app *application.Application) {
	trace := (debug.Stack())
	slog.Error("server error", "trace", trace)

	fmt.Printf("\n----DEBUG OUTPUT-----\n%s\n%s\n", err.Error(), trace)

	if *app.Debug_mode == true {
		fmt.Fprintf(w, "\n----DEBUG OUTPUT-----\n%s\n%s\n", err.Error(), trace)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}
