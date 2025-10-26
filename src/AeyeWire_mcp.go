package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/emware/aeyewire-mcp/src/analyzers"
	"github.com/emware/aeyewire-mcp/src/models"
	"github.com/emware/aeyewire-mcp/src/services"
)

const (
	VERSION     = "1.0.0"
	SERVER_NAME = "aeyewire_mcp"
)

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// MCPResponse represents an outgoing MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPServer handles MCP protocol communication
type MCPServer struct {
	llmService       *services.LLMService
	languageDetector *services.LanguageDetector
	analyzers        map[models.LanguageType]analyzers.SecurityAnalyzer
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	llmService := services.NewLLMService()
	languageDetector := services.NewLanguageDetector()

	server := &MCPServer{
		llmService:       llmService,
		languageDetector: languageDetector,
		analyzers:        make(map[models.LanguageType]analyzers.SecurityAnalyzer),
	}

	// Register analyzers
	server.analyzers[models.JAVA] = analyzers.NewJavaAnalyzer(llmService)
	server.analyzers[models.CSHARP] = analyzers.NewCSharpAnalyzer(llmService)
	server.analyzers[models.REACT_TYPESCRIPT] = analyzers.NewReactAnalyzer(llmService, models.REACT_TYPESCRIPT)
	server.analyzers[models.REACT_JAVASCRIPT] = analyzers.NewReactAnalyzer(llmService, models.REACT_JAVASCRIPT)

	return server
}

// Run starts the MCP server and processes stdio requests
func (s *MCPServer) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // 10MB buffer for large code

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			s.sendError(nil, -32700, fmt.Sprintf("Parse error: %v", err))
			continue
		}

		s.handleRequest(&request)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
	}
}

// handleRequest processes an MCP request
func (s *MCPServer) handleRequest(request *MCPRequest) {
	switch request.Method {
	case "initialize":
		s.handleInitialize(request)
	case "tools/list":
		s.handleToolsList(request)
	case "tools/call":
		s.handleToolsCall(request)
	default:
		s.sendError(request.ID, -32601, fmt.Sprintf("Method not found: %s", request.Method))
	}
}

// handleInitialize handles MCP initialize request
func (s *MCPServer) handleInitialize(request *MCPRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"serverInfo": map[string]interface{}{
			"name":    SERVER_NAME,
			"version": VERSION,
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]bool{},
		},
	}
	s.sendResponse(request.ID, result)
}

// handleToolsList handles tools/list request
func (s *MCPServer) handleToolsList(request *MCPRequest) {
	tools := []map[string]interface{}{
		{
			"name":        "analyze_security",
			"description": "Performs comprehensive security analysis on source code",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "Source code to analyze",
					},
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "File path for context and language detection (optional)",
					},
					"language": map[string]interface{}{
						"type":        "string",
						"description": "Programming language (csharp, react_typescript, react_javascript, java, auto)",
						"enum":        []string{"csharp", "react_typescript", "react_javascript", "java", "auto"},
					},
				},
				"required": []string{"code"},
			},
		},
		{
			"name":        "health_check",
			"description": "Verifies service health and dependency availability",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "list_supported_languages",
			"description": "Lists all supported programming languages and their metadata",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}
	s.sendResponse(request.ID, result)
}

// handleToolsCall handles tools/call request
func (s *MCPServer) handleToolsCall(request *MCPRequest) {
	toolName, ok := request.Params["name"].(string)
	if !ok {
		s.sendError(request.ID, -32602, "Invalid tool name")
		return
	}

	arguments, _ := request.Params["arguments"].(map[string]interface{})

	switch toolName {
	case "analyze_security":
		s.handleAnalyzeSecurity(request.ID, arguments)
	case "health_check":
		s.handleHealthCheck(request.ID)
	case "list_supported_languages":
		s.handleListSupportedLanguages(request.ID)
	default:
		s.sendError(request.ID, -32602, fmt.Sprintf("Unknown tool: %s", toolName))
	}
}

// handleAnalyzeSecurity handles the analyze_security tool
func (s *MCPServer) handleAnalyzeSecurity(requestID interface{}, args map[string]interface{}) {
	code, ok := args["code"].(string)
	if !ok || code == "" {
		s.sendError(requestID, -32602, "Missing or invalid 'code' parameter")
		return
	}

	filePath, _ := args["file_path"].(string)
	languageStr, _ := args["language"].(string)

	// Detect language
	var language models.LanguageType
	if languageStr != "" && languageStr != "auto" {
		language = models.LanguageType(languageStr)
	} else {
		language = s.languageDetector.Detect(code, filePath)
	}

	// Check if language is supported
	analyzer, ok := s.analyzers[language]
	if !ok {
		s.sendError(requestID, -32602, fmt.Sprintf("Unsupported language: %s", language))
		return
	}

	// Perform analysis
	result, err := analyzer.Analyze(code, filePath)
	if err != nil {
		s.sendError(requestID, -32603, fmt.Sprintf("Analysis failed: %v", err))
		return
	}

	// Format as markdown
	baseAnalyzer := analyzers.NewBaseAnalyzer(language, s.llmService)
	markdown := baseAnalyzer.FormatAsMarkdown(result)

	response := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": markdown,
			},
		},
	}

	s.sendResponse(requestID, response)
}

