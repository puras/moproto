# OpenAI API — Request/Response/Stream 完整 Wire-Format 结构

> 参照：
> - OpenAI 官方 API 参考（<https://platform.openai.com/docs/api-reference>）
> - 官方 Go SDK `openai-go`（`chatcompletion.go` 4253 行 + `responses/response.go` ~22818 行）
>
> 覆盖范围：ChatCompletion API + Responses API + Audio（Speech / Transcription）+ Image（Generate / Edit）+ Embedding + Moderation。
>
> 本文档为纯结构梳理，描述 wire-format JSON 形状与字段约束，**不涉及任何代码类型**。

---

## 一、ChatCompletion API（`/v1/chat/completions`）

### 1.1 请求结构（流式 + 非流式共用）

```
POST /v1/chat/completions
Content-Type: application/json
Authorization: Bearer sk-xxx

{
  // ──────────────── 必填字段 ────────────────
  "model": "gpt-4o",                                  // string
  "messages": [Message, ...],                         // 见 1.2（6 种 role）

  // ──────────────── 流式开关 ────────────────
  "stream": true,                                     // bool（指针，默认 false）
  "stream_options": { "include_usage": true },        // 仅在 stream=true 时有效

  // ──────────────── 采样参数 ────────────────
  "temperature": 0.7,                                 // float (0-2)
  "top_p": 0.9,                                       // float
  "frequency_penalty": 0,                             // float (-2, 2)
  "presence_penalty": 0,                              // float (-2, 2)
  "seed": 42,                                         // int64

  // ──────────────── 输出限制 ────────────────
  "max_tokens": 4096,                                 // int64（GPT-4 系列）
  "max_completion_tokens": 16384,                     // int64（o 系列推荐）
  "n": 1,                                             // int64（生成数）

  // ──────────────── 停止条件 ────────────────
  "stop": "END" | ["END", "STOP"],                    // string | string[]

  // ──────────────── 工具调用 ────────────────
  "tools": [Tool, ...],                               // 见 1.3（6 种 union）
  "tool_choice": ToolChoice,                          // 见 1.4（5 种策略）
  "parallel_tool_calls": true,                        // bool

  // ──────────────── 响应格式 ────────────────
  "response_format":                                  // 见 1.5
    { "type": "text" }
    | { "type": "json_object" }
    | { "type": "json_schema", "json_schema": { ... } },

  // ──────────────── Logprobs ────────────────
  "logprobs": false,                                  // bool
  "top_logprobs": 0,                                  // int64 (0-20)

  // ──────────────── Logit bias ────────────────
  "logit_bias": { "50256": -100 },                    // map[string]int64

  // ──────────────── 用户标识 / 元数据 ────────────────
  "user": "user_xxx",                                 // string（滥用检测）
  "metadata": { "key": "value" },                     // map[string]string（最多 16 对）

  // ──────────────── 存储 ────────────────
  "store": false,                                     // bool（用于评估/微调）

  // ──────────────── 服务等级 ────────────────
  "service_tier": "auto" | "default" | "priority" | "flex",

  // ──────────────── 音频输出（4.6+ GA） ────────────────
  "audio": { "voice": "alloy" | { "id": "custom_xxx" }, "format": "wav" | "mp3" | "flac" | "opus" | "pcm16" },
  "modalities": ["text", "audio"],                    // string[]

  // ──────────────── 预测内容（4.6+） ────────────────
  "prediction": { "type": "content", "content": [{ "type": "text", "text": "..." }] },

  // ──────────────── 推理（o 系列） ────────────────
  "reasoning_effort": "minimal" | "low" | "medium" | "high",

  // ──────────────── 详细度 ────────────────
  "verbosity": "low" | "medium" | "high",

  // ──────────────── Web 搜索 ────────────────
  "web_search_options": {
    "search_context_size": "low" | "medium" | "high",
    "user_location": { "type": "approximate", "city": "...", "country": "...", "region": "...", "timezone": "..." }
  },

  // ──────────────── 安全（4.6+） ────────────────
  "safety_identifier": "safe_xxx",                    // string

  // ──────────────── 提示缓存 ────────────────
  "prompt_cache_key": "cache_key_xxx"                 // string
}
```

### 1.2 Message union（6 种 role）

```
UserMessage:
{ "role": "user", "content": String | ContentParts, "name"?: "..." }

AssistantMessage:
{ "role": "assistant", "content"?: "...", "name"?: "...",
  "tool_calls"?: [ToolCall, ...],                     // Function | Custom
  "audio"?: { "id": "...", "data": "...", "expires_at": ..., "transcript": "..." },
  "refusal"?: "...",
  "function_call"?: { "name": "...", "arguments": "..." } }    // 旧版，已废弃

SystemMessage:
{ "role": "system", "content": String | ContentParts, "name"?: "..." }

DeveloperMessage:                                    // 4.6+ 新增
{ "role": "developer", "content": String | ContentParts, "name"?: "..." }

ToolMessage:                                         // 工具结果
{ "role": "tool", "content": String | ContentParts, "tool_call_id": "call_xxx" }

FunctionMessage:                                     // 旧版，已废弃
{ "role": "function", "content": "...", "name": "..." }
```

