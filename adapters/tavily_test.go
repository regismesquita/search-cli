package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTavilyProvider_Search(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify request body contains API key
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if reqBody["api_key"] != "test-key" {
			t.Errorf("Expected API key in body, got %v", reqBody["api_key"])
		}

		// Return mock response
		response := map[string]interface{}{
			"results": []map[string]interface{}{
				{
					"url":     "https://example.com",
					"content": "Example content",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Set test environment
	os.Setenv("TAVILY_API_KEY", "test-key")
	defer os.Unsetenv("TAVILY_API_KEY")

	provider, err := NewTavilyProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Override API URL for testing
	http.DefaultClient = server.Client()

	results, err := provider.Search("test query", map[string]string{"depth": "advanced"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}

	if results.Results[0].URL != "https://example.com" {
		t.Errorf("Expected URL https://example.com, got %s", results.Results[0].URL)
	}

	if results.Results[0].Content != "Example content" {
		t.Errorf("Expected content 'Example content', got %s", results.Results[0].Content)
	}
}

func TestTavilyProvider_Extract(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Return mock response
		response := ExtractResponse{
			Results: []ExtractResult{
				{
					URL:        "https://example.com",
					RawContent: "Example content",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Set test environment
	os.Setenv("TAVILY_API_KEY", "test-key")
	defer os.Unsetenv("TAVILY_API_KEY")

	provider, err := NewTavilyProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Override API URL for testing
	http.DefaultClient = server.Client()

	results, err := provider.Extract([]string{"https://example.com"})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}

	if results.Results[0].URL != "https://example.com" {
		t.Errorf("Expected URL https://example.com, got %s", results.Results[0].URL)
	}
}
