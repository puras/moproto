# `pkg/protocol/openai` — OpenAI API Wire-Format

## 概述

本包按 OpenAI 官方 API 参考 + 官方 Go SDK 验证（`chatcompletion.go` 4253 行 +
`responses/response.go` 等），完整定义 OpenAI 多端点 API 的 wire-format 类型，覆盖
ChatCompletion / Responses / Audio / Image / Embedding / Moderation 全部请求侧 +
响应侧 + 流事件 + 错误响应。

设计目标：

- **权威 wire-format 定义**：作为网关 `OpenAI -> X` 适配层的字段集基线
- **强类型 union**：多态字段（content / tools / tool_choice / output 等）通过
  interface + 私有 `isXxx()` 方法封闭实现集
- **零外部依赖**：本包仅依赖 stdlib `encoding/json`，不引 SDK 运行时

## 与 `internal/app/gateway/common/openai.go` 的关系

| 方面 | common/openai.go（保留不动） | pkg/protocol/openai（本包） |
|------|-----------------------------|----------------------------|
| 包 | `internal/app/gateway/common` | `pkg/protocol/openai` |
| 用途 | 网关适配层 wire-format（已被 adapter/provider 依赖） | OpenAI 协议权威定义 |
| Union 类型 | 部分缺失 / 用 `any` | 强类型 union 接口 |
| 字段集 | 基础字段 | 完整字段（含 4.6+ / Responses） |
| 依赖 | 不依赖 `pkg/protocol/openai` | 不依赖 `common`（避免循环） |

两者当前并存互不依赖。后续迁移路径（不在本 PR 范围）：由独立 PR 完成引用切换 + 删除旧文件。

## 文件清单（20 个 .go + 2 测试 + 2 文档）

### Request + Response 主体（按端点拆分）

| 文件 | 行数 | 职责 |
|------|------|------|
| `doc.go` | 35 | 包级注释 + 设计原则 |
| `common.go` | 100 | 共享 string alias + 命名常量（ReasoningEffort / Verbosity / ServiceTier / AudioVoice / AudioFormat / Modality / ImageDetail / ResponseFormatType / Truncation / PromptCacheRetention） |
| `chat_completion_params.go` | 257 | `OpenAIChatCompletionParams`（30+ 字段）+ `OpenAIStopUnion` + StreamOptions + Audio + Prediction + WebSearchOptions + ResponseFormatUnion + JSONSchema |
| `chat_completion_message.go` | 180 | 6 种 `OpenAIChatCompletionMessageParamUnion`（User / Assistant / System / Developer / Tool / Function）+ 各自 Content union |
| `chat_completion_content.go` | 131 | 8 种 `OpenAIChatCompletionContentPartUnion`（Text / Image / InputAudio / File / Video / Refusal / Audio） |
| `chat_completion_response.go` | 146 | `OpenAIChatCompletion` + Choice + ChoiceMessage（含 audio / annotations / logprobs）+ `OpenAICompletionUsage` + TokenLogprob |
| `chat_completion_tool.go` | 285 | 6 种 `OpenAIChatCompletionToolUnion`（Function / Custom / WebSearch / CodeInterpreter / FileSearch / ComputerUse）+ 5 种 `OpenAIChatCompletionToolChoiceOptionUnionParam` + 2 种 MessageToolCallUnion + 2 种 ChoiceMessageToolCallUnion |
| `chat_completion_chunk.go` | 75 | `OpenAIChatCompletionChunk` + ChunkChoice + Delta + DeltaToolCall + DeltaAudio |

### Responses API

