package helpers

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func ServerError(w http.ResponseWriter, err error) {
	trace := (debug.Stack())
	slog.Error("server error", "trace", trace)

	fmt.Printf("\n----DEbUG OUTPUT-----\n%s\n%s\n", err.Error(), trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}
