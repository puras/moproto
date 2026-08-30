// Package gemini 定义 Google Gemini Generative Language API 的 wire-format 类型。
//
// 范围（P0）：
//   - ChatCompletion 能力（generateContent + streamGenerateContent?alt=sse）
//   - 函数调用（FunctionCall / FunctionResponse / FunctionDeclaration）
//   - Thinking / SystemInstruction / Safety
//
// 不在 P0：
//   - Vertex AI 后端（需 OAuth/ADC/project/location）
//   - 图像生成 / Audio / Embedding / Cache / Live / Batches
//   - Server-side tools（googleSearch / codeExecution / urlContext）
//   - OpenAI/Anthropic ↔ Gemini 跨协议转换
//
// 参考 SDK：/Users/puras/workspace/source/go-genai（v1.65.0）
// 本包零外部依赖，仅 stdlib encoding/json。
//
// # 命名约定
//
//   - 类型统一使用 Gemini 前缀；
//   - wire-format 字段保留 Google 原生 camelCase（contents / systemInstruction /
//     usageMetadata），JSON tag 严格遵守；
//   - sealed union 接口仅暴露 isXxx() 私有方法，调用方按 .Type 字段分发。
package gemini
