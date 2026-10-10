package aiven

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func newTestClient(token, baseURL string) *Client {
	return &Client{
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func TestCheckStartsPoweredOffServices(t *testing.T) {
	var puts atomic.Int32
	var authHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			puts.Add(1)
			authHeader = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			return
		}

		json.NewEncoder(w).Encode(ServicesResponse{
			Services: []Service{
				{Name: "svc-on", Type: "pg", State: "RUNNING"},
				{Name: "svc-off", Type: "pg", State: "POWEROFF"},
			},
		})
	}))
	defer server.Close()

	checker := NewChecker(newTestClient("test-token", server.URL), "proj-1")

	if err := checker.Check(); err != nil {
		t.Fatalf("Check returned unexpected error: %v", err)
	}

	if puts.Load() != 1 {
		t.Fatalf("expected 1 PUT, got %d", puts.Load())
	}

	if authHeader != "aivenv1 test-token" {
		t.Fatalf("incorrect auth header: %q", authHeader)
	}
}

func TestCheckDoesNotAbortOnServiceFailure(t *testing.T) {
	var puts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			puts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ServicesResponse{
			Services: []Service{
				{Name: "svc-1", Type: "pg", State: "POWEROFF"},
				{Name: "svc-2", Type: "pg", State: "POWEROFF"},
			},
		})
	}))
	defer server.Close()

	checker := NewChecker(newTestClient("test-token", server.URL), "proj-1")

	if err := checker.Check(); err != nil {
		t.Fatalf("Check returned unexpected error: %v", err)
	}

	if puts.Load() != 2 {
		t.Fatalf("expected 2 PUTs despite the failure, got %d", puts.Load())
	}
}

func TestCheckAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	checker := NewChecker(newTestClient("bad-token", server.URL), "proj-1")

	err := checker.Check()
	if err == nil {
		t.Fatal("expected error when the API fails")
	}

	if !strings.Contains(err.Error(), "proj-1") {
		t.Fatalf("the error should mention the project, got: %v", err)
	}
}
