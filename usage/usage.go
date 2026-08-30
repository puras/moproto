package usage

// ServiceTier 服务层级。响应侧取值：standard / priority / batch。
//
// 仅承载计费相关层级（与 AnthropicUsageServiceTier 响应侧同语义）。
// Request 侧的 ServiceTier（取值空间 auto / standard_only）不入 Usage，
// 属于入参而非用量统计。
type ServiceTier string

const (
	ServiceTierStandard ServiceTier = "standard"
	ServiceTierPriority ServiceTier = "priority"
	ServiceTierBatch    ServiceTier = "batch"
)

// Base 核心计费字段，所有协议共享。
//
// 命名约定：采用 Anthropic 命名（InputTokens / OutputTokens）作为唯一权威字段。
// OpenAI 风格的 PromptTokens / CompletionTokens 不在此处出现，转换函数
// FromOpenAICompletion 内部完成映射。
//
// 三层输入 token 口径：
//   - InputTokens：协议原始值。
//       OpenAI Chat/Responses：u.PromptTokens / u.InputTokens（**含** cached_tokens）。
//       Gemini：m.PromptTokenCount（**含** cachedContentTokenCount）。
//       Anthropic：u.InputTokens（**不含** cache 部分，Anthropic 协议特性）。
//   - FreshInputTokens：实际计费的输入 token（去 cache 后）。
//       OpenAI/Gemini：max(InputTokens - CachedTokens, 0)。
//       Anthropic：等于 InputTokens（协议本身已不含 cache）。
//       计费公式（BillableInput）+ Input 桶计费 + 价格区间匹配基线均使用此字段。
//   - TotalTokens：OpenAI/Gemini 透传上游 u.TotalTokens；Anthropic 因上游不返回，
//     按"全上下文大小"自算 = Input + Output + CacheCreation + CacheRead。
//
// 字段互斥性（计费桶假设）：FreshInputTokens ⊥ CacheReadInputTokens ⊥ CacheCreationInputTokens，
// 即三个桶 token 不重叠，PriceResolver 的 CalcCostUsage 直接相加而不去重。
type Base struct {
	// InputTokens 输入 token 数（协议原始值）。
	//
	// OpenAI/Gemini 含 cached；Anthropic 不含 cache 部分。
	InputTokens int64 `json:"input_tokens"`
	// OutputTokens 输出 token 数。
	OutputTokens int64 `json:"output_tokens"`
	// TotalTokens 总 token 数。
	//
	// 填充语义：
	//   - FromAnthropic：自算 InputTokens + OutputTokens + CacheCreation + CacheRead
	//     （上游不返回 TotalTokens，按"模型实际处理的全部 token"语义计算）
	//   - FromOpenAICompletion / FromOpenAIResponsesUsage：原样透传上游 u.TotalTokens
	//     （OpenAI 上游 total_tokens 已含 cached_tokens，无需本地重算）
	//   - FromGeminiBytes：原样透传上游 m.TotalTokenCount
	TotalTokens int64 `json:"total_tokens"`
	// FreshInputTokens 计费输入 token 数（去 cache 后的口径）。
	//
	// 与 InputTokens 的区别：
	//   - InputTokens 是协议原始值（OpenAI/Gemini 含 cached_tokens）
	//   - FreshInputTokens 是实际计费的输入（OpenAI/Gemini 已扣除 cached_tokens）
	//   - Anthropic 协议下两者相等（Anthropic input_tokens 本来就不含 cache）
	//
	// 用途：
	//   - BillableInput() 计费公式基线
	//   - PriceResolver.toCostTokens 的 Input 桶
	//   - 持久化到 mo_usage_record.fresh_input_tokens 列供前端展示
	FreshInputTokens int64 `json:"fresh_input_tokens"`
}

