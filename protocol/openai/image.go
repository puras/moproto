package openai

import "encoding/json"

// OpenAIImageModel 图像模型 ID。
type OpenAIImageModel string

const (
	OpenAIImageModelGPTImage1            OpenAIImageModel = "gpt-image-1"
	OpenAIImageModelGPTImage1Mini        OpenAIImageModel = "gpt-image-1-mini"
	OpenAIImageModelGPTImage1_5          OpenAIImageModel = "gpt-image-1.5"
	OpenAIImageModelGPTImage2            OpenAIImageModel = "gpt-image-2"
	OpenAIImageModelGPTImage2_2026_04_21 OpenAIImageModel = "gpt-image-2-2026-04-21"
	OpenAIImageModelChatgptImageLatest   OpenAIImageModel = "chatgpt-image-latest"
	OpenAIImageModelDallE2               OpenAIImageModel = "dall-e-2"
	OpenAIImageModelDallE3               OpenAIImageModel = "dall-e-3"
)

// OpenAIImageGenerateParams 图像生成请求参数。
type OpenAIImageGenerateParams struct {
	// Prompt 提示词。
	Prompt string `json:"prompt"`
	// Model 模型 ID（如 "dall-e-3" / "gpt-image-1"）。
	Model OpenAIImageModel `json:"model,omitempty"`
	// N 生成数量（dall-e-2 支持多张）。
	N *int64 `json:"n,omitempty"`
	// Quality 质量（standard / hd / low / medium / high / auto）。
	Quality string `json:"quality,omitempty"`
	// ResponseFormat 响应格式（url / b64_json）。
	ResponseFormat string `json:"response_format,omitempty"`
	// Size 尺寸（1024x1024 / 1024x1792 / 1792x1024 / auto）。
	Size string `json:"size,omitempty"`
	// Style 风格（vivid / natural，仅 dall-e-3）。
	Style string `json:"style,omitempty"`
	// User 用户标识。
	User string `json:"user,omitempty"`
	// Background 后台执行（gpt-image-1+ 支持）。
	Background *bool `json:"background,omitempty"`
	// OutputCompressionQuality 压缩质量（gpt-image-1+ 支持，仅 1-100 整数）。
	OutputCompressionQuality *int64 `json:"output_compression_quality,omitempty"`
	// PromptRefinement 提示词优化（gpt-image-1+ 支持，low / high）。
	PromptRefinement string `json:"prompt_refinement,omitempty"`
}

// OpenAIImageEditParams 图像编辑请求参数。
type OpenAIImageEditParams struct {
	// Image 图像（base64 或 URL）。
	Image string `json:"image"`
	// Prompt 编辑提示词。
	Prompt string `json:"prompt"`
	// Mask 蒙版（可选）。
	Mask string `json:"mask,omitempty"`
	// Model 模型 ID。
	Model OpenAIImageModel `json:"model,omitempty"`
	// N 生成数量。
	N *int64 `json:"n,omitempty"`
	// Quality 质量。
	Quality string `json:"quality,omitempty"`
	// ResponseFormat 响应格式。
	ResponseFormat string `json:"response_format,omitempty"`
	// Size 尺寸。
	Size string `json:"size,omitempty"`
	// User 用户标识。
	User string `json:"user,omitempty"`
}

// OpenAIImageVariationParams 图像变体请求参数。
//
// 官方 SDK：ImageNewVariationParams（对应 /v1/images/generations 的 variation 端点）。
// 详见 https://platform.openai.com/docs/api-reference/images/create-variation。
type OpenAIImageVariationParams struct {
	// Image 基础图像（PNG，<4MB，必须为正方形）。
	Image string `json:"image"`
	// Model 模型（仅支持 dall-e-2）。
	Model OpenAIImageModel `json:"model,omitempty"`
	// N 生成数量（1-10）。
	N *int64 `json:"n,omitempty"`
	// ResponseFormat 响应格式（url / b64_json）。
	ResponseFormat string `json:"response_format,omitempty"`
	// Size 尺寸（256x256 / 512x512 / 1024x1024）。
	Size string `json:"size,omitempty"`
	// User 用户标识。
	User string `json:"user,omitempty"`
}

