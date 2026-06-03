package api

import (
	"testing"
)

func TestClientDefaults(t *testing.T) {
	c := NewClient()
	if c.baseURL != BaseURL {
		t.Errorf("expected baseURL %s, got %s", BaseURL, c.baseURL)
	}
	if c.userAgent != UserAgent {
		t.Errorf("expected userAgent %s, got %s", UserAgent, c.userAgent)
	}
	if c.httpClient == nil {
		t.Error("expected non-nil http client")
	}
}

func TestClientOptions(t *testing.T) {
	c := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL("https://example.com"),
	)
	if c.apiKey != "test-key" {
		t.Errorf("expected apiKey test-key, got %s", c.apiKey)
	}
	if c.baseURL != "https://example.com" {
		t.Errorf("expected baseURL https://example.com, got %s", c.baseURL)
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{Status: 404, Message: "not found"}
	if err.Error() != "API error (404): not found" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestRateLimitError(t *testing.T) {
	err := &RateLimitError{RetryAfter: 60}
	if err.Error() != "rate limited, retry after 60s" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