// CacheMetrics 缓存相关指标。
//
// 支持：
//   - Anthropic：cache_creation_input_tokens / cache_read_input_tokens + 5m/1h 分桶
//   - OpenAI：cached_tokens
type CacheMetrics struct {
	// Anthropic 缓存写入
	CreationInputTokens int64 `json:"cache_creation_input_tokens,omitempty"`
	// Anthropic 缓存读取
	ReadInputTokens int64 `json:"cache_read_input_tokens,omitempty"`
	// OpenAI 缓存命中
	CachedTokens int64 `json:"cached_tokens,omitempty"`

	// Anthropic 缓存 TTL 分桶（4.6+）
	// 短 TTL 缓存（5 分钟）
	Ephemeral5mInputTokens int64 `json:"ephemeral_5m_input_tokens,omitempty"`
	// 长 TTL 缓存（1 小时）
	Ephemeral1hInputTokens int64 `json:"ephemeral_1h_input_tokens,omitempty"`
}

// Total 返回缓存相关总 token 数（计费用，便于读取）。
func (c *CacheMetrics) Total() int64 {
	return c.CreationInputTokens + c.ReadInputTokens + c.CachedTokens +
		c.Ephemeral5mInputTokens + c.Ephemeral1hInputTokens
}

// ReasoningTokens 推理相关 token（OpenAI o 系列、DeepSeek 等）。
type ReasoningTokens struct {
	// ReasoningTokens 推理过程消耗的 token 数。
	ReasoningTokens int64 `json:"reasoning_tokens,omitempty"`
}

// ExtendedThinking 扩展思考相关（Anthropic extended thinking）。
type ExtendedThinking struct {
	// ThinkingTokens 思考过程消耗的 token 数。
	ThinkingTokens int64 `json:"thinking_tokens,omitempty"`
}

// ServerToolUsage 服务端工具调用（Anthropic web_search / web_fetch）。
type ServerToolUsage struct {
	// WebFetchRequests web_fetch 调用次数。
	WebFetchRequests int `json:"web_fetch_requests,omitempty"`
	// WebSearchRequests web_search 调用次数。
	WebSearchRequests int `json:"web_search_requests,omitempty"`
}

// AudioTokens 音频 token（OpenAI audio 模型）。
type AudioTokens struct {
	// InputAudioTokens 音频输入 token 数。
	InputAudioTokens int64 `json:"input_audio_tokens,omitempty"`
	// OutputAudioTokens 音频输出 token 数。
	OutputAudioTokens int64 `json:"output_audio_tokens,omitempty"`
}

// PredictionTokens 预测相关（OpenAI predictions 参数）。
type PredictionTokens struct {
	// AcceptedTokens 接受的预测 token 数。
	AcceptedTokens int64 `json:"accepted_prediction_tokens,omitempty"`
	// RejectedTokens 拒绝的预测 token 数。
	RejectedTokens int64 `json:"rejected_prediction_tokens,omitempty"`
}

// Usage 统一用量结构，协议无关。
//
// 使用场景：
//   - 计费系统扣费计算
//   - 用量统计与报表
//   - 渠道成本核算
//
// 字段填充原则：
//   - Base 三字段必填
//   - 嵌套分组（Cache / Reasoning / Thinking / ServerTools / Audio / Prediction）按需填充，
//     nil 表示协议不支持该维度，调用方通过 nil 检查决定是否计费
type Usage struct {
	Base

	Cache       *CacheMetrics     `json:"cache,omitempty"`
	Reasoning   *ReasoningTokens  `json:"reasoning,omitempty"`
	Thinking    *ExtendedThinking `json:"thinking,omitempty"`
	ServerTools *ServerToolUsage  `json:"server_tools,omitempty"`
	Audio       *AudioTokens      `json:"audio,omitempty"`
	Prediction  *PredictionTokens `json:"prediction,omitempty"`

	// === 多模态用量（PR12+ video/image/audio 接入）===
	// 字段单位：
	//   - ImageCount       张数
	//   - AudioInputChars  音频输入字符数
	//   - AudioOutputChars 音频输出字符数
	//   - VideoDurationSec 视频生成时长（秒）
	ImageCount       int64 `json:"image_count,omitempty"`
	AudioInputChars  int64 `json:"audio_input_chars,omitempty"`
	AudioOutputChars int64 `json:"audio_output_chars,omitempty"`
	VideoDurationSec int64 `json:"video_duration_sec,omitempty"`

	ServiceTier  ServiceTier `json:"service_tier,omitempty"`
	InferenceGeo string      `json:"inference_geo,omitempty"`

	// Protocol 来源协议标识：anthropic / openai / deepseek / ...
	Protocol string `json:"protocol,omitempty"`
}

