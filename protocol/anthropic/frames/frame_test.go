package frames

import (
	"encoding/json"
	"strings"
	"testing"

	"hegq.net/mogo/pkg/protocol/anthropic"
)

// TestMarshalSSEFrame_AllVariants 验证 8 个 stream event 变体都能正确序列化。
func TestMarshalSSEFrame_AllVariants(t *testing.T) {
	cases := []struct {
		name    string
		evt     anthropic.AnthropicStreamEvent
		evtName string
	}{
		{
			name: "message_start",
			evt: anthropic.AnthropicMessageStartEvent{
				Type: anthropic.AnthropicMessageStartEventTypeMessageStart,
				Message: &anthropic.AnthropicMessageResponse{
					ID:    "msg_test",
					Type:  anthropic.AnthropicMessageTypeMessage,
					Role:  anthropic.AnthropicMessageRoleAssistant,
					Model: "claude-test",
				},
			},
			evtName: "message_start",
		},
		{
			name:    "message_delta",
			evt:     anthropic.AnthropicMessageDeltaEvent{Type: anthropic.AnthropicMessageDeltaEventTypeMessageDelta},
			evtName: "message_delta",
		},
		{
			name:    "message_stop",
			evt:     anthropic.AnthropicMessageStopEvent{Type: anthropic.AnthropicMessageStopEventTypeMessageStop},
			evtName: "message_stop",
		},
		{
			name: "content_block_start",
			evt: anthropic.AnthropicContentBlockStartEvent{
				Type:  anthropic.AnthropicContentBlockStartEventTypeContentBlockStart,
				Index: 0,
			},
			evtName: "content_block_start",
		},
		{
			name: "content_block_delta",
			evt: anthropic.AnthropicContentBlockDeltaEvent{
				Type:  anthropic.AnthropicContentBlockDeltaEventTypeContentBlockDelta,
				Index: 0,
				Delta: func() *anthropic.RawContentBlockDeltaUnion {
					b, _ := json.Marshal(anthropic.AnthropicTextDeltaBlock{Type: anthropic.AnthropicTextDeltaBlockTypeTextDelta, Text: "hi"})
					return &anthropic.RawContentBlockDeltaUnion{Type: string(anthropic.AnthropicTextDeltaBlockTypeTextDelta), TextDelta: b}
				}(),
			},
			evtName: "content_block_delta",
		},
		{
			name: "content_block_stop",
			evt: anthropic.AnthropicContentBlockStopEvent{
				Type:  anthropic.AnthropicContentBlockStopEventTypeContentBlockStop,
				Index: 0,
			},
			evtName: "content_block_stop",
		},
		{
			name:    "error",
			evt:     anthropic.AnthropicErrorStreamEvent{Type: anthropic.AnthropicErrorStreamEventTypeError},
			evtName: "error",
		},
		{
			name:    "ping",
			evt:     anthropic.AnthropicPingStreamEvent{Type: anthropic.AnthropicPingStreamEventTypePing},
			evtName: "ping",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := MarshalSSEFrame(tc.evt)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}
			frame := string(data)
			// SSE 帧结构：event: <name>\ndata: <json>\n\n
			if !strings.HasPrefix(frame, "event: "+tc.evtName+"\ndata: ") {
				t.Errorf("frame prefix mismatch, got: %q", frame[:minLen(80, len(frame))])
			}
			if !strings.HasSuffix(frame, "\n\n") {
				t.Errorf("frame must end with \\n\\n, got suffix: %q", frame[maxLen(0, len(frame)-4):])
			}
		})
	}
}

