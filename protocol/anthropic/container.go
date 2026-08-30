package anthropic

// AnthropicContainer 容器/技能挂载（请求方向）。
//
// 顶层 container 字段接受两种形态：
//   - string  ：已存在的容器 ID（用于续接上次的 code_execution 会话）
//   - object  ：新容器声明，含 skills 列表
//
// 实际序列化时由调用方根据具体形态自行编码（json.Marshal 直接支持）。
type AnthropicContainer struct {
	ID        string           `json:"id,omitempty"`         // 响应侧：容器 ID（如 "container_xxx"）
	ExpiresAt string           `json:"expires_at,omitempty"` // 响应侧：容器过期时间
	Skills    []AnthropicSkill `json:"skills,omitempty"`     // 请求侧：技能列表
}

// SkillType 技能类型枚举。
type SkillType string

const (
	SkillTypeAnthropic SkillType = "anthropic"
	SkillTypeCustom    SkillType = "custom"
)

// AnthropicSkill 技能声明。
type AnthropicSkill struct {
	Type    SkillType `json:"type"`              // SkillTypeAnthropic | SkillTypeCustom
	SkillID string    `json:"skill_id"`          // 技能标识
	Version string    `json:"version,omitempty"` // 版本号
}
