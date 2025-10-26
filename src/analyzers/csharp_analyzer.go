package analyzers

import (
	"github.com/emware/aeyewire-mcp/src/models"
	"github.com/emware/aeyewire-mcp/src/services"
)

// CSharpAnalyzer performs security analysis on C# code
type CSharpAnalyzer struct {
	*BaseSecurityAnalyzer
}

// NewCSharpAnalyzer creates a new C# security analyzer
func NewCSharpAnalyzer(llmService *services.LLMService) *CSharpAnalyzer {
	return &CSharpAnalyzer{
		BaseSecurityAnalyzer: NewBaseAnalyzer(models.CSHARP, llmService),
	}
}

// Analyze performs security analysis on C# code
func (ca *CSharpAnalyzer) Analyze(code string, filePath string) (*models.AnalysisResult, error) {
	prompt := ca.GetSecurityRulesPrompt()
	return ca.AnalyzeWithLLM(code, filePath, prompt)
}

// GetSecurityRulesPrompt returns the security rules prompt for C#
func (ca *CSharpAnalyzer) GetSecurityRulesPrompt() string {
	return `Analyze the following C# code for security vulnerabilities. Check for these 20+ security issues:

INJECTION VULNERABILITIES:
1. SQL Injection - String concatenation in SQL queries, missing parameterized queries
2. Command Injection - Process.Start() or similar with unsanitized input
3. LDAP Injection - String concatenation in LDAP queries
4. XML Injection - Unsafe XML parsing allowing external entities

CRYPTOGRAPHIC ISSUES:
5. Weak Cryptography - DES, MD5, SHA1, hardcoded encryption keys
6. Insecure Random Number Generation - Random class for security purposes
7. Weak Password Hashing - Plain text or weak hashing algorithms

DESERIALIZATION:
8. Insecure Deserialization - BinaryFormatter, NetDataContractSerializer without validation

AUTHENTICATION & AUTHORIZATION:
9. Hardcoded Secrets - Passwords, API keys, connection strings in code
10. Authentication Bypass - Missing authorization checks, weak password policies
11. Session Management - Insecure session handling, missing timeout

PATH TRAVERSAL & FILE HANDLING:
12. Path Traversal - User input in file paths without validation
13. Insecure File Operations - Unrestricted file upload, missing validation

INPUT VALIDATION:
14. Input Validation Issues - Missing validation, regex DoS
15. Cross-Site Scripting (XSS) - Unencoded output in web applications

CODE SECURITY:
16. Code Injection - Dynamic code execution with user input (eval-like patterns)
17. Unsafe Reflection - Type.GetType() or Assembly.Load() with user input

CONFIGURATION & DEPLOYMENT:
18. Debug Mode in Production - Debug flags enabled
19. Information Disclosure - Detailed error messages, stack traces
20. Insecure Direct Object References - Missing access control checks

ADDITIONAL CONCERNS:
21. CSRF Protection - Missing anti-forgery tokens
22. Insecure Cookie Configuration - Missing HttpOnly, Secure flags
23. Open Redirect - Redirect with unvalidated user input

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
    "references": ["OWASP reference", "CWE-XXX"]
  }
]

Focus on actual vulnerabilities with specific line numbers and code snippets. If no issues are found, return an empty array [].`
}
