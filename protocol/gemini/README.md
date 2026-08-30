# pkg/protocol/gemini

Google Gemini Generative Language API wire-format 定义（仅 stdlib 依赖）。

## 范围

P0 仅 ChatCompletion（generateContent + streamGenerateContent?alt=sse）：
- text / image / file / function call Part
- Thinking / SystemInstruction / Safety
- Schema（OpenAPI 3.0 子集）
- UsageMetadata
- StreamEvent（每帧 = 完整 GenerateContentResponse）
- APIError 透传

不在 P0：Vertex AI / ImageGen / Audio / Embedding / Cache / Live / Batches / Files。

## 子包

- `pkg/protocol/gemini/frames` — SSE 帧 Marshal / 错误响应 WriteGeminiError

## 文件清单

| 文件 | 职责 |
|---|---|
| common.go | 命名常量 |
| content.go / blob.go | Content、Part、Blob |
| tool.go / schema.go / safety.go | 工具、Schema、安全 |
| generation_config.go | 生成与思考配置 |
| request.go / response.go | 请求与响应 |
| streaming.go | SSE 解析与用量提取 |
| error.go | API 错误 |

## 与 anthropic / openai 包的差异

| 维度 | Gemini | Anthropic | OpenAI |
|---|---|---|---|
| 鉴权 header | `x-goog-api-key` | `x-api-key` + `anthropic-version` | `Authorization: Bearer` |
| 模型位置 | URL path | body.model | body.model |
| 流事件 | 每帧完整 Response | 事件 union | chunk union |
| 错误格式 | `{error: {code, message, status}}` | `{type, error, request_id}` | `{error: {code, message}}` |

## 同协议流式 passthrough（raw bytes 透传）

> 起于：commit `3439d65`（Anthropic + Gemini 对称 passthrough 重构 / T3）。
> 单元测试：`pkg/protocol/gemini/streaming_test.go` `TestParseGeminiStreamChunk_PreservesRawBytes` + `frames/frame_test.go` `TestMarshalSSEFrame_RawDataPassthrough`。
> E2E 回归：`test/e2e/l2/streaming_test.go` `TestL2_Streaming_GeminiPassthrough`。

上游是 Gemini SSE 的同协议请求路径下，gateway 不做 typed struct 往返，而是直接把整帧 `data: Y\n\n` 原始字节透传给下游，避免两层 loss：

| Loss 点 | 来源 | 影响 |
|---------|------|------|
| LOSS #1（typed unmarshal） | `json.Unmarshal(data, &GeminiGenerateContentResponse)` | typed struct 未声明的未知字段（如 `groundingMetadata` / `urlContextMetadata` / `avgLogprobs`）整字段消失 |
| LOSS #2（typed marshal） | `json.Marshal(resp)` 配合 `omitempty` | `avgLogprobs` 等零值字段被剥 |

透传路径关键载体：

| 名称 | 路径 | 作用 |
|------|------|------|
| `GeminiGenerateContentResponse.RawData []byte` | `pkg/protocol/gemini/response.go` | 上游原始 `data:` 后 JSON 字节；带 `json:"-"` 标签不参与 typed marshal |
| `(*GeminiGenerateContentResponse).MarshalJSON` | `pkg/protocol/gemini/response.go` | 优先返回 `RawData`（passthrough 路径），否则 fallback typed marshal |
| `ParseGeminiStreamChunk(line, rawFrame)` | `pkg/protocol/gemini/streaming.go` | 把 rawFrame 装入 `resp.RawData`，旁路 unmarshal 仅供 usage 抽取 |
| `frames.MarshalSSEFrame(resp)` | `pkg/protocol/gemini/frames/frame.go` | `RawData` 非空时 verbatim emit `data: <raw>\n\n`，否则 typed marshal 走 fallback |

非流式路径不受影响：`RawData` 为 nil 时自动走 typed marshal（fallback 行为不变）。
