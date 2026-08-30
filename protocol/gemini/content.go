package gemini

// GeminiContent 一条对话消息（Gemini 协议 basic unit）。
//
// 字段：
//   - Parts: 一个或多个 Part；Text / InlineData / FileData / FunctionCall / FunctionResponse
//     互斥单选；
//   - Role: 仅 "user" | "model"。省略时默认 "user"。
type GeminiContent struct {
	Parts []*GeminiPart `json:"parts"`
	Role  GeminiRole    `json:"role,omitempty"`
}

// GeminiPart 内容片段（多字段互斥 union）。
//
// 七个字段互斥：Google 官方协议规定每个 Part 仅有"一种角色"。
// 调用方按需填一个字段，其余保持零值；序列化时 JSON 自动省略零值字段。
//
// Thought / ThoughtSignature 用于标记/携带 extended thinking 状态（Gemini 2.5 系列
// 支持）；不在 P0 默认使用范围。
type GeminiPart struct {
	Text             string                  `json:"text,omitempty"`
	InlineData       *GeminiBlob             `json:"inlineData,omitempty"`
	FileData         *GeminiFileData         `json:"fileData,omitempty"`
	FunctionCall     *GeminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *GeminiFunctionResponse `json:"functionResponse,omitempty"`
	Thought          bool                    `json:"thought,omitempty"`
	ThoughtSignature []byte                  `json:"thoughtSignature,omitempty"`
}

// GeminiFileData 引用已上传文件的输入。
//
// 用于 Files API 上传后传入 URI；不在 P0 默认范围（需先上传文件）。
type GeminiFileData struct {
	MIMEType string `json:"mimeType"`
	FileURI  string `json:"fileUri"`
}
