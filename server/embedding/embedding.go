package embedding

type EmbeddingModel interface {
	GetEmbedding(text string) ([]float32, error)
}

type EmbeddingResponse struct {
	Output struct {
		Embeddings []struct {
			TxtIdx    int       `json:"text_index"`
			Embedding []float32 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`

	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`

	RequestId string `json:"request_id"`
}
