package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTavilyMockServer(t *testing.T, expectedResponse interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		json.NewEncoder(w).Encode(expectedResponse)
	}))
}

func TestTavilyProvider_Search(t *testing.T) {
	mockServer := NewTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [
				{
					"url": "https://example.com",
					"content": "Example content"
				}
			]
		}`))
	}))
	defer mockServer.Close()

	// Set test environment
	oldKey := os.Getenv("TAVILY_API_KEY")
	os.Setenv("TAVILY_API_KEY", "test-key")
	defer os.Setenv("TAVILY_API_KEY", oldKey)

	provider, err := NewTavilyProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	provider.baseURL = mockServer.URL // Set mock server URL

	results, err := provider.Search("test query", map[string]string{"depth": "advanced"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Verify results
	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}

	expected := SearchResult{
		URL:     "https://example.com",
		Content: "Example content",
	}

	if results.Results[0] != expected {
		t.Errorf("Expected %+v, got %+v", expected, results.Results[0])
	}
}

func TestTavilyProvider_Extract(t *testing.T) {
	mockServer := NewTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [
				{
					"url": "https://example.com",
					"raw_content": "Example content"
				}
			]
		}`))
	}))
	defer mockServer.Close()

	// Set test environment
	oldKey := os.Getenv("TAVILY_API_KEY")
	os.Setenv("TAVILY_API_KEY", "test-key")
	defer os.Setenv("TAVILY_API_KEY", oldKey)

	provider, err := NewTavilyProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	provider.baseURL = mockServer.URL // Set mock server URL

	results, err := provider.Extract([]string{"https://example.com"})
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify results
	if len(results.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Results))
	}

	expected := ExtractResult{
		URL:        "https://example.com",
		RawContent: "Example content",
	}

	if results.Results[0] != expected {
		t.Errorf("Expected %+v, got %+v", expected, results.Results[0])
	}
}

func TestTavilyProvider_NoAPIKey(t *testing.T) {
	// Clear API key
	oldKey := os.Getenv("TAVILY_API_KEY")
	os.Unsetenv("TAVILY_API_KEY")
	defer os.Setenv("TAVILY_API_KEY", oldKey)

	_, err := NewTavilyProvider()
	if err == nil {
		t.Error("Expected error when no API key is set")
	}
}
