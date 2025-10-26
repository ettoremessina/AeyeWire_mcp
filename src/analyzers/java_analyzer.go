package analyzers

import (
	"github.com/emware/aeyewire-mcp/src/models"
	"github.com/emware/aeyewire-mcp/src/services"
)

// JavaAnalyzer performs security analysis on Java code
type JavaAnalyzer struct {
	*BaseSecurityAnalyzer
}

// NewJavaAnalyzer creates a new Java security analyzer
func NewJavaAnalyzer(llmService *services.LLMService) *JavaAnalyzer {
	return &JavaAnalyzer{
		BaseSecurityAnalyzer: NewBaseAnalyzer(models.JAVA, llmService),
	}
}

// Analyze performs security analysis on Java code
func (ja *JavaAnalyzer) Analyze(code string, filePath string) (*models.AnalysisResult, error) {
	prompt := ja.GetSecurityRulesPrompt()
	return ja.AnalyzeWithLLM(code, filePath, prompt)
}

// GetSecurityRulesPrompt returns the security rules prompt for Java
func (ja *JavaAnalyzer) GetSecurityRulesPrompt() string {
	return `Analyze the following Java code for security vulnerabilities. Check for these 25+ security issues:

INJECTION VULNERABILITIES:
1. SQL Injection - String concatenation in SQL queries, missing PreparedStatement
2. Command Injection - Runtime.exec() or ProcessBuilder with unsanitized input
3. LDAP Injection - String concatenation in LDAP filters
4. XXE (XML External Entity) - DocumentBuilderFactory without disabled external entities
5. JNDI Injection - Context.lookup() with user-controlled strings

CRYPTOGRAPHIC ISSUES:
6. Weak Cryptography - DES, 3DES, RC4, MD5, SHA1, ECB mode, hardcoded keys
7. Insecure Random Number Generation - java.util.Random or Math.random() for security
8. Insecure SSL/TLS Configuration - Trusting all certificates, disabled hostname verification

DESERIALIZATION:
9. Insecure Deserialization - ObjectInputStream.readObject() on untrusted data

AUTHENTICATION & SESSION:
10. Hardcoded Credentials - Passwords, API keys, secrets in code
11. Session Management Flaws - Session IDs in URLs, missing timeout, no regeneration
12. Authentication Bypass - Missing authentication checks, weak password policies

PATH TRAVERSAL & FILE HANDLING:
13. Path Traversal - User input in file paths without validation, ../ sequences
14. Insecure File Upload - No file type validation, missing size limits
15. Resource Leaks - Missing try-with-resources, unclosed connections

CODE EXECUTION & REFLECTION:
16. Unsafe Reflection - Class.forName() or Method.invoke() with user input
17. Expression Language Injection - Unvalidated input in JSP/JSF/Spring EL, OGNL, SpEL

SERVER-SIDE ATTACKS:
18. SSRF (Server-Side Request Forgery) - URL fetching with user-controlled destinations

INPUT VALIDATION:
19. Regex DoS (ReDoS) - Nested quantifiers causing catastrophic backtracking
20. Log Injection - Unvalidated user input in log statements
21. Mass Assignment - Direct binding to object properties without validation

ADDITIONAL CONCERNS:
22. Insecure XML Processing - Unlimited entity expansion, XML bombs
23. Unvalidated Redirects - response.sendRedirect() with user input
24. JNI Security Issues - Unchecked native method calls
25. Race Conditions & Concurrency - Check-then-act on shared resources, unsynchronized access

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
