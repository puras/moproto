package anthropic_test

import (
	"encoding/json"
	"strings"
	"testing"

	anthropic "hegq.net/mogo/pkg/protocol/anthropic"
)

// TestWireFormatContentBlocks 验证 11 个内容块的 JSON marshal 输出与重构前一致。
//
// 重构原则：Type 字段从 string 改为独立 string alias（如 AnthropicTextBlockType），
// 但 wire-format 字符串值（"text" / "tool_use" 等）保持不变。本测试锁住该契约。
func TestWireFormatContentBlocks(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "TextBlock_minimal",
			v: anthropic.AnthropicTextBlockParam{
				Type: anthropic.AnthropicTextBlockTypeText,
				Text: "hello",
			},
			want: `{"type":"text","text":"hello"}`,
		},
		{
			name: "TextBlock_with_cache_control",
			v: anthropic.AnthropicTextBlockParam{
				Type:         anthropic.AnthropicTextBlockTypeText,
				Text:         "hello",
				CacheControl: &anthropic.AnthropicCacheControlEphemeral{Type: anthropic.AnthropicCacheControlEphemeralTypeEphemeral, TTL: "5m"},
			},
			want: `{"type":"text","text":"hello","cache_control":{"type":"ephemeral","ttl":"5m"}}`,
		},
		{
			name: "ImageBlock_base64",
			v: anthropic.AnthropicImageBlockParam{
				Type:   anthropic.AnthropicImageBlockTypeImage,
				Source: anthropic.NewAnthropicImageSourceBase64("image/jpeg", "base64data"),
			},
			want: `{"type":"image","source":{"type":"base64","media_type":"image/jpeg","data":"base64data"}}`,
		},
		{
			name: "DocumentBlock_url",
			v: anthropic.AnthropicDocumentBlockParam{
				Type:   anthropic.AnthropicDocumentBlockTypeDocument,
				Source: anthropic.NewAnthropicDocumentSourceURL("https://example.com/doc.pdf"),
			},
			want: `{"type":"document","source":{"type":"url","url":"https://example.com/doc.pdf"}}`,
		},
		{
			name: "ToolUseBlock",
			v: anthropic.AnthropicToolUseBlockParam{
				Type:  anthropic.AnthropicToolUseBlockTypeToolUse,
				ID:    "tu_1",
				Name:  "get_weather",
				Input: json.RawMessage(`{"city":"SF"}`),
			},
			want: `{"type":"tool_use","id":"tu_1","name":"get_weather","input":{"city":"SF"}}`,
		},
		{
			name: "ToolResultBlock_with_error",
			v: anthropic.AnthropicToolResultBlockParam{
				Type:      anthropic.AnthropicToolResultBlockTypeToolResult,
				ToolUseID: "tu_1",
				Content:   json.RawMessage(`"not found"`),
				IsError:   true,
			},
			want: `{"type":"tool_result","tool_use_id":"tu_1","content":"not found","is_error":true}`,
		},
		{
			name: "ThinkingBlock",
			v: anthropic.AnthropicThinkingBlockParam{
				Type:      anthropic.AnthropicThinkingBlockTypeThinking,
				Thinking:  "let me think...",
				Signature: "sig_xxx",
			},
			want: `{"type":"thinking","thinking":"let me think...","signature":"sig_xxx"}`,
		},
		{
			name: "RedactedThinkingBlock",
			v: anthropic.AnthropicRedactedThinkingBlockParam{
				Type: anthropic.AnthropicRedactedThinkingBlockTypeRedactedThinking,
				Data: "encrypted_blob",
			},
			want: `{"type":"redacted_thinking","data":"encrypted_blob"}`,
		},
		{
			name: "MidConvSystemBlock",
			v: anthropic.AnthropicMidConversationSystemBlockParam{
				Type: anthropic.AnthropicMidConversationSystemBlockTypeMidConvSystem,
				Content: []anthropic.AnthropicTextBlockParam{
					{Type: anthropic.AnthropicTextBlockTypeText, Text: "remember X"},
				},
			},
			want: `{"type":"mid_conv_system","content":[{"type":"text","text":"remember X"}]}`,
		},
		{
			name: "ServerToolUseBlock_web_search",
			v: anthropic.AnthropicServerToolUseBlockParam{
				Type:  anthropic.AnthropicServerToolUseBlockTypeServerToolUse,
				ID:    "stu_1",
				Name:  anthropic.ServerToolUseBlockNameWebSearch,
				Input: json.RawMessage(`{"query":"weather"}`),
			},
			want: `{"type":"server_tool_use","id":"stu_1","name":"web_search","input":{"query":"weather"}}`,
		},
		{
			name: "WebSearchToolResultBlock",
			v: anthropic.AnthropicWebSearchToolResultBlockParam{
				Type:      anthropic.AnthropicWebSearchToolResultBlockTypeWebSearchToolResult,
				ToolUseID: "stu_1",
				Content:   json.RawMessage(`[{"type":"web_search_result","url":"https://x.com"}]`),
			},
			want: `{"type":"web_search_tool_result","tool_use_id":"stu_1","content":[{"type":"web_search_result","url":"https://x.com"}]}`,
		},
		{
			name: "ToolReferenceBlock",
			v: anthropic.AnthropicToolReferenceBlockParam{
				Type:     anthropic.AnthropicToolReferenceBlockTypeToolReference,
				ToolName: "calculator",
			},
			want: `{"type":"tool_reference","tool_name":"calculator"}`,
		},
		{
			name: "CompactionBlock",
			v: anthropic.AnthropicCompactionBlockParam{
				Type:    anthropic.AnthropicCompactionBlockTypeCompaction,
				Content: json.RawMessage(`"summary text"`),
			},
			want: `{"type":"compaction","content":"summary text"}`,
		},
		{
			name: "SearchResultBlock",
			v: anthropic.AnthropicSearchResultBlockParam{
				Type:   anthropic.AnthropicSearchResultBlockTypeSearchResult,
				Source: "https://example.com",
				Title:  "Example",
				Content: []anthropic.AnthropicTextBlockParam{
					{Type: anthropic.AnthropicTextBlockTypeText, Text: "matched"},
				},
			},
			want: `{"type":"search_result","source":"https://example.com","title":"Example","content":[{"type":"text","text":"matched"}]}`,
		},
		{
			name: "WebFetchToolResultBlock",
			v: anthropic.AnthropicWebFetchToolResultBlockParam{
				Type:      anthropic.AnthropicWebFetchToolResultBlockTypeWebFetchToolResult,
				ToolUseID: "ftu_1",
				Content:   json.RawMessage(`{"url":"https://example.com","content":"fetched text"}`),
			},
			want: `{"type":"web_fetch_tool_result","tool_use_id":"ftu_1","content":{"url":"https://example.com","content":"fetched text"}}`,
		},
		{
			name: "CodeExecutionToolResultBlock",
			v: anthropic.AnthropicCodeExecutionToolResultBlockParam{
				Type:      anthropic.AnthropicCodeExecutionToolResultBlockTypeCodeExecutionToolResult,
				ToolUseID: "ctu_1",
				Content:   json.RawMessage(`{"stdout":"hi\n"}`),
			},
			want: `{"type":"code_execution_tool_result","tool_use_id":"ctu_1","content":{"stdout":"hi\n"}}`,
		},
		{
			name: "BashCodeExecutionToolResultBlock",
			v: anthropic.AnthropicBashCodeExecutionToolResultBlockParam{
				Type:      anthropic.AnthropicBashCodeExecutionToolResultBlockTypeBashCodeExecutionToolResult,
				ToolUseID: "btu_1",
				Content:   json.RawMessage(`{"stdout":"hi","stderr":"","return_code":0}`),
			},
			want: `{"type":"bash_code_execution_tool_result","tool_use_id":"btu_1","content":{"stdout":"hi","stderr":"","return_code":0}}`,
		},
		{
			name: "TextEditorCodeExecutionToolResultBlock",
			v: anthropic.AnthropicTextEditorCodeExecutionToolResultBlockParam{
				Type:      anthropic.AnthropicTextEditorCodeExecutionToolResultBlockTypeTextEditorCodeExecutionToolResult,
				ToolUseID: "ttu_1",
				Content:   json.RawMessage(`{"type":"text_editor_code_execution_tool_result","content":"ok"}`),
			},
			want: `{"type":"text_editor_code_execution_tool_result","tool_use_id":"ttu_1","content":{"type":"text_editor_code_execution_tool_result","content":"ok"}}`,
		},
		{
			name: "ToolSearchToolResultBlock",
			v: anthropic.AnthropicToolSearchToolResultBlockParam{
				Type:      anthropic.AnthropicToolSearchToolResultBlockTypeToolSearchToolResult,
				ToolUseID: "ttu_2",
				Content:   json.RawMessage(`{"results":["calculator"]}`),
			},
			want: `{"type":"tool_search_tool_result","tool_use_id":"ttu_2","content":{"results":["calculator"]}}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("wire-format mismatch:\n  got:  %s\n  want: %s", got, c.want)
			}
		})
	}
}

// TestWireFormatTools 验证 7 个 Tool 实现的 JSON marshal 输出。
func TestWireFormatTools(t *testing.T) {
	cases := []struct {
		name string
		v    anthropic.AnthropicTool
		want string
	}{
		{
			name: "CustomTool_minimal",
			v: anthropic.AnthropicCustomTool{
				Type:        anthropic.AnthropicCustomToolTypeCustom,
				Name:        "get_weather",
				InputSchema: anthropic.AnthropicToolInputSchema{Type: "object"},
			},
			want: `{"type":"custom","name":"get_weather","input_schema":{"type":"object"}}`,
		},
		{
			name: "WebSearchTool",
			v: anthropic.AnthropicWebSearchTool{
				Type: anthropic.AnthropicWebSearchToolType20260209,
				Name: "web_search",
			},
			want: `{"type":"web_search_20260209","name":"web_search"}`,
		},
		{
			name: "WebFetchTool",
			v: anthropic.AnthropicWebFetchTool{
				Type: anthropic.AnthropicWebFetchToolType20260209,
				Name: "web_fetch",
			},
			want: `{"type":"web_fetch_20260209","name":"web_fetch"}`,
		},
		{
			name: "CodeExecutionTool",
			v: anthropic.AnthropicCodeExecutionTool{
				Type: anthropic.AnthropicCodeExecutionToolType20260120,
				Name: "code_execution",
			},
			want: `{"type":"code_execution_20260120","name":"code_execution"}`,
		},
		{
			name: "BashTool",
			v: anthropic.AnthropicBashTool{
				Type: anthropic.AnthropicBashToolType20250124,
				Name: "bash",
			},
			want: `{"type":"bash_20250124","name":"bash"}`,
		},
		{
			name: "TextEditorTool",
			v: anthropic.AnthropicTextEditorTool{
				Type: anthropic.AnthropicTextEditorToolType20250728,
				Name: "str_replace_based_edit_tool",
			},
			want: `{"type":"text_editor_20250728","name":"str_replace_based_edit_tool"}`,
		},
		{
			name: "MemoryTool",
			v: anthropic.AnthropicMemoryTool{
				Type: anthropic.AnthropicMemoryToolType20250818,
				Name: "memory",
			},
			want: `{"type":"memory_20250818","name":"memory"}`,
		},
		// ---- P0 补齐：12 个新 tool 类型 ----
		{
			name: "CodeExecutionTool20250522",
			v: anthropic.AnthropicCodeExecutionTool20250522{
				Type: anthropic.AnthropicCodeExecutionToolType20250522,
				Name: "code_execution",
			},
			want: `{"type":"code_execution_20250522","name":"code_execution"}`,
		},
		{
			name: "CodeExecutionTool20250825",
			v: anthropic.AnthropicCodeExecutionTool20250825{
				Type: anthropic.AnthropicCodeExecutionToolType20250825,
				Name: "code_execution",
			},
			want: `{"type":"code_execution_20250825","name":"code_execution"}`,
		},
		{
			name: "CodeExecutionTool20260521",
			v: anthropic.AnthropicCodeExecutionTool20260521{
				Type: anthropic.AnthropicCodeExecutionToolType20260521,
				Name: "code_execution",
			},
			want: `{"type":"code_execution_20260521","name":"code_execution"}`,
		},
		{
			name: "WebSearchTool20250305",
			v: anthropic.AnthropicWebSearchTool20250305{
				Type: anthropic.AnthropicWebSearchToolType20250305,
				Name: "web_search",
			},
			want: `{"type":"web_search_20250305","name":"web_search"}`,
		},
		{
			name: "WebSearchTool20260318_with_response_inclusion",
			v: anthropic.AnthropicWebSearchTool20260318{
				Type:              anthropic.AnthropicWebSearchToolType20260318,
				Name:              "web_search",
				ResponseInclusion: anthropic.AnthropicWebSearchToolResponseInclusionExcluded,
			},
			want: `{"type":"web_search_20260318","name":"web_search","response_inclusion":"excluded"}`,
		},
		{
			name: "WebFetchTool20250910",
			v: anthropic.AnthropicWebFetchTool20250910{
				Type: anthropic.AnthropicWebFetchToolType20250910,
				Name: "web_fetch",
			},
			want: `{"type":"web_fetch_20250910","name":"web_fetch"}`,
		},
		{
			name: "WebFetchTool20260309_with_use_cache",
			v: anthropic.AnthropicWebFetchTool20260309{
				Type:     anthropic.AnthropicWebFetchToolType20260309,
				Name:     "web_fetch",
				UseCache: true,
			},
			want: `{"type":"web_fetch_20260309","name":"web_fetch","use_cache":true}`,
		},
		{
			name: "WebFetchTool20260318",
			v: anthropic.AnthropicWebFetchTool20260318{
				Type: anthropic.AnthropicWebFetchToolType20260318,
				Name: "web_fetch",
			},
			want: `{"type":"web_fetch_20260318","name":"web_fetch"}`,
		},
		{
			name: "ToolSearchToolBm25",
			v: anthropic.AnthropicToolSearchToolBm25{
				Type: anthropic.AnthropicToolSearchToolBm25Type20251119,
				Name: "tool_search_tool_bm25",
			},
			want: `{"type":"tool_search_tool_bm25_20251119","name":"tool_search_tool_bm25"}`,
		},
		{
			name: "ToolSearchToolRegex",
			v: anthropic.AnthropicToolSearchToolRegex{
				Type: anthropic.AnthropicToolSearchToolRegexType20251119,
				Name: "tool_search_tool_regex",
			},
			want: `{"type":"tool_search_tool_regex_20251119","name":"tool_search_tool_regex"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("wire-format mismatch:\n  got:  %s\n  want: %s", got, c.want)
			}
		})
	}
}

