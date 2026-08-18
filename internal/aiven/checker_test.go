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

func TestCheckIniciaServiciosApagados(t *testing.T) {
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
		t.Fatalf("Check devolvió error inesperado: %v", err)
	}

	if puts.Load() != 1 {
		t.Fatalf("se esperaba 1 PUT, se obtuvieron %d", puts.Load())
	}

	if authHeader != "aivenv1 test-token" {
		t.Fatalf("header de auth incorrecto: %q", authHeader)
	}
}

func TestCheckNoAbortaAnteFalloDeServicio(t *testing.T) {
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
		t.Fatalf("Check devolvió error inesperado: %v", err)
	}

	if puts.Load() != 2 {
		t.Fatalf("se esperaban 2 PUT a pesar del fallo, se obtuvieron %d", puts.Load())
	}
}

func TestCheckErrorDeAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	checker := NewChecker(newTestClient("bad-token", server.URL), "proj-1")

	err := checker.Check()
	if err == nil {
		t.Fatal("se esperaba error al fallar la API")
	}

	if !strings.Contains(err.Error(), "proj-1") {
		t.Fatalf("el error debería mencionar el proyecto, se obtuvo: %v", err)
	}
}
