package formatter

import (
	"fmt"
	"strings"

	"github.com/regismesquita/search-cli/internal/adapters"
)

func FormatResults(results *adapters.SearchResponse) {
	for _, result := range results.Results {
		fmt.Printf("\n🔗 %s\n", result.URL)
		fmt.Printf("📝 %s\n", result.Content)
		fmt.Println(strings.Repeat("-", 80))
	}
}

func FormatExtractResults(results *adapters.ExtractResponse) {
	for _, result := range results.Results {
		fmt.Printf("\n🔗 %s\n", result.URL)
		fmt.Printf("📄 %s\n", result.RawContent)
		fmt.Println(strings.Repeat("-", 80))
	}

	if len(results.FailedResults) > 0 {
		fmt.Println("\n❌ Failed URLs:")
		for _, failed := range results.FailedResults {
			fmt.Printf("URL: %s\nError: %s\n", failed.URL, failed.Error)
			fmt.Println(strings.Repeat("-", 40))
		}
	}
}
