package anthropic

import "encoding/json"

// ============== ContentBlock Type 常量 ==============

// AnthropicTextBlockType text block 的 type 字段类型。
type AnthropicTextBlockType string

const AnthropicTextBlockTypeText AnthropicTextBlockType = "text"

// AnthropicImageBlockType image block 的 type 字段类型。
type AnthropicImageBlockType string

const AnthropicImageBlockTypeImage AnthropicImageBlockType = "image"

// AnthropicDocumentBlockType document block 的 type 字段类型。
type AnthropicDocumentBlockType string

const AnthropicDocumentBlockTypeDocument AnthropicDocumentBlockType = "document"

// AnthropicToolUseBlockType tool_use block 的 type 字段类型。
type AnthropicToolUseBlockType string

const AnthropicToolUseBlockTypeToolUse AnthropicToolUseBlockType = "tool_use"

// AnthropicToolResultBlockType tool_result block 的 type 字段类型。
type AnthropicToolResultBlockType string

const AnthropicToolResultBlockTypeToolResult AnthropicToolResultBlockType = "tool_result"

// AnthropicThinkingBlockType thinking block 的 type 字段类型。
type AnthropicThinkingBlockType string

const AnthropicThinkingBlockTypeThinking AnthropicThinkingBlockType = "thinking"

// AnthropicRedactedThinkingBlockType redacted_thinking block 的 type 字段类型。
type AnthropicRedactedThinkingBlockType string

const AnthropicRedactedThinkingBlockTypeRedactedThinking AnthropicRedactedThinkingBlockType = "redacted_thinking"

// AnthropicMidConversationSystemBlockType mid_conv_system block 的 type 字段类型。
type AnthropicMidConversationSystemBlockType string

const AnthropicMidConversationSystemBlockTypeMidConvSystem AnthropicMidConversationSystemBlockType = "mid_conv_system"

// AnthropicServerToolUseBlockType server_tool_use block 的 type 字段类型。
type AnthropicServerToolUseBlockType string

const AnthropicServerToolUseBlockTypeServerToolUse AnthropicServerToolUseBlockType = "server_tool_use"

// AnthropicWebSearchToolResultBlockType web_search_tool_result block 的 type 字段类型。
type AnthropicWebSearchToolResultBlockType string

const AnthropicWebSearchToolResultBlockTypeWebSearchToolResult AnthropicWebSearchToolResultBlockType = "web_search_tool_result"

// AnthropicWebSearchResultBlockType web_search_result block 的 type 字段类型。
type AnthropicWebSearchResultBlockType string

const AnthropicWebSearchResultBlockTypeWebSearchResult AnthropicWebSearchResultBlockType = "web_search_result"

// AnthropicToolReferenceBlockType tool_reference block 的 type 字段类型。
type AnthropicToolReferenceBlockType string

const AnthropicToolReferenceBlockTypeToolReference AnthropicToolReferenceBlockType = "tool_reference"

// AnthropicCompactionBlockType compaction block 的 type 字段类型。
type AnthropicCompactionBlockType string

const AnthropicCompactionBlockTypeCompaction AnthropicCompactionBlockType = "compaction"

// AnthropicSearchResultBlockType search_result block 的 type 字段类型。
type AnthropicSearchResultBlockType string

const AnthropicSearchResultBlockTypeSearchResult AnthropicSearchResultBlockType = "search_result"

// AnthropicContentBlockParam 多态内容块 union。
//
// Anthropic Messages API 的 message.content 字段是异构数组，每一项可以是文本、图片、
// 文档、工具调用、工具结果、思考块、mid_conv_system、tool_reference 等。
// 本接口配合具体类型实现，序列化由各类型自带 JSON tag 处理。
//
// 使用模式：
//
//	var b AnthropicContentBlockParam = AnthropicTextBlockParam{Text: "hello"}
//	switch v := b.(type) {
//	case AnthropicTextBlockParam:    ...
//	case AnthropicImageBlockParam:   ...
//	}
//
// 分发机制：
//   - isContentBlockParam() 私有方法封闭实现集（外部包无法新增实现）；
//   - 各实现 struct 的 Type 字段承载 wire-format 类型标签（如 "text" / "tool_use"），
//     调用方可直接 switch on Type 字段做分发，参考 anthropic-sdk-go 的
//     ContentBlockUnion.AsAny() 模式。
//
// 服务端追加新块类型时只新增实现 + 对应 Type 常量，本接口保持稳定。
type AnthropicContentBlockParam interface {
	isContentBlockParam()
}

