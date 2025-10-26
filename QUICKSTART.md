# Quick Start Guide

Get AeyeWire MCP Service running in under 5 minutes!

## Prerequisites Check

Before starting, ensure you have:
- [ ] Go 1.21+ installed (`go version`)
- [ ] LMStudio installed and running
- [ ] qwen/qwen3-coder-30b model loaded in LMStudio

## 5-Minute Setup

### 1. Build (30 seconds)

```bash
# Install dependencies and build
make install
make build
```

### 2. Start LMStudio (1 minute)

1. Open LMStudio
2. Go to "Local Server" tab
3. Select qwen/qwen3-coder-30b model
4. Click "Start Server"
5. Verify it shows: "Server running on http://localhost:1234"

### 3. Verify Installation (30 seconds)

```bash
# Check version
./build/aeyewire_mcp version
# Output: AeyeWire MCP Service v1.0.0

# Check health
./build/aeyewire_mcp health
# Should show: LLM Service: available
```

### 4. Try Your First Analysis (2 minutes)

Analyze the sample Java file:

```bash
./build/aeyewire_mcp analyze test_samples/test.java
```

You should see a security report with findings like:
- SQL Injection vulnerabilities
- Hardcoded credentials
- Insecure random number generation
- Command injection issues

## Quick Usage Examples

### Analyze Different Languages

```bash
# Java
./build/aeyewire_mcp analyze test_samples/test.java

# C#
./build/aeyewire_mcp analyze test_samples/test.cs

# React TypeScript
./build/aeyewire_mcp analyze test_samples/test.tsx
```

### List Supported Languages

```bash
./build/aeyewire_mcp languages
```

### Check Service Health

```bash
./build/aeyewire_mcp health
```

## MCP Server Mode

To use with MCP-compatible IDEs:

```bash
# Run as stdio server
./build/aeyewire_mcp
```

The server will wait for MCP protocol messages on stdin.

## IDE Integration (Claude Code)

Add to your MCP configuration file:

**macOS/Linux**: `~/.config/claude-code/mcp_config.json`

**Windows**: `%APPDATA%\claude-code\mcp_config.json`

```json
{
  "mcpServers": {
    "aeyewire": {
      "command": "/absolute/path/to/aeyewire-mcp/build/aeyewire_mcp",
      "args": []
    }
  }
}
```

Replace `/absolute/path/to/` with your actual path.

## Common Issues

### "LLM Service: unavailable"

**Fix**: Ensure LMStudio is running and the server is started.

```bash
# Test LMStudio connection
curl http://localhost:1234/v1/models
```

### "Could not detect language"

**Fix**: Use files with proper extensions:
- `.java` for Java
- `.cs` for C#
- `.tsx` for React TypeScript
- `.jsx` for React JavaScript

### Build errors

**Fix**: Ensure Go is installed and in your PATH:

```bash
# Check Go installation
go version

# If not installed, download from: https://golang.org/dl/
```

## Next Steps

✅ **You're ready to go!** Now you can:

1. **Analyze your own code**: Replace test files with your actual source files
2. **Customize configuration**: Set environment variables for different LMStudio setups
3. **Integrate with IDE**: Use MCP protocol for seamless IDE integration
4. **Read full docs**: Check [README.md](README.md) and [INSTALL.md](INSTALL.md) for advanced usage

## Daily Workflow

```bash
# 1. Start LMStudio (if not already running)
# 2. Analyze your code
./build/aeyewire_mcp analyze path/to/your/file.java

# 3. Review security findings
# 4. Fix issues
# 5. Re-analyze to verify
```

## Getting Help

- **Installation issues**: See [INSTALL.md](INSTALL.md)
- **Usage details**: See [README.md](README.md)
- **Technical specs**: See [docs/specifications.md](docs/specifications.md)
- **Build from source**: See [CLAUDE.md](CLAUDE.md)

## Test Your Setup

Run all tests to ensure everything works:

```bash
# Run unit tests
make test

# Should show: PASS
```

---

**Congratulations!** You're now ready to use AeyeWire MCP Service for automated security analysis.
