# Gemini Wire-Format

## 字段对齐表

| wire 名称 | Go 字段 | 类型 |
|---|---|---|
| `contents` | `GeminiGenerateContentRequest.Contents` | `[]*GeminiContent` |
| `parts` | `GeminiContent.Parts` | `[]*GeminiPart` |
| `inlineData` | `GeminiPart.InlineData` | `*GeminiBlob` |
| `functionCall` | `GeminiPart.FunctionCall` | `*GeminiFunctionCall` |
| `functionResponse` | `GeminiPart.FunctionResponse` | `*GeminiFunctionResponse` |
| `generationConfig` | `GeminiGenerateContentRequest.GenerationConfig` | `*GeminiGenerationConfig` |
| `usageMetadata` | `GeminiGenerateContentResponse.UsageMetadata` | `*GeminiUsageMetadata` |

## 与官方 SDK 对齐范围

- Part union: 7 个字段变体
- Content: Role 仅 user | model
- FunctionCall: ID 可选（SDK 一致）
- UsageMetadata: 6 个顶层计量字段 + 2 个 modality 明细数组
- Schema: OpenAPI 3.0 子集（递归 object + array + enum + required）
- 流协议：每个 `data:` 帧是完整 GenerateContentResponse

### 流式 passthrough：RawData carrier

同协议（Gemini → Gemini）流式响应经 gateway 转发时，`GeminiGenerateContentResponse.RawData []byte` 装载上游 `data:` 后原始 JSON 字节，`MarshalJSON` 优先返回 `RawData`（passthrough 路径），避免 typed struct 往返剥除未知字段（`groundingMetadata` / `urlContextMetadata` / `avgLogprobs` 等）与 `omitempty` 零值字段。

非流式路径（`RawData` 为 nil）自动 fallback 到 typed marshal，行为不变。详见 `README.md`「同协议流式 passthrough」章节。

## 不在 P0

- Vertex-only 字段（routingConfig / labels）
- Server-side tools（googleSearch / codeExecution / urlContext）
- speechConfig / imageConfig / audioTimestamp
