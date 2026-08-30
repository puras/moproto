package gemini

// GeminiGenerationConfig 采样参数 + 输出控制。
//
// 所有数值采样字段均为 *T（指针）以区分"未设置"和"零值"——零值在 wire-format 上
// 通常合法（如 temperature=0 表示 greedy），调用方应使用 nil 表示未设置。
type GeminiGenerationConfig struct {
	Temperature      *float32              `json:"temperature,omitempty"`
	TopP             *float32              `json:"topP,omitempty"`
	TopK             *int32                `json:"topK,omitempty"`
	MaxOutputTokens  int32                 `json:"maxOutputTokens,omitempty"`
	StopSequences    []string              `json:"stopSequences,omitempty"`
	CandidateCount   int32                 `json:"candidateCount,omitempty"`
	ResponseMIMEType string                `json:"responseMimeType,omitempty"`
	ResponseSchema   *GeminiSchema         `json:"responseSchema,omitempty"`
	PresencePenalty  *float32              `json:"presencePenalty,omitempty"`
	FrequencyPenalty *float32              `json:"frequencyPenalty,omitempty"`
	Seed             *int32                `json:"seed,omitempty"`
	ThinkingConfig   *GeminiThinkingConfig `json:"thinkingConfig,omitempty"`
}

// GeminiThinkingConfig extended thinking 配置（Gemini 2.5 系列）。
type GeminiThinkingConfig struct {
	IncludeThoughts bool   `json:"includeThoughts,omitempty"`
	ThinkingBudget  *int32 `json:"thinkingBudget,omitempty"`
}
