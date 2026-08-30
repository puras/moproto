# `pkg/protocol/anthropic` — Anthropic Messages API Wire-Format

## 概述

本包按 Anthropic 官方 API 参考 + 官方 Go SDK 验证（`anthropic-sdk-go/message.go` 10105 行），
完整定义 Anthropic Messages API 的请求侧 wire-format 类型，覆盖 Claude 4.6 / 4.7 / 4.8 / Fable 5
的字段集。

## 与 `internal/app/gateway/common/anthropic.go` 的关系

| 方面 | common/anthropic.go（保留不动） | pkg/protocol/anthropic（本包） |
|------|--------------------------------|------------------------------|
| 包 | `internal/app/gateway/common` | `pkg/protocol/anthropic` |
| 用途 | 网关适配层 wire-format（已被 adapter/provider 依赖） | Anthropic 协议权威定义 |
| ToolChoice | `any` | 强类型 union |
| ContentBlock | 单 struct + 多字段 | 各 block 独立 struct + union 接口 |
| 字段集 | 基础字段 | 完整字段（含 4.6/4.7/4.8 新增） |
| 依赖 | 不依赖 `pkg/protocol/anthropic` | 不依赖 `common`（避免循环） |

两者当前并存互不依赖。后续迁移路径（不在本 PR 范围）：由独立 PR 完成引用切换 + 删除旧文件。

## 文件清单

```
pkg/protocol/anthropic/
├── doc.go          包级注释
├── request.go      AnthropicRequest 顶层结构（20+ 字段）
├── message.go      AnthropicMessageParam + 12 种 ContentBlock
├── tool.go         Tool union（1 自定义 + 6 预定义）+ ToolChoice union（4 种）
├── thinking.go     ThinkingConfig（type=adaptive + display）
├── output.go       OutputConfig（effort + json_schema + task_budget）
├── context.go      ContextManagement（3 种 edit）
├── mcp.go          McpServerDefinition
├── container.go    Container + Skill
├── source.go       ImageSource / DocumentSource（base64|url|file|text）
└── common.go       CacheControl / Citations / Metadata + ContentBlockKind 枚举
```

## 顶层字段（AnthropicRequest）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| Model | `string` | ✓ | 如 `claude-opus-4-8` |
| MaxTokens | `int` | ✓ | 输出 token 上限 |
| Messages | `[]AnthropicMessageParam` | ✓ | 对话历史 |
| System | `json.RawMessage` | ✗ | `string \| []TextBlockParam` 双形态 |
| Tools | `[]AnthropicTool` | ✗ | 工具定义 union |
| ToolChoice | `AnthropicToolChoice` | ✗ | 强类型 union |
| Stream | `bool` | ✗ | 流式开关 |
| Thinking | `AnthropicThinkingConfig`（interface） | ✗ | 3 mode union：adaptive（4.6+ 推荐）/ enabled（传统 budget_tokens）/ disabled；构造器：`NewAnthropicThinkingConfigAdaptiveParam()` / `NewAnthropicThinkingConfigEnabledParam(budgetTokens)` / `NewAnthropicThinkingConfigDisabledParam()` |
| OutputConfig | `*AnthropicOutputConfig` | ✗ | 4.6+ GA |
| Temperature / TopP / TopK | `*float64 / *int` | ✗ | Opus 4.7+/Fable 5 已禁用 |
| StopSequences | `[]string` | ✗ | |
| Metadata | `*AnthropicMetadata` | ✗ | user_id |
| CacheControl | `*AnthropicCacheControlEphemeral` | ✗ | 顶层自动缓存 |
| ContextManagement | `*AnthropicContextManagement` | ✗ | 长上下文编辑 |
| Container | `json.RawMessage` | ✗ | `string ID \| object` 双形态 |
| McpServers | `[]AnthropicMcpServerDefinition` | ✗ | MCP 连接器（beta） |
| InferenceGeo | `string` | ✗ | 数据驻留地 |
| ServiceTier | `string` | ✗ | `auto` / `standard_only` |
| UserProfileID | `string` | ✗ | user-profiles beta |
| Speed | `string` | ✗ | `fast`（Opus 4.8） |

## ContentBlock union（12 种）

按官方 spec 与 Go SDK 对齐：

