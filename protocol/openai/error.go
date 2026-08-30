package openai

// OpenAIErrorResponse 错误响应顶层。
type OpenAIErrorResponse struct {
	// Error 错误对象。
	Error OpenAIError `json:"error"`
}

// OpenAIError 错误对象。
type OpenAIError struct {
	// Message 错误消息。
	Message string `json:"message"`
	// Type 错误类型（如 "invalid_request_error"）。
	Type string `json:"type"`
	// Code 错误代码（可选，如 "invalid_api_key"）。
	Code *string `json:"code,omitempty"`
	// Param 错误参数名（可选）。
	Param *string `json:"param,omitempty"`
	// StatusCode HTTP 状态码（可选）。
	StatusCode *int64 `json:"status_code,omitempty"`
}

// OpenAIErrorType 错误类型（union 内部使用）。
type OpenAIErrorType string

const (
	OpenAIErrorTypeInvalidRequest   OpenAIErrorType = "invalid_request_error"
	OpenAIErrorTypeAuthentication   OpenAIErrorType = "authentication_error"
	OpenAIErrorTypePermissionDenied OpenAIErrorType = "permission_denied_error"
	OpenAIErrorTypeNotFound         OpenAIErrorType = "not_found_error"
	OpenAIErrorTypeRateLimit        OpenAIErrorType = "rate_limit_error"
	OpenAIErrorTypeServer           OpenAIErrorType = "server_error"
	OpenAIErrorTypeBilling          OpenAIErrorType = "billing_error"
	OpenAIErrorTypeQuotaExceeded    OpenAIErrorType = "insufficient_quota_error"
	OpenAIErrorTypeTimeout          OpenAIErrorType = "timeout_error"
)
