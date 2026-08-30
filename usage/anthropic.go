package usage

import (
	"encoding/json"

	"hegq.net/moproto/protocol/anthropic"
)

// FromAnthropic 将 AnthropicUsage 转换为统一 Usage。
//
// 字段填充规则：
//   - Base.TotalTokens = 全上下文大小 = InputTokens + OutputTokens +
//     CacheCreationInputTokens + CacheReadInputTokens（Anthropic 上游不返回 TotalTokens，
//     按"模型实际处理的全部 token"语义自算；CumulativeCreation + Read 都是模型处理的输入）
//   - Cache: 任意 cache 字段 > 0 时创建分组并填充（注意 4.6+ 双形态：扁平顶层 +
//     嵌套 5m/1h 分桶共存，Ephemeral5m/1hInputTokens 是 CreationInputTokens 的子集，
//     不参与 TotalTokens 计算避免重复）
//   - Thinking / ServerTools / ServiceTier / InferenceGeo 按需填充
//
// nil 入参返回 nil，调用方无需 nil 检查。
func FromAnthropic(u *anthropic.AnthropicUsage) *Usage {
	if u == nil {
		return nil
	}
	result := &Usage{
		Base: Base{
			InputTokens:      int64(u.InputTokens),
			OutputTokens:     int64(u.OutputTokens),
			TotalTokens:      int64(u.InputTokens + u.OutputTokens +
				u.CacheCreationInputTokens + u.CacheReadInputTokens),
			FreshInputTokens: int64(u.InputTokens), // Anthropic 协议 input_tokens 本就不含 cache，两字段相等
		},
		Protocol: "anthropic",
	}

	// 缓存指标：顶层 cache_creation_input_tokens / cache_read_input_tokens 任一 > 0，或嵌套非 nil
	if u.CacheCreationInputTokens > 0 || u.CacheReadInputTokens > 0 || u.CacheCreation != nil {
		result.Cache = &CacheMetrics{
			CreationInputTokens: int64(u.CacheCreationInputTokens),
			ReadInputTokens:     int64(u.CacheReadInputTokens),
		}
		if u.CacheCreation != nil {
			result.Cache.Ephemeral1hInputTokens = int64(u.CacheCreation.Ephemeral1hInputTokens)
			result.Cache.Ephemeral5mInputTokens = int64(u.CacheCreation.Ephemeral5mInputTokens)
		}
	}

	// Thinking token
	if u.OutputTokensDetails != nil && u.OutputTokensDetails.ThinkingTokens > 0 {
		result.Thinking = &ExtendedThinking{
			ThinkingTokens: int64(u.OutputTokensDetails.ThinkingTokens),
		}
	}

	// 服务端工具
	if u.ServerToolUse != nil {
		result.ServerTools = &ServerToolUsage{
			WebFetchRequests:  u.ServerToolUse.WebFetchRequests,
			WebSearchRequests: u.ServerToolUse.WebSearchRequests,
		}
	}

	// 服务层级（响应侧：standard / priority / batch；Request 侧的 auto/standard_only 不进入此处）
	if u.ServiceTier != "" {
		result.ServiceTier = ServiceTier(u.ServiceTier)
	}

	// 推理地域
	if u.InferenceGeo != "" {
		result.InferenceGeo = u.InferenceGeo
	}

	return result
}

