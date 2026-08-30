package frames

import (
	"encoding/json"
	"net/http"

	"hegq.net/mogo/pkg/protocol/anthropic"
)

// WriteAnthropicError 把 errormap.ErrorType 字符串值转换为 Anthropic 协议原生错误响应，
// 写入 w 并返回写入字节数。
//
// 请求侧约定：调用方负责设置 HTTP status 与 abort —— 本函数 *不* 调 WriteHeader，
// 避免与 gin 等 framework 的 status 设置产生 "superfluous WriteHeader" 错误。
// 这是与典型反向代理模式（"中间件先 c.Status / c.Abort，工具写 body"）一致的设计：
// frames 工具只看得到 ResponseWriter，不再管控 status。
//
// 调用方负责的事：
//   - c.Writer.Header().Set("Content-Type", "application/json")（HTTP 头）
//   - c.Status(status)（调用 WriteHeader）
//   - c.Abort()
//   - frames.WriteAnthropicError(c.Writer, ...)
//   - return err
//
// errType 是 errormap.ErrorType 的字符串取值（"invalid_request"/"authentication"/
// "permission"/"not_found"/"billing"/"rate_limit"/"overloaded"/"api_error"）。
// 未知值降级为 "api_error"（Anthropic 通用 5xx）。
//
// requestID 可选（AnthropicErrorResponse.RequestID 是 omitempty）——空值时省略字段。
func WriteAnthropicError(w http.ResponseWriter, errType, message, requestID string) (int, error) {
	body := anthropic.AnthropicErrorResponse{
		Type:      anthropic.AnthropicErrorTypeError,
		RequestID: requestID,
		Error:     mapErrormapToAnthropicObject(errType, message),
	}
	return w.Write(mustMarshal(body))
}

// mapErrormapToAnthropicObject 按 errormap 字符串值分发到 9 种 AnthropicErrorObjectUnion 变体。
// 未知值降级为 AnthropicAPIError（generic 500），保证 wire-format 永远合法。
func mapErrormapToAnthropicObject(errType, message string) *anthropic.AnthropicErrorObjectUnion {
	switch errType {
	case "invalid_request":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicInvalidRequestErrorTypeInvalidRequestError),
			InvalidRequestError: mustMarshal(anthropic.AnthropicInvalidRequestError{
				Type:    anthropic.AnthropicInvalidRequestErrorTypeInvalidRequestError,
				Message: message,
			}),
		}
	case "authentication":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicAuthenticationErrorTypeAuthenticationError),
			AuthenticationError: mustMarshal(anthropic.AnthropicAuthenticationError{
				Type:    anthropic.AnthropicAuthenticationErrorTypeAuthenticationError,
				Message: message,
			}),
		}
	case "permission":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicPermissionErrorTypePermissionError),
			PermissionError: mustMarshal(anthropic.AnthropicPermissionError{
				Type:    anthropic.AnthropicPermissionErrorTypePermissionError,
				Message: message,
			}),
		}
	case "not_found":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicNotFoundErrorTypeNotFoundError),
			NotFoundError: mustMarshal(anthropic.AnthropicNotFoundError{
				Type:    anthropic.AnthropicNotFoundErrorTypeNotFoundError,
				Message: message,
			}),
		}
	case "billing":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicBillingErrorTypeBillingError),
			BillingError: mustMarshal(anthropic.AnthropicBillingError{
				Type:    anthropic.AnthropicBillingErrorTypeBillingError,
				Message: message,
			}),
		}
	case "rate_limit":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicRateLimitErrorTypeRateLimitError),
			RateLimitError: mustMarshal(anthropic.AnthropicRateLimitError{
				Type:    anthropic.AnthropicRateLimitErrorTypeRateLimitError,
				Message: message,
			}),
		}
	case "upstream_quota_exceeded":
		// 上游套餐配额 / 速率限制耗尽（与本地 "billing" = 我方账户失效区别）。
		// 复用 rate_limit 变体；客户端按"限速"处理（重试 / 切换 fallback）。
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicRateLimitErrorTypeRateLimitError),
			RateLimitError: mustMarshal(anthropic.AnthropicRateLimitError{
				Type:    anthropic.AnthropicRateLimitErrorTypeRateLimitError,
				Message: message,
			}),
		}
	case "overloaded":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicOverloadedErrorTypeOverloadedError),
			OverloadedError: mustMarshal(anthropic.AnthropicOverloadedError{
				Type:    anthropic.AnthropicOverloadedErrorTypeOverloadedError,
				Message: message,
			}),
		}
	case "api_error":
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicAPIErrorTypeAPIError),
			APIError: mustMarshal(anthropic.AnthropicAPIError{
				Type:    anthropic.AnthropicAPIErrorTypeAPIError,
				Message: message,
			}),
		}
	default:
		return &anthropic.AnthropicErrorObjectUnion{
			Type: string(anthropic.AnthropicAPIErrorTypeAPIError),
			APIError: mustMarshal(anthropic.AnthropicAPIError{
				Type:    anthropic.AnthropicAPIErrorTypeAPIError,
				Message: message,
			}),
		}
	}
}

// mustMarshal 序列化任意结构。error 不会发生（结构稳定），但保持显式 error 返回。
func mustMarshal[T any](v T) []byte {
	data, _ := json.Marshal(v)
	return data
}
