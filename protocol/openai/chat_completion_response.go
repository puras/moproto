package openai

import "encoding/json"

// OpenAIChatCompletion ChatCompletion 响应。
type OpenAIChatCompletion struct {
	// ID 响应唯一 ID（如 "chatcmpl-xxx"）。
	ID string `json:"id"`
	// Object 固定 "chat.completion"。
	Object string `json:"object,omitempty"`
	// Created 创建时间（Unix 时间戳）。
	Created int64 `json:"created,omitempty"`
	// Model 实际使用的模型。
	Model string `json:"model"`
	// Choices 完成列表。
	Choices []OpenAIChatCompletionChoice `json:"choices"`
	// Usage 用量统计（仅在 stream_options.include_usage=true 或非流式时返回）。
	Usage *OpenAICompletionUsage `json:"usage,omitempty"`
	// ServiceTier 服务等级。
	ServiceTier OpenAIServiceTier `json:"service_tier,omitempty"`
	// SystemFingerprint 系统指纹（用于路由决策追踪）。
	SystemFingerprint string `json:"system_fingerprint,omitempty"`
	// Moderation 审核结果（仅当请求中设置了 moderation 时返回）。
	Moderation *OpenAIChatCompletionResponseModeration `json:"moderation,omitempty"`
}

// OpenAIChatCompletionChoice 单个完成。
type OpenAIChatCompletionChoice struct {
	// Index 完成索引。
	Index int64 `json:"index"`
	// Message 完成消息。
	Message OpenAIChatCompletionChoiceMessage `json:"message"`
	// FinishReason 停止原因。
	FinishReason string `json:"finish_reason"`
	// Logprobs Logprobs 信息。
	Logprobs *OpenAIChatCompletionChoiceLogprobs `json:"logprobs,omitempty"`
	// ContentFilterResults 内容过滤结果（Azure OpenAI / 国内厂商扩展）。
	ContentFilterResults *OpenAIChatCompletionContentFilterResults `json:"content_filter_results,omitempty"`
}

