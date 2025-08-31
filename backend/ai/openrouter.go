package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenRouterClient struct {
	APIKey string
	Client *http.Client
}

func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		APIKey: apiKey,
		Client: &http.Client{},
	}
}

func (c *OpenRouterClient) Chat(ctx context.Context, text string) (string, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	body := map[string]interface{}{
		"model": "openai/gpt-4o",
		"messages": []map[string]string{
			{"role": "user", "content": text},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	// Debug log
	fmt.Println("Raw response:", string(respBytes))

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", nil
	}

	return parsed.Choices[0].Message.Content, nil
}
