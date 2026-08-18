package adapters

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// test rapido en caso que el archivo no se encuentre
func TestUploadState_FileNotFound(t *testing.T) {
	client := NewGitHubGistClient("fake-token", "fake-gist-id")

	err := client.UploadState(context.Background(), "archivo-que-no-existe.json")

	if err == nil {
		t.Fatal("se esperaba un error cuando el archivo no existe")
	}
}

// test rapido para verificar el exito de la funcion
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
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}
}

// test rapido para verificar el fracaso de la funcion con error 500
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
		t.Fatalf("se esperaba un error, se obtuvo: %v", err)
	}
}

// test rapido para verificar el exito de la descarga y el contenido guardado
func TestDownloadState_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fake-token" {
			t.Errorf("se esperaba el header Authorization con el token")
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
		t.Fatalf("se esperaba nil, se obtuvo: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"test":"hello"}` {
		t.Fatalf("contenido inesperado: %s", string(data))
	}
}

// test rapido para verificar el error cuando el gist no contiene el archivo
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
		t.Fatal("se esperaba un error cuando el gist no contiene el archivo")
	}
}

// test rapido para verificar el fracaso de la descarga con error 500
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
		t.Fatal("se esperaba un error con status 500")
	}
}
