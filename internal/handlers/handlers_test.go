package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/assert"
	"github.com/triplq/snippetbox/internal/server"
)

func TestPing(t *testing.T) {
	sm := scs.New()
	sm.Lifetime = 12 * time.Hour

	app := &application.Application{
		SessionManager: sm,
	}

	ts := httptest.NewTLSServer(server.Routes(app))
	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
