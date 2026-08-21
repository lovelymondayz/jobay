package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// AIClient handles calls to 9router (OpenRouter)
type AIClient struct {
	baseURL string
	apiKey  string
	model   string
}

func NewAIClient() *AIClient {
	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://9router.nousresearch.com/v1"
	}
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}
	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "poolside/laguna-m.1:free"
	}

	return &AIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
	}
}

// ChatMessage represents a single message in a chat completion
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the request body for chat completions
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse is the response from chat completions
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// Chat sends a chat completion request
func (c *AIClient) Chat(systemPrompt, userPrompt string) (string, error) {
	req := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI returned status %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in AI response")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

// ChatJSON sends a chat completion and parses the response as JSON into target
func (c *AIClient) ChatJSON(systemPrompt, userPrompt string, target interface{}) error {
	response, err := c.Chat(systemPrompt, userPrompt)
	if err != nil {
		return err
	}

	// Try to extract JSON from markdown code block
	if idx := strings.Index(response, "```"); idx != -1 {
		end := strings.Index(response[idx+3:], "```")
		if end != -1 {
			response = response[idx+3 : idx+3+end]
			// Remove language hint
			if nl := strings.Index(response, "\n"); nl != -1 {
				response = response[nl+1:]
			}
		}
	}

	if err := json.Unmarshal([]byte(response), target); err != nil {
		return fmt.Errorf("parse AI JSON response: %w (raw: %s)", err, response)
	}
	return nil
}