| 文件 | 行数 | 职责 |
|------|------|------|
| `responses_params.go` | 144 | `OpenAIResponseParams` + InputUnion + InstructionsUnion + ToolChoiceUnion（占位）+ Reasoning + TextConfig + Conversation |
| `responses_input.go` | 226 | 12 种 `OpenAIResponseInputItemUnion`（Message / FunctionCall / FunctionCallOutput / FileSearchCall / WebSearchCall / ComputerCall / Reasoning / ItemReference / ApplyPatchCall / McpCall / CodeInterpreterCall / ImageGenerationCall） |
| `responses_output.go` | 179 | 10 种 `OpenAIResponseOutputItemUnion`（Message / FunctionCall / WebSearchCall / FileSearchCall / ComputerCall / Reasoning / McpCall / McpListTools / CodeInterpreterCall / ImageGenerationCall） |
| `responses_response.go` | 103 | `OpenAIResponse` 顶层（24 字段）+ Error + IncompleteDetails + Usage + InputTokensDetails + OutputTokensDetails |
| `responses_tool.go` | 191 | 11 种 `OpenAIResponsesToolUnion`（Function / WebSearch / CodeInterpreter / FileSearch / ImageGeneration / MCP / Computer / ApplyPatch / CodeExecution / Hosted）+ 3 种 ToolChoice |
| `responses_chunk.go` | 966 | 53 种 `OpenAIResponseStreamEventUnion`（状态 / 文本 / 推理 / 工具调用 / MCP / CodeInterpreter / ImageGen / Audio）+ `ParseResponseSSEEvent` SSE 解析函数 |
| `chat_completion_stream.go` | 151 | `OpenAIChatCompletionStreamEvent` union + chunk 结构 + `ParseSSEEvent`（已补齐 Moderation/ServiceTier/Logprobs 字段） |

### 其他能力

| 文件 | 行数 | 职责 |
|------|------|------|
| `audio_speech.go` | 40 | `OpenAIAudioSpeechParams` + 2 种 AudioSpeechStreamEvent（delta / done） |
| `audio_transcription.go` | 261 | `OpenAIAudioTranscriptionParams` + ChunkingStrategy + Transcription 响应 + Word / Segment / Usage + 3 种流事件 |
| `audio_translation.go` | 71 | `OpenAITranslationParams` + 5 种 ResponseFormat 命名常量 + Translation 基础/Verbose 响应 union |
| `image.go` | 106 | `OpenAIImageGenerateParams` + `OpenAIImageEditParams` + `OpenAIImagesResponse` + ImageData + Usage + 2 种 ImageEditStreamEvent |
| `embedding.go` | 45 | `OpenAIEmbeddingParams` + `OpenAIEmbeddingResponse` + `OpenAIEmbedding` + Usage |
| `moderation.go` | 61 | `OpenAIModerationParams` + Response + `OpenAIModeration` + 11 类别 Categories / CategoryScores |
| `error.go` | 36 | `OpenAIErrorResponse` + `OpenAIError` + `OpenAIErrorType`（9 种命名常量） |

### 测试 + 文档

| 文件 | 行数 | 职责 |
|------|------|------|
| `wire_format_chat_test.go` | 699 | ChatCompletion 一致性测试（11 用例 + union 封闭性） |
| `wire_format_responses_test.go` | 541 | Responses 一致性测试（10 用例 + union 封闭性） |
| `README.md` | - | 本文档 |
| `WIRE_FORMAT.md` | - | Wire-format 详细规范 |

**Total: ~3200 行（含测试 + 文档 ~4400 行）**

## 与 anthropic 包拆分风格对比

| 维度 | anthropic 包 | openai 包 |
|------|--------------|-----------|
| 拆分依据 | 按"维度"（request / message / tool / thinking） | 按"端点 + 端点内维度"（横向） |
| 共享类型 | Message / ContentBlock / Tool 全包共享 | 几乎不跨端点共享（ChatCompletion ContentPart vs Responses InputItem 语义完全不同） |
| 原因 | Anthropic Messages 是单一端点 | OpenAI 多端点类型几乎不共用 |

## 顶层类型概览

### ChatCompletion — `OpenAIChatCompletionParams`

30+ 字段，覆盖请求侧全字段集（必填 + 可选 + 流式 + 采样 + 工具 + 响应格式 +
存储 + 服务等级 + 音频输出 + 预测 + 推理 + 详细度 + Web 搜索 + 安全 + 缓存键）。

详见 `chat_completion_params.go` 和 `WIRE_FORMAT.md` 第一章。

### ChatCompletion Response — `OpenAIChatCompletion`

含 Choices + ChoiceMessage + Usage + ServiceTier + SystemFingerprint。
详见 `chat_completion_response.go`。

### Responses — `OpenAIResponseParams`

20 字段，覆盖必填（Model / Input）+ Instructions + Tools + Stream + 推理 +
Truncation + 文本配置 + 会话引用。详见 `responses_params.go`。

### Responses Response — `OpenAIResponse`

24 字段，含 ID / Status / Error / Output / Usage / 工具回放。详见
`responses_response.go`。

