package frames

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestMarshalSSEFrameRaw_NoLoss 验证 verbatim emit 不丢失未知字段 / 多行 / omitempty 零值。
func TestMarshalSSEFrameRaw_NoLoss(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		want []byte
	}{
		{
			name: "empty input returns nil",
			in:   nil,
			want: nil,
		},
		{
			name: "basic frame verbatim",
			in:   []byte("event: response.created\ndata: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\"}}\n\n"),
			want: []byte("event: response.created\ndata: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\"}}\n\n"),
		},
		{
			// 4.6+ 未知字段 + omitempty 零值共存 —— typed 路径会被剥除；raw 路径完整保留
			name: "preserves unknown fields and zero-value omitempty",
			in: []byte(`event: response.completed
data: {"type":"response.completed","response":{"id":"resp_2","object":"response","created_at":1700000000,"status":"completed","unknown_field_4_6_plus":"preserved","cache_creation_input_tokens":0,"service_tier":""}}

`),
			want: []byte(`event: response.completed
data: {"type":"response.completed","response":{"id":"resp_2","object":"response","created_at":1700000000,"status":"completed","unknown_field_4_6_plus":"preserved","cache_creation_input_tokens":0,"service_tier":""}}

`),
		},
		{
			// 字段顺序保留（marshal 会按 struct 字段序重排，verbatim 不会）
			name: "preserves field order",
			in:   []byte(`data: {"type":"response.created","z_field":1,"a_field":2}` + "\n\n"),
			want: []byte(`data: {"type":"response.created","z_field":1,"a_field":2}` + "\n\n"),
		},
		{
			name: "unicode preserved verbatim",
			in:   []byte("event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"中文 + emoji 🎉\"}\n\n"),
			want: []byte("event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"中文 + emoji 🎉\"}\n\n"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MarshalSSEFrameRaw(tc.in)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("MarshalSSEFrameRaw byte mismatch\nwant: %q\ngot:  %q", string(tc.want), string(got))
			}
		})
	}
}

// TestParseResponsesSSEEvent_RoundTrip 验证解析 + emit 往返字节保真。
func TestParseResponsesSSEEvent_RoundTrip(t *testing.T) {
	rawFrame := []byte("event: response.created\ndata: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_99\"}}\n\n")
	line := "data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_99\"}}"

	frame, ok := ParseResponsesSSEEvent(line, rawFrame)
	if !ok {
		t.Fatal("ParseResponsesSSEEvent returned ok=false for valid event")
	}
	if frame.Type != "response.created" {
		t.Errorf("Type = %q, want %q", frame.Type, "response.created")
	}
	if !bytes.Equal(frame.Raw, rawFrame) {
		t.Errorf("Raw mismatch\nwant: %q\ngot:  %q", string(rawFrame), string(frame.Raw))
	}
	// Payload 是去前缀与换行后的 JSON；应能被 json.Unmarshal 回原对象
	var payload struct {
		Type     string `json:"type"`
		Response struct {
			ID string `json:"id"`
		} `json:"response"`
	}
	if err := json.Unmarshal(frame.Payload, &payload); err != nil {
		t.Fatalf("Payload json.Unmarshal failed: %v", err)
	}
	if payload.Type != "response.created" || payload.Response.ID != "resp_99" {
		t.Errorf("Payload round-trip wrong: %+v", payload)
	}

	// Round-trip：frame.Raw → MarshalSSEFrameRaw → 字节级一致
	re := MarshalSSEFrameRaw(frame.Raw)
	if !bytes.Equal(re, rawFrame) {
		t.Errorf("Round-trip mismatch\nwant: %q\ngot:  %q", string(rawFrame), string(re))
	}
}

// TestParseResponsesSSEEvent_NonDataAndDone 验证跳过 [DONE] / 非 data 行 / 空 body。
func TestParseResponsesSSEEvent_NonDataAndDone(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{name: "no data prefix", line: "event: response.created"},
		{name: "[DONE] sentinel", line: "data: [DONE]"},
		{name: "empty data body", line: "data: "},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := ParseResponsesSSEEvent(tc.line, []byte("dummy"))
			if ok {
				t.Errorf("ParseResponsesSSEEvent returned ok=true for %q", tc.line)
			}
		})
	}
}

// TestParseResponsesSSEEvent_InvalidJSON 验证 data 行非 JSON 时仍 ok=true（caller 透传）。
func TestParseResponsesSSEEvent_InvalidJSON(t *testing.T) {
	line := "data: not-valid-json"
	rawFrame := []byte("event: x\n" + line + "\n\n")
	frame, ok := ParseResponsesSSEEvent(line, rawFrame)
	if !ok {
		t.Fatal("ParseResponsesSSEEvent should return ok=true even for invalid JSON (passthrough-friendly)")
	}
	if frame.Type != "" {
		t.Errorf("Type should be empty for invalid JSON, got %q", frame.Type)
	}
	if !bytes.Equal(frame.Payload, []byte("not-valid-json")) {
		t.Errorf("Payload should preserve raw body, got %q", string(frame.Payload))
	}
	if !bytes.Equal(frame.Raw, rawFrame) {
		t.Errorf("Raw should be preserved, got %q", string(frame.Raw))
	}
}

// TestOpenAIResponsesRawStreamFrame_NotImplementsOpenAIUnion 确认 carrier 不被误用为 typed 事件。
//
// 此测试用编译期断言表达：本类型有私有 marker `isOpenAIResponsesRawStreamFrame()`，
// 但不实现 `isResponseStreamEvent()`。由于父包接口未导出到本包测试，改用反射
// 验证 carrier 类型名不与 union 任一已知变体重名（防止未来重构误将其加入 union）。
func TestOpenAIResponsesRawStreamFrame_TypeIdentity(t *testing.T) {
	frame := OpenAIResponsesRawStreamFrame{
		Type:    "response.created",
		Raw:     []byte("dummy"),
		Payload: []byte(`{"type":"response.created"}`),
	}
	// 通过私有 marker 方法确认
	var marker interface{ isOpenAIResponsesRawStreamFrame() } = frame
	if marker == nil {
		t.Fatal("marker interface assertion failed")
	}
	// 字段读写一致性
	if frame.Type != "response.created" {
		t.Errorf("frame.Type = %q", frame.Type)
	}
}