// TestWireFormatToolChoices 验证 4 个 ToolChoice 实现的 JSON marshal 输出。
func TestWireFormatToolChoices(t *testing.T) {
	cases := []struct {
		name string
		v    anthropic.AnthropicToolChoice
		want string
	}{
		{
			name: "Auto",
			v:    anthropic.AnthropicToolChoiceAuto{Type: anthropic.AnthropicToolChoiceAutoTypeAuto},
			want: `{"type":"auto"}`,
		},
		{
			name: "Any_disable_parallel",
			v: anthropic.AnthropicToolChoiceAny{
				Type:                   anthropic.AnthropicToolChoiceAnyTypeAny,
				DisableParallelToolUse: true,
			},
			want: `{"type":"any","disable_parallel_tool_use":true}`,
		},
		{
			name: "Tool",
			v: anthropic.AnthropicToolChoiceTool{
				Type: anthropic.AnthropicToolChoiceToolTypeTool,
				Name: "get_weather",
			},
			want: `{"type":"tool","name":"get_weather"}`,
		},
		{
			name: "None",
			v:    anthropic.AnthropicToolChoiceNone{Type: anthropic.AnthropicToolChoiceNoneTypeNone},
			want: `{"type":"none"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("wire-format mismatch:\n  got:  %s\n  want: %s", got, c.want)
			}
		})
	}
}

// TestWireFormatContextEdits 验证 3 个 ContextEdit 实现的 JSON marshal 输出。
func TestWireFormatContextEdits(t *testing.T) {
	cases := []struct {
		name string
		v    anthropic.AnthropicContextEdit
		want string
	}{
		{
			name: "ClearToolUses",
			v: anthropic.AnthropicClearToolUses20250919Edit{
				Type:             anthropic.AnthropicClearToolUses20250919EditTypeClearToolUses20250919,
				ClearToolInputs:  true,
				MinTokensToClear: 1000,
			},
			want: `{"type":"clear_tool_uses_20250919","clear_tool_inputs":true,"min_tokens_to_clear":1000}`,
		},
		{
			name: "ClearThinking",
			v:    anthropic.AnthropicClearThinking20251015Edit{Type: anthropic.AnthropicClearThinking20251015EditTypeClearThinking20251015},
			want: `{"type":"clear_thinking_20251015"}`,
		},
		{
			name: "Compact",
			v:    anthropic.AnthropicCompact20260112Edit{Type: anthropic.AnthropicCompact20260112EditTypeCompact20260112},
			want: `{"type":"compact_20260112"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("wire-format mismatch:\n  got:  %s\n  want: %s", got, c.want)
			}
		})
	}
}

