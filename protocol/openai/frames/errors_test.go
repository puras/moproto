package frames

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// helper: 模拟调用方完整流程（status + content-type + frames.Write）。
func invoke(rec *httptest.ResponseRecorder, status int, errType, message string) error {
	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(status)
	_, err := WriteOpenAIError(rec, errType, message)
	return err
}

// TestWriteOpenAIError_AllErrormapTypes 覆盖 errormap 8 个 ErrorType 字符串值
// 在 OpenAI 错误协议下的映射（含 permission → permission_denied、overloaded → server）。
func TestWriteOpenAIError_AllErrormapTypes(t *testing.T) {
	cases := []struct {
		errType   string
		wantInner string
	}{
		{"invalid_request", `"invalid_request_error"`},
		{"authentication", `"authentication_error"`},
		{"permission", `"permission_denied_error"`},
		{"not_found", `"not_found_error"`},
		{"billing", `"billing_error"`},
		{"rate_limit", `"rate_limit_error"`},
		{"upstream_quota_exceeded", `"rate_limit_error"`}, // 复用 rate_limit（429）让客户端按 Retry-After 退避
		{"overloaded", `"server_error"`},
		{"api_error", `"server_error"`},
	}
	for _, tc := range cases {
		t.Run(tc.errType, func(t *testing.T) {
			rec := httptest.NewRecorder()
			if err := invoke(rec, 400, tc.errType, "msg-x"); err != nil {
				t.Fatalf("write err: %v", err)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("content-type = %q, want application/json", got)
			}
			if rec.Code != 400 {
				t.Errorf("status = %d, want 400", rec.Code)
			}
			var outer struct {
				Error struct {
					Message string `json:"message"`
					Type    string `json:"type"`
				} `json:"error"`
			}
			if jerr := json.Unmarshal(rec.Body.Bytes(), &outer); jerr != nil {
				t.Fatalf("body not valid JSON: %v", jerr)
			}
			if outer.Error.Type != strings.Trim(tc.wantInner, `"`) {
				t.Errorf("outer.error.type = %q, want %s", outer.Error.Type, tc.wantInner)
			}
			if outer.Error.Message != "msg-x" {
				t.Errorf("message lost: got %q", outer.Error.Message)
			}
		})
	}
}

// TestWriteOpenAIError_UnknownType 未知 errType 降级为 server_error。
func TestWriteOpenAIError_UnknownType(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := invoke(rec, 500, "totally_new_type", "unknown"); err != nil {
		t.Fatal(err)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"server_error"`) {
		t.Errorf("unknown errType should fall back to server_error, got: %s", body)
	}
}

// TestWriteOpenAIError_NoRequestID 验证 OpenAI 响应不带 request_id 字段（协议差异）。
func TestWriteOpenAIError_NoRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := invoke(rec, 401, "authentication", "no req_id"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(rec.Body.String(), "request_id") {
		t.Errorf("OpenAI protocol does not carry request_id, but body contains it: %s", rec.Body.String())
	}
}
