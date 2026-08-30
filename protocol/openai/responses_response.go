package openai

// OpenAIResponse Responses 响应顶层。
type OpenAIResponse struct {
	// ID 响应唯一 ID（如 "resp_xxx"）。
	ID string `json:"id"`
	// Object 固定 "response"。
	Object string `json:"object"`
	// CreatedAt 创建时间戳。
	CreatedAt int64 `json:"created_at"`
	// Model 实际使用的模型。
	Model string `json:"model"`
	// Status 状态（completed / failed / in_progress / incomplete / queued）。
	Status string `json:"status"`
	// Error 错误（可选）。
	Error *OpenAIResponseError `json:"error,omitempty"`
	// Instructions 系统指令。
	Instructions *string `json:"instructions,omitempty"`
	// Output 输出项列表。
	Output []OpenAIResponseOutputItemUnion `json:"output"`
	// IncompleteDetails 未完成详情（可选）。
	IncompleteDetails *OpenAIResponseIncompleteDetails `json:"incomplete_details,omitempty"`
	// Usage 用量。
	Usage *OpenAIResponseUsage `json:"usage,omitempty"`
	// ParallelToolCalls 是否允许并行工具调用。
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`
	// Temperature 温度。
	Temperature *float64 `json:"temperature,omitempty"`
	// ToolChoice 工具选择策略。
	ToolChoice OpenAIResponseToolChoiceUnion `json:"tool_choice,omitempty"`
	// Tools 工具列表。
	Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
	// TopP 核采样。
	TopP *float64 `json:"top_p,omitempty"`
	// MaxOutputTokens 最大输出 token 数。
	MaxOutputTokens *int64 `json:"max_output_tokens,omitempty"`
	// Truncation 截断策略。
	Truncation OpenAIResponseTruncation `json:"truncation,omitempty"`
	// Reasoning 推理配置。
	Reasoning *OpenAIResponseReasoningParam `json:"reasoning,omitempty"`
	// Metadata 元数据。
	Metadata map[string]string `json:"metadata,omitempty"`
	// Store 是否存储。
	Store *bool `json:"store,omitempty"`
	// Background 是否后台。
	Background *bool `json:"background,omitempty"`
	// ServiceTier 服务等级。
	ServiceTier OpenAIServiceTier `json:"service_tier,omitempty"`
	// SafetyIdentifier 安全标识。
	SafetyIdentifier string `json:"safety_identifier,omitempty"`
	// PromptCacheKey 提示缓存键。
	PromptCacheKey string `json:"prompt_cache_key,omitempty"`
	// Text 文本配置。
	Text *OpenAIResponseTextConfig `json:"text,omitempty"`
	// Conversation 会话引用。
	Conversation *OpenAIResponseConversationRef `json:"conversation,omitempty"`
}

// OpenAIResponseConversationRef 会话引用（响应侧，区别于 params 侧）。
type OpenAIResponseConversationRef struct {
	// ID 会话 ID。
	ID string `json:"id"`
}

// OpenAIResponseError 响应错误。
type OpenAIResponseError struct {
	// Code 错误代码。
	Code string `json:"code"`
	// Message 错误消息。
	Message string `json:"message"`
}

// OpenAIResponseIncompleteDetails 未完成详情。
type OpenAIResponseIncompleteDetails struct {
	// Reason 未完成原因（max_output_tokens / content_filter）。
	Reason string `json:"reason"`
}

// OpenAIResponseUsage 用量。
type OpenAIResponseUsage struct {
	// InputTokens 输入 token 数。
	InputTokens int64 `json:"input_tokens"`
	// OutputTokens 输出 token 数。
	OutputTokens int64 `json:"output_tokens"`
	// TotalTokens 总 token 数。
	TotalTokens int64 `json:"total_tokens"`
	// InputTokensDetails 输入 token 详情。
	InputTokensDetails OpenAIResponseUsageInputTokensDetails `json:"input_tokens_details"`
	// OutputTokensDetails 输出 token 详情。
	OutputTokensDetails OpenAIResponseUsageOutputTokensDetails `json:"output_tokens_details"`
}

// OpenAIResponseUsageInputTokensDetails 输入 token 详情。
type OpenAIResponseUsageInputTokensDetails struct {
	// CachedTokens 缓存命中 token 数。
	CachedTokens int64 `json:"cached_tokens"`
}

// OpenAIResponseUsageOutputTokensDetails 输出 token 详情。
type OpenAIResponseUsageOutputTokensDetails struct {
	// ReasoningTokens 推理 token 数。
	ReasoningTokens int64 `json:"reasoning_tokens"`
}
