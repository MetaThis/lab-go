package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestHandlers is an example of how we would "inject" a different set of dependencies
// for tests. Here we're substituting an in-memory instance (testDB) of sqlite for tests.
func NewTestHandlers() Handlers {
	schema := NewSchema()
	h := NewHandlers(testDB, schema)
	return h
}

// Boilerplate for test requests
func SubmitRequest(json []byte, expectedStatus int, t *testing.T) string {
	req, err := http.NewRequest("POST", "/lab/instrument/1/samples", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	h := NewTestHandlers()
	router := NewRouter(h)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != expectedStatus {
		t.Errorf("Unexpected status code: got %v want %v",
			resp.Code, expectedStatus)
	}

	return resp.Body.String()
}

func TestHappyPath(t *testing.T) {
	// Note that in real life scenarios with more sophisticated JSON input, we can
	// use the Go convention of a testdata directory with .json files for multiple
	// test cases, then add a helper function to load them as needed.
	json := []byte(`[{"id":1},{"id":2},{"id":999}]`)
	resBody := SubmitRequest(json, 200, t)

	// Verify response body
	expected := `{"runId":1}`
	if resBody != expected {
		t.Errorf("Unexpected body: got %v want %v", resBody, expected)
	}
}

func TestPostBadData(t *testing.T) {
	// Use "table driven" tests to run a variety of bad data scenarios.
	payloads := []string{
		`[{"foo":"bar"}]`,
		`{"foo":"bar"}`,
		`"foo"`,
		`[]`,
		`[{"id":1},{"id":2},{"id":3},{"id":4},{"id":5},{"id":6},{"id":7},{"id":8},{"id":9},{"id":10},{"id":11}]`,
	}

	for _, payload := range payloads {
		SubmitRequest([]byte(payload), 400, t)
	}
}

func TestPostNonNumericInstrumentID(t *testing.T) {
	json := []byte(`[{"id":1},{"id":2},{"id":999}]`)

	req, err := http.NewRequest("POST", "/lab/instrument/abc/samples", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	h := NewTestHandlers()
	router := NewRouter(h)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 400 {
		t.Errorf("Unexpected status code: got %v want %v", resp.Code, 400)
	}
}
