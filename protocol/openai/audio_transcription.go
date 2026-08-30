package openai

import (
	"encoding/json"
)

// OpenAIAudioTranscriptionParams 音频转录请求参数。
//
// 字段集与 OpenAI 官方 SDK AudioTranscriptionNewParams 严格对齐。
// 详见 https://platform.openai.com/docs/api-reference/audio/create-transcription。
type OpenAIAudioTranscriptionParams struct {
	// File 音频文件（base64 编码或 URL 占位）。
	File string `json:"file"`
	// Model 模型 ID（如 "whisper-1" / "gpt-4o-transcribe"）。
	Model string `json:"model"`
	// Language 音频语言代码（如 "en"）。
	Language string `json:"language,omitempty"`
	// Prompt 提示词（可选）。
	Prompt string `json:"prompt,omitempty"`
	// ResponseFormat 响应格式（json / text / srt / verbose_json / vtt / diarized_json）。
	ResponseFormat OpenAIAudioResponseFormat `json:"response_format,omitempty"`
	// Temperature 温度。
	Temperature *float64 `json:"temperature,omitempty"`
	// TimestampGranularities 时间戳粒度（word / segment，verbose_json 时可用）。
	TimestampGranularities []string `json:"timestamp_granularities,omitempty"`
	// ChunkingStrategy 分块策略（union: "auto" 常量 | VadConfig 对象）。
	ChunkingStrategy *OpenAIAudioTranscriptionChunkingStrategyUnionImpl `json:"chunking_strategy,omitempty"`
	// Include 包含的字段列表（logprobs 等）。
	Include []OpenAITranscriptionInclude `json:"include,omitempty"`
	// Stream 是否流式。
	Stream *bool `json:"stream,omitempty"`

	// === 4.6+ diarize 字段 ===

	// KnownSpeakerNames 已知说话人姓名列表（diarize 用）。
	KnownSpeakerNames []string `json:"known_speaker_names,omitempty"`
	// KnownSpeakerReferences 已知说话人参考音频（data URLs，diarize 用）。
	KnownSpeakerReferences []string `json:"known_speaker_references,omitempty"`
}

// OpenAITranscriptionInclude 流式输出包含字段。
type OpenAITranscriptionInclude string

const (
	OpenAITranscriptionIncludeLogprobs OpenAITranscriptionInclude = "logprobs"
)

// OpenAIAudioResponseFormat 响应格式枚举。
type OpenAIAudioResponseFormat string

const (
	OpenAIAudioResponseFormatJSON         OpenAIAudioResponseFormat = "json"
	OpenAIAudioResponseFormatText         OpenAIAudioResponseFormat = "text"
	OpenAIAudioResponseFormatSRT          OpenAIAudioResponseFormat = "srt"
	OpenAIAudioResponseFormatVerboseJSON  OpenAIAudioResponseFormat = "verbose_json"
	OpenAIAudioResponseFormatVTT          OpenAIAudioResponseFormat = "vtt"
	OpenAIAudioResponseFormatDiarizedJSON OpenAIAudioResponseFormat = "diarized_json"
)

// OpenAIAudioTranscriptionChunkingStrategyUnion 分块策略 union（"auto" 常量 | VadConfig 对象）。
//
// 实现策略：OpenAIAudioTranscriptionChunkingStrategyUnionImpl 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式保留用于 type switch 分发；新增 Impl 作为 concrete 实现。
type OpenAIAudioTranscriptionChunkingStrategyUnion interface {
	isChunkingStrategyUnion()
}

// OpenAIAudioTranscriptionChunkingStrategyUnionImpl 分块策略 union 的扁平实现。
type OpenAIAudioTranscriptionChunkingStrategyUnionImpl struct {
	Type      string          `json:"type,omitempty"`
	Auto      string          `json:"-"` // "auto" 字符串，MarshalJSON 单独处理
	VadConfig json.RawMessage `json:"vad_config,omitempty"`
}

func (OpenAIAudioTranscriptionChunkingStrategyUnionImpl) isChunkingStrategyUnion() {}

