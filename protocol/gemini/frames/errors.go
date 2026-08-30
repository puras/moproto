package frames

import (
	"encoding/json"
	"net/http"

	"hegq.net/moproto/protocol/gemini"
)

// WriteGeminiError 把 errormap.ErrorType 字符串值转换为 Gemini 协议原生错误响应，
// 写入 w 并返回写入字节数。
//
// 调用方负责：
//   - c.Writer.Header().Set("Content-Type", "application/json")
//   - c.Status(status) → WriteHeader
//   - c.Abort()
//   - frames.WriteGeminiError(c.Writer, ...)
//   - return err
//
// errType 取值见 errormap.ErrorType（"invalid_request" / "authentication" /
// "permission" / "not_found" / "billing" / "rate_limit" / "overloaded" / "api_error"）。
// 未知值降级为 "INTERNAL"。
//
// requestID 可选（Gemini 错误响应无 request_id 字段——纯协议差异）。
func WriteGeminiError(w http.ResponseWriter, status int, errType, message, requestID string) (int, error) {
	_ = requestID // Gemini 协议无 request_id 字段，保留参数名以与 frames.WriteAnthropicError 对齐
	apiStatus := mapErrormapToGeminiStatus(errType)
	body := gemini.GeminiErrorResponse{
		Error: &gemini.GeminiAPIError{
			Code:    status,
			Message: message,
			Status:  apiStatus,
		},
	}
	return w.Write(mustMarshal(body))
}

// mapErrormapToGeminiStatus 按 errormap 字符串值分发到 Google API status 字符串。
//
// 未知值降级为 "INTERNAL"（generic 500）。
func mapErrormapToGeminiStatus(errType string) string {
	switch errType {
	case "invalid_request":
		return "INVALID_ARGUMENT"
	case "authentication":
		return "UNAUTHENTICATED"
	case "permission":
		return "PERMISSION_DENIED"
	case "not_found":
		return "NOT_FOUND"
	case "billing":
		return "FAILED_PRECONDITION"
	case "rate_limit":
		return "RESOURCE_EXHAUSTED"
	case "upstream_quota_exceeded":
		// 上游套餐配额 / 速率限制耗尽（与本地 "billing" = FAILED_PRECONDITION 区别）。
		// 复用 rate_limit 映射为 RESOURCE_EXHAUSTED。
		return "RESOURCE_EXHAUSTED"
	case "overloaded":
		return "UNAVAILABLE"
	case "api_error":
		return "INTERNAL"
	default:
		return "INTERNAL"
	}
}

func mustMarshal[T any](v T) []byte {
	data, _ := json.Marshal(v)
	return data
}
