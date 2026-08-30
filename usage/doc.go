// Package usage 提供协议无关的统一用量（Usage）抽象。
//
// 设计动机：
//   - 网关计费/用量统计/渠道成本核算需要单一数据结构，避免对每种协议
//     （Anthropic / OpenAI / 后续 DeepSeek / Gemini 等）做协议特化处理；
//   - 不同协议的 Usage 字段命名/嵌套结构不一致（如 Anthropic 用 InputTokens/
//     OutputTokens，OpenAI ChatCompletion 用 PromptTokens/CompletionTokens），
//     转换到统一 Usage 后，下游消费方只面对一套语义。
//
// # 依赖方向
//
//	pkg/usage  ──依赖──>  pkg/protocol/anthropic
//	pkg/usage  ──依赖──>  pkg/protocol/openai
//	pkg/protocol/*  ──不依赖──>  pkg/usage   （低层不引高层）
//
// 禁止反向依赖：pkg/protocol/* 不得 import 本包。
//
// # 命名约定
//
//   - 核心计费字段采用 Anthropic 命名（InputTokens / OutputTokens / TotalTokens），
//     作为唯一权威字段；OpenAI 的 PromptTokens / CompletionTokens 在 FromOpenAICompletion
//     转换函数内部一次性映射到 InputTokens / OutputTokens，调用方不再看到别名。
//   - 嵌套分组（Cache / Reasoning / Thinking / ServerTools / Audio / Prediction）
//     用指针可空，避免对协议不支持维度伪造数据。
//
// # 流式聚合
//
// Usage.Add(other *Usage) 累加所有可累加字段，供 protocol/stream_handler.go
// 在 WriteStream 阶段逐 delta 合并用量。
package usage
