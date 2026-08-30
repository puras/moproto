package anthropic

import "encoding/json"

// AnthropicErrorType error_response.type 字段类型。
//
// 顶层错误响应 type 固定为 "error"。
type AnthropicErrorType string

const AnthropicErrorTypeError AnthropicErrorType = "error"

// AnthropicErrorObject 错误对象 union（error_response.error 字段）。
//
// 9 种变体对应 9 种 ErrorType，每种仅 Type + Message 两个字段。
// 通过 isErrorObject() 私有方法封闭实现集。
type AnthropicErrorObject interface {
	isErrorObject()
}

// ============== 9 个错误变体 ==============

// AnthropicInvalidRequestError 请求参数错误（HTTP 400）。
type AnthropicInvalidRequestError struct {
	Type    AnthropicInvalidRequestErrorType `json:"type"` // 固定 "invalid_request_error"
	Message string                           `json:"message"`
}

func (AnthropicInvalidRequestError) isErrorObject() {}

type AnthropicInvalidRequestErrorType string

const AnthropicInvalidRequestErrorTypeInvalidRequestError AnthropicInvalidRequestErrorType = "invalid_request_error"

// AnthropicAuthenticationError 认证错误（HTTP 401）。
type AnthropicAuthenticationError struct {
	Type    AnthropicAuthenticationErrorType `json:"type"` // 固定 "authentication_error"
	Message string                           `json:"message"`
}

func (AnthropicAuthenticationError) isErrorObject() {}

type AnthropicAuthenticationErrorType string

const AnthropicAuthenticationErrorTypeAuthenticationError AnthropicAuthenticationErrorType = "authentication_error"

// AnthropicPermissionError 权限错误（HTTP 403）。
type AnthropicPermissionError struct {
	Type    AnthropicPermissionErrorType `json:"type"` // 固定 "permission_error"
	Message string                       `json:"message"`
}

func (AnthropicPermissionError) isErrorObject() {}

type AnthropicPermissionErrorType string

const AnthropicPermissionErrorTypePermissionError AnthropicPermissionErrorType = "permission_error"

// AnthropicNotFoundError 资源未找到（HTTP 404）。
type AnthropicNotFoundError struct {
	Type    AnthropicNotFoundErrorType `json:"type"` // 固定 "not_found_error"
	Message string                     `json:"message"`
}

func (AnthropicNotFoundError) isErrorObject() {}

type AnthropicNotFoundErrorType string

const AnthropicNotFoundErrorTypeNotFoundError AnthropicNotFoundErrorType = "not_found_error"

// AnthropicRateLimitError 速率限制（HTTP 429）。
type AnthropicRateLimitError struct {
	Type    AnthropicRateLimitErrorType `json:"type"` // 固定 "rate_limit_error"
	Message string                      `json:"message"`
}

func (AnthropicRateLimitError) isErrorObject() {}

type AnthropicRateLimitErrorType string

const AnthropicRateLimitErrorTypeRateLimitError AnthropicRateLimitErrorType = "rate_limit_error"

// AnthropicTimeoutError 超时（HTTP 408 / 504）。
type AnthropicTimeoutError struct {
	Type    AnthropicTimeoutErrorType `json:"type"` // 固定 "timeout_error"
	Message string                    `json:"message"`
}

func (AnthropicTimeoutError) isErrorObject() {}

type AnthropicTimeoutErrorType string

const AnthropicTimeoutErrorTypeTimeoutError AnthropicTimeoutErrorType = "timeout_error"

// AnthropicOverloadedError 服务端过载（HTTP 529）。
type AnthropicOverloadedError struct {
	Type    AnthropicOverloadedErrorType `json:"type"` // 固定 "overloaded_error"
	Message string                       `json:"message"`
}

func (AnthropicOverloadedError) isErrorObject() {}

type AnthropicOverloadedErrorType string

const AnthropicOverloadedErrorTypeOverloadedError AnthropicOverloadedErrorType = "overloaded_error"

// AnthropicAPIError 通用 API 错误（HTTP 500）。
type AnthropicAPIError struct {
	Type    AnthropicAPIErrorType `json:"type"` // 固定 "api_error"
	Message string                `json:"message"`
}

func (AnthropicAPIError) isErrorObject() {}

type AnthropicAPIErrorType string

const AnthropicAPIErrorTypeAPIError AnthropicAPIErrorType = "api_error"

// AnthropicBillingError 账单错误（HTTP 402）。
type AnthropicBillingError struct {
	Type    AnthropicBillingErrorType `json:"type"` // 固定 "billing_error"
	Message string                    `json:"message"`
}

func (AnthropicBillingError) isErrorObject() {}

type AnthropicBillingErrorType string

const AnthropicBillingErrorTypeBillingError AnthropicBillingErrorType = "billing_error"

// AnthropicErrorResponse HTTP 非 2xx 响应体（Anthropic API 标准错误结构）。
type AnthropicErrorResponse struct {
	Type      AnthropicErrorType         `json:"type"`                 // 固定 "error"
	Error     *AnthropicErrorObjectUnion `json:"error"`                // 9 种变体 union（扁平实现，解决 interface unmarshal 问题）
	RequestID string                     `json:"request_id,omitempty"` // 如 "req_xxx"
}

