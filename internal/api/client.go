package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	BaseURL        = "https://api.modrinth.com/v2"
	UserAgent      = "Mint/0.1.0 (terminal-modrinth-client)"
	DefaultTimeout = 30 * time.Second
	MaxRetries     = 3
	RetryDelay     = time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	apiKey     string
}

type ClientOption func(*Client)

func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:   BaseURL,
		userAgent: UserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) doRequest(method, path string, query url.Values, body io.Reader) (*http.Response, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", c.apiKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) get(path string, query url.Values, result interface{}) error {
	resp, err := c.doRequest(http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("not found: %s", string(body))
	case http.StatusTooManyRequests:
		retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return &RateLimitError{RetryAfter: retryAfter}
	default:
		var apiErr struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
			return &APIError{Status: resp.StatusCode, Message: apiErr.Error}
		}
		return &APIError{Status: resp.StatusCode, Message: string(body)}
	}
}

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.Status, e.Message)
}

type RateLimitError struct {
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited, retry after %ds", e.RetryAfter)
}
