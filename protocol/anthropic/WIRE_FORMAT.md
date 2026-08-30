# Anthropic Messages API — 流式/非流式 Request/Response 完整结构

> 参照:
> - Anthropic 官方 Messages API 文档
> - 官方 Go SDK（`anthropic-sdk-go/message.go`，10105 行，作为字段集基线）
> - ClaudeCode TypeScript 源码（流式累积模式与思考块回传策略）
>
> 覆盖范围: Claude Opus 4.6 / 4.7 / 4.8 / Fable 5。
>
> 本文档为纯结构梳理,描述 wire-format JSON 形状与字段约束,不涉及任何代码类型。
>
> Request 与 Response 的两种形态(**流式 vs 非流式**)共享同一份请求体,仅由 `stream` 字段区分。

---

## 一、Request 对象(流式 & 非流式共用)

### 1.1 顶层结构

```
POST /v1/messages
Content-Type: application/json
x-api-key: <API_KEY>                 # 鉴权(亦可用 Authorization: Bearer)
anthropic-version: 2023-06-01        # API 版本
anthropic-beta: <id1>,<id2>          # 多 beta 用英文逗号分隔

{
  // ──────────────── 必填字段 ────────────────
  "model": "claude-opus-4-8",                // string,如 claude-opus-4-8 / claude-fable-5
  "max_tokens": 16000,                       // int,输出 token 上限;4.6+ 支持到 128K(流式)
  "messages": [MessageParam, ...],           // 见 1.2

  // ──────────────── 系统提示 ────────────────
  "system": string | TextBlockParam[],       // string 简写 或 结构化块数组
                                              // 结构化时每块可带 cache_control

  // ──────────────── 工具 ────────────────
  "tools": [Tool, ...],                      // 见 1.3,1 自定义 + 6 预定义 union
  "tool_choice": ToolChoice,                 // 见 1.4,默认 auto

  // ──────────────── 流式开关(唯一差异点) ────────────────
  "stream": true | false,                    // 省略或 false = 非流式;true = SSE

  // ──────────────── 思考 / 输出控制(4.6+) ────────────────
  "thinking": {                              // 见 1.5
    "type": "adaptive",
    "display": "summarized" | "omitted"      // 4.7/4.8/Fable 5 控制 thinking 文本可见性
  },
  "output_config": {                         // 见 1.6,GA on 4.6
    "effort": "low" | "medium" | "high" | "xhigh" | "max",
    "format": JsonOutputFormat,              // 结构化输出
    "task_budget": { "type": "tokens", "total": N }   // beta: agentic 任务预算
  },

  // ──────────────── 采样参数(⚠️ Opus 4.7+/Fable 5 已禁用) ────────────────
  "temperature": 0.7,                        // float,4.7+ 传则 400
  "top_p": 0.9,                              // float,4.7+ 传则 400
  "top_k": 5,                                // int,  4.7+ 传则 400

  // ──────────────── 行为控制 ────────────────
  "stop_sequences": ["END"],                 // string[],自定义停止序列
  "metadata": { "user_id": "u_123" },        // 请求级元数据
  "cache_control": {                         // 顶层自动缓存(标最后一个可缓存块)
    "type": "ephemeral",
    "ttl":   "5m" | "1h"                     // 默认 5m,1h 双倍写入费
  },

  // ──────────────── 长上下文管理(beta) ────────────────
  "context_management": {
    "edits": [
      { "type": "clear_tool_uses_20250919", "clear_tool_inputs": true? },
      { "type": "clear_thinking_20251015" },
      { "type": "compact_20260112" }
    ]
  },

  // ──────────────── MCP / 容器 / 服务控制 ────────────────
  "mcp_servers": [McpServerDefinition, ...],  // beta: 服务端 MCP 连接
  "container": string | {                     // 代码执行容器,双形态
    "id":   "container_xxx",                  //   string = 已存在容器 ID
    "skills": [{ "type": "anthropic", "skill_id": "pptx", "version": "latest" }, ...]
  },
  "inference_geo": "us",                     // 数据驻留地(1P only)
  "service_tier": "auto" | "standard_only",  // 优先级
  "user_profile_id": "up_xxx",               // user-profiles beta
  "speed": "fast",                           // Opus 4.8 fast mode(beta)

  // ──────────────── 兼容性已禁用 ────────────────
  // "output_format": { ... }                // ⚠️ 旧字段;改用 output_config.format
  // 末位 assistant 消息 prefill             // ⚠️ Opus 4.6+ / Fable 5 传则 400
}
```