**UserMessage.Content 双形态**：

```
String:        "Hello"                               // 纯文本
ContentParts:  [ContentPart, ...]                    // 多模态
```

### 1.3 ContentPart union（7 种）

```
TextContentPart:        { "type": "text", "text": "..." }
ImageContentPart:       { "type": "image_url", "image_url": { "url": "https://...", "detail": "auto" | "low" | "high" } }
InputAudioContentPart:  { "type": "input_audio", "input_audio": { "data": "<b64>", "format": "wav" | "mp3" } }
FileContentPart:        { "type": "file", "file": { "file_data": "<b64>", "filename": "..." } }
VideoContentPart:       { "type": "video_url", "video_url": { "url": "https://...", "detail": "..." } }
RefusalContentPart:     { "type": "refusal", "refusal": "..." }                            // 助手消息
AudioContentPart:       { "type": "audio", "audio": { "id": "audio_xxx", "data": "<b64>" } }  // 助手消息输出
```

### 1.4 Tool union（6 种）

| 工具 | type | 客户端/服务端 |
|------|------|--------------|
| **Function** | `function` | 客户端 |
| **Custom** | `custom`（4.6+） | 客户端 |
| **WebSearch** | `web_search`（hosted） | 服务端 |
| **CodeInterpreter** | `code_interpreter`（hosted） | 服务端 |
| **FileSearch** | `file_search`（hosted） | 服务端 |
| **ComputerUse** | `computer_use_preview`（hosted） | 服务端 |

```
FunctionTool:
{ "type": "function",
  "function": { "name": "get_weather", "description": "...",
                "parameters": { /* JSON Schema */ }, "strict"?: true } }

CustomTool:                                           // 4.6+
{ "type": "custom",
  "custom": { "name": "...", "description": "...",
              "format"?: { "type": "text" } |
                         { "type": "grammar", "syntax": "lark" | "regex", "definition": "..." } } }

WebSearchTool:
{ "type": "web_search",
  "web_search_options"?: { "search_context_size": "...", "user_location": {...} } }

CodeInterpreterTool:
{ "type": "code_interpreter", "container"?: "cntr_xxx" }

FileSearchTool:
{ "type": "file_search",
  "vector_store_ids"?: [...],
  "max_num_results"?: 10,
  "ranking_options"?: { "ranker": "...", "score_threshold": 0.5 },
  "filters"?: { /* JSON */ } }

ComputerUseTool:
{ "type": "computer_use_preview",
  "display_width": 1024, "display_height": 768, "environment": "mac" | "windows" | "linux" | ... }
```

### 1.5 ToolChoice union（5 种策略）

```
{ "type": "auto" }                                                  // 默认
{ "type": "required" }                                              // 必须调用某个工具
{ "type": "none" }                                                  // 禁止调用
{ "type": "allowed_tools", "allowed_tools": [...], "mode": "auto" | "required" }   // 4.6+
{ "type": "function", "function": { "name": "web_search" } }       // 指定 hosted 工具
```

### 1.6 ResponseFormat union（3 种）

```
{ "type": "text" }
{ "type": "json_object" }
{ "type": "json_schema",
  "json_schema": { "name": "schema", "description": "...", "schema": { /* JSON Schema */ }, "strict": true } }
```

### 1.7 响应结构

```
ChatCompletionResponse:
{
  "id": "chatcmpl-xxx",
  "object": "chat.completion",                       // 固定
  "created": 1718000000,                             // Unix 时间戳
  "model": "gpt-4o-2024-08-06",
  "choices": [Choice, ...],
  "usage"?: Usage,                                   // 仅 stream_options.include_usage=true 或非流式
  "service_tier"?: "default" | "priority" | "flex",
  "system_fingerprint"?: "fp_xxx"                    // 系统指纹（路由决策追踪）
}

Choice:
{
  "index": 0,
  "message": {
    "role": "assistant",
    "content": "..."  | "",                          // 工具调用时可能为空
    "refusal"?: null | "...",
    "tool_calls"?: [ToolCall, ...],
    "audio"?: { "id": "audio_xxx", "data": "<b64>", "expires_at": ..., "transcript": "..." },
    "annotations"?: [URLCitation, ...],
    "function_call"?: {...}                          // 旧版
  },
  "finish_reason": "stop" | "length" | "tool_calls" | "content_filter" | "function_call",
  "logprobs"?: { "content": [...], "refusal": [...] }
}

Usage:
{
  "prompt_tokens": 100,
  "completion_tokens": 200,
  "total_tokens": 300,
  "prompt_tokens_details"?: { "cached_tokens": 80, "audio_tokens": 0 },
  "completion_tokens_details"?: { "reasoning_tokens": 50, "audio_tokens": 0,
                                    "accepted_prediction_tokens": 10, "rejected_prediction_tokens": 5 }
}
```

### 1.8 流式 chunk（SSE）

每个 chunk 是一个 ChatCompletion 对象，`choices` 数组中含 `Delta` 字段：

