package gemini

import "encoding/json"

// GeminiGenerateContentResponse generateContent / streamGenerateContent 顶层响应。
//
// 每个流式帧（streamGenerateContent）都是一个完整的 GeminiGenerateContentResponse，
// 这是与 Anthropic / OpenAI 流事件协议的根本区别。
//
// RawData 字段用于同协议流式 passthrough：保留上游 data: 后的原始 JSON 字节，
// 自定义 MarshalJSON 优先返回 RawData（避免 typed struct 重新 marshal 剥离
// omitempty 零值字段 / 未知字段如 groundingMetadata / avgLogprobs 等）。
// 非流式路径不填 RawData，自动走 typed marshal，行为不变。
type GeminiGenerateContentResponse struct {
	RawData []byte `json:"-"` // 上游原始 data: 字节；passthrough 时填

	Candidates     []*GeminiCandidate    `json:"candidates,omitempty"`
	PromptFeedback *GeminiPromptFeedback `json:"promptFeedback,omitempty"`
	UsageMetadata  *GeminiUsageMetadata  `json:"usageMetadata,omitempty"`
	ModelVersion   string                `json:"modelVersion,omitempty"`
	ResponseID     string                `json:"responseId,omitempty"`
	CreateTime     string                `json:"createTime,omitempty"`
}

// MarshalJSON 优先返回上游原始字节（passthrough 路径），否则走 typed marshal。
//
// 触发 RawData 路径的条件：RawData 非空（流式 parse 已填）。
// 不影响非流式（RawData 为零值）与 typed unmarshal 后的二次序列化。
func (r *GeminiGenerateContentResponse) MarshalJSON() ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	if len(r.RawData) > 0 {
		return r.RawData, nil
	}
	type alias GeminiGenerateContentResponse
	return json.Marshal((*alias)(r))
}

// GeminiCandidate 单个候选结果。
type GeminiCandidate struct {
	Content       *GeminiContent        `json:"content,omitempty"`
	FinishReason  string                `json:"finishReason,omitempty"`
	Index         int32                 `json:"index,omitempty"`
	SafetyRatings []*GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

// GeminiPromptFeedback 整体 prompt 反馈（用于 block 诊断）。
type GeminiPromptFeedback struct {
	BlockReason        string                `json:"blockReason,omitempty"`
	BlockReasonMessage string                `json:"blockReasonMessage,omitempty"`
	SafetyRatings      []*GeminiSafetyRating `json:"safetyRatings,omitempty"`
}

// GeminiUsageMetadata 用量元数据（Google 协议命名约定）。
//
// 字段名严格保留 camelCase（与官方 SDK 一致），不做命名转换。
type GeminiUsageMetadata struct {
	PromptTokenCount        int32                       `json:"promptTokenCount"`
	CandidatesTokenCount    int32                       `json:"candidatesTokenCount"`
	TotalTokenCount         int32                       `json:"totalTokenCount"`
	CachedContentTokenCount int32                       `json:"cachedContentTokenCount"`
	ThoughtsTokenCount      int32                       `json:"thoughtsTokenCount"`
	ToolUsePromptTokenCount int32                       `json:"toolUsePromptTokenCount"`
	PromptTokensDetails     []*GeminiModalityTokenCount `json:"promptTokensDetails,omitempty"`
	CandidatesTokensDetails []*GeminiModalityTokenCount `json:"candidatesTokensDetails,omitempty"`
}

// GeminiModalityTokenCount 多模态分项统计（input/output text/image/audio 等）。
type GeminiModalityTokenCount struct {
	Modality   string `json:"modality,omitempty"`
	TokenCount int32  `json:"tokenCount"`
}