// IsZero 检查是否为空用量（Base + 所有分组 + 多模态字段均无值）。
//
// 单独存在 Cache.CachedTokens=1000 但 Base 全 0 的场景属于合法用法
// （只统计缓存命中），因此遍历所有分组。
// 多模态场景同样：纯视频/图片/音频请求可能 Base 全 0 但 multimodal 非零。
func (u *Usage) IsZero() bool {
	if u == nil {
		return true
	}
	if u.InputTokens != 0 || u.OutputTokens != 0 || u.TotalTokens != 0 || u.FreshInputTokens != 0 {
		return false
	}
	if u.ImageCount != 0 || u.AudioInputChars != 0 || u.AudioOutputChars != 0 || u.VideoDurationSec != 0 {
		return false
	}
	return u.Cache == nil && u.Reasoning == nil && u.Thinking == nil &&
		u.ServerTools == nil && u.Audio == nil && u.Prediction == nil
}

// Add 聚合 other 到当前 Usage（原地修改）。
//
// 用于 protocol/stream_handler.go 流式响应中逐 delta 合并用量。
// nil receiver / nil other 自动处理：nil receiver 等价于 zero Usage；
// nil other 直接返回。
//
// 聚合规则：
//   - Base 三字段直接相加
//   - 嵌套分组：双方 nil → 保持 nil；一方非 nil → 创建/复用分组并按字段相加
func (u *Usage) Add(other *Usage) {
	if u == nil || other == nil {
		return
	}
	u.InputTokens += other.InputTokens
	u.OutputTokens += other.OutputTokens
	u.TotalTokens += other.TotalTokens
	u.FreshInputTokens += other.FreshInputTokens

	if other.Cache != nil {
		if u.Cache == nil {
			u.Cache = &CacheMetrics{}
		}
		u.Cache.CreationInputTokens += other.Cache.CreationInputTokens
		u.Cache.ReadInputTokens += other.Cache.ReadInputTokens
		u.Cache.CachedTokens += other.Cache.CachedTokens
		u.Cache.Ephemeral5mInputTokens += other.Cache.Ephemeral5mInputTokens
		u.Cache.Ephemeral1hInputTokens += other.Cache.Ephemeral1hInputTokens
	}

	if other.Reasoning != nil {
		if u.Reasoning == nil {
			u.Reasoning = &ReasoningTokens{}
		}
		u.Reasoning.ReasoningTokens += other.Reasoning.ReasoningTokens
	}

	if other.Thinking != nil {
		if u.Thinking == nil {
			u.Thinking = &ExtendedThinking{}
		}
		u.Thinking.ThinkingTokens += other.Thinking.ThinkingTokens
	}

	if other.ServerTools != nil {
		if u.ServerTools == nil {
			u.ServerTools = &ServerToolUsage{}
		}
		u.ServerTools.WebFetchRequests += other.ServerTools.WebFetchRequests
		u.ServerTools.WebSearchRequests += other.ServerTools.WebSearchRequests
	}

	if other.Audio != nil {
		if u.Audio == nil {
			u.Audio = &AudioTokens{}
		}
		u.Audio.InputAudioTokens += other.Audio.InputAudioTokens
		u.Audio.OutputAudioTokens += other.Audio.OutputAudioTokens
	}

	if other.Prediction != nil {
		if u.Prediction == nil {
			u.Prediction = &PredictionTokens{}
		}
		u.Prediction.AcceptedTokens += other.Prediction.AcceptedTokens
		u.Prediction.RejectedTokens += other.Prediction.RejectedTokens
	}

	// 多模态用量累加（PR12+ video/image/audio 接入）
	u.ImageCount += other.ImageCount
	u.AudioInputChars += other.AudioInputChars
	u.AudioOutputChars += other.AudioOutputChars
	u.VideoDurationSec += other.VideoDurationSec

	// ServiceTier / InferenceGeo / Protocol 非累加语义，后者覆盖前者（流式通常只有一次 metadata）
	if other.ServiceTier != "" {
		u.ServiceTier = other.ServiceTier
	}
	if other.InferenceGeo != "" {
		u.InferenceGeo = other.InferenceGeo
	}
	if other.Protocol != "" {
		u.Protocol = other.Protocol
	}
}

