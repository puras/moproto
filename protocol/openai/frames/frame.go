package frames

import "strings"

// NormalizeSSEFrame 把上游透传的原始字节规整成 OpenAI Chat Completions SSE 帧格式
// "data: ...\n\n"。
//
// OpenAI Chat Completions 流协议要求每条消息帧符合：
//
//	data: <json_payload>\n
//	\n
//
// 上游 Provider 在透传模式下发的字节可能有 4 种形态：
//  1. 完整帧："data: {json}\n\n"  → 原样写出
//  2. 仅 data 行："data: {json}\n" → 补一个 \n
//  3. 仅 json："{json}"            → 包成 "data: {json}\n\n"
//  4. 裸 JSON + 单换行："{json}\n" → 包成 "data: {json}\n\n"
//
// 规整策略：
//   - 缺 "data: " 前缀则补
//   - 末尾统一以 \n\n 结尾（SSE 帧终止符）
//   - 空行/纯空白行原样保留作为 SSE 事件分隔符
func NormalizeSSEFrame(raw string) string {
	if raw == "" {
		return raw
	}
	trimmed := strings.TrimRight(raw, " \t")
	// 空行（纯 \n 或空白）：保留作为 SSE 事件分隔符
	if strings.TrimSpace(trimmed) == "" {
		return "\n"
	}

	const prefix = "data: "
	body := trimmed
	if !strings.HasPrefix(body, prefix) {
		body = prefix + body
	}
	// 末尾统一两个换行（SSE 帧终止）
	if !strings.HasSuffix(body, "\n\n") {
		if strings.HasSuffix(body, "\n") {
			body += "\n"
		} else {
			body += "\n\n"
		}
	}
	return body
}
