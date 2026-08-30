package anthropic

import "encoding/json"

// ============== RawContentBlockDelta union ==============

// RawContentBlockDelta content_block_delta.delta 字段的 union。
//
// 5 种 delta 类型，承载文本块 / thinking 块 / 工具输入 / signature / 引用位置的增量。
// 兼容旧 sealed interface；RawContentBlockDeltaUnion 为扁平实现（解决 unmarshal 问题）。
type RawContentBlockDelta interface {
	isRawContentBlockDelta()
}

type AnthropicTextDeltaBlockType string

const AnthropicTextDeltaBlockTypeTextDelta AnthropicTextDeltaBlockType = "text_delta"

// AnthropicTextDeltaBlock 文本块增量。
type AnthropicTextDeltaBlock struct {
	Type AnthropicTextDeltaBlockType `json:"type"` // 固定 "text_delta"
	Text string                      `json:"text"`
}

func (AnthropicTextDeltaBlock) isRawContentBlockDelta() {}

type AnthropicThinkingDeltaBlockType string

const AnthropicThinkingDeltaBlockTypeThinkingDelta AnthropicThinkingDeltaBlockType = "thinking_delta"

// AnthropicThinkingDeltaBlock thinking 块增量。
type AnthropicThinkingDeltaBlock struct {
	Type     AnthropicThinkingDeltaBlockType `json:"type"` // 固定 "thinking_delta"
	Thinking string                          `json:"thinking"`
}

func (AnthropicThinkingDeltaBlock) isRawContentBlockDelta() {}

type AnthropicInputJSONDeltaBlockType string

const AnthropicInputJSONDeltaBlockTypeInputJSONDelta AnthropicInputJSONDeltaBlockType = "input_json_delta"

// AnthropicInputJSONDeltaBlock 工具 input JSON 增量（流式拼接）。
type AnthropicInputJSONDeltaBlock struct {
	Type        AnthropicInputJSONDeltaBlockType `json:"type"` // 固定 "input_json_delta"
	PartialJSON string                           `json:"partial_json"`
}

func (AnthropicInputJSONDeltaBlock) isRawContentBlockDelta() {}

type AnthropicSignatureDeltaBlockType string

const AnthropicSignatureDeltaBlockTypeSignatureDelta AnthropicSignatureDeltaBlockType = "signature_delta"

// AnthropicSignatureDeltaBlock thinking 签名增量（用于审计回放）。
type AnthropicSignatureDeltaBlock struct {
	Type      AnthropicSignatureDeltaBlockType `json:"type"` // 固定 "signature_delta"
	Signature string                           `json:"signature"`
}

func (AnthropicSignatureDeltaBlock) isRawContentBlockDelta() {}

type AnthropicCitationsDeltaBlockType string

const AnthropicCitationsDeltaBlockTypeCitationsDelta AnthropicCitationsDeltaBlockType = "citations_delta"

// AnthropicCitationsDeltaBlock 引用位置增量。
//
// Citation 字段详细 5 种 TextCitation 子类型
// （CharLocation / PageLocation / ContentBlockLocation /
// SearchResultLocation / WebSearchResultLocation）由 P1 补齐，本轮用 RawMessage 占位。
type AnthropicCitationsDeltaBlock struct {
	Type     AnthropicCitationsDeltaBlockType `json:"type"` // 固定 "citations_delta"
	Citation json.RawMessage                  `json:"citation"`
}

func (AnthropicCitationsDeltaBlock) isRawContentBlockDelta() {}

// RawContentBlockDeltaUnion 扁平实现（解决 interface unmarshal 问题）。
type RawContentBlockDeltaUnion struct {
	Type           string          `json:"type"`
	TextDelta      json.RawMessage `json:"text_delta,omitempty"`
	ThinkingDelta  json.RawMessage `json:"thinking_delta,omitempty"`
	InputJSONDelta json.RawMessage `json:"input_json_delta,omitempty"`
	SignatureDelta json.RawMessage `json:"signature_delta,omitempty"`
	CitationsDelta json.RawMessage `json:"citations_delta,omitempty"`
}

