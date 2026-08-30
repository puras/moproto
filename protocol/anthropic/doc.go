// Package anthropic 定义 Anthropic Messages API 的 wire-format 类型。
//
// 本包基于 Anthropic 官方 API 参考（https://docs.anthropic.com/en/api/messages），
// 覆盖 Claude 4.6 / 4.7 / 4.8 / Fable 5 的全部请求字段与内容块类型。
//
// 与 internal/app/gateway/common/anthropic.go 的关系：
//   - common 包是网关适配层的精简版 wire-format（已被现有 adapter/provider 依赖）；
//   - 本包是 Anthropic 协议的权威定义，字段更完整、类型更严格（强类型 union）。
//
// 两者当前并存互不依赖，后续迁移在独立 PR 中完成。
//
// # 命名约定
//
//   - 类型统一使用 Anthropic 前缀，与 common 包同名前缀共存时无歧义；
//   - Param 后缀表示"客户端发出"类型（与 Anthropic 官方 SDK 对齐）；
//   - 多态 union 通过接口 + 私有 isXxx() 方法封闭实现集。
//
// # 类型分发机制（union）
//
// 旧版设计参考 common.ContentPart 的 partKind() 模式，由 union 接口内嵌公共
// Kind() 方法用于 type switch 之外的反射场景。该模式存在"双重事实源"问题：
// struct 上的 Type 字段（wire-format 标签）与 Kind() 返回值需要手动保持一致。
//
// 新版设计对齐官方 anthropic-sdk-go 的 ContentBlockUnion.AsAny() 模式：
//
//   - union 接口仅保留 isXxx() 私有方法封闭实现集；
//   - 各实现 struct 的 Type 字段承载 wire-format 类型标签（如 "text" / "tool_use"），
//     字段类型为独立的 string alias（如 AnthropicTextBlockType）+ 命名常量
//     （如 AnthropicTextBlockTypeText），编译期即保证 wire-format 字符串不被误传；
//   - 调用方通过 type switch on interface 或直接 switch on .Type 字段分发，
//     无需 Kind() 桥接方法；
//   - 新增类型只需新增 struct + Type string alias + 命名常量，零冗余样板代码。
//
// 与 common.ContentPart 的关系：common 包保留 partKind() 模式不变（本包重构
// 不影响 common），它是 gateway 内部通用多态模型（文本 / 图片 / 工具调用 /
// 工具结果）。两套并行存在，后续迁移在独立 PR 中完成。
package anthropic
