package gemini

// GeminiGenerateContentRequest Gemini 协议顶层请求 body（generateContent / streamGenerateContent）。
//
// 字段命名严格遵循 Google 官方 camelCase：contents / systemInstruction /
// tools / toolConfig / safetySettings / generationConfig / cachedContent。
type GeminiGenerateContentRequest struct {
	Contents          []*GeminiContent        `json:"contents"`
	SystemInstruction *GeminiContent          `json:"systemInstruction,omitempty"`
	Tools             []*GeminiTool           `json:"tools,omitempty"`
	ToolConfig        *GeminiToolConfig       `json:"toolConfig,omitempty"`
	SafetySettings    []*GeminiSafetySetting  `json:"safetySettings,omitempty"`
	GenerationConfig  *GeminiGenerationConfig `json:"generationConfig,omitempty"`
	CachedContent     string                  `json:"cachedContent,omitempty"`
}

// GeminiRequest 适配器类型 — 持有 GeminiGenerateContentRequest 并实现 common.Request。
//
// MoGo 协议层的请求统一通过 common.Request 接口传递；本类型做两件事：
//  1. 嵌入 *GeminiGenerateContentRequest（方便类型断言访问顶层字段）；
//  2. 实现 GetFormat / GetModel 方法。
//
// GetModel 返回空字符串：Gemini 协议的 model 名在 URL path 而非 body，
// 由 handler.ExtractModelMeta 从 c.Param("model") 提取并存入 extras.Model。
type GeminiRequest struct {
	*GeminiGenerateContentRequest
}

// GetFormat 返回协议标识 "gemini"。
func (r *GeminiRequest) GetFormat() string {
	return "gemini"
}

// GetModel 返回请求中的 model（始终空，由 handler 路径提取覆盖）。
func (r *GeminiRequest) GetModel() string {
	return ""
}
