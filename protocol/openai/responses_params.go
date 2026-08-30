package openai

import "encoding/json"

// OpenAIResponseInputUnionFlat 输入 union 扁平化实现（2 变体）。
type OpenAIResponseInputUnionFlat struct {
	Type  string                             `json:"type,omitempty"` // "string" / "items"
	Text  string                             `json:"text,omitempty"`
	Items []OpenAIResponseInputItemUnionFlat `json:"items,omitempty"`
}

func (OpenAIResponseInputUnionFlat) isResponseInputUnion() {}

func (u OpenAIResponseInputUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.Text == "" {
			return nil, nil
		}
		// Return as raw JSON string value (with quotes)
		return []byte("\"" + u.Text + "\""), nil
	case "items":
		if u.Items == nil {
			return nil, nil
		}
		return json.Marshal(u.Items)
	}
	return nil, nil
}

func (u *OpenAIResponseInputUnionFlat) UnmarshalJSON(data []byte) error {
	// Check first char: '"' means raw string, '[' means items array
	if len(data) > 0 && data[0] == '"' {
		// Raw string: "hello"
		u.Type = "string"
		return json.Unmarshal(data, &u.Text)
	}
	if len(data) > 0 && data[0] == '[' {
		// Items array
		u.Type = "items"
		var items []OpenAIResponseInputItemUnionFlat
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		u.Items = items
		return nil
	}
	// Object with explicit type field
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	if tmp.Type == "items" {
		var items []OpenAIResponseInputItemUnionFlat
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		u.Items = items
	}
	return nil
}

// OpenAIResponseParams 是 Responses API 的请求顶层类型（替代官方 SDK
// ResponseNewParams 的 New 中缀）。
//
// 字段集与 OpenAI 官方 Go SDK ResponseNewParams 严格对齐。
// 详见 https://platform.openai.com/docs/api-reference/responses/create。
//
// 实现 common.Request 接口（GetModel / GetFormat）。
func (r *OpenAIResponseParams) GetModel() string  { return r.Model }
func (r *OpenAIResponseParams) GetFormat() string { return "responses" }

