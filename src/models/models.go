package models

// LanguageType represents supported programming languages
type LanguageType string

const (
	CSHARP             LanguageType = "csharp"
	REACT_TYPESCRIPT   LanguageType = "react_typescript"
	REACT_JAVASCRIPT   LanguageType = "react_javascript"
	JAVA               LanguageType = "java"
	UNKNOWN            LanguageType = "unknown"
)

// SeverityLevel represents the severity of a security issue
type SeverityLevel string

const (
	LOW      SeverityLevel = "LOW"
	MEDIUM   SeverityLevel = "MEDIUM"
	HIGH     SeverityLevel = "HIGH"
	CRITICAL SeverityLevel = "CRITICAL"
)

// SecurityIssue represents an individual security vulnerability
type SecurityIssue struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Severity     SeverityLevel `json:"severity"`
	LineNumber   int           `json:"line_number"`
	ColumnNumber int           `json:"column_number"`
	FilePath     string        `json:"file_path"`
	CodeSnippet  string        `json:"code_snippet"`
	Remediation  string        `json:"remediation"`
	References   []string      `json:"references"`
}

// AnalysisRequest represents input for security analysis
type AnalysisRequest struct {
	Code     string       `json:"code"`
	FilePath string       `json:"file_path,omitempty"`
	Language LanguageType `json:"language,omitempty"`
}

// AnalysisMetadata contains metadata about the analysis
type AnalysisMetadata struct {
	AnalysisTime    string            `json:"analysis_time"`
	IssuesFound     int               `json:"issues_found"`
	CriticalCount   int               `json:"critical_count"`
	HighCount       int               `json:"high_count"`
	MediumCount     int               `json:"medium_count"`
	LowCount        int               `json:"low_count"`
	DetectedLanguage LanguageType     `json:"detected_language"`
	Errors          []string          `json:"errors,omitempty"`
}

// AnalysisResult represents the output of security analysis
type AnalysisResult struct {
	Language          LanguageType     `json:"language"`
	Issues            []SecurityIssue  `json:"issues"`
	Summary           string           `json:"summary"`
	AnalysisMetadata  AnalysisMetadata `json:"analysis_metadata"`
}

// HealthCheckResponse represents the health status of the service
type HealthCheckResponse struct {
	Status            string       `json:"status"`
	Version           string       `json:"version"`
	LLMServiceStatus  string       `json:"llm_service_status"`
	SupportedLanguages []string    `json:"supported_languages"`
}

// LanguageInfo represents metadata about a supported language
type LanguageInfo struct {
	Identifier   string   `json:"identifier"`
	Description  string   `json:"description"`
	Extensions   []string `json:"extensions"`
}
