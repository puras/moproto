package usage

import (
	"encoding/json"

	pkgopenai "github.com/puras/moproto/protocol/openai"
)

// FromOpenAICompletion 将 OpenAICompletionUsage 转换为统一 Usage。
//
// 字段填充规则：
//   - Base.InputTokens = u.PromptTokens（**含** cached_tokens，协议原始值）
//   - Base.FreshInputTokens = max(u.PromptTokens - cached, 0)（去 cache 后的计费输入）
//   - Base.OutputTokens = u.CompletionTokens
//   - Base.TotalTokens = u.TotalTokens（原样透传，OpenAI 上游 total_tokens 已含 cached）
//   - Cache.CachedTokens = u.PromptTokensDetails.CachedTokens
//
// OpenAI 协议特性：cached_tokens ⊂ prompt_tokens（不是 Anthropic 的 cache_read ⊥ input）。
// 计费时 Input 桶用 FreshInputTokens（已剔除 cached），CacheRead 桶用 Cache.CachedTokens，
// 两桶不重叠——避免 CalcCostUsage 7 桶乘法重复计费。
// 展示 InputTokens 时按 FreshInputTokens（已剔除缓存命中的"非缓存输入"）。
//
// 注意：OpenAI 协议无 cache_creation_input_tokens 字段（OpenAI prompt caching 是隐式管理），
// 所以 CacheMetrics 的 CreationInputTokens / ReadInputTokens 在此路径下始终为 0。
//
// nil 入参返回 nil，调用方无需 nil 检查。
func FromOpenAICompletion(u *pkgopenai.OpenAICompletionUsage) *Usage {
	if u == nil {
		return nil
	}
	var cached int64
	if u.PromptTokensDetails != nil {
		cached = u.PromptTokensDetails.CachedTokens
	}
	result := &Usage{
		Base: Base{
			InputTokens:      u.PromptTokens,                 // 协议原始值（含 cached）
			OutputTokens:     u.CompletionTokens,
			TotalTokens:      u.TotalTokens,                  // 透传上游
			FreshInputTokens: max(u.PromptTokens-cached, 0),  // 去 cache 后
		},
		Protocol: "openai",
	}

	// 缓存指标：OpenAI 仅有 CachedTokens + AudioTokens(input)
	if u.PromptTokensDetails != nil {
		cache := &CacheMetrics{
			CachedTokens: u.PromptTokensDetails.CachedTokens,
		}
		if cache.CachedTokens > 0 || u.PromptTokensDetails.AudioTokens > 0 {
			result.Cache = cache
		}
	}

	// 推理 token
	if u.CompletionTokensDetails != nil && u.CompletionTokensDetails.ReasoningTokens > 0 {
		result.Reasoning = &ReasoningTokens{
			ReasoningTokens: u.CompletionTokensDetails.ReasoningTokens,
		}
	}

	// 预测 token
	if u.CompletionTokensDetails != nil &&
		(u.CompletionTokensDetails.AcceptedPredictionTokens > 0 ||
			u.CompletionTokensDetails.RejectedPredictionTokens > 0) {
		result.Prediction = &PredictionTokens{
			AcceptedTokens: u.CompletionTokensDetails.AcceptedPredictionTokens,
			RejectedTokens: u.CompletionTokensDetails.RejectedPredictionTokens,
		}
	}

	// 音频 token（input + output）
	if (u.PromptTokensDetails != nil && u.PromptTokensDetails.AudioTokens > 0) ||
		(u.CompletionTokensDetails != nil && u.CompletionTokensDetails.AudioTokens > 0) {
		result.Audio = &AudioTokens{}
		if u.PromptTokensDetails != nil {
			result.Audio.InputAudioTokens = u.PromptTokensDetails.AudioTokens
		}
		if u.CompletionTokensDetails != nil {
			result.Audio.OutputAudioTokens = u.CompletionTokensDetails.AudioTokens
		}
	}

	return result
}

