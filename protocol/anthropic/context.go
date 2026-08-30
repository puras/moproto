package anthropic

import "encoding/json"

// AnthropicClearToolUses20250919EditType clear_tool_uses_20250919 edit 的 type 字段类型。
type AnthropicClearToolUses20250919EditType string

const AnthropicClearToolUses20250919EditTypeClearToolUses20250919 AnthropicClearToolUses20250919EditType = "clear_tool_uses_20250919"

// AnthropicClearThinking20251015EditType clear_thinking_20251015 edit 的 type 字段类型。
type AnthropicClearThinking20251015EditType string

const AnthropicClearThinking20251015EditTypeClearThinking20251015 AnthropicClearThinking20251015EditType = "clear_thinking_20251015"

// AnthropicCompact20260112EditType compact_20260112 edit 的 type 字段类型。
type AnthropicCompact20260112EditType string

const AnthropicCompact20260112EditTypeCompact20260112 AnthropicCompact20260112EditType = "compact_20260112"

// AnthropicContextEdit 上下文编辑动作 union。
//
// 服务端按 edits 数组顺序对请求消息做处理：
//   - clear_tool_uses : 清除旧的 tool_use/tool_result 配对
//   - clear_thinking  : 清除 thinking 块
//   - compact         : 压缩为摘要
//
// 通过 isContextEdit() 私有方法封闭实现集（外部包无法新增实现）。
// 兼容旧 sealed interface；AnthropicContextEditUnion 为扁平实现（解决 unmarshal 问题）。
type AnthropicContextEdit interface {
	isContextEdit()
}

// AnthropicContextEditUnion 扁平实现（解决 interface unmarshal 问题）。
type AnthropicContextEditUnion struct {
	Type          string          `json:"type"`
	ClearToolUses json.RawMessage `json:"clear_tool_uses_20250919,omitempty"`
	ClearThinking json.RawMessage `json:"clear_thinking_20251015,omitempty"`
	Compact       json.RawMessage `json:"compact_20260112,omitempty"`
}

func (AnthropicContextEditUnion) isContextEdit() {}

func (u AnthropicContextEditUnion) MarshalJSON() ([]byte, error) {
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
	case string(AnthropicCompact20260112EditTypeCompact20260112):
		if u.Compact == nil {
			return nil, nil
		}
		return u.Compact, nil
	default:
		return nil, nil
	}
}

func (u *AnthropicContextEditUnion) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type

	switch tmp.Type {
	case string(AnthropicClearToolUses20250919EditTypeClearToolUses20250919):
		var b AnthropicClearToolUses20250919Edit
		if err := json.Unmarshal(data, &b); err == nil {
			u.ClearToolUses, _ = json.Marshal(b)
		}
	case string(AnthropicClearThinking20251015EditTypeClearThinking20251015):
		var b AnthropicClearThinking20251015Edit
		if err := json.Unmarshal(data, &b); err == nil {
			u.ClearThinking, _ = json.Marshal(b)
		}
	case string(AnthropicCompact20260112EditTypeCompact20260112):
		var b AnthropicCompact20260112Edit
		if err := json.Unmarshal(data, &b); err == nil {
			u.Compact, _ = json.Marshal(b)
		}
	}
	return nil
}

// AnthropicContextManagement 长上下文管理（beta）。
//
// 服务端按 edits 数组顺序对请求消息做处理。
type AnthropicContextManagement struct {
	Edits []AnthropicContextEditUnion `json:"edits"`
}

// AnthropicClearToolUses20250919Edit 清除 tool_use/tool_result 配对。
//
// 可选参数控制清除阈值与保留范围。
type AnthropicClearToolUses20250919Edit struct {
	Type             AnthropicClearToolUses20250919EditType `json:"type"`                        // 固定 "clear_tool_uses_20250919"
	Keep             string                                 `json:"keep,omitempty"`              // 保留规则
	ClearToolInputs  bool                                   `json:"clear_tool_inputs,omitempty"` // 是否清空 tool input 内容
	MinTokensToClear int                                    `json:"min_tokens_to_clear,omitempty"`
}

func (AnthropicClearToolUses20250919Edit) isContextEdit() {}

// AnthropicClearThinking20251015Edit 清除 thinking 块。
type AnthropicClearThinking20251015Edit struct {
	Type AnthropicClearThinking20251015EditType `json:"type"`           // 固定 "clear_thinking_20251015"
	Keep string                                 `json:"keep,omitempty"` // 保留规则
}

func (AnthropicClearThinking20251015Edit) isContextEdit() {}

// AnthropicCompact20260112Edit 压缩长上下文为摘要。
type AnthropicCompact20260112Edit struct {
	Type AnthropicCompact20260112EditType `json:"type"`           // 固定 "compact_20260112"
	Keep string                           `json:"keep,omitempty"` // 保留规则
}

func (AnthropicCompact20260112Edit) isContextEdit() {}
