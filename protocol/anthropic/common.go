package anthropic

// AnthropicCacheControlEphemeralType cache_control.type 字段类型。
//
// 固定取值为 "ephemeral"；类型化为独立 string alias，避免与协议中其他 type
// 标签（如 content block 的 "text"）混用。
type AnthropicCacheControlEphemeralType string

// AnthropicCacheControlEphemeralTypeEphemeral 唯一的合法取值。
const AnthropicCacheControlEphemeralTypeEphemeral AnthropicCacheControlEphemeralType = "ephemeral"

// AnthropicCacheControlEphemeral 提示缓存控制（Anthropic prompt caching）。
//
// 4.6+ 支持 ttl 取值："5m"（默认，5 分钟）或 "1h"（1 小时）。
// 挂载在文本块 / 工具 / 顶层 request 的最后一个可缓存块上。
type AnthropicCacheControlEphemeral struct {
	Type AnthropicCacheControlEphemeralType `json:"type"`          // 固定 "ephemeral"
	TTL  string                             `json:"ttl,omitempty"` // "5m" | "1h"
}

// AnthropicCitationsConfig 文本块引文行为配置。
type AnthropicCitationsConfig struct {
	Enabled bool `json:"enabled"`
}

// AnthropicMetadata 请求级元数据（透传给计费/审计链路）。
type AnthropicMetadata struct {
	UserID string `json:"user_id,omitempty"`
}

// AnthropicMessageType message.type 字段类型。
//
// 顶层 Message 响应的 type 固定为 "message"。
type AnthropicMessageType string

// AnthropicMessageTypeMessage 唯一合法取值。
const AnthropicMessageTypeMessage AnthropicMessageType = "message"

// AnthropicMessageRole 消息角色。
//
// 响应中 assistant 角色固定；请求中可取 user / assistant / system。
type AnthropicMessageRole string

const (
	AnthropicMessageRoleUser      AnthropicMessageRole = "user"
	AnthropicMessageRoleAssistant AnthropicMessageRole = "assistant"
	AnthropicMessageRoleSystem    AnthropicMessageRole = "system"
)

// AnthropicToolCallerDirectType caller.type = "direct"。
type AnthropicToolCallerDirectType string

// AnthropicToolCallerDirectTypeDirect 唯一合法取值。
const AnthropicToolCallerDirectTypeDirect AnthropicToolCallerDirectType = "direct"

// AnthropicToolCallerCodeExecution20250825Type caller.type = "code_execution_20250825"。
type AnthropicToolCallerCodeExecution20250825Type string

// AnthropicToolCallerCodeExecution20250825TypeCodeExecution20250825 唯一合法取值。
const AnthropicToolCallerCodeExecution20250825TypeCodeExecution20250825 AnthropicToolCallerCodeExecution20250825Type = "code_execution_20250825"

// AnthropicToolCallerCodeExecution20260120Type caller.type = "code_execution_20260120"。
type AnthropicToolCallerCodeExecution20260120Type string

// AnthropicToolCallerCodeExecution20260120TypeCodeExecution20260120 唯一合法取值。
const AnthropicToolCallerCodeExecution20260120TypeCodeExecution20260120 AnthropicToolCallerCodeExecution20260120Type = "code_execution_20260120"
