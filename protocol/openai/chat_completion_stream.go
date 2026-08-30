package openai

import (
	"encoding/json"
	"strings"
)

// ============== ChatCompletion SSE 事件类型常量 ==============

// OpenAIChatCompletionStreamChunkObject SSE chunk 事件的 object 字段值。
const OpenAIChatCompletionStreamChunkObject = "chat.completion.chunk"

// ============== OpenAIChatCompletionStreamEvent union ==============

// OpenAIChatCompletionStreamEvent Chat Completion SSE 事件 union。
//
// 两种事件：
//   - chunk：data: {"id":"...","object":"chat.completion.chunk","choices":[...]}
//   - done：data: [DONE]
//
// 通过 isChatCompletionStreamEvent() 私有方法封闭实现集。
type OpenAIChatCompletionStreamEvent interface {
	isChatCompletionStreamEvent()
}

// OpenAIChatCompletionStreamChunkEvent chunk 事件，对齐 OpenAI 官方 SSE 格式。
type OpenAIChatCompletionStreamChunkEvent struct {
	ID                string                                     `json:"id"`
	Object            string                                     `json:"object"`
	Created           int64                                      `json:"created"`
	Model             string                                     `json:"model"`
	Choices           []OpenAIChatCompletionStreamChunkChoice    `json:"choices"`
	Usage             *OpenAICompletionUsage                     `json:"usage,omitempty"`
	SystemFingerprint string                                     `json:"system_fingerprint,omitempty"`
	ServiceTier       string                                     `json:"service_tier,omitempty"`
	Moderation        *OpenAIChatCompletionStreamChunkModeration `json:"moderation,omitempty"`
}

func (OpenAIChatCompletionStreamChunkEvent) isChatCompletionStreamEvent() {}

// OpenAIChatCompletionStreamEventUnion 扁平实现（解决 interface unmarshal 问题）。
type OpenAIChatCompletionStreamEventUnion struct {
	Type  string          `json:"type"`
	Chunk json.RawMessage `json:"chunk,omitempty"` // OpenAIChatCompletionStreamChunkEvent
}

func (OpenAIChatCompletionStreamEventUnion) isChatCompletionStreamEvent() {}

func (u OpenAIChatCompletionStreamEventUnion) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	if u.Chunk == nil {
		return nil, nil
	}
	return u.Chunk, nil
}

