package openai

// OpenAIChatCompletionChunk 流式 chunk。
type OpenAIChatCompletionChunk struct {
	// ID 响应唯一 ID。
	ID string `json:"id"`
	// Object 固定 "chat.completion.chunk"。
	Object string `json:"object,omitempty"`
	// Created 创建时间戳。
	Created int64 `json:"created,omitempty"`
	// Model 模型。
	Model string `json:"model"`
	// Choices chunk 选择列表。
	Choices []OpenAIChatCompletionChunkChoice `json:"choices"`
	// Usage 用量（仅当 stream_options.include_usage=true 时在最后一个 chunk 返回）。
	Usage *OpenAICompletionUsage `json:"usage,omitempty"`
	// ServiceTier 服务等级。
	ServiceTier OpenAIServiceTier `json:"service_tier,omitempty"`
	// SystemFingerprint 系统指纹。
	SystemFingerprint string `json:"system_fingerprint,omitempty"`
	// Moderation 审核结果（仅当请求中设置了 moderation 时返回）。
	Moderation *OpenAIChatCompletionResponseModeration `json:"moderation,omitempty"`
}

// OpenAIChatCompletionChunkChoice chunk 选择。
type OpenAIChatCompletionChunkChoice struct {
	// Index 选择索引。
	Index int64 `json:"index"`
	// ID 响应 ID。
	ID string `json:"id,omitempty"`
	// Delta 增量内容。
	Delta OpenAIChatCompletionChunkChoiceDelta `json:"delta"`
	// FinishReason 停止原因（最后一个 chunk 设置）。
	FinishReason *string `json:"finish_reason,omitempty"`
	// Logprobs Logprobs 信息。
	Logprobs *OpenAIChatCompletionChoiceLogprobs `json:"logprobs,omitempty"`
	// ContentFilterResults 内容过滤结果（Azure OpenAI / 国内厂商扩展）。
	ContentFilterResults *OpenAIChatCompletionContentFilterResults `json:"content_filter_results,omitempty"`
}

// OpenAIChatCompletionChunkChoiceDelta chunk 增量内容。
type OpenAIChatCompletionChunkChoiceDelta struct {
	// Role 角色（第一个 chunk 设置）。
	Role string `json:"role,omitempty"`
	// Content 文本增量。
	Content string `json:"content,omitempty"`
	// FunctionCall 旧版函数调用增量（已废弃）。
	FunctionCall *OpenAIChatCompletionMessageFunctionCallParam `json:"function_call,omitempty"`
	// Refusal 拒绝增量。
	Refusal *string `json:"refusal,omitempty"`
	// ToolCalls 工具调用增量列表。
	ToolCalls []OpenAIChatCompletionChunkChoiceDeltaToolCall `json:"tool_calls,omitempty"`
	// Audio 音频输出增量。
	Audio *OpenAIChatCompletionChunkChoiceDeltaAudio `json:"audio,omitempty"`
	// Reasoning 推理内容增量（DeepSeek 等厂商扩展）。
	Reasoning string `json:"reasoning,omitempty"`
	// ReasoningContent 推理内容增量（DeepSeek / Qwen 等厂商扩展）。
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

// OpenAIChatCompletionChunkChoiceDeltaToolCall 工具调用增量。
type OpenAIChatCompletionChunkChoiceDeltaToolCall struct {
	// Index 工具调用索引。
	Index int64 `json:"index"`
	// ID 工具调用 ID（第一个 chunk 设置）。
	ID string `json:"id,omitempty"`
	// Type 类型（第一个 chunk 设置）。
	Type string `json:"type,omitempty"`
	// Function 函数调用增量详情。
	Function *OpenAIChatCompletionMessageToolCallFunction `json:"function,omitempty"`
}

// OpenAIChatCompletionChunkChoiceDeltaAudio 音频输出增量。
type OpenAIChatCompletionChunkChoiceDeltaAudio struct {
	// ID 音频引用 ID。
	ID string `json:"id,omitempty"`
	// Data base64 编码音频数据（增量）。
	Data string `json:"data,omitempty"`
	// Transcript 音频转录文本（增量）。
	Transcript string `json:"transcript,omitempty"`
	// ExpiresAt 音频过期时间。
	ExpiresAt int64 `json:"expires_at,omitempty"`
}
