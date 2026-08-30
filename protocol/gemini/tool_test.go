package gemini

import (
	"encoding/json"
	"testing"
)

func TestGeminiFunctionCall_Marshal(t *testing.T) {
	fc := GeminiFunctionCall{
		ID:   "call_1",
		Name: "get_weather",
		Args: map[string]any{"city": "Beijing"},
	}
	data, _ := json.Marshal(fc)
	want := `{"id":"call_1","name":"get_weather","args":{"city":"Beijing"}}`
	if string(data) != want {
		t.Errorf("got %s, want %s", string(data), want)
	}
}

func TestGeminiTool_Marshal(t *testing.T) {
	tool := GeminiTool{
		FunctionDeclarations: []*GeminiFunctionDeclaration{
			{Name: "f1", Description: "do f1", Parameters: &GeminiSchema{Type: "object"}},
		},
	}
	data, _ := json.Marshal(tool)
	if !contains(string(data), `"functionDeclarations"`) {
		t.Errorf("missing functionDeclarations: %s", string(data))
	}
}

func TestGeminiFunctionCallingConfig_Mode(t *testing.T) {
	cfg := GeminiFunctionCallingConfig{Mode: GeminiFunctionCallingModeAny, AllowedFunctionNames: []string{"f1"}}
	data, _ := json.Marshal(cfg)
	if !contains(string(data), `"mode":"ANY"`) {
		t.Errorf("missing mode ANY: %s", string(data))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf([]byte(s), []byte(substr)) >= 0
}
