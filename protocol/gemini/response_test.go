package gemini

import (
	"encoding/json"
	"testing"
)

func TestGeminiGenerateContentResponse_Marshal(t *testing.T) {
	resp := GeminiGenerateContentResponse{
		Candidates: []*GeminiCandidate{
			{
				Content:      &GeminiContent{Role: GeminiRoleModel, Parts: []*GeminiPart{{Text: "hi"}}},
				FinishReason: "STOP",
				Index:        0,
			},
		},
		UsageMetadata: &GeminiUsageMetadata{
			PromptTokenCount:     10,
			CandidatesTokenCount: 5,
			TotalTokenCount:      15,
		},
		ModelVersion: "gemini-2.5-pro",
	}
	data, _ := json.Marshal(resp)
	if !contains(string(data), `"finishReason":"STOP"`) {
		t.Errorf("missing finishReason: %s", string(data))
	}
	if !contains(string(data), `"promptTokenCount":10`) {
		t.Errorf("missing promptTokenCount: %s", string(data))
	}
	if !contains(string(data), `"modelVersion":"gemini-2.5-pro"`) {
		t.Errorf("missing modelVersion: %s", string(data))
	}
}

func TestGeminiUsageMetadata_AllFields(t *testing.T) {
	u := GeminiUsageMetadata{
		PromptTokenCount:        10,
		CandidatesTokenCount:    5,
		TotalTokenCount:         15,
		CachedContentTokenCount: 3,
		ThoughtsTokenCount:      2,
		ToolUsePromptTokenCount: 1,
	}
	data, _ := json.Marshal(u)
	for _, s := range []string{`"promptTokenCount":10`, `"candidatesTokenCount":5`, `"totalTokenCount":15`,
		`"cachedContentTokenCount":3`, `"thoughtsTokenCount":2`, `"toolUsePromptTokenCount":1`} {
		if !contains(string(data), s) {
			t.Errorf("missing %s in %s", s, string(data))
		}
	}
}
