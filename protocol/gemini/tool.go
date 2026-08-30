package gemini

// GeminiFunctionCall 模型对工具的调用请求（响应侧 functionCall 字段）。
//
// 与 OpenAI 的 tool_call 类似，但语义不同：
//   - ID 可选（Gemini 不强制 ID；OpenAI 必有 tool_call_id）；
//   - Args 为 map[string]any（官方 SDK 一致）。
type GeminiFunctionCall struct {
	ID   string         `json:"id,omitempty"`
	Name string         `json:"name,omitempty"`
	Args map[string]any `json:"args,omitempty"`
}

// GeminiFunctionResponse 工具执行结果回传（请求侧 functionResponse 字段）。
//
// 触发新一轮模型调用；Name 必须匹配先前 FunctionCall 的 Name。
type GeminiFunctionResponse struct {
	ID       string         `json:"id,omitempty"`
	Name     string         `json:"name,omitempty"`
	Response map[string]any `json:"response,omitempty"`
}

// GeminiFunctionDeclaration 单个工具的 schema 定义。
//
// 与 OpenAI function description 等价；Parameters 是 OpenAPI 3.0 子集 schema。
type GeminiFunctionDeclaration struct {
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Parameters  *GeminiSchema `json:"parameters,omitempty"`
}

// GeminiTool 工具容器（Gemini 协议支持多种工具类型，P0 仅函数声明）。
type GeminiTool struct {
	FunctionDeclarations []*GeminiFunctionDeclaration `json:"functionDeclarations,omitempty"`
}

// GeminiToolConfig 工具调用配置。
type GeminiToolConfig struct {
	FunctionCallingConfig *GeminiFunctionCallingConfig `json:"functionCallingConfig,omitempty"`
}

// GeminiFunctionCallingConfig 函数调用行为配置。
type GeminiFunctionCallingConfig struct {
	Mode                 GeminiFunctionCallingMode `json:"mode,omitempty"`
	AllowedFunctionNames []string                  `json:"allowedFunctionNames,omitempty"`
}