// AnthropicErrorObjectUnion 错误对象 union 的扁平实现。
//
// 解决 json.Unmarshal 无法反序列化到 interface 字段的问题。
// 9 个 RawMessage 字段分别存各变体完整 JSON（含 type）。
type AnthropicErrorObjectUnion struct {
	Type string `json:"type"`

	InvalidRequestError json.RawMessage `json:"invalid_request_error,omitempty"`
	AuthenticationError json.RawMessage `json:"authentication_error,omitempty"`
	PermissionError     json.RawMessage `json:"permission_error,omitempty"`
	NotFoundError       json.RawMessage `json:"not_found_error,omitempty"`
	RateLimitError      json.RawMessage `json:"rate_limit_error,omitempty"`
	TimeoutError        json.RawMessage `json:"timeout_error,omitempty"`
	OverloadedError     json.RawMessage `json:"overloaded_error,omitempty"`
	APIError            json.RawMessage `json:"api_error,omitempty"`
	BillingError        json.RawMessage `json:"billing_error,omitempty"`
}

func (AnthropicErrorObjectUnion) isErrorObject() {}

func (u AnthropicErrorObjectUnion) MarshalJSON() ([]byte, error) {
	// 空 Type 或 nil RawMessage 返回 nil 以触发 omitempty suppression
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case string(AnthropicInvalidRequestErrorTypeInvalidRequestError):
		if u.InvalidRequestError == nil {
			return nil, nil
		}
		return u.InvalidRequestError, nil
	case string(AnthropicAuthenticationErrorTypeAuthenticationError):
		if u.AuthenticationError == nil {
			return nil, nil
		}
		return u.AuthenticationError, nil
	case string(AnthropicPermissionErrorTypePermissionError):
		if u.PermissionError == nil {
			return nil, nil
		}
		return u.PermissionError, nil
	case string(AnthropicNotFoundErrorTypeNotFoundError):
		if u.NotFoundError == nil {
			return nil, nil
		}
		return u.NotFoundError, nil
	case string(AnthropicRateLimitErrorTypeRateLimitError):
		if u.RateLimitError == nil {
			return nil, nil
		}
		return u.RateLimitError, nil
	case string(AnthropicTimeoutErrorTypeTimeoutError):
		if u.TimeoutError == nil {
			return nil, nil
		}
		return u.TimeoutError, nil
	case string(AnthropicOverloadedErrorTypeOverloadedError):
		if u.OverloadedError == nil {
			return nil, nil
		}
		return u.OverloadedError, nil
	case string(AnthropicAPIErrorTypeAPIError):
		if u.APIError == nil {
			return nil, nil
		}
		return u.APIError, nil
	case string(AnthropicBillingErrorTypeBillingError):
		if u.BillingError == nil {
			return nil, nil
		}
		return u.BillingError, nil
	default:
		return nil, nil
	}
}

func (u *AnthropicErrorObjectUnion) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type

	switch tmp.Type {
	case string(AnthropicInvalidRequestErrorTypeInvalidRequestError):
		var b AnthropicInvalidRequestError
		if err := json.Unmarshal(data, &b); err == nil {
			u.InvalidRequestError, _ = json.Marshal(b)
		}
	case string(AnthropicAuthenticationErrorTypeAuthenticationError):
		var b AnthropicAuthenticationError
		if err := json.Unmarshal(data, &b); err == nil {
			u.AuthenticationError, _ = json.Marshal(b)
		}
	case string(AnthropicPermissionErrorTypePermissionError):
		var b AnthropicPermissionError
		if err := json.Unmarshal(data, &b); err == nil {
			u.PermissionError, _ = json.Marshal(b)
		}
	case string(AnthropicNotFoundErrorTypeNotFoundError):
		var b AnthropicNotFoundError
		if err := json.Unmarshal(data, &b); err == nil {
			u.NotFoundError, _ = json.Marshal(b)
		}
	case string(AnthropicRateLimitErrorTypeRateLimitError):
		var b AnthropicRateLimitError
		if err := json.Unmarshal(data, &b); err == nil {
			u.RateLimitError, _ = json.Marshal(b)
		}
	case string(AnthropicTimeoutErrorTypeTimeoutError):
		var b AnthropicTimeoutError
		if err := json.Unmarshal(data, &b); err == nil {
			u.TimeoutError, _ = json.Marshal(b)
		}
	case string(AnthropicOverloadedErrorTypeOverloadedError):
		var b AnthropicOverloadedError
		if err := json.Unmarshal(data, &b); err == nil {
			u.OverloadedError, _ = json.Marshal(b)
		}
	case string(AnthropicAPIErrorTypeAPIError):
		var b AnthropicAPIError
		if err := json.Unmarshal(data, &b); err == nil {
			u.APIError, _ = json.Marshal(b)
		}
	case string(AnthropicBillingErrorTypeBillingError):
		var b AnthropicBillingError
		if err := json.Unmarshal(data, &b); err == nil {
			u.BillingError, _ = json.Marshal(b)
		}
	}
	return nil
}