// AnthropicMessageParam Anthropic 单条消息（请求方向）。
//
// Role 取值："user" | "assistant" | "system"（system 用于 mid-conversation，
// 见 AnthropicMidConversationSystemBlockParam）。
//
// Content 支持两种形态（服务端均接受）：
//   - 简单文本：{"content": "hello"}   → 解码为 string
//   - 结构化：  {"content": [{...}, ...]} → 解码为 []AnthropicContentBlockParam
//
// 使用 json.RawMessage 保留两种形态，反序列化由调用方自行区分。
type AnthropicMessageParam struct {
	Role    AnthropicMessageRole `json:"role"`
	Content json.RawMessage      `json:"content"`
}

// ============== 文本块 ==============

// AnthropicTextBlockParam 文本内容块（user 角色最常用，也出现在 system / tool_result）。
type AnthropicTextBlockParam struct {
	Type         AnthropicTextBlockType          `json:"type"` // 固定 "text"
	Text         string                          `json:"text"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
	Citations    *AnthropicCitationsConfig       `json:"citations,omitempty"`
}

func (AnthropicTextBlockParam) isContentBlockParam() {}

// ============== 图片块 ==============

// AnthropicImageBlockParam 图片内容块。
type AnthropicImageBlockParam struct {
	Type   AnthropicImageBlockType `json:"type"`   // 固定 "image"
	Source AnthropicImageSource    `json:"source"` // interface：AnthropicImageSourceBase64 / AnthropicImageSourceURL
}

func (AnthropicImageBlockParam) isContentBlockParam() {}

// ============== 文档块 ==============

// AnthropicDocumentBlockParam 文档内容块（PDF / 文本 / Office 等）。
type AnthropicDocumentBlockParam struct {
	Type      AnthropicDocumentBlockType `json:"type"`   // 固定 "document"
	Source    AnthropicDocumentSource    `json:"source"` // interface：5 种实现
	Title     string                     `json:"title,omitempty"`
	Citations *AnthropicCitationsConfig  `json:"citations,omitempty"`
}

func (AnthropicDocumentBlockParam) isContentBlockParam() {}

// ============== 工具调用块（assistant 角色专用） ==============

// AnthropicToolCaller 调用者来源（用于 PTC：Programmatic Tool Calling）。
//
// 取值：
//   - "direct"                    — 直接由模型调用
//   - "code_execution_20250825"   — code_execution 容器内部调用
//   - "code_execution_20260120"   — code_execution 容器内部调用（新版）
//
// 对齐官方 anthropic-sdk-go ToolUseBlockParamCallerUnion / ServerToolUseBlockParamCallerUnion
// / WebFetchToolResultBlockParamCallerUnion / WebSearchToolResultBlockParamCallerUnion
// （均为相同 3 态）。本包提供 union interface + 3 个具体 struct。
type AnthropicToolCallerUnion interface {
	isToolCaller()
}

// AnthropicToolCallerDirect 调用者 = 直接模型调用（无 ToolID）。
type AnthropicToolCallerDirect struct {
	Type AnthropicToolCallerDirectType `json:"type"` // 固定 "direct"
}

func (AnthropicToolCallerDirect) isToolCaller() {}

// AnthropicToolCallerUnionFlat 扁平实现（解决 interface unmarshal 问题）。
type AnthropicToolCallerUnionFlat struct {
	Type                  string          `json:"type"`
	Direct                json.RawMessage `json:"direct,omitempty"`
	CodeExecution20250825 json.RawMessage `json:"code_execution_20250825,omitempty"`
	CodeExecution20260120 json.RawMessage `json:"code_execution_20260120,omitempty"`
}

func (AnthropicToolCallerUnionFlat) isToolCaller() {}

func (u AnthropicToolCallerUnionFlat) MarshalJSON() ([]byte, error) {
	// 空 Type 或 nil RawMessage 返回 nil 以触发 omitempty suppression
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "direct":
		if u.Direct == nil {
			return nil, nil
		}
		return u.Direct, nil
	case "code_execution_20250825":
		if u.CodeExecution20250825 == nil {
			return nil, nil
		}
		return u.CodeExecution20250825, nil
	case "code_execution_20260120":
		if u.CodeExecution20260120 == nil {
			return nil, nil
		}
		return u.CodeExecution20260120, nil
	default:
		return nil, nil
	}
}

func (u *AnthropicToolCallerUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type

	switch tmp.Type {
	case "direct":
		var b AnthropicToolCallerDirect
		if err := json.Unmarshal(data, &b); err == nil {
			u.Direct, _ = json.Marshal(b)
		}
	case "code_execution_20250825":
		var b AnthropicToolCallerCodeExecution20250825
		if err := json.Unmarshal(data, &b); err == nil {
			u.CodeExecution20250825, _ = json.Marshal(b)
		}
	case "code_execution_20260120":
		var b AnthropicToolCallerCodeExecution20260120
		if err := json.Unmarshal(data, &b); err == nil {
			u.CodeExecution20260120, _ = json.Marshal(b)
		}
	}
	return nil
}

// AnthropicToolCallerCodeExecution20250825 调用者 = code_execution 容器 20250825 版本。
type AnthropicToolCallerCodeExecution20250825 struct {
	Type   AnthropicToolCallerCodeExecution20250825Type `json:"type"` // 固定 "code_execution_20250825"
	ToolID string                                       `json:"tool_id"`
}

func (AnthropicToolCallerCodeExecution20250825) isToolCaller() {}

// AnthropicToolCallerCodeExecution20260120 调用者 = code_execution 容器 20260120 版本。
type AnthropicToolCallerCodeExecution20260120 struct {
	Type   AnthropicToolCallerCodeExecution20260120Type `json:"type"` // 固定 "code_execution_20260120"
	ToolID string                                       `json:"tool_id"`
}

func (AnthropicToolCallerCodeExecution20260120) isToolCaller() {}

// AnthropicToolUseBlockParam assistant 角色：决定调用工具。
//
// Input 字段在 Param 方向（请求侧）可接受任意结构，由调用方按 tool input_schema 自行反序列化。
// 在响应方向（ToolUseBlock）Input 是 json.RawMessage，由 SDK 消费者自己解析。
type AnthropicToolUseBlockParam struct {
	Type         AnthropicToolUseBlockType       `json:"type"` // 固定 "tool_use"
	ID           string                          `json:"id"`
	Name         string                          `json:"name"`
	Input        json.RawMessage                 `json:"input"`
	Caller       *AnthropicToolCallerUnionFlat   `json:"caller,omitempty"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicToolUseBlockParam) isContentBlockParam() {}