## 关键 Union 清单（17 个）

| Union | 文件 | 实现数 | 说明 |
|-------|------|--------|------|
| `OpenAIChatCompletionMessageParamUnion` | message.go | 6 | User / Assistant / System / Developer / Tool / Function |
| `OpenAIChatCompletionUserMessageContentUnion` | message.go | 2 | String / ContentParts |
| `OpenAIChatCompletionSystemMessageContentUnion` | message.go | 2 | String / ContentParts |
| `OpenAIChatCompletionDeveloperMessageContentUnion` | message.go | 2 | String / ContentParts |
| `OpenAIChatCompletionToolMessageContentUnion` | message.go | 2 | String / ContentParts |
| `OpenAIChatCompletionContentPartUnion` | content.go | 4 | Text / Image / InputAudio / File（拆分后仅 Request 侧 4 变体；Refusal/Audio 在 AssistantMessage 顶层字段） |
| `OpenAIChatCompletionToolUnion` | tool.go | 6 | Function / Custom / WebSearch / CodeInterpreter / FileSearch / ComputerUse |
| `OpenAIChatCompletionToolChoiceOptionUnionParam` | tool.go | 5 | Auto / Required / None / Allowed / HostedNamed |
| `OpenAIChatCompletionMessageToolCallUnion` | tool.go | 2 | Function / Custom |
| `OpenAIChatCompletionChoiceMessageToolCallUnion` | tool.go | 2 | Function / Custom |
| `OpenAIStopUnion` | params.go | 2 | String / Strings |
| `OpenAIResponseFormatUnion` | params.go | 3 | Text / JSONObject / JSONSchema |
| `OpenAIAudioVoiceUnion` | params.go | 2 | String / ID |
| `OpenAIResponseInputUnion` | responses_params.go | 2 | String / Items |
| `OpenAIResponseInstructionsUnion` | responses_params.go | 2 | String / Items |
| `OpenAIResponseToolChoiceUnion` | responses_tool.go | 10 | Auto / Required / None / Allowed / ApplyPatch / Custom / Function / Mcp / Shell / Types |
| `OpenAIResponsesToolUnion` | responses_tool.go | 17 | Function / WebSearch / CodeInterpreter / FileSearch / ImageGen / MCP / Computer / ApplyPatch / CodeExecution / Hosted / LocalShell / Shell / Custom / ToolSearch / WebSearchPreview / Namespace |
| `OpenAIResponseInputItemUnion` | responses_input.go | 24 | 见文件清单（4.6+ 新增 12 变体） |
| `OpenAIResponseOutputItemUnion` | responses_output.go | 23 | 见文件清单（4.6+ 新增 13 变体） |
| `OpenAIResponseInputMessageContentUnion` | responses_input.go | 2 | String / Parts（占位 `[]any`） |
| `OpenAIResponseStreamEventUnion` | responses_chunk.go | 53 | 见 `WIRE_FORMAT.md` 第三章 |
| `OpenAIChatCompletionToolCustomFormatUnion` | tool.go | 2 | Text / Grammar |
| `OpenAIAudioSpeechStreamEventUnion` | audio_speech.go | 2 | AudioDelta / AudioDone |
| `OpenAIImageEditStreamEventUnion` | image.go | 2 | PartialImage / Completed |

## 设计原则

1. **强类型 union 接口 + 私有 `isXxx()`**：每个 union 仅暴露私有封闭方法，
   类型分发通过 `switch v := b.(type)` 完成。零反射。

2. **`Param` 后缀 = 客户端发出**：Request 侧类型统一 `Param` 后缀，与官方 SDK
   命名对齐；Response 侧无后缀。

3. **`stop` / `messages.content` / `input` / `instructions` 等双形态用 union**：
   `string | []Item` 等形态用 union 接口（`OpenAIStopUnion` /
   `OpenAIResponseInputUnion` 等）实现，对应 JSON `string` / 数组两种形态。

4. **嵌入字段用结构化类型**：枚举字段（Tool Type / Message Role /
   ReasoningEffort / Verbosity 等）使用独立 string alias + 命名常量，编译期即
   防止 wire-format 字符串拼写错误。

5. **零外部依赖**：本包仅依赖 `encoding/json`，不引 `apijson` / `param` /
   `respjson` / `constant` 等 SDK 运行时。

## 使用示例