| 类型 | type | 关键字段 |
|------|------|----------|
| `AnthropicTextBlockParam` | `text` | Text, CacheControl?, Citations? |
| `AnthropicImageBlockParam` | `image` | Source: ImageSource |
| `AnthropicDocumentBlockParam` | `document` | Source, Title?, Citations? |
| `AnthropicToolUseBlockParam` | `tool_use` | ID, Name, Input (RawMessage), Caller? |
| `AnthropicToolResultBlockParam` | `tool_result` | ToolUseID, Content (RawMessage union), IsError?, CacheControl? |
| `AnthropicThinkingBlockParam` | `thinking` | Thinking, Signature |
| `AnthropicRedactedThinkingBlockParam` | `redacted_thinking` | Data |
| `AnthropicMidConversationSystemBlockParam` | `mid_conv_system` | Content []TextBlock |
| `AnthropicServerToolUseBlockParam` | `server_tool_use` | ID, Name (7 enum), Input, Caller? |
| `AnthropicWebSearchToolResultBlockParam` | `web_search_tool_result` | ToolUseID, Content (Result[] / Error) |
| `AnthropicToolReferenceBlockParam` | `tool_reference` | ToolName (tool search 动态加载) |
| `AnthropicCompactionBlockParam` | `compaction` | Content |

使用模式：

```go
var b anthropic.AnthropicContentBlockParam = anthropic.AnthropicTextBlockParam{
    Type: "text", Text: "hello",
}
switch v := b.(type) {
case anthropic.AnthropicTextBlockParam:    // ...
case anthropic.AnthropicImageBlockParam:   // ...
}
```

## Tool union（1 自定义 + 6 预定义）

```go
type AnthropicTool interface {
    isTool()
    ToolKind() ToolKind
}
```

- **自定义工具**：`AnthropicCustomTool`（含 InputSchema + Strict?）
- **预定义工具**：
  - `AnthropicWebSearchTool` — `web_search_20260209`
  - `AnthropicWebFetchTool` — `web_fetch_20260209`
  - `AnthropicCodeExecutionTool` — `code_execution_20260120`
  - `AnthropicBashTool` — `bash_20250124`
  - `AnthropicTextEditorTool` — `text_editor_20250728`
  - `AnthropicMemoryTool` — `memory_20250818`

## ToolChoice union（4 种策略）

| 类型 | type | 字段 |
|------|------|------|
| `AnthropicToolChoiceAuto` | `auto` | DisableParallelToolUse? |
| `AnthropicToolChoiceAny` | `any` | DisableParallelToolUse? |
| `AnthropicToolChoiceTool` | `tool` | Name, DisableParallelToolUse? |
| `AnthropicToolChoiceNone` | `none` | — |

## ContextManagement（3 种 edit）

- `AnthropicClearToolUses20250919Edit` — `clear_tool_uses_20250919`（清除 tool_use/tool_result）
- `AnthropicClearThinking20251015Edit` — `clear_thinking_20251015`（清除 thinking 块）
- `AnthropicCompact20260112Edit` — `compact_20260112`（压缩为摘要）

## 关键设计决策

1. **ContentBlock 强类型 union**：每个 block 独立 struct + 共享 `AnthropicContentBlockParam` 接口 + `BlockKind()` 方法。
   序列化由各类型自带 JSON tag 处理。`type` 字段在各 struct 内是固定值字面量字符串，编译期可查。

2. **Caller 字段同时存在于 ToolUse / ServerToolUse**：`AnthropicToolCaller` 含 Type + ToolID，
   PTC（Programmatic Tool Calling）场景下 type=`code_execution_20260120` + ToolID 指向调用方 server_tool_use。

3. **ToolUseBlock.Input 使用 `json.RawMessage`**：Go SDK 验证服务端返回的是原始 JSON 字节，
   SDK 消费者必须自行 `json.Unmarshal(block.Input, &v)` 才能拿到结构化数据。

4. **System / Container 双形态**：用 `json.RawMessage` 保留 `string` / object 两种序列化形态，
   反序列化由调用方按需区分。

5. **ToolResultBlockParam.Content 用 RawMessage**：支持 `string | []ContentBlock` 双形态。

6. **WebSearchToolResult.Content 是 union**：`[]WebSearchResultBlock`（成功）| 单个 Error 对象（失败）。
   6 种 ErrorCode：`invalid_tool_input` / `unavailable` / `max_uses_exceeded` /
   `too_many_requests` / `query_too_long` / `request_too_large`。

7. **WebSearch/WebFetchTool 接受 UserLocation**：仅 WebSearch，含 Type / City / Country / Region / Timezone。

