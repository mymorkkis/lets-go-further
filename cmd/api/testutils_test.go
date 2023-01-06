package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	testConfig := config{
		port: 9999,
		env:  "testing",
	}

	return &application{
		version: version,
		logger:  log.New(io.Discard, "", 0),
		config:  &testConfig,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
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

func createExpectedBodyResponse(t *testing.T, code int, data any) string {
	response := make(map[string]any)
	response["status"] = Status{Code: code, Message: http.StatusText(code)}
	response["systemInfo"] = SystemInfo{Environment: "testing", Version: version}
	response["data"] = data

	json_, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	return string(json_)
}
