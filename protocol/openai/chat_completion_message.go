package openai

import "encoding/json"

// OpenAIChatCompletionMessageParamUnion 消息 union 接口。
//
// 6 种实现：User / Assistant / System / Developer / Tool / Function。
type OpenAIChatCompletionMessageParamUnion interface {
	isChatCompletionMessageParam()
}

// OpenAIChatCompletionMessageParamUnionFlat 消息 union 扁平化实现。
type OpenAIChatCompletionMessageParamUnionFlat struct {
	// Type 标识当前激活的变体：user/system/developer/assistant/tool/function。
	Type string `json:"type,omitempty"`
	// User 用户消息变体。
	User *OpenAIChatCompletionUserMessageParam `json:"user,omitempty"`
	// Assistant 助手消息变体。
	Assistant *OpenAIChatCompletionAssistantMessageParam `json:"assistant,omitempty"`
	// System 系统消息变体。
	System *OpenAIChatCompletionSystemMessageParam `json:"system,omitempty"`
	// Developer 开发者消息变体。
	Developer *OpenAIChatCompletionDeveloperMessageParam `json:"developer,omitempty"`
	// Tool 工具结果消息变体。
	Tool *OpenAIChatCompletionToolMessageParam `json:"tool,omitempty"`
	// Function 旧版函数消息变体。
	Function *OpenAIChatCompletionFunctionMessageParam `json:"function,omitempty"`
}

func (OpenAIChatCompletionMessageParamUnionFlat) isChatCompletionMessageParam() {}

func (u OpenAIChatCompletionMessageParamUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "user":
		if u.User == nil {
			return nil, nil
		}
		return json.Marshal(u.User)
	case "assistant":
		if u.Assistant == nil {
			return nil, nil
		}
		return json.Marshal(u.Assistant)
	case "system":
		if u.System == nil {
			return nil, nil
		}
		return json.Marshal(u.System)
	case "developer":
		if u.Developer == nil {
			return nil, nil
		}
		return json.Marshal(u.Developer)
	case "tool":
		if u.Tool == nil {
			return nil, nil
		}
		return json.Marshal(u.Tool)
	case "function":
		if u.Function == nil {
			return nil, nil
		}
		return json.Marshal(u.Function)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionMessageParamUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "user":
		var v OpenAIChatCompletionUserMessageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.User = &v
	case "assistant":
		var v OpenAIChatCompletionAssistantMessageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Assistant = &v
	case "system":
		var v OpenAIChatCompletionSystemMessageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.System = &v
	case "developer":
		var v OpenAIChatCompletionDeveloperMessageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Developer = &v
	case "tool":
		var v OpenAIChatCompletionToolMessageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Tool = &v
	case "function":
		var v OpenAIChatCompletionFunctionMessageParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Function = &v
	}
	return nil
}

// OpenAIChatCompletionMessageRole 消息角色。
type OpenAIChatCompletionMessageRole string

const (
	OpenAIChatCompletionMessageRoleUser      OpenAIChatCompletionMessageRole = "user"
	OpenAIChatCompletionMessageRoleAssistant OpenAIChatCompletionMessageRole = "assistant"
	OpenAIChatCompletionMessageRoleSystem    OpenAIChatCompletionMessageRole = "system"
	OpenAIChatCompletionMessageRoleTool      OpenAIChatCompletionMessageRole = "tool"
	OpenAIChatCompletionMessageRoleDeveloper OpenAIChatCompletionMessageRole = "developer"
	OpenAIChatCompletionMessageRoleFunction  OpenAIChatCompletionMessageRole = "function"
)

// OpenAIChatCompletionUserMessageParam 用户消息。
type OpenAIChatCompletionUserMessageParam struct {
	// Role 固定 "user"。
	Role OpenAIChatCompletionMessageRole `json:"role"`
	// Content 消息内容（ContentPart union 或 []ContentPart）。
	Content OpenAIChatCompletionUserMessageContentUnion `json:"content"`
	// Name 用户名（可选）。
	Name string `json:"name,omitempty"`
}

func (OpenAIChatCompletionUserMessageParam) isChatCompletionMessageParam() {}

// OpenAIChatCompletionUserMessageContentUnion 用户消息内容 union。
type OpenAIChatCompletionUserMessageContentUnion interface {
	isChatCompletionUserMessageContent()
}

// OpenAIChatCompletionUserMessageContentUnionFlat 用户消息内容 union 扁平化实现。
type OpenAIChatCompletionUserMessageContentUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"string"/"parts"。
	Type string `json:"type,omitempty"`
	// String 纯文本变体。
	String *OpenAIChatCompletionUserMessageContentString `json:"string,omitempty"`
	// Parts 多模态内容块数组变体。
	Parts *OpenAIChatCompletionUserMessageContentParts `json:"parts,omitempty"`
}

func (OpenAIChatCompletionUserMessageContentUnionFlat) isChatCompletionUserMessageContent() {}

func (u OpenAIChatCompletionUserMessageContentUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "parts":
		if u.Parts == nil {
			return nil, nil
		}
		return json.Marshal(u.Parts)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionUserMessageContentUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIChatCompletionUserMessageContentString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "parts":
		var v OpenAIChatCompletionUserMessageContentParts
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Parts = &v
	}
	return nil
}

