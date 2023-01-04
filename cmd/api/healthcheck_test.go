package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, header, body := ts.get(t, "/v1/healthcheck")

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "application/json", header["Content-Type"][0])
	expected_body := createExpectedBodyResponse(t, http.StatusOK, nil)

	assert.JSONEq(t, expected_body, body)
}
