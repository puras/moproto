// Package frames 提供 Anthropic SSE 流帧的构造与序列化助手。
//
// 本包是 pkg/protocol/anthropic 的姊妹包，封装了对 AnthropicStreamEvent /
// AnthropicUsage 等 wire-format 类型的"按 SSE 协议拼装字节"操作。
//
// 设计原则：
//   - 与父包 anthropic 严格分离：父包只声明 wire-format 类型，本包只声明操作。
//   - 对父包类型零依赖添加：未来若 OpenAI / Gemini 也需要 SSE 帧工具，
//     可以平行建立 frames/openai/、frames/gemini/。
//   - 纯函数，无 HTTP/gin 依赖，可被任何 caller 安全复用。
package frames

import (
	"encoding/json"
	"fmt"
	"strings"

	"hegq.net/moproto/protocol/anthropic"
)

// MarshalSSEFrame 把一个 AnthropicStreamEvent 序列化为完整 SSE 帧字节。
//
// 帧格式为 Anthropic SSE 协议标准：
//
//	event: <event_name>\n
//	data: <json_payload>\n
//	\n
//
// 使用顺序：
//
//	data, err := frames.MarshalSSEFrame(evt)
//	if err != nil { return err }
//	sw.Write(string(data))
//
// 事件名（event: 后面的取值）通过 type switch + sseEventName 派生，保证与 struct 的
// Type 字段保持一致（任何一处新增事件都要同时改两个地方，避免事实源分散）。
//
// 注意：本函数返回的字节不含 http.Flusher 写入动作；调用方负责 flush。
func MarshalSSEFrame(e anthropic.AnthropicStreamEvent) ([]byte, error) {
	name := sseEventName(e)
	if name == "" {
		return nil, fmt.Errorf("anthropic: unknown AnthropicStreamEvent type %T", e)
	}
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("anthropic: marshal SSE frame: %w", err)
	}
	return []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", name, data)), nil
}

// MarshalSSEFrameRaw 直接 emit 上游原始 SSE 帧字节（双行格式）。
//
// 与 MarshalSSEFrame 的区别：本函数**不做任何 json.Unmarshal / json.Marshal 往返**，
// 因此上游的 omitempty 零值字段（cache_creation_input_tokens:0 / service_tier:"" /
// stop_reason:null 等）与未知字段字节保真。
//
// caller 传入从 "event:" 到 "\n\n" 的完整 raw bytes（即
// anthropic.AnthropicRawStreamFrame.Raw）。空输入返回 nil，由 caller 跳过写出。
//
// 使用顺序：
//
//	if out := frames.MarshalSSEFrameRaw(rf.Raw); out != nil {
//	    sw.Write(string(out))
//	}
func MarshalSSEFrameRaw(rawFrame []byte) []byte {
	if len(rawFrame) == 0 {
		return nil
	}
	return rawFrame
}

// sseEventName 返回 stream event 在 SSE 协议中对应的 event 字段值。
//
// 8 个变体（与父包 stream_event.go 的 8 个 struct 一一对应）：
//
//	message_start / message_delta / message_stop
//	content_block_start / content_block_delta / content_block_stop
//	error / ping
//
// 调用方应有默认分支处理未知类型（本函数对未知类型返回 ""，由 MarshalSSEFrame 转错误）。
func sseEventName(e anthropic.AnthropicStreamEvent) string {
	switch e.(type) {
	case anthropic.AnthropicMessageStartEvent:
		return string(anthropic.AnthropicMessageStartEventTypeMessageStart)
	case anthropic.AnthropicMessageDeltaEvent:
		return string(anthropic.AnthropicMessageDeltaEventTypeMessageDelta)
	case anthropic.AnthropicMessageStopEvent:
		return string(anthropic.AnthropicMessageStopEventTypeMessageStop)
	case anthropic.AnthropicContentBlockStartEvent:
		return string(anthropic.AnthropicContentBlockStartEventTypeContentBlockStart)
	case anthropic.AnthropicContentBlockDeltaEvent:
		return string(anthropic.AnthropicContentBlockDeltaEventTypeContentBlockDelta)
	case anthropic.AnthropicContentBlockStopEvent:
		return string(anthropic.AnthropicContentBlockStopEventTypeContentBlockStop)
	case anthropic.AnthropicErrorStreamEvent:
		return string(anthropic.AnthropicErrorStreamEventTypeError)
	case anthropic.AnthropicPingStreamEvent:
		return string(anthropic.AnthropicPingStreamEventTypePing)
	default:
		return ""
	}
}

// NormalizeSSEFrame 把上游透传的原始字节规整成 Anthropic Messages API SSE 帧格式。
//
// Anthropic SSE 帧标准（与 OpenAI 区别在于可选 event 字段）：
//
//	event: <event_name>\n
//	data: <json_payload>\n
//	\n
//
// 上游 Provider 在透传模式下发的字节可能有 4 种形态：
//  1. 完整帧："event: xxx\ndata: {...}\n\n"     （最常见）
//  2. 缺尾部空行："event: xxx\ndata: {...}\n"
//  3. 仅 data 行："data: {...}\n\n"
//  4. 裸 JSON："{...}"
//
// 规整策略：
//   - 缺 data 前缀则补；缺 event 字段不补（Anthropic 协议 event 可选）
//   - 末尾统一以 \n\n 结尾（SSE 帧终止符）
//   - 空行/纯空白行原样保留作为 SSE 事件分隔符
//
// 设计选择：本函数不主动补 event 字段。从裸 JSON 推 event 类型需要 unmarshal + 读
// type 字段（昂贵），而调用方（handler）通常已知上游事件类型。event 字段缺失的边界
// 场景留待调用方按需处理。
func NormalizeSSEFrame(raw string) string {
	if raw == "" {
		return raw
	}
	trimmed := strings.TrimRight(raw, " \t")
	// 空行（纯 \n 或空白）：保留作为 SSE 事件分隔符
	if strings.TrimSpace(trimmed) == "" {
		return "\n"
	}

	hasData := strings.HasPrefix(trimmed, "data: ") || strings.Contains(trimmed, "\ndata: ")

	// 补 data 前缀（缺前缀时把整段当 data body）
	if !hasData {
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

	// 已有 data：只补尾部 \n\n
	if strings.HasSuffix(trimmed, "\n\n") {
		return trimmed
	}
	if strings.HasSuffix(trimmed, "\n") {
		return trimmed + "\n"
	}
	return trimmed + "\n\n"
}
