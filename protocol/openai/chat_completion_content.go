package openai

import "encoding/json"

// OpenAIChatCompletionContentPartUnion 内容块 union 接口（Request 侧）。
//
// 4 种实现（对齐官方 SDK ChatCompletionContentPartUnionParam）：Text / Image / InputAudio / File。
//
// 注：Assistant 响应输出的 Refusal / Audio 不在 Request 侧 ContentPartUnion 中；
// 它们作为 AssistantMessage 顶层字段（Refusal *string / Audio *MessageAudioParam）
// 在 chat_completion_message.go 中表达。
type OpenAIChatCompletionContentPartUnion interface {
	isChatCompletionContentPart()
}

// OpenAIChatCompletionContentPartUnionFlat 内容块 union 扁平化实现。
type OpenAIChatCompletionContentPartUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"text"/"image_url"/"input_audio"/"file"。
	Type string `json:"type,omitempty"`
	// Text 文本内容块。
	Text *OpenAIChatCompletionContentPartTextParam `json:"text,omitempty"`
	// Image 图片内容块。
	Image *OpenAIChatCompletionContentPartImageParam `json:"image_url,omitempty"`
	// InputAudio 输入音频内容块。
	InputAudio *OpenAIChatCompletionContentPartInputAudioParam `json:"input_audio,omitempty"`
	// File 文件内容块。
	File *OpenAIChatCompletionContentPartFileParam `json:"file,omitempty"`
}

func (OpenAIChatCompletionContentPartUnionFlat) isChatCompletionContentPart() {}

func (u OpenAIChatCompletionContentPartUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "text":
		if u.Text == nil {
			return nil, nil
		}
		return json.Marshal(u.Text)
	case "image_url":
		if u.Image == nil {
			return nil, nil
		}
		return json.Marshal(u.Image)
	case "input_audio":
		if u.InputAudio == nil {
			return nil, nil
		}
		return json.Marshal(u.InputAudio)
	case "file":
		if u.File == nil {
			return nil, nil
		}
		return json.Marshal(u.File)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionContentPartUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "text":
		var v OpenAIChatCompletionContentPartTextParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Text = &v
	case "image_url":
		var v OpenAIChatCompletionContentPartImageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Image = &v
	case "input_audio":
		var v OpenAIChatCompletionContentPartInputAudioParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.InputAudio = &v
	case "file":
		var v OpenAIChatCompletionContentPartFileParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.File = &v
	}
	return nil
}

// OpenAIChatCompletionContentPartType 内容块 type 字段类型。
type OpenAIChatCompletionContentPartType string

const (
	OpenAIChatCompletionContentPartTypeText       OpenAIChatCompletionContentPartType = "text"
	OpenAIChatCompletionContentPartTypeImageURL   OpenAIChatCompletionContentPartType = "image_url"
	OpenAIChatCompletionContentPartTypeInputAudio OpenAIChatCompletionContentPartType = "input_audio"
	OpenAIChatCompletionContentPartTypeFile       OpenAIChatCompletionContentPartType = "file"
)

// OpenAIChatCompletionContentPartTextParam 文本内容块。
type OpenAIChatCompletionContentPartTextParam struct {
	// Type 固定 "text"。
	Type OpenAIChatCompletionContentPartType `json:"type"`
	// Text 文本内容。
	Text string `json:"text"`
}

func (OpenAIChatCompletionContentPartTextParam) isChatCompletionContentPart() {}

// OpenAIChatCompletionContentPartImageParam 图片内容块。
type OpenAIChatCompletionContentPartImageParam struct {
	// Type 固定 "image_url"。
	Type OpenAIChatCompletionContentPartType `json:"type"`
	// ImageURL 图片 URL 或 base64。
	ImageURL OpenAIChatCompletionContentPartImageURL `json:"image_url"`
}

func (OpenAIChatCompletionContentPartImageParam) isChatCompletionContentPart() {}

// OpenAIChatCompletionContentPartImageURL 图片 URL 结构。
type OpenAIChatCompletionContentPartImageURL struct {
	// URL 图片 URL（HTTP URL 或 data:base64 格式）。
	URL string `json:"url"`
	// Detail 图片细节级别。
	Detail OpenAIImageDetail `json:"detail,omitempty"`
}

// OpenAIChatCompletionContentPartInputAudioParam 输入音频内容块。
type OpenAIChatCompletionContentPartInputAudioParam struct {
	// Type 固定 "input_audio"。
	Type OpenAIChatCompletionContentPartType `json:"type"`
	// InputAudio 音频数据。
	InputAudio OpenAIChatCompletionContentPartInputAudio `json:"input_audio"`
}

func (OpenAIChatCompletionContentPartInputAudioParam) isChatCompletionContentPart() {}

// OpenAIChatCompletionContentPartInputAudio 输入音频结构。
type OpenAIChatCompletionContentPartInputAudio struct {
	// Data base64 编码的音频数据。
	Data string `json:"data"`
	// Format 音频格式（wav / mp3）。
	Format OpenAIAudioFormat `json:"format"`
}

// OpenAIChatCompletionContentPartFileParam 文件内容块。
type OpenAIChatCompletionContentPartFileParam struct {
	// Type 固定 "file"。
	Type OpenAIChatCompletionContentPartType `json:"type"`
	// File 文件信息。
	File OpenAIChatCompletionContentPartFile `json:"file"`
}

func (OpenAIChatCompletionContentPartFileParam) isChatCompletionContentPart() {}

// OpenAIChatCompletionContentPartFile 文件结构。
type OpenAIChatCompletionContentPartFile struct {
	// FileData base64 编码的文件数据。
	FileData string `json:"file_data"`
	// Filename 文件名。
	Filename string `json:"filename,omitempty"`
}