8. **ToolReferenceBlock**：tool search 工具动态加载后注入 messages / tool_result.content 的标记块。

9. **ServerToolUseBlock.Name 7 种枚举值**：
   `web_search` / `web_fetch` / `code_execution` / `bash_code_execution` /
   `text_editor_code_execution` / `tool_search_tool_regex` / `tool_search_tool_bm25`。

10. **Container 双形态**：顶层 container 字段接受 string（已存在容器 ID）或 object（含 skills 列表）。

11. **同协议流式 passthrough（raw bytes 透传）**：上游 Anthropic SSE（`event: X\ndata: Y\n\n` 双行格式）经 gateway 转发时，**不**做 typed struct 往返（避免 LOSS #1 未知字段丢失 + LOSS #2 `omitempty` 零值字段剥除）。由 `AnthropicRawStreamFrame` carrier（私有 `isAnthropicRawStreamFrame()` 标记，**不**实现 `AnthropicStreamEvent` 接口，强制 type switch 显式分流）+ `frames.MarshalSSEFrameRaw` verbatim emit 字节完成。typed 路径仍保留 8 个 `AnthropicStreamEvent` 变体供 synthetic 帧构造（`BuildMessageStartFrame` 等）。详见 `WIRE_FORMAT.md` §3.3 与 E2E `TestL2_Streaming_AnthropicPassthrough`。

## 验证

```bash
cd backend
go build ./pkg/protocol/anthropic/...   # 通过
go vet ./pkg/protocol/anthropic/...     # 通过
go build ./...                          # 全量回归通过（不影响其他包）
gofmt -l pkg/protocol/anthropic/        # 无输出（格式干净）
```

## 未纳入本 PR 的项（按需后续补齐）

下列项属于响应侧 / 流式 / 工具结果细分 / 错误模型 / 高级元数据 / 鉴权头，与本 PR 的
"请求侧 Request 定义"范围正交，列出供后续 PR 参考。

### 1. 响应侧 Message 顶层

```go
// 服务端返回的 Message（非流式）
type AnthropicMessage struct {
    ID           string                          // msg_xxx
    Type         string                          // "message"
    Role         string                          // "assistant"
    Content      []AnthropicContentBlock         // 响应块（含 signature）
    Model        string
    StopReason   string                          // end_turn|max_tokens|stop_sequence|tool_use|pause_turn|refusal|model_context_window_exceeded
    StopSequence string                          // 命中的自定义 stop_sequence
    Container    AnthropicResponseContainer      // 代码执行容器信息
    Usage        AnthropicUsage                  // 完整 Usage
}

// 服务端返回的 Container
type AnthropicResponseContainer struct {
    ID        string    `json:"id"`
    ExpiresAt time.Time `json:"expires_at"`
}
```

`Container` 信息也会出现在非流式响应 + 流式 `message_delta` 事件的 delta 里。

### 2. 完整 Usage（响应侧，含 4 个补全字段）

Go SDK 显示 Usage 比本包请求侧关心的字段多 4 个：

```go
type AnthropicUsage struct {
    InputTokens              int64
    OutputTokens             int64
    CacheCreationInputTokens int64
    CacheReadInputTokens     int64
    InferenceGeo             string

    // 以下为本 PR 未纳入的字段（响应侧）
    CacheCreation       AnthropicCacheCreationBreakdown  // 5m / 1h 分桶
    OutputTokensDetails AnthropicOutputTokensDetails    // 含 ThinkingTokens
    ServerToolUse       AnthropicServerToolUsage        // web_search / web_fetch 拆分计数
    ServiceTier         string                          // "standard" | "priority" | "batch"
}

type AnthropicCacheCreationBreakdown struct {
    Ephemeral5mInputTokens int64  // 5m TTL 写入（费率 1.25x）
    Ephemeral1hInputTokens int64  // 1h TTL 写入（费率 2x）
}

type AnthropicOutputTokensDetails struct {
    ThinkingTokens int64  // 思考文本 token（含分隔符，原始 token，非返回给用户可见 thinking）
}

type AnthropicServerToolUsage struct {
    WebSearchRequests int64
    WebFetchRequests  int64
}
```

流式专属的 `MessageDeltaUsage` 字段更少（无 cache_creation / output_tokens_details /
server_tool_use / service_tier 四个字段）——要拿完整 Usage 必须等非流式响应或累积 stream。

### 3. 响应侧 ContentBlock（ContentBlock 响应对应）

请求侧 ContentBlockParam 有 `Param` 后缀；响应侧是对应的"无 Param"版本，差异：