// AnthropicToolResultBlockContentUnion tool_result.content 的元素类型。
//
// 服务端对 tool_result.content 的接受形式：
//   - 字符串（AnthropicToolResultBlockContentText）：纯文本结果
//   - 块数组：富内容，每一项可是 text / image / search_result / document / tool_reference
//
// 对齐官方 anthropic-sdk-go ToolResultBlockParamContentUnion（message.go:7033-7040）。
// 本包提供 5 个具体 struct + 1 个 string 简化形态，全部实现 isToolResultContent()。
// AnthropicToolResultBlockParam.Content 仍保留为 json.RawMessage 以兼容"嵌套 string"
// 形态（由调用方按需做 switch on type）。
type AnthropicToolResultBlockContentUnion interface {
	isToolResultContent()
}

// AnthropicToolResultBlockContentText 工具结果为纯字符串。
//
// 对齐官方 SDK 中 ToolResultBlockParamContentUnion 的字符串形态（被 inline 进
// json:"content" 字段时序列化为 string 而非对象）。
type AnthropicToolResultBlockContentText struct {
	Text string `json:"-"`
}

func (AnthropicToolResultBlockContentText) isToolResultContent() {}

// AnthropicToolResultBlockContentImage 工具结果为图片块。
type AnthropicToolResultBlockContentImage struct {
	Image AnthropicImageBlockParam `json:"image"`
}

func (AnthropicToolResultBlockContentImage) isToolResultContent() {}

// AnthropicToolResultBlockContentSearchResult 工具结果为搜索结果块。
type AnthropicToolResultBlockContentSearchResult struct {
	SearchResult AnthropicSearchResultBlockParam `json:"search_result"`
}

func (AnthropicToolResultBlockContentSearchResult) isToolResultContent() {}

// AnthropicToolResultBlockContentDocument 工具结果为文档块。
type AnthropicToolResultBlockContentDocument struct {
	Document AnthropicDocumentBlockParam `json:"document"`
}

func (AnthropicToolResultBlockContentDocument) isToolResultContent() {}