// MarshalJSON 实现 OpenAIAudioTranscriptionChunkingStrategyUnionImpl 的自定义序列化。
//
// "auto" 场景输出为 JSON string "auto"，匹配 wire-format。
// VadConfig 对象输出为普通 JSON 对象。
func (u OpenAIAudioTranscriptionChunkingStrategyUnionImpl) MarshalJSON() ([]byte, error) {
	if u.Type == "" && u.Auto == "" && u.VadConfig == nil {
		return nil, nil
	}
	if u.Auto == "auto" {
		// 输出 JSON string "auto"
		return json.Marshal("auto")
	}
	if u.VadConfig != nil {
		return u.VadConfig, nil
	}
	return nil, nil
}

// UnmarshalJSON 实现 OpenAIAudioTranscriptionChunkingStrategyUnionImpl 的自定义反序列化。
func (u *OpenAIAudioTranscriptionChunkingStrategyUnionImpl) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	// 检查是否为字符串 "auto"
	var autoStr string
	if err := json.Unmarshal(data, &autoStr); err == nil && autoStr == "auto" {
		u.Type = "auto"
		u.Auto = "auto"
		return nil
	}
	// 否则为对象 VadConfig
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	if tmp.Type == "server_vad" {
		var b OpenAIAudioTranscriptionChunkingStrategyVadConfig
		if err := json.Unmarshal(data, &b); err == nil {
			u.VadConfig, _ = json.Marshal(b)
		}
	}
	return nil
}

// OpenAIAudioTranscriptionChunkingStrategyAuto "auto" 常量。
type OpenAIAudioTranscriptionChunkingStrategyAuto string

const OpenAIAudioTranscriptionChunkingStrategyAutoValue OpenAIAudioTranscriptionChunkingStrategyAuto = "auto"

func (OpenAIAudioTranscriptionChunkingStrategyAuto) isChunkingStrategyUnion() {}

// OpenAIAudioTranscriptionChunkingStrategyVadConfig 服务端 VAD 配置（"server_vad"）。
type OpenAIAudioTranscriptionChunkingStrategyVadConfig struct {
	// Type 固定为 "server_vad"。
	Type string `json:"type"`
	// PrefixPadMs 前缀填充毫秒数（官方命名 prefix_padding_ms）。
	PrefixPadMs *int64 `json:"prefix_padding_ms,omitempty"`
	// SilenceDurationMs 静音持续毫秒数。
	SilenceDurationMs *int64 `json:"silence_duration_ms,omitempty"`
	// Threshold VAD 阈值。
	Threshold *float64 `json:"threshold,omitempty"`
}

func (OpenAIAudioTranscriptionChunkingStrategyVadConfig) isChunkingStrategyUnion() {}

// === 响应 union：基础 Transcription | Verbose Transcription（带 Duration/Language/Segments/Words） ===

// OpenAIAudioTranscriptionResponseUnion 音频转录响应 union。
//
// 实现策略：OpenAIAudioTranscriptionResponseUnionImpl 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式保留用于 type switch 分发；新增 Impl 作为 concrete 实现。
// verbose_json 响应通过 language 字段区分。
type OpenAIAudioTranscriptionResponseUnion interface {
	isAudioTranscriptionResponseUnion()
}

// OpenAIAudioTranscriptionResponseUnionImpl 音频转录响应 union 的扁平实现。
type OpenAIAudioTranscriptionResponseUnionImpl struct {
	Type    string          `json:"type,omitempty"`
	Basic   json.RawMessage `json:"basic,omitempty"`
	Verbose json.RawMessage `json:"verbose,omitempty"`
}

func (OpenAIAudioTranscriptionResponseUnionImpl) isAudioTranscriptionResponseUnion() {}

// MarshalJSON 实现 OpenAIAudioTranscriptionResponseUnionImpl 的自定义序列化。
func (u OpenAIAudioTranscriptionResponseUnionImpl) MarshalJSON() ([]byte, error) {
	if u.Type == "" && u.Basic == nil && u.Verbose == nil {
		return nil, nil
	}
	if u.Verbose != nil {
		return u.Verbose, nil
	}
	if u.Basic != nil {
		return u.Basic, nil
	}
	return nil, nil
}

