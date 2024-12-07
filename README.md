# Search CLI

A command-line interface for searching using Serper and Tavily APIs.

## Features

- Search using Google (via Serper API)
- Search and extract content using Tavily API
- JSON output support
- Configurable search depth (for Tavily)

## Installation

### Using Go Install
```bash
go install github.com/regismesquita/search-cli/cmd/search@latest
```

### Using mise (formerly rtx)
```bash
mise use -g go@latest
mise install github.com/regismesquita/search-cli/cmd/search@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/regismesquita/search-cli
cd search-cli

# Build and install
make build
make install      # Installs to ~/go/bin
# or
make local-install # Installs to /usr/local/bin (requires sudo)
```

## Configuration

Set your API keys as environment variables:
```bash
# For Serper
export SERPER_API_KEY=your_key_here

# For Tavily
export TAVILY_API_KEY=your_key_here
```

## Usage

### Basic Search (using Serper)
```bash
# Default search using Serper
search "your query"
# or explicitly
search -s "your query"
```

### Using Tavily
```bash
# Basic search
search -t "your query"

# Advanced search
search -t -depth advanced "your query"

# Extract content from URLs
search -t -e "https://example.com"
```

### JSON Output
```bash
search -json "your query"
search -t -json "your query"
search -t -e -json "https://example.com"
```

## Options

- `-s`: Use Serper (default)
- `-t`: Use Tavily
- `-e`: Extract content (Tavily only)
- `-json`: Output in JSON format
- `-depth`: Search depth for Tavily (basic or advanced)

## Development

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Clean build artifacts
make clean
```

## Project Structure

```
.
├── cmd/
│   └── search/        # Main application
│       └── main.go
├── internal/
│   ├── adapters/      # API providers
│   │   ├── serper.go
│   │   ├── tavily.go
│   │   └── types.go
│   └── formatter/     # Output formatting
│       └── formatter.go
├── go.mod
└── Makefile
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request