```
Chunk:
{
  "id": "chatcmpl-xxx",
  "object": "chat.completion.chunk",
  "created": 1718000000,
  "model": "gpt-4o",
  "choices": [ChunkChoice, ...],
  "usage"?: Usage,                                   // 仅最后一个 chunk（include_usage=true）
  "service_tier"?: "...",
  "system_fingerprint"?: "..."
}

ChunkChoice:
{
  "index": 0,
  "id"?: "...",
  "delta": {
    "role"?: "assistant",                            // 第一个 chunk 设置
    "content"?: "Hello",                             // 文本增量
    "refusal"?: "...",
    "tool_calls"?: [{                                // 工具调用增量
      "index": 0,
      "id"?: "call_xxx",                             // 第一个 chunk 设置
      "type"?: "function",                           // 第一个 chunk 设置
      "function"?: { "name"?: "...", "arguments"?: "..." }   // 增量 JSON 字符串片段
    }],
    "audio"?: { "id"?: "...", "data"?: "<b64>", "transcript"?: "...", "expires_at"?: ... },
    "function_call"?: {...}                          // 旧版
  },
  "finish_reason"?: "stop" | ...,                    // 最后一个 chunk 设置
  "logprobs"?: {...}
}
```

---

## 二、Responses API（`/v1/responses`）

> Responses API 是 4.6+ 主推的"统一协议"接口，整合 ChatCompletion + Tools +
> ImageGeneration + FileSearch + MCP + CodeInterpreter + Reasoning。

### 2.1 请求结构

```
POST /v1/responses
Content-Type: application/json
Authorization: Bearer sk-xxx

{
  // ──────────────── 必填字段 ────────────────
  "model": "gpt-4o",
  "input": String | [InputItem, ...],                // 见 2.2

  // ──────────────── 系统指令 ────────────────
  "instructions": String | [InputItem, ...],

  // ──────────────── 工具 ────────────────
  "tools": [Tool, ...],                              // 见 2.3（10 种 union）
  "tool_choice": "auto" | "required" | "none",
  "parallel_tool_calls"?: true,

  // ──────────────── 流式 ────────────────
  "stream"?: true,
  // 流模式返回 OpenAIResponseStreamEventUnion（45 种事件变体）

  // ──────────────── 输出控制 ────────────────
  "include"?: ["file_search_call.results", "web_search_call.results", ...],   // 需要回放的字段
  "store"?: false,
  "truncation"?: "auto" | "disabled",
  "max_output_tokens"?: 4096,

  // ──────────────── 推理（o 系列） ────────────────
  "reasoning"?: { "effort": "minimal" | "low" | "medium" | "high",
                  "summary"?: "auto" | "concise" | "detailed" },

  // ──────────────── 元数据 ────────────────
  "metadata"?: { "key": "value" },

  // ──────────────── 文本输出配置 ────────────────
  "text"?: { "format"?: { "type": "text" } | { "type": "json_schema", ... },
             "verbosity"?: "low" | "medium" | "high" },

  // ──────────────── 后台执行 ────────────────
  "background"?: false,

  // ──────────────── 安全标识 ────────────────
  "safety_identifier"?: "safe_xxx",

  // ──────────────── 提示缓存 ────────────────
  "prompt_cache_key"?: "cache_key_xxx",

  // ──────────────── 多轮会话 ────────────────
  "conversation"?: { "id": "conv_xxx" }
}
```

### 2.2 InputItem union（12 种）

```
MessageItem:                                         // user / system / assistant
{ "type": "message", "role": "user" | "system" | "assistant",
  "content": String | [ContentPart, ...],
  "status"?: "in_progress" | "completed" | "incomplete" }

FunctionToolCallItem:                                // assistant 工具调用
{ "type": "function_call", "id"?: "fc_xxx", "call_id"?: "call_xxx",
  "name": "get_weather", "arguments": "{...}",
  "status"?: "in_progress" | "completed" | "incomplete" }

FunctionToolCallOutputItem:                          // 工具结果
{ "type": "function_call_output", "call_id": "call_xxx",
  "output": "...", "status"?: "..." }

FileSearchToolCallItem:
{ "type": "file_search_call", "id"?: "...", "queries"?: [...], "status"?: "...",
  "results"?: [{ "file_id": "...", "filename": "...", "score": 0.95, "text": "..." }] }

WebSearchToolCallItem:
{ "type": "web_search_call", "id"?: "...", "status"?: "...",
  "action"?: { "type": "search" | "open_page" | "find_in_page",
               "query"?: "...", "url"?: "..." } }

ComputerToolCallItem:
{ "type": "computer_call", "id"?: "...", "call_id"?: "...",
  "action"?: { "type": "click" | "type" | "screenshot" | "double_click" | "scroll" | "keypress" | "move" | ... ,
               "x"?: ..., "y"?: ..., "button"?: "left" | "right" | "middle",
               "text"?: "...", "keys"?: [...], "scroll_x"?: ..., "scroll_y"?: ... },
  "status"?: "..." }

ReasoningItem:                                       // o 系列
{ "type": "reasoning", "id"?: "...", "summary"?: [...] }

ItemReferenceItem:                                   // 引用已有 item
{ "type": "item_reference", "id": "msg_xxx" }

ApplyPatchToolCallItem:                              // 文件编辑
{ "type": "apply_patch_call", "id"?: "...", "call_id"?: "...", "status"?: "...",
  "operation"?: { "type": "create_file" | "delete_file" | "update_file",
                  "path"?: "...", "diff"?: "...", "old_text"?: "...", "new_content"?: "..." } }

McpToolCallItem:                                     // MCP 工具调用
{ "type": "mcp_call", "id"?: "...", "name"?: "...", "server_label"?: "...",
  "arguments"?: "...", "output"?: "...",
  "error"?: { "code": 0, "message": "..." },
  "status"?: "..." }

CodeInterpreterToolCallItem:
{ "type": "code_interpreter_call", "id"?: "...", "call_id"?: "...", "code"?: "...",
  "outputs"?: [{ "type": "logs" | "image", "logs"?: "...", "url"?: "..." }],
  "status"?: "...", "container_id"?: "..." }

ImageGenerationCallItem:
{ "type": "image_generation_call", "id"?: "...", "result"?: "<b64>", "status"?: "..." }
```

