package supabase

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBaseURL = "https://api.supabase.com"

type Client struct {
	token      string
	httpClient *http.Client
	baseUrl string
}

type Project struct {
	ID     string `json:"id"`
	Ref    string `json:"ref"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func NewClient(token string) *Client {
	
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		baseUrl: apiBaseURL,
	}
}

func (c *Client) ListProjects() ([]Project, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/v1/projects", c.baseUrl),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"supabase API returned status %d",
			resp.StatusCode,
		)
	}

	var projects []Project

	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) GetProject(ref string) (*Project, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/v1/projects/%s", c.baseUrl, ref),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"supabase API returned status %d for project %s",
			resp.StatusCode,
			ref,
		)
	}

	var project Project

	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) RestoreProject(ref string) error {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/projects/%s/restore", c.baseUrl, ref),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf(
			"supabase API returned status %d while restoring project %s",
			resp.StatusCode,
			ref,
		)
	}

	return nil
}