package anthropic

import (
	"encoding/json"
	"strings"
)

// ============== 类型化 Type 字段 ==============

type AnthropicMessageStartEventType string

const AnthropicMessageStartEventTypeMessageStart AnthropicMessageStartEventType = "message_start"

type AnthropicMessageDeltaEventType string

const AnthropicMessageDeltaEventTypeMessageDelta AnthropicMessageDeltaEventType = "message_delta"

type AnthropicMessageStopEventType string

const AnthropicMessageStopEventTypeMessageStop AnthropicMessageStopEventType = "message_stop"

type AnthropicContentBlockStartEventType string

const AnthropicContentBlockStartEventTypeContentBlockStart AnthropicContentBlockStartEventType = "content_block_start"

type AnthropicContentBlockDeltaEventType string

const AnthropicContentBlockDeltaEventTypeContentBlockDelta AnthropicContentBlockDeltaEventType = "content_block_delta"

type AnthropicContentBlockStopEventType string

const AnthropicContentBlockStopEventTypeContentBlockStop AnthropicContentBlockStopEventType = "content_block_stop"

type AnthropicErrorStreamEventType string

const AnthropicErrorStreamEventTypeError AnthropicErrorStreamEventType = "error"

type AnthropicPingStreamEventType string

const AnthropicPingStreamEventTypePing AnthropicPingStreamEventType = "ping"

// ============== AnthropicStreamEvent union ==============

// AnthropicStreamEvent SSE 事件 union。
//
// 6 个核心事件来自官方 API：
//
//	message_start / message_delta / message_stop
//	content_block_start / content_block_delta / content_block_stop
//
// 2 个扩展事件（error / ping）由网关审计需要补齐。
//
// 通过 isStreamEvent() 私有方法封闭实现集；调用方通过 type switch 或 .Type 分发。
type AnthropicStreamEvent interface {
	isStreamEvent()
}

// AnthropicContentBlockStartEventContentBlock 是 content_block_start.content_block 的类型。
//
// 使用扁平 struct AnthropicContentBlockUnion 替代 sealed interface，
// 解决 json.Unmarshal 无法反序列化到 interface 字段的问题。
type AnthropicContentBlockStartEventContentBlock = AnthropicContentBlockUnion

// ============== 8 个 event struct ==============

// AnthropicMessageStartEvent message_start 事件。
type AnthropicMessageStartEvent struct {
	Type    AnthropicMessageStartEventType `json:"type"` // 固定 "message_start"
	Message *AnthropicMessageResponse      `json:"message"`
}

func (AnthropicMessageStartEvent) isStreamEvent() {}

// AnthropicMessageDeltaEvent message_delta 事件。
type AnthropicMessageDeltaEvent struct {
	Type  AnthropicMessageDeltaEventType `json:"type"` // 固定 "message_delta"
	Delta *MessageDeltaEventDelta        `json:"delta,omitempty"`
	Usage *MessageDeltaUsage             `json:"usage,omitempty"`
}

func (AnthropicMessageDeltaEvent) isStreamEvent() {}

// AnthropicMessageStopEvent message_stop 事件。
type AnthropicMessageStopEvent struct {
	Type AnthropicMessageStopEventType `json:"type"` // 固定 "message_stop"
}

func (AnthropicMessageStopEvent) isStreamEvent() {}

// AnthropicContentBlockStartEvent content_block_start 事件。
type AnthropicContentBlockStartEvent struct {
	Type         AnthropicContentBlockStartEventType `json:"type"` // 固定 "content_block_start"
	Index        int                                 `json:"index"`
	ContentBlock *AnthropicContentBlockUnion         `json:"content_block"`
}

func (AnthropicContentBlockStartEvent) isStreamEvent() {}

// AnthropicContentBlockDeltaEvent content_block_delta 事件。
type AnthropicContentBlockDeltaEvent struct {
	Type  AnthropicContentBlockDeltaEventType `json:"type"` // 固定 "content_block_delta"
	Index int                                 `json:"index"`
	Delta *RawContentBlockDeltaUnion          `json:"delta"`
}

func (AnthropicContentBlockDeltaEvent) isStreamEvent() {}

// AnthropicContentBlockStopEvent content_block_stop 事件。
type AnthropicContentBlockStopEvent struct {
	Type  AnthropicContentBlockStopEventType `json:"type"` // 固定 "content_block_stop"
	Index int                                `json:"index"`
}

func (AnthropicContentBlockStopEvent) isStreamEvent() {}

// AnthropicErrorStreamEvent error 事件（SSE 流式错误，复用 HTTP ErrorResponse 结构）。
type AnthropicErrorStreamEvent struct {
	Type  AnthropicErrorStreamEventType `json:"type"` // 固定 "error"
	Error *AnthropicErrorResponse       `json:"error"`
}

func (AnthropicErrorStreamEvent) isStreamEvent() {}

// AnthropicPingStreamEvent ping 事件（心跳保活）。
type AnthropicPingStreamEvent struct {
	Type AnthropicPingStreamEventType `json:"type"` // 固定 "ping"
}

func (AnthropicPingStreamEvent) isStreamEvent() {}