### 1.2 MessageParam

```
MessageParam
├── role: "user" | "assistant" | "system"     // system 仅支持中段插入(Opus 4.8+)
└── content: string | ContentBlockParam[]
```

`content` 可以是字符串(纯文本)或块数组。块数组见 1.7。

### 1.3 Tool union(1 自定义 + 6 预定义)

**自定义工具**(需 `input_schema`):

```
{
  "name": "get_weather",
  "description": "Get current weather for a city.",
  "input_schema": {                            // JSON Schema
    "type": "object",
    "properties": { ... },
    "required": [...],
    "additionalProperties": false               // strict 模式必填
  },
  "strict": true,                               // GA: 严格 schema 校验
  "cache_control": { "type": "ephemeral" }
}
```

**预定义工具**(无需 schema,声明 `type` + `name`):

| 工具 | type | name | 客户端/服务端 |
|------|------|------|--------------|
| Web Search | `web_search_20260209` | `web_search` | 服务端 |
| Web Search 旧版 | `web_search_20250305` | `web_search` | 服务端 |
| Web Fetch | `web_fetch_20260209` | `web_fetch` | 服务端 |
| Web Fetch 旧版 | `web_fetch_20250910` | `web_fetch` | 服务端 |
| Code Execution | `code_execution_20260120` | `code_execution` | 服务端 |
| Bash | `bash_20250124` | `bash` | 客户端 |
| Text Editor | `text_editor_20250728` | `str_replace_based_edit_tool` | 客户端 |
| Memory | `memory_20250818` | `memory` | 客户端 |

预定义工具的字段差异(均可选):

```
web_search_* : {
  "max_uses": 5,
  "allowed_domains": ["example.com"],
  "blocked_domains": ["spam.com"],
  "user_location": { "type": "approximate", "city": "SF", "country": "US", ... },
  "defer_loading": false,
  "cache_control": { ... }
}

web_fetch_* : {
  "max_uses": 5,
  "max_content_tokens": 100000,
  "allowed_domains": [...],
  "blocked_domains": [...],
  "citations": { "enabled": true },
  "defer_loading": false,
  "cache_control": { ... }
}

code_execution_20260120 : { /* 无字段 */ }

bash_20250124 : { /* 无字段 */ }

text_editor_20250728 : {
  "max_characters": 100000,                     // view 输出上限
  "cache_control": { ... }
}

memory_20250818 : { /* 无字段 */ }
```

### 1.4 ToolChoice union(4 种策略)

```
{ "type": "auto", "disable_parallel_tool_use": false }    // 默认,模型自决
{ "type": "any",  "disable_parallel_tool_use": false }    // 必须调用某个工具
{ "type": "tool", "name": "get_weather", "disable_parallel_tool_use": false }
{ "type": "none" }                                        // 禁用工具
```

### 1.5 ThinkingConfig

```
// 4.6+ 推荐
{ "type": "adaptive", "display": "summarized" }   // 4.7/4.8/Fable 5:display 控制 thinking 文本可见
{ "type": "adaptive", "display": "omitted" }      // 4.7/4.8/Fable 5 默认,thinking 块文本为空

// Opus 4.6 过渡(旧模型兼容),新代码不应使用
{ "type": "enabled", "budget_tokens": 8192 }      // ⚠️ 4.7+/Fable 5 传则 400

// 旧模型
{ "type": "disabled" }                            // Fable 5 传则 400,Opus 4.8/4.7 接受
```

### 1.6 OutputConfig

```
{
  "effort": "high",                            // 低/中/高/max;xhigh 是 Opus 4.7+ 专属
  "format": {                                  // 结构化输出,GA on 4.6
    "type": "json_schema",
    "schema": { ... }                          // 必须 additionalProperties: false
  },
  "task_budget": {                             // beta,agentic 任务预算
    "type": "tokens", "total": 64000           // 最小 20000
  }
}
```

### 1.7 ContentBlockParam union(12 种请求侧块)

按 type 字段路由,序列化形式以 type 字面量区分。

#### 文本 / 媒体