// handleHealthCheck handles the health_check tool
func (s *MCPServer) handleHealthCheck(requestID interface{}) {
	llmHealthy, _ := s.llmService.HealthCheck()

	llmStatus := "unavailable"
	if llmHealthy {
		llmStatus = "available"
	}

	supportedLanguages := []string{}
	for lang := range s.analyzers {
		supportedLanguages = append(supportedLanguages, string(lang))
	}

	healthResponse := models.HealthCheckResponse{
		Status:             "healthy",
		Version:            VERSION,
		LLMServiceStatus:   llmStatus,
		SupportedLanguages: supportedLanguages,
	}

	jsonData, _ := json.MarshalIndent(healthResponse, "", "  ")

	response := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(jsonData),
			},
		},
	}

	s.sendResponse(requestID, response)
}

// handleListSupportedLanguages handles the list_supported_languages tool
func (s *MCPServer) handleListSupportedLanguages(requestID interface{}) {
	languages := s.languageDetector.GetSupportedLanguages()

	jsonData, _ := json.MarshalIndent(languages, "", "  ")

	response := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(jsonData),
			},
		},
	}

	s.sendResponse(requestID, response)
}

// sendResponse sends an MCP response
func (s *MCPServer) sendResponse(id interface{}, result interface{}) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	jsonData, _ := json.Marshal(response)
	fmt.Println(string(jsonData))
}

// sendError sends an MCP error response
func (s *MCPServer) sendError(id interface{}, code int, message string) {
	response := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}

	jsonData, _ := json.Marshal(response)
	fmt.Println(string(jsonData))
}

func main() {
	// Check for command-line usage
	if len(os.Args) > 1 {
		handleCommandLine()
		return
	}

	// Run as MCP stdio server
	server := NewMCPServer()
	server.Run()
}

// handleCommandLine handles command-line tool usage
func handleCommandLine() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "analyze":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing file path")
			printUsage()
			os.Exit(1)
		}
		analyzeFile(os.Args[2])
	case "health":
		checkHealth()
	case "languages":
		listLanguages()
	case "version":
		fmt.Printf("AeyeWire MCP Service v%s\n", VERSION)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("AeyeWire MCP Service")
	fmt.Println("\nUsage:")
	fmt.Println("  aeyewire_mcp                  # Run as MCP stdio server")
	fmt.Println("  aeyewire_mcp analyze <file>   # Analyze a file")
	fmt.Println("  aeyewire_mcp health           # Check service health")
	fmt.Println("  aeyewire_mcp languages        # List supported languages")
	fmt.Println("  aeyewire_mcp version          # Show version")
}

func analyzeFile(filePath string) {
	// Read file
	codeBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	code := string(codeBytes)

	// Create server and analyze
	server := NewMCPServer()
	language := server.languageDetector.Detect(code, filePath)

	if language == models.UNKNOWN {
		fmt.Println("Error: Could not detect language")
		os.Exit(1)
	}

	analyzer, ok := server.analyzers[language]
	if !ok {
		fmt.Printf("Error: Unsupported language: %s\n", language)
		os.Exit(1)
	}

	fmt.Printf("Analyzing %s as %s...\n\n", filePath, language)

	result, err := analyzer.Analyze(code, filePath)
	if err != nil {
		fmt.Printf("Analysis failed: %v\n", err)
		os.Exit(1)
	}

	// Format and print
	baseAnalyzer := analyzers.NewBaseAnalyzer(language, server.llmService)
	markdown := baseAnalyzer.FormatAsMarkdown(result)
	fmt.Println(markdown)
}

func checkHealth() {
	server := NewMCPServer()
	llmHealthy, err := server.llmService.HealthCheck()

	fmt.Printf("Service Status: healthy\n")
	fmt.Printf("Version: %s\n", VERSION)

	if llmHealthy {
		fmt.Printf("LLM Service: available\n")
	} else {
		fmt.Printf("LLM Service: unavailable")
		if err != nil {
			fmt.Printf(" (%v)", err)
		}
		fmt.Println()
	}

	fmt.Println("\nSupported Languages:")
	for lang := range server.analyzers {
		fmt.Printf("  - %s\n", lang)
	}
}

func listLanguages() {
	server := NewMCPServer()
	languages := server.languageDetector.GetSupportedLanguages()

	fmt.Println("Supported Languages:")
	for _, lang := range languages {
		fmt.Printf("\n%s (%s)\n", strings.ToUpper(lang.Identifier), lang.Description)
		fmt.Printf("  Extensions: %s\n", strings.Join(lang.Extensions, ", "))
	}
}