// ============== raw bytes carrier（同协议 passthrough） ==============

// AnthropicRawStreamFrame 同协议 SSE 透传载体，保留上游 event/data 两行原始字节。
//
// 设计动机：typed union（8 个 variant）在 `bytes → struct → bytes` 往返中会丢字段：
//   - LOSS #1：上游未知字段在 json.Unmarshal 时被丢弃
//   - LOSS #2：omitempty 零值字段（cache_creation_input_tokens:0 /
//     service_tier:"" / stop_reason:null 等）在 json.Marshal 时被剥离
//
// 因此同协议（Anthropic → Anthropic）流式路径改走字节透传：上游整帧 bytes 原样
// 存入 Raw，WriteStream 直接 emit；TypedEvent 仅作旁路，供 usage 抽取 / 状态机分发。
//
// 本类型**刻意不实现** isStreamEvent()，因此不满足 sealed 的 AnthropicStreamEvent
// 接口——调用方必须用独立的 `case *AnthropicRawStreamFrame:` 分支显式处理，
// 避免误落入 typed 路径导致重新 marshal。
//
// typed 路径的 8 个 variant 仍然保留，供 frames.BuildMessageStartFrame 等
// synthetic 帧构造、以及未来跨协议转换使用。
type AnthropicRawStreamFrame struct {
	// Type SSE 事件类型（"message_start" / "message_delta" / ...），
	// 从 data JSON 的 type 字段抽取；未知类型时为上游原值。
	Type string
	// Raw 上游整帧原始字节（`event: X\ndata: Y\n\n` 格式，含终止空行）。
	Raw []byte
	// TypedEvent 同步解析出的 typed 事件，供 usage 抽取与状态机分发；
	// 上游事件类型未知时为 nil。
	TypedEvent AnthropicStreamEvent
}

// isAnthropicRawStreamFrame 私有标记方法，用于区别于 sealed 的 AnthropicStreamEvent。
//
// 刻意与 isStreamEvent() 不同名：本类型不得满足 AnthropicStreamEvent 接口。
func (AnthropicRawStreamFrame) isAnthropicRawStreamFrame() {}

// ParseSSERawFrame 解析 SSE 行为 AnthropicRawStreamFrame（raw bytes + typed 旁路）。
//
// 与 ParseSSEEvent 的区别：本函数额外携带上游整帧原始字节，供同协议字节透传。
//
// 返回零值（Type == "" && Raw == nil && TypedEvent == nil）的场景：
//   - line 无 "data: " 前缀（如 event: 行、空行）
//   - data body 为空或 [DONE] 哨兵（Anthropic 协议无 [DONE]，此处兼容）
//
// TypedEvent 为 nil 表示事件类型未知（ParseSSEEvent 的 default 分支）；
// 调用方应据此决定是否丢帧。
func ParseSSERawFrame(line string, rawFrame []byte) (AnthropicRawStreamFrame, error) {
	if !strings.HasPrefix(line, "data: ") {
		return AnthropicRawStreamFrame{}, nil
	}
	body := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
	if body == "" || body == "[DONE]" {
		return AnthropicRawStreamFrame{}, nil
	}

	// typed 旁路：复用既有 dispatch，未知类型返回 (nil, nil)
	evt, err := ParseSSEEvent(line)
	if err != nil {
		return AnthropicRawStreamFrame{}, err
	}

	var typeOnly struct {
		Type string `json:"type"`
	}
	_ = json.Unmarshal([]byte(body), &typeOnly)

	return AnthropicRawStreamFrame{
		Type:       typeOnly.Type,
		Raw:        rawFrame,
		TypedEvent: evt,
	}, nil
}

// ParseSSEEvent 解析 SSE 行（"data: {...}"）为 AnthropicStreamEvent。
//
// 在 pkg 内部实现，可以返回 sealed interface 的具体实现。
// 供 adapter 层复用，不需要自己定义 SSE 行 struct。
//
// 不处理空行、非 data: 前缀行；由调用方保证 line 格式正确。
func ParseSSEEvent(line string) (AnthropicStreamEvent, error) {
	data := strings.TrimPrefix(line, "data: ")
	data = strings.TrimSpace(data)
	if data == "" {
		return nil, nil
	}

	// 先拿 type 字段判断事件类型
	var typeOnly struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(data), &typeOnly); err != nil {
		return nil, err
	}

	switch typeOnly.Type {
	case string(AnthropicMessageStartEventTypeMessageStart):
		var e AnthropicMessageStartEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicContentBlockStartEventTypeContentBlockStart):
		var e AnthropicContentBlockStartEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicContentBlockDeltaEventTypeContentBlockDelta):
		var e AnthropicContentBlockDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicContentBlockStopEventTypeContentBlockStop):
		var e AnthropicContentBlockStopEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicMessageDeltaEventTypeMessageDelta):
		var e AnthropicMessageDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicMessageStopEventTypeMessageStop):
		var e AnthropicMessageStopEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicErrorStreamEventTypeError):
		var e AnthropicErrorStreamEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	case string(AnthropicPingStreamEventTypePing):
		var e AnthropicPingStreamEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil

	default:
		return nil, nil
	}
}