| 请求侧（Param） | 响应侧 | 差异 |
|----------------|--------|------|
| `AnthropicToolUseBlockParam` | `AnthropicToolUseBlock` | 响应侧 Caller 必填（`api:"required"`），Input 用 RawMessage |
| `AnthropicServerToolUseBlockParam` | `AnthropicServerToolUseBlock` | 同上 |
| `AnthropicWebSearchToolResultBlockParam` | `AnthropicWebSearchToolResultBlock` | 响应侧 Caller 必填 |
| — | `AnthropicWebFetchBlock` | 响应侧：Content + RetrievedAt + URL |

### 4. 流式事件类型（5 种 delta）

流式 `content_block_delta` 事件的 delta 字段有 5 种类型（本 PR 未纳入）：

```go
type AnthropicContentBlockDelta interface{ isDelta(); DeltaKind() DeltaKind }

type AnthropicTextDelta struct{ Type string; Text string }                     // text_delta
type AnthropicInputJSONDelta struct{ Type string; PartialJSON string }         // input_json_delta
type AnthropicCitationsDelta struct{ Type string; Citation CitationsDeltaCitationUnion } // citations_delta
type AnthropicThinkingDelta struct{ Type string; Thinking string }              // thinking_delta
type AnthropicSignatureDelta struct{ Type string; Signature string }           // signature_delta
```

`signature_delta` 是 thinking 块闭合的最后一段 delta（携带 signature 字段），用于客户端
完整还原 thinking 块并回传。

### 5. 流式事件枚举（SSE event 类型）

```go
// 流式事件类型（每种 event 有独立 struct）
const (
    StreamEventMessageStart    = "message_start"     // 含初始 Message + MessageUsage
    StreamEventMessageStop     = "message_stop"
    StreamEventMessageDelta    = "message_delta"     // 含 StopReason / StopSequence / Container / Usage
    StreamEventContentBlockStart  = "content_block_start"   // 含起始 ContentBlock
    StreamEventContentBlockDelta  = "content_block_delta"   // 5 种 delta
    StreamEventContentBlockStop   = "content_block_stop"
    StreamEventPing          = "ping"
    StreamEventError         = "error"             // ErrorResponse
)
```

### 6. Citation *LocationParam（5 种引文位置）

`AnthropicCitationsParam` 是 union 接口，对齐官方 anthropic-sdk-go 的
`TextCitationParamUnion`（message.go:5784-5836）。

通过 `isCitationsParam()` 私有方法封闭实现集；调用方通过 type switch 或各实现的
Type 字段分发。

```go
type AnthropicCitationsParam interface {
    isCitationsParam()
}

// 5 种实现按 type 字符串路由
//   - AnthropicCitationCharLocationParam            （"char_location"）
//   - AnthropicCitationPageLocationParam            （"page_location"）
//   - AnthropicCitationContentBlockLocationParam    （"content_block_location"）
//   - AnthropicCitationSearchResultLocationParam    （"search_result_location"）
//   - AnthropicCitationWebSearchResultLocationParam （"web_search_result_location"）
```

各实现字段（对齐官方 SDK `Citation*LocationParam`）：

| Param | 必填字段 | 可选字段 |
|-------|---------|---------|
| `AnthropicCitationCharLocationParam` | Type / CitedText / StartCharIndex / EndCharIndex | DocumentIndex / DocumentTitle |
| `AnthropicCitationPageLocationParam` | Type / CitedText / StartPageNumber / EndPageNumber | DocumentIndex / DocumentTitle / FileID |
| `AnthropicCitationContentBlockLocationParam` | Type / CitedText / StartBlockIndex / EndBlockIndex | DocumentIndex / DocumentTitle / FileID |
| `AnthropicCitationSearchResultLocationParam` | Type / CitedText / StartBlockIndex / EndBlockIndex / SearchResultIndex | Title / Source |
| `AnthropicCitationWebSearchResultLocationParam` | Type / CitedText / EncryptedIndex | Title / URL |

构造器（避免手填 Type 字段导致 wire-format 字符串错误）：

```go
NewAnthropicCitationCharLocationParam(citedText, startChar, endChar)
NewAnthropicCitationPageLocationParam(citedText, startPage, endPage)
NewAnthropicCitationContentBlockLocationParam(citedText, startBlock, endBlock)
NewAnthropicCitationSearchResultLocationParam(citedText, startBlock, endBlock, searchResultIndex)
NewAnthropicCitationWebSearchResultLocationParam(citedText, encryptedIndex)
```