// HasInput 检查是否存在任意形式的输入 token。
//
// 以下任一条件即视为上游返回了有效输入：
//   - InputTokens > 0（Anthropic 官方语义：input_tokens 已包含 cache 部分）
//   - Cache.ReadInputTokens > 0（Anthropic cache_read / OpenAI cached_tokens 缓存命中）
//   - Cache.CreationInputTokens > 0（Anthropic cache_creation 缓存写入；
//     非标上游如 MiniMax 会把全部输入 token 计入 cache_creation 而非 input_tokens）
//   - Thinking.ThinkingTokens > 0（extended thinking 必有非零输入）
//
// 用于 ValidUsage 判定：上游返回任意输入 token 即视为有效响应。
func (u *Usage) HasInput() bool {
	if u == nil {
		return false
	}
	if u.InputTokens > 0 {
		return true
	}
	if u.Cache != nil && (u.Cache.ReadInputTokens > 0 || u.Cache.CachedTokens > 0 || u.Cache.CreationInputTokens > 0) {
		return true
	}
	if u.Thinking != nil && u.Thinking.ThinkingTokens > 0 {
		return true
	}
	return false
}

// BillableInput 返回应计费的输入 token 数。
//
// 计费逻辑：
//   - cache_read 命中不计费（来自历史缓存，不消耗实时算力）
//   - cache_creation 写入按输入价计费（占用了实时算力建缓存）
//   - 公式：FreshInputTokens + CacheCreationInputTokens
//
// FreshInputTokens 已是"去 cache 后的输入"：
//   - OpenAI/Gemini：InputTokens - CachedTokens（缓存命中部分剔除）
//   - Anthropic：等于 InputTokens（Anthropic 协议 input_tokens 本就不含 cache）
//
// Ephemeral5m / Ephemeral1h 是 CreationInputTokens 的 4.6+ TTL 子集分桶，已包含在
// CreationInputTokens 顶层扁平字段中，不参与 BillableInput 公式避免重复计入。
//
// 注意：此为近似逻辑，实际计费由上游模型定价表决定。
func (u *Usage) BillableInput() int64 {
	if u == nil {
		return 0
	}
	billable := u.FreshInputTokens
	if u.Cache != nil {
		billable += u.Cache.CreationInputTokens
	}
	return billable
}

// Clone 返回深拷贝，避免转换函数间共享指针状态。
func (u *Usage) Clone() *Usage {
	if u == nil {
		return nil
	}
	c := &Usage{
		Base:             u.Base,
		ImageCount:       u.ImageCount,
		AudioInputChars:  u.AudioInputChars,
		AudioOutputChars: u.AudioOutputChars,
		VideoDurationSec: u.VideoDurationSec,
		ServiceTier:      u.ServiceTier,
		InferenceGeo:     u.InferenceGeo,
		Protocol:         u.Protocol,
	}
	if u.Cache != nil {
		cc := *u.Cache
		c.Cache = &cc
	}
	if u.Reasoning != nil {
		rr := *u.Reasoning
		c.Reasoning = &rr
	}
	if u.Thinking != nil {
		tt := *u.Thinking
		c.Thinking = &tt
	}
	if u.ServerTools != nil {
		ss := *u.ServerTools
		c.ServerTools = &ss
	}
	if u.Audio != nil {
		aa := *u.Audio
		c.Audio = &aa
	}
	if u.Prediction != nil {
		pp := *u.Prediction
		c.Prediction = &pp
	}
	return c
}
