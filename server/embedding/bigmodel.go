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
	API_KEY_BIGMODEL    = "8a34011840e7457f8f8e6a61410ac8b1.oqXualgmrDHMqOCd"
	URL_API_BIGMODEL    = "https://open.bigmodel.cn/api/paas/v4/embeddings"
	MODEL_NAME_BIGMODEL = "embedding-3"
)

type BigmodelEmbedding struct {
	EmbeddingModel
}

func (embedding *BigmodelEmbedding) GetEmbedding(text string) ([]float32, error) {
	payload := map[string]any{
		"model": MODEL_NAME_BIGMODEL,
		"input": map[string]any{
			"texts": []string{text},
		},
		"dimensions": 1024, // 2048、1024、512、256；修改维度变量时，milvus collection的schema也需要修改
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_API_BIGMODEL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+API_KEY_BIGMODEL)

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
