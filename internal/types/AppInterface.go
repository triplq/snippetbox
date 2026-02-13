package types

import "net/http"

type AppInterface interface {
	IsAuthenticated(r *http.Request) bool
	GetFlash(r *http.Request) string
}
