# MoProto AI 协议解析层迁移实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 把 mogo 后端的 AI 协议解析层（openai / anthropic / gemini 三大 provider + 用量解析 usage）整体迁移到 MoProto，作为独立可 import 的 Go 解析库。

**Architecture:** 单 Go module（`hegq.net/moproto`），顶层布局 `protocol/{openai,anthropic,gemini}` + `usage/`，所有源文件由 mogo 镜像复制，仅替换 import 路径（`hegq.net/mogo/pkg/` → `hegq.net/moproto/`）与移除子目录内的 `go.mod`/`go.sum`；包声明（`package openai` 等）保持不变；测试文件全部随源文件复制。

**Tech Stack:** Go 1.26、`github.com/stretchr/testify`（测试用，必要时追加 `github.com/leanovate/gopter`）。

**Spec:** `docs/superpowers/specs/2026-08-28-moproto-migration-design.md`

## Global Constraints

- **MoGo 工程（`/Users/puras/workspace/proj/106hz/mogo`）只读**：本计划全程不得修改、commit、checkout、restore mogo 内任何文件。所有 `cp`/`sed`/`find`/`go` 命令的目标路径均在 `moproto/`，源路径仅作为读取端参数。
- **零业务改动**：除包名/import 路径外不修改任何 mogo 源文件逻辑；包声明 `package openai/anthropic/gemini/usage` 保持原样；`mcp.go` 等含 mogo 业务注释的源文件保留原注释不动。
- **不复制清单**：`pkg/protocol/capability.go`、`pkg/parser/error_parser/`、`pkg/adapter/`、`pkg/translate/`、`pkg/types/types.go`、`pkg/protocol/video/`、`pkg/protocol/seedance/` 一律不复制。
- **复制清单**：`pkg/protocol/openai/`、`pkg/protocol/anthropic/`、`pkg/protocol/gemini/`（含各自 `frames/` 子目录）、`pkg/usage/` 整目录（含所有 `.go` 与 `_test.go`）。
- **不引入 mogo 模块**：moproto 的 `go.mod` 不允许出现 `hegq.net/...` require；`go list -m all` 应不出现 mogo。
- **commit 频率**：每个 Task 完成后立即 git commit；commit message 用中文描述本次动作。
- **Mogo 工程目录不能放 moproto 的工作目录**：moproto 是 `/Users/puras/workspace/proj/hegq/moproto`，mogo 是 `/Users/puras/workspace/proj/106hz/mogo/backend`，两者完全独立。

---

## 文件结构概览

```
moproto/
├── go.mod                                # Task 1 创建
├── go.sum                                # Task 4 生成
├── README.md                             # Task 5 创建
├── protocol/openai/                      # Task 2 复制（22 个非测试 .go + frames/）
├── protocol/anthropic/                   # Task 2 复制（23 个非测试 .go + frames/）
├── protocol/gemini/                      # Task 2 复制（15 个非测试 .go + frames/）
└── usage/                                # Task 2 复制（4 个非测试 .go + _test.go）
```

Task 1：搭骨架（仅 `go.mod`）→ Task 2：复制四个目录 → Task 3：替换 import 路径 → Task 4：`go mod tidy` + 编译验证 → Task 5：测试运行 + README。每个 Task 是一次独立提交。

---

### Task 1：创建 MoProto 顶层 go.mod

**Files:**
- Create: `/Users/puras/workspace/proj/hegq/moproto/go.mod`

**Interfaces:**
- Produces: moproto 的 Go module 根，所有 `go build`/`go test` 从此目录运行

- [ ] **Step 1：写入 go.mod**

创建文件 `/Users/puras/workspace/proj/hegq/moproto/go.mod`，内容：

```go
module hegq.net/moproto

go 1.26

require (
    github.com/stretchr/testify v1.11.1
)
```

