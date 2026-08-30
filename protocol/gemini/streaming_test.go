package gemini

import (
	"encoding/json"
	"strings"
	"testing"

	pkgusage "hegq.net/mogo/pkg/usage"
)

// jsonMarshal 简写：测试中调用 GeminiGenerateContentResponse.MarshalJSON。
func jsonMarshal(r *GeminiGenerateContentResponse) ([]byte, error) {
	return json.Marshal(r)
}

func TestExtractGeminiUsage_NilMetadata(t *testing.T) {
	if u := ExtractGeminiUsage(nil); u != nil {
		t.Errorf("expected nil, got %+v", u)
	}
}

func TestExtractGeminiUsage_AllFields(t *testing.T) {
	u := ExtractGeminiUsage(&GeminiUsageMetadata{
		PromptTokenCount:        10,
		CandidatesTokenCount:    5,
		TotalTokenCount:         15,
		CachedContentTokenCount: 3,
		ThoughtsTokenCount:      2,
		ToolUsePromptTokenCount: 1,
	})
	if u == nil {
		t.Fatal("nil usage")
	}
	if u.InputTokens != 10 {
		t.Errorf("input: got %d, want 10", u.InputTokens)
	}
	if u.OutputTokens != 5 {
		t.Errorf("output: got %d, want 5", u.OutputTokens)
	}
	if u.TotalTokens != 15 {
		t.Errorf("total: got %d, want 15", u.TotalTokens)
	}
	if u.Protocol != "gemini" {
		t.Errorf("protocol: got %q, want gemini", u.Protocol)
	}
	if u.Cache == nil || u.Cache.CachedTokens != 3 {
		t.Errorf("cache.CachedTokens: got %+v, want 3", u.Cache)
	}
	if u.Thinking == nil || u.Thinking.ThinkingTokens != 2 {
		t.Errorf("thinking: got %+v, want 2", u.Thinking)
	}
}

func TestParseGeminiStreamChunk_ValidLine(t *testing.T) {
	line := `data: {"candidates":[{"content":{"parts":[{"text":"hi"}]},"finishReason":"STOP"}],"usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1,"totalTokenCount":2}}`
	rawFrame := []byte(line + "\n\n")
	evt, done, err, usage := ParseGeminiStreamChunk(line, rawFrame)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if done {
		t.Error("should not be done")
	}
	if evt == nil {
		t.Fatal("nil event")
	}
	if usage == nil || usage.InputTokens != 1 {
		t.Errorf("usage: got %+v, want input=1", usage)
	}
	resp, ok := evt.(*GeminiGenerateContentResponse)
	if !ok {
		t.Fatalf("evt type: got %T, want *GeminiGenerateContentResponse", evt)
	}
	if string(resp.RawData) == "" {
		t.Error("RawData should be set for passthrough")
	}
}

func TestParseGeminiStreamChunk_EmptyLine(t *testing.T) {
	_, _, _, u := ParseGeminiStreamChunk("", nil)
	if u != nil {
		t.Errorf("empty line should return nil usage")
	}
}

func TestParseGeminiStreamChunk_NoDataPrefix(t *testing.T) {
	_, _, _, u := ParseGeminiStreamChunk(`{"foo":1}`, nil)
	if u != nil {
		t.Errorf("no data prefix should return nil")
	}
}

func TestExtractGeminiUsage_ProtocolField(_ *testing.T) {
	_ = pkgusage.Usage{} // 引用确保 import 不丢
}

