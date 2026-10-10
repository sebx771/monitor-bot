package adapters

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// quick test for the case where the file is not found
func TestUploadState_FileNotFound(t *testing.T) {
	client := NewGitHubGistClient("fake-token", "fake-gist-id")

	err := client.UploadState(context.Background(), "non-existent-file.json")

	if err == nil {
		t.Fatal("expected an error when the file does not exist")
	}
}

// quick test to verify the success of the function
func TestUploadState_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewGitHubGistClient("fake-token", "fake-gist-id")
	client.baseURL = server.URL

	dir := t.TempDir()

	filePath := filepath.Join(dir, "state.json")
	err := os.WriteFile(filePath, []byte(`{"test":"hello"}`), 0644)

	if err != nil {
		t.Fatal(err)
	}
	err = client.UploadState(context.Background(), filePath)

	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

// quick test to verify the failure of the function with error 500
func TestUploadState_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewGitHubGistClient("fake-token", "fake-gist-id")
	client.baseURL = server.URL

	dir := t.TempDir()

	filePath := filepath.Join(dir, "state.json")
	err := os.WriteFile(filePath, []byte(`{"test":"hello"}`), 0644)

	if err != nil {
		t.Fatal(err)
	}
	err = client.UploadState(context.Background(), filePath)

	if err == nil {
		t.Fatalf("expected an error, got: %v", err)
	}
}

// quick test to verify the success of the download and the saved content
func TestDownloadState_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fake-token" {
			t.Errorf("expected the Authorization header with the token")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"files":{"state.json":{"content":"{\"test\":\"hello\"}"}}}`))
	}))
	defer server.Close()

	client := NewGitHubGistClient("fake-token", "fake-gist-id")
	client.baseURL = server.URL

	dir := t.TempDir()
	filePath := filepath.Join(dir, "state.json")

	err := client.DownloadState(context.Background(), filePath)
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"test":"hello"}` {
		t.Fatalf("unexpected content: %s", string(data))
	}
}

// quick test to verify the error when the gist does not contain the file
func TestDownloadState_FileNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"files":{"otro.json":{"content":"x"}}}`))
	}))
	defer server.Close()

	client := NewGitHubGistClient("fake-token", "fake-gist-id")
	client.baseURL = server.URL

	dir := t.TempDir()
	filePath := filepath.Join(dir, "state.json")

	err := client.DownloadState(context.Background(), filePath)
	if err == nil {
		t.Fatal("expected an error when the gist does not contain the file")
	}
}

// quick test to verify the failure of the download with error 500
func TestDownloadState_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewGitHubGistClient("fake-token", "fake-gist-id")
	client.baseURL = server.URL

	dir := t.TempDir()
	filePath := filepath.Join(dir, "state.json")

	err := client.DownloadState(context.Background(), filePath)
	if err == nil {
		t.Fatal("expected an error with status 500")
	}
}