// FromOpenAIBytes 从原始响应 JSON 字节提取 usage。
//
// OpenAI 响应中 usage 是普通 struct，直接 unmarshal 无 union interface 问题。
// 本函数从原始响应字节中提取 usage 字段，解析失败返回 nil。
func FromOpenAIBytes(data []byte) *Usage {
	if len(data) == 0 {
		return nil
	}
	var wrapper struct {
		Usage *pkgopenai.OpenAICompletionUsage `json:"usage"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil
	}
	return FromOpenAICompletion(wrapper.Usage)
}

// FromOpenAIResponsesUsage 将 OpenAIResponseUsage 转换为统一 Usage。
//
// 与 FromOpenAICompletion 的差异（OpenAI Responses API 字段命名不同）：
//   - u.InputTokens → Base.InputTokens（协议原始值，**含** cached_tokens）
//   - max(InputTokens - CachedTokens, 0) → Base.FreshInputTokens（计费输入）
//   - u.OutputTokens → Base.OutputTokens
//   - u.TotalTokens → Base.TotalTokens（透传上游）
//   - InputTokensDetails.CachedTokens → Cache.CachedTokens
//   - OutputTokensDetails.ReasoningTokens → Reasoning.ReasoningTokens
//
// 语义同 FromOpenAICompletion：cached ⊂ input，Input 桶用 FreshInputTokens。
//
// nil 入参返回 nil，调用方无需 nil 检查。
func FromOpenAIResponsesUsage(u *pkgopenai.OpenAIResponseUsage) *Usage {
	if u == nil {
		return nil
	}
	result := &Usage{
		Base: Base{
			InputTokens:      u.InputTokens,
			OutputTokens:     u.OutputTokens,
			TotalTokens:      u.TotalTokens,
			FreshInputTokens: max(u.InputTokens-u.InputTokensDetails.CachedTokens, 0),
		},
		Protocol: "openai",
	}

	// 缓存命中（OpenAI Responses API 沿用 CachedTokens 命名）
	if u.InputTokensDetails.CachedTokens > 0 {
		result.Cache = &CacheMetrics{
			CachedTokens: u.InputTokensDetails.CachedTokens,
		}
	}

	// 推理 token（reasoning model 场景：o-series / GPT-5 等）
	if u.OutputTokensDetails.ReasoningTokens > 0 {
		result.Reasoning = &ReasoningTokens{
			ReasoningTokens: u.OutputTokensDetails.ReasoningTokens,
		}
	}

	return result
}

// FromOpenAIResponsesBytes 从原始响应 JSON 字节提取 Responses API usage。
//
// 与 FromOpenAIBytes 的差异：本函数兼容 Responses API 的两种 usage 位置：
//   - 非流响应顶层：{"usage": {...}}（OpenAIResponse 类型）
//   - 流式 response.completed 事件：{"response": {"usage": {...}}}（OpenAIResponseCompletedEvent）
//
// 解析失败返回 nil。优先尝试嵌套位置（流式事件形态），失败回退顶层位置。
func FromOpenAIResponsesBytes(data []byte) *Usage {
	if len(data) == 0 {
		return nil
	}
	// 优先尝试流式事件形态：{"response": {"usage": {...}}}
	var wrapperStream struct {
		Response *struct {
			Usage *pkgopenai.OpenAIResponseUsage `json:"usage"`
		} `json:"response"`
	}
	if err := json.Unmarshal(data, &wrapperStream); err == nil &&
		wrapperStream.Response != nil && wrapperStream.Response.Usage != nil {
		return FromOpenAIResponsesUsage(wrapperStream.Response.Usage)
	}
	// 兜底：非流响应顶层 {"usage": {...}}
	var wrapperTop struct {
		Usage *pkgopenai.OpenAIResponseUsage `json:"usage"`
	}
	if err := json.Unmarshal(data, &wrapperTop); err != nil {
		return nil
	}
	return FromOpenAIResponsesUsage(wrapperTop.Usage)
}