// TestMarshalSSEFrame_NilInterface 验证 nil 接口也能被 type switch 的 default 命中，
// 返回 "unknown" 错误而非 panic。
//
// AnthropicStreamEvent 是 sealed 接口（私有 isStreamEvent()），frames 包无法自定义 mock，
// 这里用 nil 接口做非法入参测试——nil 走 type switch 的 default，与未知类型行为等价。
func TestMarshalSSEFrame_NilInterface(t *testing.T) {
	var nilEvt anthropic.AnthropicStreamEvent
	_, err := MarshalSSEFrame(nilEvt)
	if err == nil {
		t.Fatal("expected error for nil event, got nil")
	}
	if !strings.Contains(err.Error(), "unknown") {
		t.Errorf("error message should mention 'unknown', got: %v", err)
	}
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxLen(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ============== NormalizeSSEFrame 测试 ==============

func TestNormalizeSSEFrame_Empty(t *testing.T) {
	if got := NormalizeSSEFrame(""); got != "" {
		t.Errorf("空字符串应保持空，实际 %q", got)
	}
}

func TestNormalizeSSEFrame_BlankLine(t *testing.T) {
	// SSE 用空行做事件分隔符，必须原样保留
	tests := []string{"\n", "   ", "\t\n", "\n\n"}
	for _, in := range tests {
		got := NormalizeSSEFrame(in)
		if got != "\n" {
			t.Errorf("空行 %q 应原样保留为 \\n，实际 %q", in, got)
		}
	}
}

func TestNormalizeSSEFrame_CompleteFrame(t *testing.T) {
	// 完整 Anthropic 帧原样保留
	in := "event: message_start\ndata: {\"type\":\"message_start\"}\n\n"
	got := NormalizeSSEFrame(in)
	if got != in {
		t.Errorf("完整帧应原样保留，got %q want %q", got, in)
	}
}

func TestNormalizeSSEFrame_MissingTrailingBlank(t *testing.T) {
	// 缺尾部空行 → 补一个 \n
	in := "event: message_start\ndata: {\"type\":\"message_start\"}\n"
	want := "event: message_start\ndata: {\"type\":\"message_start\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("缺尾行应补 \\n\\n，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_DataLineOnly(t *testing.T) {
	// 仅 data 行 → 原样保留
	in := "data: {\"type\":\"content_block_delta\"}\n\n"
	got := NormalizeSSEFrame(in)
	if got != in {
		t.Errorf("data-only 帧应原样保留，got %q want %q", got, in)
	}
}

func TestNormalizeSSEFrame_BareJSON(t *testing.T) {
	// 裸 JSON → 包成 data: ...\n\n
	in := `{"type":"content_block_delta"}`
	want := "data: {\"type\":\"content_block_delta\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("裸 JSON 应包成 data 帧，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_BareJSONWithNewline(t *testing.T) {
	// 裸 JSON + 单换行 → 包成 data: ...\n\n
	in := `{"type":"content_block_delta"}` + "\n"
	want := "data: {\"type\":\"content_block_delta\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("裸 JSON+换行 应包成 data 帧，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_NoEventField(t *testing.T) {
	// Anthropic 协议 event 字段可选，本函数不主动补 event
	// （handler 通常已知上游事件类型）
	in := "data: {\"type\":\"message_start\"}"
	want := "data: {\"type\":\"message_start\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("无 event 字段应只补 \\n\\n，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_TrailingSpaces(t *testing.T) {
	in := "event: ping\ndata: {\"type\":\"ping\"}   "
	want := "event: ping\ndata: {\"type\":\"ping\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("尾部空格应 trim，got %q want %q", got, want)
	}
}

// ============== MarshalSSEFrameRaw 测试（同协议 passthrough） ==============

// TestMarshalSSEFrameRaw_Verbatim 验证 raw 帧字节透传：不做 json.Marshal 往返。
//
// 与 MarshalSSEFrame 关键区别：本函数保留上游 omitempty 零值字段
// （cache_creation_input_tokens:0 / service_tier:"" / stop_reason:null 等）
// 与未知字段字节保真。
func TestMarshalSSEFrameRaw_Verbatim(t *testing.T) {
	cases := []struct {
		name string
		raw  []byte
	}{
		{
			name: "message_delta 含 omitempty 字段",
			raw:  []byte("event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":null,\"stop_sequence\":null,\"stop_details\":{\"refusal_category\":\"\"}},\"usage\":{\"cache_creation_input_tokens\":0,\"cache_read_input_tokens\":0,\"input_tokens\":10,\"output_tokens\":5,\"service_tier\":\"\"}}\n\n"),
		},
		{
			name: "未知事件类型（含未来扩展字段）",
			raw:  []byte("event: future_event\ndata: {\"type\":\"future_event\",\"new_field_4_99\":{\"x\":1},\"unknown_metadata\":\"foo\"}\n\n"),
		},
		{
			name: "嵌套结构与 Unicode",
			raw:  []byte("event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"model\":\"claude-sonnet-4-20250514\",\"content\":[{\"type\":\"text\",\"text\":\"你好\"}]}}\n\n"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := MarshalSSEFrameRaw(tc.raw)
			if string(out) != string(tc.raw) {
				t.Errorf("byte mismatch:\n got: %q\nwant: %q", out, tc.raw)
			}
		})
	}
}

// TestMarshalSSEFrameRaw_EmptyReturnsNil 验证空输入返回 nil（caller 跳过写出）。
func TestMarshalSSEFrameRaw_EmptyReturnsNil(t *testing.T) {
	if got := MarshalSSEFrameRaw(nil); got != nil {
		t.Errorf("nil should return nil, got %q", got)
	}
	if got := MarshalSSEFrameRaw([]byte{}); got != nil {
		t.Errorf("empty slice should return nil, got %q", got)
	}
}
