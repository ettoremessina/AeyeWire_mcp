package analyzers

import (
	"github.com/emware/aeyewire-mcp/src/models"
	"github.com/emware/aeyewire-mcp/src/services"
)

// ReactAnalyzer performs security analysis on React code (TypeScript and JavaScript)
type ReactAnalyzer struct {
	*BaseSecurityAnalyzer
}

// NewReactAnalyzer creates a new React security analyzer
func NewReactAnalyzer(llmService *services.LLMService, language models.LanguageType) *ReactAnalyzer {
	return &ReactAnalyzer{
		BaseSecurityAnalyzer: NewBaseAnalyzer(language, llmService),
	}
}

// Analyze performs security analysis on React code
func (ra *ReactAnalyzer) Analyze(code string, filePath string) (*models.AnalysisResult, error) {
	prompt := ra.GetSecurityRulesPrompt()
	return ra.AnalyzeWithLLM(code, filePath, prompt)
}

// GetSecurityRulesPrompt returns the security rules prompt for React
func (ra *ReactAnalyzer) GetSecurityRulesPrompt() string {
	basePrompt := `Analyze the following React code for security vulnerabilities. Check for these security issues:

XSS (CROSS-SITE SCRIPTING):
1. Dangerous HTML Rendering - dangerouslySetInnerHTML without sanitization
2. Unescaped User Input - Direct rendering of user input in JSX
3. URL Injection - Unsafe href or src attributes with user input
4. Unsafe Attribute Binding - User-controlled event handlers

STATE & PROPS SECURITY:
5. Insecure State Management - Sensitive data in client-side state
6. Props Validation - Missing PropTypes or TypeScript types for security-critical props
7. State Mutation - Direct state mutations bypassing security checks

API & DATA HANDLING:
8. Insecure API Calls - Hardcoded API keys, credentials in code
9. CSRF Protection - Missing CSRF tokens in API requests
10. API Endpoint Exposure - Sensitive endpoints or data exposed
11. Insecure Data Storage - Sensitive data in localStorage/sessionStorage

AUTHENTICATION & AUTHORIZATION:
12. Client-Side Auth Logic - Authentication decisions made purely on client
13. Token Storage - Insecure JWT or token storage
14. Missing Authorization Checks - Routes/components without proper access control

INPUT VALIDATION:
15. Form Validation - Missing or client-only validation
16. File Upload Security - Unrestricted file uploads
17. Regex DoS - Vulnerable regular expressions

CONFIGURATION:
18. Debug Code - console.log with sensitive data, debug flags in production
19. Error Handling - Detailed error messages exposing system information
20. Insecure Dependencies - Known vulnerabilities in npm packages

REACT-SPECIFIC:
21. Unsafe Refs - Direct DOM manipulation bypassing React security
22. Third-Party Components - Untrusted or unvalidated component usage
23. Code Injection - eval(), Function constructor, or dynamic code execution`

	if ra.Language == models.REACT_TYPESCRIPT {
		basePrompt += `

TYPESCRIPT-SPECIFIC:
24. Type Safety Bypass - 'any' type for security-critical data
25. Type Assertions - Unsafe type casting that bypasses security checks
26. Missing Null Checks - Potential null/undefined without proper guards`
	}

	basePrompt += `

Return findings as a JSON array of security issues with this structure:
[
  {
    "id": "unique-id",
    "title": "Issue title",
    "description": "Detailed description",
    "severity": "CRITICAL|HIGH|MEDIUM|LOW",
    "line_number": 0,
    "column_number": 0,
    "code_snippet": "vulnerable code",
    "remediation": "How to fix",
    "references": ["OWASP reference", "React Security Best Practices"]
  }
]

Focus on actual vulnerabilities with specific line numbers and code snippets. If no issues are found, return an empty array [].`

	return basePrompt
}
