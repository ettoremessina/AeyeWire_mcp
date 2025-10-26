package services

import (
	"testing"

	"github.com/emware/aeyewire-mcp/src/models"
)

func TestDetectFromExtension(t *testing.T) {
	detector := NewLanguageDetector()

	tests := []struct {
		name     string
		filePath string
		expected models.LanguageType
	}{
		{"Java file", "Example.java", models.JAVA},
		{"C# file", "Example.cs", models.CSHARP},
		{"TypeScript React file", "Component.tsx", models.REACT_TYPESCRIPT},
		{"JavaScript React file", "Component.jsx", models.REACT_JAVASCRIPT},
		{"TypeScript file", "utils.ts", models.REACT_TYPESCRIPT},
		{"JavaScript file", "utils.js", models.REACT_JAVASCRIPT},
		{"Unknown extension", "example.py", models.UNKNOWN},
		{"No extension", "README", models.UNKNOWN},
		{"Empty path", "", models.UNKNOWN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.DetectFromExtension(tt.filePath)
			if result != tt.expected {
				t.Errorf("DetectFromExtension(%s) = %v, want %v", tt.filePath, result, tt.expected)
			}
		})
	}
}

func TestDetectFromContent(t *testing.T) {
	detector := NewLanguageDetector()

	tests := []struct {
		name     string
		code     string
		expected models.LanguageType
	}{
		{
			name: "Java code",
			code: `package com.example;
import java.util.*;
public class Example {
    public static void main(String[] args) {}
}`,
			expected: models.JAVA,
		},
		{
			name: "C# code",
			code: `using System;
namespace Example {
    public class Test {
        public async Task<int> Method() { return 0; }
    }
}`,
			expected: models.CSHARP,
		},
		{
			name: "React TypeScript",
			code: `import React from 'react';
interface Props {
    name: string;
}
const Component: React.FC<Props> = ({ name }) => {
    return <div>{name}</div>;
};`,
			expected: models.REACT_TYPESCRIPT,
		},
		{
			name: "React JavaScript",
			code: `import React from 'react';
const Component = ({ name }) => {
    const [state, setState] = useState('');
    return <div>{name}</div>;
};`,
			expected: models.REACT_JAVASCRIPT,
		},
		{
			name:     "Unknown code",
			code:     `print("Hello World")`,
			expected: models.UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.DetectFromContent(tt.code)
			if result != tt.expected {
				t.Errorf("DetectFromContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	detector := NewLanguageDetector()

	tests := []struct {
		name     string
		code     string
		filePath string
		expected models.LanguageType
	}{
		{
			name:     "Java with extension",
			code:     "public class Test {}",
			filePath: "Test.java",
			expected: models.JAVA,
		},
		{
			name:     "C# with extension",
			code:     "public class Test {}",
			filePath: "Test.cs",
			expected: models.CSHARP,
		},
		{
			name: "React TypeScript with extension",
			code: `import React from 'react';
const Component: React.FC = () => <div />;`,
			filePath: "Component.tsx",
			expected: models.REACT_TYPESCRIPT,
		},
		{
			name:     "Java without extension",
			code:     "package com.example; public class Test {}",
			filePath: "",
			expected: models.JAVA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.Detect(tt.code, tt.filePath)
			if result != tt.expected {
				t.Errorf("Detect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	detector := NewLanguageDetector()
	languages := detector.GetSupportedLanguages()

	if len(languages) != 4 {
		t.Errorf("Expected 4 supported languages, got %d", len(languages))
	}

	expectedLanguages := map[string]bool{
		"csharp":             false,
		"react_typescript":   false,
		"react_javascript":   false,
		"java":               false,
	}

	for _, lang := range languages {
		if _, ok := expectedLanguages[lang.Identifier]; ok {
			expectedLanguages[lang.Identifier] = true
		}
	}

	for lang, found := range expectedLanguages {
		if !found {
			t.Errorf("Expected language %s not found in supported languages", lang)
		}
	}
}
