package gemini

import (
	"encoding/json"
	"testing"
)

func TestGeminiErrorResponse_Unmarshal(t *testing.T) {
	raw := `{"error":{"code":401,"message":"bad key","status":"UNAUTHENTICATED"}}`
	var r GeminiErrorResponse
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r.Error == nil {
		t.Fatal("Error nil")
	}
	if r.Error.Code != 401 || r.Error.Status != "UNAUTHENTICATED" || r.Error.Message != "bad key" {
		t.Errorf("unexpected: %+v", r.Error)
	}
}

func TestGeminiAPIError_Marshal(t *testing.T) {
	e := GeminiAPIError{Code: 429, Message: "quota", Status: "RESOURCE_EXHAUSTED"}
	data, _ := json.Marshal(e)
	if !contains(string(data), `"status":"RESOURCE_EXHAUSTED"`) {
		t.Errorf("missing status: %s", string(data))
	}
}
