package gemini

import (
	"encoding/json"
	"testing"
)

func TestGeminiRequest_GetFormat(t *testing.T) {
	r := &GeminiRequest{}
	if r.GetFormat() != "gemini" {
		t.Errorf("GetFormat: got %s, want gemini", r.GetFormat())
	}
}

func TestGeminiRequest_GetModel(t *testing.T) {
	// model 不在 body，由 handler 注入 extras.ResolvedModel；GetModel 返回 ""
	r := &GeminiRequest{GeminiGenerateContentRequest: &GeminiGenerateContentRequest{}}
	if r.GetModel() != "" {
		t.Errorf("GetModel: got %q, want empty", r.GetModel())
	}
}

func TestGeminiGenerateContentRequest_FullMarshal(t *testing.T) {
	req := GeminiGenerateContentRequest{
		Contents: []*GeminiContent{
			{Role: GeminiRoleUser, Parts: []*GeminiPart{{Text: "hi"}}},
		},
		SystemInstruction: &GeminiContent{Parts: []*GeminiPart{{Text: "You are helpful"}}},
		GenerationConfig:  &GeminiGenerationConfig{Temperature: float32Ptr(0.7)},
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !contains(string(data), `"systemInstruction":{"parts":[{"text":"You are helpful"}]}`) {
		t.Errorf("missing systemInstruction: %s", string(data))
	}
	if !contains(string(data), `"temperature":0.7`) {
		t.Errorf("missing temperature: %s", string(data))
	}
}

func float32Ptr(f float32) *float32 { return &f }
