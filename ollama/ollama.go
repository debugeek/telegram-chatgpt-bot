package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestBody struct {
	Model     string  `json:"model"`
	Prompt    string  `json:"prompt"`
	Temp      float32 `json:"temperature"`
	MaxTokens int     `json:"max_tokens"`
	Stream    bool    `json:"stream"`
}

type ResponseBody struct {
	Response string `json:"response"`
}

func Chat(endpoint string, model string, prompt string, temp float32, maxTokens int) string {
	reqBody, err := json.Marshal(RequestBody{
		Model:     model,
		Prompt:    prompt,
		Temp:      temp,
		MaxTokens: maxTokens,
		Stream:    false,
	})
	if err != nil {
		return err.Error()
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err.Error()
	}

	req.Header.Add("Content-Type", "application/json")

	http := &http.Client{}
	resp, err := http.Do(req)
	if err != nil {
		return err.Error()
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var respBody ResponseBody
		if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return err.Error()
		}
		return respBody.Response
	} else {
		return fmt.Sprintf("%v", resp)
	}
}