type OpenAIResponseParams struct {
	// === 必填字段 ===

	// Model 使用的模型 ID（如 "gpt-4o"）。
	Model string `json:"model"`
	// Input 输入内容（string 或 []InputItem）。
	Input *OpenAIResponseInputUnionFlat `json:"input"`

	// === 可选字段 ===

	// Instructions 系统指令（string 或 []InputItem）。
	Instructions OpenAIResponseInstructionsUnion `json:"instructions,omitempty"`

	// === 工具 ===

	// Tools 工具列表。
	Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
	// ToolChoice 工具选择策略。
	ToolChoice OpenAIResponseToolChoiceUnion `json:"tool_choice,omitempty"`
	// ParallelToolCalls 是否允许并行工具调用。
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`

	// === 流式 ===

	// Stream 是否流式。
	Stream *bool `json:"stream,omitempty"`

	// === 输出控制 ===

	// Include 需要包含的字段列表（如 "file_search_call.results"）。
	Include []string `json:"include,omitempty"`
	// Store 是否存储用于评估。
	Store *bool `json:"store,omitempty"`
	// Truncation 截断策略。
	Truncation OpenAIResponseTruncation `json:"truncation,omitempty"`
	// MaxOutputTokens 最大输出 token 数。
	MaxOutputTokens *int64 `json:"max_output_tokens,omitempty"`

	// === 推理（o 系列）===

	// Reasoning 推理配置。
	Reasoning *OpenAIResponseReasoningParam `json:"reasoning,omitempty"`

	// === 元数据 ===

	// Metadata 自定义元数据。
	Metadata map[string]string `json:"metadata,omitempty"`

	// === 文本输出配置 ===

	// Text 文本输出配置。
	Text *OpenAIResponseTextConfig `json:"text,omitempty"`

	// === 后台执行 ===

	// Background 是否后台运行。
	Background *bool `json:"background,omitempty"`

	// === 安全标识 ===

	// SafetyIdentifier 安全标识（4.6+）。
	SafetyIdentifier string `json:"safety_identifier,omitempty"`

	// === 提示缓存 ===

	// PromptCacheKey 提示缓存键。
	PromptCacheKey string `json:"prompt_cache_key,omitempty"`

	// === 多轮会话 ===

	// Conversation 会话引用。
	Conversation *OpenAIResponseConversationParam `json:"conversation,omitempty"`

	// === P1 补齐字段（2026-07-01）===

	// MaxToolCalls 最大工具调用次数（4.6+）。
	MaxToolCalls *int64 `json:"max_tool_calls,omitempty"`
	// PreviousResponseID 上一轮响应 ID（多轮会话）。
	PreviousResponseID string `json:"previous_response_id,omitempty"`
	// Prompt 提示模板引用。
	Prompt *OpenAIResponsePrompt `json:"prompt,omitempty"`
	// Moderation 审核配置。
	Moderation *OpenAIResponseModeration `json:"moderation,omitempty"`
	// StreamOptions 流式选项。
	StreamOptions *OpenAIResponseStreamOptions `json:"stream_options,omitempty"`
	// TopLogprobs 返回每个位置的最可能 token 数（0-20）。
	TopLogprobs *int64 `json:"top_logprobs,omitempty"`
	// ContextManagement 上下文管理配置。
	ContextManagement []OpenAIResponseContextManagement `json:"context_management,omitempty"`
}

// === 辅助 union 类型 ===

// OpenAIResponseInputUnion 输入 union（string | []InputItem）。
type OpenAIResponseInputUnion interface {
	isResponseInputUnion()
}

// OpenAIResponseInputUnionString 纯文本输入。
type OpenAIResponseInputUnionString string

func (OpenAIResponseInputUnionString) isResponseInputUnion() {}

// OpenAIResponseInputUnionItems 输入项列表。
type OpenAIResponseInputUnionItems []OpenAIResponseInputItemUnion

func (OpenAIResponseInputUnionItems) isResponseInputUnion() {}

// OpenAIResponseInstructionsUnion instructions union（string | []InputItem）。
type OpenAIResponseInstructionsUnion interface {
	isResponseInstructionsUnion()
}

// OpenAIResponseInstructionsUnionString 纯文本 instructions。
type OpenAIResponseInstructionsUnionString string

func (OpenAIResponseInstructionsUnionString) isResponseInstructionsUnion() {}

// OpenAIResponseInstructionsUnionItems instructions 输入项列表。
type OpenAIResponseInstructionsUnionItems []OpenAIResponseInputItemUnion

func (OpenAIResponseInstructionsUnionItems) isResponseInstructionsUnion() {}

// OpenAIResponseInstructionsUnionFlat instructions union 扁平化实现（2 变体）。
type OpenAIResponseInstructionsUnionFlat struct {
	Type  string                             `json:"type,omitempty"` // "string" / "items"
	Text  string                             `json:"text,omitempty"`
	Items []OpenAIResponseInputItemUnionFlat `json:"items,omitempty"`
}

func (OpenAIResponseInstructionsUnionFlat) isResponseInstructionsUnion() {}

func (u OpenAIResponseInstructionsUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.Text == "" {
			return nil, nil
		}
		return []byte("\"" + u.Text + "\""), nil
	case "items":
		if u.Items == nil {
			return nil, nil
		}
		return json.Marshal(u.Items)
	}
	return nil, nil
}

func (u *OpenAIResponseInstructionsUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		u.Type = "string"
		return json.Unmarshal(data, &u.Text)
	}
	if len(data) > 0 && data[0] == '[' {
		u.Type = "items"
		var items []OpenAIResponseInputItemUnionFlat
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		u.Items = items
		return nil
	}
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	if tmp.Type == "items" {
		var items []OpenAIResponseInputItemUnionFlat
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		u.Items = items
	}
	return nil
}

// OpenAIResponseToolChoiceUnion 工具选择策略 union（10 变体）。
type OpenAIResponseToolChoiceUnion interface {
	isResponseToolChoiceUnion()
}

// OpenAIResponseToolChoiceUnionFlat 工具选择策略 union 扁平化实现（10 变体）。
type OpenAIResponseToolChoiceUnionFlat struct {
	// Type 区分变体："auto" / "required" / "none" / "allowed_tools" / "apply_patch" / "custom" / "function" / "mcp" / "shell" / "types"
	Type string `json:"type,omitempty"`

	// auto / required / none（空对象）
	Extra json.RawMessage `json:"extra,omitempty"`

	// allowed_tools
	Mode  string                                           `json:"mode,omitempty"`
	Tools []OpenAIResponseToolChoiceOptionAllowedToolEntry `json:"tools,omitempty"`

	// apply_patch / shell（空对象，用 Extra）

	// custom
	Name string `json:"name,omitempty"`

	// function
	FunctionName string `json:"function_name,omitempty"`

	// mcp
	ServerLabel string `json:"server_label,omitempty"`
	McpName     string `json:"mcp_name,omitempty"`

	// types
	TypeList []string `json:"type_list,omitempty"`
}

func (OpenAIResponseToolChoiceUnionFlat) isResponseToolChoiceUnion() {}

func (u OpenAIResponseToolChoiceUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "auto", "required", "none", "apply_patch", "shell":
		if u.Extra == nil {
			return nil, nil
		}
		return u.Extra, nil
	case "allowed_tools":
		if u.Tools == nil {
			return nil, nil
		}
		return json.Marshal(struct {
			Type  string                                           `json:"type"`
			Mode  string                                           `json:"mode"`
			Tools []OpenAIResponseToolChoiceOptionAllowedToolEntry `json:"tools"`
		}{
			Type:  "allowed_tools",
			Mode:  u.Mode,
			Tools: u.Tools,
		})
	case "custom":
		return json.Marshal(struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}{
			Type: "custom",
			Name: u.Name,
		})
	case "function":
		return json.Marshal(struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}{
			Type: "function",
			Name: u.FunctionName,
		})
	case "mcp":
		return json.Marshal(struct {
			Type        string `json:"type"`
			ServerLabel string `json:"server_label"`
			Name        string `json:"name"`
		}{
			Type:        "mcp",
			ServerLabel: u.ServerLabel,
			Name:        u.McpName,
		})
	case "types":
		return json.Marshal(struct {
			Type []string `json:"type"`
		}{
			Type: u.TypeList,
		})
	}
	return nil, nil
}

func (u *OpenAIResponseToolChoiceUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '{' {
		var tmp struct {
			Type         string                                           `json:"type"`
			Mode         string                                           `json:"mode"`
			Tools        []OpenAIResponseToolChoiceOptionAllowedToolEntry `json:"tools"`
			Name         string                                           `json:"name"`
			FunctionName string                                           `json:"function_name"`
			ServerLabel  string                                           `json:"server_label"`
			McpName      string                                           `json:"name"`
			TypeList     []string                                         `json:"type"`
		}
		if err := json.Unmarshal(data, &tmp); err != nil {
			return err
		}
		u.Type = tmp.Type
		switch tmp.Type {
		case "auto", "required", "none", "apply_patch", "shell":
			u.Extra = json.RawMessage(data)
		case "allowed_tools":
			u.Mode = tmp.Mode
			u.Tools = tmp.Tools
		case "custom":
			u.Name = tmp.Name
		case "function":
			u.FunctionName = tmp.FunctionName
		case "mcp":
			u.ServerLabel = tmp.ServerLabel
			u.McpName = tmp.Name
		case "types":
			u.TypeList = tmp.TypeList
		}
		return nil
	}
	// 字符串形式：auto / required / none
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		u.Type = s
	}
	return nil
}

// OpenAIResponseReasoningParam 推理配置。
type OpenAIResponseReasoningParam struct {
	// Effort 推理强度。
	Effort OpenAIReasoningEffort `json:"effort,omitempty"`
	// Summary 摘要模式（auto / concise / detailed）。
	Summary string `json:"summary,omitempty"`
}

// OpenAIResponseTextConfig 文本输出配置。
type OpenAIResponseTextConfig struct {
	// Format 文本格式（text / json_schema）。
	Format *OpenAIResponseTextFormatUnionFlat `json:"format,omitempty"`
	// Verbosity 详细度。
	Verbosity OpenAIVerbosity `json:"verbosity,omitempty"`
}

// OpenAIResponseTextFormatUnion 文本格式 union（2 变体：text / json_schema）。
type OpenAIResponseTextFormatUnion interface {
	isResponseTextFormatUnion()
}

// OpenAIResponseTextFormatUnionFlat 文本格式 union 扁平化实现（2 变体）。
type OpenAIResponseTextFormatUnionFlat struct {
	// Type 区分变体："text" / "json_schema"
	Type string `json:"type,omitempty"`
	// JSONSchema JSON Schema 定义（当 Type = "json_schema" 时）。
	JSONSchema *OpenAIResponseTextFormatJSONSchema `json:"json_schema,omitempty"`
}

func (OpenAIResponseTextFormatUnionFlat) isResponseTextFormatUnion() {}

func (u OpenAIResponseTextFormatUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "text":
		return []byte(`"text"`), nil
	case "json_schema":
		if u.JSONSchema == nil {
			return nil, nil
		}
		return json.Marshal(u.JSONSchema)
	}
	return nil, nil
}

func (u *OpenAIResponseTextFormatUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		u.Type = s
		return nil
	}
	var tmp struct {
		Type       string                              `json:"type"`
		JSONSchema *OpenAIResponseTextFormatJSONSchema `json:"json_schema,omitempty"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	u.JSONSchema = tmp.JSONSchema
	return nil
}

