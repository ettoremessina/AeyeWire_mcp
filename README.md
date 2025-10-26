# AeyeWire MCP Service

A Model Context Protocol (MCP) service that performs automated security analysis on source code using LMStudio with the qwen/qwen3-coder-30b language model.

## Features

- **Multi-Language Support**: C#, Java, React (TypeScript/JavaScript)
- **MCP Integration**: Standard stdio interface for IDE integration
- **CLI Mode**: Command-line tool for standalone usage
- **Comprehensive Analysis**: 20-25+ security rules per language
- **LLM-Powered**: Leverages advanced AI for intelligent security detection

## Prerequisites

- Go 1.21 or higher
- LMStudio running with qwen/qwen3-coder-30b model
- LMStudio server accessible at `http://localhost:1234` (default)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd aeyewire-mcp
```

2. Install dependencies:
```bash
make install
```

3. Build the binary:
```bash
make build
```

## Configuration

Configure via environment variables (all optional with defaults):

```bash
export LMSTUDIO_BASE_URL="http://localhost:1234"  # LMStudio server URL
export LMSTUDIO_MODEL="qwen/qwen3-coder-30b"      # Model name
export LMSTUDIO_API_KEY=""                         # API key (if required)
export MCP_SERVER_NAME="aeyewire_mcp"            # Server identifier
export MCP_SERVER_VERSION="1.0.0"                 # Service version
```

## Usage

### MCP Server Mode

Run as an MCP stdio server for IDE integration:

```bash
./build/aeyewire_mcp
```

Or with make:

```bash
make run
```

### Command-Line Mode

Analyze a specific file:

```bash
./build/aeyewire_mcp analyze path/to/file.java
```

Check service health:

```bash
./build/aeyewire_mcp health
# or
make health
```

List supported languages:

```bash
./build/aeyewire_mcp languages
# or
make languages
```

Show version:

```bash
./build/aeyewire_mcp version
# or
make version
```

## MCP Tools

The service exposes three MCP tools:

### 1. analyze_security

Performs comprehensive security analysis on source code.

**Parameters**:
- `code` (string, required): Source code to analyze
- `file_path` (string, optional): File path for context
- `language` (string, optional): Language override (csharp, java, react_typescript, react_javascript, auto)

**Returns**: Markdown-formatted security report

### 2. health_check

Verifies service health and dependency availability.

**Parameters**: None

**Returns**: JSON health status

### 3. list_supported_languages

Lists all supported programming languages.

**Parameters**: None

**Returns**: JSON array of language metadata

## Supported Languages

- **C#** (.cs) - 20+ security rules
- **Java** (.java) - 25+ security rules
- **React TypeScript** (.tsx, .ts) - 20+ security rules
- **React JavaScript** (.jsx, .js) - 20+ security rules

## Development

Format code:
```bash
make fmt
```

Run linter:
```bash
make lint
```

Run tests:
```bash
make test
```

Clean build artifacts:
```bash
make clean
```

## Project Structure

```
.
├── src/
│   ├── AeyeWire_mcp.go           # Main MCP service
│   ├── models/
│   │   └── models.go              # Data models
│   ├── services/
│   │   ├── language_detector.go   # Language detection
│   │   └── llm_service.go         # LLM integration
│   └── analyzers/
│       ├── base_analyzer.go       # Base analyzer
│       ├── java_analyzer.go       # Java analyzer
│       ├── csharp_analyzer.go     # C# analyzer
│       └── react_analyzer.go      # React analyzer
├── docs/
│   └── specifications.md          # Detailed specifications
├── go.mod                         # Go module file
├── Makefile                       # Build automation
└── README.md                      # This file
```

## Security Analysis Coverage

### Common Vulnerabilities Detected

- SQL Injection
- Command Injection
- XSS (Cross-Site Scripting)
- Path Traversal
- Insecure Deserialization
- Weak Cryptography
- Hardcoded Credentials
- Authentication Bypass
- CSRF Vulnerabilities
- And many more...

See [docs/specifications.md](docs/specifications.md) for complete security rule details.

## Troubleshooting

**LLM Service Unavailable**:
- Ensure LMStudio is running
- Verify the model is loaded
- Check `LMSTUDIO_BASE_URL` configuration

**Language Detection Issues**:
- Provide `file_path` parameter for extension-based detection
- Use `language` parameter to override detection

**Analysis Timeout**:
- Large files may take longer to analyze
- Default timeout is 120 seconds

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please see CONTRIBUTING.md for guidelines.
