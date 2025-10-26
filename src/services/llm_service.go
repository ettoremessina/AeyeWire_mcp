package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// LLMService handles communication with LMStudio
type LLMService struct {
	baseURL string
	model   string
	apiKey  string
	client  *http.Client
}

// LLMRequest represents a request to the LLM
type LLMRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse represents the response from LLM
type LLMResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewLLMService creates a new LLM service with configuration
func NewLLMService() *LLMService {
	baseURL := os.Getenv("LMSTUDIO_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:1234"
	}

	model := os.Getenv("LMSTUDIO_MODEL")
	if model == "" {
		model = "qwen/qwen3-coder-30b"
	}

	apiKey := os.Getenv("LMSTUDIO_API_KEY")

	return &LLMService{
		baseURL: baseURL,
		model:   model,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Analyze sends code to LLM for security analysis
func (llm *LLMService) Analyze(code string, prompt string) (string, error) {
	messages := []Message{
		{
			Role:    "system",
			Content: "You are a security analysis expert. Analyze the provided code and return findings in JSON format.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("%s\n\nCode to analyze:\n```\n%s\n```", prompt, code),
		},
	}

	request := LLMRequest{
		Model:       llm.model,
		Messages:    messages,
		Temperature: 0.1, // Low temperature for consistent analysis
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", llm.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if llm.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", llm.apiKey))
	}

	resp, err := llm.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM service returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var llmResp LLMResponse
	if err := json.Unmarshal(bodyBytes, &llmResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(llmResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in LLM response")
	}

	return llmResp.Choices[0].Message.Content, nil
}

// HealthCheck verifies LLM service availability
func (llm *LLMService) HealthCheck() (bool, error) {
	url := fmt.Sprintf("%s/v1/models", llm.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	if llm.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", llm.apiKey))
	}

	resp, err := llm.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
