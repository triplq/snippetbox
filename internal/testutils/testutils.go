package testutils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/models/mocks"
	"github.com/triplq/snippetbox/internal/templates"
)

func NewTestApp(t *testing.T) *application.Application {
	sm := scs.New()
	sm.Lifetime = time.Hour * 12
	sm.Cookie.Secure = true

	tc, err := templates.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	return &application.Application{
		Snippets:       &mocks.SnippetModel{},
		Users:          &mocks.UserModel{},
		SessionManager: sm,
		FormDecoder:    formDecoder,
		TemplateCache:  tc,
	}
}

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewTLSServer(h)
	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
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