```
TextBlockParam:
{ "type": "text", "text": "...", "cache_control"?: ..., "citations"?: CitationsConfig }

ImageBlockParam:
{ "type": "image", "source": ImageSource, "cache_control"?: ... }
  ImageSource:
    { "type": "base64", "media_type": "image/png" | "image/jpeg" | "image/gif" | "image/webp", "data": "<b64>" }
    { "type": "url", "url": "https://..." }
    { "type": "file", "file_id": "file_xxx" }      // Files API beta

DocumentBlockParam:
{ "type": "document", "source": DocumentSource, "title"?: "...", "citations"?: CitationsConfig }
  DocumentSource:
    { "type": "base64", "media_type": "application/pdf", "data": "<b64>" }
    { "type": "url",   "url": "https://.../file.pdf" }
    { "type": "file",  "file_id": "file_xxx" }
    { "type": "text",  "media_type": "text/plain", "data": "..." }    // Go SDK 含此形态
```

#### 工具调用 / 工具结果

```
ToolUseBlockParam:                                              // assistant 角色
{ "type": "tool_use", "id": "toolu_xxx", "name": "...", "input": {...}, "caller"?: ToolCaller, "cache_control"?: ... }

ToolResultBlockParam:                                           // user 角色
{ "type": "tool_result", "tool_use_id": "toolu_xxx",
  "content": string | ContentBlockParam[],                      // 双形态
  "is_error": false, "cache_control"?: ... }
  // content 数组里允许再嵌套 text / image / tool_reference

ToolReferenceBlockParam:                                        // tool search 动态加载标记
{ "type": "tool_reference", "tool_name": "...", "cache_control"?: ... }
```

#### 思考块

```
ThinkingBlockParam:                                              // 4.6+ 回传用
{ "type": "thinking", "thinking": "...", "signature": "..." }    // 4.7+ 默认 thinking 文本为空
                                                                  // signature 必须原样保留

RedactedThinkingBlockParam:                                      // 加密的思考
{ "type": "redacted_thinking", "data": "..." }

MidConversationSystemBlockParam:                                 // Opus 4.8+ 中段系统消息
{ "type": "mid_conv_system", "content": [TextBlockParam, ...] }  // 必须在 user 之后;不能作为 messages[0]
```

#### 服务端工具调用 / 结果(请求侧即允许标注历史)

```
ServerToolUseBlockParam:                                         // Claude 自己调用 web_search 等
{ "type": "server_tool_use",
  "id": "srvtoolu_xxx",
  "name": "web_search" | "web_fetch" | "code_execution"
       | "bash_code_execution" | "text_editor_code_execution"
       | "tool_search_tool_regex" | "tool_search_tool_bm25",
  "input": {...}, "caller"?: ToolCaller }

WebSearchToolResultBlockParam:                                   // web_search 的结果回放
{ "type": "web_search_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": WebSearchResult[] | WebSearchToolError }              // union:成功列表 / 错误对象
  WebSearchResult: {
    "type": "web_search_result",
    "title": "...", "url": "...",
    "encrypted_content": "...",                                  // 用于客户端继续对话
    "page_age": "..."                                             // 可选
  }
  WebSearchToolError: {
    "type": "web_search_tool_result_error",
    "error_code": "invalid_tool_input" | "unavailable"
               | "max_uses_exceeded" | "too_many_requests"
               | "query_too_long" | "request_too_large"
  }
```

#### 文件 / 压缩

```
ContainerUploadBlockParam:                                       // 用户上传文件到 code execution 容器
{ "type": "container_upload", "file_id": "file_xxx", "cache_control"?: ... }

CompactionBlockParam:                                            // 长上下文压缩回放
{ "type": "compaction", "content": string | null }               // null = 压缩失败
```

#### ToolCaller(同时用于 ToolUse / ServerToolUse)

```
ToolCaller:
{ "type": "direct" }                                              // 用户直接请求
{ "type": "code_execution_20250825", "tool_id": "srvtoolu_xxx" }  // PTC:从 code_execution 调用
{ "type": "code_execution_20260120", "tool_id": "srvtoolu_xxx" }
```

#### CitationsConfig

```
CitationsConfig:
{ "enabled": true }                                                // 文档引用开关
```

---

## 二、非流式 Response

### 2.1 顶层结构

```json
{
  "id": "msg_01ABC...",
  "type": "message",                            // 固定
  "role": "assistant",                          // 固定
  "model": "claude-opus-4-8",
  "content": [ContentBlock, ...],               // 见 2.2
  "stop_reason": "end_turn",
  "stop_sequence": "",                          // 命中的自定义 stop_sequence,无则空串
  "container": {                                // 代码执行容器信息(若使用)
    "id": "container_xxx",
    "expires_at": "2026-07-01T00:00:00Z"
  },
  "usage": Usage,                               // 见 2.3
  "diagnostics": CacheDiagnostics,               // beta:缓存诊断
  "context_management": {                       // beta:上下文编辑元数据
    "applied_edits": [
      { "type": "clear_tool_uses_20250919", "cleared_input_tokens": 1024,
        "cleared_tool_uses": 5 },
      { "type": "clear_thinking_20251015", "cleared_input_tokens": 512,
        "cleared_thinking_blocks": 2 }
    ]
  }
}
```

