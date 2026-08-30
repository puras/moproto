package anthropic

import (
	"encoding/json"
)

// AnthropicContentBlock 响应内容块 union。
//
// 12 个变体对应 Anthropic API 全部响应内容类型。
//
// 实现策略：AnthropicContentBlockUnion 扁平 struct + 自定义 JSON marshal/unmarshal。
// 旧 sealed interface 模式（AnthropicContentBlock + 12 个 isContentBlock 实现）保留用于
// type switch 分发；新增 AnthropicContentBlockUnion 作为唯一的 concrete 实现，
// 解决 json.Unmarshal 无法反序列化到 interface 字段的问题。
//
// 调用方通过 type switch 或 .Type 字段分发。
type AnthropicContentBlock interface {
	isContentBlock()
}

// AnthropicContentBlockUnion 响应内容块 union 的扁平实现。
//
// 所有 12 个变体的字段平铺为 json.RawMessage（延迟反序列化），
// Type 字段作为 discriminator，MarshalJSON/UnmarshalJSON 负责派发。
//
// 设计依据：官方 SDK anthropic-sdk-go 的 ContentBlockUnion 扁平 struct 模式，
// 区别于旧 sealed interface 方案。零外部依赖，仅 stdlib encoding/json。
type AnthropicContentBlockUnion struct {
	Type                              string          `json:"type"`
	Text                              json.RawMessage `json:"text,omitempty"`
	Thinking                          json.RawMessage `json:"thinking,omitempty"`
	RedactedThinking                  json.RawMessage `json:"redacted_thinking,omitempty"`
	ToolUse                           json.RawMessage `json:"tool_use,omitempty"`
	ServerToolUse                     json.RawMessage `json:"server_tool_use,omitempty"`
	WebSearchToolResult               json.RawMessage `json:"web_search_tool_result,omitempty"`
	WebFetchToolResult                json.RawMessage `json:"web_fetch_tool_result,omitempty"`
	CodeExecutionToolResult           json.RawMessage `json:"code_execution_tool_result,omitempty"`
	BashCodeExecutionToolResult       json.RawMessage `json:"bash_code_execution_tool_result,omitempty"`
	TextEditorCodeExecutionToolResult json.RawMessage `json:"text_editor_code_execution_tool_result,omitempty"`
	ToolSearchToolResult              json.RawMessage `json:"tool_search_tool_result,omitempty"`
	ContainerUpload                   json.RawMessage `json:"container_upload,omitempty"`
}

func (AnthropicContentBlockUnion) isContentBlock() {}

// MarshalJSON 实现 AnthropicContentBlockUnion 的自定义序列化。
//
// 每个 RawMessage 字段存的是完整变体 JSON（含 type 字段），直接返回。
// 空 Type 或 nil RawMessage 返回 nil 以触发 omitempty suppression。
func (u AnthropicContentBlockUnion) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case string(AnthropicTextBlockTypeText):
		if u.Text == nil {
			return nil, nil
		}
		return u.Text, nil
	case string(AnthropicThinkingBlockTypeThinking):
		if u.Thinking == nil {
			return nil, nil
		}
		return u.Thinking, nil
	case string(AnthropicRedactedThinkingBlockTypeRedactedThinking):
		if u.RedactedThinking == nil {
			return nil, nil
		}
		return u.RedactedThinking, nil
	case string(AnthropicToolUseBlockTypeToolUse):
		if u.ToolUse == nil {
			return nil, nil
		}
		return u.ToolUse, nil
	case string(AnthropicServerToolUseBlockTypeServerToolUse):
		if u.ServerToolUse == nil {
			return nil, nil
		}
		return u.ServerToolUse, nil
	case string(AnthropicWebSearchToolResultBlockTypeWebSearchToolResult):
		if u.WebSearchToolResult == nil {
			return nil, nil
		}
		return u.WebSearchToolResult, nil
	case string(AnthropicWebFetchToolResultBlockTypeWebFetchToolResult):
		if u.WebFetchToolResult == nil {
			return nil, nil
		}
		return u.WebFetchToolResult, nil
	case string(AnthropicCodeExecutionToolResultBlockTypeCodeExecutionToolResult):
		if u.CodeExecutionToolResult == nil {
			return nil, nil
		}
		return u.CodeExecutionToolResult, nil
	case string(AnthropicBashCodeExecutionToolResultBlockTypeBashCodeExecutionToolResult):
		if u.BashCodeExecutionToolResult == nil {
			return nil, nil
		}
		return u.BashCodeExecutionToolResult, nil
	case string(AnthropicTextEditorCodeExecutionToolResultBlockTypeTextEditorCodeExecutionToolResult):
		if u.TextEditorCodeExecutionToolResult == nil {
			return nil, nil
		}
		return u.TextEditorCodeExecutionToolResult, nil
	case string(AnthropicToolSearchToolResultBlockTypeToolSearchToolResult):
		if u.ToolSearchToolResult == nil {
			return nil, nil
		}
		return u.ToolSearchToolResult, nil
	case string(AnthropicContainerUploadBlockTypeContainerUpload):
		if u.ContainerUpload == nil {
			return nil, nil
		}
		return u.ContainerUpload, nil
	default:
		return nil, nil
	}
}

