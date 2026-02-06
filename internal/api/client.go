package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.infomaniak.com"

// Client communicates with the Infomaniak API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// ClientConfig holds configuration for creating a Client.
type ClientConfig struct {
	Token   string
	BaseURL string
}

// NewClient creates a new Infomaniak API client.
func NewClient(cfg ClientConfig) *Client {
	base := cfg.BaseURL
	if base == "" {
		base = defaultBaseURL
	}
	return &Client{
		baseURL: base,
		token:   cfg.Token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request %s %s: %w", method, path, err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request %s %s: %w", method, path, err)
	}

	return resp, nil
}

func decodeResponse[T any](resp *http.Response) (*Response[T], error) {
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var result Response[T]
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response (status %d): %w", resp.StatusCode, err)
	}

	if result.Result == "error" {
		if result.Error != nil {
			return nil, fmt.Errorf("api error %s: %s", result.Error.Code, result.Error.Description)
		}
		return nil, fmt.Errorf("api error (status %d): unknown error", resp.StatusCode)
	}

	return &result, nil
}
