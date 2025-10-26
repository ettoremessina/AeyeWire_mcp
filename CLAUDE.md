# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AeyeWire MCP Service is a Model Context Protocol (MCP) service that performs automated security analysis on source code using LMStudio with the qwen/qwen3-coder-30b language model. The service exposes a standard stdio interface for integration with IDEs and can also be used as a command-line tool.

## Architecture

### Core Component Structure

The service is built around four main components:

1. **MCP Service** (`src/AeyeWire_mcp.go`)
   - Entry point handling tool registration and execution
   - Manages stdio-based MCP integration
   - Exposes three tools: `analyze_security`, `health_check`, `list_supported_languages`

2. **Language Detection Service**
   - Detects programming languages from code content and file extensions
   - Supports: C# (.cs), React TypeScript (.tsx, .ts), React JavaScript (.jsx, .js), Java (.java)
   - Uses pattern matching as fallback when extension-based detection is unavailable

3. **LLM Service**
   - Integrates with LMStudio using langchain
   - Communicates with qwen/qwen3-coder-30b model
   - Handles prompt engineering for security analysis and response parsing

4. **Specialized Analyzers**
   - Base analyzer class (`BaseSecurityAnalyzer`) with common functionality
   - Language-specific analyzers (C#, React, Java) inheriting from base
   - Each analyzer implements 20-25+ security rules specific to the language

### Data Flow

1. Code submission → Language detection (extension or pattern-based)
2. Preprocessing (comment removal, structure preservation)
3. LLM analysis with language-specific security rule prompts
4. Response parsing into `SecurityIssue` objects
5. Report generation (structured markdown with severity categorization)

## Data Models

**SecurityIssue**: Individual vulnerability representation
- Fields: id, title, description, severity (LOW/MEDIUM/HIGH/CRITICAL), line_number, column_number, file_path, code_snippet, remediation, references

**AnalysisRequest**: Input model
- Fields: code, file_path (optional), language (optional, defaults to "auto")

**AnalysisResult**: Output model
- Fields: language, issues[], summary, analysis_metadata

**LanguageType**: Enum of CSHARP, REACT_TYPESCRIPT, REACT_JAVASCRIPT, JAVA, UNKNOWN

## Configuration

Default environment variables:
- `LMSTUDIO_BASE_URL`: http://localhost:1234
- `LMSTUDIO_MODEL`: qwen/qwen3-coder-30b
- `LMSTUDIO_API_KEY`: (optional)
- `MCP_SERVER_NAME`: aeyewire_mcp
- `MCP_SERVER_VERSION`: 1.0.0

## Security Analysis Rules

### Common Vulnerability Categories Across Languages

**Injection Vulnerabilities** (CRITICAL/HIGH):
- SQL injection (parameterized queries required)
- Command injection (input sanitization, allowlists)
- LDAP injection (Java), XPath injection, JNDI injection (Java)
- XXE (XML External Entity) attacks (Java)

**Cryptographic Issues** (HIGH/MEDIUM):
- Weak algorithms (DES, 3DES, RC4, MD5, SHA1)
- Insecure random number generation
- Hardcoded encryption keys
- Insecure SSL/TLS configuration (Java)

**Authentication & Session** (CRITICAL/HIGH):
- Hardcoded credentials
- Session management flaws
- Authentication bypass vulnerabilities

**Path Traversal & File Handling** (HIGH):
- Directory traversal attacks
- Insecure file uploads
- Resource leaks

**Code Execution** (CRITICAL/HIGH):
- Insecure deserialization
- Unsafe reflection (Java)
- Expression Language injection (Java - JSP/JSF/Spring/OGNL)

**Web-Specific** (varies):
- XSS prevention (React)
- SSRF (Server-Side Request Forgery)
- Unvalidated redirects
- CSRF protection

**Input Validation** (MEDIUM/LOW):
- Regex DoS (ReDoS)
- Log injection
- Mass assignment vulnerabilities

### Language-Specific Notes

**Java**: 25+ rules covering JNDI injection, JNI security, expression language injection, concurrency/race conditions
**React**: Focus on XSS, insecure state management, props validation, API security, type safety
**C#**: 20+ rules for .NET-specific vulnerabilities

## Adding New Language Support

To extend support for a new language:

1. Add new value to `LanguageType` enum
2. Create analyzer class inheriting from `BaseSecurityAnalyzer`
3. Implement language-specific security rules (20+ recommended)
4. Update language detector patterns
5. Register analyzer in MCP service
6. Add language-specific preprocessing rules
7. Update `list_supported_languages` tool response

## Testing Strategy

Required test coverage areas:
- Unit tests for language detection (extension and pattern-based)
- Integration tests for each analyzer with language-specific test cases
- End-to-end tests for MCP service stdio interface
- Installation verification scripts

Quality metrics to track:
- Security rule coverage per language
- Analysis accuracy (false positives/negatives)
- Performance benchmarks for code analysis
- Error rate monitoring

## Implementation Stack

- **Language**: Go (primary implementation language)
- **LLM Integration**: langchain for LMStudio communication
- **Protocol**: Model Context Protocol (MCP) with stdio interface
- **External Dependency**: LMStudio running qwen/qwen3-coder-30b model

## Error Handling Philosophy

Graceful degradation approach:
- LLM unavailable → return error with health status
- Language detection fails → fallback to pattern matching
- Analysis failures → return partial results with error metadata
- Invalid input → clear, actionable error messages

Include structured error reporting in analysis metadata for all failure modes.