```go
import openai "hegq.net/mogo/pkg/protocol/openai"

// 反序列化 ChatCompletion 请求
body := []byte(`{
    "model": "gpt-4o",
    "messages": [
        {"role": "user", "content": "Hello"}
    ],
    "stream": true,
    "tools": [
        {"type": "function", "function": {"name": "get_weather"}}
    ]
}`)
var req openai.OpenAIChatCompletionParams
if err := json.Unmarshal(body, &req); err != nil {
    log.Fatal(err)
}
// req.Model == "gpt-4o"
// len(req.Messages) == 1
// len(req.Tools) == 1

// 类型分发（content union）
switch m := req.Messages[0].(type) {
case openai.OpenAIChatCompletionUserMessageParam:
    switch c := m.Content.(type) {
    case openai.OpenAIChatCompletionUserMessageContentString:
        fmt.Println("Text:", string(c))
    case openai.OpenAIChatCompletionUserMessageContentParts:
        // 多模态
    }
case openai.OpenAIChatCompletionAssistantMessageParam:
    // 助手消息
}

// Responses API 反序列化
respBody := []byte(`{
    "id": "resp_xxx",
    "object": "response",
    "status": "completed",
    "output": [...]
}`)
var resp openai.OpenAIResponse
if err := json.Unmarshal(respBody, &resp); err != nil {
    log.Fatal(err)
}

// 类型分发 output items
for _, item := range resp.Output {
    switch v := item.(type) {
    case openai.OpenAIResponseOutputMessage:
        // 文本消息
    case openai.OpenAIResponseOutputFunctionCall:
        // 函数调用
    case openai.OpenAIResponseOutputReasoning:
        // 推理（o 系列）
    }
}
```

## 设计取舍

1. **ChatCompletion ContentPartUnion 拆分（4.6+）**：对齐官方 SDK，仅保留 Request 侧
   4 变体（Text / Image / InputAudio / File）。Refusal / Audio 原本混入 ContentPart，
   已被提取到 `OpenAIChatCompletionAssistantMessageParam` 的顶层 `Refusal *string`
   和 `Audio *OpenAIChatCompletionMessageAudioParam` 字段（语义与官方 SDK 的
   AssistantMessageParamContentArrayOfContentPartUnion 一致）。mogo 多余的
   `video_url` ContentPart 已删除（官方不存在该 ContentPart 变体）。

2. **Responses StreamEvent SSE 解析**：`ParseResponseSSEEvent(line string)` 函数
   在 `responses_chunk.go` 末尾实现，支持解析全部 53 种 SSE 事件类型。

3. **Responses InputMessage Content Parts**：`OpenAIResponseInputMessageContentParts`
   当前以 `[]any` 占位（简化实现），后续 PR 替换为强类型 `[]ContentPartUnion`。

4. **ChatCompletion Stream Chunk 增量字段**：`OpenAIChatCompletionChunkChoiceDelta`
   中的 `FunctionCall` 增量字段保留兼容（官方已废弃）。

## 已补齐项（P0，2026-07-01）

| 项 | 状态 | 说明 |
|----|------|------|
| `OpenAIReasoningEffort` 补 `none` / `xhigh` | ✅ | gpt-5.1 / codex-max 必需 |
| `OpenAIServiceTier` 补 `scale` | ✅ | 4.6+ 服务等级 |
| `OpenAIAudioVoice` 补 `marin` / `cedar` | ✅ | 4.6+ 新增语音 |
| AudioTranslation 整个端点 | ✅ | `OpenAITranslationParams` + `Translation` / `TranslationVerbose` 响应 union |
| Responses Tool 补 6 变体 | ✅ | LocalShell / Shell / Custom / ToolSearch / WebSearchPreview / Namespace |
| Responses ToolChoice 补 7 变体 | ✅ | Allowed / ApplyPatch / Custom / Function / Mcp / Shell / Types |
| Responses InputItem 补 12 变体 | ✅ | 24 变体（对齐官方主线 31 中核心子集） |
| Responses OutputItem 补 13 变体 | ✅ | 23 变体（对齐官方主线 26 中核心子集） |
| ChatCompletion ContentPartUnion 拆分 | ✅ | 7→4 变体，Refusal/Audio 抽到 AssistantMessage |
| `OpenAIResponse.ToolChoice/Tools` any 占位替换 | ✅ | 强类型 union |
| Responses StreamEvent 53 变体完整覆盖 | ✅ | `OpenAIResponseStreamEventUnion` 53 种事件 + `ParseResponseSSEEvent` |
| ChatCompletion SSE 补 Moderation/ServiceTier/Logprobs | ✅ | `OpenAIChatCompletionStreamChunkEvent` 补齐三个缺失字段 |

