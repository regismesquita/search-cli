package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type SerperProvider struct {
	apiKey  string
	baseURL string
}

type serperRequest struct {
	Query string `json:"q"`
}

type serperResponse struct {
	Organic []struct {
		Link        string `json:"link"`
		Description string `json:"snippet"`
	} `json:"organic"`
}

func NewSerperProvider() (*SerperProvider, error) {
	apiKey := os.Getenv("SERPER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SERPER_API_KEY environment variable is not set")
	}
	return &SerperProvider{
		apiKey:  apiKey,
		baseURL: "https://google.serper.dev/search",
	}, nil
}

func (s *SerperProvider) Search(query string, options map[string]string) (*SearchResponse, error) {
	payload := serperRequest{Query: query}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var serperResp serperResponse
	if err := json.NewDecoder(resp.Body).Decode(&serperResp); err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(serperResp.Organic))
	for _, result := range serperResp.Organic {
		results = append(results, SearchResult{
			URL:     result.Link,
			Content: result.Description,
		})
	}

	return &SearchResponse{Results: results}, nil
}
