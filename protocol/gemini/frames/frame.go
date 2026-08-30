// Package frames 提供 Gemini SSE 流帧的构造与归一化助手。
//
// 与 pkg/protocol/gemini 严格分离：父包只声明 wire-format 类型，本包只声明操作。
// 零外部依赖，纯函数，可被任何 caller 安全复用。
//
// Gemini 协议 SSE 帧特殊：仅 data 行，无 event 字段（每帧都是一个完整的
// GenerateContentResponse）。本包 Marshaller 直接生成 "data: <json>\n\n" 格式。
package frames

import (
	"encoding/json"
	"fmt"
	"strings"

	"hegq.net/mogo/pkg/protocol/gemini"
)

// MarshalSSEFrame 把 *gemini.GeminiGenerateContentResponse 序列化为 SSE 帧字节。
//
// 帧格式：
//
//	data: <json_payload>\n
//	\n
//
// Gemini 协议无 event 字段。
//
// Passthrough 旁路：resp.RawData 非空时直接 emit 原始字节（"data: <raw>\n\n"），
// 避免 json.Marshal 往返剥离 omitempty 零值字段与未知字段
// （groundingMetadata / urlContextMetadata / avgLogprobs 等）。
func MarshalSSEFrame(value any) ([]byte, error) {
	resp, ok := value.(*gemini.GeminiGenerateContentResponse)
	if !ok || resp == nil {
		return nil, fmt.Errorf("gemini: expected non-nil response")
	}
	if len(resp.RawData) > 0 {
		return []byte(fmt.Sprintf("data: %s\n\n", resp.RawData)), nil
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("gemini: marshal SSE frame: %w", err)
	}
	return []byte(fmt.Sprintf("data: %s\n\n", data)), nil
}

// NormalizeSSEFrame 把上游透传的字节规整成标准 SSE 帧格式。
//
// 上游可能发的字节形态：
//  1. 完整帧："data: {...}\n\n"      （最常见）
//  2. 缺尾部空行："data: {...}\n"
//  3. 裸 JSON："{...}"
//
// 规整策略：缺 data 前缀则补；末尾统一以 \n\n 结尾；空行/纯空白行原样保留。
func NormalizeSSEFrame(raw string) string {
	if raw == "" {
		return raw
	}
	trimmed := strings.TrimRight(raw, " \t")
	if strings.TrimSpace(trimmed) == "" {
		return "\n"
	}

	if !strings.HasPrefix(trimmed, "data: ") {
		body := trimmed
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		body = "data: " + body
		if !strings.HasSuffix(body, "\n\n") {
			body += "\n"
		}
		return body
	}

	if strings.HasSuffix(trimmed, "\n\n") {
		return trimmed
	}
	if strings.HasSuffix(trimmed, "\n") {
		return trimmed + "\n"
	}
	return trimmed + "\n\n"
}
