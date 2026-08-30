package anthropic

import (
	"bytes"
	"testing"
)

// TestParseSSERawFrame_PreservesRawBytes 验证 ParseSSERawFrame 把上游整帧 raw
// bytes 装载到 carrier.Raw,且 Type 字段从 data JSON 的 type 抽取。
//
// 镜像 OpenAI TestParseSSEEvent_PreservesRawBytes 模式 — 保证 passthrough 路径下
// omitempty 零值字段（cache_creation_input_tokens:0 / service_tier:"" /
// stop_reason:null 等）与未知字段字节保真。
func TestParseSSERawFrame_PreservesRawBytes(t *testing.T) {
	cases := []struct {
		name      string
		rawFrame  []byte
		dataLine  string
		wantType  string
		wantRawOK bool
	}{
		{
			name: "message_start 含 cache_creation_input_tokens:0",
			rawFrame: []byte(`event: message_start
data: {"type":"message_start","message":{"id":"msg_test","model":"claude-x","usage":{"cache_creation_input_tokens":0,"cache_read_input_tokens":0,"input_tokens":10,"output_tokens":0,"service_tier":""}}}

`),
			dataLine:  `data: {"type":"message_start","message":{"id":"msg_test","model":"claude-x","usage":{"cache_creation_input_tokens":0,"cache_read_input_tokens":0,"input_tokens":10,"output_tokens":0,"service_tier":""}}}`,
			wantType:  "message_start",
			wantRawOK: true,
		},
		{
			name: "message_delta 含 stop_reason:null",
			rawFrame: []byte(`event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":null,"stop_sequence":null},"usage":{"input_tokens":0,"output_tokens":5}}

`),
			dataLine:  `data: {"type":"message_delta","delta":{"stop_reason":null,"stop_sequence":null},"usage":{"input_tokens":0,"output_tokens":5}}`,
			wantType:  "message_delta",
			wantRawOK: true,
		},
		{
			name: "未知事件类型仍保留 raw bytes",
			rawFrame: []byte(`event: future_event
data: {"type":"future_event","unknown_field":{"x":1}}

`),
			dataLine:  `data: {"type":"future_event","unknown_field":{"x":1}}`,
			wantType:  "future_event",
			wantRawOK: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rf, err := ParseSSERawFrame(tc.dataLine, tc.rawFrame)
			if err != nil {
				t.Fatalf("ParseSSERawFrame: %v", err)
			}
			if rf.Type != tc.wantType {
				t.Errorf("Type: got %q, want %q", rf.Type, tc.wantType)
			}
			if tc.wantRawOK && !bytes.Equal(rf.Raw, tc.rawFrame) {
				t.Errorf("Raw mismatch:\n got: %q\nwant: %q", rf.Raw, tc.rawFrame)
			}
		})
	}
}

// TestParseSSERawFrame_EmptyOrDone 返回零值(空 frame / 非 data 前缀)。
func TestParseSSERawFrame_EmptyOrDone(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		rawFrame []byte
	}{
		{"非 data: 前缀", "event: message_start", []byte("event: message_start\n\n")},
		{"data: 空 body", "data: ", nil},
		{"data: [DONE] 哨兵", "data: [DONE]", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rf, err := ParseSSERawFrame(tc.line, tc.rawFrame)
			if err != nil {
				t.Fatalf("ParseSSERawFrame: %v", err)
			}
			if rf.Type != "" || len(rf.Raw) != 0 || rf.TypedEvent != nil {
				t.Errorf("expected zero value, got %+v", rf)
			}
		})
	}
}

// TestAnthropicRawStreamFrame_NotAnthropicStreamEvent 验证 sealed 接口隔离:
// AnthropicRawStreamFrame 不得满足 pkg.AnthropicStreamEvent 接口。
//
// 设计动机:调用方必须用独立的 `case *AnthropicRawStreamFrame:` 分支显式处理,
// 避免误落入 typed 路径导致重新 marshal(参见 stream_event.go 类型注释)。
//
// 若未来重构误给 AnthropicRawStreamFrame 添加 isStreamEvent() 方法,
// 本测试会 fail,迫使设计者重新审视 union 边界。
func TestAnthropicRawStreamFrame_NotAnthropicStreamEvent(t *testing.T) {
	rf := AnthropicRawStreamFrame{}
	if _, ok := any(&rf).(AnthropicStreamEvent); ok {
		t.Fatal("AnthropicRawStreamFrame must NOT implement AnthropicStreamEvent (sealed interface violated)")
	}
}
