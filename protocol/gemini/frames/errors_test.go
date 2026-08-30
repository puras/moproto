package frames

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteGeminiError_InvalidRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	n, err := WriteGeminiError(rr, http.StatusBadRequest, "invalid_request", "bad body", "")
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}
	var body GeminiErrorResponseWire
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error.Code != 400 || body.Error.Status != "INVALID_ARGUMENT" || body.Error.Message != "bad body" {
		t.Errorf("unexpected body: %+v", body.Error)
	}
}

func TestWriteGeminiError_RateLimit(t *testing.T) {
	rr := httptest.NewRecorder()
	_, _ = WriteGeminiError(rr, 429, "rate_limit", "quota exceeded", "req_123")
	var body GeminiErrorResponseWire
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body.Error.Status != "RESOURCE_EXHAUSTED" {
		t.Errorf("rate_limit should map to RESOURCE_EXHAUSTED, got %s", body.Error.Status)
	}
}

func TestWriteGeminiError_UpstreamQuotaExceeded(t *testing.T) {
	rr := httptest.NewRecorder()
	_, _ = WriteGeminiError(rr, 429, "upstream_quota_exceeded", "upstream quota hit limit", "")
	var body GeminiErrorResponseWire
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body.Error.Status != "RESOURCE_EXHAUSTED" {
		t.Errorf("upstream_quota_exceeded should map to RESOURCE_EXHAUSTED, got %s", body.Error.Status)
	}
}

func TestWriteGeminiError_UnknownErrType_Defaults(t *testing.T) {
	rr := httptest.NewRecorder()
	_, _ = WriteGeminiError(rr, 500, "foo_bar", "weird", "")
	var body GeminiErrorResponseWire
	_ = json.Unmarshal(rr.Body.Bytes(), &body)
	if body.Error.Status != "INTERNAL" {
		t.Errorf("unknown errtype should map to INTERNAL, got %s", body.Error.Status)
	}
}

// GeminiErrorResponseWire 用于测试 wire-format；与生产路径的 GeminiErrorResponse 等价。
type GeminiErrorResponseWire struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}
