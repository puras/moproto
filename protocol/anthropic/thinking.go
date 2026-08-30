package anthropic

// ThinkingDisplay 控制 adaptive 模式下 thinking 内容对客户端的可见性（4.7/4.8 新增）。
//
// 仅适用于 Adaptive 模式；Enabled 模式有独立的 ThinkingConfigEnabledDisplay（见下）。
type ThinkingDisplay string

const (
	ThinkingDisplaySummarized ThinkingDisplay = "summarized"
	ThinkingDisplayOmitted    ThinkingDisplay = "omitted"
)

// ThinkingConfigEnabledDisplay 控制 enabled 模式下 thinking 内容对客户端的可见性（4.7/4.8 新增）。
//
// 仅适用于 Enabled 模式；与 Adaptive 模式的 ThinkingDisplay 取值集合相同但类型独立，
// 编译期防止两种模式互相串用。
type ThinkingConfigEnabledDisplay string

const (
	ThinkingConfigEnabledDisplaySummarized ThinkingConfigEnabledDisplay = "summarized"
	ThinkingConfigEnabledDisplayOmitted    ThinkingConfigEnabledDisplay = "omitted"
)

// ============== AnthropicThinkingConfig union ==============
//
// AnthropicThinkingConfig 是请求级 thinking 配置的 union 接口，对齐官方
// anthropic-sdk-go 的 ThinkingConfigParamUnion（message.go:6610）。
//
// 通过 isThinkingConfig() 私有方法封闭实现集；调用方通过 type switch 或各
// 实现的 Type 字段分发，详见 doc.go。
type AnthropicThinkingConfig interface {
	isThinkingConfig()
}

// ============== 类型化 Type 字段 ==============

// AnthropicThinkingConfigAdaptiveType adaptive 模式的 type 字段类型。
type AnthropicThinkingConfigAdaptiveType string

// AnthropicThinkingConfigAdaptiveTypeAdaptive adaptive 模式唯一合法取值。
const AnthropicThinkingConfigAdaptiveTypeAdaptive AnthropicThinkingConfigAdaptiveType = "adaptive"

// AnthropicThinkingConfigEnabledType enabled 模式的 type 字段类型。
type AnthropicThinkingConfigEnabledType string

// AnthropicThinkingConfigEnabledTypeEnabled enabled 模式唯一合法取值。
const AnthropicThinkingConfigEnabledTypeEnabled AnthropicThinkingConfigEnabledType = "enabled"

// AnthropicThinkingConfigDisabledType disabled 模式的 type 字段类型。
type AnthropicThinkingConfigDisabledType string

// AnthropicThinkingConfigDisabledTypeDisabled disabled 模式唯一合法取值。
const AnthropicThinkingConfigDisabledTypeDisabled AnthropicThinkingConfigDisabledType = "disabled"

// ============== 3 个实现 ==============

// AnthropicThinkingConfigAdaptive 4.6+ 推荐模式：Claude 自适应思考深度，无需指定 budget。
//
// Display 控制 thinking 内容对客户端的可见性：
//   - ThinkingDisplaySummarized（默认）：返回精简摘要
//   - ThinkingDisplayOmitted：不返回 thinking 文本（仅保留 block 占位）
type AnthropicThinkingConfigAdaptive struct {
	Type    AnthropicThinkingConfigAdaptiveType `json:"type"` // 固定 "adaptive"
	Display ThinkingDisplay                     `json:"display,omitempty"`
}

func (AnthropicThinkingConfigAdaptive) isThinkingConfig() {}

// AnthropicThinkingConfigEnabled 传统 extended-thinking 模式（4.5- 之前）。
//
// 必填字段：Type + BudgetTokens。
// BudgetTokens 约束（≥1024 且 < max_tokens）由 Go 层不校验，留给网关 validation 链路。
type AnthropicThinkingConfigEnabled struct {
	Type         AnthropicThinkingConfigEnabledType `json:"type"` // 固定 "enabled"
	BudgetTokens int64                              `json:"budget_tokens"`
	Display      ThinkingConfigEnabledDisplay       `json:"display,omitempty"`
}

func (AnthropicThinkingConfigEnabled) isThinkingConfig() {}

// AnthropicThinkingConfigDisabled 显式禁用 thinking。
//
// 仅 type 字段，无其他配置。
type AnthropicThinkingConfigDisabled struct {
	Type AnthropicThinkingConfigDisabledType `json:"type"` // 固定 "disabled"
}

func (AnthropicThinkingConfigDisabled) isThinkingConfig() {}

// ============== 构造器 ==============
//
// 构造器函数避免调用方手填 Type 字段导致的 wire-format 字符串错误。
// 对齐官方 anthropic-sdk-go 的 NewThinkingConfig*Param / ThinkingConfigParamOf*。

// NewAnthropicThinkingConfigAdaptiveParam 构造 adaptive 模式（4.6+ 推荐）。
func NewAnthropicThinkingConfigAdaptiveParam() AnthropicThinkingConfigAdaptive {
	return AnthropicThinkingConfigAdaptive{
		Type: AnthropicThinkingConfigAdaptiveTypeAdaptive,
	}
}

// NewAnthropicThinkingConfigEnabledParam 构造 enabled 模式（传统 extended-thinking）。
//
// BudgetTokens 必须 ≥1024 且 < max_tokens；Go 层不校验，由网关 validation 链路负责。
func NewAnthropicThinkingConfigEnabledParam(budgetTokens int64) AnthropicThinkingConfigEnabled {
	return AnthropicThinkingConfigEnabled{
		Type:         AnthropicThinkingConfigEnabledTypeEnabled,
		BudgetTokens: budgetTokens,
	}
}

// NewAnthropicThinkingConfigDisabledParam 构造 disabled 模式。
func NewAnthropicThinkingConfigDisabledParam() AnthropicThinkingConfigDisabled {
	return AnthropicThinkingConfigDisabled{
		Type: AnthropicThinkingConfigDisabledTypeDisabled,
	}
}
