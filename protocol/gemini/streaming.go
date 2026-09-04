package gemini

import (
	"encoding/json"
	"strings"

	pkgusage "github.com/puras/moproto/usage"
)

// ParseGeminiStreamChunk 解析 Gemini SSE 行为 *GeminiGenerateContentResponse。
//
// Gemini 流式协议特殊：每个 `data:` 行都是一个完整的 GenerateContentResponse（不是
// 事件 union），与 Anthropic/OpenAI 流协议差异显著。SSE 帧之间不需要事件名，
// 仅 `data:` + JSON body + 空行分隔。
//
// 返回 (event, done, error, usage)：
//   - event: *GeminiGenerateContentResponse
//   - done: 流是否结束（Gemini 协议无显式 [DONE] 哨兵，channel 关闭即结束；返回 false）
//   - error: 解析错误
//   - usage: 从此帧 usageMetadata 提取的 *pkgusage.Usage
//
// 参数：
//   - line: SSE data: 行原文（如 `data: {...}`）
//   - rawFrame: 还原后的完整 SSE 帧字节（`data: Y\n\n`）。存到 resp.RawData 用于
//     MarshalJSON 旁路透传；typed unmarshal 失败不再杀帧（RawData 已保真）。
func ParseGeminiStreamChunk(line string, rawFrame []byte) (any, bool, error, *pkgusage.Usage) {
	if !strings.HasPrefix(line, "data: ") {
		return nil, false, nil, nil
	}
	body := strings.TrimPrefix(line, "data: ")
	if strings.TrimSpace(body) == "" {
		return nil, false, nil, nil
	}

	resp := &GeminiGenerateContentResponse{
		RawData: []byte(body),
	}

	// 旁路 typed unmarshal 仅供 usage 抽取与元数据补齐；失败不再杀帧
	// （RawData 已保留原始字节，MarshalJSON 走 passthrough）
	_ = json.Unmarshal(resp.RawData, resp)

	usage := ExtractGeminiUsage(resp.UsageMetadata)
	return resp, false, nil, usage
}

// ExtractGeminiUsage 从 *GeminiUsageMetadata 转 *pkgusage.Usage。
//
// 字段映射（与 pkgusage.Usage 对齐）：
//   - PromptTokenCount      → Base.InputTokens
//   - CandidatesTokenCount  → Base.OutputTokens
//   - TotalTokenCount       → Base.TotalTokens（原样透传，官方已计算好）
//   - CachedContentTokenCount → Cache.CachedTokens（OpenAI cached_tokens 字段）
//   - ThoughtsTokenCount    → Thinking.ThinkingTokens
//
// nil 入参返回 nil（调用方无需 nil 检查）。
func ExtractGeminiUsage(m *GeminiUsageMetadata) *pkgusage.Usage {
	if m == nil {
		return nil
	}

	u := &pkgusage.Usage{
		Base: pkgusage.Base{
			InputTokens:  int64(m.PromptTokenCount),
			OutputTokens: int64(m.CandidatesTokenCount),
			TotalTokens:  int64(m.TotalTokenCount),
		},
		Protocol: "gemini",
	}

	if m.CachedContentTokenCount > 0 {
		u.Cache = &pkgusage.CacheMetrics{
			CachedTokens: int64(m.CachedContentTokenCount),
		}
	}

	if m.ThoughtsTokenCount > 0 {
		u.Thinking = &pkgusage.ExtendedThinking{
			ThinkingTokens: int64(m.ThoughtsTokenCount),
		}
	}

	return u
}
