package supabase

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/sebx771/monitor-bot/internal/logger"
)

func newTestClient(token, baseURL string) *Client {

	return &Client{
		token:      token,
		baseUrl:    baseURL,
		httpClient: &http.Client{},
	}
}

func TestCheck_Checker_ListProjectError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {

			res.WriteHeader(http.StatusInternalServerError)
			return
		}

	}))
	defer server.Close()
	checker := InitChecker(server.URL)
	err := checker.Check()

	if err == nil {
		t.Fatalf("want error , got nil")
	}

}

func TestCheck_Checker_ActiveProjects(t *testing.T) {
	var pots atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode([]Project{
				{ID: "1", Ref: "123", Name: "mysql", Status: "ACTIVE"},
			})
			return
		}

		if r.Method == http.MethodPost {
			pots.Add(1)
			w.WriteHeader(http.StatusAccepted)
		}
	}))
	defer server.Close()
	checker := InitChecker(server.URL)
	checker.Check()

	if pots.Load() != 0 {
		t.Fatalf("expected 0 POST, got %d", pots.Load())
	}

}

func TestCheck_Checker_InactiveProjects(t *testing.T) {
	var pots atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode([]Project{
				{ID: "2", Ref: "1234", Name: "mysql", Status: "INACTIVE"},
			})
			return
		}

		if r.Method == http.MethodPost {
			pots.Add(1)
			w.WriteHeader(http.StatusAccepted)
		}
	}))
	defer server.Close()
	checker := InitChecker(server.URL)
	checker.Check()

	if pots.Load() != 1 {
		t.Fatalf("expected 1 PUT, got %d", pots.Load())
	}

}

func TestCheck_Checker_RestoreProjectError(t *testing.T) {
	var posts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode([]Project{
				{
					ID:     "2",
					Ref:    "1234",
					Name:   "mysql",
					Status: "INACTIVE",
				},
			})
			return
		}

		if r.Method == http.MethodPost {
			posts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))

	defer server.Close()

	checker := InitChecker(server.URL)

	err := checker.Check()

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if posts.Load() != 1 {
		t.Fatalf("expected 1 POST, got %d", posts.Load())
	}
}

func InitChecker(URL string) Checker {
	logg := logger.NewLogger("TEST")
	client := newTestClient("tk", URL)
	checker := NewChecker(client, logg)
	return *checker
}