// OpenAIChatCompletionUserMessageContentString 纯文本用户消息。
type OpenAIChatCompletionUserMessageContentString string

func (OpenAIChatCompletionUserMessageContentString) isChatCompletionUserMessageContent() {}

// OpenAIChatCompletionUserMessageContentParts 多模态内容块数组。
type OpenAIChatCompletionUserMessageContentParts []OpenAIChatCompletionContentPartUnion

func (OpenAIChatCompletionUserMessageContentParts) isChatCompletionUserMessageContent() {}

// OpenAIChatCompletionAssistantMessageParam 助手消息。
type OpenAIChatCompletionAssistantMessageParam struct {
	// Role 固定 "assistant"。
	Role OpenAIChatCompletionMessageRole `json:"role"`
	// Content 文本内容（可选，工具调用时可为空）。
	Content *string `json:"content,omitempty"`
	// Name 助手名。
	Name string `json:"name,omitempty"`
	// ToolCalls 工具调用列表。
	ToolCalls []OpenAIChatCompletionMessageToolCallUnion `json:"tool_calls,omitempty"`
	// Audio 助手消息音频输出（4.6+）。
	Audio *OpenAIChatCompletionMessageAudioParam `json:"audio,omitempty"`
	// Refusal 拒绝消息。
	Refusal *string `json:"refusal,omitempty"`
	// FunctionCall 旧版函数调用（已废弃但保留兼容）。
	FunctionCall *OpenAIChatCompletionMessageFunctionCallParam `json:"function_call,omitempty"`
	// ReasoningContent 推理内容（DeepSeek / Kimi 等厂商扩展；多轮 thinking 模式必须带回）。
	//ReasoningContent *string `json:"reasoning_content,omitempty"`
}

func (OpenAIChatCompletionAssistantMessageParam) isChatCompletionMessageParam() {}

// OpenAIChatCompletionMessageAudioParam 助手消息音频。
type OpenAIChatCompletionMessageAudioParam struct {
	// ID 音频引用 ID。
	ID string `json:"id,omitempty"`
	// Data base64 编码的音频数据。
	Data string `json:"data,omitempty"`
	// ExpiresAt 音频过期时间。
	ExpiresAt int64 `json:"expires_at,omitempty"`
	// Transcript 音频转录文本。
	Transcript string `json:"transcript,omitempty"`
}

// OpenAIChatCompletionMessageFunctionCallParam 旧版函数调用。
type OpenAIChatCompletionMessageFunctionCallParam struct {
	// Arguments 参数 JSON 字符串。
	Arguments string `json:"arguments,omitempty"`
	// Name 函数名。
	Name string `json:"name,omitempty"`
}

// OpenAIChatCompletionSystemMessageParam 系统消息。
type OpenAIChatCompletionSystemMessageParam struct {
	// Role 固定 "system"。
	Role OpenAIChatCompletionMessageRole `json:"role"`
	// Content 文本或内容块数组。
	Content OpenAIChatCompletionSystemMessageContentUnion `json:"content"`
	// Name 系统名（可选）。
	Name string `json:"name,omitempty"`
}

func (OpenAIChatCompletionSystemMessageParam) isChatCompletionMessageParam() {}

// OpenAIChatCompletionSystemMessageContentUnion 系统消息内容 union。
type OpenAIChatCompletionSystemMessageContentUnion interface {
	isChatCompletionSystemMessageContent()
}

// OpenAIChatCompletionSystemMessageContentUnionFlat 系统消息内容 union 扁平化实现。
type OpenAIChatCompletionSystemMessageContentUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"string"/"parts"。
	Type string `json:"type,omitempty"`
	// String 纯文本变体。
	String *OpenAIChatCompletionSystemMessageContentString `json:"string,omitempty"`
	// Parts 多模态内容块数组变体。
	Parts *OpenAIChatCompletionSystemMessageContentParts `json:"parts,omitempty"`
}

func (OpenAIChatCompletionSystemMessageContentUnionFlat) isChatCompletionSystemMessageContent() {}

func (u OpenAIChatCompletionSystemMessageContentUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "parts":
		if u.Parts == nil {
			return nil, nil
		}
		return json.Marshal(u.Parts)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionSystemMessageContentUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIChatCompletionSystemMessageContentString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "parts":
		var v OpenAIChatCompletionSystemMessageContentParts
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Parts = &v
	}
	return nil
}

// OpenAIChatCompletionSystemMessageContentString 纯文本系统消息。
type OpenAIChatCompletionSystemMessageContentString string

func (OpenAIChatCompletionSystemMessageContentString) isChatCompletionSystemMessageContent() {}

// OpenAIChatCompletionSystemMessageContentParts 多模态系统消息内容块数组。
type OpenAIChatCompletionSystemMessageContentParts []OpenAIChatCompletionContentPartUnion

func (OpenAIChatCompletionSystemMessageContentParts) isChatCompletionSystemMessageContent() {}

