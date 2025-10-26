package analyzers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/emware/aeyewire-mcp/src/models"
	"github.com/emware/aeyewire-mcp/src/services"
)

// BaseSecurityAnalyzer provides common functionality for all analyzers
type BaseSecurityAnalyzer struct {
	Language   models.LanguageType
	LLMService *services.LLMService
}

// SecurityAnalyzer interface that all analyzers must implement
type SecurityAnalyzer interface {
	Analyze(code string, filePath string) (*models.AnalysisResult, error)
	GetSecurityRulesPrompt() string
}

// NewBaseAnalyzer creates a new base analyzer
func NewBaseAnalyzer(language models.LanguageType, llmService *services.LLMService) *BaseSecurityAnalyzer {
	return &BaseSecurityAnalyzer{
		Language:   language,
		LLMService: llmService,
	}
}

// PreprocessCode removes comments while maintaining line structure
func (ba *BaseSecurityAnalyzer) PreprocessCode(code string, language models.LanguageType) string {
	switch language {
	case models.JAVA, models.CSHARP, models.REACT_TYPESCRIPT, models.REACT_JAVASCRIPT:
		return ba.removeComments(code)
	default:
		return code
	}
}

// removeComments removes single-line and multi-line comments
func (ba *BaseSecurityAnalyzer) removeComments(code string) string {
	// Remove multi-line comments (/* */ and /** */)
	multiLineComment := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	code = multiLineComment.ReplaceAllString(code, "")

	// Remove single-line comments (//)
	singleLineComment := regexp.MustCompile(`//.*`)
	code = singleLineComment.ReplaceAllString(code, "")

	return code
}

