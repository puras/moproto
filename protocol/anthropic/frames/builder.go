package frames

import (
	"fmt"

	"github.com/puras/moproto/protocol/anthropic"
)

// MessageStartOptions 构造 message_start 帧所需参数。
//
// MessageStart 是流式响应的首事件，必须携带 message.usage（Anthropic 协议要求）；
// 缺 usage 时下游 SDK（如 claude-code）会抛 "cannot read property 'input_tokens' of undefined"。
type MessageStartOptions struct {
	// Model 下游请求的模型名，写入 message.model（必填）。
	Model string

	// Usage 必填；InputTokens/OutputTokens 必填 ≥ 0；CacheCreation*/CacheRead*
	// 可选，omitempty 约定：>0 时填充，0 时省略，避免多余零字段触发下游 SDK 解析异常。
	Usage *anthropic.AnthropicUsage

	// MessageID 可选；不传时调用 GenerateMessageID() 自动生成。
	MessageID string
}

// MessageDeltaOptions 构造 message_delta 帧所需参数。
//
// message_delta 通常是流的最后非停止事件，携带 delta.stop_reason 与本轮累计 usage。
type MessageDeltaOptions struct {
	// StopReason 可选；空值时省略字段（Anthropic 协议容错约定）。
	// 已知有效值：AnthropicStopReasonEndTurn / AnthropicStopReasonMaxTokens /
	// AnthropicStopReasonStopSequence / AnthropicStopReasonToolUse / AnthropicStopReasonRefusal。
	StopReason anthropic.AnthropicStopReason

	// Usage 可选；与 message_start 用 MessageDeltaUsage 类型承载。
	// 0 值字段会被 omitempty 省略。
	Usage *anthropic.MessageDeltaUsage
}

// BuildMessageStartFrame 生成 message_start 的完整 SSE 帧字节。
//
// 错误：opts.Usage 为 nil → 返回 error；opts.Model 为空 → 返回 error（协议要求 model 必填）。
//
// 与父包 anthropic.AnthropicMessageResponse 的 omitempty 约定一致：message_start 阶段的
// stop_reason/stop_sequence 字段不主动填，调用方对最终 wire-format 不必强求 1:1 等价
// Anthropic 官方 SDK 行为。
func BuildMessageStartFrame(opts MessageStartOptions) ([]byte, error) {
	if opts.Usage == nil {
		return nil, fmt.Errorf("anthropic: BuildMessageStartFrame: Usage is required")
	}
	if opts.Model == "" {
		return nil, fmt.Errorf("anthropic: BuildMessageStartFrame: Model is required")
	}

	msgID := opts.MessageID
	if msgID == "" {
		msgID = GenerateMessageID()
	}

	return MarshalSSEFrame(anthropic.AnthropicMessageStartEvent{
		Type: anthropic.AnthropicMessageStartEventTypeMessageStart,
		Message: &anthropic.AnthropicMessageResponse{
			ID:      msgID,
			Type:    anthropic.AnthropicMessageTypeMessage,
			Role:    anthropic.AnthropicMessageRoleAssistant,
			Model:   opts.Model,
			Content: []*anthropic.AnthropicContentBlockUnion{}, // 与原 handler 行为一致：输出 "content": []
			Usage:   opts.Usage,
		},
	})
}

// BuildMessageDeltaFrame 生成 message_delta 的完整 SSE 帧字节。
//
// StopReason 为空时不输出 delta.stop_reason 字段（omitempty），
// 这是与 Anthropic 官方 SDK 行为一致的选择。
func BuildMessageDeltaFrame(opts MessageDeltaOptions) ([]byte, error) {
	delta := &anthropic.MessageDeltaEventDelta{}
	if opts.StopReason != "" {
		sr := opts.StopReason
		delta.StopReason = &sr
	}

	usage := opts.Usage
	if usage == nil {
		usage = &anthropic.MessageDeltaUsage{}
	}

	return MarshalSSEFrame(anthropic.AnthropicMessageDeltaEvent{
		Type:  anthropic.AnthropicMessageDeltaEventTypeMessageDelta,
		Delta: delta,
		Usage: usage,
	})
}

// BuildMessageStopFrame 生成 message_stop 的完整 SSE 帧字节。
//
// message_stop 是流的终止事件，结构固定（仅 {"type":"message_stop"}）。
// 无填充逻辑，调用即返回；不会返回 error（保证调用方零检查）。
func BuildMessageStopFrame() []byte {
	frame, _ := MarshalSSEFrame(anthropic.AnthropicMessageStopEvent{
		Type: anthropic.AnthropicMessageStopEventTypeMessageStop,
	})
	return frame
}