- [ ] **Step 2：验证 go.mod 解析正确**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && go mod download github.com/stretchr/testify
```
Expected: 成功下载（或提示 already added to go.sum，无错误）。

- [ ] **Step 3：确认 mogo 工程未被修改**

Run:
```bash
cd /Users/puras/workspace/proj/106hz/mogo/backend && git status --short
```
Expected: 输出与计划开始前一致（即 mogo 仍处于用户原本的工作状态，本 Task 不应触碰它）。

- [ ] **Step 4：commit**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  git add go.mod go.sum && \
  git -c user.email=claude@moproto.local -c user.name=Claude commit -m "chore(moproto): 初始化 module hegq.net/moproto"
```
Expected: 一次干净提交，工作树恢复为空。

---

### Task 2：复制 mogo 协议目录到 moproto（保留文件结构）

**Files:**
- Create: `/Users/puras/workspace/proj/hegq/moproto/protocol/openai/`（含全部 .go 与 _test.go）
- Create: `/Users/puras/workspace/proj/hegq/moproto/protocol/anthropic/`（含全部 .go 与 _test.go）
- Create: `/Users/puras/workspace/proj/hegq/moproto/protocol/gemini/`（含全部 .go 与 _test.go）
- Create: `/Users/puras/workspace/proj/hegq/moproto/usage/`（含全部 .go 与 _test.go）

**Interfaces:**
- Consumes: mogo 的 `pkg/protocol/openai/`、`pkg/protocol/anthropic/`、`pkg/protocol/gemini/`、`pkg/usage/`（只读源）
- Produces: moproto 内四组同名目录的完整文件副本（尚未修改 import）

- [ ] **Step 1：复制四个目录**

Run:
```bash
SRC=/Users/puras/workspace/proj/106hz/mogo/backend/pkg
DST=/Users/puras/workspace/proj/hegq/moproto

mkdir -p "$DST/protocol"
cp -R "$SRC/protocol/openai"    "$DST/protocol/openai"
cp -R "$SRC/protocol/anthropic" "$DST/protocol/anthropic"
cp -R "$SRC/protocol/gemini"    "$DST/protocol/gemini"
cp -R "$SRC/usage"              "$DST/usage"
```

- [ ] **Step 2：验证复制完整性**

Run:
```bash
DST=/Users/puras/workspace/proj/hegq/moproto
SRC=/Users/puras/workspace/proj/106hz/mogo/backend/pkg
for d in protocol/openai protocol/anthropic protocol/gemini usage; do
  rel=${d#protocol/}
  cnt_dst=$(find "$DST/$d" -name '*.go' | wc -l)
  cnt_src=$(find "$SRC/$d" -name '*.go' | wc -l)
  echo "$d: src=$cnt_src dst=$cnt_dst"
done
```
Expected: 每个目录 src 与 dst 的 `.go` 文件数量相同（用 `wc -l` 取得实际值后比对；如差异非零则说明复制未完成）。

- [ ] **Step 3：删除子目录内被复制过来的 go.mod / go.sum**

mogo 各子包内通常没有独立 go.mod，但 frames 子目录或 vendor 元数据若被复制应清理。为安全起见只扫描 `protocol/` 与 `usage/` 子树，避免误删 moproto 顶层 `go.mod`/`go.sum`：

Run:
```bash
DST=/Users/puras/workspace/proj/hegq/moproto
find "$DST/protocol" "$DST/usage" -name 'go.mod' -delete
find "$DST/protocol" "$DST/usage" -name 'go.sum' -delete
```

- [ ] **Step 4：确认 mogo 仍未被修改**

Run:
```bash
cd /Users/puras/workspace/proj/106hz/mogo/backend && git status --short
```
Expected: 与 Task 1 Step 3 输出完全一致。

- [ ] **Step 5：commit**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  git add protocol/ usage/ && \
  git -c user.email=claude@moproto.local -c user.name=Claude commit -m "feat(moproto): 从 mogo 复制 protocol/openai/anthropic/gemini + usage 目录"