## 已补齐项（P1，2026-07-01）

| 项 | 状态 | 说明 |
|----|------|------|
| Responses `MaxToolCalls` | ✅ | `OpenAIResponseParams.MaxToolCalls` |
| Responses `PreviousResponseID` | ✅ | `OpenAIResponseParams.PreviousResponseID` |
| Responses `Prompt`（`ResponsePromptParam`） | ✅ | `OpenAIResponsePrompt` + `OpenAIResponsePromptVariableUnion`（含 input_text / input_image / input_file） |
| Responses `Moderation`（请求侧） | ✅ | `OpenAIResponseModeration` |
| Responses `StreamOptions` | ✅ | `OpenAIResponseStreamOptions{IncludeObfuscation}` |
| Responses `TopLogprobs` | ✅ | `OpenAIResponseParams.TopLogprobs` |
| Responses `ContextManagement` | ✅ | `OpenAIResponseContextManagement`（含 `compact_threshold`） |
| Image `OpenAIImageModel` 类型 + 8 常量 | ✅ | gpt-image-1 / gpt-image-1-mini / gpt-image-1.5 / gpt-image-2 / chatgpt-image-latest / dall-e-2 / dall-e-3 |
| Image `ImageGenerateParams` gpt-image-1+ 字段 | ✅ | `Background` / `OutputCompressionQuality` / `PromptRefinement` |
| `ImageVariationParams` 端点 | ✅ | `OpenAIImageVariationParams`（dall-e-2 衍生图） |

## 已补齐项（P2，2026-07-01）

| Responses 流事件 8 个次要变体 | ✅ | reasoning_summary_text / custom_tool_call_input / computer_call.* / annotation.* / queued_item.* 全部补齐 |
| `OpenAIResponseInputMessageContentParts []any` 强类型化 | ✅ | `OpenAIResponseInputMessageContentPartUnion`（含 input_text / input_image / input_file） |
| ChatCompletion 旧版 `Functions` / `FunctionCall` | ✅ | 已有 `FunctionCall` + `FunctionMessageParam` + `FunctionToolParam` + Tools 中的 function 类型 |
| Responses `InputMessageContentPartUnion` | ✅ | 与 P2-2 配套的三变体 union |

## 已补齐项（P3，2026-07-01）

| 项 | 状态 | 说明 |
|----|------|------|
| CompactionItem 字段深化 | ✅ | `OpenAIResponseOutputCompactionItem` 新增 `CreatedBy` |
| `AdditionalTools` InputItem / OutputItem | ✅ | `OpenAIResponseInputAdditionalToolsItem` + `OpenAIResponseOutputAdditionalTools` |
| ResponseInputItem 缺 6 变体 | ✅ | `CustomToolCall` / `CustomToolCallOutput` / `AdditionalTools` / `FunctionCallOutput` / `FunctionWebSearch` / `ToolSearchOutputItem` |
| ResponseOutputItem 缺 FunctionWebSearch | ✅ | `OpenAIResponseOutputFunctionWebSearch` |
| ResponseOutputText `Logprobs` 字段 | ✅ | `OpenAIResponseOutputMessageContent.Logprobs` + Logprob + TopLogprob 子结构 |

## 验证

```bash
cd backend
go test ./pkg/protocol/openai/... -v -count=1   # 全部 PASS
go vet ./pkg/protocol/openai/...                # 无错误
go build ./...                                  # 全量回归通过
gofmt -l pkg/protocol/openai/                  # 无输出（格式干净）
```

## 参考

- OpenAI 官方 API：<https://platform.openai.com/docs/api-reference>
- 官方 Go SDK：`github.com/openai/openai-go`（作为字段集基线）
  - `chatcompletion.go` 4253 行
  - `responses/response.go` ~22818 行（Responses + StreamEvent）
  - `audio/` + `image/` + `embeddings/` + `moderation.go`
- 现有 `internal/app/gateway/common/openai.go`（保留不动，作为网关适配层 wire-format 现状）
