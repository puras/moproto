# MoProto：AI 协议解析层迁移设计

**日期**：2026-08-28
**作者**：Claude（brainstorming 会话）
**状态**：待用户审阅
**强约束**：MoGo 工程（`/Users/puras/workspace/proj/106hz/mogo`）**只读**，本迁移不得修改 mogo 任何文件。

## 1. 目标

将 mogo 工程的 AI 协议解析层（openai / anthropic / gemini 三大 provider 的协议结构、流、事件、帧；以及跨 provider 的用量解析 usage）整体迁移到 MoProto 项目，作为**独立的、可被外部项目 import 的 Go 解析库**，为后续 AI 协议解析需求提供基础能力。

迁移原则：

- **零业务改动**：mogo 原文件除包名/import 路径外不修改任何逻辑
- **零 mogo 依赖**：moproto 不引入 `hegq.net/mogo/...` 任何路径
- **可独立编译与测试**：迁移完成后 `go build ./...` 与 `go test ./...` 在 moproto 内可独立跑通

## 2. 范围

### 2.1 复制（目标清单）

| 源（mogo） | 目标（moproto） | 说明 |
|---|---|---|
| `pkg/protocol/openai/` | `protocol/openai/` | OpenAI 协议结构/事件/帧 |
| `pkg/protocol/anthropic/` | `protocol/anthropic/` | Anthropic 协议结构/事件/帧 |
| `pkg/protocol/gemini/` | `protocol/gemini/` | Gemini 协议结构/事件/帧 |
| `pkg/usage/` | `usage/` | 各 provider 用量解析 |

每个子包内的所有 `.go` 与 `_test.go` 文件**全部复制**。

### 2.2 不复制

| 路径 | 不复制原因 |
|---|---|
| `pkg/protocol/capability.go` | 业务 SSOT，含 mogo 业务术语 |
| `pkg/parser/error_parser/` | 直接 `import hegq.net/mogo/internal/app/gateway/...` |
| `pkg/adapter/` | 当前为 `ErrNotImplemented` stub，无实际解析能力 |
| `pkg/translate/` | 同上 stub |
| `pkg/types/types.go` | 仅一个两常量枚举，不属于协议解析 |
| `pkg/protocol/video/`, `pkg/protocol/seedance/` | 范围外（非三大 LLM 协议） |

## 3. 目录结构

```
/Users/puras/workspace/proj/hegq/moproto/
├── go.mod                            # module hegq.net/moproto
├── go.sum                            # go mod tidy 生成
├── README.md                         # 简述本库、来源、用法
├── docs/
│   └── superpowers/
│       └── specs/
│           └── 2026-08-28-moproto-migration-design.md
├── protocol/
│   ├── openai/                       # 自 mogo/pkg/protocol/openai
│   │   ├── chat_completion_params.go
│   │   ├── chat_completion_response.go
│   │   ├── chat_completion_chunk.go
│   │   ├── chat_completion_message.go
│   │   ├── chat_completion_stream.go
│   │   ├── chat_completion_content.go
│   │   ├── chat_completion_tool.go
│   │   ├── common.go
│   │   ├── error.go
│   │   ├── audio_speech.go
│   │   ├── audio_transcription.go
│   │   ├── audio_translation.go
│   │   ├── embedding.go
│   │   ├── image.go
│   │   ├── moderation.go
│   │   ├── responses_params.go
│   │   ├── responses_input.go
│   │   ├── responses_output.go
│   │   ├── responses_response.go
│   │   ├── responses_chunk.go
│   │   ├── responses_compaction.go
│   │   ├── responses_tool.go
│   │   ├── doc.go
│   │   ├── frames/
│   │   │   ├── frame.go
│   │   │   ├── errors.go
│   │   │   └── responses_raw_stream.go
│   │   └── *_test.go
│   ├── anthropic/                    # 自 mogo/pkg/protocol/anthropic
│   │   ├── request.go
│   │   ├── response.go
│   │   ├── message.go
│   │   ├── context.go
│   │   ├── context_response.go
│   │   ├── stream_event.go
│   │   ├── stream_event_delta.go
│   │   ├── citation.go
│   │   ├── citation_getter_test.go
│   │   ├── content_block.go
│   │   ├── content_block_subresult.go
│   │   ├── tool.go
│   │   ├── thinking.go
│   │   ├── output.go
│   │   ├── source.go
│   │   ├── container.go
│   │   ├── mcp.go
│   │   ├── usage.go
│   │   ├── error.go
│   │   ├── common.go
│   │   ├── doc.go
│   │   ├── frames/
│   │   │   ├── frame.go
│   │   │   ├── errors.go
│   │   │   ├── builder.go
│   │   │   └── message_id.go
│   │   └── *_test.go
│   └── gemini/                       # 自 mogo/pkg/protocol/gemini
│       ├── request.go
│       ├── response.go
│       ├── content.go
│       ├── generation_config.go
│       ├── safety.go
│       ├── schema.go
│       ├── tool.go
│       ├── streaming.go
│       ├── blob.go
│       ├── error.go
│       ├── common.go
│       ├── doc.go
│       ├── frames/
│       │   ├── frame.go
│       │   └── errors.go
│       └── *_test.go
└── usage/                            # 自 mogo/pkg/usage
    ├── usage.go                      # 公共 Usage 与工具
    ├── openai.go
    ├── anthropic.go
    ├── gemini.go
    ├── doc.go
    └── *_test.go
```

