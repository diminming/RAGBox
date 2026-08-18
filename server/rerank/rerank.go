package rerank

type RerankModel interface {
	Rerank(query string, candidates []string) ([]string, error)
}

type RerankResponse struct {
	Output struct {
		Results []struct {
			Document struct {
				Text string `json:"text"`
			} `json:"document"`
			DocIdx int     `json:"index"`
			Score  float32 `json:"relevance_score"`
		} `json:"results"`
	} `json:"output"`

	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`

	RequestId string `json:"request_id"`
}