### 2.2 ContentBlock 响应侧(13 种)

**文本 / 思考**

```
TextBlock:
{ "type": "text", "text": "...", "citations"?: [Citation, ...] }

ThinkingBlock:                                                          // 4.7/4.8/Fable 5 默认 thinking 文本为空
{ "type": "thinking", "thinking": "", "signature": "..." }              // signature 必须原样回传

RedactedThinkingBlock:
{ "type": "redacted_thinking", "data": "..." }
```

**工具调用**(响应侧 Caller 必填)

```
ToolUseBlock:
{ "type": "tool_use", "id": "toolu_xxx", "name": "...", "input": {...}, "caller": ToolCaller }
  // input 是已解析 JSON 对象,而非 RawMessage
```

**服务端工具调用**

```
ServerToolUseBlock:
{ "type": "server_tool_use", "id": "srvtoolu_xxx", "name": <7 种 enum>, "input": {...}, "caller": ToolCaller }

WebSearchToolResultBlock:
{ "type": "web_search_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": [WebSearchResult, ...] | WebSearchToolError }              // 成功列表 或 错误对象

WebFetchBlock:                                                           // 嵌套 DocumentBlock
{ "type": "web_fetch_result", "url": "https://...",
  "retrieved_at": "2026-06-29T10:00:00Z",                                // ISO 8601
  "content": DocumentBlock }                                              // 复用 Document 块

WebFetchToolResultBlock:
{ "type": "web_fetch_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": WebFetchBlock | WebFetchToolResultError }

CodeExecutionToolResultBlock:                                            // code_execution 顶层(beta 旧版)
{ "type": "code_execution_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": CodeExecutionResult | Error }

BashCodeExecutionToolResultBlock:                                        // code_execution_20260120
{ "type": "bash_code_execution_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": BashCodeExecutionResult | BashCodeExecutionToolResultError }

TextEditorCodeExecutionToolResultBlock:                                  // text_editor 工具结果
{ "type": "text_editor_code_execution_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": ViewResult | CreateResult | StrReplaceResult | Error }

ToolSearchToolResultBlock:                                               // tool search 动态加载结果
{ "type": "tool_search_tool_result", "tool_use_id": "srvtoolu_xxx",
  "content": ToolSearchResult[] | Error }
```

**代码执行细分结果**

```
BashCodeExecutionResult:
{ "type": "bash_code_execution_result",
  "stdout": "...", "stderr": "...", "return_code": 0,
  "content": [BashCodeExecutionOutputBlock, ...] }                       // 引用生成的 file_id

BashCodeExecutionOutputBlock:
{ "type": "bash_code_execution_output", "file_id": "file_xxx" }

TextEditorCodeExecutionViewResultBlock:
{ "type": "text_editor_code_execution_view_result", "content": "...",
  "start_line": 0, "num_lines": 50 }

TextEditorCodeExecutionCreateResultBlock:
{ "type": "text_editor_code_execution_create_result" }

TextEditorCodeExecutionStrReplaceResultBlock:
{ "type": "text_editor_code_execution_str_replace_result",
  "old_start": 10, "old_end": 12, "new_start": 10, "new_end": 11 }

TextEditorCodeExecutionToolResultError:
{ "type": "text_editor_code_execution_tool_result_error",
  "error_code": "...", "error_message": "..." }
```

**文件 / 压缩**

```
ContainerUploadBlock:                                                   // Files API 上传后由客户端发送给 Claude 的容器文件引用
{ "type": "container_upload", "file_id": "file_xxx" }

CompactionBlock:                                                        // 长上下文压缩
{ "type": "compaction", "content": string | null }                       // null 表示压缩失败
```

**回退标记(Fable 5)**

```
FallbackBlock:                                                           // server-side fallback 切换点
{ "type": "fallback", "from": { "model": "claude-fable-5" }, "to": { "model": "claude-opus-4-8" } }
```

### 2.3 Usage(完整结构)