// OpenAIChatCompletionChoiceMessage 完成消息。
type OpenAIChatCompletionChoiceMessage struct {
	// Role 固定 "assistant"。
	Role string `json:"role"`
	// Content 文本内容（可能为空字符串当有 tool_calls 时）。
	Content string `json:"content"`
	// Refusal 拒绝消息。
	Refusal *string `json:"refusal,omitempty"`
	// ToolCalls 工具调用列表。
	ToolCalls []OpenAIChatCompletionChoiceMessageToolCallUnion `json:"tool_calls,omitempty"`
	// Audio 助手消息音频输出（4.6+）。
	Audio *OpenAIChatCompletionChoiceMessageAudio `json:"audio,omitempty"`
	// Annotations 注释列表（如 URL 引用）。
	Annotations []OpenAIChatCompletionChoiceMessageAnnotation `json:"annotations,omitempty"`
	// FunctionCall 旧版函数调用（已废弃）。
	FunctionCall *OpenAIChatCompletionMessageFunctionCallParam `json:"function_call,omitempty"`
	// Reasoning 推理内容（DeepSeek 等厂商扩展；与 reasoning_content 字段等价）。
	Reasoning string `json:"reasoning,omitempty"`
	// ReasoningContent 推理内容（DeepSeek / Qwen 等厂商扩展）。
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

// OpenAIChatCompletionResponseModeration 响应侧审核结果。
type OpenAIChatCompletionResponseModeration struct {
	// Input 输入审核结果。
	Input *OpenAIChatCompletionModerationInput `json:"input,omitempty"`
	// Output 输出审核结果。
	Output *OpenAIChatCompletionModerationOutput `json:"output,omitempty"`
}

// OpenAIChatCompletionModerationInput 输入审核结果。
type OpenAIChatCompletionModerationInput struct {
	// Flagged 是否命中审核。
	Flagged bool `json:"flagged"`
	// Categories 各分类是否命中。
	Categories map[string]bool `json:"categories,omitempty"`
	// CategoryScores 各分类得分。
	CategoryScores map[string]float64 `json:"category_scores,omitempty"`
	// Model 审核模型。
	Model string `json:"model,omitempty"`
}

// OpenAIChatCompletionModerationOutput 输出审核结果。
type OpenAIChatCompletionModerationOutput struct {
	// Flagged 是否命中审核。
	Flagged bool `json:"flagged"`
	// Categories 各分类是否命中。
	Categories map[string]bool `json:"categories,omitempty"`
	// CategoryScores 各分类得分。
	CategoryScores map[string]float64 `json:"category_scores,omitempty"`
	// Model 审核模型。
	Model string `json:"model,omitempty"`
}

// OpenAIChatCompletionContentFilterResults 内容过滤结果（Azure OpenAI 扩展）。
type OpenAIChatCompletionContentFilterResults struct {
	// Hate 仇恨内容。
	Hate *OpenAIChatCompletionContentFilterResult `json:"hate,omitempty"`
	// HateThreatening 仇恨威胁内容。
	HateThreatening *OpenAIChatCompletionContentFilterResult `json:"hate_threatening,omitempty"`
	// SelfHarm 自残内容。
	SelfHarm *OpenAIChatCompletionContentFilterResult `json:"self_harm,omitempty"`
	// Sexual 性相关内容。
	Sexual *OpenAIChatCompletionContentFilterResult `json:"sexual,omitempty"`
	// Violence 暴力内容。
	Violence *OpenAIChatCompletionContentFilterResult `json:"violence,omitempty"`
	// Jailbreak 越狱尝试。
	Jailbreak *OpenAIChatCompletionContentFilterResult `json:"jailbreak,omitempty"`
	// Profanity 脏话。
	Profanity *OpenAIChatCompletionContentFilterResult `json:"profanity,omitempty"`
}

// OpenAIChatCompletionContentFilterResult 单项内容过滤结果。
type OpenAIChatCompletionContentFilterResult struct {
	// Filtered 是否被过滤。
	Filtered bool `json:"filtered"`
	// Severity 严重程度。
	Severity string `json:"severity,omitempty"`
}

// OpenAIChatCompletionChoiceMessageAudio 完成消息音频。
type OpenAIChatCompletionChoiceMessageAudio struct {
	// ID 音频引用 ID。
	ID string `json:"id"`
	// Data base64 编码的音频数据。
	Data string `json:"data"`
	// ExpiresAt 音频过期时间。
	ExpiresAt int64 `json:"expires_at"`
	// Transcript 音频转录文本。
	Transcript string `json:"transcript"`
}

// OpenAIChatCompletionChoiceMessageAnnotation 完成消息注释。
type OpenAIChatCompletionChoiceMessageAnnotation struct {
	// Type 注释类型（如 "url_citation"）。
	Type string `json:"type"`
	// URL 引用 URL（如 "url_citation"）。
	URLCitation *OpenAIChatCompletionMessageAnnotationURLCitation `json:"url_citation,omitempty"`
}

// OpenAIChatCompletionMessageAnnotationURLCitation URL 引用注释。
type OpenAIChatCompletionMessageAnnotationURLCitation struct {
	// URL 引用的 URL。
	URL string `json:"url"`
	// Title 引用标题。
	Title string `json:"title"`
	// StartIndex 内容中引用起始位置。
	StartIndex int64 `json:"start_index"`
	// EndIndex 内容中引用结束位置。
	EndIndex int64 `json:"end_index"`
}

// OpenAIChatCompletionChoiceLogprobs Logprobs 信息。
type OpenAIChatCompletionChoiceLogprobs struct {
	// Content token logprobs 列表。
	Content []OpenAIChatCompletionTokenLogprob `json:"content,omitempty"`
	// Refusal 拒绝 token logprobs。
	Refusal []json.RawMessage `json:"refusal,omitempty"`
}

// OpenAIChatCompletionTokenLogprob 单 token logprob。
type OpenAIChatCompletionTokenLogprob struct {
	// Token token 字符串。
	Token string `json:"token"`
	// Bytes token 字节表示。
	Bytes []int64 `json:"bytes,omitempty"`
	// Logprob logprob 值。
	Logprob float64 `json:"logprob"`
	// TopLogprobs top-N logprobs。
	TopLogprobs []OpenAIChatCompletionTokenLogprobTop `json:"top_logprobs"`
}

// OpenAIChatCompletionTokenLogprobTop top logprob。
type OpenAIChatCompletionTokenLogprobTop struct {
	Token   string  `json:"token"`
	Bytes   []int64 `json:"bytes,omitempty"`
	Logprob float64 `json:"logprob"`
}

// OpenAICompletionUsage 用量统计。
type OpenAICompletionUsage struct {
	// PromptTokens 输入 token 数。
	PromptTokens int64 `json:"prompt_tokens"`
	// CompletionTokens 输出 token 数。
	CompletionTokens int64 `json:"completion_tokens"`
	// TotalTokens 总 token 数。
	TotalTokens int64 `json:"total_tokens"`
	// PromptTokensDetails 输入 token 详细信息。
	PromptTokensDetails *OpenAICompletionUsagePromptTokensDetails `json:"prompt_tokens_details,omitempty"`
	// CompletionTokensDetails 输出 token 详细信息。
	CompletionTokensDetails *OpenAICompletionUsageCompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

// OpenAICompletionUsagePromptTokensDetails 输入 token 详情。
type OpenAICompletionUsagePromptTokensDetails struct {
	// CachedTokens 缓存命中 token 数。
	CachedTokens int64 `json:"cached_tokens,omitempty"`
	// AudioTokens 音频输入 token 数。
	AudioTokens int64 `json:"audio_tokens,omitempty"`
}

// OpenAICompletionUsageCompletionTokensDetails 输出 token 详情。
type OpenAICompletionUsageCompletionTokensDetails struct {
	// ReasoningTokens 推理 token 数。
	ReasoningTokens int64 `json:"reasoning_tokens,omitempty"`
	// AudioTokens 音频输出 token 数。
	AudioTokens int64 `json:"audio_tokens,omitempty"`
	// AcceptedPredictionTokens 接受的预测 token 数。
	AcceptedPredictionTokens int64 `json:"accepted_prediction_tokens,omitempty"`
	// RejectedPredictionTokens 拒绝的预测 token 数。
	RejectedPredictionTokens int64 `json:"rejected_prediction_tokens,omitempty"`
}