### 2.3 Tool union（10 种）

| 工具 | type | 用途 |
|------|------|------|
| **Function** | `function` | 客户端函数 |
| **WebSearch** | `web_search` | 服务端搜索 |
| **CodeInterpreter** | `code_interpreter` | 服务端代码沙箱 |
| **FileSearch** | `file_search` | 服务端文件检索 |
| **ImageGeneration** | `image_generation` | 服务端图像生成 |
| **MCP** | `mcp` | 远程 MCP 工具 |
| **Computer** | `computer_use_preview` | 电脑使用 |
| **ApplyPatch** | `apply_patch` | 文件编辑 |
| **CodeExecution** | `code_execution` | 通用代码执行 |
| **Hosted** | `web_search_preview` 等 | 历史 hosted 工具兜底 |

```
FunctionTool:
{ "type": "function", "name": "get_weather", "description": "...",
  "strict"?: true, "parameters": { /* JSON Schema */ } }

WebSearchTool:
{ "type": "web_search",
  "filters"?: { "allowed_domains": [...] },
  "search_context_size"?: "low" | "medium" | "high",
  "user_location"?: { "type": "approximate", "city": "...", ... } }

CodeInterpreterTool:
{ "type": "code_interpreter",
  "container": { "type": "auto" } }

FileSearchTool:
{ "type": "file_search",
  "vector_store_ids": ["vs_xxx"],                     // 必填
  "max_num_results"?: 10,
  "ranking_options"?: { "ranker": "...", "score_threshold": 0.5 },
  "filters"?: { /* JSON */ } }

ImageGenerationTool:
{ "type": "image_generation",
  "background"?: "transparent" | "opaque" | "auto",
  "input_fidelity"?: "high" | "low",
  "model"?: "gpt-image-1",
  "moderation"?: "auto" | "low",
  "quality"?: "low" | "medium" | "high" | "auto",
  "output_format"?: "png" | "webp" | "jpeg",
  "size"?: "1024x1024" | "1024x1536" | "1536x1024" | "auto" }

McpTool:
{ "type": "mcp",
  "server_label": "...", "server_url"?: "...",
  "headers"?: { "Authorization": "..." },
  "allowed_tools"?: [...],
  "require_approval"?: "always" | "never" }

ComputerTool:
{ "type": "computer_use_preview",
  "display_width": 1024, "display_height": 768, "environment": "mac" | "windows" | "linux" }

ApplyPatchTool:
{ "type": "apply_patch" }                            // 无附加字段

CodeExecutionTool:
{ "type": "code_execution" }                         // 无附加字段

HostedTool:                                          // 兜底
{ "type": "web_search_preview" | ... }
```

### 2.4 响应结构

```
Response:
{
  "id": "resp_xxx",
  "object": "response",                              // 固定
  "created_at": 1718000000,
  "model": "gpt-4o-2024-08-06",
  "status": "completed" | "failed" | "in_progress" | "incomplete" | "queued",
  "error"?: { "code": "...", "message": "..." },
  "instructions"?: null | "...",
  "output": [OutputItem, ...],                       // 见 2.5（10 种 union）
  "incomplete_details"?: { "reason": "max_output_tokens" | "content_filter" },
  "usage": {
    "input_tokens": 100,
    "output_tokens": 200,
    "total_tokens": 300,
    "input_tokens_details": { "cached_tokens": 80 },
    "output_tokens_details": { "reasoning_tokens": 50 }
  },
  "parallel_tool_calls"?: true,
  "temperature"?: 0.7,
  "tool_choice"?: "...",
  "tools"?: [...],
  "top_p"?: 0.9,
  "max_output_tokens"?: 4096,
  "truncation"?: "auto" | "disabled",
  "reasoning"?: {...},
  "metadata"?: {...},
  "store"?: false,
  "background"?: false,
  "service_tier"?: "...",
  "safety_identifier"?: "...",
  "prompt_cache_key"?: "...",
  "text"?: {...},
  "conversation"?: { "id": "conv_xxx" }
}
```

### 2.5 OutputItem union（10 种）

