package anthropic

// AnthropicStopReason stop_reason 字段类型。
//
// 6 个取值覆盖正常结束 / 超限 / 工具调用 / 拒绝 / 暂停。
type AnthropicStopReason string

const (
	AnthropicStopReasonEndTurn      AnthropicStopReason = "end_turn"
	AnthropicStopReasonMaxTokens    AnthropicStopReason = "max_tokens"
	AnthropicStopReasonStopSequence AnthropicStopReason = "stop_sequence"
	AnthropicStopReasonToolUse      AnthropicStopReason = "tool_use"
	AnthropicStopReasonRefusal      AnthropicStopReason = "refusal"
	AnthropicStopReasonPauseTurn    AnthropicStopReason = "pause_turn"
)

// AnthropicRefusalCategory refusal 分类（4.6+ 引入）。
type AnthropicRefusalCategory string

const (
	AnthropicRefusalCategoryCyber               AnthropicRefusalCategory = "cyber"
	AnthropicRefusalCategoryBio                 AnthropicRefusalCategory = "bio"
	AnthropicRefusalCategoryFrontierLLM         AnthropicRefusalCategory = "frontier_llm"
	AnthropicRefusalCategoryReasoningExtraction AnthropicRefusalCategory = "reasoning_extraction"
	AnthropicRefusalCategoryMilitaryWeapons     AnthropicRefusalCategory = "military_weapons"
)

// AnthropicRefusalDetails 拒绝原因详情（stop_reason=refusal 时携带）。
type AnthropicRefusalDetails struct {
	Category AnthropicRefusalCategory `json:"category,omitempty"`
	Reason   string                   `json:"reason,omitempty"`
}

// AnthropicStopDetails stop_details 字段。
//
// 当 stop_reason=refusal 时，Refusal 字段填充具体拒绝原因和分类。
type AnthropicStopDetails struct {
	Refusal *AnthropicRefusalDetails `json:"refusal,omitempty"`
}

// AnthropicMessageResponse 顶层 Message 响应。
//
// 字段顺序对齐官方 SDK Message 结构；Container 仅在 code_execution 路径有值。
// Content 是 AnthropicContentBlock union 切片（见 content_block.go）。
// StopReason / StopSequence / StopDetails 通常在最后一个 message_delta 中出现。
//
// 实现 common.Response 接口的 4 个方法（GetFormat/GetID/GetModel/GetUsage），
// 但用 string / AnthropicUsage 而非 common.Format / common.Usage（pkg 不能引
// internal）；common.Format 已经是 string alias，类型签名匹配。
type AnthropicMessageResponse struct {
	ID           string                        `json:"id"`   // msg_xxx
	Type         AnthropicMessageType          `json:"type"` // 固定 "message"
	Role         AnthropicMessageRole          `json:"role"` // 固定 "assistant"
	Model        string                        `json:"model"`
	Container    *AnthropicContainer           `json:"container,omitempty"` // code_execution 容器
	Content      []*AnthropicContentBlockUnion `json:"content"`
	StopReason   *AnthropicStopReason          `json:"stop_reason,omitempty"`
	StopSequence *string                       `json:"stop_sequence,omitempty"`
	StopDetails  *AnthropicStopDetails         `json:"stop_details,omitempty"` // refusal 时携带
	Usage        *AnthropicUsage               `json:"usage,omitempty"`
}
