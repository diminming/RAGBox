package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ragbox/config"
	"time"
)

var (
	URL_API_BAILIAN    = config.Config.Model.Embedding.APIURL
	MODEL_NAME_BAILIAN = config.Config.Model.Embedding.ModelName
	API_KEY_BAILIAN    = config.Config.Model.Embedding.APIKey
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
