package testutils

import (
	"bytes"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/triplq/snippetbox/internal/application"
	"github.com/triplq/snippetbox/internal/models/mocks"
	"github.com/triplq/snippetbox/internal/templates"
)

var csrfTokenRX = regexp.MustCompile(`name=["']csrf_token["'][^>]*value=["']([^"']+)["']`)

func ExtractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

func (ts *TestServer) PostForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	req, err := http.NewRequest(
		http.MethodPost,
		ts.URL+urlPath,
		strings.NewReader(string(form.Encode())),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ts.URL)

	rs, err := ts.Client().Do(req)
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
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &TestServer{ts}
}

func (ts *TestServer) Get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
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
