package openai

import (
	"encoding/json"
)

// === Audio Translation 请求参数 ===

// OpenAITranslationParams 音频翻译请求参数。
//
// 字段集与 OpenAI 官方 SDK AudioTranslationNewParams 对齐。
// 详见 https://platform.openai.com/docs/api-reference/audio/create-translation。
//
// 复用 OpenAITranscriptionSegment / OpenAITranscriptionWord 等子结构（末尾 alias 已声明）。
type OpenAITranslationParams struct {
	// File 音频文件（base64 编码或 URL 占位）。
	File string `json:"file"`
	// Model 模型 ID（仅支持 whisper-1）。
	Model string `json:"model"`
	// Prompt 可选提示词（应使用英文）。
	Prompt string `json:"prompt,omitempty"`
	// Temperature 温度（0-1）。
	Temperature *float64 `json:"temperature,omitempty"`
	// ResponseFormat 响应格式（json / text / srt / verbose_json / vtt）。
	ResponseFormat OpenAITranslationResponseFormat `json:"response_format,omitempty"`
}

// OpenAITranslationResponseFormat 翻译响应格式枚举。
type OpenAITranslationResponseFormat string

const (
	OpenAITranslationResponseFormatJSON        OpenAITranslationResponseFormat = "json"
	OpenAITranslationResponseFormatText        OpenAITranslationResponseFormat = "text"
	OpenAITranslationResponseFormatSRT         OpenAITranslationResponseFormat = "srt"
	OpenAITranslationResponseFormatVerboseJSON OpenAITranslationResponseFormat = "verbose_json"
	OpenAITranslationResponseFormatVTT         OpenAITranslationResponseFormat = "vtt"
)

// === 响应 union：基础 Translation | Verbose Translation ===

// OpenAITranslationResponseUnion 音频翻译响应 union。
//
// 实现策略：OpenAITranslationResponseUnionImpl 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式保留用于 type switch 分发；新增 Impl 作为 concrete 实现。
// verbose_json 响应通过 language 字段的存在来区分（区别于基础响应只有 text 字段）。
type OpenAITranslationResponseUnion interface {
	isTranslationResponseUnion()
}

// OpenAITranslationResponseUnionImpl 音频翻译响应 union 的扁平实现。
type OpenAITranslationResponseUnionImpl struct {
	Type     string          `json:"type,omitempty"`
	Text     json.RawMessage `json:"text,omitempty"`
	Language json.RawMessage `json:"language,omitempty"`
	Duration json.RawMessage `json:"duration,omitempty"`
	Segments json.RawMessage `json:"segments,omitempty"`
}

func (OpenAITranslationResponseUnionImpl) isTranslationResponseUnion() {}

// MarshalJSON 实现 OpenAITranslationResponseUnionImpl 的自定义序列化。
func (u OpenAITranslationResponseUnionImpl) MarshalJSON() ([]byte, error) {
	if u.Type == "" && u.Text == nil && u.Language == nil {
		return nil, nil
	}
	// 优先检查 Language 字段（区分 verbose vs basic）
	if u.Language != nil {
		return u.Language, nil
	}
	if u.Text != nil {
		return u.Text, nil
	}
	return nil, nil
}

// UnmarshalJSON 实现 OpenAITranslationResponseUnionImpl 的自定义反序列化。
func (u *OpenAITranslationResponseUnionImpl) UnmarshalJSON(data []byte) error {
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
		var b OpenAITranslationVerbose
		if err := json.Unmarshal(data, &b); err == nil {
			u.Language, _ = json.Marshal(b)
		}
	} else {
		u.Type = "basic"
		var b OpenAITranslation
		if err := json.Unmarshal(data, &b); err == nil {
			u.Text, _ = json.Marshal(b)
		}
	}
	return nil
}

// OpenAITranslation 基础响应（json / text / srt / vtt）。
type OpenAITranslation struct {
	// Text 翻译文本（目标语言为英文）。
	Text string `json:"text"`
}

func (OpenAITranslation) isTranslationResponseUnion() {}

// OpenAITranslationVerbose Verbose JSON 响应（verbose_json 时返回）。
type OpenAITranslationVerbose struct {
	// Text 翻译文本。
	Text string `json:"text"`
	// Language 源语言代码（ISO-639-1）。
	Language string `json:"language,omitempty"`
	// Duration 音频时长（秒）。
	Duration *float64 `json:"duration,omitempty"`
	// Segments 段级时间戳。
	Segments []OpenAITranslationSegment `json:"segments,omitempty"`
}

func (OpenAITranslationVerbose) isTranslationResponseUnion() {}

// 复用 Transcription 子结构（段级时间戳 / 词级时间戳）。
type (
	// OpenAITranslationSegment 翻译段级时间戳（别名复用 OpenAITranscriptionSegment）。
	OpenAITranslationSegment = OpenAITranscriptionSegment
	// OpenAITranslationWord 翻译词级时间戳（别名复用 OpenAITranscriptionWord）。
	OpenAITranslationWord = OpenAITranscriptionWord
)