```json
{
  "input_tokens": 1024,                            // 未缓存输入 token
  "output_tokens": 256,                            // 输出 token(含 thinking)
  "cache_creation_input_tokens": 8192,             // 缓存写入 token(写入费率 ~1.25x)
  "cache_read_input_tokens": 4096,                 // 缓存读取 token(读取费率 ~0.1x)
  "inference_geo": "us",                           // 实际推理所在地

  "cache_creation": {                              // TTL 分桶
    "ephemeral_5m_input_tokens": 8192,             // 5m TTL(费率 1.25x)
    "ephemeral_1h_input_tokens": 0                 // 1h TTL(费率 2x)
  },
  "output_tokens_details": {
    "thinking_tokens": 180                         // 思考文本 token(含分隔符,原始 token)
  },                                              //   output_tokens - thinking_tokens ≈ 可见输出
  "server_tool_use": {
    "web_search_requests": 1,
    "web_fetch_requests": 0
  },
  "service_tier": "standard"                       // "standard" | "priority" | "batch"

  // Fable 5 server-side fallback 专属
  "iterations": [
    { "type": "message",           ...Usage 子项 },  // 被拒的尝试
    { "type": "fallback_message", ...Usage 子项 }   // 实际服务的尝试(同顶层 usage)
  ]
}
```

**流式专属**:`message_delta` 事件的 Usage **只含** `output_tokens`(更新) + `input_tokens` + `cache_creation_input_tokens` + `cache_read_input_tokens` + `web_search_requests` + `inference_geo`,**不包含** `cache_creation` 分桶 / `output_tokens_details` / `server_tool_use` / `service_tier` / `iterations`。要拿完整 Usage 必须等待非流式响应或在客户端累积。

### 2.4 StopReason(响应结束原因)

```
"end_turn"                          // 自然停止
"max_tokens"                        // 达到 max_tokens 上限
"stop_sequence"                     // 命中自定义 stop_sequence(具体值在 stop_sequence 字段)
"tool_use"                          // 模型调用了工具,需要执行
"pause_turn"                        // 长任务暂停,可原样回传继续
"refusal"                           // 安全策略拒绝,见 stop_details
"model_context_window_exceeded"     // 上下文窗口耗尽
```

### 2.5 StopDetails(仅 `stop_reason == "refusal"` 时存在)

```
{
  "type": "refusal",
  "category": "cyber" | "bio" | "reasoning_extraction" | "frontier_llm" | null,
  "explanation": "..." | null,
  "fallback_credit_token": "fc_xxx" | null,        // Fable 5:用于后续 fallback credit 重试
  "fallback_has_prefill_claim": false,              // Fable 5:是否能用 prefill 重试
  "recommended_model": "claude-opus-4-8" | null    // Fable 5:rate-limited 时推荐重试的模型
}
```

**关键行为**:HTTP 200 + `stop_reason: "refusal"`。
- 预输出(pre-output):`content` 为空,**不计费**。
- 流中(mid-stream):已流出的部分**按正常费率计费**,应丢弃。

### 2.6 Citation(5 种位置类型)

```
CitationCharLocation:                                                    // 纯文本字符位置
{ "type": "char_location", "cited_text": "...",
  "document_index": 0, "document_title": "...",
  "file_id": "file_xxx",
  "start_char_index": 100, "end_char_index": 150 }

CitationPageLocation:                                                    // PDF 页码
{ "type": "page_location", "cited_text": "...",
  "document_index": 0, "document_title": "...",
  "start_page_number": 5, "end_page_number": 6 }

CitationContentBlockLocation:                                            // 自定义内容块
{ "type": "content_block_location", "cited_text": "...",
  "start_block_index": 0, "end_block_index": 1 }

CitationsWebSearchResultLocation:                                       // web_search 来源
{ "type": "web_search_result_location", "cited_text": "...",
  "encrypted_index": "...", "title": "...", "url": "..." }

CitationsSearchResultLocation:                                          // tool search 来源
{ "type": "search_result_location", "cited_text": "...",
  "search_result_index": 0, "source": "...", "title": "..." }
```

### 2.7 错误响应

```json
{
  "type": "error",
  "error": {
    "type": "api_error" | "authentication_error" | "permission_error"
         | "not_found_error" | "rate_limit_error" | "overloaded_error"
         | "invalid_request_error" | "billing_error",
    "message": "Detailed error description."
  },
  "request_id": "req_011CSHoEeqs5C35K2UUqR7Fy"        // 排错时附在工单里
}
```

---

## 三、流式 Response(SSE 事件流)

`Content-Type: text/event-stream`,每个事件以 `data: {JSON}\n\n` 形式发出。事件类型 8 种。