// AnalyzeWithLLM performs LLM-based security analysis
func (ba *BaseSecurityAnalyzer) AnalyzeWithLLM(code string, filePath string, securityRulesPrompt string) (*models.AnalysisResult, error) {
	startTime := time.Now()

	// Preprocess code
	preprocessed := ba.PreprocessCode(code, ba.Language)

	// Perform LLM analysis
	response, err := ba.LLMService.Analyze(preprocessed, securityRulesPrompt)
	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Parse LLM response
	issues, err := ba.parseIssuesFromResponse(response, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// Generate metadata
	metadata := ba.generateMetadata(issues, ba.Language, time.Since(startTime))

	// Generate summary
	summary := ba.generateSummary(issues)

	return &models.AnalysisResult{
		Language:         ba.Language,
		Issues:           issues,
		Summary:          summary,
		AnalysisMetadata: metadata,
	}, nil
}

// parseIssuesFromResponse parses SecurityIssue objects from LLM response
func (ba *BaseSecurityAnalyzer) parseIssuesFromResponse(response string, filePath string) ([]models.SecurityIssue, error) {
	// Extract JSON from response (it might be wrapped in markdown code blocks)
	jsonStr := ba.extractJSON(response)

	var issues []models.SecurityIssue
	if err := json.Unmarshal([]byte(jsonStr), &issues); err != nil {
		// If direct parsing fails, try wrapping in an object
		var wrapper struct {
			Issues []models.SecurityIssue `json:"issues"`
		}
		if err2 := json.Unmarshal([]byte(jsonStr), &wrapper); err2 != nil {
			return nil, fmt.Errorf("failed to parse issues: %w", err)
		}
		issues = wrapper.Issues
	}

	// Enrich issues with file path
	for i := range issues {
		if issues[i].FilePath == "" {
			issues[i].FilePath = filePath
		}
		if issues[i].ID == "" {
			issues[i].ID = fmt.Sprintf("ISSUE-%d", i+1)
		}
	}

	return issues, nil
}

// extractJSON extracts JSON content from markdown code blocks or plain text
func (ba *BaseSecurityAnalyzer) extractJSON(response string) string {
	// Try to extract from markdown code block
	codeBlockRegex := regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)```")
	matches := codeBlockRegex.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// If no code block, look for JSON array or object
	jsonRegex := regexp.MustCompile(`(\[[\s\S]*\]|\{[\s\S]*\})`)
	matches = jsonRegex.FindStringSubmatch(response)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return response
}

// generateMetadata creates analysis metadata
func (ba *BaseSecurityAnalyzer) generateMetadata(issues []models.SecurityIssue, language models.LanguageType, duration time.Duration) models.AnalysisMetadata {
	metadata := models.AnalysisMetadata{
		AnalysisTime:     duration.String(),
		IssuesFound:      len(issues),
		DetectedLanguage: language,
		CriticalCount:    0,
		HighCount:        0,
		MediumCount:      0,
		LowCount:         0,
	}

	for _, issue := range issues {
		switch issue.Severity {
		case models.CRITICAL:
			metadata.CriticalCount++
		case models.HIGH:
			metadata.HighCount++
		case models.MEDIUM:
			metadata.MediumCount++
		case models.LOW:
			metadata.LowCount++
		}
	}

	return metadata
}

// generateSummary creates a summary of the analysis
func (ba *BaseSecurityAnalyzer) generateSummary(issues []models.SecurityIssue) string {
	if len(issues) == 0 {
		return "No security issues detected."
	}

	critical := 0
	high := 0
	medium := 0
	low := 0

	for _, issue := range issues {
		switch issue.Severity {
		case models.CRITICAL:
			critical++
		case models.HIGH:
			high++
		case models.MEDIUM:
			medium++
		case models.LOW:
			low++
		}
	}

	return fmt.Sprintf("Found %d security issue(s): %d critical, %d high, %d medium, %d low.",
		len(issues), critical, high, medium, low)
}

// FormatAsMarkdown formats the analysis result as markdown
func (ba *BaseSecurityAnalyzer) FormatAsMarkdown(result *models.AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString("# Security Analysis Report\n\n")
	sb.WriteString(fmt.Sprintf("**Language**: %s\n\n", result.Language))
	sb.WriteString(fmt.Sprintf("**Analysis Time**: %s\n\n", result.AnalysisMetadata.AnalysisTime))
	sb.WriteString(fmt.Sprintf("## Summary\n\n%s\n\n", result.Summary))

	if len(result.Issues) == 0 {
		sb.WriteString("No security issues found.\n")
		return sb.String()
	}

	// Group issues by severity
	criticalIssues := []models.SecurityIssue{}
	highIssues := []models.SecurityIssue{}
	mediumIssues := []models.SecurityIssue{}
	lowIssues := []models.SecurityIssue{}

	for _, issue := range result.Issues {
		switch issue.Severity {
		case models.CRITICAL:
			criticalIssues = append(criticalIssues, issue)
		case models.HIGH:
			highIssues = append(highIssues, issue)
		case models.MEDIUM:
			mediumIssues = append(mediumIssues, issue)
		case models.LOW:
			lowIssues = append(lowIssues, issue)
		}
	}

	// Write issues by severity
	if len(criticalIssues) > 0 {
		sb.WriteString("## Critical Issues\n\n")
		for _, issue := range criticalIssues {
			ba.writeIssue(&sb, issue)
		}
	}

	if len(highIssues) > 0 {
		sb.WriteString("## High Severity Issues\n\n")
		for _, issue := range highIssues {
			ba.writeIssue(&sb, issue)
		}
	}

	if len(mediumIssues) > 0 {
		sb.WriteString("## Medium Severity Issues\n\n")
		for _, issue := range mediumIssues {
			ba.writeIssue(&sb, issue)
		}
	}

	if len(lowIssues) > 0 {
		sb.WriteString("## Low Severity Issues\n\n")
		for _, issue := range lowIssues {
			ba.writeIssue(&sb, issue)
		}
	}

	return sb.String()
}

// writeIssue writes a single issue to the markdown builder
func (ba *BaseSecurityAnalyzer) writeIssue(sb *strings.Builder, issue models.SecurityIssue) {
	sb.WriteString(fmt.Sprintf("### %s\n\n", issue.Title))
	sb.WriteString(fmt.Sprintf("**Severity**: %s\n\n", issue.Severity))
	sb.WriteString(fmt.Sprintf("**Description**: %s\n\n", issue.Description))

	if issue.LineNumber > 0 {
		sb.WriteString(fmt.Sprintf("**Location**: Line %d", issue.LineNumber))
		if issue.ColumnNumber > 0 {
			sb.WriteString(fmt.Sprintf(", Column %d", issue.ColumnNumber))
		}
		sb.WriteString("\n\n")
	}

	if issue.CodeSnippet != "" {
		sb.WriteString("**Code Snippet**:\n```\n")
		sb.WriteString(issue.CodeSnippet)
		sb.WriteString("\n```\n\n")
	}

	if issue.Remediation != "" {
		sb.WriteString(fmt.Sprintf("**Remediation**: %s\n\n", issue.Remediation))
	}

	if len(issue.References) > 0 {
		sb.WriteString("**References**:\n")
		for _, ref := range issue.References {
			sb.WriteString(fmt.Sprintf("- %s\n", ref))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
}
