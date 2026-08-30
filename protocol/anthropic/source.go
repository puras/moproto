package anthropic

import "encoding/json"

// ============== AnthropicImageSource union ==============
//
// 对齐官方 anthropic-sdk-go ImageBlockParamSourceUnion（message.go:3339-3343）。
// 2 种实现：Base64 / URL（与官方 SDK 一致；本包历史曾有 file 类型，已移除以对齐官方）。
//
// 通过 isImageSource() 私有方法封闭实现集；调用方通过 type switch 或各实现的
// Type 字段分发，详见 doc.go。
type AnthropicImageSource interface {
	isImageSource()
}

// ============== 类型化 Type 字段 ==============

// AnthropicImageSourceBase64Type base64 类型的 type 字段类型。
type AnthropicImageSourceBase64Type string

// AnthropicImageSourceBase64TypeBase64 base64 唯一合法取值。
const AnthropicImageSourceBase64TypeBase64 AnthropicImageSourceBase64Type = "base64"

// AnthropicImageSourceURLType url 类型的 type 字段类型。
type AnthropicImageSourceURLType string

// AnthropicImageSourceURLTypeURL url 唯一合法取值。
const AnthropicImageSourceURLTypeURL AnthropicImageSourceURLType = "url"

// ============== 2 个实现 ==============

// AnthropicImageSourceBase64 内联 base64 编码图片（对齐官方 Base64ImageSourceParam）。
type AnthropicImageSourceBase64 struct {
	Type      AnthropicImageSourceBase64Type `json:"type"`       // 固定 "base64"
	MediaType string                         `json:"media_type"` // 必填 MIME（如 image/jpeg）
	Data      string                         `json:"data"`       // 必填 base64 编码数据
}

func (AnthropicImageSourceBase64) isImageSource() {}

// AnthropicImageSourceURL 远程 HTTP(S) 链接图片（对齐官方 URLImageSourceParam）。
type AnthropicImageSourceURL struct {
	Type AnthropicImageSourceURLType `json:"type"` // 固定 "url"
	URL  string                      `json:"url"`  // 必填 HTTP(S) URL
}

func (AnthropicImageSourceURL) isImageSource() {}

// ============== 构造器 ==============

// NewAnthropicImageSourceBase64 构造 base64 图片源。
func NewAnthropicImageSourceBase64(mediaType, data string) AnthropicImageSourceBase64 {
	return AnthropicImageSourceBase64{
		Type:      AnthropicImageSourceBase64TypeBase64,
		MediaType: mediaType,
		Data:      data,
	}
}

// NewAnthropicImageSourceURL 构造 URL 图片源。
func NewAnthropicImageSourceURL(url string) AnthropicImageSourceURL {
	return AnthropicImageSourceURL{
		Type: AnthropicImageSourceURLTypeURL,
		URL:  url,
	}
}

// ============== AnthropicDocumentSource union ==============
//
// 对齐官方 anthropic-sdk-go DocumentBlockParamSourceUnion（message.go:3183-3189）。
// 5 种实现：Base64 / URL / File / Text / Content。
//
// 与 Image Source 的差异：Document 多 2 个 type（file + content），保留 file
// 向后兼容（Files API 历史上支持）；新增 content 支持嵌套内容块（官方对齐）。
type AnthropicDocumentSource interface {
	isDocumentSource()
}

// ============== 类型化 Type 字段 ==============

// AnthropicDocumentSourceBase64Type base64 类型的 type 字段类型。
type AnthropicDocumentSourceBase64Type string

// AnthropicDocumentSourceBase64TypeBase64 base64 唯一合法取值。
const AnthropicDocumentSourceBase64TypeBase64 AnthropicDocumentSourceBase64Type = "base64"

// AnthropicDocumentSourceURLType url 类型的 type 字段类型。
type AnthropicDocumentSourceURLType string

// AnthropicDocumentSourceURLTypeURL url 唯一合法取值。
const AnthropicDocumentSourceURLTypeURL AnthropicDocumentSourceURLType = "url"

// AnthropicDocumentSourceFileType file 类型的 type 字段类型（本包扩展）。
type AnthropicDocumentSourceFileType string

// AnthropicDocumentSourceFileTypeFile file 唯一合法取值。
const AnthropicDocumentSourceFileTypeFile AnthropicDocumentSourceFileType = "file"

// AnthropicDocumentSourceTextType text 类型的 type 字段类型。
type AnthropicDocumentSourceTextType string