### 3.1 事件时序

```
message_start
   ↓
content_block_start
content_block_delta  ━━━ 多次,直到该块结束
content_block_stop
   ↓ (对每个 content_block 重复 2→3→4)
   ↓
message_delta                                        // 仅一次
   ↓
message_stop
   ↓
[连接关闭]
// 任何时点可能穿插 event: ping(Keep-Alive)
```

### 3.2 事件完整结构

#### `message_start`

```
{
  "type": "message_start",
  "message": {
    "id": "msg_01ABC...",
    "type": "message",
    "role": "assistant",
    "model": "claude-opus-4-8",
    "content": [],
    "stop_reason": null,
    "stop_details": null,
    "container": null,
    "usage": {
      "input_tokens": 1024,
      "output_tokens": 1,                            // ← 固定为 1,作为输入统计哨兵
      "cache_creation_input_tokens": 8192,
      "cache_read_input_tokens": 4096,
      "inference_geo": "us"
    },
    "diagnostics": null,
    "context_management": null
  }
}
```

#### `content_block_start`

```
{
  "type": "content_block_start",
  "index": 0,                                       // 内容块在 content 数组中的下标
  "content_block": {                                // 起始状态,各类型不同
    "type": "text",
    "text": ""                                       // text 块起始为空

    // type:"thinking"
    // { "type":"thinking", "thinking":"", "signature":"" }

    // type:"tool_use"
    // { "type":"tool_use", "id":"toolu_xxx", "name":"...", "input":{}, "caller":{...} }

    // type:"server_tool_use"
    // { "type":"server_tool_use", "id":"srvtoolu_xxx", "name":"web_search", "input":{}, "caller":{...} }

    // type:"web_search_tool_result" / "web_fetch_tool_result"
    // { "type":"...", "tool_use_id":"srvtoolu_xxx", "content":[...] | {...} }

    // type:"bash_code_execution_tool_result" / "text_editor_code_execution_tool_result"
    // { "type":"...", "tool_use_id":"srvtoolu_xxx", "content":{...} }

    // type:"tool_search_tool_result"
    // { "type":"tool_search_tool_result", "tool_use_id":"srvtoolu_xxx", "content":[...] | {...} }

    // type:"container_upload"
    // { "type":"container_upload", "file_id":"file_xxx" }

    // type:"compaction"
    // { "type":"compaction", "content": null }

    // type:"redacted_thinking"
    // { "type":"redacted_thinking", "data":"" }
  }
}
```

#### `content_block_delta`(5 种类型)

`delta` 字段按 `type` 路由到 5 种结构:

```
// 1. 文本块增量
TextDelta:
{ "type": "text_delta", "text": "片段字符串" }

// 2. 工具调用 input(流式 JSON,需自行拼接 + json.Unmarshal)
InputJSONDelta:
{ "type": "input_json_delta", "partial_json": "{\"loc" }    // 增量 JSON 片段

// 3. 思考块增量
ThinkingDelta:
{ "type": "thinking_delta", "thinking": "思考片段" }          // 4.7+ 默认 text=空字符串

// 4. 思考块签名增量(闭合标记)
SignatureDelta:
{ "type": "signature_delta", "signature": "..." }           // thinking 块的最后一段 delta

// 5. 引用增量
CitationsDelta:
{ "type": "citations_delta", "citation": Citation }            // Citation 见 2.6
```

完整事件形态:

```
{
  "type": "content_block_delta",
  "index": 0,
  "delta": <上述 5 种之一>
}
```

#### `content_block_stop`

```
{
  "type": "content_block_stop",
  "index": 0
}
```

#### `message_delta`(仅一次)

```
{
  "type": "message_delta",
  "delta": {
    "stop_reason": "end_turn" | ...,                  // 首次出现 stop_reason
    "stop_details": RefusalStopDetails | null,        // 仅 refusal 时存在
    "stop_sequence": "" | "...",                      // 命中的自定义 stop_sequence
    "container": { "id":"container_xxx", "expires_at":"..." } | null
  },
  "usage": {
    "output_tokens": 256,                             // 实时累计更新
    "input_tokens": 1024,
    "cache_creation_input_tokens": 8192,
    "cache_read_input_tokens": 4096,
    "web_search_requests": 1,
    "inference_geo": "us"
  },                                                  // 流式 usage 字段比非流式少 4 个
  "context_management": {                             // beta:上下文编辑元数据
    "applied_edits": [...]
  }
}
```

#### `message_stop`

```
{
  "type": "message_stop"
}
```

