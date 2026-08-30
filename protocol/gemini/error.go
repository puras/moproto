package gemini

// GeminiErrorResponse HTTP 4xx/5xx 响应体（Google API 标准错误结构）。
type GeminiErrorResponse struct {
	Error *GeminiAPIError `json:"error"`
}

// GeminiAPIError 单个错误对象。
//
// 字段命名遵循 Google API 标准：code（HTTP status）、message、status（Google API
// 状态字符串，如 UNAUTHENTICATED / RESOURCE_EXHAUSTED）。
type GeminiAPIError struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Status  string           `json:"status"`
	Details []map[string]any `json:"details,omitempty"`
}
