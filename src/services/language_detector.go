package services

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/emware/aeyewire-mcp/src/models"
)

// LanguageDetector detects programming languages from code and file extensions
type LanguageDetector struct {
	patterns map[models.LanguageType][]*regexp.Regexp
}

// NewLanguageDetector creates a new language detector with initialized patterns
func NewLanguageDetector() *LanguageDetector {
	detector := &LanguageDetector{
		patterns: make(map[models.LanguageType][]*regexp.Regexp),
	}
	detector.initializePatterns()
	return detector
}

// initializePatterns sets up regex patterns for language detection
func (ld *LanguageDetector) initializePatterns() {
	// Java patterns
	ld.patterns[models.JAVA] = []*regexp.Regexp{
		regexp.MustCompile(`package\s+[a-z][a-z0-9_.]*;`),
		regexp.MustCompile(`import\s+(java|javax|org)\.`),
		regexp.MustCompile(`public\s+(class|interface|enum)\s+\w+`),
		regexp.MustCompile(`@(Override|Autowired|Entity|Controller|Service|Repository)`),
		regexp.MustCompile(`(public|private|protected)\s+(static\s+)?(void|int|String|boolean)`),
	}

	// C# patterns
	ld.patterns[models.CSHARP] = []*regexp.Regexp{
		regexp.MustCompile(`using\s+System`),
		regexp.MustCompile(`namespace\s+\w+`),
		regexp.MustCompile(`(public|private|protected)\s+(class|interface|struct)\s+\w+`),
		regexp.MustCompile(`\[assembly:\s*\w+`),
		regexp.MustCompile(`(async\s+)?Task<`),
	}

	// React TypeScript patterns
	ld.patterns[models.REACT_TYPESCRIPT] = []*regexp.Regexp{
		regexp.MustCompile(`import\s+.*from\s+['"]react['"]`),
		regexp.MustCompile(`interface\s+\w+\s*{`),
		regexp.MustCompile(`type\s+\w+\s*=`),
		regexp.MustCompile(`:\s*(React\.)?FC<`),
		regexp.MustCompile(`<\w+[^>]*>`), // JSX/TSX tags
	}

	// React JavaScript patterns
	ld.patterns[models.REACT_JAVASCRIPT] = []*regexp.Regexp{
		regexp.MustCompile(`import\s+.*from\s+['"]react['"]`),
		regexp.MustCompile(`React\.createElement`),
		regexp.MustCompile(`useState|useEffect|useCallback|useMemo`),
		regexp.MustCompile(`<\w+[^>]*>`), // JSX tags
		regexp.MustCompile(`export\s+(default\s+)?(function|const)\s+\w+`),
	}
}

// DetectFromExtension detects language based on file extension
func (ld *LanguageDetector) DetectFromExtension(filePath string) models.LanguageType {
	if filePath == "" {
		return models.UNKNOWN
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".java":
		return models.JAVA
	case ".cs":
		return models.CSHARP
	case ".tsx":
		return models.REACT_TYPESCRIPT
	case ".ts":
		// Could be React or plain TypeScript, will need content analysis
		return models.REACT_TYPESCRIPT
	case ".jsx":
		return models.REACT_JAVASCRIPT
	case ".js":
		// Could be React or plain JavaScript, will need content analysis
		return models.REACT_JAVASCRIPT
	default:
		return models.UNKNOWN
	}
}

// DetectFromContent detects language based on code content using pattern matching
func (ld *LanguageDetector) DetectFromContent(code string) models.LanguageType {
	scores := make(map[models.LanguageType]int)

	// Check each language's patterns
	for lang, patterns := range ld.patterns {
		for _, pattern := range patterns {
			if pattern.MatchString(code) {
				scores[lang]++
			}
		}
	}

	// Return language with highest score
	maxScore := 0
	detectedLang := models.UNKNOWN

	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			detectedLang = lang
		}
	}

	return detectedLang
}

// Detect detects language using extension first, then falls back to content analysis
func (ld *LanguageDetector) Detect(code string, filePath string) models.LanguageType {
	// Try extension-based detection first
	if filePath != "" {
		lang := ld.DetectFromExtension(filePath)
		if lang != models.UNKNOWN {
			// For .ts and .js files, verify it's actually React code
			if lang == models.REACT_TYPESCRIPT || lang == models.REACT_JAVASCRIPT {
				// Check if code contains React patterns
				contentLang := ld.DetectFromContent(code)
				if contentLang == lang {
					return lang
				}
			} else {
				return lang
			}
		}
	}

	// Fall back to content-based detection
	return ld.DetectFromContent(code)
}

// GetSupportedLanguages returns a list of all supported languages with metadata
func (ld *LanguageDetector) GetSupportedLanguages() []models.LanguageInfo {
	return []models.LanguageInfo{
		{
			Identifier:  string(models.CSHARP),
			Description: "C# programming language",
			Extensions:  []string{".cs"},
		},
		{
			Identifier:  string(models.REACT_TYPESCRIPT),
			Description: "React with TypeScript",
			Extensions:  []string{".tsx", ".ts"},
		},
		{
			Identifier:  string(models.REACT_JAVASCRIPT),
			Description: "React with JavaScript",
			Extensions:  []string{".jsx", ".js"},
		},
		{
			Identifier:  string(models.JAVA),
			Description: "Java programming language",
			Extensions:  []string{".java"},
		},
	}
}