// OpenAIResponseTextFormatJSONSchema JSON Schema 定义。
type OpenAIResponseTextFormatJSONSchema struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      any    `json:"schema,omitempty"`
	Strict      *bool  `json:"strict,omitempty"`
}

// OpenAIResponseConversationParam 会话引用。
type OpenAIResponseConversationParam struct {
	// ID 会话 ID。
	ID string `json:"id"`
}

// === P1 补齐字段（2026-07-01）===

// OpenAIResponsePrompt 引用提示模板。
type OpenAIResponsePrompt struct {
	// ID 提示模板 ID。
	ID string `json:"id"`
	// Variables 变量映射。
	Variables map[string]*OpenAIResponsePromptVariableUnionFlat `json:"variables,omitempty"`
	// Version 模板版本。
	Version string `json:"version,omitempty"`
}

// OpenAIResponsePromptVariableUnion 提示模板变量 union（4 变体）。
type OpenAIResponsePromptVariableUnion interface {
	isResponsePromptVariable()
}

type openAIResponsePromptVariableString string

func (openAIResponsePromptVariableString) isResponsePromptVariable() {}

// OpenAIResponsePromptVariableInputText 提示模板文本变量。
type OpenAIResponsePromptVariableInputText struct {
	Type string `json:"type"` // "input_text"
	Text string `json:"text"`
}

