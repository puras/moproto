package openai

import "encoding/json"

// OpenAIResponsesCompactionRequest OpenAI Responses API 压缩端点请求体（/v1/responses/compact）
//
// 用途：客户端调用 /v1/responses/compact 把长上下文（history + tool 状态 + file refs）压缩为
// 单条 summary message，缩短后续 input 长度、节省 token 成本。
//
// 字段集（轻量；与 new-api 一致；不强行对齐 OpenAI 官方——第三方厂商多未实现此端点，
// 本期按 new-api 形态覆盖 OpenAI 官方 + 未来第三方实现）：
//   - Model:              实际模型（gpt-4o / gpt-5 等）
//   - Input:              待压缩内容（string | []InputItem）；保留 json.RawMessage 由调用方按 type switch 处理
//   - Instructions:       可选系统提示
//   - PreviousResponseID: 上一次 response ID；可与 Input 互斥
//
// 实现 common.Request 接口（GetModel / GetFormat）。
func (r *OpenAIResponsesCompactionRequest) GetModel() string  { return r.Model }
func (r *OpenAIResponsesCompactionRequest) GetFormat() string { return "responses_compact" }

type OpenAIResponsesCompactionRequest struct {
	// Model 实际模型名（必填；gpt-4o / gpt-5 / o3 等）。
	Model string `json:"model"`
	// Input 待压缩内容（string | []InputItem）。保留 json.RawMessage 避免
	// 与 OpenAIResponseInputUnion（30 变体）耦合；handler 走 RawBody 透传
	// 路径时不会进入 typed unmarshal。
	Input json.RawMessage `json:"input,omitempty"`
	// Instructions 可选系统提示；压缩时可作为压缩目标的语义指导。
	Instructions string `json:"instructions,omitempty"`
	// PreviousResponseID 上一次 response ID；可与 Input 互斥
	// （OpenAI 文档：两者同时给时 Input 优先；网关侧不强制二选一）。
	PreviousResponseID string `json:"previous_response_id,omitempty"`
}
