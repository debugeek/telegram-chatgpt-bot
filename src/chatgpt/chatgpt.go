package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice `json:"choices"`
}

type ChatCompletionChoice struct {
	Message ChatCompletionMessage `json:"message"`
}

type ChatCompetitionResponseError struct {
	Error ChatCompetitionResponseErrorInfo `json:"error"`
}

type ChatCompetitionResponseErrorInfo struct {
	Message string `json:"message"`
}

func SendText(text string, apikey string, model string) string {
	message := ChatCompletionMessage{
		Role:    "user",
		Content: text,
	}
	request := ChatCompletionRequest{
		Model:    model,
		Messages: []ChatCompletionMessage{message},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return err.Error()
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return err.Error()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
	req.Header.Add("Content-Type", "application/json")

	http := &http.Client{}
	resp, err := http.Do(req)
	if err != nil {
		return err.Error()
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		response := ChatCompletionResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return err.Error()
		}
		return response.Choices[0].Message.Content
	} else if resp.StatusCode >= 400 {
		error := ChatCompetitionResponseError{}
		if err := json.NewDecoder(resp.Body).Decode(&error); err != nil {
			return err.Error()
		}
		return error.Error.Message
	} else {
		return fmt.Sprintf("%v", resp)
	}
}