func (u *OpenAIChatCompletionStreamEventUnion) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var raw struct {
		Object string `json:"object"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	u.Type = raw.Object

	// chunk 事件:保留上游原始字节,不做二次 marshal。
	//
	// 设计意图:Union.MarshalJSON 直接返回 u.Chunk;UnmarshalJSON 只要把上游 data
	// 原样放入 u.Chunk,下游 WriteStream 透传出去就是 byte-for-byte 保真。
	//
	// 此前这里用 json.Unmarshal→json.Marshal 往返,会把 reasoning_content:"" /
	// finish_reason:null / logprobs:null / token_ids:null / prompt_token_ids:null /
	// stop_reason:null 等 omitempty 零值字段全部剥离,导致客户端拿到的 SSE 与上游
	// 原始流差异显著、出现 "Null" 显示问题。
	if raw.Object == OpenAIChatCompletionStreamChunkObject {
		u.Chunk = data
	}
	return nil
}

// OpenAIChatCompletionStreamChunkChoice SSE chunk 的 choice。
type OpenAIChatCompletionStreamChunkChoice struct {
	Index        int64                                          `json:"index"`
	ID           string                                         `json:"id,omitempty"`
	Delta        OpenAIChatCompletionStreamChunkChoiceDelta     `json:"delta"`
	FinishReason *string                                        `json:"finish_reason,omitempty"`
	Logprobs     *OpenAIChatCompletionStreamChunkChoiceLogprobs `json:"logprobs,omitempty"`
}

// OpenAIChatCompletionStreamChunkChoiceDelta SSE chunk 的 delta。
type OpenAIChatCompletionStreamChunkChoiceDelta struct {
	Role             string                                               `json:"role,omitempty"`
	Content          string                                               `json:"content,omitempty"`
	Refusal          *string                                              `json:"refusal,omitempty"`
	ToolCalls        []OpenAIChatCompletionStreamChunkChoiceDeltaToolCall `json:"tool_calls,omitempty"`
	Audio            *OpenAIChatCompletionStreamChunkChoiceDeltaAudio     `json:"audio,omitempty"`
	Reasoning        string                                               `json:"reasoning,omitempty"`
	ReasoningContent string                                               `json:"reasoning_content,omitempty"`
}

// OpenAIChatCompletionStreamChunkChoiceDeltaToolCall SSE chunk 的工具调用。
type OpenAIChatCompletionStreamChunkChoiceDeltaToolCall struct {
	Index    int64                                                       `json:"index"`
	ID       string                                                      `json:"id,omitempty"`
	Type     string                                                      `json:"type,omitempty"`
	Function *OpenAIChatCompletionStreamChunkChoiceDeltaToolCallFunction `json:"function,omitempty"`
}

// OpenAIChatCompletionStreamChunkChoiceDeltaToolCallFunction SSE chunk 的函数调用。
type OpenAIChatCompletionStreamChunkChoiceDeltaToolCallFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// OpenAIChatCompletionStreamChunkChoiceDeltaAudio SSE chunk 的音频输出。
type OpenAIChatCompletionStreamChunkChoiceDeltaAudio struct {
	ID         string `json:"id,omitempty"`
	Data       string `json:"data,omitempty"`
	Transcript string `json:"transcript,omitempty"`
	ExpiresAt  int64  `json:"expires_at,omitempty"`
}

// OpenAIChatCompletionStreamChunkChoiceLogprobs chunk choice 的 logprobs。
type OpenAIChatCompletionStreamChunkChoiceLogprobs struct {
	Content      []OpenAIChatCompletionStreamChunkChoiceLogprobsContent `json:"content,omitempty"`
	Role         []string                                               `json:"role,omitempty"`
	FunctionCall []any                                                  `json:"function_call,omitempty"`
	ToolCall     []any                                                  `json:"tool_call,omitempty"`
}

// OpenAIChatCompletionStreamChunkChoiceLogprobsContent logprobs content 项。
type OpenAIChatCompletionStreamChunkChoiceLogprobsContent struct {
	Token        string  `json:"token,omitempty"`
	LogProb      float64 `json:"logprob,omitempty"`
	Prob         float64 `json:"prob,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
	TopLogProbs  []struct {
		Token   string  `json:"token,omitempty"`
		LogProb float64 `json:"logprob,omitempty"`
	} `json:"top_logprobs,omitempty"`
}

// OpenAIChatCompletionStreamChunkModeration chunk 的 moderation 字段。
type OpenAIChatCompletionStreamChunkModeration struct {
	ID      string `json:"id,omitempty"`
	Output  string `json:"output,omitempty"`
	Version string `json:"version,omitempty"`
}

// ============== ParseSSEEvent ==============

// ParseSSEEvent 解析 Chat Completion SSE 行（"data: {...}" 或 "data: [DONE]"）为 OpenAIChatCompletionStreamEvent。
//
// 在 pkg 内部实现，可以返回 sealed interface 的具体实现。
// 供 adapter 层复用，不需要自己定义 SSE 行 struct。
//
// 返回值类型为 OpenAIChatCompletionStreamEventUnion（扁平实现），
// 供 protocol 层直接做 marshal 输出。
func ParseSSEEvent(line string) (*OpenAIChatCompletionStreamEventUnion, error) {
	data := line
	hasPrefix := strings.HasPrefix(line, "data: ")
	if hasPrefix {
		data = strings.TrimPrefix(line, "data: ")
	}
	data = strings.TrimSpace(data)
	if data == "" {
		return nil, nil
	}

	// [DONE] 结束信号：返回空 union，Type 为空
	if data == "[DONE]" {
		return &OpenAIChatCompletionStreamEventUnion{}, nil
	}

	if !hasPrefix {
		return nil, nil
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, err
	}

	obj, _ := raw["object"].(string)
	if obj != OpenAIChatCompletionStreamChunkObject {
		return nil, nil
	}

	// chunk 事件:把上游 data 原始字节装进 Chunk,不做任何二次 marshal。
	//
	// ParseSSEEvent 自身也是 passthrough 实现——不做 json.Unmarshal→Marshal 往返,
	// 避免 omitempty 零值字段(reasoning_content:"" / finish_reason:null /
	// logprobs:null / token_ids:null / prompt_token_ids:null / stop_reason:null)
	// 被剥离。obj 已在上方 raw["object"] 处确认是 "chat.completion.chunk",
	// 这里无需再 unmarshal 一遍 chunk struct。
	return &OpenAIChatCompletionStreamEventUnion{
		Type:  obj,
		Chunk: json.RawMessage(data),
	}, nil
}