func (RawContentBlockDeltaUnion) isRawContentBlockDelta() {}

func (u RawContentBlockDeltaUnion) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case string(AnthropicTextDeltaBlockTypeTextDelta):
		if u.TextDelta == nil {
			return nil, nil
		}
		return u.TextDelta, nil
	case string(AnthropicThinkingDeltaBlockTypeThinkingDelta):
		if u.ThinkingDelta == nil {
			return nil, nil
		}
		return u.ThinkingDelta, nil
	case string(AnthropicInputJSONDeltaBlockTypeInputJSONDelta):
		if u.InputJSONDelta == nil {
			return nil, nil
		}
		return u.InputJSONDelta, nil
	case string(AnthropicSignatureDeltaBlockTypeSignatureDelta):
		if u.SignatureDelta == nil {
			return nil, nil
		}
		return u.SignatureDelta, nil
	case string(AnthropicCitationsDeltaBlockTypeCitationsDelta):
		if u.CitationsDelta == nil {
			return nil, nil
		}
		return u.CitationsDelta, nil
	default:
		return nil, nil
	}
}

func (u *RawContentBlockDeltaUnion) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type

	switch tmp.Type {
	case string(AnthropicTextDeltaBlockTypeTextDelta):
		var b AnthropicTextDeltaBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.TextDelta, _ = json.Marshal(b)
		}
	case string(AnthropicThinkingDeltaBlockTypeThinkingDelta):
		var b AnthropicThinkingDeltaBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.ThinkingDelta, _ = json.Marshal(b)
		}
	case string(AnthropicInputJSONDeltaBlockTypeInputJSONDelta):
		var b AnthropicInputJSONDeltaBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.InputJSONDelta, _ = json.Marshal(b)
		}
	case string(AnthropicSignatureDeltaBlockTypeSignatureDelta):
		var b AnthropicSignatureDeltaBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.SignatureDelta, _ = json.Marshal(b)
		}
	case string(AnthropicCitationsDeltaBlockTypeCitationsDelta):
		var b AnthropicCitationsDeltaBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.CitationsDelta, _ = json.Marshal(b)
		}
	}
	return nil
}

// ============== MessageDeltaEventDelta ==============

// MessageDeltaEventDelta message_delta.delta 字段。
//
// 携带本轮 delta 后的停止原因 + 容器 + 停止序列等"终止信息"。
// 这些信息通常在最后一个 message_delta 里出现。
type MessageDeltaEventDelta struct {
	StopReason   *AnthropicStopReason  `json:"stop_reason,omitempty"`  // 仅最后一个 delta 出现
	StopDetails  *AnthropicStopDetails `json:"stop_details,omitempty"` // 仅 stop_reason=refusal 时
	StopSequence *string               `json:"stop_sequence,omitempty"`
	Container    *AnthropicContainer   `json:"container,omitempty"` // code execution 容器续接
}

// ============== MessageDeltaUsage ==============

// MessageDeltaUsage 流式 delta 期间附带的 usage 子集。
//
// 比顶层 AnthropicUsage 少了 CacheCreation 嵌套明细（5m/1h）+ InferenceGeo + ServiceTier，
// 因为这些在 message_start 已经发过。
type MessageDeltaUsage struct {
	InputTokens              int                           `json:"input_tokens"`
	OutputTokens             int                           `json:"output_tokens"`
	CacheCreationInputTokens int                           `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int                           `json:"cache_read_input_tokens,omitempty"`
	OutputTokensDetails      *AnthropicOutputTokensDetails `json:"output_tokens_details,omitempty"`
	ServerToolUse            *AnthropicServerToolUsage     `json:"server_tool_use,omitempty"`
}