// UnmarshalJSON 实现 AnthropicContentBlockUnion 的自定义反序列化。
//
// strategy：直接 json.Unmarshal 到目标变体 struct（延迟到具体变体层级），
// 比手动抽取字段更安全。
func (u *AnthropicContentBlockUnion) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type

	switch tmp.Type {
	case string(AnthropicTextBlockTypeText):
		var b AnthropicTextBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.Text, _ = json.Marshal(b)
		}
	case string(AnthropicThinkingBlockTypeThinking):
		var b AnthropicThinkingBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.Thinking, _ = json.Marshal(b)
		}
	case string(AnthropicRedactedThinkingBlockTypeRedactedThinking):
		var b AnthropicRedactedThinkingBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.RedactedThinking, _ = json.Marshal(b)
		}
	case string(AnthropicToolUseBlockTypeToolUse):
		var b AnthropicToolUseBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.ToolUse, _ = json.Marshal(b)
		}
	case string(AnthropicServerToolUseBlockTypeServerToolUse):
		var b AnthropicServerToolUseBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.ServerToolUse, _ = json.Marshal(b)
		}
	case string(AnthropicWebSearchToolResultBlockTypeWebSearchToolResult):
		var b AnthropicWebSearchToolResultBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.WebSearchToolResult, _ = json.Marshal(b)
		}
	case string(AnthropicWebFetchToolResultBlockTypeWebFetchToolResult):
		var b AnthropicWebFetchToolResultBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.WebFetchToolResult, _ = json.Marshal(b)
		}
	case string(AnthropicCodeExecutionToolResultBlockTypeCodeExecutionToolResult):
		var b AnthropicCodeExecutionToolResultBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.CodeExecutionToolResult, _ = json.Marshal(b)
		}
	case string(AnthropicBashCodeExecutionToolResultBlockTypeBashCodeExecutionToolResult):
		var b AnthropicBashCodeExecutionToolResultBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.BashCodeExecutionToolResult, _ = json.Marshal(b)
		}
	case string(AnthropicTextEditorCodeExecutionToolResultBlockTypeTextEditorCodeExecutionToolResult):
		var b AnthropicTextEditorCodeExecutionToolResultBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.TextEditorCodeExecutionToolResult, _ = json.Marshal(b)
		}
	case string(AnthropicToolSearchToolResultBlockTypeToolSearchToolResult):
		var b AnthropicToolSearchToolResultBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.ToolSearchToolResult, _ = json.Marshal(b)
		}
	case string(AnthropicContainerUploadBlockTypeContainerUpload):
		var b AnthropicContainerUploadBlock
		if err := json.Unmarshal(data, &b); err == nil {
			u.ContainerUpload, _ = json.Marshal(b)
		}
	}
	return nil
}

