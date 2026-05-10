package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequestBody struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResponseBody struct {
	Message ChatMessage `json:"message"`
}

func Chat(endpoint string, apiKey string, model string, prompt string, temp float32, maxTokens int) string {
	reqBody, err := json.Marshal(ChatRequestBody{
		Model: model,
		Messages: []ChatMessage{{
			Role:    "user",
			Content: prompt,
		}},
		Stream: false,
	})
	if err != nil {
		return err.Error()
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err.Error()
	}

	req.Header.Add("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var respBody ChatResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return err.Error()
		}
		return respBody.Message.Content
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("%s", resp.Status)
	}
	return strings.TrimSpace(fmt.Sprintf("%s: %s", resp.Status, body))
}
