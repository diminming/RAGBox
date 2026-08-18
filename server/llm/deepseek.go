package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"ragbox/config"
	"time"
)

var (
	API_KEY = config.Config.Model.LLM.APIKey
	URL     = config.Config.Model.LLM.APIURL
	METHOD  = "POST"
)

type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`

	// Data struct {
	// 	Choices []struct {
	// 		Message struct {
	// 			Content string `json:"content"`
	// 			Role    string `json:"role"`
	// 		} `json:"message"`
	// 		FinishReason string `json:"finish_reason"`
	// 		Index        int    `json:"index"`
	// 	} `json:"choices"`
	// } `json:"data"`

	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`

	// Msg string `json:"msg"`
}

func Send2LLM(query string) (*LLMResponse, error) {
	type chatRequest struct {
		Messages []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"messages"`
		Model           string `json:"model"`
		ReasoningEffort string `json:"reasoning_effort"`
		MaxTokens       int    `json:"max_tokens"`
		Temperature     int    `json:"temperature"`
		TopP            int    `json:"top_p"`
		Stream          bool   `json:"stream"`
	}

	payload := chatRequest{
		Messages: []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		}{
			{Content: "你是一个有帮助的助手", Role: "system"},
			{Content: query, Role: "user"},
		},
		Model:           config.Config.Model.LLM.ModelName,
		ReasoningEffort: "low",
		MaxTokens:       4096,
		Temperature:     1,
		TopP:            1,
		Stream:          false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error("request marshal error.", "err", err)
		return nil, err
	}
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(METHOD, URL, bytes.NewReader(body))
	if err != nil {
		slog.Error("request creation error.", "err", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	res, err := client.Do(req)
	if err != nil {
		slog.Error("request execution error.", "err", err)
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("response read error.", "err", err)
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed: status=%d body=%s", res.StatusCode, string(respBody))
	}

	resp := new(LLMResponse)
	err = json.Unmarshal(respBody, resp)
	if err != nil {
		slog.Error("response unmarshal error.", "err", err)
		return nil, err
	}
	return resp, nil
}
