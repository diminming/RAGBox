package rerank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	URL_API_BAILIAN    = "https://llm-qmh1u0iznlorduyx.cn-beijing.maas.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank"
	MODEL_NAME_BAILIAN = "qwen3-rerank"
	API_KEY_BAILIAN    = "sk-896400710cc2412e9f6f1afe1447227b"
)

type BailianRerank struct {
	RerankModel
}

func (rerank *BailianRerank) Rerank(query string, textLst []string) ([]string, error) {
	payload := map[string]any{
		"model": MODEL_NAME_BAILIAN,
		"input": map[string]any{
			"query":     query,
			"documents": textLst,
		},
		"parameters": map[string]any{
			"return_documents": true,
			"top_n":            5,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_API_BAILIAN, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_KEY_BAILIAN)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	outcome := RerankResponse{}
	err = json.Unmarshal(respBody, &outcome)
	if err != nil {
		return nil, err
	}

	if len(outcome.Output.Results) == 0 {
		return nil, fmt.Errorf("no reranked documents returned")
	}

	results := make([]string, len(outcome.Output.Results))
	for i, doc := range outcome.Output.Results {
		results[i] = fmt.Sprintf("DocIdx: %d, Score: %f, Text: %s", doc.DocIdx, doc.Score, doc.Document.Text)
	}
	return results, nil
}