// AnthropicToolResultBlockContentToolReference 工具结果为工具引用块（tool_search 场景）。
type AnthropicToolResultBlockContentToolReference struct {
	ToolReference AnthropicToolReferenceBlockParam `json:"tool_reference"`
}

func (AnthropicToolResultBlockContentToolReference) isToolResultContent() {}

// AnthropicToolResultBlockParam user 角色：返回工具执行结果。
//
// Content 字段是 union：string | []AnthropicContentBlockParam。用 json.RawMessage 保留双形态。
type AnthropicToolResultBlockParam struct {
	Type         AnthropicToolResultBlockType    `json:"type"` // 固定 "tool_result"
	ToolUseID    string                          `json:"tool_use_id"`
	Content      json.RawMessage                 `json:"content"` // string | []AnthropicContentBlockParam
	IsError      bool                            `json:"is_error,omitempty"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicToolResultBlockParam) isContentBlockParam() {}

// ============== 思考块（4.6/4.7/4.8 回传） ==============

// AnthropicThinkingBlockParam 思考内容（assistant 回传给客户端，client 再回传时使用）。
type AnthropicThinkingBlockParam struct {
	Type      AnthropicThinkingBlockType `json:"type"` // 固定 "thinking"
	Thinking  string                     `json:"thinking"`
	Signature string                     `json:"signature"`
}

func (AnthropicThinkingBlockParam) isContentBlockParam() {}

// AnthropicRedactedThinkingBlockParam 加密的思考块（服务器出于安全原因替换时返回）。
type AnthropicRedactedThinkingBlockParam struct {
	Type AnthropicRedactedThinkingBlockType `json:"type"` // 固定 "redacted_thinking"
	Data string                             `json:"data"`
}

func (AnthropicRedactedThinkingBlockParam) isContentBlockParam() {}

// ============== Mid-Conversation System 块 ==============

// AnthropicMidConversationSystemBlockParam 对应 type="mid_conv_system"。
//
// 用于在 messages 数组中插入或更新 system 级指令（区别于顶层 system 字段）。
// 当 MessageParam.Role == "system" 时，Content 字段的每一项都应是这个块。
type AnthropicMidConversationSystemBlockParam struct {
	Type         AnthropicMidConversationSystemBlockType `json:"type"` // 固定 "mid_conv_system"
	Content      []AnthropicTextBlockParam               `json:"content"`
	CacheControl *AnthropicCacheControlEphemeral         `json:"cache_control,omitempty"`
}

func (AnthropicMidConversationSystemBlockParam) isContentBlockParam() {}

// ============== 服务端工具调用块 ==============

// ServerToolUseBlockName 服务端工具名枚举。
type ServerToolUseBlockName string

const (
	ServerToolUseBlockNameWebSearch               ServerToolUseBlockName = "web_search"
	ServerToolUseBlockNameWebFetch                ServerToolUseBlockName = "web_fetch"
	ServerToolUseBlockNameCodeExecution           ServerToolUseBlockName = "code_execution"
	ServerToolUseBlockNameBashCodeExecution       ServerToolUseBlockName = "bash_code_execution"
	ServerToolUseBlockNameTextEditorCodeExecution ServerToolUseBlockName = "text_editor_code_execution"
	ServerToolUseBlockNameToolSearchToolRegex     ServerToolUseBlockName = "tool_search_tool_regex"
	ServerToolUseBlockNameToolSearchToolBm25      ServerToolUseBlockName = "tool_search_tool_bm25"
)

// AnthropicServerToolUseBlockParam 服务端工具调用（web_search / web_fetch / code_execution 等）。
//
// 7 种 name 取值见 ServerToolUseBlockName 常量。
type AnthropicServerToolUseBlockParam struct {
	Type         AnthropicServerToolUseBlockType `json:"type"` // 固定 "server_tool_use"
	ID           string                          `json:"id"`
	Name         ServerToolUseBlockName          `json:"name"`
	Input        json.RawMessage                 `json:"input"`
	Caller       *AnthropicToolCallerUnionFlat   `json:"caller,omitempty"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicServerToolUseBlockParam) isContentBlockParam() {}

// ============== Web Search Tool 结果块 ==============