#### `ping`(可选 Keep-Alive)

```
{
  "type": "ping"
}
```

#### `error`(异常终止)

```
{
  "type": "error",
  "error": {
    "type": "...",
    "message": "..."
  }
}
```

### 3.3 同协议 passthrough（raw bytes 透传）

> 起于：commit `a952246`（Anthropic + Gemini 对称 passthrough 重构 / T2）。
> 单元测试：`pkg/protocol/anthropic/stream_event_test.go` `TestParseSSERawFrame_PreservesRawBytes` + `frames/frame_test.go` `TestMarshalSSEFrameRaw_Verbatim`。
> E2E 回归：`test/e2e/l2/streaming_test.go` `TestL2_Streaming_AnthropicPassthrough`。

上游是 Anthropic SSE 的同协议请求路径下，gateway 不做 typed struct 往返，而是直接把整帧 `event: X\ndata: Y\n\n` 原始字节透传给下游，避免两层 loss：

| Loss 点 | 来源 | 影响 |
|---------|------|------|
| LOSS #1（typed unmarshal） | `json.Unmarshal(data, &typedStruct)` | typed struct 未声明的未知字段（如 4.6+ `stop_details.refusal_category` / 未来扩展字段）整字段消失 |
| LOSS #2（typed marshal） | `json.Marshal(typedStruct)` 配合 `omitempty` | 上游零值字段被剥（如 `cache_creation_input_tokens:0` / `cache_read_input_tokens:0` / `service_tier:""` / `stop_reason:null` / `stop_sequence:null`） |

透传路径关键载体：

| 名称 | 路径 | 作用 |
|------|------|------|
| `AnthropicRawStreamFrame` | `pkg/protocol/anthropic/stream_event.go` | sealed carrier（带私有 `isAnthropicRawStreamFrame()` 标记，**不**实现 `AnthropicStreamEvent` 接口，强制 type switch 显式分流） |
| `ParseSSERawFrame(line, rawFrame)` | `pkg/protocol/anthropic/stream_event.go` | 抽 type + 装载整帧 raw bytes |
| `frames.MarshalSSEFrameRaw(raw)` | `pkg/protocol/anthropic/frames/frame.go` | verbatim 透传字节，**不**走 `json.Marshal` |
| DoStream 累加器 | `internal/app/gateway/provider/base.go` | `bufio.Scanner` 行扫描；`event:` 行缓存 → `data:` 行触发时拼出 `event: X\ndata: Y\n\n` 整帧 |

模型重写仍生效：byte-level 替换 `"model":"<from>"` → `"model":"<to>"`（仅 Anthropic model 字段，与 OpenAI 对称）。

typed 路径仍保留：8 个 `AnthropicStreamEvent` 变体继续服务 synthetic 帧（`BuildMessageStartFrame` 等内部构造），由 `MarshalSSEFrame` 输出。

---

## 四、流式 vs 非流式 核心差异

| 维度 | 非流式 | 流式 |
|------|--------|------|
| 请求开关 | `stream` 字段省略或 `false` | `stream: true` |
| 响应 Content-Type | `application/json` | `text/event-stream` |
| 返回形态 | 单一 JSON 对象 | SSE 事件流(8 种事件) |
| 拿到 `stop_reason` 时机 | 响应末尾 | `message_delta` 事件 |
| `output_tokens` 累计 | 完整统计 | `message_start` 时固定为 **1**(哨兵值),`message_delta` 时累计更新 |
| Usage 完整性 | 全字段(含 cache_creation 分桶 / output_tokens_details / server_tool_use / service_tier / iterations) | 仅 6 字段(缺失分桶与细分) |
| 适用场景 | 短输出(≤16K tokens) | 长输出(≥16K tokens) / 实时 UI |
| 4.6+ 推荐 | 默认走流式避免 SDK HTTP 超时 | 唯一推荐路径 |
| 工具调用 input | 完整 JSON 对象 | 增量 `input_json_delta` JSON 片段,客户端拼接 + parse |
| thinking 块 | 完整文本(4.7+ 默认空) | 同上,delta 是 `thinking_delta` + 闭合 `signature_delta` |
| 累积复杂度 | 一次性解析 | 需客户端 accumulator(ClaudeCode 用 `stream.finalMessage()` / Go SDK `stream.Accumulate`) |

---

## 五、回传复用注意事项

流式累积得到的 finalMessage 与非流式 Response 在语义上完全等价——可直接送回 `messages[]`(保留 `tool_use`、`thinking`、`signature`、`compaction`、`server_tool_use`、`web_search_tool_result` 等所有块原样)作为下一轮请求。

