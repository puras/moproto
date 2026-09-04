package openai_test

import (
	"encoding/json"
	"testing"

	openai "github.com/puras/moproto/protocol/openai"
)

// === A. ResponseParams 顶层字段 marshal 测试（6 case） ===

func TestWireFormatResponseParams(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Params_minimal",
			v: openai.OpenAIResponseParams{
				Model: "gpt-4o",
				Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
			},
			want: `{"model":"gpt-4o","input":"hello"}`,
		},
		{
			name: "Params_with_input_items",
			v: openai.OpenAIResponseParams{
				Model: "gpt-4o",
				Input: &openai.OpenAIResponseInputUnionFlat{
					Type: "items",
					Items: []openai.OpenAIResponseInputItemUnionFlat{
						{
							Type:           "message",
							MessageRole:    openai.OpenAIResponseInputMessageItemRoleUser,
							MessageContent: &openai.OpenAIResponseInputMessageContentUnionFlat{Type: "string", Text: "hi"},
						},
					},
				},
			},
			want: `{"model":"gpt-4o","input":[{"type":"message","role":"user","content":"hi"}]}`,
		},
		{
			name: "Params_with_tools_and_choice",
			v: openai.OpenAIResponseParams{
				Model: "gpt-4o",
				Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hi"},
				Tools: []openai.OpenAIResponsesToolUnionFlat{
					openai.OpenAIResponsesToolUnionFlat{
						Type: "function",
						Name: "get_weather",
					},
				},
				ToolChoice: openai.OpenAIResponseToolChoiceOptionAuto{Type: "auto"},
			},
			want: `{"model":"gpt-4o","input":"hi","tools":[{"type":"function","name":"get_weather"}],"tool_choice":{"type":"auto"}}`,
		},
		{
			name: "Params_with_reasoning_config",
			v: openai.OpenAIResponseParams{
				Model: "o1",
				Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hi"},
				Reasoning: &openai.OpenAIResponseReasoningParam{
					Effort:  openai.OpenAIReasoningEffortHigh,
					Summary: "auto",
				},
			},
			want: `{"model":"o1","input":"hi","reasoning":{"effort":"high","summary":"auto"}}`,
		},
		{
			name: "Params_with_text_config",
			v: openai.OpenAIResponseParams{
				Model: "gpt-4o",
				Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hi"},
				Text: &openai.OpenAIResponseTextConfig{
					Verbosity: openai.OpenAIVerbosityHigh,
				},
			},
			want: `{"model":"gpt-4o","input":"hi","text":{"verbosity":"high"}}`,
		},
		{
			name: "Params_with_conversation",
			v: openai.OpenAIResponseParams{
				Model: "gpt-4o",
				Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hi"},
				Conversation: &openai.OpenAIResponseConversationParam{
					ID: "conv_123",
				},
			},
			want: `{"model":"gpt-4o","input":"hi","conversation":{"id":"conv_123"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === B. InputItem union marshal 测试（4 case，只测关键变体） ===

func TestWireFormatResponseInputItemUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "MessageItem",
			v: openai.OpenAIResponseInputMessageItem{
				Type:    "message",
				Role:    openai.OpenAIResponseInputMessageItemRoleUser,
				Content: &openai.OpenAIResponseInputMessageContentUnionFlat{Type: "string", Text: "hi"},
			},
			want: `{"type":"message","role":"user","content":"hi"}`,
		},
		{
			name: "FunctionToolCallItem",
			v: openai.OpenAIResponseInputFunctionToolCallItem{
				Type:      "function_call",
				CallID:    "call_123",
				Name:      "get_weather",
				Arguments: `{"city":"SF"}`,
			},
			want: `{"type":"function_call","call_id":"call_123","name":"get_weather","arguments":"{\"city\":\"SF\"}"}`,
		},
		{
			name: "ItemReferenceItem",
			v: openai.OpenAIResponseInputItemReferenceItem{
				Type: "item_reference",
				ID:   "msg_123",
			},
			want: `{"type":"item_reference","id":"msg_123"}`,
		},
		{
			name: "ReasoningItem",
			v: openai.OpenAIResponseInputReasoningItem{
				Type: "reasoning",
				ID:   "rs_123",
			},
			want: `{"type":"reasoning","id":"rs_123"}`,
		},
		// 4.6+ 新增变体
		{
			name: "LocalShellCallItem",
			v: openai.OpenAIResponseInputLocalShellCallItem{
				Type:   "local_shell_call",
				ID:     "lsc_123",
				CallID: "call_123",
				Action: &openai.OpenAIResponseInputLocalShellAction{
					Type:    "exec",
					Command: "ls -la",
				},
				Status: "in_progress",
			},
			want: `{"type":"local_shell_call","id":"lsc_123","call_id":"call_123","action":{"type":"exec","command":"ls -la"},"status":"in_progress"}`,
		},
		{
			name: "McpApprovalRequestItem",
			v: openai.OpenAIResponseInputMcpApprovalRequestItem{
				Type:              "mcp_approval_request",
				ID:                "mcap_123",
				ServerLabel:       "github",
				Name:              "create_issue",
				Arguments:         `{"title":"bug"}`,
				ApprovalRequestID: "apr_123",
			},
			want: `{"type":"mcp_approval_request","id":"mcap_123","server_label":"github","name":"create_issue","arguments":"{\"title\":\"bug\"}","approval_request_id":"apr_123"}`,
		},
		{
			name: "CompactionTriggerItem",
			v: openai.OpenAIResponseInputCompactionTriggerItem{
				Type: "compaction_trigger",
			},
			want: `{"type":"compaction_trigger"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === C. OutputItem union marshal 测试（3 case） ===

func TestWireFormatResponseOutputItemUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Message",
			v: openai.OpenAIResponseOutputMessage{
				Type:   "message",
				ID:     "msg_123",
				Role:   "assistant",
				Status: "completed",
				Content: []openai.OpenAIResponseOutputMessageContent{
					{Type: "output_text", Text: "Hello!"},
				},
			},
			want: `{"type":"message","id":"msg_123","role":"assistant","status":"completed","content":[{"type":"output_text","text":"Hello!"}]}`,
		},
		{
			name: "FunctionCall",
			v: openai.OpenAIResponseOutputFunctionCall{
				Type:      "function_call",
				ID:        "fc_123",
				CallID:    "call_123",
				Name:      "get_weather",
				Arguments: `{"city":"SF"}`,
			},
			want: `{"type":"function_call","id":"fc_123","call_id":"call_123","name":"get_weather","arguments":"{\"city\":\"SF\"}"}`,
		},
		{
			name: "Reasoning",
			v: openai.OpenAIResponseOutputReasoning{
				Type:    "reasoning",
				ID:      "rs_123",
				Status:  "completed",
				Content: []string{"Let me think..."},
			},
			want: `{"type":"reasoning","id":"rs_123","status":"completed","content":["Let me think..."]}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === D. ResponsesTool union marshal 测试（3 case） ===

func TestWireFormatResponsesToolUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "FunctionTool",
			v: openai.OpenAIResponsesFunctionToolParam{
				Type: "function",
				Name: "get_weather",
			},
			want: `{"type":"function","name":"get_weather"}`,
		},
		{
			name: "WebSearchTool",
			v: openai.OpenAIResponsesWebSearchToolParam{
				Type:              "web_search",
				SearchContextSize: "high",
			},
			want: `{"type":"web_search","search_context_size":"high"}`,
		},
		{
			name: "ImageGenerationTool",
			v: openai.OpenAIResponsesImageGenerationToolParam{
				Type:         "image_generation",
				Model:        "gpt-image-1",
				Quality:      "high",
				OutputFormat: "png",
			},
			want: `{"type":"image_generation","model":"gpt-image-1","quality":"high","output_format":"png"}`,
		},
		// 4.6+ 新增变体
		{
			name: "LocalShellTool",
			v: openai.OpenAIResponsesLocalShellToolParam{
				Type: "local_shell",
			},
			want: `{"type":"local_shell"}`,
		},
		{
			name: "CustomTool_grammar",
			v: openai.OpenAIResponsesCustomToolParam{
				Type:        "custom",
				Name:        "extract_entities",
				Description: "Extract named entities",
				Format: &openai.OpenAIResponsesCustomToolFormatGrammar{
					Type:       "grammar",
					Syntax:     "lark",
					Definition: `start: "x"`,
				},
			},
			want: `{"type":"custom","name":"extract_entities","description":"Extract named entities","format":{"type":"grammar","syntax":"lark","definition":"start: \"x\""}}`,
		},
		{
			name: "ToolSearchTool",
			v: openai.OpenAIResponsesToolSearchToolParam{
				Type:      "tool_search",
				Execution: "auto",
			},
			want: `{"type":"tool_search","execution":"auto"}`,
		},
		{
			name: "NamespaceTool",
			v: openai.OpenAIResponsesNamespaceToolParam{
				Type:        "namespace",
				Name:        "math",
				Description: "Math tools",
			},
			want: `{"type":"namespace","name":"math","description":"Math tools"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === E. Response 响应 marshal 测试（3 case） ===

func TestWireFormatResponse(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Response_minimal",
			v: openai.OpenAIResponse{
				ID:        "resp_123",
				Object:    "response",
				CreatedAt: 1234567890,
				Model:     "gpt-4o",
				Status:    "completed",
				Output:    []openai.OpenAIResponseOutputItemUnion{},
			},
			want: `{"id":"resp_123","object":"response","created_at":1234567890,"model":"gpt-4o","status":"completed","output":[]}`,
		},
		{
			name: "Response_with_usage_breakdown",
			v: openai.OpenAIResponse{
				ID: "resp_123", Object: "response", CreatedAt: 1234567890, Model: "gpt-4o",
				Status: "completed",
				Output: []openai.OpenAIResponseOutputItemUnion{},
				Usage: &openai.OpenAIResponseUsage{
					InputTokens:  10,
					OutputTokens: 5,
					TotalTokens:  15,
					InputTokensDetails: openai.OpenAIResponseUsageInputTokensDetails{
						CachedTokens: 8,
					},
					OutputTokensDetails: openai.OpenAIResponseUsageOutputTokensDetails{
						ReasoningTokens: 2,
					},
				},
			},
			want: `{"id":"resp_123","object":"response","created_at":1234567890,"model":"gpt-4o","status":"completed","output":[],"usage":{"input_tokens":10,"output_tokens":5,"total_tokens":15,"input_tokens_details":{"cached_tokens":8},"output_tokens_details":{"reasoning_tokens":2}}}`,
		},
		{
			name: "Response_with_incomplete_details",
			v: openai.OpenAIResponse{
				ID: "resp_123", Object: "response", CreatedAt: 1234567890, Model: "gpt-4o",
				Status: "incomplete",
				Output: []openai.OpenAIResponseOutputItemUnion{},
				IncompleteDetails: &openai.OpenAIResponseIncompleteDetails{
					Reason: "max_output_tokens",
				},
			},
			want: `{"id":"resp_123","object":"response","created_at":1234567890,"model":"gpt-4o","status":"incomplete","output":[],"incomplete_details":{"reason":"max_output_tokens"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === F. ResponseStreamEvent union marshal 测试（3 case，只测主线） ===

func TestWireFormatResponseStreamEvent(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "CreatedEvent",
			v: openai.OpenAIResponseCreatedEvent{
				Type:     "response.created",
				Sequence: 0,
				Response: &openai.OpenAIResponse{
					ID: "resp_123", Object: "response", CreatedAt: 1234567890, Model: "gpt-4o", Status: "in_progress",
					Output: []openai.OpenAIResponseOutputItemUnion{},
				},
			},
			want: `{"type":"response.created","sequence_number":0,"response":{"id":"resp_123","object":"response","created_at":1234567890,"model":"gpt-4o","status":"in_progress","output":[]}}`,
		},
		{
			name: "TextDeltaEvent",
			v: openai.OpenAIResponseTextDeltaEvent{
				Type:         "response.output_text.delta",
				Sequence:     5,
				ItemID:       "msg_123",
				OutputIndex:  0,
				ContentIndex: 0,
				Delta:        "hel",
			},
			want: `{"type":"response.output_text.delta","sequence_number":5,"item_id":"msg_123","output_index":0,"content_index":0,"delta":"hel"}`,
		},
		{
			name: "CompletedEvent",
			v: openai.OpenAIResponseCompletedEvent{
				Type:     "response.completed",
				Sequence: 10,
				Response: &openai.OpenAIResponse{
					ID: "resp_123", Object: "response", CreatedAt: 1234567890, Model: "gpt-4o", Status: "completed",
					Output: []openai.OpenAIResponseOutputItemUnion{},
				},
			},
			want: `{"type":"response.completed","sequence_number":10,"response":{"id":"resp_123","object":"response","created_at":1234567890,"model":"gpt-4o","status":"completed","output":[]}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === G. 其他能力 marshal 测试 ===

func TestWireFormatOtherCapabilities(t *testing.T) {
	t.Run("EmbeddingParams_minimal", func(t *testing.T) {
		v := openai.OpenAIEmbeddingParams{
			Input: &openai.OpenAIEmbeddingInputUnionFlat{
				Type:   "string",
				String: func() *openai.OpenAIEmbeddingInputString { s := openai.OpenAIEmbeddingInputString("hi"); return &s }(),
			},
			Model: "text-embedding-3-small",
		}
		want := `{"input":"hi","model":"text-embedding-3-small"}`
		got, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(got) != want {
			t.Errorf("got %s\nwant %s", got, want)
		}
	})

	t.Run("ImageGenerateParams", func(t *testing.T) {
		v := openai.OpenAIImageGenerateParams{
			Prompt:         "a cat",
			Model:          "dall-e-3",
			Size:           "1024x1024",
			ResponseFormat: "url",
		}
		want := `{"prompt":"a cat","model":"dall-e-3","response_format":"url","size":"1024x1024"}`
		got, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(got) != want {
			t.Errorf("got %s\nwant %s", got, want)
		}
	})

	t.Run("AudioSpeechParams", func(t *testing.T) {
		speed := 1.5
		v := openai.OpenAIAudioSpeechParams{
			Model:          "tts-1",
			Input:          "Hello world",
			Voice:          openai.OpenAIAudioVoiceAlloy,
			ResponseFormat: openai.OpenAIAudioFormatMP3,
			Speed:          &speed,
		}
		want := `{"model":"tts-1","input":"Hello world","voice":"alloy","response_format":"mp3","speed":1.5}`
		got, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(got) != want {
			t.Errorf("got %s\nwant %s", got, want)
		}
	})

	t.Run("ModerationParams", func(t *testing.T) {
		v := openai.OpenAIModerationParams{
			Input: &openai.OpenAIModerationInputUnionFlat{
				Type: "string",
				String: func() *openai.OpenAIModerationInputString {
					s := openai.OpenAIModerationInputString("test text")
					return &s
				}(),
			},
			Model: "omni-moderation-latest",
		}
		want := `{"input":"test text","model":"omni-moderation-latest"}`
		got, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(got) != want {
			t.Errorf("got %s\nwant %s", got, want)
		}
	})
}

// === H. Error marshal 测试（3 case） ===

func TestWireFormatError(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "InvalidRequest",
			v: openai.OpenAIErrorResponse{
				Error: openai.OpenAIError{
					Message: "Invalid parameter",
					Type:    "invalid_request_error",
				},
			},
			want: `{"error":{"message":"Invalid parameter","type":"invalid_request_error"}}`,
		},
		{
			name: "Authentication",
			v: openai.OpenAIErrorResponse{
				Error: openai.OpenAIError{
					Message: "Invalid API key",
					Type:    "authentication_error",
					Code:    ptrString("invalid_api_key"),
				},
			},
			want: `{"error":{"message":"Invalid API key","type":"authentication_error","code":"invalid_api_key"}}`,
		},
		{
			name: "RateLimit",
			v: openai.OpenAIErrorResponse{
				Error: openai.OpenAIError{
					Message: "Rate limit exceeded",
					Type:    "rate_limit_error",
					Code:    ptrString("rate_limit_exceeded"),
				},
			},
			want: `{"error":{"message":"Rate limit exceeded","type":"rate_limit_error","code":"rate_limit_exceeded"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === I. Union 封闭性测试 ===

func TestResponseInputItemUnionClosure(t *testing.T) {
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputMessageItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputFunctionToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputFunctionToolCallOutputItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputFileSearchToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputWebSearchToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputComputerToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputReasoningItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputItemReferenceItem{}
	// 4.6+ 新增变体
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputApplyPatchToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputMcpToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputCodeInterpreterToolCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputImageGenerationCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputToolSearchCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputToolSearchOutputItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputComputerCallOutputItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputLocalShellCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputLocalShellCallOutputItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputShellCallItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputShellCallOutputItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputApplyPatchCallOutputItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputMcpApprovalRequestItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputMcpApprovalResponseItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputCompactionTriggerItem{}
	var _ openai.OpenAIResponseInputItemUnion = openai.OpenAIResponseInputEasyInputMessageItem{}
}

func TestResponsesToolUnionClosure(t *testing.T) {
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesFunctionToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesWebSearchToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesCodeInterpreterToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesFileSearchToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesImageGenerationToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesMcpToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesComputerToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesApplyPatchToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesCodeExecutionToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesHostedToolParam{}
	// 4.6+ 新增变体
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesLocalShellToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesShellToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesCustomToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesToolSearchToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesWebSearchPreviewToolParam{}
	var _ openai.OpenAIResponsesToolUnion = openai.OpenAIResponsesNamespaceToolParam{}
}

func TestResponseOutputItemUnionClosure(t *testing.T) {
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputMessage{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputFunctionCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputWebSearchCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputFileSearchCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputComputerCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputReasoning{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputMcpCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputMcpListTools{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputCodeInterpreterCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputImageGenerationCall{}
	// 4.6+ 新增变体
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputToolSearchCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputToolSearchOutput{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputLocalShellCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputLocalShellCallOutput{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputShellCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputShellCallOutput{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputApplyPatchToolCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputApplyPatchToolCallOutput{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputMcpApprovalRequest{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputMcpApprovalResponse{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputCustomToolCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputCustomToolCallOutput{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputCompactionItem{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputFunctionShellToolCall{}
	var _ openai.OpenAIResponseOutputItemUnion = openai.OpenAIResponseOutputFunctionShellToolCallOutput{}
}

// === M. ToolChoice union marshal 测试 ===

func TestWireFormatResponseToolChoiceUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "ToolChoice_auto",
			v:    openai.OpenAIResponseToolChoiceOptionAuto{Type: "auto"},
			want: `{"type":"auto"}`,
		},
		{
			name: "ToolChoice_required",
			v:    openai.OpenAIResponseToolChoiceOptionRequired{Type: "required"},
			want: `{"type":"required"}`,
		},
		{
			name: "ToolChoice_none",
			v:    openai.OpenAIResponseToolChoiceOptionNone{Type: "none"},
			want: `{"type":"none"}`,
		},
		{
			name: "ToolChoice_allowed",
			v: openai.OpenAIResponseToolChoiceOptionAllowedParam{
				Type: "allowed_tools",
				Mode: openai.OpenAIResponseToolChoiceOptionModeRequired,
				Tools: []openai.OpenAIResponseToolChoiceOptionAllowedToolEntry{
					{Type: "function", Name: "get_weather"},
					{Type: "function", Name: "get_time"},
				},
			},
			want: `{"type":"allowed_tools","mode":"required","tools":[{"type":"function","name":"get_weather"},{"type":"function","name":"get_time"}]}`,
		},
		{
			name: "ToolChoice_function",
			v:    openai.OpenAIResponseToolChoiceOptionFunctionParam{Type: "function", Name: "get_weather"},
			want: `{"type":"function","name":"get_weather"}`,
		},
		{
			name: "ToolChoice_mcp",
			v: openai.OpenAIResponseToolChoiceOptionMcpParam{
				Type:        "mcp",
				ServerLabel: "github",
				Name:        "create_issue",
			},
			want: `{"type":"mcp","server_label":"github","name":"create_issue"}`,
		},
		{
			name: "ToolChoice_custom",
			v:    openai.OpenAIResponseToolChoiceOptionCustomParam{Type: "custom", Name: "extract"},
			want: `{"type":"custom","name":"extract"}`,
		},
		{
			name: "ToolChoice_types",
			v: openai.OpenAIResponseToolChoiceOptionTypesParam{
				Type: []string{"function", "web_search"},
			},
			want: `{"type":["function","web_search"]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === N. OutputItem 新变体 marshal 测试 ===

func TestWireFormatResponseOutputItemUnionP1(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "LocalShellCall",
			v: openai.OpenAIResponseOutputLocalShellCall{
				Type:   "local_shell_call",
				ID:     "lsc_123",
				CallID: "call_123",
				Status: "completed",
				Action: &openai.OpenAIResponseInputLocalShellAction{
					Type:    "exec",
					Command: "pwd",
				},
			},
			want: `{"type":"local_shell_call","id":"lsc_123","call_id":"call_123","status":"completed","action":{"type":"exec","command":"pwd"}}`,
		},
		{
			name: "ShellCallOutput",
			v: openai.OpenAIResponseOutputShellCallOutput{
				Type:   "shell_call_output",
				ID:     "sco_123",
				Output: "hello\n",
			},
			want: `{"type":"shell_call_output","id":"sco_123","output":"hello\n"}`,
		},
		{
			name: "CustomToolCall",
			v: openai.OpenAIResponseOutputCustomToolCall{
				Type:   "custom_tool_call",
				ID:     "ctc_123",
				CallID: "call_123",
				Name:   "extract",
				Input:  `{"text":"x"}`,
				Status: "completed",
			},
			want: `{"type":"custom_tool_call","id":"ctc_123","call_id":"call_123","name":"extract","input":"{\"text\":\"x\"}","status":"completed"}`,
		},
		{
			name: "CompactionItem",
			v: openai.OpenAIResponseOutputCompactionItem{
				Type: "compaction",
				ID:   "cmp_123",
			},
			want: `{"type":"compaction","id":"cmp_123"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestResponseToolChoiceUnionClosure(t *testing.T) {
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionAuto{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionRequired{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionNone{}
	// 4.6+ 新增变体
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionAllowedParam{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionApplyPatchParam{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionCustomParam{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionFunctionParam{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionMcpParam{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionShellParam{}
	var _ openai.OpenAIResponseToolChoiceUnion = openai.OpenAIResponseToolChoiceOptionTypesParam{}
}

// ptrString 辅助函数定义在 wire_format_chat_test.go 中

// === K. 补齐的 6 个流事件变体 marshal 测试 ===

func TestWireFormatResponseStreamEventP1(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "RefusalDeltaEvent",
			v: openai.OpenAIResponseRefusalDeltaEvent{
				Type:         "response.refusal.delta",
				Sequence:     7,
				ItemID:       "msg_001",
				OutputIndex:  0,
				ContentIndex: 0,
				Delta:        "I cannot help",
			},
			want: `{"type":"response.refusal.delta","sequence_number":7,"item_id":"msg_001","output_index":0,"content_index":0,"delta":"I cannot help"}`,
		},
		{
			name: "ReasoningSummaryTextDeltaEvent",
			v: openai.OpenAIResponseReasoningSummaryTextDeltaEvent{
				Type:         "response.reasoning_summary_text.delta",
				Sequence:     8,
				ItemID:       "rs_001",
				OutputIndex:  0,
				SummaryIndex: 0,
				Delta:        "thinking step",
			},
			want: `{"type":"response.reasoning_summary_text.delta","sequence_number":8,"item_id":"rs_001","output_index":0,"summary_index":0,"delta":"thinking step"}`,
		},
		{
			name: "ReasoningSummaryTextDoneEvent",
			v: openai.OpenAIResponseReasoningSummaryTextDoneEvent{
				Type:         "response.reasoning_summary_text.done",
				Sequence:     9,
				ItemID:       "rs_001",
				OutputIndex:  0,
				SummaryIndex: 0,
				Text:         "final summary",
			},
			want: `{"type":"response.reasoning_summary_text.done","sequence_number":9,"item_id":"rs_001","output_index":0,"summary_index":0,"text":"final summary"}`,
		},
		{
			name: "OutputTextAnnotationAddedEvent",
			v: openai.OpenAIResponseOutputTextAnnotationAddedEvent{
				Type:            "response.output_text.annotation.added",
				Sequence:        11,
				ItemID:          "msg_002",
				OutputIndex:     0,
				ContentIndex:    0,
				AnnotationIndex: 0,
				Annotation:      json.RawMessage(`{"type":"url_citation","url":"https://example.com","title":"example"}`),
			},
			want: `{"type":"response.output_text.annotation.added","sequence_number":11,"item_id":"msg_002","output_index":0,"content_index":0,"annotation_index":0,"annotation":{"type":"url_citation","url":"https://example.com","title":"example"}}`,
		},
		{
			name: "CustomToolCallInputDeltaEvent",
			v: openai.OpenAIResponseCustomToolCallInputDeltaEvent{
				Type:        "response.custom_tool_call_input.delta",
				Sequence:    12,
				ItemID:      "ctc_001",
				OutputIndex: 0,
				Delta:       "{\"q\":\"par",
			},
			want: `{"type":"response.custom_tool_call_input.delta","sequence_number":12,"item_id":"ctc_001","output_index":0,"delta":"{\"q\":\"par"}`,
		},
		{
			name: "CustomToolCallInputDoneEvent",
			v: openai.OpenAIResponseCustomToolCallInputDoneEvent{
				Type:        "response.custom_tool_call_input.done",
				Sequence:    13,
				ItemID:      "ctc_001",
				OutputIndex: 0,
				Input:       "{\"q\":\"partial\"}",
			},
			want: `{"type":"response.custom_tool_call_input.done","sequence_number":13,"item_id":"ctc_001","output_index":0,"input":"{\"q\":\"partial\"}"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === L. 6 个补齐变体的 union 封闭性测试 ===

func TestResponseStreamEventUnionClosureP1(t *testing.T) {
	var _ openai.OpenAIResponseStreamEventUnion = openai.OpenAIResponseRefusalDeltaEvent{}
	var _ openai.OpenAIResponseStreamEventUnion = openai.OpenAIResponseReasoningSummaryTextDeltaEvent{}
	var _ openai.OpenAIResponseStreamEventUnion = openai.OpenAIResponseReasoningSummaryTextDoneEvent{}
	var _ openai.OpenAIResponseStreamEventUnion = openai.OpenAIResponseOutputTextAnnotationAddedEvent{}
	var _ openai.OpenAIResponseStreamEventUnion = openai.OpenAIResponseCustomToolCallInputDeltaEvent{}
	var _ openai.OpenAIResponseStreamEventUnion = openai.OpenAIResponseCustomToolCallInputDoneEvent{}
}

// === K. Embedding InputUnion marshal 测试（4 case） ===

func TestWireFormatEmbeddingInputUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "InputString",
			v: openai.OpenAIEmbeddingParams{
				Input: &openai.OpenAIEmbeddingInputUnionFlat{
					Type:   "string",
					String: func() *openai.OpenAIEmbeddingInputString { s := openai.OpenAIEmbeddingInputString("hi"); return &s }(),
				},
				Model: string(openai.OpenAIEmbeddingModelTextEmbedding3Small),
			},
			want: `{"input":"hi","model":"text-embedding-3-small"}`,
		},
		{
			name: "InputStringArray",
			v: openai.OpenAIEmbeddingParams{
				Input: &openai.OpenAIEmbeddingInputUnionFlat{
					Type: "string_array",
					StringArray: func() *openai.OpenAIEmbeddingInputStringArray {
						s := openai.OpenAIEmbeddingInputStringArray{"hi", "hello"}
						return &s
					}(),
				},
				Model: string(openai.OpenAIEmbeddingModelTextEmbedding3Large),
			},
			want: `{"input":["hi","hello"],"model":"text-embedding-3-large"}`,
		},
		{
			name: "InputTokenArray",
			v: openai.OpenAIEmbeddingParams{
				Input: &openai.OpenAIEmbeddingInputUnionFlat{
					Type: "token_array",
					TokenArray: func() *openai.OpenAIEmbeddingInputTokenArray {
						s := openai.OpenAIEmbeddingInputTokenArray{1, 2, 3}
						return &s
					}(),
				},
				Model: string(openai.OpenAIEmbeddingModelTextEmbeddingAda002),
			},
			want: `{"input":[1,2,3],"model":"text-embedding-ada-002"}`,
		},
		{
			name: "InputTokenArrayArray",
			v: openai.OpenAIEmbeddingParams{
				Input: &openai.OpenAIEmbeddingInputUnionFlat{
					Type: "token_array_array",
					TokenArrayArray: func() *openai.OpenAIEmbeddingInputTokenArrayArray {
						s := openai.OpenAIEmbeddingInputTokenArrayArray{{1, 2, 3}, {4, 5}}
						return &s
					}(),
				},
				Model: string(openai.OpenAIEmbeddingModelTextEmbedding3Small),
			},
			want: `{"input":[[1,2,3],[4,5]],"model":"text-embedding-3-small"}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("got %s\nwant %s", got, c.want)
			}
		})
	}
}

// === L. Moderation InputUnion + MultiModal marshal 测试（6 case） ===

func TestWireFormatModerationInputUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "InputString",
			v: openai.OpenAIModerationParams{
				Input: &openai.OpenAIModerationInputUnionFlat{
					Type: "string",
					String: func() *openai.OpenAIModerationInputString {
						s := openai.OpenAIModerationInputString("test text")
						return &s
					}(),
				},
				Model: string(openai.OpenAIModerationModelOmniModerationLatest),
			},
			want: `{"input":"test text","model":"omni-moderation-latest"}`,
		},
		{
			name: "InputStringArray",
			v: openai.OpenAIModerationParams{
				Input: &openai.OpenAIModerationInputUnionFlat{
					Type: "string_array",
					StringArray: func() *openai.OpenAIModerationInputStringArray {
						s := openai.OpenAIModerationInputStringArray{"a", "b"}
						return &s
					}(),
				},
				Model: string(openai.OpenAIModerationModelOmniModerationLatest),
			},
			want: `{"input":["a","b"],"model":"omni-moderation-latest"}`,
		},
		{
			name: "InputMultiModalArray_text_only",
			v: openai.OpenAIModerationParams{
				Input: &openai.OpenAIModerationInputUnionFlat{
					Type: "multi_modal_array",
					MultiModalArray: func() *openai.OpenAIModerationInputMultiModalArray {
						s := openai.OpenAIModerationInputMultiModalArray{
							openai.OpenAIModerationMultiModalInputTextParam{
								Type: openai.OpenAIModerationMultiModalInputTypeText,
								Text: "hello",
							},
						}
						return &s
					}(),
				},
				Model: string(openai.OpenAIModerationModelOmniModerationLatest),
			},
			want: `{"input":[{"type":"text","text":"hello"}],"model":"omni-moderation-latest"}`,
		},
		{
			name: "InputMultiModalArray_image_url",
			v: openai.OpenAIModerationParams{
				Input: &openai.OpenAIModerationInputUnionFlat{
					Type: "multi_modal_array",
					MultiModalArray: func() *openai.OpenAIModerationInputMultiModalArray {
						s := openai.OpenAIModerationInputMultiModalArray{
							openai.OpenAIModerationMultiModalInputImageURLParam{
								Type: openai.OpenAIModerationMultiModalInputTypeImageURL,
								ImageURL: openai.OpenAIModerationImageURL{
									URL: "https://example.com/cat.jpg",
								},
							},
						}
						return &s
					}(),
				},
				Model: string(openai.OpenAIModerationModelOmniModerationLatest),
			},
			want: `{"input":[{"type":"image_url","image_url":{"url":"https://example.com/cat.jpg"}}],"model":"omni-moderation-latest"}`,
		},
		{
			name: "InputMultiModalArray_mixed",
			v: openai.OpenAIModerationParams{
				Input: &openai.OpenAIModerationInputUnionFlat{
					Type: "multi_modal_array",
					MultiModalArray: func() *openai.OpenAIModerationInputMultiModalArray {
						s := openai.OpenAIModerationInputMultiModalArray{
							openai.OpenAIModerationMultiModalInputTextParam{
								Type: openai.OpenAIModerationMultiModalInputTypeText,
								Text: "describe",
							},
							openai.OpenAIModerationMultiModalInputImageURLParam{
								Type: openai.OpenAIModerationMultiModalInputTypeImageURL,
								ImageURL: openai.OpenAIModerationImageURL{
									URL: "data:image/jpeg;base64,AAA",
								},
							},
						}
						return &s
					}(),
				},
				Model: string(openai.OpenAIModerationModelOmniModerationLatest),
			},
			want: `{"input":[{"type":"text","text":"describe"},{"type":"image_url","image_url":{"url":"data:image/jpeg;base64,AAA"}}],"model":"omni-moderation-latest"}`,
		},
		{
			name: "InputMultiModal_empty_array",
			v: openai.OpenAIModerationParams{
				Input: &openai.OpenAIModerationInputUnionFlat{
					Type: "multi_modal_array",
					MultiModalArray: func() *openai.OpenAIModerationInputMultiModalArray {
						s := openai.OpenAIModerationInputMultiModalArray{}
						return &s
					}(),
				},
				Model: string(openai.OpenAIModerationModelTextModerationLatest),
			},
			want: `{"input":[],"model":"text-moderation-latest"}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("got %s\nwant %s", got, c.want)
			}
		})
	}
}

// === M. Moderation 响应含 illicit 分类 marshal 测试 ===

func TestWireFormatModerationResponseWithIllicit(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "CategoriesWithIllicit",
			v: openai.OpenAIModerationCategories{
				Hate:           true,
				Illicit:        true,
				IllicitViolent: true,
				Violence:       true,
			},
			want: `{"hate":true,"violence":true,"illicit":true,"illicit/violent":true}`,
		},
		{
			name: "CategoryScoresWithIllicit",
			v: openai.OpenAIModerationCategoryScores{
				Hate:           0.9,
				Illicit:        0.7,
				IllicitViolent: 0.3,
				Violence:       0.5,
			},
			want: `{"hate":0.9,"violence":0.5,"illicit":0.7,"illicit/violent":0.3}`,
		},
		{
			name: "CategoriesEmptyOmitted",
			v:    openai.OpenAIModerationCategories{},
			want: `{}`,
		},
		{
			name: "Moderation_full_with_illicit",
			v: openai.OpenAIModeration{
				Flagged: true,
				Categories: openai.OpenAIModerationCategories{
					Illicit:        true,
					IllicitViolent: true,
				},
				CategoryScores: openai.OpenAIModerationCategoryScores{
					Illicit:        0.85,
					IllicitViolent: 0.65,
				},
			},
			want: `{"flagged":true,"categories":{"illicit":true,"illicit/violent":true},"category_scores":{"illicit":0.85,"illicit/violent":0.65}}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("got %s\nwant %s", got, c.want)
			}
		})
	}
}

// === N. Union 封闭性测试 ===

func TestEmbeddingInputUnionClosure(t *testing.T) {
	var _ openai.OpenAIEmbeddingInputUnion = openai.OpenAIEmbeddingInputString("")
	var _ openai.OpenAIEmbeddingInputUnion = openai.OpenAIEmbeddingInputStringArray{}
	var _ openai.OpenAIEmbeddingInputUnion = openai.OpenAIEmbeddingInputTokenArray{}
	var _ openai.OpenAIEmbeddingInputUnion = openai.OpenAIEmbeddingInputTokenArrayArray{}
}

func TestModerationInputUnionClosure(t *testing.T) {
	var _ openai.OpenAIModerationInputUnion = openai.OpenAIModerationInputString("")
	var _ openai.OpenAIModerationInputUnion = openai.OpenAIModerationInputStringArray{}
	var _ openai.OpenAIModerationInputUnion = openai.OpenAIModerationInputMultiModalArray{}
}

func TestModerationMultiModalInputUnionClosure(t *testing.T) {
	var _ openai.OpenAIModerationMultiModalInputUnionParam = openai.OpenAIModerationMultiModalInputImageURLParam{}
	var _ openai.OpenAIModerationMultiModalInputUnionParam = openai.OpenAIModerationMultiModalInputTextParam{}
}

// === P1: Responses 新增字段测试 ===

func TestWireFormatResponseParamsP1(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "MaxToolCalls",
			v: func() openai.OpenAIResponseParams {
				max := int64(5)
				return openai.OpenAIResponseParams{
					Model:        "gpt-4o",
					Input:        &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
					MaxToolCalls: &max,
				}
			}(),
			want: `{"model":"gpt-4o","input":"hello","max_tool_calls":5}`,
		},
		{
			name: "PreviousResponseID",
			v: openai.OpenAIResponseParams{
				Model:              "gpt-4o",
				Input:              &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
				PreviousResponseID: "resp_abc123",
			},
			want: `{"model":"gpt-4o","input":"hello","previous_response_id":"resp_abc123"}`,
		},
		{
			name: "Prompt_template",
			v: openai.OpenAIResponseParams{
				Model: "gpt-4o",
				Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
				Prompt: &openai.OpenAIResponsePrompt{
					ID:      "tmpl_001",
					Version: "1.0",
				},
			},
			want: `{"model":"gpt-4o","input":"hello","prompt":{"id":"tmpl_001","version":"1.0"}}`,
		},
		{
			name: "Moderation",
			v: openai.OpenAIResponseParams{
				Model:      "gpt-4o",
				Input:      &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
				Moderation: &openai.OpenAIResponseModeration{Model: "omni-moderation-latest"},
			},
			want: `{"model":"gpt-4o","input":"hello","moderation":{"model":"omni-moderation-latest"}}`,
		},
		{
			name: "StreamOptions",
			v: func() openai.OpenAIResponseParams {
				obf := false
				return openai.OpenAIResponseParams{
					Model:         "gpt-4o",
					Input:         &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
					StreamOptions: &openai.OpenAIResponseStreamOptions{IncludeObfuscation: &obf},
				}
			}(),
			want: `{"model":"gpt-4o","input":"hello","stream_options":{"include_obfuscation":false}}`,
		},
		{
			name: "TopLogprobs",
			v: func() openai.OpenAIResponseParams {
				top := int64(5)
				return openai.OpenAIResponseParams{
					Model:       "gpt-4o",
					Input:       &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
					TopLogprobs: &top,
				}
			}(),
			want: `{"model":"gpt-4o","input":"hello","top_logprobs":5}`,
		},
		{
			name: "ContextManagement",
			v: func() openai.OpenAIResponseParams {
				threshold := int64(10000)
				return openai.OpenAIResponseParams{
					Model: "gpt-4o",
					Input: &openai.OpenAIResponseInputUnionFlat{Type: "string", Text: "hello"},
					ContextManagement: []openai.OpenAIResponseContextManagement{
						{Type: "compaction", CompactThreshold: &threshold},
					},
				}
			}(),
			want: `{"model":"gpt-4o","input":"hello","context_management":[{"type":"compaction","compact_threshold":10000}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestWireFormatImageVariationParams(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "minimal",
			v: openai.OpenAIImageVariationParams{
				Image: "base64data",
			},
			want: `{"image":"base64data"}`,
		},
		{
			name: "full_with_gpt_image_1",
			v: func() openai.OpenAIImageVariationParams {
				n := int64(1)
				return openai.OpenAIImageVariationParams{
					Image:          "base64data",
					Model:          openai.OpenAIImageModelGPTImage1,
					N:              &n,
					ResponseFormat: "url",
					Size:           "1024x1024",
				}
			}(),
			want: `{"image":"base64data","model":"gpt-image-1","n":1,"response_format":"url","size":"1024x1024"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestWireFormatImageGenerateParamsP1(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "gpt_image_1_with_background",
			v: func() openai.OpenAIImageGenerateParams {
				bg := true
				q := int64(85)
				return openai.OpenAIImageGenerateParams{
					Prompt:                   "A cute cat",
					Model:                    openai.OpenAIImageModelGPTImage1,
					Background:               &bg,
					OutputCompressionQuality: &q,
					PromptRefinement:         "high",
				}
			}(),
			want: `{"prompt":"A cute cat","model":"gpt-image-1","background":true,"output_compression_quality":85,"prompt_refinement":"high"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestImageModelConstants(t *testing.T) {
	models := []openai.OpenAIImageModel{
		openai.OpenAIImageModelGPTImage1,
		openai.OpenAIImageModelGPTImage1Mini,
		openai.OpenAIImageModelGPTImage1_5,
		openai.OpenAIImageModelGPTImage2,
		openai.OpenAIImageModelGPTImage2_2026_04_21,
		openai.OpenAIImageModelChatgptImageLatest,
		openai.OpenAIImageModelDallE2,
		openai.OpenAIImageModelDallE3,
	}
	want := []string{
		"gpt-image-1", "gpt-image-1-mini", "gpt-image-1.5",
		"gpt-image-2", "gpt-image-2-2026-04-21",
		"chatgpt-image-latest", "dall-e-2", "dall-e-3",
	}
	for i, m := range models {
		if string(m) != want[i] {
			t.Errorf("models[%d] = %q, want %q", i, string(m), want[i])
		}
	}
}

// === P2: InputMessageContentPartUnion 强类型化 ===

func TestWireFormatInputMessageContentPartUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "input_text",
			v: openai.OpenAIResponseInputMessageContentPartText{
				Type: "input_text",
				Text: "Hello world",
			},
			want: `{"type":"input_text","text":"Hello world"}`,
		},
		{
			name: "input_image",
			v: openai.OpenAIResponseInputMessageContentPartImage{
				Type:     "input_image",
				ImageURL: "https://example.com/image.png",
				Detail:   "high",
			},
			want: `{"type":"input_image","image_url":"https://example.com/image.png","detail":"high"}`,
		},
		{
			name: "input_file",
			v: openai.OpenAIResponseInputMessageContentPartFile{
				Type:   "input_file",
				FileID: "file_abc123",
			},
			want: `{"type":"input_file","file_id":"file_abc123"}`,
		},
		{
			name: "InputMessageContentParts_array",
			v: openai.OpenAIResponseInputMessageContentParts{
				openai.OpenAIResponseInputMessageContentPartText{Type: "input_text", Text: "Hello"},
				openai.OpenAIResponseInputMessageContentPartImage{Type: "input_image", ImageURL: "https://example.com/img.png"},
				openai.OpenAIResponseInputMessageContentPartFile{Type: "input_file", FileID: "file_xyz"},
			},
			want: `[{"type":"input_text","text":"Hello"},{"type":"input_image","image_url":"https://example.com/img.png"},{"type":"input_file","file_id":"file_xyz"}]`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestInputMessageContentPartUnionClosure(t *testing.T) {
	var _ openai.OpenAIResponseInputMessageContentPartUnion = openai.OpenAIResponseInputMessageContentPartText{}
	var _ openai.OpenAIResponseInputMessageContentPartUnion = openai.OpenAIResponseInputMessageContentPartImage{}
	var _ openai.OpenAIResponseInputMessageContentPartUnion = openai.OpenAIResponseInputMessageContentPartFile{}
}

// === P3: CompactionItem + AdditionalTools + Logprobs ===

func TestWireFormatResponseOutputItemP3(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "CompactionItem_with_created_by",
			v: openai.OpenAIResponseOutputCompactionItem{
				Type:             "compaction",
				ID:               "comp_001",
				EncryptedContent: "abc123",
				CreatedBy:        "developer",
			},
			want: `{"type":"compaction","id":"comp_001","encrypted_content":"abc123","created_by":"developer"}`,
		},
		{
			name: "AdditionalTools_output",
			v: openai.OpenAIResponseOutputAdditionalTools{
				Type: "additional_tools",
				ID:   "at_001",
				Role: "developer",
				Tools: []openai.OpenAIResponsesToolUnionFlat{
					openai.OpenAIResponsesToolUnionFlat{Type: "function", Name: "get_weather"},
				},
			},
			want: `{"type":"additional_tools","id":"at_001","role":"developer","tools":[{"type":"function","name":"get_weather"}]}`,
		},
		{
			name: "FunctionWebSearch_output",
			v: openai.OpenAIResponseOutputFunctionWebSearch{
				Type:   "web_search_call",
				ID:     "ws_001",
				Status: "completed",
				Action: &openai.OpenAIResponseInputWebSearchAction{Type: "search", Query: "weather"},
			},
			want: `{"type":"web_search_call","id":"ws_001","status":"completed","action":{"type":"search","query":"weather"}}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestWireFormatInputItemP3(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "AdditionalTools_input",
			v: openai.OpenAIResponseInputAdditionalToolsItem{
				Type: "additional_tools",
				ID:   "at_001",
				Role: "developer",
				Tools: []openai.OpenAIResponsesToolUnionFlat{
					openai.OpenAIResponsesToolUnionFlat{Type: "function", Name: "get_weather"},
				},
			},
			want: `{"type":"additional_tools","id":"at_001","role":"developer","tools":[{"type":"function","name":"get_weather"}]}`,
		},
		{
			name: "CustomToolCall_input",
			v: openai.OpenAIResponseInputCustomToolCallItem{
				Type:   "custom_tool_call",
				ID:     "ct_001",
				CallID: "call_001",
				Name:   "my_tool",
				Input:  `{"arg":"value"}`,
				Status: "in_progress",
			},
			want: `{"type":"custom_tool_call","id":"ct_001","call_id":"call_001","name":"my_tool","input":"{\"arg\":\"value\"}","status":"in_progress"}`,
		},
		{
			name: "CustomToolCallOutput_input",
			v: openai.OpenAIResponseInputCustomToolCallOutputItem{
				Type:   "custom_tool_call_output",
				ID:     "cto_001",
				CallID: "call_001",
				Output: `{"result":"ok"}`,
			},
			want: `{"type":"custom_tool_call_output","id":"cto_001","call_id":"call_001","output":"{\"result\":\"ok\"}"}`,
		},
		{
			name: "FunctionCallOutput_with_text",
			v: openai.OpenAIResponseInputFunctionCallOutputItem{
				Type:   "function_call_output",
				ID:     "fco_001",
				CallID: "call_001",
				Output: &openai.OpenAIResponseInputFunctionCallOutputUnionFlat{Type: "output_text", Text: "sunny"},
				Status: "completed",
			},
			want: `{"type":"function_call_output","id":"fco_001","call_id":"call_001","output":{"type":"output_text","text":"sunny"},"status":"completed"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

func TestOutputMessageContentLogprobs(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "with_logprobs",
			v: openai.OpenAIResponseOutputMessageContent{
				Type: "output_text",
				Text: "Hello",
				Logprobs: []openai.OpenAIResponseOutputMessageContentLogprob{
					{
						Token:   "Hello",
						Bytes:   []int64{72, 101},
						Logprob: -0.5,
						TopLogprobs: []openai.OpenAIResponseOutputMessageContentLogprobTopLogprob{
							{Token: "Hi", Bytes: []int64{72, 105}, Logprob: -0.6},
						},
					},
				},
			},
			want: `{"type":"output_text","text":"Hello","logprobs":[{"token":"Hello","bytes":[72,101],"logprob":-0.5,"top_logprobs":[{"token":"Hi","bytes":[72,105],"logprob":-0.6}]}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}