// OpenAIChatCompletionDeveloperMessageParam 开发者消息。
type OpenAIChatCompletionDeveloperMessageParam struct {
	// Role 固定 "developer"。
	Role OpenAIChatCompletionMessageRole `json:"role"`
	// Content 文本或内容块数组。
	Content OpenAIChatCompletionDeveloperMessageContentUnion `json:"content"`
	// Name 开发者名（可选）。
	Name string `json:"name,omitempty"`
}

func (OpenAIChatCompletionDeveloperMessageParam) isChatCompletionMessageParam() {}

// OpenAIChatCompletionDeveloperMessageContentUnion 开发者消息内容 union。
type OpenAIChatCompletionDeveloperMessageContentUnion interface {
	isChatCompletionDeveloperMessageContent()
}

// OpenAIChatCompletionDeveloperMessageContentUnionFlat 开发者消息内容 union 扁平化实现。
type OpenAIChatCompletionDeveloperMessageContentUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"string"/"parts"。
	Type string `json:"type,omitempty"`
	// String 纯文本变体。
	String *OpenAIChatCompletionDeveloperMessageContentString `json:"string,omitempty"`
	// Parts 多模态内容块数组变体。
	Parts *OpenAIChatCompletionDeveloperMessageContentParts `json:"parts,omitempty"`
}

func (OpenAIChatCompletionDeveloperMessageContentUnionFlat) isChatCompletionDeveloperMessageContent() {
}

func (u OpenAIChatCompletionDeveloperMessageContentUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "parts":
		if u.Parts == nil {
			return nil, nil
		}
		return json.Marshal(u.Parts)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionDeveloperMessageContentUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIChatCompletionDeveloperMessageContentString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "parts":
		var v OpenAIChatCompletionDeveloperMessageContentParts
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Parts = &v
	}
	return nil
}

// OpenAIChatCompletionDeveloperMessageContentString 纯文本开发者消息。
type OpenAIChatCompletionDeveloperMessageContentString string

func (OpenAIChatCompletionDeveloperMessageContentString) isChatCompletionDeveloperMessageContent() {}

// OpenAIChatCompletionDeveloperMessageContentParts 多模态开发者消息内容块数组。
type OpenAIChatCompletionDeveloperMessageContentParts []OpenAIChatCompletionContentPartUnion

func (OpenAIChatCompletionDeveloperMessageContentParts) isChatCompletionDeveloperMessageContent() {}

// OpenAIChatCompletionToolMessageParam 工具结果消息。
type OpenAIChatCompletionToolMessageParam struct {
	// Role 固定 "tool"。
	Role OpenAIChatCompletionMessageRole `json:"role"`
	// Content 工具返回内容。
	Content OpenAIChatCompletionToolMessageContentUnion `json:"content"`
	// ToolCallID 工具调用 ID。
	ToolCallID string `json:"tool_call_id"`
}

func (OpenAIChatCompletionToolMessageParam) isChatCompletionMessageParam() {}

// OpenAIChatCompletionToolMessageContentUnion 工具消息内容 union。
type OpenAIChatCompletionToolMessageContentUnion interface {
	isChatCompletionToolMessageContent()
}

// OpenAIChatCompletionToolMessageContentUnionFlat 工具消息内容 union 扁平化实现。
type OpenAIChatCompletionToolMessageContentUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"string"/"parts"。
	Type string `json:"type,omitempty"`
	// String 纯文本变体。
	String *OpenAIChatCompletionToolMessageContentString `json:"string,omitempty"`
	// Parts 多模态内容块数组变体。
	Parts *OpenAIChatCompletionToolMessageContentParts `json:"parts,omitempty"`
}

func (OpenAIChatCompletionToolMessageContentUnionFlat) isChatCompletionToolMessageContent() {}

func (u OpenAIChatCompletionToolMessageContentUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "parts":
		if u.Parts == nil {
			return nil, nil
		}
		return json.Marshal(u.Parts)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionToolMessageContentUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIChatCompletionToolMessageContentString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "parts":
		var v OpenAIChatCompletionToolMessageContentParts
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Parts = &v
	}
	return nil
}

// OpenAIChatCompletionToolMessageContentString 纯文本工具消息。
type OpenAIChatCompletionToolMessageContentString string

func (OpenAIChatCompletionToolMessageContentString) isChatCompletionToolMessageContent() {}

// OpenAIChatCompletionToolMessageContentParts 多模态工具消息内容块数组。
type OpenAIChatCompletionToolMessageContentParts []OpenAIChatCompletionContentPartUnion

func (OpenAIChatCompletionToolMessageContentParts) isChatCompletionToolMessageContent() {}

// OpenAIChatCompletionFunctionMessageParam 旧版函数调用结果消息。
type OpenAIChatCompletionFunctionMessageParam struct {
	// Role 固定 "function"。
	Role OpenAIChatCompletionMessageRole `json:"role"`
	// Content 函数返回文本。
	Content string `json:"content"`
	// Name 函数名。
	Name string `json:"name"`
}

func (OpenAIChatCompletionFunctionMessageParam) isChatCompletionMessageParam() {}