> 注：上表 `_test.go` 行为占位，实际复制以 mogo 源目录为准。

## 4. 包名与 import 路径变更

### 4.1 包声明

**所有 `.go` 文件的包声明保持不变**（`package openai` / `package anthropic` / `package gemini` / `package usage` 都不改名）。包名属于公开 API 的一部分，重命名会破坏调用方。

### 4.2 跨包 import 路径替换

| 旧（mogo） | 新（moproto） |
|---|---|
| `"hegq.net/mogo/pkg/protocol/openai"` | `"hegq.net/moproto/protocol/openai"` |
| `"hegq.net/mogo/pkg/protocol/anthropic"` | `"hegq.net/moproto/protocol/anthropic"` |
| `"hegq.net/mogo/pkg/protocol/gemini"` | `"hegq.net/moproto/protocol/gemini"` |
| `"hegq.net/mogo/pkg/usage"` | `"hegq.net/moproto/usage"` |

替换规则：所有 `.go` 文件中出现的 `hegq.net/mogo/pkg/` 一次性替换为 `hegq.net/moproto/`。

### 4.3 文件内 doc/注释中的 import 路径

mogo 注释里有多处形如 `// 用于满足 internal/app/gateway/common.Request 接口` 或 `// 协议层（如 internal/app/gateway/protocol）...` 的引用。这些**是注释，不影响编译**，但仍建议替换为 moproto 等价路径（`internal/app/...` 是 mogo 业务概念，迁移后语义变空；最简方案是保留原注释不动以最大化"零改动"原则——本设计选择**保留**）。

## 5. go.mod 与三方依赖

### 5.1 go.mod

```go
module hegq.net/moproto

go 1.26

require (
    github.com/stretchr/testify v1.11.1
)
```

如果 property-based 测试（`leanovate/gopter`）实际被使用，再追加：

```go
    github.com/leanovate/gopter v0.2.11
```

`go mod tidy` 后由工具自动补全 `go.sum`。

### 5.2 依赖策略

- 仅引入 mogo 子包**实际使用**的三方依赖（实际只 `testify`，且为测试用）
- **不引入** `hegq.net/mogo/...` 任何路径
- **不引入** mogo 的业务三方依赖（gin、redis、stripe、zap 等）

### 5.3 验证

`go mod why` + `go list -m all` 应不出现 `hegq.net/...`。

## 6. 复制流程

执行步骤（按顺序）：

1. **复制目录**：
   ```bash
   SRC=/Users/puras/workspace/proj/106hz/mogo/backend/pkg
   DST=/Users/puras/workspace/proj/hegq/moproto

   mkdir -p $DST/protocol $DST/usage
   cp -R $SRC/protocol/openai    $DST/protocol/openai
   cp -R $SRC/protocol/anthropic $DST/protocol/anthropic
   cp -R $SRC/protocol/gemini    $DST/protocol/gemini
   cp -R $SRC/usage              $DST/usage
   ```

2. **删除 mogo 内的 go.mod / go.sum**（如被复制过来）：
   ```bash
   find $DST -name 'go.mod' -not -path "$DST/go.mod" -delete
   find $DST -name 'go.sum' -delete
   ```

3. **替换 import 路径**（在所有 .go 与 _test.go 内）：
   ```bash
   find $DST -name '*.go' -print0 | xargs -0 sed -i '' \
       's|hegq.net/mogo/pkg/|hegq.net/moproto/|g'
   ```

4. **创建 moproto 顶层 go.mod**（见 5.1）。

5. **生成 go.sum 与验证**：
   ```bash
   cd $DST
   go mod tidy
   go build ./...
   go vet ./...
   go test ./...
   ```

