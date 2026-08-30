package anthropic

import "encoding/json"

// AnthropicEffort 思考/行动深度枚举（4.6+）。
type AnthropicEffort string

const (
	AnthropicEffortLow    AnthropicEffort = "low"
	AnthropicEffortMedium AnthropicEffort = "medium"
	AnthropicEffortHigh   AnthropicEffort = "high"
	AnthropicEffortXHigh  AnthropicEffort = "xhigh"
	AnthropicEffortMax    AnthropicEffort = "max"
)

// AnthropicOutputConfig 输出控制（4.6+ GA）。
//
// 三部分组合使用：
//   - Effort     ：思考/行动深度（4.6+）
//   - Format     ：结构化输出（json_schema）
//   - TaskBudget ：Agentic 任务 token 预算（beta）
type AnthropicOutputConfig struct {
	Effort     AnthropicEffort            `json:"effort,omitempty"`      // low|medium|high|xhigh|max
	Format     *AnthropicJsonOutputFormat `json:"format,omitempty"`      // 结构化输出
	TaskBudget *AnthropicTaskBudget       `json:"task_budget,omitempty"` // Agentic 任务预算（beta）
}

// AnthropicJsonOutputFormatType output_config.format.type 字段类型。
type AnthropicJsonOutputFormatType string

const AnthropicJsonOutputFormatTypeJsonSchema AnthropicJsonOutputFormatType = "json_schema"

// AnthropicJsonOutputFormat 结构化输出格式（json_schema）。
//
// Schema 是 JSON Schema 对象（任意结构），用 RawMessage 保留原始字节。
type AnthropicJsonOutputFormat struct {
	Type   AnthropicJsonOutputFormatType `json:"type"` // 固定 "json_schema"
	Schema json.RawMessage               `json:"schema"`
}

// AnthropicTaskBudgetType output_config.task_budget.type 字段类型。
type AnthropicTaskBudgetType string

const AnthropicTaskBudgetTypeTokens AnthropicTaskBudgetType = "tokens"

// AnthropicTaskBudget Agentic 任务的 token 预算（beta）。
type AnthropicTaskBudget struct {
	Type  AnthropicTaskBudgetType `json:"type"`  // 固定 "tokens"
	Total int                     `json:"total"` // 预算 token 数
}
