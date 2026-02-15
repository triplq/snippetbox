package testutils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/triplq/snippetbox/internal/application"
)

func NewTestApp(t *testing.T) *application.Application {
	sm := scs.New()
	sm.Lifetime = time.Hour * 12

	return &application.Application{
		SessionManager: sm,
	}
}

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)
	return &TestServer{ts}
}

func (ts *TestServer) Get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