```
OutputMessage:
{ "type": "message", "id": "msg_xxx", "role": "assistant", "status": "completed",
  "content": [OutputContent, ...] }
  // OutputContent: { "type": "output_text", "text": "...", "annotations"?: [...] }
  //               | { "type": "refusal", "refusal": "..." }

OutputFunctionCall:
{ "type": "function_call", "id": "fc_xxx", "call_id": "call_xxx",
  "name": "get_weather", "arguments": "{...}",
  "status"?: "in_progress" | "completed" | "incomplete" }

OutputWebSearchCall:
{ "type": "web_search_call", "id": "ws_xxx", "status": "...",
  "action"?: { "type": "search" | ... } }

OutputFileSearchCall:
{ "type": "file_search_call", "id": "fs_xxx", "status": "...",
  "queries"?: [...], "results"?: [...] }

OutputComputerCall:
{ "type": "computer_call", "id": "cc_xxx", "call_id"?: "...", "status"?: "...",
  "action"?: {...},
  "outputs"?: [{ "type": "computer_screenshot", "file_id"?: "...", "image_url"?: "..." }] }

OutputReasoning:                                     // o 系列推理
{ "type": "reasoning", "id": "rs_xxx", "status"?: "...",
  "summary"?: [{ "type": "summary_text", "text": "..." }],
  "content"?: ["..."] }

OutputMcpCall:
{ "type": "mcp_call", "id": "mcp_xxx",
  "name"?: "...", "server_label"?: "...", "arguments"?: "...", "output"?: "...",
  "error"?: {...}, "status"?: "..." }

OutputMcpListTools:
{ "type": "mcp_list_tools", "id": "mlt_xxx", "server_label": "...",
  "tools": [{ "name": "...", "description": "...", "input_schema": {...}, "annotations": {...} }],
  "error"?: {...} }

OutputCodeInterpreterCall:
{ "type": "code_interpreter_call", "id": "ci_xxx", "status": "...",
  "container_id"?: "...", "code"?: "...", "outputs"?: [...] }

OutputImageGenerationCall:
{ "type": "image_generation_call", "id": "img_xxx", "status": "...",
  "result"?: "<b64>", "revised_prompt"?: "..." }
```

### 2.6 注释（Annotation）

```
OutputMessage.content[].annotations:
{
  "type": "url_citation" | "file_citation" | "container_file_citation",
  "url_citation"?:     { "url": "...", "title": "...", "start_index": 100, "end_index": 150 },
  "file_citation"?:    { "file_id": "...", "filename": "...", "index": 0 }
}
```

---

## 三、Responses 流事件（45 种）

流式响应每个事件为 SSE 形式，每个事件有一个 `type` 字段。按类别归类：

### 3.1 状态变更（10 种）

```
response.created              { type, sequence_number, response }
response.queued               { type, sequence_number, response }
response.in_progress          { type, sequence_number, response }
response.completed            { type, sequence_number, response }
response.failed               { type, sequence_number, response }
response.incomplete           { type, sequence_number, response }
response.refusal.done         { type, sequence_number, item_id, output_index, content_index, refusal }
response.error                { type, sequence_number, error }
response.output_item.added    { type, sequence_number, output_index, item: OutputItemUnion }
response.output_item.done     { type, sequence_number, output_index, item: OutputItemUnion }
```

### 3.2 文本增量（4 种）

```
response.output_text.delta                { type, sequence_number, item_id, output_index, content_index, delta, logprobs? }
response.output_text.done                 { type, sequence_number, item_id, output_index, content_index, text }
response.content_part.added               { type, sequence_number, item_id, output_index, content_index, part }
response.content_part.done                { type, sequence_number, item_id, output_index, content_index, part }
```

### 3.3 推理增量（4 种）

```
response.reasoning_text.delta                       { type, sequence_number, item_id, output_index, content_index, delta }
response.reasoning_text.done                        { type, sequence_number, item_id, output_index, content_index, text }
response.reasoning_summary_part.added               { type, sequence_number, item_id, output_index, summary_index, part }
response.reasoning_summary_part.done                { type, sequence_number, item_id, output_index, summary_index, part }
```

### 3.4 工具调用（6 种）

```
response.function_call_arguments.delta              { type, sequence_number, item_id, output_index, delta }
response.function_call_arguments.done               { type, sequence_number, item_id, output_index, arguments }
response.file_search_call.in_progress               { type, sequence_number, output_index, item_id }
response.file_search_call.searching                 { type, sequence_number, output_index, item_id }
response.file_search_call.completed                 { type, sequence_number, output_index, item_id }
response.web_search_call.in_progress                { type, sequence_number, output_index, item_id }
response.web_search_call.searching                  { type, sequence_number, output_index, item_id }
response.web_search_call.completed                  { type, sequence_number, output_index, item_id }
```

### 3.5 MCP 事件（9 种）

```
response.mcp_call.in_progress                       { type, sequence_number, output_index, item_id }
response.mcp_call_arguments.delta                   { type, sequence_number, output_index, item_id, delta }
response.mcp_call_arguments.done                    { type, sequence_number, output_index, item_id, arguments }
response.mcp_call.completed                         { type, sequence_number, output_index, item_id }
response.mcp_call.failed                            { type, sequence_number, output_index, item_id }
response.mcp_list_tools.in_progress                 { type, sequence_number, item_id }
response.mcp_list_tools.completed                   { type, sequence_number, item_id }
response.mcp_list_tools.failed                      { type, sequence_number, item_id }
```