// UnmarshalJSON 实现 OpenAIAudioTranscriptionResponseUnionImpl 的自定义反序列化。
func (u *OpenAIAudioTranscriptionResponseUnionImpl) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var tmp struct {
		Language string `json:"language,omitempty"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if tmp.Language != "" {
		u.Type = "verbose"
		var b OpenAIAudioTranscriptionVerbose
		if err := json.Unmarshal(data, &b); err == nil {
			u.Verbose, _ = json.Marshal(b)
		}
	} else {
		u.Type = "basic"
		var b OpenAIAudioTranscription
		if err := json.Unmarshal(data, &b); err == nil {
			u.Basic, _ = json.Marshal(b)
		}
	}
	return nil
}

// OpenAIAudioTranscription 基础响应（json / text）。
type OpenAIAudioTranscription struct {
	// Text 转录文本。
	Text string `json:"text"`
	// Logprobs token log probabilities（仅当 include[]=logprobs 时出现）。
	Logprobs []OpenAITranscriptionLogprob `json:"logprobs,omitempty"`
	// Usage 用量（按 tokens 计费）。
	Usage *OpenAITranscriptionUsageUnionImpl `json:"usage,omitempty"`
}

func (OpenAIAudioTranscription) isAudioTranscriptionResponseUnion() {}

// OpenAIAudioTranscriptionVerbose Verbose JSON 响应（带 Duration/Language/Segments/Words）。
type OpenAIAudioTranscriptionVerbose struct {
	// Text 转录文本。
	Text string `json:"text"`
	// Language 语言（verbose_json 必填）。
	Language string `json:"language,omitempty"`
	// Duration 时长（verbose_json 必填）。
	Duration *float64 `json:"duration,omitempty"`
	// Words 词级时间戳。
	Words []OpenAIAudioTranscriptionWord `json:"words,omitempty"`
	// Segments 段级时间戳。
	Segments []OpenAIAudioTranscriptionSegment `json:"segments,omitempty"`
	// Usage 用量（按 duration 秒计费）。
	Usage *OpenAITranscriptionUsageUnionImpl `json:"usage,omitempty"`
}

func (OpenAIAudioTranscriptionVerbose) isAudioTranscriptionResponseUnion() {}

// === Usage Union（tokens | duration）===

// OpenAITranscriptionUsageUnion 用量 union（按 tokens | 按 duration 计费）。
//
// 实现策略：OpenAITranscriptionUsageUnionImpl 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式保留用于 type switch 分发；新增 Impl 作为 concrete 实现。
type OpenAITranscriptionUsageUnion interface {
	isTranscriptionUsageUnion()
}

// OpenAITranscriptionUsageUnionImpl 用量 union 的扁平实现。
type OpenAITranscriptionUsageUnionImpl struct {
	Type    string          `json:"type,omitempty"`
	Tokens  json.RawMessage `json:"tokens,omitempty"`
	Seconds json.RawMessage `json:"seconds,omitempty"`
}

func (OpenAITranscriptionUsageUnionImpl) isTranscriptionUsageUnion() {}

// MarshalJSON 实现 OpenAITranscriptionUsageUnionImpl 的自定义序列化。
func (u OpenAITranscriptionUsageUnionImpl) MarshalJSON() ([]byte, error) {
	if u.Type == "" && u.Tokens == nil && u.Seconds == nil {
		return nil, nil
	}
	if u.Tokens != nil {
		return u.Tokens, nil
	}
	if u.Seconds != nil {
		return u.Seconds, nil
	}
	return nil, nil
}

// UnmarshalJSON 实现 OpenAITranscriptionUsageUnionImpl 的自定义反序列化。
func (u *OpenAITranscriptionUsageUnionImpl) UnmarshalJSON(data []byte) error {
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
	case "tokens":
		var b OpenAITranscriptionUsageTokens
		if err := json.Unmarshal(data, &b); err == nil {
			u.Tokens, _ = json.Marshal(b)
		}
	case "duration":
		var b OpenAITranscriptionUsageDuration
		if err := json.Unmarshal(data, &b); err == nil {
			u.Seconds, _ = json.Marshal(b)
		}
	}
	return nil
}

// OpenAITranscriptionUsageTokens tokens 计费模式。
type OpenAITranscriptionUsageTokens struct {
	// Type 固定为 "tokens"。
	Type string `json:"type"`
	// InputTokens 输入 token 数。
	InputTokens int64 `json:"input_tokens"`
	// OutputTokens 输出 token 数。
	OutputTokens int64 `json:"output_tokens"`
	// TotalTokens 总 token 数。
	TotalTokens int64 `json:"total_tokens"`
	// InputTokenDetails 输入 token 明细（含 audio_tokens / text_tokens）。
	InputTokenDetails *OpenAITranscriptionUsageTokensInputTokenDetails `json:"input_token_details,omitempty"`
}

func (OpenAITranscriptionUsageTokens) isTranscriptionUsageUnion() {}

// OpenAITranscriptionUsageTokensInputTokenDetails tokens 输入明细。
type OpenAITranscriptionUsageTokensInputTokenDetails struct {
	// AudioTokens 音频 token 数。
	AudioTokens int64 `json:"audio_tokens,omitempty"`
	// TextTokens 文本 token 数。
	TextTokens int64 `json:"text_tokens,omitempty"`
}

// OpenAITranscriptionUsageDuration duration 计费模式。
type OpenAITranscriptionUsageDuration struct {
	// Type 固定为 "duration"。
	Type string `json:"type"`
	// Seconds 音频秒数。
	Seconds float64 `json:"seconds"`
}

func (OpenAITranscriptionUsageDuration) isTranscriptionUsageUnion() {}

// === 流事件 union（3 变体）===

// OpenAITranscriptionStreamEvent 流事件 union。
//
// 实现策略：OpenAITranscriptionStreamEventUnionImpl 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式保留用于 type switch 分发；新增 Impl 作为 concrete 实现。
type OpenAITranscriptionStreamEvent interface {
	isTranscriptionStreamEvent()
}

// OpenAITranscriptionStreamEventUnionImpl 流事件 union 的扁平实现。
type OpenAITranscriptionStreamEventUnionImpl struct {
	Type    string          `json:"type,omitempty"`
	Segment json.RawMessage `json:"segment,omitempty"`
	Delta   json.RawMessage `json:"delta,omitempty"`
	Done    json.RawMessage `json:"done,omitempty"`
}

func (OpenAITranscriptionStreamEventUnionImpl) isTranscriptionStreamEvent() {}

// MarshalJSON 实现 OpenAITranscriptionStreamEventUnionImpl 的自定义序列化。
func (u OpenAITranscriptionStreamEventUnionImpl) MarshalJSON() ([]byte, error) {
	if u.Type == "" && u.Segment == nil && u.Delta == nil && u.Done == nil {
		return nil, nil
	}
	if u.Segment != nil {
		return u.Segment, nil
	}
	if u.Delta != nil {
		return u.Delta, nil
	}
	if u.Done != nil {
		return u.Done, nil
	}
	return nil, nil
}

// UnmarshalJSON 实现 OpenAITranscriptionStreamEventUnionImpl 的自定义反序列化。
func (u *OpenAITranscriptionStreamEventUnionImpl) UnmarshalJSON(data []byte) error {
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
	case string(OpenAITranscriptionStreamEventTypeSegment):
		var b OpenAITranscriptionTextSegmentEvent
		if err := json.Unmarshal(data, &b); err == nil {
			u.Segment, _ = json.Marshal(b)
		}
	case string(OpenAITranscriptionStreamEventTypeDelta):
		var b OpenAITranscriptionTextDeltaEvent
		if err := json.Unmarshal(data, &b); err == nil {
			u.Delta, _ = json.Marshal(b)
		}
	case string(OpenAITranscriptionStreamEventTypeDone):
		var b OpenAITranscriptionTextDoneEvent
		if err := json.Unmarshal(data, &b); err == nil {
			u.Done, _ = json.Marshal(b)
		}
	}
	return nil
}

// OpenAITranscriptionTextSegmentEvent 说话人分段事件（diarize 流事件）。
type OpenAITranscriptionTextSegmentEvent struct {
	// Type 固定 "transcript.text.segment"。
	Type OpenAITranscriptionStreamEventType `json:"type"`
	// ID 段标识。
	ID string `json:"id"`
	// End 段结束时间（秒）。
	End float64 `json:"end"`
	// Speaker 说话人标签。
	Speaker string `json:"speaker"`
	// Start 段开始时间（秒）。
	Start float64 `json:"start"`
	// Text 段转录文本。
	Text string `json:"text"`
}

func (OpenAITranscriptionTextSegmentEvent) isTranscriptionStreamEvent() {}

// OpenAITranscriptionTextDeltaEvent 文本增量事件。
type OpenAITranscriptionTextDeltaEvent struct {
	// Type 固定 "transcript.text.delta"。
	Type OpenAITranscriptionStreamEventType `json:"type"`
	// Delta 增量文本。
	Delta string `json:"delta"`
	// Logprobs 增量 token logprobs（仅当 include[]=logprobs 时出现）。
	Logprobs []OpenAITranscriptionLogprob `json:"logprobs,omitempty"`
	// SegmentID 所属段标识（diarize 时）。
	SegmentID string `json:"segment_id,omitempty"`
}

func (OpenAITranscriptionTextDeltaEvent) isTranscriptionStreamEvent() {}

// OpenAITranscriptionTextDoneEvent 文本完成事件。
type OpenAITranscriptionTextDoneEvent struct {
	// Type 固定 "transcript.text.done"。
	Type OpenAITranscriptionStreamEventType `json:"type"`
	// Text 完整文本。
	Text string `json:"text"`
	// Logprobs 完整 token logprobs（仅当 include[]=logprobs 时出现）。
	Logprobs []OpenAITranscriptionLogprob `json:"logprobs,omitempty"`
	// Usage 用量（按 tokens）。
	Usage *OpenAITranscriptionUsageTokens `json:"usage,omitempty"`
}

func (OpenAITranscriptionTextDoneEvent) isTranscriptionStreamEvent() {}

// OpenAITranscriptionStreamEventType 流事件类型（type 字符串枚举）。
type OpenAITranscriptionStreamEventType string

const (
	OpenAITranscriptionStreamEventTypeSegment OpenAITranscriptionStreamEventType = "transcript.text.segment"
	OpenAITranscriptionStreamEventTypeDelta   OpenAITranscriptionStreamEventType = "transcript.text.delta"
	OpenAITranscriptionStreamEventTypeDone    OpenAITranscriptionStreamEventType = "transcript.text.done"
)

// === 流事件子结构 ===

// OpenAITranscriptionLogprob 单个 token log 概率。
type OpenAITranscriptionLogprob struct {
	// Token token 字符串。
	Token string `json:"token"`
	// Bytes token 字节。
	Bytes []float64 `json:"bytes"`
	// Logprob 对数概率。
	Logprob float64 `json:"logprob"`
}

// OpenAITranscriptionWord 词级时间戳。
type OpenAITranscriptionWord struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// OpenAITranscriptionSegment 段级时间戳。
type OpenAITranscriptionSegment struct {
	ID               int64   `json:"id"`
	Seek             float64 `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int64 `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

// 保留旧名 type alias 兼容（避免破坏外部引用）
type (
	// OpenAIAudioTranscriptionWord 旧名 → 新名 alias。
	OpenAIAudioTranscriptionWord = OpenAITranscriptionWord
	// OpenAIAudioTranscriptionSegment 旧名 → 新名 alias。
	OpenAIAudioTranscriptionSegment = OpenAITranscriptionSegment
)
