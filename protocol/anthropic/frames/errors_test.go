package frames

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

// helper: 模拟调用方完整流程（status + content-type + frames.Write）。
func invoke(rec *httptest.ResponseRecorder, status int, errType, message, requestID string) error {
	rec.Header().Set("Content-Type", "application/json")
	rec.WriteHeader(status)
	_, err := WriteAnthropicError(rec, errType, message, requestID)
	return err
}

// TestWriteAnthropicError_AllErrormapTypes 覆盖 errormap 8 个 ErrorType 字符串值。
func TestWriteAnthropicError_AllErrormapTypes(t *testing.T) {
	cases := []struct {
		errType   string
		wantInner string
	}{
		{"invalid_request", `"invalid_request_error"`},
		{"authentication", `"authentication_error"`},
		{"permission", `"permission_error"`},
		{"not_found", `"not_found_error"`},
		{"billing", `"billing_error"`},
		{"rate_limit", `"rate_limit_error"`},
		{"upstream_quota_exceeded", `"rate_limit_error"`}, // 复用 rate_limit 变体
		{"overloaded", `"overloaded_error"`},
		{"api_error", `"api_error"`},
	}
	for _, tc := range cases {
		t.Run(tc.errType, func(t *testing.T) {
			rec := httptest.NewRecorder()
			if err := invoke(rec, 400, tc.errType, "msg-x", "req-1"); err != nil {
				t.Fatalf("write err: %v", err)
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("content-type = %q, want application/json", got)
			}
			if rec.Code != 400 {
				t.Errorf("status = %d, want 400", rec.Code)
			}
			var outer struct {
				Type      string `json:"type"`
				Error     any    `json:"error"`
				RequestID string `json:"request_id"`
			}
			if jerr := json.Unmarshal(rec.Body.Bytes(), &outer); jerr != nil {
				t.Fatalf("body not valid JSON: %v, body=%s", jerr, rec.Body.String())
			}
			if outer.Type != "error" {
				t.Errorf("outer.type = %q, want error", outer.Type)
			}
			if outer.RequestID != "req-1" {
				t.Errorf("request_id = %q, want req-1", outer.RequestID)
			}
			if !strings.Contains(rec.Body.String(), tc.wantInner) {
				t.Errorf("body missing %s, got: %s", tc.wantInner, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), `"msg-x"`) {
				t.Errorf("body missing message, got: %s", rec.Body.String())
			}
		})
	}
}

// TestWriteAnthropicError_UnknownType 未知 errType 降级为 api_error。
func TestWriteAnthropicError_UnknownType(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := invoke(rec, 500, "totally_new_type", "unknown", ""); err != nil {
		t.Fatal(err)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"api_error"`) {
		t.Errorf("unknown errType should fall back to api_error, got: %s", body)
	}
	if strings.Contains(body, `"request_id"`) {
		t.Errorf("empty request_id should omit field, got: %s", body)
	}
}