### 3.6 Code Interpreter（5 种）

```
response.code_interpreter_call.in_progress          { type, sequence_number, output_index, item_id }
response.code_interpreter_call.interpreting         { type, sequence_number, output_index, item_id }
response.code_interpreter_call_code.delta           { type, sequence_number, output_index, item_id, delta }
response.code_interpreter_call_code.done            { type, sequence_number, output_index, item_id, code }
response.code_interpreter_call.completed            { type, sequence_number, output_index, item_id }
```

### 3.7 Image Generation（4 种）

```
response.image_generation_call.in_progress          { type, sequence_number, output_index, item_id }
response.image_generation_call.generating           { type, sequence_number, output_index, item_id }
response.image_generation_call.partial_image        { type, sequence_number, output_index, item_id, partial_image_index, b64_json }
response.image_generation_call.completed            { type, sequence_number, output_index, item_id }
```

### 3.8 Audio（4 种）

```
response.audio.delta                                { type, sequence_number, item_id, output_index, content_index, delta }
response.audio.done                                 { type, sequence_number, item_id, output_index, content_index, audio }
response.audio_transcript.delta                     { type, sequence_number, item_id, output_index, content_index, delta }
response.audio_transcript.done                      { type, sequence_number, item_id, output_index, content_index, transcript }
```

### 3.9 事件流时序

```
response.created
  ↓
(response.queued)  ← 仅后台执行模式
  ↓
response.in_progress
  ↓
[对每个 output_item 重复]
  response.output_item.added
    ↓
  [对每个 content_part 重复]
    response.content_part.added
    [对 text 增量累积]
      response.output_text.delta  ━━━多次
      response.output_text.done
    response.content_part.done
  response.output_item.done
  ↓
response.completed  /  response.failed  /  response.incomplete
↓
// 任何时点可能穿插 response.error（异常终止）
```

### 3.10 TODO 预留扩展点（8 个未实现的次要事件）

下列 8 类次要变体未在本 PR 实现，待官方 SDK 稳定后按需补齐：

```
1. response.reasoning_summary_text.delta            — 推理摘要文本增量
2. response.reasoning_summary_text.done             — 推理摘要文本完成
3. response.custom_tool_call_input.delta            — 自定义工具输入增量
4. response.custom_tool_call_input.done             — 自定义工具输入完成
5. response.computer_call.*（5 变体）                 — Computer Use 事件簇
   - response.computer_call.created
   - response.computer_call.delta
   - response.computer_call.done
   - response.computer_call.screenshot
   - response.computer_call.output_item.added
6. response.annotation.added                        — 注释添加
7. response.annotation.done                         — 注释完成
8. response.queued_item.*（3 变体）                  — 队列项事件簇
   - response.queued_item.added
   - response.queued_item.done
   - response.queued_item.failed
```

---

## 四、其他能力端点

### 4.1 Audio Speech（`/v1/audio/speech`）

```
POST /v1/audio/speech
{
  "model": "tts-1" | "gpt-4o-mini-tts",
  "input": "Hello world",
  "voice": "alloy" | "echo" | "fable" | "onyx" | "nova" | "shimmer" | "ash" | "ballad" | "coral" | "sage" | "verse",
  "response_format"?: "mp3" | "opus" | "aac" | "flac" | "wav" | "pcm16",
  "speed"?: 1.0,                                       // 0.25 - 4.0
  "instructions"?: "...",                              // gpt-4o-mini-tts
  "stream_format"?: "sse" | "audio"                    // 流式模式
}
// 响应：二进制音频（mp3 等），流式时为 SSE 事件
//   speech.audio.delta   { type, audio }
//   speech.audio.done    { type, audio }
```

### 4.2 Audio Transcription（`/v1/audio/transcriptions`）

```
POST /v1/audio/transcriptions  (multipart/form-data)
{
  "file": "<二进制音频>",
  "model": "whisper-1" | "gpt-4o-transcribe" | "gpt-4o-mini-transcribe",
  "language"?: "en",
  "prompt"?: "...",
  "response_format"?: "json" | "text" | "srt" | "verbose_json" | "vtt",
  "temperature"?: 0.0,
  "timestamp_granularities"?: ["word", "segment"],     // verbose_json 时
  "chunking_strategy"?: { "type": "auto" | "vad",
                          "prefix_pad_ms"?: 0, "silence_duration_ms"?: 200, "threshold"?: 0.5 },
  "include"?: [...],
  "stream"?: false
}
// 响应：
{
  "text": "Hello world",
  "language"?: "en",                                   // verbose_json
  "duration"?: 1.5,
  "words"?: [{ "word": "Hello", "start": 0.0, "end": 0.5 }],
  "segments"?: [{ "id": 0, "seek": 0.0, "start": 0.0, "end": 1.5, "text": "...",
                  "tokens": [...], "temperature": 0.0, "avg_logprob": -0.5,
                  "compression_ratio": 1.0, "no_speech_prob": 0.01 }],
  "usage"?: { "type": "tokens", "input_tokens": 100, "output_tokens": 50, "total_tokens": 150 }
}
```

### 4.3 Image Generation（`/v1/images/generations`）

