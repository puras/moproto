package anthropic

// AnthropicMcpServerType mcp_servers[].type 字段类型。
//
// 当前唯一取值为 "url"（URL 形式的 MCP 服务器）。
type AnthropicMcpServerType string

const AnthropicMcpServerTypeURL AnthropicMcpServerType = "url"

// AnthropicMcpServerDefinition MCP 连接器定义（beta）。
//
// 把远程 MCP 服务器挂载到 Claude 工具集。type 固定 "url"。
type AnthropicMcpServerDefinition struct {
	Type               AnthropicMcpServerType  `json:"type"` // 固定 "url"
	URL                string                  `json:"url"`
	Name               string                  `json:"name"`
	ToolConfiguration  *AnthropicMcpToolConfig `json:"tool_configuration,omitempty"`
	AuthorizationToken string                  `json:"authorization_token,omitempty"`
}

// AnthropicMcpToolConfig MCP 工具可见性配置。
type AnthropicMcpToolConfig struct {
	Enabled    bool     `json:"enabled,omitempty"`
	Allowed    []string `json:"allowed_tools,omitempty"`
	Disallowed []string `json:"disallowed_tools,omitempty"`
}
