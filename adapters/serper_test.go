package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSerperProvider_Search(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Errorf("Expected API key header, got %s", r.Header.Get("X-API-KEY"))
		}

		// Return mock response
		response := map[string]interface{}{
			"organic": []map[string]interface{}{
				{
					"link":    "https://example.com",
					"snippet": "Example content",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Set test environment
	os.Setenv("SERPER_API_KEY", "test-key")
	defer os.Unsetenv("SERPER_API_KEY")

	provider, err := NewSerperProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Override API URL for testing
	originalURL := "https://google.serper.dev/search"
	http.DefaultClient = server.Client()

	results, err := provider.Search("test query", nil)
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

	// Restore original URL
	_ = originalURL
}