// AnthropicWebSearchResultBlockParam 单条 web 搜索结果。
//
// 成功分支：content 是 []AnthropicWebSearchResultBlockParam
// 失败分支：content 是 *AnthropicWebSearchToolRequestErrorParam
type AnthropicWebSearchResultBlockParam struct {
	Type             AnthropicWebSearchResultBlockType `json:"type"` // 固定 "web_search_result"
	EncryptedContent string                            `json:"encrypted_content,omitempty"`
	Title            string                            `json:"title,omitempty"`
	URL              string                            `json:"url,omitempty"`
	PageAge          string                            `json:"page_age,omitempty"`
}

// AnthropicWebSearchToolRequestErrorType web_search_tool_result_error 的 type 字段类型。
type AnthropicWebSearchToolRequestErrorType string

const AnthropicWebSearchToolRequestErrorTypeWebSearchToolResultError AnthropicWebSearchToolRequestErrorType = "web_search_tool_result_error"

// AnthropicWebSearchToolRequestErrorParam web_search 工具错误。
type AnthropicWebSearchToolRequestErrorParam struct {
	Type      AnthropicWebSearchToolRequestErrorType `json:"type"` // 固定 "web_search_tool_result_error"
	ErrorCode string                                 `json:"error_code"`
}

// WebSearchToolResultErrorCode web_search 错误码枚举。
type WebSearchToolResultErrorCode string

const (
	WebSearchToolResultErrorCodeInvalidToolInput WebSearchToolResultErrorCode = "invalid_tool_input"
	WebSearchToolResultErrorCodeUnavailable      WebSearchToolResultErrorCode = "unavailable"
	WebSearchToolResultErrorCodeMaxUsesExceeded  WebSearchToolResultErrorCode = "max_uses_exceeded"
	WebSearchToolResultErrorCodeTooManyRequests  WebSearchToolResultErrorCode = "too_many_requests"
	WebSearchToolResultErrorCodeQueryTooLong     WebSearchToolResultErrorCode = "query_too_long"
	WebSearchToolResultErrorCodeRequestTooLarge  WebSearchToolResultErrorCode = "request_too_large"
)

// AnthropicWebSearchToolResultBlockParam web_search 工具结果（union 容器）。
//
// Content 字段是 union：[]AnthropicWebSearchResultBlockParam（成功）
// | *AnthropicWebSearchToolRequestErrorParam（失败）。
// 实际序列化形态通过 Content 的 JSON 结构由服务端判定。
type AnthropicWebSearchToolResultBlockParam struct {
	Type         AnthropicWebSearchToolResultBlockType `json:"type"` // 固定 "web_search_tool_result"
	ToolUseID    string                                `json:"tool_use_id"`
	Content      json.RawMessage                       `json:"content"` // []Result | ErrorObject
	Caller       *AnthropicToolCallerUnionFlat         `json:"caller,omitempty"`
	CacheControl *AnthropicCacheControlEphemeral       `json:"cache_control,omitempty"`
}

func (AnthropicWebSearchToolResultBlockParam) isContentBlockParam() {}

// ============== Tool Reference 块（tool search 动态加载） ==============

// AnthropicToolReferenceBlockParam tool_reference 块。
//
// 用于在 messages / tool_result.content 中标记某个工具已被 tool search 动态加载。
type AnthropicToolReferenceBlockParam struct {
	Type         AnthropicToolReferenceBlockType `json:"type"` // 固定 "tool_reference"
	ToolName     string                          `json:"tool_name"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicToolReferenceBlockParam) isContentBlockParam() {}

// ============== Compaction 块（长上下文压缩） ==============

// AnthropicCompactionBlockParam 长上下文压缩块（用于恢复历史摘要）。
type AnthropicCompactionBlockParam struct {
	Type    AnthropicCompactionBlockType `json:"type"` // 固定 "compaction"
	Content json.RawMessage              `json:"content"`
}

func (AnthropicCompactionBlockParam) isContentBlockParam() {}

// ============== Search Result 块（search 工具返回） ==============

// AnthropicSearchResultBlockParam 单条 search 工具结果（出现在 tool_result.content 数组中）。
//
// 对齐官方 anthropic-sdk-go SearchResultBlockParam（message.go:5244-5254）。
// 必填：Type + Content + Source + Title；CacheControl / Citations 可选。
type AnthropicSearchResultBlockParam struct {
	Type         AnthropicSearchResultBlockType  `json:"type"` // 固定 "search_result"
	Source       string                          `json:"source"`
	Title        string                          `json:"title"`
	Content      []AnthropicTextBlockParam       `json:"content"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
	Citations    *AnthropicCitationsConfig       `json:"citations,omitempty"`
}