// ============== 类型化 Type 字段（Response 侧新增） ==============
//
// 注意：以下 6 个 Type 别名是 Response 侧新引入，不在 Request 侧的
// message.go 中。已被 Request 侧复用的 6 个别名（Text / Thinking /
// RedactedThinking / ToolUse / ServerToolUse / WebSearchToolResult）
// 直接 import 复用，不再重复声明。

// AnthropicWebFetchToolResultBlockType web_fetch_tool_result 内容块的 type 字段类型。
type AnthropicWebFetchToolResultBlockType string

const AnthropicWebFetchToolResultBlockTypeWebFetchToolResult AnthropicWebFetchToolResultBlockType = "web_fetch_tool_result"

// AnthropicCodeExecutionToolResultBlockType code_execution_tool_result 内容块的 type 字段类型。
type AnthropicCodeExecutionToolResultBlockType string

const AnthropicCodeExecutionToolResultBlockTypeCodeExecutionToolResult AnthropicCodeExecutionToolResultBlockType = "code_execution_tool_result"

// AnthropicBashCodeExecutionToolResultBlockType bash_code_execution_tool_result 内容块的 type 字段类型。
type AnthropicBashCodeExecutionToolResultBlockType string

const AnthropicBashCodeExecutionToolResultBlockTypeBashCodeExecutionToolResult AnthropicBashCodeExecutionToolResultBlockType = "bash_code_execution_tool_result"

// AnthropicTextEditorCodeExecutionToolResultBlockType text_editor_code_execution_tool_result 内容块的 type 字段类型。
type AnthropicTextEditorCodeExecutionToolResultBlockType string

const AnthropicTextEditorCodeExecutionToolResultBlockTypeTextEditorCodeExecutionToolResult AnthropicTextEditorCodeExecutionToolResultBlockType = "text_editor_code_execution_tool_result"

// AnthropicToolSearchToolResultBlockType tool_search_tool_result 内容块的 type 字段类型。
type AnthropicToolSearchToolResultBlockType string

const AnthropicToolSearchToolResultBlockTypeToolSearchToolResult AnthropicToolSearchToolResultBlockType = "tool_search_tool_result"

// AnthropicContainerUploadBlockType container_upload 内容块的 type 字段类型。
type AnthropicContainerUploadBlockType string

const AnthropicContainerUploadBlockTypeContainerUpload AnthropicContainerUploadBlockType = "container_upload"

// ============== 12 个变体 struct ==============

// AnthropicTextBlock 文本块。
type AnthropicTextBlock struct {
	Type      AnthropicTextBlockType `json:"type"` // 固定 "text"
	Text      string                 `json:"text"`
	Citations []json.RawMessage      `json:"citations,omitempty"` // P1 补齐 5 种 TextCitation
}

func (AnthropicTextBlock) isContentBlock() {}

// AnthropicThinkingBlock thinking 块（4.6+ extended thinking 输出）。
type AnthropicThinkingBlock struct {
	Type      AnthropicThinkingBlockType `json:"type"` // 固定 "thinking"
	Thinking  string                     `json:"thinking"`
	Signature string                     `json:"signature"`
}

func (AnthropicThinkingBlock) isContentBlock() {}

// AnthropicRedactedThinkingBlock 加密的 thinking 块（用于安全审计场景，服务端加密后回传）。
type AnthropicRedactedThinkingBlock struct {
	Type AnthropicRedactedThinkingBlockType `json:"type"` // 固定 "redacted_thinking"
	Data string                             `json:"data"`
}

func (AnthropicRedactedThinkingBlock) isContentBlock() {}

// AnthropicToolUseBlock 工具调用块（模型主动调用客户端工具）。
type AnthropicToolUseBlock struct {
	Type  AnthropicToolUseBlockType `json:"type"` // 固定 "tool_use"
	ID    string                    `json:"id"`
	Name  string                    `json:"name"`
	Input json.RawMessage           `json:"input"`
}

func (AnthropicToolUseBlock) isContentBlock() {}

