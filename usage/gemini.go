package usage

import (
	"encoding/json"
)

// geminiUsageMetadata Gemini 协议 usageMetadata 子字段的本地投影。
//
// 故意不引用 pkg/protocol/gemini.GeminiUsageMetadata：pkg/protocol/gemini 通过
// streaming.go 反向引用 pkg/usage，引入会形成 import cycle。本地 struct 字段名严格
// 保留 camelCase（与 Google 协议命名约定一致），与 GeminiUsageMetadata 一一对应。
type geminiUsageMetadata struct {
	PromptTokenCount        int32 `json:"promptTokenCount"`
	CandidatesTokenCount    int32 `json:"candidatesTokenCount"`
	TotalTokenCount         int32 `json:"totalTokenCount"`
	CachedContentTokenCount int32 `json:"cachedContentTokenCount"`
	ThoughtsTokenCount      int32 `json:"thoughtsTokenCount"`
}

// FromGeminiBytes 从 Gemini GenerateContentResponse 原始响应字节提取 usage。
//
// 仅 unmarshal usageMetadata 子字段（普通 struct，无 union），绕过 Gemini Content union
// 的 unmarshal 问题（部分上游返回非标准 content 字段）。
//
// 解析失败或 usageMetadata 缺失时返回 nil；调用方应使用估算值兜底。
func FromGeminiBytes(data []byte) *Usage {
	if len(data) == 0 {
		return nil
	}
	var wrapper struct {
		UsageMetadata *geminiUsageMetadata `json:"usageMetadata"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil
	}
	if wrapper.UsageMetadata == nil {
		return nil
	}
	m := wrapper.UsageMetadata

	// Gemini 协议特性：cachedContentTokenCount ⊂ promptTokenCount
	// （官方 UsageMetadata 明确 "promptTokenCount includes cached content tokens"）。
	//   - Base.InputTokens = m.PromptTokenCount（**含** cached，协议原始值）
	//   - Base.FreshInputTokens = max(PromptTokenCount - CachedContentTokenCount, 0)（去 cache 后）
	//   - Base.TotalTokens = m.TotalTokenCount（透传上游）
	//   - Cache.CachedTokens = m.CachedContentTokenCount
	//
	// 注：Gemini 协议无 cache_creation 字段（caching 是隐式的），Cache.CreationInputTokens
	// 永远为 0；Cache.CachedTokens 与 OpenAI Chat Completions 的 cached_tokens 同语义。
	u := &Usage{
		Base: Base{
			InputTokens:      int64(m.PromptTokenCount),                                       // 协议原始值（含 cached）
			OutputTokens:     int64(m.CandidatesTokenCount),
			TotalTokens:      int64(m.TotalTokenCount),                                        // 透传上游
			FreshInputTokens: max(int64(m.PromptTokenCount-m.CachedContentTokenCount), 0),    // 去 cache 后
		},
		Protocol: "gemini",
	}
	if m.CachedContentTokenCount > 0 {
		u.Cache = &CacheMetrics{
			CachedTokens: int64(m.CachedContentTokenCount),
		}
	}
	if m.ThoughtsTokenCount > 0 {
		u.Thinking = &ExtendedThinking{
			ThinkingTokens: int64(m.ThoughtsTokenCount),
		}
	}
	return u
}