// AnthropicDocumentSourceTextTypeText text 唯一合法取值。
const AnthropicDocumentSourceTextTypeText AnthropicDocumentSourceTextType = "text"

// AnthropicDocumentSourceContentType content 类型的 type 字段类型（官方新增）。
type AnthropicDocumentSourceContentType string

// AnthropicDocumentSourceContentTypeContent content 唯一合法取值。
const AnthropicDocumentSourceContentTypeContent AnthropicDocumentSourceContentType = "content"

// ============== 5 个实现 ==============

// AnthropicDocumentSourceBase64 内联 base64 文档（PDF / 文本等）。
type AnthropicDocumentSourceBase64 struct {
	Type      AnthropicDocumentSourceBase64Type `json:"type"`       // 固定 "base64"
	MediaType string                            `json:"media_type"` // 必填 MIME（如 application/pdf）
	Data      string                            `json:"data"`       // 必填 base64 编码数据
}

func (AnthropicDocumentSourceBase64) isDocumentSource() {}

// AnthropicDocumentSourceURL 远程 HTTP(S) 链接文档。
type AnthropicDocumentSourceURL struct {
	Type AnthropicDocumentSourceURLType `json:"type"` // 固定 "url"
	URL  string                         `json:"url"`  // 必填 HTTP(S) URL
}

func (AnthropicDocumentSourceURL) isDocumentSource() {}

// AnthropicDocumentSourceFile Files API 上传的文件引用（本包扩展保留）。
type AnthropicDocumentSourceFile struct {
	Type   AnthropicDocumentSourceFileType `json:"type"`    // 固定 "file"
	FileID string                          `json:"file_id"` // 必填 文件 ID
}

func (AnthropicDocumentSourceFile) isDocumentSource() {}

// AnthropicDocumentSourceText 内联纯文本文档。
type AnthropicDocumentSourceText struct {
	Type      AnthropicDocumentSourceTextType `json:"type"`       // 固定 "text"
	MediaType string                          `json:"media_type"` // 必填 MIME（text/plain 等）
	Data      string                          `json:"data"`       // 必填 原始文本
}

func (AnthropicDocumentSourceText) isDocumentSource() {}

// AnthropicDocumentSourceContent 嵌套内容块文档（官方 SDK ContentBlockSourceParam）。
//
// Content 字段承载递归的 AnthropicContentBlockParam union；调用方按 type switch 处理。
type AnthropicDocumentSourceContent struct {
	Type    AnthropicDocumentSourceContentType `json:"type"`    // 固定 "content"
	Content json.RawMessage                    `json:"content"` // 必填 嵌套 content blocks
}

func (AnthropicDocumentSourceContent) isDocumentSource() {}

// ============== 构造器 ==============

// NewAnthropicDocumentSourceBase64 构造 base64 文档源。
func NewAnthropicDocumentSourceBase64(mediaType, data string) AnthropicDocumentSourceBase64 {
	return AnthropicDocumentSourceBase64{
		Type:      AnthropicDocumentSourceBase64TypeBase64,
		MediaType: mediaType,
		Data:      data,
	}
}

// NewAnthropicDocumentSourceURL 构造 URL 文档源。
func NewAnthropicDocumentSourceURL(url string) AnthropicDocumentSourceURL {
	return AnthropicDocumentSourceURL{
		Type: AnthropicDocumentSourceURLTypeURL,
		URL:  url,
	}
}

// NewAnthropicDocumentSourceFile 构造 file 文档源。
func NewAnthropicDocumentSourceFile(fileID string) AnthropicDocumentSourceFile {
	return AnthropicDocumentSourceFile{
		Type:   AnthropicDocumentSourceFileTypeFile,
		FileID: fileID,
	}
}

// NewAnthropicDocumentSourceText 构造 text 文档源。
func NewAnthropicDocumentSourceText(mediaType, data string) AnthropicDocumentSourceText {
	return AnthropicDocumentSourceText{
		Type:      AnthropicDocumentSourceTextTypeText,
		MediaType: mediaType,
		Data:      data,
	}
}

// NewAnthropicDocumentSourceContent 构造嵌套 content 文档源。
//
// content 字段承载递归的 content blocks（json.RawMessage 透传）。
func NewAnthropicDocumentSourceContent(content json.RawMessage) AnthropicDocumentSourceContent {
	return AnthropicDocumentSourceContent{
		Type:    AnthropicDocumentSourceContentTypeContent,
		Content: content,
	}
}