func (AnthropicSearchResultBlockParam) isContentBlockParam() {}

// ============== Web Fetch 工具结果块 ==============

// AnthropicWebFetchToolResultBlockParam web_fetch 工具结果（user 角色回传）。
//
// 对齐官方 anthropic-sdk-go WebFetchToolResultBlockParam（message.go:8914-8925）。
// 必填：Type + ToolUseID + Content；Content 是 union（Result | Error），用 json.RawMessage 透传。
type AnthropicWebFetchToolResultBlockParam struct {
	Type         AnthropicWebFetchToolResultBlockType `json:"type"` // 固定 "web_fetch_tool_result"
	ToolUseID    string                               `json:"tool_use_id"`
	Content      json.RawMessage                      `json:"content"` // 双形态：Result | Error
	Caller       *AnthropicToolCallerUnionFlat        `json:"caller,omitempty"`
	CacheControl *AnthropicCacheControlEphemeral      `json:"cache_control,omitempty"`
}

func (AnthropicWebFetchToolResultBlockParam) isContentBlockParam() {}

// ============== Code Execution 工具结果块（3 种） ==============

// AnthropicCodeExecutionToolResultBlockParam code_execution 工具结果（user 角色回传）。
//
// 对齐官方 anthropic-sdk-go CodeExecutionToolResultBlockParam（message.go:1313-1323）。
type AnthropicCodeExecutionToolResultBlockParam struct {
	Type         AnthropicCodeExecutionToolResultBlockType `json:"type"` // 固定 "code_execution_tool_result"
	ToolUseID    string                                    `json:"tool_use_id"`
	Content      json.RawMessage                           `json:"content"` // 双形态：Result | Error
	CacheControl *AnthropicCacheControlEphemeral           `json:"cache_control,omitempty"`
}

func (AnthropicCodeExecutionToolResultBlockParam) isContentBlockParam() {}

// AnthropicBashCodeExecutionToolResultBlockParam bash_code_execution 工具结果（user 角色回传）。
//
// 对齐官方 anthropic-sdk-go BashCodeExecutionToolResultBlockParam（message.go:338-352）。
type AnthropicBashCodeExecutionToolResultBlockParam struct {
	Type         AnthropicBashCodeExecutionToolResultBlockType `json:"type"` // 固定 "bash_code_execution_tool_result"
	ToolUseID    string                                        `json:"tool_use_id"`
	Content      json.RawMessage                               `json:"content"` // 双形态：Result | Error
	CacheControl *AnthropicCacheControlEphemeral               `json:"cache_control,omitempty"`
}

func (AnthropicBashCodeExecutionToolResultBlockParam) isContentBlockParam() {}

// AnthropicTextEditorCodeExecutionToolResultBlockParam text_editor_code_execution 工具结果（user 角色回传）。
//
// 对齐官方 anthropic-sdk-go TextEditorCodeExecutionToolResultBlockParam（message.go:6182-6192）。
type AnthropicTextEditorCodeExecutionToolResultBlockParam struct {
	Type         AnthropicTextEditorCodeExecutionToolResultBlockType `json:"type"` // 固定 "text_editor_code_execution_tool_result"
	ToolUseID    string                                              `json:"tool_use_id"`
	Content      json.RawMessage                                     `json:"content"` // 双形态：Result | Error
	CacheControl *AnthropicCacheControlEphemeral                     `json:"cache_control,omitempty"`
}

func (AnthropicTextEditorCodeExecutionToolResultBlockParam) isContentBlockParam() {}

// ============== Tool Search 工具结果块 ==============

// AnthropicToolSearchToolResultBlockParam tool_search 工具结果（user 角色回传）。
//
// 对齐官方 anthropic-sdk-go ToolSearchToolResultBlockParam（message.go:7413-7422）。
type AnthropicToolSearchToolResultBlockParam struct {
	Type         AnthropicToolSearchToolResultBlockType `json:"type"` // 固定 "tool_search_tool_result"
	ToolUseID    string                                 `json:"tool_use_id"`
	Content      json.RawMessage                        `json:"content"` // 双形态：Result | Error
	CacheControl *AnthropicCacheControlEphemeral        `json:"cache_control,omitempty"`
}

func (AnthropicToolSearchToolResultBlockParam) isContentBlockParam() {}
