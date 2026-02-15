package handlers_test

import (
	"net/http"
	"testing"

	"github.com/triplq/snippetbox/internal/assert"
	"github.com/triplq/snippetbox/internal/server"
	"github.com/triplq/snippetbox/internal/testutils"
)

func TestPing(t *testing.T) {
	app := testutils.NewTestApp(t)

	ts := testutils.NewTestServer(t, server.Routes(app))
	defer ts.Close()

	code, _, body := ts.Get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
