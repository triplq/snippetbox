package helpers

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func ServerError(w http.ResponseWriter, err error, errLog *log.Logger) {
	trace := fmt.Sprintf("%s/n%s", err, debug.Stack())
	errLog.Print(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func NotFound(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}