```
Expected: 提交包含 protocol/ 与 usage/ 全部新文件。

---

### Task 3：替换所有 .go 文件中的 mogo import 路径

**Files:**
- Modify: `/Users/puras/workspace/proj/hegq/moproto/protocol/**/*.go`（含 _test.go）
- Modify: `/Users/puras/workspace/proj/hegq/moproto/usage/**/*.go`（含 _test.go）

**Interfaces:**
- Consumes: Task 2 复制的全部 .go 文件
- Produces: 所有文件中 `hegq.net/mogo/pkg/` 字面量替换为 `hegq.net/moproto/`；其它内容（含包声明、注释、`internal/app/...` 注释引用）保持原样

- [ ] **Step 1：扫描当前存在的 mogo 路径引用，确认替换面**

Run:
```bash
DST=/Users/puras/workspace/proj/hegq/moproto
grep -rl 'hegq.net/mogo/pkg/' "$DST" --include='*.go' | tee /tmp/moproto_mogo_imports.txt
echo "---total files---"
wc -l /tmp/moproto_mogo_imports.txt
```
Expected: 多个 .go 文件路径，非空列表。

- [ ] **Step 2：执行 sed 替换**

Run:
```bash
DST=/Users/puras/workspace/proj/hegq/moproto
find "$DST" -name '*.go' -type f -print0 | \
  xargs -0 sed -i '' 's|hegq.net/mogo/pkg/|hegq.net/moproto/|g'
```

- [ ] **Step 3：验证替换结果（必须零残留）**

Run:
```bash
DST=/Users/puras/workspace/proj/hegq/moproto
hits=$(grep -r 'hegq.net/mogo' "$DST" --include='*.go' | wc -l)
echo "remaining hegq.net/mogo references in .go files: $hits"
```
Expected: 输出 0。如果非零，列出具体文件并回滚 Task 3 重做。

- [ ] **Step 4：抽查关键 import 替换正确**

Run:
```bash
DST=/Users/puras/workspace/proj/hegq/moproto
grep -rh 'hegq.net/moproto' "$DST" --include='*.go' | sort -u
```
Expected: 至少出现 `hegq.net/moproto/protocol/openai`、`hegq.net/moproto/protocol/anthropic`、`hegq.net/moproto/protocol/gemini`、`hegq.net/moproto/usage` 四个 import 路径。

- [ ] **Step 5：确认 mogo 工程依然只读**

Run:
```bash
cd /Users/puras/workspace/proj/106hz/mogo/backend && git status --short
```
Expected: 与 Task 2 Step 4 输出完全一致（mogo 无任何变化）。

- [ ] **Step 6：commit**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  git add -A && \
  git -c user.email=claude@moproto.local -c user.name=Claude commit -m "refactor(moproto): 将 import 路径由 hegq.net/mogo/pkg/* 改为 hegq.net/moproto/*"
```
Expected: 大量文件被本次提交纳入。

---

### Task 4：go mod tidy + 编译 + vet 验证

**Files:**
- Modify: `/Users/puras/workspace/proj/hegq/moproto/go.mod`（tidy 补全 indirect）
- Create: `/Users/puras/workspace/proj/hegq/moproto/go.sum`（tidy 生成）

**Interfaces:**
- Consumes: Task 3 后的全部 .go 文件
- Produces: 一个能 `go build ./...` 成功的 moproto module；`go vet ./...` 零警告；`go.mod`/`go.sum` 含全部实际使用的依赖

- [ ] **Step 1：运行 go mod tidy**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && go mod tidy
```
Expected: 退出码 0；如果失败，记录错误信息并继续。

- [ ] **Step 2：检查 go.mod / go.sum 未引入 mogo**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  grep -E 'hegq\.net/mogo' go.mod go.sum && echo "FAIL: mogo leak" || echo "OK: no mogo dep"
```
Expected: 输出 `OK: no mogo dep`。

