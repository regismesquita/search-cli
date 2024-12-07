package formatter

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/regismesquita/search-cli/internal/adapters"
)

func TestFormatResults(t *testing.T) {
	results := &adapters.SearchResponse{
		Results: []adapters.SearchResult{
			{
				URL:     "https://example.com",
				Content: "Test content",
			},
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	FormatResults(results)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "https://example.com") {
		t.Error("Expected output to contain URL")
	}
	if !strings.Contains(output, "Test content") {
		t.Error("Expected output to contain content")
	}
}

func TestFormatExtractResults(t *testing.T) {
	results := &adapters.ExtractResponse{
		Results: []adapters.ExtractResult{
			{
				URL:        "https://example.com",
				RawContent: "Test content",
			},
		},
		FailedResults: []adapters.FailedResult{
			{
				URL:   "https://failed.com",
				Error: "Failed to extract",
			},
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	FormatExtractResults(results)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "https://example.com") {
		t.Error("Expected output to contain URL")
	}
	if !strings.Contains(output, "Test content") {
		t.Error("Expected output to contain content")
	}
	if !strings.Contains(output, "Failed URLs") {
		t.Error("Expected output to contain failed URLs section")
	}
	if !strings.Contains(output, "https://failed.com") {
		t.Error("Expected output to contain failed URL")
	}
}