6. **写入 README.md**：简述来源（mogo v0.x）、用法、版本同步策略。

## 7. 错误处理与风险

### 7.1 风险 1：mogo 子包偷偷依赖 mogo 内部包

**检测**：复制前先在 mogo 目录内验证：

```bash
grep -rl 'hegq.net/mogo/internal' \
    /Users/puras/workspace/proj/106hz/mogo/backend/pkg/protocol \
    /Users/puras/workspace/proj/106hz/mogo/backend/pkg/usage \
    --include='*.go' | grep -v _test
```

（已在 2026-08-28 探索阶段确认：仅 `pkg/parser/error_parser/error_parser.go` 引用 `internal/`，目标范围 protocol/ 与 usage/ **无此问题**。）

**应对**：如发现新的引用文件，单独审视：

- 若仅是注释：保留注释不动
- 若是非测试代码引用 internal：评估是否需要将该文件从迁移清单中剔除

### 7.2 风险 2：测试文件有 mogo 特有断言

**检测**：`grep` 测试文件中是否有 `mogo/` 路径引用或 mogo 业务包符号。

**应对**：基本无（已确认）。如有，单独处理。

### 7.3 风险 3：三方依赖被遗漏

**检测**：`go mod tidy` 报错即缺包。

**应对**：补 go.mod；记录到 README 三方依赖说明。

### 7.4 风险 4：复制过程中 mogo 目录被误改

**检测**：所有 `cp -R` 都是写入 moproto；mogo 目录只在只读模式下访问（ls / grep / sed 都不会修改源）。

**应对**：所有命令的目标路径都在 `moproto/`，源路径只作为参数出现在 `cp` / `sed -i` 的输入端，**不修改 mogo**。

## 8. 测试策略

- 全部 `_test.go` 复制到 moproto，**不重写**测试逻辑
- 包声明 `package xxx` 与 `package xxx_test` 均保留
- 跨包 import 路径走相同的 sed 替换
- 运行 `go test ./...` 验证：所有原 mogo 测试在 moproto 编译并通过
- 验证手段：
  - `go build ./...` 零错误
  - `go vet ./...` 零警告
  - `go test ./...` 全过
  - `go test -race ./...` 无竞态告警

## 9. 验证标准

迁移完成的判定（全部满足）：

1. `go build ./...` 在 moproto 内成功
2. `go vet ./...` 零警告
3. `go test ./...` 全部通过
4. `go test -race ./...` 无竞态告警
5. mogo 仓库工作树无任何改动（`git status` clean）
6. moproto 新文件清单覆盖 4.1 节的全部文件
7. `grep -r 'hegq.net/mogo' /Users/puras/workspace/proj/hegq/moproto` 零命中

## 10. 后续（非本设计范围）

下列事项**不在本次迁移范围**，留作未来工作：

- 是否将 MoProto 拆为独立 module / 独立 git 仓库发布
- 与 mogo 的版本同步策略（定期 diff 与 rebase）
- `pkg/adapter/` 真正实现跨协议转换
- 增加新的 provider（如 Mistral、Cohere）
- 给 MoProto 加 `examples/` 目录与使用文档

## 附录 A：盘点明细（2026-08-28 探索阶段记录）

### A.1 mogo PKG 内 AI 协议相关文件统计

| 子包 | 非测试 .go | 总行数（含测试） |
|---|---|---|
| `pkg/protocol/openai/` | 22 | — |
| `pkg/protocol/anthropic/` | 23 | — |
| `pkg/protocol/gemini/` | 15 | — |
| `pkg/protocol/{openai,anthropic,gemini}/frames/` | 各 1-4 | — |
| `pkg/protocol/` 合计（含 capability.go） | 约 65 | 约 20,324 |
| `pkg/usage/` | 4 | 约 1,623（含 953 行测试） |
| `pkg/adapter/` | 2 | 62（全部 stub） |
| `pkg/translate/` | 1 | 31（全部 stub） |
| `pkg/types/` | 1 | 9 |
| `pkg/parser/error_parser/` | 1 | 236 |
| **复制范围合计** | **约 65** | **约 22,000** |

### A.2 三方依赖检查

`pkg/protocol/` 与 `pkg/usage/` 内非测试代码的三方 import：

- `github.com/google/uuid` — 实际未在 protocol/usage 内非测试代码直接使用
- `github.com/stretchr/testify` — 仅测试

最终 go.mod 仅需 `testify`（测试用）；`go mod tidy` 后会按需补全 indirect 依赖。
