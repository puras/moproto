package openai

import (
	"encoding/json"
)

// OpenAIAudioSpeechParams 语音合成请求参数。
type OpenAIAudioSpeechParams struct {
	// Model 模型 ID（如 "tts-1" / "gpt-4o-mini-tts"）。
	Model string `json:"model"`
	// Input 要合成的文本。
	Input string `json:"input"`
	// Voice 语音类型。
	Voice OpenAIAudioVoice `json:"voice"`
	// ResponseFormat 输出格式（mp3 / opus / aac / flac / wav / pcm16）。
	ResponseFormat OpenAIAudioFormat `json:"response_format,omitempty"`
	// Speed 语速（0.25 - 4.0）。
	Speed *float64 `json:"speed,omitempty"`
	// Instructions 语音风格指令（gpt-4o-mini-tts 支持）。
	Instructions string `json:"instructions,omitempty"`
	// StreamFormat 流式输出格式（sse / audio）。
	StreamFormat string `json:"stream_format,omitempty"`
}

// OpenAIAudioSpeechStreamEventUnion 语音合成流事件 union。
//
// 实现策略：OpenAIAudioSpeechStreamEventUnion 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式（OpenAIAudioSpeechStreamEventUnion + 2 个 isAudioSpeechStreamEvent 实现）
// 保留用于 type switch 分发；新增 OpenAIAudioSpeechStreamEventUnion 作为唯一的 concrete 实现，
// 解决 json.Unmarshal 无法反序列化到 interface 字段的问题。
type OpenAIAudioSpeechStreamEventUnion interface {
	isAudioSpeechStreamEvent()
}

// OpenAIAudioSpeechStreamEventUnionImpl 语音合成流事件 union 的扁平实现。
type OpenAIAudioSpeechStreamEventUnionImpl struct {
	Type  string          `json:"type"`
	Delta json.RawMessage `json:"delta,omitempty"`
	Done  json.RawMessage `json:"done,omitempty"`
}

func (OpenAIAudioSpeechStreamEventUnionImpl) isAudioSpeechStreamEvent() {}

// MarshalJSON 实现 OpenAIAudioSpeechStreamEventUnionImpl 的自定义序列化。
func (u OpenAIAudioSpeechStreamEventUnionImpl) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "speech.audio.delta":
		if u.Delta == nil {
			return nil, nil
		}
		return u.Delta, nil
	case "speech.audio.done":
		if u.Done == nil {
			return nil, nil
		}
		return u.Done, nil
	default:
		return nil, nil
	}
}

// UnmarshalJSON 实现 OpenAIAudioSpeechStreamEventUnionImpl 的自定义反序列化。
func (u *OpenAIAudioSpeechStreamEventUnionImpl) UnmarshalJSON(data []byte) error {
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
	case "speech.audio.delta":
		var b OpenAIAudioSpeechAudioDeltaEvent
		if err := json.Unmarshal(data, &b); err == nil {
			u.Delta, _ = json.Marshal(b)
		}
	case "speech.audio.done":
		var b OpenAIAudioSpeechAudioDoneEvent
		if err := json.Unmarshal(data, &b); err == nil {
			u.Done, _ = json.Marshal(b)
		}
	}
	return nil
}

// OpenAIAudioSpeechAudioDeltaEvent 音频增量事件。
type OpenAIAudioSpeechAudioDeltaEvent struct {
	Type  string `json:"type"` // "speech.audio.delta"
	Audio string `json:"audio"`
}

func (OpenAIAudioSpeechAudioDeltaEvent) isAudioSpeechStreamEvent() {}

// OpenAIAudioSpeechAudioDoneEvent 音频完成事件。
type OpenAIAudioSpeechAudioDoneEvent struct {
	Type  string `json:"type"` // "speech.audio.done"
	Audio string `json:"audio"`
}

func (OpenAIAudioSpeechAudioDoneEvent) isAudioSpeechStreamEvent() {}
