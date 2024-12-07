package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type TavilyProvider struct {
	apiKey string
}

func NewTavilyProvider() (*TavilyProvider, error) {
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEY environment variable is not set")
	}
	return &TavilyProvider{apiKey: apiKey}, nil
}

func (t *TavilyProvider) Extract(urls []string) (*ExtractResponse, error) {
	params := struct {
		URLs   []string `json:"urls"`
		ApiKey string   `json:"api_key"`
	}{
		URLs:   urls,
		ApiKey: t.apiKey,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.tavily.com/extract", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

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

	var result ExtractResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (t *TavilyProvider) Search(query string, options map[string]string) (*SearchResponse, error) {
	params := struct {
		Query          string   `json:"query"`
		ApiKey         string   `json:"api_key"`
		SearchDepth    string   `json:"search_depth"`
		IncludeDomains []string `json:"include_domains"`
		ExcludeDomains []string `json:"exclude_domains"`
	}{
		Query:          query,
		ApiKey:         t.apiKey,
		SearchDepth:    "basic",
		IncludeDomains: []string{},
		ExcludeDomains: []string{},
	}

	// Apply options if provided
	if depth, ok := options["depth"]; ok {
		params.SearchDepth = depth
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.tavily.com/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

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

	var result struct {
		Results []struct {
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Convert to common SearchResponse format
	searchResults := make([]SearchResult, 0, len(result.Results))
	for _, r := range result.Results {
		searchResults = append(searchResults, SearchResult{
			URL:     r.URL,
			Content: r.Content,
		})
	}

	return &SearchResponse{Results: searchResults}, nil
}
