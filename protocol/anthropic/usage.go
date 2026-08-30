package anthropic

// AnthropicUsageServiceTier usage.service_tier 字段类型。
//
// 标识本次请求所命中的服务层级（按定价/可用性区分），用于计费与配额分账。
//
// 响应侧取值空间：standard / priority / batch。
// Request 侧（顶层 ServiceTier 字段）取值空间：auto / standard_only。
// 两者是不同类型，互不冲突。
type AnthropicUsageServiceTier string

const (
	// AnthropicUsageServiceTierStandard 标准服务层级（默认）。
	// 对齐官方 UsageServiceTierStandard（message.go:8495）。
	AnthropicUsageServiceTierStandard AnthropicUsageServiceTier = "standard"
	// AnthropicUsageServiceTierPriority 优先服务层级（更快但溢价）。
	// 对齐官方 UsageServiceTierPriority（message.go:8496）。
	AnthropicUsageServiceTierPriority AnthropicUsageServiceTier = "priority"
	// AnthropicUsageServiceTierBatch 批处理层级（离线任务，折扣计价）。
	// 对齐官方 UsageServiceTierBatch（message.go:8497）。
	AnthropicUsageServiceTierBatch AnthropicUsageServiceTier = "batch"

	// ---- 兼容旧 wire-format（去前缀别名，便于习惯） ----
	// AnthropicUsageServiceTierAuto 由系统自动选择（标准优先级）。
	AnthropicUsageServiceTierAuto AnthropicUsageServiceTier = "auto"
	// AnthropicUsageServiceTierStandardOnly 强制仅使用 standard 层级。
	AnthropicUsageServiceTierStandardOnly AnthropicUsageServiceTier = "standard_only"
)

// AnthropicCacheCreation prompt caching 写入细节（嵌套于 usage.cache_creation）。
//
// 4.6+ 支持按 5m（短 TTL）与 1h（长 TTL）分别计量缓存写入 token，
// 用于区分短期/长期缓存写成本。
type AnthropicCacheCreation struct {
	Ephemeral1hInputTokens int `json:"ephemeral_1h_input_tokens,omitempty"`
	Ephemeral5mInputTokens int `json:"ephemeral_5m_input_tokens,omitempty"`
}

// AnthropicOutputTokensDetails output_tokens 细节（嵌套于 usage.output_tokens_details）。
//
// ThinkingTokens 用于启用 extended thinking 的模型（如 Sonnet 4.6+），
// 记录其中归类为"思考"的 token 数量，便于分账/审计。
type AnthropicOutputTokensDetails struct {
	ThinkingTokens int `json:"thinking_tokens,omitempty"`
}

// AnthropicServerToolUsage 服务端工具调用计数（嵌套于 usage.server_tool_use）。
//
// 统计由 Anthropic 服务端代为执行的工具（如 web_search / web_fetch）调用次数，
// 用于按工具维度计费与限额追踪。
type AnthropicServerToolUsage struct {
	WebFetchRequests  int `json:"web_fetch_requests,omitempty"`
	WebSearchRequests int `json:"web_search_requests,omitempty"`
}

// AnthropicUsage 顶层 usage 字段。
//
// 汇总本次请求的输入/输出/缓存/服务端工具/服务层级/推理地域等指标。
// CacheCreation / OutputTokensDetails / ServerToolUse 为嵌套对象，
// 分别承载更细粒度的统计（5m/1h 缓存桶、thinking token 数、web_search/fetch 调用次数）。
type AnthropicUsage struct {
	InputTokens              int                           `json:"input_tokens"`
	OutputTokens             int                           `json:"output_tokens"`
	CacheCreationInputTokens int                           `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int                           `json:"cache_read_input_tokens,omitempty"`
	CacheCreation            *AnthropicCacheCreation       `json:"cache_creation,omitempty"`
	OutputTokensDetails      *AnthropicOutputTokensDetails `json:"output_tokens_details,omitempty"`
	ServerToolUse            *AnthropicServerToolUsage     `json:"server_tool_use,omitempty"`
	ServiceTier              AnthropicUsageServiceTier     `json:"service_tier,omitempty"`
	InferenceGeo             string                        `json:"inference_geo,omitempty"`
}
