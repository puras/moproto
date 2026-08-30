package anthropic

import "encoding/json"

// AnthropicAppliedEditUnion 上下文压缩结果 union。
//
// 服务端在 context_management 处理完成后，通过 applied_edits 字段报告实际执行的压缩操作。
// 通过 isAppliedEdit() 私有方法封闭实现集（外部包无法新增实现）。
// 兼容旧 sealed interface；AnthropicAppliedEditUnionUnion 为扁平实现（解决 unmarshal 问题）。
type AnthropicAppliedEditUnion interface {
	isAppliedEdit()
}

// AnthropicAppliedEditUnionUnion 扁平实现（解决 interface unmarshal 问题）。
type AnthropicAppliedEditUnionUnion struct {
	Type          string          `json:"type"`
	ClearToolUses json.RawMessage `json:"clear_tool_uses_20250919,omitempty"`
	ClearThinking json.RawMessage `json:"clear_thinking_20251015,omitempty"`
}

func (AnthropicAppliedEditUnionUnion) isAppliedEdit() {}

func (u AnthropicAppliedEditUnionUnion) MarshalJSON() ([]byte, error) {
	// 空 Type 或 nil RawMessage 返回 nil 以触发 omitempty suppression
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case string(AnthropicClearToolUses20250919EditTypeClearToolUses20250919):
		if u.ClearToolUses == nil {
			return nil, nil
		}
		return u.ClearToolUses, nil
	case string(AnthropicClearThinking20251015EditTypeClearThinking20251015):
		if u.ClearThinking == nil {
			return nil, nil
		}
		return u.ClearThinking, nil
	default:
		return nil, nil
	}
}

func (u *AnthropicAppliedEditUnionUnion) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type

	switch tmp.Type {
	case string(AnthropicClearToolUses20250919EditTypeClearToolUses20250919):
		var b AnthropicClearToolUses20250919EditResponse
		if err := json.Unmarshal(data, &b); err == nil {
			u.ClearToolUses, _ = json.Marshal(b)
		}
	case string(AnthropicClearThinking20251015EditTypeClearThinking20251015):
		var b AnthropicClearThinking20251015EditResponse
		if err := json.Unmarshal(data, &b); err == nil {
			u.ClearThinking, _ = json.Marshal(b)
		}
	}
	return nil
}

// AnthropicClearToolUses20250919EditResponse 清除 tool_use/tool_result 配对的结果。
//
// 报告实际清除的 tool_use 数量和释放的 token 数。
type AnthropicClearToolUses20250919EditResponse struct {
	Type               AnthropicClearToolUses20250919EditType `json:"type"`                           // 固定 "clear_tool_uses_20250919"
	ClearedInputTokens int64                                  `json:"cleared_input_tokens,omitempty"` // 释放的输入 token
	ClearedToolUses    int64                                  `json:"cleared_tool_uses,omitempty"`    // 清除的 tool_use 数量
}

func (AnthropicClearToolUses20250919EditResponse) isAppliedEdit() {}

// AnthropicClearThinking20251015EditResponse 清除 thinking 块的结果。
//
// 报告实际清除的 thinking 轮次和释放的 token 数。
type AnthropicClearThinking20251015EditResponse struct {
	Type                 AnthropicClearThinking20251015EditType `json:"type"`                             // 固定 "clear_thinking_20251015"
	ClearedInputTokens   int64                                  `json:"cleared_input_tokens,omitempty"`   // 释放的输入 token
	ClearedThinkingTurns int64                                  `json:"cleared_thinking_turns,omitempty"` // 清除的 thinking 轮次
}

func (AnthropicClearThinking20251015EditResponse) isAppliedEdit() {}

// AnthropicContextManagementResponse 上下文管理的服务端响应。
//
// 报告 context_management.edits 中各操作实际执行的结果。
type AnthropicContextManagementResponse struct {
	AppliedEdits []AnthropicAppliedEditUnionUnion `json:"applied_edits"`
}