// TestUnionSealedImplementation 验证 union 接口的封闭性。
//
// AnthropicContentBlockParam 只能由本包内的类型实现（isContentBlockParam 是私有方法）。
// 本测试通过 type switch 验证接口断言通过——证明接口签名本身有效。
func TestUnionSealedImplementation(t *testing.T) {
	// 编译期验证：所有实现都满足接口
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicTextBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicImageBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicDocumentBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicToolUseBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicToolResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicThinkingBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicRedactedThinkingBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicMidConversationSystemBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicServerToolUseBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicWebSearchToolResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicToolReferenceBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicCompactionBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicSearchResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicWebFetchToolResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicCodeExecutionToolResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicBashCodeExecutionToolResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicTextEditorCodeExecutionToolResultBlockParam{}
	var _ anthropic.AnthropicContentBlockParam = anthropic.AnthropicToolSearchToolResultBlockParam{}

	// ---- P0 补齐：ToolResultBlockContentUnion 5 个具体实现 ----
	var _ anthropic.AnthropicToolResultBlockContentUnion = anthropic.AnthropicToolResultBlockContentText{}
	var _ anthropic.AnthropicToolResultBlockContentUnion = anthropic.AnthropicToolResultBlockContentImage{}
	var _ anthropic.AnthropicToolResultBlockContentUnion = anthropic.AnthropicToolResultBlockContentSearchResult{}
	var _ anthropic.AnthropicToolResultBlockContentUnion = anthropic.AnthropicToolResultBlockContentDocument{}
	var _ anthropic.AnthropicToolResultBlockContentUnion = anthropic.AnthropicToolResultBlockContentToolReference{}

	var _ anthropic.AnthropicTool = anthropic.AnthropicCustomTool{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebSearchTool{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebFetchTool{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicCodeExecutionTool{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicBashTool{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicTextEditorTool{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicMemoryTool{}
	// ---- P0 补齐：12 个新 tool 实现 ----
	var _ anthropic.AnthropicTool = anthropic.AnthropicCodeExecutionTool20250522{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicCodeExecutionTool20250825{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicCodeExecutionTool20260521{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebSearchTool20250305{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebSearchTool20260318{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebFetchTool20250910{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebFetchTool20260309{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicWebFetchTool20260318{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicToolSearchToolBm25{}
	var _ anthropic.AnthropicTool = anthropic.AnthropicToolSearchToolRegex{}

	var _ anthropic.AnthropicToolChoice = anthropic.AnthropicToolChoiceAuto{}
	var _ anthropic.AnthropicToolChoice = anthropic.AnthropicToolChoiceAny{}
	var _ anthropic.AnthropicToolChoice = anthropic.AnthropicToolChoiceTool{}
	var _ anthropic.AnthropicToolChoice = anthropic.AnthropicToolChoiceNone{}

	var _ anthropic.AnthropicContextEdit = anthropic.AnthropicClearToolUses20250919Edit{}
	var _ anthropic.AnthropicContextEdit = anthropic.AnthropicClearThinking20251015Edit{}
	var _ anthropic.AnthropicContextEdit = anthropic.AnthropicCompact20260112Edit{}

	var _ anthropic.AnthropicThinkingConfig = anthropic.AnthropicThinkingConfigAdaptive{}
	var _ anthropic.AnthropicThinkingConfig = anthropic.AnthropicThinkingConfigEnabled{}
	var _ anthropic.AnthropicThinkingConfig = anthropic.AnthropicThinkingConfigDisabled{}

	var _ anthropic.AnthropicCitationsParam = anthropic.AnthropicCitationCharLocationParam{}
	var _ anthropic.AnthropicCitationsParam = anthropic.AnthropicCitationPageLocationParam{}
	var _ anthropic.AnthropicCitationsParam = anthropic.AnthropicCitationContentBlockLocationParam{}
	var _ anthropic.AnthropicCitationsParam = anthropic.AnthropicCitationSearchResultLocationParam{}
	var _ anthropic.AnthropicCitationsParam = anthropic.AnthropicCitationWebSearchResultLocationParam{}

	// type switch 验证分发机制
	blocks := []anthropic.AnthropicContentBlockParam{
		anthropic.AnthropicTextBlockParam{Type: anthropic.AnthropicTextBlockTypeText, Text: "hi"},
		anthropic.AnthropicImageBlockParam{Type: anthropic.AnthropicImageBlockTypeImage},
	}
	for _, b := range blocks {
		// 通过 Type 字段分发（替代旧的 Kind() 方法）
		var typeStr string
		switch v := b.(type) {
		case anthropic.AnthropicTextBlockParam:
			typeStr = string(v.Type)
		case anthropic.AnthropicImageBlockParam:
			typeStr = string(v.Type)
		default:
			t.Errorf("unexpected block type %T", b)
		}
		if typeStr == "" {
			t.Errorf("empty type string for %T", b)
		}
	}
}

// TestRequestGetModel 验证顶层 Request 的 GetModel 方法（与 common 包接口对齐）。
func TestRequestGetModel(t *testing.T) {
	r := &anthropic.AnthropicRequest{Model: "claude-fable-5"}
	if got := r.GetModel(); got != "claude-fable-5" {
		t.Errorf("GetModel = %q, want claude-fable-5", got)
	}
}

// TestRequestMarshalSample 验证 Request 的完整 marshal 输出。
func TestRequestMarshalSample(t *testing.T) {
	req := anthropic.AnthropicRequest{
		Model:     "claude-fable-5",
		MaxTokens: 1024,
		Messages: []anthropic.AnthropicMessageParam{
			{
				Role:    anthropic.AnthropicMessageRoleUser,
				Content: json.RawMessage(`"hi"`),
			},
		},
		System:      json.RawMessage(`"you are helpful"`),
		Stream:      true,
		ServiceTier: anthropic.ServiceTierAuto,
		Speed:       anthropic.SpeedFast,
	}
	got, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"model":"claude-fable-5","max_tokens":1024,"messages":[{"role":"user","content":"hi"}],"system":"you are helpful","stream":true,"service_tier":"auto","speed":"fast"}`
	if string(got) != want {
		t.Errorf("Request marshal mismatch:\n  got:  %s\n  want: %s", got, want)
	}
}

// TestWireFormatThinkingConfigs 验证 AnthropicThinkingConfig union 3 个实现的 JSON marshal。
//
// 覆盖 adaptive / enabled / disabled 三种 mode 的 wire-format 输出。
func TestWireFormatThinkingConfigs(t *testing.T) {
	cases := []struct {
		name string
		v    anthropic.AnthropicThinkingConfig
		want string
	}{
		{
			name: "Adaptive_minimal",
			v: anthropic.AnthropicThinkingConfigAdaptive{
				Type: anthropic.AnthropicThinkingConfigAdaptiveTypeAdaptive,
			},
			want: `{"type":"adaptive"}`,
		},
		{
			name: "Adaptive_with_display",
			v: anthropic.AnthropicThinkingConfigAdaptive{
				Type:    anthropic.AnthropicThinkingConfigAdaptiveTypeAdaptive,
				Display: anthropic.ThinkingDisplaySummarized,
			},
			want: `{"type":"adaptive","display":"summarized"}`,
		},
		{
			name: "Enabled_with_budget",
			v: anthropic.AnthropicThinkingConfigEnabled{
				Type:         anthropic.AnthropicThinkingConfigEnabledTypeEnabled,
				BudgetTokens: 2048,
			},
			want: `{"type":"enabled","budget_tokens":2048}`,
		},
		{
			name: "Enabled_with_display",
			v: anthropic.AnthropicThinkingConfigEnabled{
				Type:         anthropic.AnthropicThinkingConfigEnabledTypeEnabled,
				BudgetTokens: 4096,
				Display:      anthropic.ThinkingConfigEnabledDisplayOmitted,
			},
			want: `{"type":"enabled","budget_tokens":4096,"display":"omitted"}`,
		},
		{
			name: "Disabled",
			v: anthropic.AnthropicThinkingConfigDisabled{
				Type: anthropic.AnthropicThinkingConfigDisabledTypeDisabled,
			},
			want: `{"type":"disabled"}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s: marshal mismatch:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestWireFormatCustomToolInputSchema 验证 AnthropicToolInputSchema 的结构化序列化与
// Properties / ExtraFields 字段捕获行为。
//
// 因 Go map 序列化顺序不固定，复杂用例（Properties / ExtraFields 含 map）改用
// wantContains 子串断言；简单用例（仅标准字段）仍用全等断言。
func TestWireFormatCustomToolInputSchema(t *testing.T) {
	cases := []struct {
		name         string
		v            anthropic.AnthropicCustomTool
		want         string
		wantContains []string
	}{
		{
			name: "minimal_type_only",
			v: anthropic.AnthropicCustomTool{
				Type:        anthropic.AnthropicCustomToolTypeCustom,
				Name:        "get_weather",
				InputSchema: anthropic.AnthropicToolInputSchema{Type: "object"},
			},
			want: `{"type":"custom","name":"get_weather","input_schema":{"type":"object"}}`,
		},
		{
			name: "with_properties_and_required",
			v: anthropic.AnthropicCustomTool{
				Type: anthropic.AnthropicCustomToolTypeCustom,
				Name: "search",
				InputSchema: anthropic.AnthropicToolInputSchema{
					Type: "object",
					Properties: map[string]any{
						"city": map[string]any{"type": "string"},
						"unit": map[string]any{"type": "string", "enum": []string{"c", "f"}},
					},
					Required: []string{"city"},
				},
			},
			// Go map 序列化顺序不固定；测试只断言 wire-format 完整性（unmarshal 后字段一致）。
			wantContains: []string{
				`"type":"custom"`,
				`"name":"search"`,
				`"input_schema":{`,
				`"type":"object"`,
				`"properties":{`,
				`"city":{"type":"string"}`,
				`"unit":{`,
				`"type":"string"`,
				`"enum":["c","f"]`,
				`"required":["city"]`,
			},
		},
		{
			name: "with_extra_fields",
			v: anthropic.AnthropicCustomTool{
				Type: anthropic.AnthropicCustomToolTypeCustom,
				Name: "f",
				InputSchema: anthropic.AnthropicToolInputSchema{
					Type: "object",
					ExtraFields: map[string]any{
						"$schema": "https://json-schema.org/draft/2020-12/schema",
					},
				},
			},
			// Go map 序列化顺序不固定；ExtraFields 应与标准字段一起出现在 input_schema 内。
			wantContains: []string{
				`"type":"custom"`,
				`"name":"f"`,
				`"input_schema":{`,
				`"type":"object"`,
				`"$schema":"https://json-schema.org/draft/2020-12/schema"`,
			},
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal: %v", c.name, err)
			continue
		}
		if c.want != "" {
			if string(got) != c.want {
				t.Errorf("%s: marshal mismatch:\n  got:  %s\n  want: %s", c.name, got, c.want)
			}
		}
		for _, sub := range c.wantContains {
			if !strings.Contains(string(got), sub) {
				t.Errorf("%s: marshal missing substring %q\n  got: %s", c.name, sub, got)
			}
		}
	}
}

// TestWireFormatCitationsParams 验证 5 种 AnthropicCitationsParam 实现的 JSON marshal。
//
// 对齐官方 anthropic-sdk-go Citation*LocationParam（message.go:582-783）。
func TestWireFormatCitationsParams(t *testing.T) {
	idx0 := int64(0)
	idx5 := int64(5)
	title := "doc title"
	fid := "file_xxx"
	src := "src"
	u := "https://example.com"

	cases := []struct {
		name string
		v    anthropic.AnthropicCitationsParam
		want string
	}{
		{
			name: "CharLocation_minimal",
			v: anthropic.AnthropicCitationCharLocationParam{
				Type:           anthropic.AnthropicCitationCharLocationTypeCharLocation,
				CitedText:      "hello",
				StartCharIndex: 0,
				EndCharIndex:   5,
			},
			want: `{"type":"char_location","cited_text":"hello","start_char_index":0,"end_char_index":5}`,
		},
		{
			name: "CharLocation_with_doc",
			v: anthropic.AnthropicCitationCharLocationParam{
				Type:           anthropic.AnthropicCitationCharLocationTypeCharLocation,
				CitedText:      "hello",
				DocumentIndex:  &idx0,
				DocumentTitle:  &title,
				StartCharIndex: 10,
				EndCharIndex:   15,
			},
			want: `{"type":"char_location","cited_text":"hello","document_index":0,"document_title":"doc title","start_char_index":10,"end_char_index":15}`,
		},
		{
			name: "PageLocation_with_file",
			v: anthropic.AnthropicCitationPageLocationParam{
				Type:            anthropic.AnthropicCitationPageLocationTypePageLocation,
				CitedText:       "section 1",
				FileID:          &fid,
				StartPageNumber: 2,
				EndPageNumber:   3,
			},
			want: `{"type":"page_location","cited_text":"section 1","file_id":"file_xxx","start_page_number":2,"end_page_number":3}`,
		},
		{
			name: "ContentBlockLocation_with_file",
			v: anthropic.AnthropicCitationContentBlockLocationParam{
				Type:            anthropic.AnthropicCitationContentBlockLocationTypeContentBlockLocation,
				CitedText:       "block content",
				FileID:          &fid,
				StartBlockIndex: 0,
				EndBlockIndex:   10,
			},
			want: `{"type":"content_block_location","cited_text":"block content","file_id":"file_xxx","start_block_index":0,"end_block_index":10}`,
		},
		{
			name: "SearchResultLocation",
			v: anthropic.AnthropicCitationSearchResultLocationParam{
				Type:              anthropic.AnthropicCitationSearchResultLocationTypeSearchResultLocation,
				CitedText:         "matched",
				StartBlockIndex:   0,
				EndBlockIndex:     5,
				SearchResultIndex: idx5,
				Title:             &title,
				Source:            &src,
			},
			want: `{"type":"search_result_location","cited_text":"matched","start_block_index":0,"end_block_index":5,"search_result_index":5,"title":"doc title","source":"src"}`,
		},
		{
			name: "WebSearchResultLocation",
			v: anthropic.AnthropicCitationWebSearchResultLocationParam{
				Type:           anthropic.AnthropicCitationWebSearchResultLocationTypeWebSearchResultLocation,
				CitedText:      "matched text",
				EncryptedIndex: "enc_blob",
				Title:          &title,
				URL:            &u,
			},
			want: `{"type":"web_search_result_location","cited_text":"matched text","encrypted_index":"enc_blob","title":"doc title","url":"https://example.com"}`,
		},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.v)
		if err != nil {
			t.Errorf("%s: marshal: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s: marshal mismatch:\n  got:  %s\n  want: %s", c.name, got, c.want)
		}
	}
}

// TestWireFormatConstructors 验证 8 个构造器函数生成的 wire-format 与字面量构造一致。
func TestWireFormatConstructors(t *testing.T) {
	t.Run("ThinkingConfig_Adaptive", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicThinkingConfigAdaptiveParam())
		want := `{"type":"adaptive"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("ThinkingConfig_Enabled", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicThinkingConfigEnabledParam(2048))
		want := `{"type":"enabled","budget_tokens":2048}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("ThinkingConfig_Disabled", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicThinkingConfigDisabledParam())
		want := `{"type":"disabled"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("Citation_CharLocation", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCitationCharLocationParam("hello", 0, 5))
		want := `{"type":"char_location","cited_text":"hello","start_char_index":0,"end_char_index":5}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("Citation_PageLocation", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCitationPageLocationParam("section 1", 2, 3))
		want := `{"type":"page_location","cited_text":"section 1","start_page_number":2,"end_page_number":3}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("Citation_ContentBlockLocation", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCitationContentBlockLocationParam("block", 0, 10))
		want := `{"type":"content_block_location","cited_text":"block","start_block_index":0,"end_block_index":10}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("Citation_SearchResultLocation", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCitationSearchResultLocationParam("matched", 0, 5, 5))
		want := `{"type":"search_result_location","cited_text":"matched","start_block_index":0,"end_block_index":5,"search_result_index":5}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("Citation_WebSearchResultLocation", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCitationWebSearchResultLocationParam("matched text", "enc_blob"))
		want := `{"type":"web_search_result_location","cited_text":"matched text","encrypted_index":"enc_blob"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})
}

// TestWireFormatCustomToolExtensions 验证 AnthropicCustomTool 4 个扩展字段的
// omitempty 行为 + NewAnthropicCustomTool 构造器。
func TestWireFormatCustomToolExtensions(t *testing.T) {
	t.Run("minimal_no_extensions", func(t *testing.T) {
		// 无扩展字段：与原 CustomTool_minimal wire-format 一致
		got, _ := json.Marshal(anthropic.NewAnthropicCustomTool("get_weather", anthropic.AnthropicToolInputSchema{Type: "object"}))
		want := `{"type":"custom","name":"get_weather","input_schema":{"type":"object"}}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("with_defer_loading", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCustomTool(
			"get_weather",
			anthropic.AnthropicToolInputSchema{Type: "object"},
			anthropic.WithDeferLoading(true),
		))
		want := `{"type":"custom","name":"get_weather","input_schema":{"type":"object"},"defer_loading":true}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("with_eager_input_streaming_and_allowed_callers", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCustomTool(
			"get_weather",
			anthropic.AnthropicToolInputSchema{Type: "object"},
			anthropic.WithEagerInputStreaming(true),
			anthropic.WithAllowedCallers("direct", "code_execution_20260120"),
			anthropic.WithDescription("Get weather"),
		))
		want := `{"type":"custom","name":"get_weather","description":"Get weather","input_schema":{"type":"object"},"eager_input_streaming":true,"allowed_callers":["direct","code_execution_20260120"]}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("with_input_examples", func(t *testing.T) {
		got, _ := json.Marshal(anthropic.NewAnthropicCustomTool(
			"get_weather",
			anthropic.AnthropicToolInputSchema{Type: "object"},
			anthropic.WithInputExamples(
				map[string]any{"input": map[string]any{"city": "SF"}, "output": "72F"},
			),
		))
		want := `{"type":"custom","name":"get_weather","input_schema":{"type":"object"},"input_examples":[{"input":{"city":"SF"},"output":"72F"}]}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("validate_allowed_callers", func(t *testing.T) {
		// 合法取值
		for _, c := range []string{"direct", "code_execution_20250825", "code_execution_20260120", "code_execution_20260521"} {
			if !anthropic.IsAllowedCallerValid(c) {
				t.Errorf("IsAllowedCallerValid(%q) = false, want true", c)
			}
		}
		// 不合法取值
		for _, c := range []string{"nonexistent", "code_execution_2099", ""} {
			if anthropic.IsAllowedCallerValid(c) {
				t.Errorf("IsAllowedCallerValid(%q) = true, want false", c)
			}
		}

		// 自定义工具的 ValidateAllowedCallers
		tool := anthropic.NewAnthropicCustomTool(
			"bash_run", anthropic.AnthropicToolInputSchema{Type: "object"},
			anthropic.WithAllowedCallers("direct", "code_execution_20260120"),
		)
		if bad, ok := tool.ValidateAllowedCallers(); !ok {
			t.Errorf("ValidateAllowedCallers: unexpected invalid %q", bad)
		}

		tool2 := anthropic.NewAnthropicCustomTool(
			"bad_run", anthropic.AnthropicToolInputSchema{Type: "object"},
			anthropic.WithAllowedCallers("direct", "bogus"),
		)
		bad, ok := tool2.ValidateAllowedCallers()
		if ok {
			t.Errorf("ValidateAllowedCallers: should have failed, got ok=true")
		}
		if bad != "bogus" {
			t.Errorf("ValidateAllowedCallers: bad = %q, want \"bogus\"", bad)
		}
	})
}

// TestWireFormatTextEditorHistorical 测试 text_editor 工具的 2 个历史版本
// （20250124 / 20250429）的 marshal wire-format。
func TestWireFormatTextEditorHistorical(t *testing.T) {
	t.Run("text_editor_20250124_minimal", func(t *testing.T) {
		tool := anthropic.AnthropicTextEditorTool20250124{
			Type: anthropic.AnthropicTextEditorToolType20250124,
			Name: "str_replace_editor",
		}
		got, _ := json.Marshal(tool)
		want := `{"type":"text_editor_20250124","name":"str_replace_editor"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_20250124_with_options", func(t *testing.T) {
		deferTrue := true
		tool := anthropic.AnthropicTextEditorTool20250124{
			Type:           anthropic.AnthropicTextEditorToolType20250124,
			Name:           "str_replace_editor",
			DeferLoading:   &deferTrue,
			AllowedCallers: []string{"direct"},
		}
		got, _ := json.Marshal(tool)
		want := `{"type":"text_editor_20250124","name":"str_replace_editor","defer_loading":true,"allowed_callers":["direct"]}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_20250429_minimal", func(t *testing.T) {
		tool := anthropic.AnthropicTextEditorTool20250429{
			Type: anthropic.AnthropicTextEditorToolType20250429,
			Name: "str_replace_based_edit_tool",
		}
		got, _ := json.Marshal(tool)
		want := `{"type":"text_editor_20250429","name":"str_replace_based_edit_tool"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_20250728_unchanged", func(t *testing.T) {
		// 既有的 20250728 版本仍能正常工作
		tool := anthropic.AnthropicTextEditorTool{
			Type: anthropic.AnthropicTextEditorToolType20250728,
			Name: "str_replace_based_edit_tool",
		}
		got, _ := json.Marshal(tool)
		want := `{"type":"text_editor_20250728","name":"str_replace_based_edit_tool"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})
}

// TestWireFormatUsageServiceTierResponses 测试 usage.service_tier 响应侧 3 档取值
// （standard / priority / batch）的 wire-format + 保留旧 auto / standard_only 兼容。
func TestWireFormatUsageServiceTierResponses(t *testing.T) {
	t.Run("marshal_three_response_values", func(t *testing.T) {
		for _, c := range []struct {
			s    anthropic.AnthropicUsageServiceTier
			want string
		}{
			{anthropic.AnthropicUsageServiceTierStandard, `"standard"`},
			{anthropic.AnthropicUsageServiceTierPriority, `"priority"`},
			{anthropic.AnthropicUsageServiceTierBatch, `"batch"`},
			{anthropic.AnthropicUsageServiceTierAuto, `"auto"`},
			{anthropic.AnthropicUsageServiceTierStandardOnly, `"standard_only"`},
		} {
			got, _ := json.Marshal(c.s)
			if string(got) != c.want {
				t.Errorf("marshal(%v) = %s, want %s", c.s, got, c.want)
			}
		}
	})
}

// TestWireFormatContentBlockSubresult 测试 code_execution / bash_code_execution /
// text_editor_code_execution 三类工具响应块内部子结构的 marshal wire-format。
func TestWireFormatContentBlockSubresult(t *testing.T) {
	t.Run("code_execution_result_success", func(t *testing.T) {
		block := anthropic.AnthropicCodeExecutionResultBlock{
			Type:       anthropic.CodeExecutionResultTypeCodeExecutionResult,
			ReturnCode: 0,
			Stderr:     "",
			Stdout:     "Hello from code\n",
		}
		got, _ := json.Marshal(block)
		want := `{"type":"code_execution_result","return_code":0,"stderr":"","stdout":"Hello from code\n"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("code_execution_result_with_outputs", func(t *testing.T) {
		block := anthropic.AnthropicCodeExecutionResultBlock{
			Type:       anthropic.CodeExecutionResultTypeCodeExecutionResult,
			Content:    []anthropic.AnthropicCodeExecutionOutputBlock{{Type: anthropic.CodeExecutionOutputTypeCodeExecutionOutput, FileID: "file_abc"}},
			ReturnCode: 1,
			Stderr:     "err",
			Stdout:     "",
		}
		got, _ := json.Marshal(block)
		want := `{"type":"code_execution_result","content":[{"type":"code_execution_output","file_id":"file_abc"}],"return_code":1,"stderr":"err","stdout":""}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("code_execution_result_error", func(t *testing.T) {
		block := anthropic.AnthropicCodeExecutionToolResultError{
			Type:      anthropic.CodeExecutionToolResultErrorTypeCodeExecutionToolResultError,
			ErrorCode: anthropic.AnthropicCodeExecutionToolResultErrorCodeTimeLimit,
		}
		got, _ := json.Marshal(block)
		want := `{"type":"code_execution_tool_result_error","error_code":"time_limit_exceeded"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("bash_code_execution_result_success", func(t *testing.T) {
		block := anthropic.AnthropicBashCodeExecutionResultBlock{
			Type:       anthropic.BashCodeExecutionResultTypeBashCodeExecutionResult,
			ReturnCode: 0,
			Stderr:     "",
			Stdout:     "out",
		}
		got, _ := json.Marshal(block)
		want := `{"type":"bash_code_execution_result","return_code":0,"stderr":"","stdout":"out"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_view_result", func(t *testing.T) {
		start := int64(10)
		block := anthropic.AnthropicTextEditorCodeExecutionViewResultBlock{
			Type:      anthropic.TextEditorCodeExecutionViewResultTypeTextEditorCodeExecutionViewResult,
			Content:   "line content...",
			StartLine: &start,
		}
		got, _ := json.Marshal(block)
		want := `{"type":"text_editor_code_execution_view_result","content":"line content...","start_line":10}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_create_result", func(t *testing.T) {
		block := anthropic.AnthropicTextEditorCodeExecutionCreateResultBlock{
			Type:         anthropic.TextEditorCodeExecutionCreateResultTypeTextEditorCodeExecutionCreateResult,
			IsFileUpdate: false,
		}
		got, _ := json.Marshal(block)
		want := `{"type":"text_editor_code_execution_create_result","is_file_update":false}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_str_replace_result", func(t *testing.T) {
		newStart := int64(5)
		newLines := int64(3)
		block := anthropic.AnthropicTextEditorCodeExecutionStrReplaceResultBlock{
			Type:     anthropic.TextEditorCodeExecutionStrReplaceResultTypeTextEditorCodeExecutionStrReplaceResult,
			NewStart: &newStart,
			NewLines: &newLines,
		}
		got, _ := json.Marshal(block)
		want := `{"type":"text_editor_code_execution_str_replace_result","new_lines":3,"new_start":5}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("text_editor_error_with_message", func(t *testing.T) {
		msg := "File not found"
		block := anthropic.AnthropicTextEditorCodeExecutionToolResultError{
			Type:         anthropic.TextEditorCodeExecutionToolResultErrorTypeTextEditorCodeExecutionToolResultError,
			ErrorCode:    anthropic.AnthropicTextEditorCodeExecutionToolResultErrorCodeFileNotFound,
			ErrorMessage: &msg,
		}
		got, _ := json.Marshal(block)
		want := `{"type":"text_editor_code_execution_tool_result_error","error_code":"file_not_found","error_message":"File not found"}`
		if string(got) != want {
			t.Errorf("marshal mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})
}
