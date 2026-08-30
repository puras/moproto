package frames

import (
	"encoding/json"
	"strings"
	"testing"

	"hegq.net/mogo/pkg/protocol/anthropic"
)

// TestBuildMessageStartFrame_RequiredFields 验证 error 路径。
func TestBuildMessageStartFrame_RequiredFields(t *testing.T) {
	// 缺 Usage
	if _, err := BuildMessageStartFrame(MessageStartOptions{Model: "claude-x"}); err == nil {
		t.Error("expected error when Usage is nil")
	}
	// 缺 Model
	if _, err := BuildMessageStartFrame(MessageStartOptions{
		Usage: &anthropic.AnthropicUsage{InputTokens: 10, OutputTokens: 5},
	}); err == nil {
		t.Error("expected error when Model is empty")
	}
}

// TestBuildMessageStartFrame_WireFormat 验证输出的 JSON 形状与 SDK 兼容。
func TestBuildMessageStartFrame_WireFormat(t *testing.T) {
	opts := MessageStartOptions{
		Model: "claude-test",
		Usage: &anthropic.AnthropicUsage{
			InputTokens:              100,
			OutputTokens:             50,
			CacheCreationInputTokens: 0,
			CacheReadInputTokens:     0,
		},
	}
	frame, err := BuildMessageStartFrame(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := string(frame)

	if !strings.HasPrefix(raw, "event: message_start\ndata: ") {
		t.Errorf("missing frame prefix, got head: %q", raw[:minLen(60, len(raw))])
	}
	if !strings.HasSuffix(raw, "\n\n") {
		t.Errorf("frame must end with \\n\\n")
	}

	body := raw[len("event: message_start\ndata: "):]
	body = body[:len(body)-2]

	var parsed struct {
		Type    string `json:"type"`
		Message struct {
			ID      string                    `json:"id"`
			Type    string                    `json:"type"`
			Role    string                    `json:"role"`
			Model   string                    `json:"model"`
			Content []any                     `json:"content"`
			Usage   *anthropic.AnthropicUsage `json:"usage"`
		} `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("data body is not valid JSON: %v", err)
	}

	if parsed.Type != "message_start" {
		t.Errorf("event.type want message_start, got %q", parsed.Type)
	}
	if parsed.Message.Type != "message" {
		t.Errorf("message.type want message, got %q", parsed.Message.Type)
	}
	if parsed.Message.Role != "assistant" {
		t.Errorf("message.role want assistant, got %q", parsed.Message.Role)
	}
	if parsed.Message.Model != "claude-test" {
		t.Errorf("message.model want claude-test, got %q", parsed.Message.Model)
	}
	if parsed.Message.Content == nil {
		t.Error("message.content must be [] not null")
	}
	if len(parsed.Message.Content) != 0 {
		t.Errorf("message.content should be empty, got %d items", len(parsed.Message.Content))
	}
	if parsed.Message.ID == "" {
		t.Error("message.id should be auto-generated non-empty")
	}
	if parsed.Message.Usage == nil || parsed.Message.Usage.InputTokens != 100 {
		t.Errorf("message.usage.input_tokens want 100, got %+v", parsed.Message.Usage)
	}
}

// TestBuildMessageStartFrame_OmitsZeroCacheFields 验证 cache_*_input_tokens 为 0 时被省略。
func TestBuildMessageStartFrame_OmitsZeroCacheFields(t *testing.T) {
	opts := MessageStartOptions{
		Model: "claude-x",
		Usage: &anthropic.AnthropicUsage{InputTokens: 10, OutputTokens: 5},
	}
	frame, _ := BuildMessageStartFrame(opts)
	body := string(frame[len("event: message_start\ndata: "):])
	body = body[:len(body)-2]
	if strings.Contains(body, "cache_creation_input_tokens") {
		t.Errorf("zero cache_creation should be omitted, got: %s", body)
	}
	if strings.Contains(body, "cache_read_input_tokens") {
		t.Errorf("zero cache_read should be omitted, got: %s", body)
	}
}

// TestBuildMessageDeltaFrame_StopReasonAndUsage 验证 stop_reason 与 usage 的字段映射。
func TestBuildMessageDeltaFrame_StopReasonAndUsage(t *testing.T) {
	opts := MessageDeltaOptions{
		StopReason: anthropic.AnthropicStopReasonToolUse,
		Usage: &anthropic.MessageDeltaUsage{
			InputTokens:              200,
			OutputTokens:             80,
			CacheCreationInputTokens: 10,
			CacheReadInputTokens:     5,
		},
	}
	frame, err := BuildMessageDeltaFrame(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := string(frame)

	if !strings.HasPrefix(raw, "event: message_delta\ndata: ") {
		t.Errorf("missing frame prefix: %q", raw[:minLen(60, len(raw))])
	}
	body := raw[len("event: message_delta\ndata: "):]
	body = body[:len(body)-2]

	var parsed struct {
		Type  string                            `json:"type"`
		Delta *anthropic.MessageDeltaEventDelta `json:"delta"`
		Usage *anthropic.MessageDeltaUsage      `json:"usage"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("data body invalid JSON: %v", err)
	}
	if parsed.Type != "message_delta" {
		t.Errorf("event.type want message_delta, got %q", parsed.Type)
	}
	if parsed.Delta == nil || parsed.Delta.StopReason == nil {
		t.Fatal("delta.stop_reason should be set")
	}
	if *parsed.Delta.StopReason != anthropic.AnthropicStopReasonToolUse {
		t.Errorf("delta.stop_reason want tool_use, got %q", *parsed.Delta.StopReason)
	}
	if parsed.Usage == nil || parsed.Usage.InputTokens != 200 {
		t.Errorf("usage.input_tokens want 200, got %+v", parsed.Usage)
	}
	if parsed.Usage.CacheCreationInputTokens != 10 || parsed.Usage.CacheReadInputTokens != 5 {
		t.Errorf("cache fields not propagated: %+v", parsed.Usage)
	}
}

// TestBuildMessageDeltaFrame_EmptyStopReason 验证空 StopReason 不输出 delta.stop_reason。
func TestBuildMessageDeltaFrame_EmptyStopReason(t *testing.T) {
	opts := MessageDeltaOptions{Usage: &anthropic.MessageDeltaUsage{OutputTokens: 5}}
	frame, _ := BuildMessageDeltaFrame(opts)
	body := string(frame[len("event: message_delta\ndata: "):])
	body = body[:len(body)-2]
	if strings.Contains(body, `"stop_reason"`) {
		t.Errorf("empty stop_reason should not appear, got: %s", body)
	}
}

// TestBuildMessageStopFrame_Constant 验证 message_stop 结构固定。
func TestBuildMessageStopFrame_Constant(t *testing.T) {
	frame := BuildMessageStopFrame()
	raw := string(frame)
	if !strings.HasPrefix(raw, "event: message_stop\ndata: ") {
		t.Errorf("frame prefix wrong: %q", raw[:minLen(50, len(raw))])
	}
	body := raw[len("event: message_stop\ndata: "):]
	body = body[:len(body)-2]
	var parsed struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("data not valid JSON: %v", err)
	}
	if parsed.Type != "message_stop" {
		t.Errorf("type want message_stop, got %q", parsed.Type)
	}
}
