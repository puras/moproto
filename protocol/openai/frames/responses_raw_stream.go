// Package frames 提供 OpenAI Responses SSE 流帧的 raw-bytes carrier 与工具。
//
// 设计动机：typed union（53 个 variant）在 `bytes → struct → bytes` 往返中会丢字段：
//   - LOSS #1：上游未知字段在 json.Unmarshal 时被丢弃
//   - LOSS #2：omitempty 零值字段在 json.Marshal 时被剥离
//
// 因此同协议（OpenAI Responses → OpenAI Responses）流式路径改走字节透传：
// 上游整帧 bytes 原样存入 OpenAIResponsesRawStreamFrame.Raw，WriteStream 直接 emit。
// 协议层（如 internal/app/gateway/protocol）在需要时按 Type 字段自行 typed dispatch。
//
// 与父包 openai 的关系：
//   - 本包刻意不引入 sealed `OpenAIResponseStreamEventUnion` 接口（53 变体 typed dispatch
//     留给父包 openai 与 protocol 层），避免 frames → openai 单向依赖演变为循环。
//   - carrier 私有 marker `isOpenAIResponsesRawStreamFrame()` 同样用于密封。
//   - 调用方必须用独立的 `case *OpenAIResponsesRawStreamFrame:` 分支显式处理，
//     避免误落入 typed 路径导致重新 marshal。
package frames

import (
	"encoding/json"
	"strings"
)

// OpenAIResponsesRawStreamFrame 同协议 SSE 透传载体，保留上游 event/data 原始字节。
//
// 与 Anthropic 同模式（参考 pkg/protocol/anthropic/frames 包），区别在于 OpenAI
// Responses 的 SSE 协议同时支持：
//
//	event: response.created\n
//	data: {...}\n
//	\n
//
// （Anthropic 形态：双行 + 强制 event: 行）
// 与 OpenAI Chat Completions 的纯 `data: {...}\n\n` 单行形态也不同——
// Responses 引入 event 字段区分事件类型（created/in_progress/completed/output_text.delta 等）。
//
// carrier 字段：
//   - Type:    从 data JSON 的 type 字段抽取；未知类型时为上游原值（透传时保留）
//   - Raw:     上游整帧原始字节（`event: X\ndata: Y\n\n` 完整帧，含终止空行）
//   - Payload: 仅 data 行的 JSON payload（不含 data: 前缀与换行）；供 typed dispatch 旁路使用
type OpenAIResponsesRawStreamFrame struct {
	// Type SSE 事件类型（"response.created" / "response.completed" / ...），
	// 从 data JSON 的 type 字段抽取；未知类型或非 JSON 时为上游原值。
	Type string
	// Raw 上游整帧原始字节（`event: X\ndata: Y\n\n` 格式，含终止空行）。
	Raw []byte
	// Payload data 行去前缀与换行后的 JSON 字节（用于 typed 旁路）。
	Payload []byte
}

// isOpenAIResponsesRawStreamFrame 私有标记方法，用于区别于 typed union。
//
// 刻意不实现父包 sealed 的 isResponseStreamEvent()：本类型不得满足
// OpenAIResponseStreamEventUnion 接口，强制调用方用独立分支显式分流。
func (OpenAIResponsesRawStreamFrame) isOpenAIResponsesRawStreamFrame() {}

// MarshalSSEFrameRaw 直接 emit 上游原始 SSE 帧字节。
//
// 与 marshal typed event 的区别：本函数不做任何 json.Unmarshal / json.Marshal 往返，
// 因此上游的 omitempty 零值字段与未知字段字节保真。
//
// caller 传入从 "event:" 到 "\n\n" 的完整 raw bytes
// （即 OpenAIResponsesRawStreamFrame.Raw）。空输入返回 nil，由 caller 跳过写出。
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

// ParseResponsesSSEEvent 解析 SSE 单行 + 还原后的完整 rawFrame 为 carrier。
//
// 参数：
//   - line:     SSE 单行（典型 "data: {...}" 形态，data 前缀可选）
//   - rawFrame: 还原后的完整 SSE 帧字节（`event: X\ndata: Y\n\n`），存入 carrier.Raw 供 WriteStream 字节透传
//
// 返回：
//   - frame: 解析后的 carrier
//   - ok:    是否成功解析出有效事件；false 表示应跳过该行（空行 / 无 data 前缀 / [DONE] 哨兵）
//
// 行为：
//   - 无 "data: " 前缀 → ok=false（不报错）
//   - data body 为空 / "[DONE]" → ok=false
//   - data body 是合法 JSON → Type 从 type 字段抽取，Payload 保留原文
//   - data body 非 JSON → Type 为空，Payload 仍保留（供 caller 透传），ok=true
//
// 刻意不做 typed dispatch（避免 frames → openai 反向依赖）；caller 在需要时按
// Payload 自行调用父包 openai 的 typed dispatch 逻辑。
func ParseResponsesSSEEvent(line string, rawFrame []byte) (frame OpenAIResponsesRawStreamFrame, ok bool) {
	if !strings.HasPrefix(line, "data: ") {
		return OpenAIResponsesRawStreamFrame{}, false
	}
	body := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
	if body == "" || body == "[DONE]" {
		return OpenAIResponsesRawStreamFrame{}, false
	}

	// 轻量探测 type 字段；失败时 Payload 仍保留
	var typeOnly struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(body), &typeOnly); err != nil {
		return OpenAIResponsesRawStreamFrame{
			Raw:     rawFrame,
			Payload: []byte(body),
		}, true
	}

	return OpenAIResponsesRawStreamFrame{
		Type:    typeOnly.Type,
		Raw:     rawFrame,
		Payload: []byte(body),
	}, true
}