```
POST /v1/images/generations
{
  "prompt": "A white cat",
  "model"?: "dall-e-2" | "dall-e-3" | "gpt-image-1",
  "n"?: 1,                                             // dall-e-2 可多张
  "quality"?: "standard" | "hd" | "low" | "medium" | "high" | "auto",
  "response_format"?: "url" | "b64_json",
  "size"?: "1024x1024" | "1024x1792" | "1792x1024" | "auto",
  "style"?: "vivid" | "natural",                       // 仅 dall-e-3
  "user"?: "user_xxx"
}
// 响应：
{
  "created": 1718000000,
  "data": [{ "url"?: "...", "b64_json"?: "...", "revised_prompt"?: "..." }],
  "usage"?: {                                          // gpt-image-1 支持
    "input_tokens": 100, "input_tokens_details": { "image_tokens": 80, "text_tokens": 20 },
    "output_tokens": 50, "total_tokens": 150
  }
}
```

### 4.4 Image Edit（`/v1/images/edits`）

```
POST /v1/images/edits (multipart/form-data)
{
  "image": "<base64 或 URL>",
  "prompt": "Add a hat",
  "mask"?: "<base64>",
  "model"?: "dall-e-2" | "gpt-image-1",
  "n"?: 1,
  "quality"?: "standard" | "hd" | "low" | "medium" | "high" | "auto",
  "response_format"?: "url" | "b64_json",
  "size"?: "1024x1024" | ...,
  "user"?: "user_xxx"
}
// 流式响应事件：
//   image_edit.partial_image   { type, partial_image_index, b64_json }
//   image_edit.completed       { type, usage?, b64_json? }
```

### 4.5 Embedding（`/v1/embeddings`）

```
POST /v1/embeddings
{
  "input": "Hello",                                    // string | string[]
  "model": "text-embedding-3-small" | "text-embedding-3-large" | "text-embedding-ada-002",
  "encoding_format"?: "float" | "base64",
  "dimensions"?: 1536,                                 // 覆盖默认维度
  "user"?: "user_xxx"
}
// 响应：
{
  "object": "list",
  "data": [{ "object": "embedding", "index": 0, "embedding": [0.1, -0.2, ...] }],
  "model": "text-embedding-3-small",
  "usage": { "prompt_tokens": 5, "total_tokens": 5 }
}
```

### 4.6 Moderation（`/v1/moderations`）

```
POST /v1/moderations
{
  "input": "...",
  "model"?: "omni-moderation-latest" | "text-moderation-latest" | "omni-moderation-2024-09-26"
}
// 响应：
{
  "id": "modr-xxx",
  "model": "omni-moderation-latest",
  "results": [{
    "flagged": false,
    "categories": {
      "hate": false, "hate/threatening": false,
      "harassment": false, "harassment/threatening": false,
      "self-harm": false, "self-harm/intent": false, "self-harm/instructions": false,
      "sexual": false, "sexual/minors": false,
      "violence": false, "violence/graphic": false
    },
    "category_scores": {
      "hate": 0.01, "hate/threatening": 0.001,
      "harassment": 0.02, ...,                            // 11 个类别分数
      "violence": 0.001, "violence/graphic": 0.0
    },
    "applied_input_risk_level"?: "low" | "medium" | "high"   // omni-moderation
  }]
}
```

---

## 五、错误响应（统一形态）

```
{
  "error": {
    "message": "Invalid API key provided.",
    "type": "invalid_request_error" | "authentication_error" | "permission_denied_error"
          | "not_found_error" | "rate_limit_error" | "server_error" | "billing_error"
          | "insufficient_quota_error" | "timeout_error",
    "code"?: "invalid_api_key",                         // 可选，错误代码
    "param"?: "messages",                               // 可选，错误参数名
    "status_code"?: 401                                 // 可选，HTTP 状态码
  }
}
```

HTTP 状态码语义：

| Code | 场景 |
|------|------|
| 400 | invalid_request_error — 请求体不合法 |
| 401 | authentication_error — 鉴权失败 |
| 403 | permission_denied_error — 权限不足 |
| 404 | not_found_error — 资源不存在 |
| 429 | rate_limit_error / insufficient_quota_error — 限速 / 配额 |
| 500 | server_error — 服务端错误 |
| 503 | server_error（暂时不可用） |
| 524 | timeout_error — 超时 |

---

## 六、ChatCompletion 流式 vs Responses 流式

| 维度 | ChatCompletion 流式 | Responses 流式 |
|------|-------------------|----------------|
| 事件格式 | 单一 ChatCompletion 对象（`object: "chat.completion.chunk"`），无独立 type 字段 | 独立 StreamEvent 对象（`type: "response.created"` 等 45 种） |
| 累积策略 | 客户端按 `id` + `index` 拼装 delta | 服务端按 `sequence_number` 严格保证顺序 |
| 终止 | 最后一个 chunk `finish_reason` 设置 | 必到 `response.completed` / `failed` / `incomplete` |
| Usage | 仅最后 chunk 携带（`include_usage=true` 时） | `response.completed` 事件 `response.usage` |
| 工具调用 | delta 中 `tool_calls` 增量 JSON 字符串 | 独立 `response.function_call_arguments.delta/done` 事件 |
| 推理 | 无内置（仅 reasoning_tokens 计数） | 独立 `response.reasoning_*` 事件 + 增量 summary |
| 注释/引用 | assistant message 中 `annotations` 数组 | OutputMessage.content[].annotations |
| 适合场景 | 简单 ChatCompletion / 工具调用 | 复杂代理任务（工具 / 推理 / MCP / 多模态） |