// TestParseGeminiStreamChunk_PreservesRawBytes 验证 raw bytes passthrough：
//
// 1. RawData 字段填充为上游 data: 后原始 JSON
// 2. MarshalJSON 优先返回 RawData（passthrough 路径）
// 3. typed unmarshal 失败不再杀帧（RawData 已保真）
// 4. unknown field（groundingMetadata / urlContextMetadata / avgLogprobs 等）字节保真
func TestParseGeminiStreamChunk_PreservesRawBytes(t *testing.T) {
	// 含未知字段 groundingMetadata（typed struct 未声明）+ omitempty 零值 avgLogprobs
	body := `{"candidates":[{"content":{"parts":[{"text":"hi"}],"role":"model"},"finishReason":"STOP","index":0,"avgLogprobs":0.4}],"modelVersion":"gemini-2.5-pro","usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1,"totalTokenCount":2},"groundingMetadata":{"webSearchQueries":["q1","q2"],"searchEntryPoint":{}},"urlContextMetadata":{"urlMetadata":[{"url":"https://x.com"}]}}`
	line := `data: ` + body
	rawFrame := []byte(line + "\n\n")

	evt, _, err, usage := ParseGeminiStreamChunk(line, rawFrame)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	resp, ok := evt.(*GeminiGenerateContentResponse)
	if !ok {
		t.Fatalf("evt type: got %T, want *GeminiGenerateContentResponse", evt)
	}

	// RawData 必须 == 上游原始 body（passthrough 字节保真）
	if string(resp.RawData) != body {
		t.Errorf("RawData mismatch:\n got: %q\nwant: %q", resp.RawData, body)
	}

	// MarshalJSON 必须返回 RawData（即 passthrough 输出）
	out, mErr := jsonMarshal(resp)
	if mErr != nil {
		t.Fatalf("MarshalJSON: %v", mErr)
	}
	if string(out) != body {
		t.Errorf("MarshalJSON output mismatch:\n got: %q\nwant: %q", out, body)
	}

	// 抽 usage 仍然工作（typed 旁路）
	if usage == nil || usage.InputTokens != 1 {
		t.Errorf("usage: got %+v, want input=1", usage)
	}

	// 元数据抽取
	if resp.ModelVersion != "gemini-2.5-pro" {
		t.Errorf("ModelVersion: got %q, want %q", resp.ModelVersion, "gemini-2.5-pro")
	}
}

// TestParseGeminiStreamChunk_TypedMarshalFallback 验证 RawData 为空时
// MarshalJSON 走 typed marshal（非流式路径或 typed 转换后路径）。
func TestParseGeminiStreamChunk_TypedMarshalFallback(t *testing.T) {
	resp := &GeminiGenerateContentResponse{
		Candidates: []*GeminiCandidate{
			{Content: &GeminiContent{Role: "model"}, FinishReason: "STOP", Index: 0},
		},
		ModelVersion: "gemini-2.5-pro",
	}
	// RawData 为空 → 走 typed marshal
	if len(resp.RawData) != 0 {
		t.Fatal("RawData should be empty initially")
	}
	out, err := jsonMarshal(resp)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	// typed 输出应包含 modelVersion 与 candidates
	str := string(out)
	if !strings.Contains(str, `"modelVersion":"gemini-2.5-pro"`) {
		t.Errorf("typed output should contain modelVersion, got: %q", str)
	}
	if !strings.Contains(str, `"finishReason":"STOP"`) {
		t.Errorf("typed output should contain finishReason, got: %q", str)
	}
}

// TestParseGeminiStreamChunk_BadJSONSurvives 验证 typed unmarshal 失败
// 不再杀帧 — RawData 已保真。
//
// 注:完整 emit 路径测试见 frames/frame_test.go 的 TestMarshalSSEFrame_RawDataPassthrough
// (本包无法 import frames 包以避免 import cycle)。
func TestParseGeminiStreamChunk_BadJSONSurvives(t *testing.T) {
	// 非 JSON body（如 upstream bug / mock 错误）
	body := `not-json-at-all-but-still-passthrough`
	line := `data: ` + body
	rawFrame := []byte(line + "\n\n")

	evt, _, err, _ := ParseGeminiStreamChunk(line, rawFrame)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	resp, ok := evt.(*GeminiGenerateContentResponse)
	if !ok {
		t.Fatalf("evt type: got %T", evt)
	}
	if string(resp.RawData) != body {
		t.Errorf("RawData mismatch: got %q, want %q", resp.RawData, body)
	}
}
