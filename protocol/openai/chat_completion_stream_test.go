package openai_test

import (
	"encoding/json"
	"strings"
	"testing"

	openai "hegq.net/mogo/pkg/protocol/openai"
)

// TestParseSSEEvent_PreservesRawBytes 验证上游 chunk 的 omitempty 零值字段
// 在 ParseSSEEvent → MarshalJSON 往返后仍保留（passthrough 保真）。
//
// 修复背景:OpenAIChatCompletionStreamEventUnion.UnmarshalJSON 之前对 chunk
// 做二次 json.Marshal,把 reasoning_content:"" / finish_reason:null /
// logprobs:null / token_ids:null / prompt_token_ids:null / stop_reason:null
// 等字段全部剥离,导致前端 SSE 与上游原始流差异显著。
func TestParseSSEEvent_PreservesRawBytes(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		mustContainAll []string // MarshalJSON 后必须出现的字段
	}{
		{
			name:  "first chunk 含 reasoning_content 与 prompt_token_ids",
			input: `data: {"choices":[{"delta":{"content":"","reasoning_content":"","role":"assistant"},"finish_reason":null,"index":0,"logprobs":null}],"created":1786331614,"id":"chatcmpl-test","model":"Qwen35","object":"chat.completion.chunk","prompt_token_ids":null}`,
			mustContainAll: []string{
				`"reasoning_content":""`,
				`"prompt_token_ids":null`,
				`"logprobs":null`,
				`"finish_reason":null`,
			},
		},
		{
			name:  "中间 chunk 含 finish_reason:null 与 token_ids:null",
			input: `data: {"choices":[{"delta":{"content":"你好","reasoning_content":""},"finish_reason":null,"index":0,"logprobs":null,"token_ids":null}],"created":1786331614,"id":"chatcmpl-test","model":"Qwen35","object":"chat.completion.chunk"}`,
			mustContainAll: []string{
				`"reasoning_content":""`,
				`"finish_reason":null`,
				`"logprobs":null`,
				`"token_ids":null`,
			},
		},
		{
			name:  "末尾 chunk 含 stop_reason",
			input: `data: {"choices":[{"delta":{"content":"","reasoning_content":""},"finish_reason":"stop","index":0,"logprobs":null,"stop_reason":null,"token_ids":null}],"created":1786331614,"id":"chatcmpl-test","model":"Qwen35","object":"chat.completion.chunk"}`,
			mustContainAll: []string{
				`"stop_reason":null`,
				`"finish_reason":"stop"`,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evt, err := openai.ParseSSEEvent(tc.input)
			if err != nil {
				t.Fatalf("ParseSSEEvent: %v", err)
			}
			if evt == nil {
				t.Fatal("nil event")
			}
			out, err := json.Marshal(evt)
			if err != nil {
				t.Fatalf("MarshalJSON: %v", err)
			}
			for _, must := range tc.mustContainAll {
				if !strings.Contains(string(out), must) {
					t.Errorf("missing field %q in marshaled output: %s", must, string(out))
				}
			}
		})
	}
}

// TestParseSSEEvent_RawChunkIsOriginal 验证 evt.Chunk 字段就是上游 data: 后
// 的原始字节（不做任何二次 marshal/格式化）。
func TestParseSSEEvent_RawChunkIsOriginal(t *testing.T) {
	raw := `{"choices":[{"delta":{"content":"你好","reasoning_content":""},"finish_reason":null,"index":0,"logprobs":null,"token_ids":null}],"created":1786331614,"id":"chatcmpl-test","model":"Qwen35","object":"chat.completion.chunk"}`
	input := "data: " + raw
	evt, err := openai.ParseSSEEvent(input)
	if err != nil {
		t.Fatalf("ParseSSEEvent: %v", err)
	}
	if evt == nil {
		t.Fatal("nil event")
	}
	if string(evt.Chunk) != raw {
		t.Errorf("Chunk mismatch:\n got: %s\nwant: %s", string(evt.Chunk), raw)
	}
}

// TestParseSSEEvent_DoneSentinel 验证 [DONE] 哨兵返回空 Type（与既有契约一致）。
func TestParseSSEEvent_DoneSentinel(t *testing.T) {
	evt, err := openai.ParseSSEEvent("data: [DONE]")
	if err != nil {
		t.Fatalf("ParseSSEEvent: %v", err)
	}
	if evt == nil {
		t.Fatal("nil event")
	}
	if evt.Type != "" {
		t.Errorf("DONE sentinel should have empty Type, got %q", evt.Type)
	}
	if len(evt.Chunk) != 0 {
		t.Errorf("DONE sentinel should have empty Chunk, got %q", string(evt.Chunk))
	}
}

// TestParseSSEEvent_NonChunkObjectIgnored 验证非 chat.completion.chunk 的 object
// 字段被忽略（不当作 chunk 返回）。
func TestParseSSEEvent_NonChunkObjectIgnored(t *testing.T) {
	evt, err := openai.ParseSSEEvent(`data: {"object":"not.a.chunk","x":1}`)
	if err != nil {
		t.Fatalf("ParseSSEEvent: %v", err)
	}
	if evt != nil {
		t.Errorf("non-chunk object should return nil, got %+v", evt)
	}
}
