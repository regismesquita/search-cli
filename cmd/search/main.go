package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"search-cli/internal/adapters"
	"search-cli/internal/formatter"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Output in JSON format")
	tavilyMode := flag.Bool("tavily", false, "Use Tavily as provider")
	serperMode := flag.Bool("serper", false, "Use Serper as provider (default)")
	extractMode := flag.Bool("extract", false, "Extract content from URLs (Tavily only)")
	depth := flag.String("depth", "basic", "Search depth (basic or advanced) - only for Tavily")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: Query or URLs required")
		os.Exit(1)
	}

	// Validate provider flags
	if *tavilyMode && *serperMode {
		fmt.Fprintln(os.Stderr, "Error: Cannot use both -tavily and -serper")
		os.Exit(1)
	}

	// Validate extract mode
	if *extractMode && *serperMode {
		fmt.Fprintln(os.Stderr, "Error: Extract mode is only available with Tavily")
		os.Exit(1)
	}

	if *extractMode {
		provider, err := adapters.NewTavilyProvider()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		urls := strings.Fields(strings.Join(flag.Args(), " "))
		results, err := provider.Extract(urls)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if *jsonOutput {
			json.NewEncoder(os.Stdout).Encode(results)
		} else {
			formatter.FormatExtractResults(results)
		}
	} else {
		var searchProvider adapters.SearchProvider
		var err error

		// Default to Serper unless Tavily is explicitly requested
		if *tavilyMode {
			searchProvider, err = adapters.NewTavilyProvider()
		} else {
			searchProvider, err = adapters.NewSerperProvider()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		options := map[string]string{}
		if *tavilyMode {
			options["depth"] = *depth
		}

		query := flag.Arg(0)
		results, err := searchProvider.Search(query, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if *jsonOutput {
			json.NewEncoder(os.Stdout).Encode(results)
		} else {
			formatter.FormatResults(results)
		}
	}
}