package aiven

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	token      string
	httpClient *http.Client
}

type Service struct {
	Name  string `json:"service_name"`
	Type  string `json:"service_type"`
	State string `json:"state"`
}

type ServicesResponse struct {
	Services []Service `json:"services"`
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetServices(project string) ([]Service, error) {
	url := fmt.Sprintf(
		"https://api.aiven.io/v1/project/%s/service",
		project,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "aivenv1 "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"Aiven API returned status: %s",
			res.Status,
		)
	}

	var response ServicesResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Services, nil
}

func (c *Client) StartService(project string, service Service) error {
	url := fmt.Sprintf(
		"https://api.aiven.io/v1/project/%s/service/%s",
		project,
		service.Name,
	)

	body := []byte(`{"powered":true}`)

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "aivenv1 "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf(
			"Aiven API returned status: %s",
			res.Status,
		)
	}

	return nil
}