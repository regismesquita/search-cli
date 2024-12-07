package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupMockServer(t *testing.T, expectedResponse interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Errorf("Expected API key header, got %s", r.Header.Get("X-API-KEY"))
		}

		json.NewEncoder(w).Encode(expectedResponse)
	}))
}

func TestSerperProvider_Search(t *testing.T) {
	mockServer := NewTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Errorf("Expected API key header, got %s", r.Header.Get("X-API-KEY"))
		}

		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"organic": [
				{
					"link": "https://example.com",
					"snippet": "Example content"
				}
			]
		}`))
	}))
	defer mockServer.Close()

	// Set test environment
	oldKey := os.Getenv("SERPER_API_KEY")
	os.Setenv("SERPER_API_KEY", "test-key")
	defer os.Setenv("SERPER_API_KEY", oldKey)

	provider, err := NewSerperProvider()
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	provider.baseURL = mockServer.URL

	results, err := provider.Search("test query", nil)
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

func TestSerperProvider_NoAPIKey(t *testing.T) {
	// Clear API key
	oldKey := os.Getenv("SERPER_API_KEY")
	os.Unsetenv("SERPER_API_KEY")
	defer os.Setenv("SERPER_API_KEY", oldKey)

	_, err := NewSerperProvider()
	if err == nil {
		t.Error("Expected error when no API key is set")
	}
}