---

## 七、其他注意事项

### 7.1 ChatCompletion 双形态 Stop

```
"stop": "END"                                           // 字符串单形态
"stop": ["END", "STOP"]                                 // 字符串数组多形态
```

反序列化时通过 `OpenAIStopUnion` 接口分发。

### 7.2 旧版 vs 新版字段

| 旧版（已废弃） | 新版 |
|--------------|------|
| `function_call`（AssistantMessage 中） | `tool_calls[].function` |
| `function` 角色 message | `tool` 角色 + `tool_calls` |
| `temperature` 在 o 系列禁用 | 用 `reasoning_effort` 替代 |
| `max_tokens` 在 o 系列禁用 | 用 `max_completion_tokens` 替代 |
| `logprobs` 仅 content | content + refusal 双类别 |

### 7.3 ChatCompletion 旧版 Function 角色仍保留兼容

```
{ "role": "function", "content": "...", "name": "get_weather" }
```

仅在 `function_call` 调用场景下作为结果回传，新代码不应主动使用。

### 7.4 字段版本兼容性

| 字段 | 适用范围 | 备注 |
|------|----------|------|
| `developer` 角色 | 4.6+ | 优先级高于 `system` |
| `prediction` | 4.6+ | speculative decoding 加速 |
| `custom` 工具 | 4.6+ | 替代部分 hosted 工具场景 |
| `background` | 4.6+ | 长任务后台模式 |
| `conversation` | 4.6+ | 多轮会话状态引用 |
| `apply_patch` 工具 | 4.6+ | 文件编辑 |
| `computer_use_preview` | 4.6+ | 电脑使用（仍 preview） |
| `responses.reasoning.summary` | o 系列 | 推理摘要 |
| `mcp` 工具 | 4.6+ | 远程 MCP server 连接 |
| `safety_identifier` | 4.6+ | 滥用追踪（区别于 user） |

### 7.5 ChatCompletion Prompt Caching

- `prompt_cache_key`：手工指定缓存键，跨请求复用
- `prompt_tokens_details.cached_tokens`：缓存命中 token 数（费率折扣）
- `safety_identifier`：滥用检测 ID，可用于安全策略豁免

### 7.6 Responses 与 ChatCompletion 路由选择

| 场景 | 推荐端点 |
|------|----------|
| 简单 ChatCompletion | `/v1/chat/completions` |
| 工具调用 + 多模态 + 复杂推理 | `/v1/responses` |
| 后台长任务 | `/v1/responses` + `background: true` |
| File Search / Image Gen / MCP | `/v1/responses`（ChatCompletion 仍可用 hosted 工具） |
| 多轮状态管理 | `/v1/responses` + `conversation.id` |

---

## 八、字段基线参考

| 类型 | 字段数 | 来源 |
|------|--------|------|
| ChatCompletion Params | 30+ | `openai-go/chatcompletion.go` |
| ChatCompletion Message (6 种) | 6×5 ≈ 30 | 同上 |
| ChatCompletion ContentPart (7 种) | 7×3 ≈ 21 | 同上 |
| ChatCompletion Tool (6 种) | 6×5 ≈ 30 | 同上 |
| ChatCompletion ToolChoice (5 种) | 5 | 同上 |
| Responses Params | 20 | `openai-go/responses/` |
| Responses InputItem (12 种) | 12×5 ≈ 60 | 同上 |
| Responses OutputItem (10 种) | 10×5 ≈ 50 | 同上 |
| Responses Tool (10 种) | 10×5 ≈ 50 | 同上 |
| Responses StreamEvent (45 种) | 45×3 ≈ 135 | `openai-go/responses/stream.go` |
| Audio Speech | 8 | `openai-go/audio/speech.go` |
| Audio Transcription | 11 | `openai-go/audio/transcription.go` |
| Image Generate | 8 | `openai-go/image/generate.go` |
| Image Edit | 9 | `openai-go/image/edit.go` |
| Embedding | 5 | `openai-go/embeddings.go` |
| Moderation | 26（含 categories / scores 各 11） | `openai-go/moderation.go` |
| Error | 5 | `openai-go/error.go` |

---

## 九、参考

- OpenAI 官方 API 参考：<https://platform.openai.com/docs/api-reference>
- 官方 Go SDK：`github.com/openai/openai-go`
  - `chatcompletion.go` 4253 行
  - `responses/`（response.go + stream.go + input.go + output.go + tools.go，~22818 行）
  - `audio/speech.go` + `audio/transcription.go`
  - `image/generate.go` + `image/edit.go`
  - `embeddings.go` + `moderation.go`
- 现有 `internal/app/gateway/common/openai.go`（保留不动，作为网关适配层 wire-format 现状）
- 字段版本演进：4.5+ ↔ 4.6+ 重要变更（developer / prediction / custom tool /
  apply_patch / mcp / background / conversation）
