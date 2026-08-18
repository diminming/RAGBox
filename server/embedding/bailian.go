package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	URL_API_BAILIAN    = "https://llm-qmh1u0iznlorduyx.cn-beijing.maas.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	MODEL_NAME_BAILIAN = "qwen3.7-text-embedding"
	API_KEY_BAILIAN    = "sk-896400710cc2412e9f6f1afe1447227b"
)

type BailianEmbedding struct {
	EmbeddingModel
}

func (embedding *BailianEmbedding) GetEmbedding(text string) ([]float32, error) {
	payload := map[string]any{
		"model": MODEL_NAME_BAILIAN,
		"input": map[string]any{
			"texts": []string{text},
		},
		"dimension": 1536, // 2560、2048、1536、1024、768、512、256；修改纬度变量时，milvus collection的schema也需要修改
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

	client := &http.Client{Timeout: 30 * time.Second}
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

	outcome := EmbeddingResponse{}
	err = json.Unmarshal(respBody, &outcome)
	if err != nil {
		return nil, err
	}

	if len(outcome.Output.Embeddings) != 1 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	vectors := outcome.Output.Embeddings[0].Embedding
	return vectors, nil
}
