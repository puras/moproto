package anthropic

import "encoding/json"

// AnthropicRequest Anthropic Messages API 请求（wire-format）。
//
// 顶层字段基于官方 API 参考 + Go SDK MessageNewParams，覆盖 Claude 4.6/4.7/4.8/Fable 5
// 全部请求字段。强类型 union（Tool / ToolChoice / ContentBlock）替代 common 包里的
// 单 struct + any 模式。
type AnthropicRequest struct {
	// ---- 必填 ----
	Model     string                  `json:"model"`
	MaxTokens int                     `json:"max_tokens"`
	Messages  []AnthropicMessageParam `json:"messages"`

	// ---- 系统提示 ----
	// 服务端同时接受 string 和 []AnthropicTextBlockParam 两种形态，用 RawMessage 保留。
	System json.RawMessage `json:"system,omitempty"`

	// ---- 工具 ----
	Tools      []AnthropicTool     `json:"tools,omitempty"`
	ToolChoice AnthropicToolChoice `json:"tool_choice,omitempty"`

	// ---- 流式 ----
	Stream bool `json:"stream,omitempty"`

	// ---- 思考 / 输出控制 ----
	// Thinking 是 3 mode union：adaptive（4.6+ 推荐）/ enabled（传统 budget_tokens）/ disabled
	Thinking     AnthropicThinkingConfig `json:"thinking,omitempty"`
	OutputConfig *AnthropicOutputConfig  `json:"output_config,omitempty"`

	// ---- 采样参数（Opus 4.7+/Fable 5 已禁用） ----
	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
	TopK        *int     `json:"top_k,omitempty"`

	// ---- 停止 / 元数据 ----
	StopSequences []string           `json:"stop_sequences,omitempty"`
	Metadata      *AnthropicMetadata `json:"metadata,omitempty"`

	// ---- 顶层缓存（标记最后一个可缓存块） ----
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`

	// ---- 长上下文管理（beta） ----
	ContextManagement *AnthropicContextManagement `json:"context_management,omitempty"`

	// ---- 容器 / 技能（string ID 或 object 形态） ----
	// 用 RawMessage：string 用于续接，object 含 skills 列表。
	Container json.RawMessage `json:"container,omitempty"`

	// ---- MCP 连接器（beta） ----
	McpServers []AnthropicMcpServerDefinition `json:"mcp_servers,omitempty"`

	// ---- 地理 / 服务等级 / 用户档案 ----
	InferenceGeo  string      `json:"inference_geo,omitempty"`   // 数据驻留地，如 "us"（语义自由，保持 string）
	ServiceTier   ServiceTier `json:"service_tier,omitempty"`    // "auto" | "standard_only"
	UserProfileID string      `json:"user_profile_id,omitempty"` // user-profiles beta

	// ---- Speed（Opus 4.8 fast 模式） ----
	Speed Speed `json:"speed,omitempty"` // "fast"
}

// GetModel 返回模型名（与 common.AnthropicRequest 接口对齐，便于上层接入）。
func (r *AnthropicRequest) GetModel() string { return r.Model }

// GetFormat 返回协议格式标识 "anthropic"。
//
// 用于满足 internal/app/gateway/common.Request 接口。
// 注：返回 string 而非 common.Format 是为了避免 pkg → internal 反向引用；
// common.Format 已被定义为 string alias（见 common/format.go），所以方法签名匹配。
func (*AnthropicRequest) GetFormat() string { return "anthropic" }

// ServiceTier 取值。
type ServiceTier string

const (
	ServiceTierAuto         ServiceTier = "auto"
	ServiceTierStandardOnly ServiceTier = "standard_only"
)

// Speed 取值。
type Speed string

const SpeedFast Speed = "fast"