// FromAnthropicBytes 从原始响应 JSON 字节提取 usage。
//
// 绕过了整个 AnthropicMessageResponse 的 union interface unmarshal 问题
// （部分渠道如 MiniMax 返回非标准 content 字段，导致 Content unmarshal 失败）。
// usage 子字段是普通 struct，直接 unmarshal 无此问题。
//
// 解析失败返回 nil，调用方应使用估算值兜底。
func FromAnthropicBytes(data []byte) *Usage {
	if len(data) == 0 {
		return nil
	}
	var wrapper struct {
		Usage *anthropic.AnthropicUsage `json:"usage"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil
	}
	return FromAnthropic(wrapper.Usage)
}

// FromAnthropicMessageDelta 将 Anthropic 流式 message_delta 的 MessageDeltaUsage
// 转换为统一 Usage。
//
// 与 FromAnthropic 的差异：
//   - MessageDeltaUsage 是流式精简版（无 5m/1h CacheCreation 嵌套明细 + InferenceGeo +
//     ServiceTier，这些字段在 message_start 已发出）
//   - 保留 OutputTokensDetails.ThinkingTokens + ServerToolUse（流式可累加）
//   - Base.TotalTokens = 全上下文大小 = InputTokens + OutputTokens +
//     CacheCreationInputTokens + CacheReadInputTokens（Anthropic 上游不返回
//     TotalTokens，自算"模型实际处理的全部 token"语义）
//
// nil 入参返回 nil。
func FromAnthropicMessageDelta(u *anthropic.MessageDeltaUsage) *Usage {
	if u == nil {
		return nil
	}
	result := &Usage{
		Base: Base{
			InputTokens:      int64(u.InputTokens),
			OutputTokens:     int64(u.OutputTokens),
			TotalTokens:      int64(u.InputTokens + u.OutputTokens +
				u.CacheCreationInputTokens + u.CacheReadInputTokens),
			FreshInputTokens: int64(u.InputTokens), // 流式语义同顶层
		},
		Protocol: "anthropic",
	}

	// 缓存指标
	if u.CacheCreationInputTokens > 0 || u.CacheReadInputTokens > 0 {
		result.Cache = &CacheMetrics{
			CreationInputTokens: int64(u.CacheCreationInputTokens),
			ReadInputTokens:     int64(u.CacheReadInputTokens),
		}
	}

	// Thinking token
	if u.OutputTokensDetails != nil && u.OutputTokensDetails.ThinkingTokens > 0 {
		result.Thinking = &ExtendedThinking{
			ThinkingTokens: int64(u.OutputTokensDetails.ThinkingTokens),
		}
	}

	// 服务端工具
	if u.ServerToolUse != nil {
		result.ServerTools = &ServerToolUsage{
			WebFetchRequests:  u.ServerToolUse.WebFetchRequests,
			WebSearchRequests: u.ServerToolUse.WebSearchRequests,
		}
	}

	return result
}

// ToAnthropic 把网关统一的 *Usage 转换为 Anthropic 响应用 *AnthropicUsage。
//
// 用途：handler 在自构 Anthropic 协议响应（如 message_start 帧内嵌 usage）时，
// 把累积的 *Usage 投影回 Anthropic 协议层类型。
//
// 字段映射：
//   - Base.InputTokens / OutputTokens → AnthropicUsage.InputTokens / OutputTokens
//   - Cache.CreationInputTokens / ReadInputTokens → 顶层扁平字段（4.0/4.5 兼容）
//   - Cache.Ephemeral5mInputTokens / Ephemeral1hInputTokens → CacheCreation 嵌套（4.6+）
//   - Thinking.ThinkingTokens → OutputTokensDetails.ThinkingTokens 嵌套
//   - ServerTools.Web{Fetch,Search}Requests → ServerToolUse 嵌套
//   - ServiceTier / InferenceGeo 顶层字段
//
// 0 值字段 / nil 分组由 json omitempty 跳过（不进入 wire-format payload）。
//
// 始终返回 non-nil *AnthropicUsage（即使 u 为 nil），方便 handler 直接传入 frames.BuildMessageStartFrame
// （该函数要求 Usage 必填）。
func ToAnthropic(u *Usage) *anthropic.AnthropicUsage {
	if u == nil {
		return &anthropic.AnthropicUsage{}
	}
	out := &anthropic.AnthropicUsage{
		InputTokens:  int(u.InputTokens),
		OutputTokens: int(u.OutputTokens),
		ServiceTier:  anthropic.AnthropicUsageServiceTier(u.ServiceTier),
		InferenceGeo: u.InferenceGeo,
	}

	// Cache 分组：扁平 + 嵌套双形态输出。
	if u.Cache != nil {
		out.CacheCreationInputTokens = int(u.Cache.CreationInputTokens)
		out.CacheReadInputTokens = int(u.Cache.ReadInputTokens)
		if u.Cache.Ephemeral5mInputTokens > 0 || u.Cache.Ephemeral1hInputTokens > 0 {
			out.CacheCreation = &anthropic.AnthropicCacheCreation{
				Ephemeral1hInputTokens: int(u.Cache.Ephemeral1hInputTokens),
				Ephemeral5mInputTokens: int(u.Cache.Ephemeral5mInputTokens),
			}
		}
	}

	// Thinking 分组 → OutputTokensDetails 嵌套（extended thinking）。
	if u.Thinking != nil && u.Thinking.ThinkingTokens > 0 {
		out.OutputTokensDetails = &anthropic.AnthropicOutputTokensDetails{
			ThinkingTokens: int(u.Thinking.ThinkingTokens),
		}
	}

	// ServerTools 分组 → ServerToolUse 嵌套。
	if u.ServerTools != nil && (u.ServerTools.WebFetchRequests > 0 || u.ServerTools.WebSearchRequests > 0) {
		out.ServerToolUse = &anthropic.AnthropicServerToolUsage{
			WebFetchRequests:  u.ServerTools.WebFetchRequests,
			WebSearchRequests: u.ServerTools.WebSearchRequests,
		}
	}

	return out
}
