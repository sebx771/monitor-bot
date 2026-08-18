package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	requestTimeout = 30 * time.Second
	gistAPIBaseURL = "https://api.github.com"
)

type gistFile struct {
	Content string `json:"content"`
}

type gistResponse struct {
	Files map[string]gistFile `json:"files"`
}

// GitHubGistClient implementa port.StateStorage usando la API de GitHub
// para almacenar el estado de sesión en un Gist.
type GitHubGistClient struct {
	token   string
	gistID  string
	baseURL string

	httpClient *http.Client
}

func NewGitHubGistClient(token, gistID string) *GitHubGistClient {
	return &GitHubGistClient{
		token:   token,
		gistID:  gistID,
		baseURL: gistAPIBaseURL,
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// DownloadState descarga el archivo de estado desde el Gist y lo guarda en
// destinationPath, creando los directorios necesarios si la ruta lo requiere.
func (g *GitHubGistClient) DownloadState(ctx context.Context, destinationPath string) error {
	url := fmt.Sprintf("%s/gists/%s", g.baseURL, g.gistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creando request al Gist: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("descargando estado desde el Gist: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API devolvió status: %s", res.Status)
	}

	var gist gistResponse
	if err := json.NewDecoder(res.Body).Decode(&gist); err != nil {
		return fmt.Errorf("decodificando respuesta del Gist: %w", err)
	}

	fileName := filepath.Base(destinationPath)
	file, ok := gist.Files[fileName]
	if !ok {
		return fmt.Errorf("el Gist no contiene el archivo %s", fileName)
	}

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return fmt.Errorf("creando el directorio de storage: %w", err)
	}

	if err := os.WriteFile(destinationPath, []byte(file.Content), 0o644); err != nil {
		return fmt.Errorf("escribiendo el archivo de estado: %w", err)
	}

	return nil
}

// UploadState todavía no está implementado.
type gistUpdateFile struct {
	Content string `json:"content"`
}

type gistUpdateRequest struct {
	Files map[string]gistUpdateFile `json:"files"`
}

func (g *GitHubGistClient) UploadState(ctx context.Context, sourcePath string) error {
	file, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("leyendo archivo de estado: %w", err)
	}

	fileName := filepath.Base(sourcePath)

	payload := gistUpdateRequest{
		Files: map[string]gistUpdateFile{
			fileName: {
				Content: string(file),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("serializando estado del Gist: %w", err)
	}

	url := fmt.Sprintf("%s/gists/%s", g.baseURL, g.gistID)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("creando request al Gist: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("actualizando estado en el Gist: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API devolvió status: %s", res.Status)
	}

	return nil
}
