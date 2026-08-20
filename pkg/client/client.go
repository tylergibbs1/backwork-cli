package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type ErrorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Hint    string                 `json:"hint,omitempty"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"error"`
}

func New(apiKey, baseURL string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Request(method, path string, body interface{}, result interface{}) error {
	return c.RequestWithHeaders(method, path, body, result, nil)
}

func (c *Client) RequestWithHeaders(method, path string, body interface{}, result interface{}, headers map[string]string) error {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("User-Agent", "backwork-cli/1.0.0")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && !errResp.Success {
			return fmt.Errorf("%s: %s", errResp.Error.Code, errResp.Error.Message)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(path string, result interface{}) error {
	return c.Request("GET", path, nil, result)
}

func (c *Client) Post(path string, body interface{}, result interface{}) error {
	return c.Request("POST", path, body, result)
}

func (c *Client) PostWithHeaders(path string, body interface{}, result interface{}, headers map[string]string) error {
	return c.RequestWithHeaders("POST", path, body, result, headers)
}
