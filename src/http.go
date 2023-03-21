package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTTPClient struct {
	baseURL string
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int          `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatResponseUsage      `json:"usage"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

type ChatResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatResponseError struct {
	Error ChatResponseErrorInfo `json:"error"`
}

type ChatResponseErrorInfo struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"`
	Code    interface{} `json:"code"`
}

func (c *HTTPClient) buildURL(path string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	u.Path = strings.Join([]string{
		strings.TrimRight(u.Path, "/"),
		strings.TrimLeft(path, "/"),
	}, "/")
	return u.String(), nil
}

func (c *HTTPClient) build(ctx context.Context, path string, method string, headers map[string]string, params interface{}) (req *http.Request, err error) {
	url, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	reader, contentType, err := c.buildReader(params)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

func (c *HTTPClient) buildReader(params interface{}) (io.Reader, string, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, "", err
	}
	reader := bytes.NewBuffer(b)
	return reader, "application/json", nil
}

func (c *HTTPClient) request(ctx context.Context, path string, method string, headers map[string]string, params interface{}, result interface{}) error {
	req, err := c.build(ctx, path, method, headers, params)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
		return nil
	} else if resp.StatusCode >= 400 {
		err := ChatResponseError{}
		if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
			return err
		}
		return fmt.Errorf("%s", err.Error.Message)
	} else {
		return fmt.Errorf("%v", resp)
	}
}

func request[T any](c *HTTPClient, ctx context.Context, path string, method string, headers map[string]string, params interface{}, result T) (T, error) {
	err := c.request(ctx, path, method, headers, params, &result)
	return result, err
}

func (c *HTTPClient) SendChatMessage(message string, key string, model string) (resp ChatCompletionResponse, err error) {
	return request(c,
		context.Background(),
		"chat/completions",
		"POST",
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", key),
		},
		ChatCompletionRequest{
			Model: model,
			Messages: []ChatCompletionMessage{
				{Role: "user", Content: message},
			},
		}, resp)
}
