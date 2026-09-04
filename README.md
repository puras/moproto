# MoProto

AI 协议解析层 Go 库。

## 来源

本项目代码镜像自 [mogo](https://example.invalid/mogo) 后端的
`pkg/protocol/{openai,anthropic,gemini}` 与 `pkg/usage/` 目录，
迁移日期 2026-08-28。

迁移原则：

- 仅替换 import 路径（`hegq.net/mogo/pkg/` → `hegq.net/moproto/`）
- 包声明（`package openai` 等）保持原样
- 不引入 mogo 模块依赖

## 范围

### 已迁移 provider

- `protocol/openai` — OpenAI ChatCompletion / Responses / Audio / Embedding / Image / Moderation
- `protocol/anthropic` — Anthropic Messages（含 MCP、Thinking、ContentBlock）
- `protocol/gemini` — Google Gemini GenerateContent / StreamGenerateContent
- `usage` — 各 provider 用量（token）解析

### 未迁移

- 错误解析器（依赖 mogo 内部包）
- 跨协议转换 stub（mogo 中尚未实现）
- 业务 SSOT（capability 常量）

## 用法

```go
import (
    "hegq.net/moproto/protocol/openai"
    "hegq.net/moproto/usage"
)

// 解析 OpenAI ChatCompletion 响应
var resp openai.OpenAIChatCompletion
if err := json.Unmarshal(rawBody, &resp); err != nil { /* ... */ }

// 提取用量
u, err := usage.ParseOpenAI(&resp)
```

## 三方依赖

无第三方依赖。所有测试均使用标准库 `testing` 包。

（历史：本项目早期使用 `github.com/stretchr/testify` 作断言库，后续重构为标准库 `testing`，go.mod 已清空。）

## 版本

当前发布版本：`v0.1.0`（发布于 2026-09-04）。

可通过根包常量读取版本号：

```go
import "github.com/puras/moproto"

fmt.Println(moproto.Version)            // "0.1.0"
fmt.Println(moproto.VersionMeta.ReleasedAt)
// moproto.VersionMeta.CompatibleWith 列出已验证可双向转换的源模块路径
```

或通过 `-ldflags` 在编译期注入：

```bash
go build -ldflags "-X github.com/puras/moproto.Version=0.1.0" ./...
```

## 版本同步策略

本项目作为 mogo 协议层的独立镜像发布。当 mogo 内协议结构有更新时，
按以下步骤同步：

1. 在 mogo 内执行目标 provider 单元测试，记录基线
2. 从 mogo 复制新版本目录到 moproto 对应路径
3. 重跑 Task 3 的 sed 替换
4. 在 moproto 内执行 `go mod tidy && go test ./...`
5. 对比 mogo 与 moproto 测试结果应一致
