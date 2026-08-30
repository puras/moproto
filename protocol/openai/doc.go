// Package openai 定义 OpenAI API 的 wire-format 类型。
//
// 本包基于 OpenAI 官方 API 参考（https://platform.openai.com/docs/api-reference）
// + 官方 Go SDK 验证（chatcompletion.go 4253 行 + responses/response.go 22818 行），
// 覆盖 ChatCompletion + Responses + Audio + Image + Embedding + Moderation 全部
// 请求/响应/流事件类型。
//
// 与 internal/app/gateway/common/openai.go 的关系：
//   - common 包是网关适配层的精简版 wire-format（保留不动）；
//   - 本包是 OpenAI 协议的权威定义，字段更完整、类型更严格（强类型 union）。
//
// 两者当前并存互不依赖，后续迁移在独立 PR 中完成。
//
// # 命名约定
//
//   - 类型统一使用 OpenAI 前缀；
//   - Param 后缀表示"客户端发出"类型（Request 侧）；
//   - 多态 union 通过接口 + 私有 isXxx() 方法封闭实现集。
//
// # 类型分发机制（union）
//
// 与 anthropic-sdk-go 的 ContentBlockUnion.AsAny() 模式对齐：
//
//   - union 接口仅保留 isXxx() 私有方法封闭实现集；
//   - 各实现 struct 的 Type 字段承载 wire-format 类型标签（如 "text" / "function"），
//     字段类型为独立的 string alias（如 OpenAIChatCompletionMessageRole）+ 命名常量
//     （如 OpenAIChatCompletionMessageRoleUser），编译期即保证 wire-format 字符串不被误传；
//   - 调用方通过 type switch on interface 或直接 switch on .Type 字段分发。
//
// # 与 anthropic 包拆分风格对比
//
// 本包按"端点 + 端点内维度"横向拆分，区别于 anthropic 包的"按维度"拆分，
// 原因是 OpenAI 多端点类型几乎不共用（ChatCompletion ContentPart vs Responses
// InputItem 语义完全不同），按端点拆能让所有相关文件聚合。
package openai