// AnthropicServerToolUseBlock 服务端工具调用块（如 web_search / code_execution）。
type AnthropicServerToolUseBlock struct {
	Type   AnthropicServerToolUseBlockType `json:"type"` // 固定 "server_tool_use"
	ID     string                          `json:"id"`
	Name   json.RawMessage                 `json:"name"` // 双形态：string | WebSearchToolResult 等
	Input  json.RawMessage                 `json:"input"`
	Caller *AnthropicToolCallerUnionFlat   `json:"caller,omitempty"`
}

func (AnthropicServerToolUseBlock) isContentBlock() {}

// AnthropicWebSearchToolResultBlock web_search 结果块。
type AnthropicWebSearchToolResultBlock struct {
	Type      AnthropicWebSearchToolResultBlockType `json:"type"` // 固定 "web_search_tool_result"
	ToolUseID string                                `json:"tool_use_id"`
	Content   json.RawMessage                       `json:"content"` // 双形态：Result[] | Error
	Caller    *AnthropicToolCallerUnionFlat         `json:"caller,omitempty"`
}

func (AnthropicWebSearchToolResultBlock) isContentBlock() {}

// AnthropicWebFetchToolResultBlock web_fetch 结果块。
type AnthropicWebFetchToolResultBlock struct {
	Type      AnthropicWebFetchToolResultBlockType `json:"type"` // 固定 "web_fetch_tool_result"
	ToolUseID string                               `json:"tool_use_id"`
	Content   json.RawMessage                      `json:"content"` // 双形态：Result | Error
	Caller    *AnthropicToolCallerUnionFlat        `json:"caller,omitempty"`
}

func (AnthropicWebFetchToolResultBlock) isContentBlock() {}

// AnthropicCodeExecutionToolResultBlock code_execution 结果块。
type AnthropicCodeExecutionToolResultBlock struct {
	Type      AnthropicCodeExecutionToolResultBlockType `json:"type"` // 固定 "code_execution_tool_result"
	ToolUseID string                                    `json:"tool_use_id"`
	Content   json.RawMessage                           `json:"content"` // 双形态：Result | Error
}

func (AnthropicCodeExecutionToolResultBlock) isContentBlock() {}

// AnthropicBashCodeExecutionToolResultBlock bash code execution 结果块。
type AnthropicBashCodeExecutionToolResultBlock struct {
	Type      AnthropicBashCodeExecutionToolResultBlockType `json:"type"` // 固定 "bash_code_execution_tool_result"
	ToolUseID string                                        `json:"tool_use_id"`
	Content   json.RawMessage                               `json:"content"` // 双形态：Result | Error
}

func (AnthropicBashCodeExecutionToolResultBlock) isContentBlock() {}

// AnthropicTextEditorCodeExecutionToolResultBlock text_editor code execution 结果块。
type AnthropicTextEditorCodeExecutionToolResultBlock struct {
	Type      AnthropicTextEditorCodeExecutionToolResultBlockType `json:"type"` // 固定 "text_editor_code_execution_tool_result"
	ToolUseID string                                              `json:"tool_use_id"`
	Content   json.RawMessage                                     `json:"content"` // 双形态：Result | Error
}

func (AnthropicTextEditorCodeExecutionToolResultBlock) isContentBlock() {}

// AnthropicToolSearchToolResultBlock tool_search 结果块。
type AnthropicToolSearchToolResultBlock struct {
	Type      AnthropicToolSearchToolResultBlockType `json:"type"` // 固定 "tool_search_tool_result"
	ToolUseID string                                 `json:"tool_use_id"`
	Content   json.RawMessage                        `json:"content"` // 双形态：Result | Error
}

func (AnthropicToolSearchToolResultBlock) isContentBlock() {}

// AnthropicContainerUploadBlock 容器上传块（code_execution 上传产物到容器）。
type AnthropicContainerUploadBlock struct {
	Type   AnthropicContainerUploadBlockType `json:"type"` // 固定 "container_upload"
	FileID string                            `json:"file_id"`
}

func (AnthropicContainerUploadBlock) isContentBlock() {}
