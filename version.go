// Package moproto 是 AI 协议解析层 Go 库的镜像发布根。
//
// 本文件集中维护模块版本号，作为外部消费者通过 -ldflags 注入或直接 import
// 引用时的唯一权威源。
//
// 版本同步策略：
//   - 本版本号与 git tag 一一对应；
//   - 每次发布由 maintainer 在 tag 同步更新本常量；
//   - 消费者可通过 go build -ldflags "-X github.com/puras/moproto.Version=vX.Y.Z"
//     在编译期覆盖。
package moproto

// Version 是 moproto 模块的当前发布版本。
//
// 该值必须与最近的 git tag（vX.Y.Z 形式，去掉前缀 v 与 Version 一致）保持一致。
const Version = "0.1.0"

// VersionMeta 汇总版本相关元数据，供外部探针/诊断/上报使用。
//
// 所有字段由 maintainer 在 tag 时人工填写，禁止在版本号 bump 之外修改。
var VersionMeta = struct {
	// Version 与 Version 常量一致，便于序列化场景（避免反射常量）。
	Version string
	// CompatibleWith 列出已验证可双向转换的源模块版本（mogo pkg 路径）。
	CompatibleWith []string
	// ReleasedAt 为 ISO 8601 发布日期（UTC）。
	ReleasedAt string
}{
	Version: Version,
	CompatibleWith: []string{
		"hegq.net/mogo/pkg/protocol/openai",
		"hegq.net/mogo/pkg/protocol/anthropic",
		"hegq.net/mogo/pkg/protocol/gemini",
		"hegq.net/mogo/pkg/usage",
	},
	ReleasedAt: "2026-09-04",
}