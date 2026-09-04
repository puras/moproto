package anthropic_test

import (
	"encoding/json"
	"testing"

	anthropic "github.com/puras/moproto/protocol/anthropic"
)

// TestWireFormatResponseMessage 验证顶层 AnthropicMessageResponse 的 JSON marshal 输出。
func TestResponseWireFormatMessage(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Message_minimal",
			v: anthropic.AnthropicMessageResponse{
				ID:    "msg_01",
				Type:  anthropic.AnthropicMessageTypeMessage,
				Role:  anthropic.AnthropicMessageRoleAssistant,
				Model: "claude-sonnet-4-6",
			},
			want: `{"id":"msg_01","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":null}`,
		},
		{
			name: "Message_with_content_and_usage",
			v: anthropic.AnthropicMessageResponse{
				ID:    "msg_02",
				Type:  anthropic.AnthropicMessageTypeMessage,
				Role:  anthropic.AnthropicMessageRoleAssistant,
				Model: "claude-sonnet-4-6",
				Content: func() []*anthropic.AnthropicContentBlockUnion {
					b, _ := json.Marshal(anthropic.AnthropicTextBlock{Type: anthropic.AnthropicTextBlockTypeText, Text: "hello"})
					return []*anthropic.AnthropicContentBlockUnion{
						{Type: string(anthropic.AnthropicTextBlockTypeText), Text: b},
					}
				}(),
				StopReason: ptrStopReason(anthropic.AnthropicStopReasonEndTurn),
				Usage: &anthropic.AnthropicUsage{
					InputTokens:  10,
					OutputTokens: 5,
				},
			},
			want: `{"id":"msg_02","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":[{"type":"text","text":"hello"}],"stop_reason":"end_turn","usage":{"input_tokens":10,"output_tokens":5}}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal error: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestWireFormatContentBlocks 验证 12 个内容块的 JSON marshal 输出（部分覆盖）。
func TestResponseWireFormatContentBlocks(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "TextBlock_minimal",
			v: anthropic.AnthropicTextBlock{
				Type: anthropic.AnthropicTextBlockTypeText,
				Text: "hello",
			},
			want: `{"type":"text","text":"hello"}`,
		},
		{
			name: "ThinkingBlock",
			v: anthropic.AnthropicThinkingBlock{
				Type:      anthropic.AnthropicThinkingBlockTypeThinking,
				Thinking:  "Let me think...",
				Signature: "sig_xyz",
			},
			want: `{"type":"thinking","thinking":"Let me think...","signature":"sig_xyz"}`,
		},
		{
			name: "RedactedThinkingBlock",
			v: anthropic.AnthropicRedactedThinkingBlock{
				Type: anthropic.AnthropicRedactedThinkingBlockTypeRedactedThinking,
				Data: "encrypted_blob",
			},
			want: `{"type":"redacted_thinking","data":"encrypted_blob"}`,
		},
		{
			name: "ToolUseBlock",
			v: anthropic.AnthropicToolUseBlock{
				Type:  anthropic.AnthropicToolUseBlockTypeToolUse,
				ID:    "tu_1",
				Name:  "get_weather",
				Input: json.RawMessage(`{"city":"SF"}`),
			},
			want: `{"type":"tool_use","id":"tu_1","name":"get_weather","input":{"city":"SF"}}`,
		},
		{
			name: "WebSearchToolResultBlock",
			v: anthropic.AnthropicWebSearchToolResultBlock{
				Type:      anthropic.AnthropicWebSearchToolResultBlockTypeWebSearchToolResult,
				ToolUseID: "tu_1",
				Content:   json.RawMessage(`[{"type":"web_search_result","title":"x","url":"https://example.com"}]`),
			},
			want: `{"type":"web_search_tool_result","tool_use_id":"tu_1","content":[{"type":"web_search_result","title":"x","url":"https://example.com"}]}`,
		},
		{
			name: "ContainerUploadBlock",
			v: anthropic.AnthropicContainerUploadBlock{
				Type:   anthropic.AnthropicContainerUploadBlockTypeContainerUpload,
				FileID: "file_01",
			},
			want: `{"type":"container_upload","file_id":"file_01"}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal error: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestWireFormatUsage 验证 AnthropicUsage + 4 个嵌套类型。
func TestWireFormatUsage(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Usage_minimal",
			v: anthropic.AnthropicUsage{
				InputTokens:  10,
				OutputTokens: 20,
			},
			want: `{"input_tokens":10,"output_tokens":20}`,
		},
		{
			name: "Usage_full",
			v: anthropic.AnthropicUsage{
				InputTokens:              100,
				OutputTokens:             200,
				CacheCreationInputTokens: 30,
				CacheReadInputTokens:     40,
				CacheCreation: &anthropic.AnthropicCacheCreation{
					Ephemeral1hInputTokens: 10,
					Ephemeral5mInputTokens: 20,
				},
				OutputTokensDetails: &anthropic.AnthropicOutputTokensDetails{
					ThinkingTokens: 50,
				},
				ServerToolUse: &anthropic.AnthropicServerToolUsage{
					WebFetchRequests:  1,
					WebSearchRequests: 2,
				},
				ServiceTier:  anthropic.AnthropicUsageServiceTierAuto,
				InferenceGeo: "us",
			},
			want: `{"input_tokens":100,"output_tokens":200,"cache_creation_input_tokens":30,"cache_read_input_tokens":40,"cache_creation":{"ephemeral_1h_input_tokens":10,"ephemeral_5m_input_tokens":20},"output_tokens_details":{"thinking_tokens":50},"server_tool_use":{"web_fetch_requests":1,"web_search_requests":2},"service_tier":"auto","inference_geo":"us"}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal error: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestWireFormatErrorResponses 验证 9 种错误变体 + 顶层 ErrorResponse。
func TestWireFormatErrorResponses(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "InvalidRequestError",
			v: anthropic.AnthropicInvalidRequestError{
				Type:    anthropic.AnthropicInvalidRequestErrorTypeInvalidRequestError,
				Message: "invalid x",
			},
			want: `{"type":"invalid_request_error","message":"invalid x"}`,
		},
		{
			name: "RateLimitError",
			v: anthropic.AnthropicRateLimitError{
				Type:    anthropic.AnthropicRateLimitErrorTypeRateLimitError,
				Message: "rate limited",
			},
			want: `{"type":"rate_limit_error","message":"rate limited"}`,
		},
		{
			name: "ErrorResponse_full",
			v: anthropic.AnthropicErrorResponse{
				Type: anthropic.AnthropicErrorTypeError,
				Error: func() *anthropic.AnthropicErrorObjectUnion {
					errBytes, _ := json.Marshal(anthropic.AnthropicInvalidRequestError{
						Type:    anthropic.AnthropicInvalidRequestErrorTypeInvalidRequestError,
						Message: "bad model",
					})
					return &anthropic.AnthropicErrorObjectUnion{
						Type:                string(anthropic.AnthropicInvalidRequestErrorTypeInvalidRequestError),
						InvalidRequestError: errBytes,
					}
				}(),
				RequestID: "req_01",
			},
			want: `{"type":"error","error":{"type":"invalid_request_error","message":"bad model"},"request_id":"req_01"}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal error: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestWireFormatStreamEvents 验证 8 个 SSE event marshal（部分覆盖）。
func TestWireFormatStreamEvents(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "MessageStartEvent",
			v: anthropic.AnthropicMessageStartEvent{
				Type: anthropic.AnthropicMessageStartEventTypeMessageStart,
				Message: &anthropic.AnthropicMessageResponse{
					ID:    "msg_01",
					Type:  anthropic.AnthropicMessageTypeMessage,
					Role:  anthropic.AnthropicMessageRoleAssistant,
					Model: "claude-sonnet-4-6",
				},
			},
			want: `{"type":"message_start","message":{"id":"msg_01","type":"message","role":"assistant","model":"claude-sonnet-4-6","content":null}}`,
		},
		{
			name: "ContentBlockDeltaEvent_text",
			v: anthropic.AnthropicContentBlockDeltaEvent{
				Type:  anthropic.AnthropicContentBlockDeltaEventTypeContentBlockDelta,
				Index: 0,
				Delta: func() *anthropic.RawContentBlockDeltaUnion {
					b, _ := json.Marshal(anthropic.AnthropicTextDeltaBlock{
						Type: anthropic.AnthropicTextDeltaBlockTypeTextDelta,
						Text: "hi",
					})
					return &anthropic.RawContentBlockDeltaUnion{Type: string(anthropic.AnthropicTextDeltaBlockTypeTextDelta), TextDelta: b}
				}(),
			},
			want: `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"}}`,
		},
		{
			name: "PingEvent",
			v: anthropic.AnthropicPingStreamEvent{
				Type: anthropic.AnthropicPingStreamEventTypePing,
			},
			want: `{"type":"ping"}`,
		},
		{
			name: "MessageStopEvent",
			v: anthropic.AnthropicMessageStopEvent{
				Type: anthropic.AnthropicMessageStopEventTypeMessageStop,
			},
			want: `{"type":"message_stop"}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal error: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestUnionSealedImplementation 验证 4 个 union 接口的封闭性。
//
// 通过类型断言 + 编译期检查（无法从外部包新增实现）。
func TestResponseUnionSealedImplementation(t *testing.T) {
	// AnthropicContentBlock
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicTextBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicThinkingBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicRedactedThinkingBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicToolUseBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicServerToolUseBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicWebSearchToolResultBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicWebFetchToolResultBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicCodeExecutionToolResultBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicBashCodeExecutionToolResultBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicTextEditorCodeExecutionToolResultBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicToolSearchToolResultBlock{}
	var _ anthropic.AnthropicContentBlock = anthropic.AnthropicContainerUploadBlock{}

	// AnthropicErrorObject
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicInvalidRequestError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicAuthenticationError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicPermissionError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicNotFoundError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicRateLimitError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicTimeoutError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicOverloadedError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicAPIError{}
	var _ anthropic.AnthropicErrorObject = anthropic.AnthropicBillingError{}

	// AnthropicStreamEvent
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicMessageStartEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicMessageDeltaEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicMessageStopEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicContentBlockStartEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicContentBlockDeltaEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicContentBlockStopEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicErrorStreamEvent{}
	var _ anthropic.AnthropicStreamEvent = anthropic.AnthropicPingStreamEvent{}

	// RawContentBlockDelta
	var _ anthropic.RawContentBlockDelta = anthropic.AnthropicTextDeltaBlock{}
	var _ anthropic.RawContentBlockDelta = anthropic.AnthropicThinkingDeltaBlock{}
	var _ anthropic.RawContentBlockDelta = anthropic.AnthropicInputJSONDeltaBlock{}
	var _ anthropic.RawContentBlockDelta = anthropic.AnthropicSignatureDeltaBlock{}
	var _ anthropic.RawContentBlockDelta = anthropic.AnthropicCitationsDeltaBlock{}

	// AnthropicAppliedEditUnion
	var _ anthropic.AnthropicAppliedEditUnion = anthropic.AnthropicClearToolUses20250919EditResponse{}
	var _ anthropic.AnthropicAppliedEditUnion = anthropic.AnthropicClearThinking20251015EditResponse{}
}

// TestWireFormatAppliedEditsResponse 验证 applied_edits 响应类型的 JSON marshal。
func TestWireFormatAppliedEditsResponse(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "clear_tool_uses_response",
			v: anthropic.AnthropicClearToolUses20250919EditResponse{
				Type:               anthropic.AnthropicClearToolUses20250919EditTypeClearToolUses20250919,
				ClearedInputTokens: 1500,
				ClearedToolUses:    3,
			},
			want: `{"type":"clear_tool_uses_20250919","cleared_input_tokens":1500,"cleared_tool_uses":3}`,
		},
		{
			name: "clear_thinking_response",
			v: anthropic.AnthropicClearThinking20251015EditResponse{
				Type:                 anthropic.AnthropicClearThinking20251015EditTypeClearThinking20251015,
				ClearedInputTokens:   800,
				ClearedThinkingTurns: 2,
			},
			want: `{"type":"clear_thinking_20251015","cleared_input_tokens":800,"cleared_thinking_turns":2}`,
		},
		{
			name: "context_management_response",
			v: anthropic.AnthropicContextManagementResponse{
				AppliedEdits: func() []anthropic.AnthropicAppliedEditUnionUnion {
					a, _ := json.Marshal(anthropic.AnthropicClearToolUses20250919EditResponse{
						Type:               anthropic.AnthropicClearToolUses20250919EditTypeClearToolUses20250919,
						ClearedInputTokens: 1500,
						ClearedToolUses:    3,
					})
					b, _ := json.Marshal(anthropic.AnthropicClearThinking20251015EditResponse{
						Type:                 anthropic.AnthropicClearThinking20251015EditTypeClearThinking20251015,
						ClearedInputTokens:   800,
						ClearedThinkingTurns: 2,
					})
					return []anthropic.AnthropicAppliedEditUnionUnion{
						{Type: string(anthropic.AnthropicClearToolUses20250919EditTypeClearToolUses20250919), ClearToolUses: a},
						{Type: string(anthropic.AnthropicClearThinking20251015EditTypeClearThinking20251015), ClearThinking: b},
					}
				}(),
			},
			want: `{"applied_edits":[{"type":"clear_tool_uses_20250919","cleared_input_tokens":1500,"cleared_tool_uses":3},{"type":"clear_thinking_20251015","cleared_input_tokens":800,"cleared_thinking_turns":2}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("json.Marshal error: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s, want %s", string(got), tc.want)
			}
		})
	}
}

// ptrStopReason 辅助函数：把 AnthropicStopReason 转成 *AnthropicStopReason。
func ptrStopReason(r anthropic.AnthropicStopReason) *anthropic.AnthropicStopReason {
	return &r
}