- [ ] **Step 3：go build 全部包**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && go build ./...
```
Expected: 退出码 0，零错误。如有错误，记录第一个错误的文件:行号。

- [ ] **Step 4：go vet 全部包**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && go vet ./...
```
Expected: mogo 原代码继承的"repeats json tag"警告可作为遗留接受，不视为 fail（spec §1 零业务改动优先于 §9 零警告）。**迁移不得引入新的 vet 警告**——若仅是 mogo 原代码已存在的重复 json tag，迁移仍视为通过。

- [ ] **Step 5：commit（go.mod / go.sum 可能被 tidy 改动）**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  git add go.mod go.sum && \
  git -c user.email=claude@moproto.local -c user.name=Claude commit -m "chore(moproto): go mod tidy 补全依赖"
```
Expected: 干净提交。

---

### Task 5：运行全部测试 + 写入 README

**Files:**
- Modify: `/Users/puras/workspace/proj/hegq/moproto/README.md`（创建）
- Read-only: 全部 moproto 测试文件

**Interfaces:**
- Consumes: Task 4 后的可编译 moproto
- Produces: `go test ./...` 全过；`go test -race ./...` 无竞态告警；README.md 说明库用途、来源、用法

- [ ] **Step 1：运行全部测试（普通模式）**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && go test ./...
```
Expected: 所有测试通过；如有失败，记录失败的包与测试名。

- [ ] **Step 2：运行 race detector**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && go test -race ./...
```
Expected: 退出码 0，无 DATA RACE 报告。

- [ ] **Step 3：写入 README.md**

创建 `/Users/puras/workspace/proj/hegq/moproto/README.md`，内容：

```markdown
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

- `github.com/stretchr/testify` — 测试用
- 必要时追加 `github.com/leanovate/gopter` — property-based 测试

## 版本同步策略

本项目作为 mogo 协议层的独立镜像发布。当 mogo 内协议结构有更新时，
按以下步骤同步：

1. 在 mogo 内执行目标 provider 单元测试，记录基线
2. 从 mogo 复制新版本目录到 moproto 对应路径
3. 重跑 Task 3 的 sed 替换
4. 在 moproto 内执行 `go mod tidy && go test ./...`
5. 对比 mogo 与 moproto 测试结果应一致
```

- [ ] **Step 4：最终验证清单**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  echo "=== mogo working tree ===" && \
  (cd /Users/puras/workspace/proj/106hz/mogo/backend && git status --short | head) && \
  echo "=== moproto no mogo leak ===" && \
  (grep -r 'hegq.net/mogo' . --include='*.go' --include='go.mod' --include='go.sum' | head || echo "OK: zero hits") && \
  echo "=== build ===" && \
  go build ./... && echo "build OK" && \
  echo "=== vet ===" && \
  go vet ./... && echo "vet OK" && \
  echo "=== test ===" && \
  go test ./... 2>&1 | tail -20
```
Expected:
- mogo working tree 输出与计划开始前一致
- moproto 中 mogo 字面量零命中
- `build OK`、`vet OK`、测试汇总显示 `PASS`

- [ ] **Step 5：commit README**

Run:
```bash
cd /Users/puras/workspace/proj/hegq/moproto && \
  git add README.md && \
  git -c user.email=claude@moproto.local -c user.name=Claude commit -m "docs(moproto): README 说明来源/范围/用法/同步策略"
```
Expected: 一次干净提交。

---

## 完成标志

完成所有 5 个 Task 后，迁移成功。判定条件（与 spec §9 一致，§9 中"go vet 零警告"由用户 2026-08-30 裁决让步为"go vet 无新引入警告"）：

1. ✅ `go build ./...` 零错误
2. ✅ `go vet ./...` 退出码可非 0，但**不得有 mogo 未引入的新警告**（mogo 原代码继承的"repeats json tag"接受为遗留）
3. ✅ `go test ./...` 全部通过
4. ✅ `go test -race ./...` 无竞态告警
5. ✅ mogo 仓库工作树无任何改动（git status 与计划开始前一致）
6. ✅ `grep -r 'hegq.net/mogo' /Users/puras/workspace/proj/hegq/moproto` 零命中
