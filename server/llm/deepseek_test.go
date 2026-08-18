package llm_test

import (
	"encoding/json"
	"fmt"
	"ragbox/llm"
	"testing"
)

func TestSend2LLM(t *testing.T) {
	content := `{
    "data":
    {
        "choices":
        [
            {
                "finish_reason": "stop",
                "index": 0,
                "logprobs": null,
                "message":
                {
                    "content": "你好！我是DeepSeek，一个由深度求索公司开发的AI助手。我旨在通过对话帮助你解决问题、提供信息、激发灵感，或者在你需要的时候陪聊。无论是学习、工作、生活中的疑问，还是想聊聊天，我都非常乐意提供支持！有任何问题，尽管问我吧！😊",
                    "role": "assistant"
                }
            }
        ],
        "created": 1786005777,
        "id": "db8be1ea-43ae-45ef-9410-98a53e825462",
        "model": "deepseek-v4-flash",
        "object": "chat.completion",
        "system_fingerprint": "fp_a18b46594c_prod0820_fp8_kvcache_20260402",
        "usage":
        {
            "completion_tokens": 65,
            "prompt_cache_hit_tokens": 0,
            "prompt_cache_miss_tokens": 10,
            "prompt_tokens": 10,
            "prompt_tokens_details":
            {
                "cached_tokens": 0
            },
            "total_tokens": 75
        }
    },
    "msg": "success"
}`

	var resp llm.LLMResponse
	err := json.Unmarshal([]byte(content), &resp)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}

func TestDeepSeekDirectResponseUnmarshal(t *testing.T) {
	content := `{
  "id": "demo-id",
  "object": "chat.completion",
  "created": 1786005777,
  "model": "deepseek-v4-flash",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "direct-format"
      },
      "finish_reason": "stop"
    }
  ]
}`

	var resp llm.LLMResponse
	err := json.Unmarshal([]byte(content), &resp)
	if err != nil {
		t.Fatalf("Failed to unmarshal direct format JSON: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice in top-level choices, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "direct-format" {
		t.Fatalf("unexpected content: %s", resp.Choices[0].Message.Content)
	}
}
