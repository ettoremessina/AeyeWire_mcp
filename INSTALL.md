# AeyeWire MCP Service - Installation Guide

## Prerequisites

### Required Software

1. **Go** (version 1.21 or higher)
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **LMStudio** with qwen/qwen3-coder-30b model
   - Download from: https://lmstudio.ai/
   - Install and launch LMStudio
   - Download the qwen/qwen3-coder-30b model
   - Start the local server (default port: 1234)

### Optional Software

- **Make** - For using Makefile commands (usually pre-installed on Unix systems)
- **golangci-lint** - For code linting (optional, for development)

## Installation Steps

### Step 1: Clone/Download the Repository

```bash
# If using git
git clone <repository-url>
cd aeyewire-mcp

# Or extract from archive
unzip aeyewire-mcp.zip
cd aeyewire-mcp
```

### Step 2: Install Go Dependencies

```bash
# Using make
make install

# Or manually
go mod download
go mod tidy
```

### Step 3: Build the Binary

```bash
# Using make
make build

# Or manually
mkdir -p build
go build -o build/aeyewire_mcp src/AeyeWire_mcp.go
```

The binary will be created at `build/aeyewire_mcp`.

### Step 4: Configure Environment (Optional)

Create a `.env` file or set environment variables:

```bash
export LMSTUDIO_BASE_URL="http://localhost:1234"
export LMSTUDIO_MODEL="qwen/qwen3-coder-30b"
export LMSTUDIO_API_KEY=""  # Leave empty if not required
```

### Step 5: Verify Installation

Test the installation:

```bash
# Check version
./build/aeyewire_mcp version

# List supported languages
./build/aeyewire_mcp languages

# Check health (requires LMStudio running)
./build/aeyewire_mcp health
```

## LMStudio Setup

### Starting LMStudio Server

1. Open LMStudio
2. Go to the "Local Server" tab
3. Select the qwen/qwen3-coder-30b model
4. Click "Start Server"
5. Verify it's running at `http://localhost:1234`

### Verifying LMStudio Connection

```bash
# Check if LMStudio is accessible
curl http://localhost:1234/v1/models

# Or use the health check
./build/aeyewire_mcp health
```

Expected output for health check:
```
Service Status: healthy
Version: 1.0.0
LLM Service: available

Supported Languages:
  - java
  - csharp
  - react_typescript
  - react_javascript
```

## Usage Examples

### Command-Line Usage

Analyze a Java file:
```bash
./build/aeyewire_mcp analyze test_samples/test.java
```

Analyze a C# file:
```bash
./build/aeyewire_mcp analyze test_samples/test.cs
```

Analyze a React TypeScript file:
```bash
./build/aeyewire_mcp analyze test_samples/test.tsx
```

### MCP Server Mode

For IDE integration, run as an MCP stdio server:

```bash
./build/aeyewire_mcp
```

The server will listen on stdin/stdout for MCP protocol messages.

## IDE Integration

### Claude Code / MCP-Compatible IDEs

Add to your MCP configuration:

```json
{
  "mcpServers": {
    "aeyewire": {
      "command": "/path/to/aeyewire-mcp/build/aeyewire_mcp",
      "args": [],
      "env": {
        "LMSTUDIO_BASE_URL": "http://localhost:1234",
        "LMSTUDIO_MODEL": "qwen/qwen3-coder-30b"
      }
    }
  }
}
```

## Troubleshooting

### Build Fails

**Error**: `go: command not found`
- **Solution**: Install Go from https://golang.org/dl/

**Error**: Module download fails
- **Solution**: Check internet connection and Go proxy settings
- Try: `export GOPROXY=https://proxy.golang.org,direct`

### LMStudio Connection Issues

**Error**: `LLM Service: unavailable`
- **Solution**:
  1. Verify LMStudio is running
  2. Check the model is loaded
  3. Confirm server URL: `http://localhost:1234`
  4. Test with: `curl http://localhost:1234/v1/models`

**Error**: Connection timeout
- **Solution**:
  1. Check firewall settings
  2. Verify LMStudio server port (default 1234)
  3. Try setting: `export LMSTUDIO_BASE_URL="http://127.0.0.1:1234"`

### Analysis Issues

**Error**: `Could not detect language`
- **Solution**:
  1. Ensure file has proper extension (.java, .cs, .tsx, .jsx, .ts, .js)
  2. Check file contains valid code
  3. Use language override parameter

**Error**: Analysis timeout
- **Solution**:
  1. Reduce code size (analyze in chunks)
  2. Ensure LMStudio model is loaded and responsive
  3. Check system resources (CPU, memory)

## Uninstallation

```bash
# Remove build artifacts
make clean

# Remove binary
rm -rf build/

# Optional: Remove entire directory
cd ..
rm -rf aeyewire-mcp/
```

## Updates and Maintenance

### Updating Dependencies

```bash
go get -u ./...
go mod tidy
make build
```

### Rebuilding

```bash
make clean
make build
```

## Support

For issues and questions:
- Check the [README.md](README.md) for usage examples
- Review [docs/specifications.md](docs/specifications.md) for technical details
- Submit issues to the project repository

## Next Steps

After installation:
1. Test with sample files in `test_samples/`
2. Review the [README.md](README.md) for detailed usage
3. Configure IDE integration for MCP protocol
4. Customize security rules if needed

## Development Setup

If you plan to contribute or modify the code:

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```
