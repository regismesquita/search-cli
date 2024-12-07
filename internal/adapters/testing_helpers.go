package adapters

import (
	"net/http"
	"net/http/httptest"
)

// TestServer represents a mock server for testing
type TestServer struct {
	*httptest.Server
	OriginalClient *http.Client
}

// NewTestServer creates a new test server and patches the default HTTP client
func NewTestServer(handler http.HandlerFunc) *TestServer {
	server := httptest.NewServer(handler)
	originalClient := http.DefaultClient
	http.DefaultClient = server.Client()

	return &TestServer{
		Server:         server,
		OriginalClient: originalClient,
	}
}

// Close closes the server and restores the original HTTP client
func (s *TestServer) Close() {
	http.DefaultClient = s.OriginalClient
	s.Server.Close()
}