**关键约束**:
- 4.7+ 的 `thinking` 块即便 `thinking` 文本为空也必须**原样回传**,修改会被 400。
- Fable 5 的 thinking 块在不同模型间传递时被**静默丢弃**(unbilled);同模型内必须原样回传。
- 缓存块(`cache_control`)重新进入 messages 会刷新缓存键,影响 prompt cache 命中。
- `signature` 字段是签名认证,任何修改都会导致 400。
- `encrypted_content`(web_search_result)用于跨轮次复用搜索结果;保留原值即可,服务端解密。

---

## 六、字段版本兼容性

| 字段 | 适用版本 | 备注 |
|------|----------|------|
| `output_config` | 4.6+ GA | 替换旧 `output_format` |
| `output_config.effort` | 4.5+ | Opus 4.6+ 支持 `max`;Opus 4.7+ 支持 `xhigh` |
| `output_config.task_budget` | 4.7+ beta | header `task-budgets-2026-03-13` |
| `thinking: adaptive` | 4.6+ | 替换 `budget_tokens` |
| `thinking.display` | 4.7+ / Fable 5 | 4.6 不识别 |
| `compaction` 块 | 4.7+ beta | header `compact-2026-01-12` |
| `mid_conv_system` 块 | Opus 4.8+ | 无 beta header;必须 user 之后;不能 messages[0] |
| `context_management` | 4.6+ beta | header `context-management-2025-06-27` |
| `mcp_servers` | 4.6+ beta | header `mcp-client-2025-11-20` |
| `temperature` / `top_p` / `top_k` | ⚠️ Opus 4.7+/Fable 5 禁用 | 传则 400 |
| `budget_tokens` | ⚠️ Opus 4.7+/Fable 5 禁用 | 传则 400;4.6 过渡可用 |
| `thinking: enabled` | ⚠️ Opus 4.7+/Fable 5 禁用 | 同上 |
| `thinking: disabled` | ⚠️ Fable 5 禁用(传则 400) | Opus 4.8/4.7 接受 |
| 末位 assistant prefill | ⚠️ Opus 4.6+/Fable 5 禁用 | 传则 400 |
| `speed: "fast"` | Opus 4.8 | 4.6/4.7 fast mode 已 deprecated |
| `web_search_20260209` 动态过滤 | Opus 4.6+/Sonnet 4.6+ | 旧模型用 `web_search_20250305` |
| `web_fetch_20260209` 动态过滤 | Opus 4.6+/Sonnet 4.6+ | 旧模型用 `web_fetch_20250910` |
| `code_execution_20260521` | GA | 旧版 `code_execution_20260120` 仍可用 |
| `text_editor_20250728` | GA | name 必为 `str_replace_based_edit_tool`(旧 `str_replace_editor` 已弃用) |
| `memory_20250818` | GA | |
| `tool_use.caller` | 4.6+ | 非 beta header |
| `fallbacks` / `fallback_credit_token` | Fable 5 | beta header `server-side-fallback-2026-06-01` |

---

## 七、鉴权头(HTTP Header)

```
Content-Type:        application/json             # 必填
x-api-key:           sk-ant-xxx                   # 鉴权(必填)
Authorization:       Bearer sk-ant-oat01-...      # OAuth bearer 替代 x-api-key
anthropic-version:   2023-06-01                   # API 版本
anthropic-beta:      beta-id-1,beta-id-2          # 多 beta 用换行或逗号分隔
anthropic-user-profile-id: up_xxx                  # user-profiles beta
```

---

## 八、参考

- Anthropic 官方 Messages API: <https://docs.anthropic.com/en/api/messages>
- 官方 Go SDK 字段基线: `anthropic-sdk-go/message.go`(10105 行)
- ClaudeCode TypeScript: 流式累积模式与思考块回传策略
- 官方 SDK Repositories:
  - Python: <https://github.com/anthropics/anthropic-sdk-python>
  - TypeScript: <https://github.com/anthropics/anthropic-sdk-typescript>
  - Go: <https://github.com/anthropics/anthropic-sdk-go>
  - Java: <https://github.com/anthropics/anthropic-sdk-java>
  - C#: <https://github.com/anthropics/anthropic-sdk-csharp>
  - PHP: <https://github.com/anthropics/anthropic-sdk-php>
  - Ruby: <https://github.com/anthropics/anthropic-sdk-ruby>