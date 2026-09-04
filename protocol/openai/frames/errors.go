// Package frames 提供 OpenAI ChatCompletions 协议错误响应的构造工具。
//
// 当前子包仅承载 errors.go（错误响应写入）；未来 SSE 帧构造等需求若有可在此追加。
// 命名 / 边界与 pkg/protocol/anthropic/frames 平行——一旦 OpenAI SSE 流式需要
// 类似的 UnmarshalSSEEvent/MarshalSSEFrame 工具时，再在本目录补充。
package frames

import (
	"encoding/json"
	"net/http"

	"github.com/puras/moproto/protocol/openai"
)

// WriteOpenAIError 把 errormap.ErrorType 字符串值转换为 OpenAI 协议原生错误响应，
// 写入 w 并返回写入字节数。
//
// 请求侧约定：调用方负责设置 HTTP status 与 Content-Type header —— 本函数 *不* 调
// WriteHeader / Header().Set，避免与 gin 等 framework 的 status 设置产生 "superfluous
// WriteHeader" 错误。
//
// 调用方负责的事：
//   - c.Writer.Header().Set("Content-Type", "application/json")
//   - c.Status(status)
//   - c.Abort()
//   - frames.WriteOpenAIError(c.Writer, ...)
//   - return err
//
// errType 是 errormap.ErrorType 的字符串取值（"invalid_request"/"authentication"/
// "permission"/"not_found"/"billing"/"rate_limit"/"overloaded"/"api_error"）。
// 未知值降级为 "server_error"。
//
// 与 Anthropic frames 的差异：
//   - OpenAI 错误响应顶层是 {"error":{"message","type"}}，不带 request_id 字段。
//   - OpenAI 错误的 "type" 字段与 errormap 枚举值不完全一一对应（例如 permission →
//     permission_denied_error，overloaded/api_error 都映射到 server_error）。
func WriteOpenAIError(w http.ResponseWriter, errType, message string) (int, error) {
	body := openai.OpenAIErrorResponse{
		Error: openai.OpenAIError{
			Message: message,
			Type:    mapErrormapToOpenAI(errType),
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}

// mapErrormapToOpenAI 把 errormap ErrorType 字符串值映射到 OpenAI 官方错误类型名。
// 未知值降级为 server_error（HTTP 5xx）。
func mapErrormapToOpenAI(errType string) string {
	switch errType {
	case "invalid_request":
		return string(openai.OpenAIErrorTypeInvalidRequest)
	case "authentication":
		return string(openai.OpenAIErrorTypeAuthentication)
	case "permission":
		return string(openai.OpenAIErrorTypePermissionDenied)
	case "not_found":
		return string(openai.OpenAIErrorTypeNotFound)
	case "billing":
		return string(openai.OpenAIErrorTypeBilling)
	case "rate_limit":
		return string(openai.OpenAIErrorTypeRateLimit)
	case "upstream_quota_exceeded":
		// 上游套餐配额 / 速率限制耗尽：复用 rate_limit（429 语义）让客户端按 Retry-After 退避
		// 不复用 billing：billing_error 暗示"网关支付失效"，会误导客户端走错路径
		return string(openai.OpenAIErrorTypeRateLimit)
	case "overloaded":
		return string(openai.OpenAIErrorTypeServer)
	case "api_error":
		return string(openai.OpenAIErrorTypeServer)
	default:
		return string(openai.OpenAIErrorTypeServer)
	}
}
