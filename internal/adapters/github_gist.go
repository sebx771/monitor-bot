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

// GitHubGistClient implements port.StateStorage using the GitHub API to
// store the session state in a Gist.
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

// DownloadState downloads the state file from the Gist and saves it to the
// destinationPath, creating the necessary directories if the path requires it.
func (g *GitHubGistClient) DownloadState(ctx context.Context, destinationPath string) error {
	url := fmt.Sprintf("%s/gists/%s", g.baseURL, g.gistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request to the Gist: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading state from the Gist: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status: %s", res.Status)
	}

	var gist gistResponse
	if err := json.NewDecoder(res.Body).Decode(&gist); err != nil {
		return fmt.Errorf("decoding the Gist response: %w", err)
	}

	fileName := filepath.Base(destinationPath)
	file, ok := gist.Files[fileName]
	if !ok {
		return fmt.Errorf("the Gist does not contain the file %s", fileName)
	}

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return fmt.Errorf("creating the storage directory: %w", err)
	}

	if err := os.WriteFile(destinationPath, []byte(file.Content), 0o644); err != nil {
		return fmt.Errorf("writing the state file: %w", err)
	}

	return nil
}

// UploadState is not implemented yet.
type gistUpdateFile struct {
	Content string `json:"content"`
}

type gistUpdateRequest struct {
	Files map[string]gistUpdateFile `json:"files"`
}

func (g *GitHubGistClient) UploadState(ctx context.Context, sourcePath string) error {
	file, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("reading the state file: %w", err)
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
		return fmt.Errorf("serializing the Gist state: %w", err)
	}

	url := fmt.Sprintf("%s/gists/%s", g.baseURL, g.gistID)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("creating request to the Gist: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("updating the state in the Gist: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status: %s", res.Status)
	}

	return nil
}