// OpenAIImagesResponse 图像响应。
type OpenAIImagesResponse struct {
	// Created 创建时间戳。
	Created int64 `json:"created"`
	// Data 图像数据列表。
	Data []OpenAIImageData `json:"data"`
	// Usage 用量（gpt-image-1 支持）。
	Usage *OpenAIImageUsage `json:"usage,omitempty"`
}

// OpenAIImageData 单张图像数据。
type OpenAIImageData struct {
	// URL 图像 URL（response_format=url 时）。
	URL string `json:"url,omitempty"`
	// B64JSON base64 编码（response_format=b64_json 时）。
	B64JSON string `json:"b64_json,omitempty"`
	// RevisedPrompt 修订后的提示词（dall-e-3）。
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// OpenAIImageUsage 用量。
type OpenAIImageUsage struct {
	// InputTokens 输入 token 数。
	InputTokens int64 `json:"input_tokens"`
	// InputTokensDetails 输入 token 详情。
	InputTokensDetails OpenAIImageUsageInputTokensDetails `json:"input_tokens_details"`
	// OutputTokens 输出 token 数。
	OutputTokens int64 `json:"output_tokens"`
	// TotalTokens 总 token 数。
	TotalTokens int64 `json:"total_tokens"`
}

// OpenAIImageUsageInputTokensDetails 输入 token 详情。
type OpenAIImageUsageInputTokensDetails struct {
	// ImageTokens 图像 token 数。
	ImageTokens int64 `json:"image_tokens"`
	// TextTokens 文本 token 数。
	TextTokens int64 `json:"text_tokens"`
}

// OpenAIImageEditStreamEventUnion 图像编辑流事件 union。
type OpenAIImageEditStreamEventUnion interface {
	isImageEditStreamEvent()
}

// OpenAIImageEditStreamEventUnionFlat 图像编辑流事件 union 扁平化实现。
type OpenAIImageEditStreamEventUnionFlat struct {
	// Type 标识当前激活的变体：image_edit.partial_image/image_edit.completed。
	Type string `json:"type,omitempty"`
	// PartialImage 部分图像事件变体。
	PartialImage *OpenAIImageEditPartialImageEvent `json:"image_edit.partial_image,omitempty"`
	// Completed 图像编辑完成事件变体。
	Completed *OpenAIImageEditCompletedEvent `json:"image_edit.completed,omitempty"`
}

func (OpenAIImageEditStreamEventUnionFlat) isImageEditStreamEvent() {}

func (u OpenAIImageEditStreamEventUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "image_edit.partial_image":
		if u.PartialImage == nil {
			return nil, nil
		}
		return json.Marshal(u.PartialImage)
	case "image_edit.completed":
		if u.Completed == nil {
			return nil, nil
		}
		return json.Marshal(u.Completed)
	}
	return nil, nil
}

func (u *OpenAIImageEditStreamEventUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "image_edit.partial_image":
		var v OpenAIImageEditPartialImageEvent
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.PartialImage = &v
	case "image_edit.completed":
		var v OpenAIImageEditCompletedEvent
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Completed = &v
	}
	return nil
}

// OpenAIImageEditPartialImageEvent 部分图像事件。
type OpenAIImageEditPartialImageEvent struct {
	Type              string `json:"type"` // "image_edit.partial_image"
	PartialImageIndex int64  `json:"partial_image_index"`
	B64JSON           string `json:"b64_json"`
}

func (OpenAIImageEditPartialImageEvent) isImageEditStreamEvent() {}

// OpenAIImageEditCompletedEvent 图像编辑完成事件。
type OpenAIImageEditCompletedEvent struct {
	Type    string            `json:"type"` // "image_edit.completed"
	Usage   *OpenAIImageUsage `json:"usage,omitempty"`
	B64JSON string            `json:"b64_json,omitempty"`
}

func (OpenAIImageEditCompletedEvent) isImageEditStreamEvent() {}