func (OpenAIResponsePromptVariableInputText) isResponsePromptVariable() {}

// OpenAIResponsePromptVariableInputImage 提示模板图像变量。
type OpenAIResponsePromptVariableInputImage struct {
	Type     string `json:"type"` // "input_image"
	ImageURL string `json:"image_url,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

func (OpenAIResponsePromptVariableInputImage) isResponsePromptVariable() {}

// OpenAIResponsePromptVariableInputFile 提示模板文件变量。
type OpenAIResponsePromptVariableInputFile struct {
	Type   string `json:"type"` // "input_file"
	FileID string `json:"file_id"`
}

func (OpenAIResponsePromptVariableInputFile) isResponsePromptVariable() {}

// OpenAIResponsePromptVariableUnionFlat 提示模板变量 union 扁平化实现（4 变体）。
type OpenAIResponsePromptVariableUnionFlat struct {
	// Type 区分变体："string" / "input_text" / "input_image" / "input_file"
	Type string `json:"type,omitempty"`

	// string
	String string `json:"string,omitempty"`

	// input_text
	Text string `json:"text,omitempty"`

	// input_image
	ImageURL string `json:"image_url,omitempty"`
	Detail   string `json:"detail,omitempty"`

	// input_file
	FileID string `json:"file_id,omitempty"`
}

func (OpenAIResponsePromptVariableUnionFlat) isResponsePromptVariable() {}

func (u OpenAIResponsePromptVariableUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == "" {
			return nil, nil
		}
		return []byte("\"" + u.String + "\""), nil
	case "input_text":
		return json.Marshal(struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{
			Type: "input_text",
			Text: u.Text,
		})
	case "input_image":
		return json.Marshal(struct {
			Type     string `json:"type"`
			ImageURL string `json:"image_url,omitempty"`
			Detail   string `json:"detail,omitempty"`
		}{
			Type:     "input_image",
			ImageURL: u.ImageURL,
			Detail:   u.Detail,
		})
	case "input_file":
		return json.Marshal(struct {
			Type   string `json:"type"`
			FileID string `json:"file_id"`
		}{
			Type:   "input_file",
			FileID: u.FileID,
		})
	}
	return nil, nil
}

func (u *OpenAIResponsePromptVariableUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		u.Type = "string"
		return json.Unmarshal(data, &u.String)
	}
	var tmp struct {
		Type     string `json:"type"`
		Text     string `json:"text"`
		ImageURL string `json:"image_url"`
		Detail   string `json:"detail"`
		FileID   string `json:"file_id"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "input_text":
		u.Text = tmp.Text
	case "input_image":
		u.ImageURL = tmp.ImageURL
		u.Detail = tmp.Detail
	case "input_file":
		u.FileID = tmp.FileID
	}
	return nil
}

// OpenAIResponseModeration 审核配置。
type OpenAIResponseModeration struct {
	// Model 审核模型。
	Model string `json:"model"`
}

// OpenAIResponseStreamOptions 流式选项。
type OpenAIResponseStreamOptions struct {
	// IncludeObfuscation 是否包含混淆（缓解侧信道攻击）。
	IncludeObfuscation *bool `json:"include_obfuscation,omitempty"`
}

// OpenAIResponseContextManagement 上下文管理配置。
type OpenAIResponseContextManagement struct {
	// Type 类型（目前仅 "compaction"）。
	Type string `json:"type"`
	// CompactThreshold 压缩触发阈值（token 数）。
	CompactThreshold *int64 `json:"compact_threshold,omitempty"`
}