**与 `AnthropicCitationsConfig`（common.go）的关系**：Config 是块级开关（"是否为本
block 生成 citations"），本 union 是具体引文元数据。两者作用层级不同，因此
`message.go` 现有 `Citations *AnthropicCitationsConfig` 字段保持不变；本 union
主要用于响应侧（`content_block.go` 未来的 citations 元数据解析）。

**字段约定**：与官方 SDK 严格对齐——`Start*Index` / `End*Index` 用于字符/页/块索引，
`EncryptedIndex` 用于 web 搜索的加密索引（非 `EncryptedContent`），`FileID` 用于
Files API 文档标识。

### 7. 错误响应

```go
type AnthropicErrorResponse struct {
    Type  string // "error"
    Error AnthropicError
}

type AnthropicError struct {
    Type    string // "api_error" | "authentication_error" | "permission_error" |
                   // "not_found_error" | "rate_limit_error" | "overloaded_error" |
                   // "invalid_request_error" | "billing_error"
    Message string
}
```

### 8. Bash / TextEditor Code Execution 结果细分块

```go
// Bash 工具结果（成功 / 失败 两条分支）
type AnthropicBashCodeExecutionResultBlockParam struct {
    Type       string // "bash_code_execution_result"
    Content    []AnthropicBashCodeExecutionOutputBlockParam  // 引用的文件
    Stdout     string
    Stderr     string
    ReturnCode int64
}

type AnthropicBashCodeExecutionToolResultErrorParam struct {
    Type      string // "bash_code_execution_tool_result_error"
    ErrorCode string
}

// TextEditor 工具三个子结果
type AnthropicTextEditorCodeExecutionViewResultBlockParam      // view
type AnthropicTextEditorCodeExecutionCreateResultBlockParam    // create
type AnthropicTextEditorCodeExecutionStrReplaceResultBlockParam // str_replace
type AnthropicTextEditorCodeExecutionToolResultErrorParam      // 错误分支
```

这些是 `tool_result.content` 数组里的细分块，与 `AnthropicWebSearchToolResultBlockParam`
同级别。

### 9. WebFetch 工具结果（union）

WebFetch 工具的结果块也是 union（成功 / 失败）：

```go
type AnthropicWebFetchToolResultBlockParam struct {
    Type      string // "web_fetch_tool_result"
    ToolUseID string
    Content   json.RawMessage  // 成功: AnthropicWebFetchBlock; 失败: AnthropicWebFetchToolResultError
}

type AnthropicWebFetchBlock struct {
    Type        string // "web_fetch_result"
    URL         string
    RetrievedAt string  // ISO 8601
    Content     AnthropicDocumentBlock  // 文档块
}
```

### 10. 鉴权头（HTTP header）

服务端通过 header 而非 body 接收的部分（不属于请求 wire-format JSON）：

```go
const (
    HeaderAPIKey         = "x-api-key"                    // 必填
    HeaderAuthToken      = "Authorization"                // 替代方案：Bearer
    HeaderAnthropicVersion = "anthropic-version"          // 如 "2023-06-01"
    HeaderAnthropicBeta  = "anthropic-beta"               // 多个 beta 用换行分隔
    HeaderUserProfileID  = "anthropic-user-profile-id"    // user-profiles beta
)
```

### 11. CountTokens 子端点

`POST /v1/messages/count_tokens` 返回 `MessageTokensCount{ InputTokens }`，用于预算预估。
本 PR 不涉及 MessageTokensCount 类型。

### 12. 单元测试

`pkg/protocol/anthropic/*_test.go` 套件未纳入本 PR。建议覆盖：

- Marshal / Unmarshal 往返：每个 ContentBlock kind、Tool kind、ToolChoice kind
- AnthropicRequest 完整序列化 + 各 omitempty 字段验证
- Union 类型断言（`switch v := b.(type)`）
- 基于 Anthropic 官方 API 文档样例的 fixture 测试

## 参考

- Anthropic 官方 Messages API：<https://docs.anthropic.com/en/api/messages>
- 官方 Go SDK：`/Users/puras/workspace/source/anthropic-sdk-go/message.go`（10105 行，作为字段集基线）
- 现有 `internal/app/gateway/common/anthropic.go`（保留不动，作为网关适配层 wire-format 现状）
- 现有 `common/content.go` ContentPart 接口模式作为 union 类型设计模板