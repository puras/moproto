package frames

import (
	"encoding/json"
	"strings"
	"testing"

	"hegq.net/moproto/protocol/gemini"
)

func TestMarshalSSEFrame_GeminiResponse(t *testing.T) {
	resp := &gemini.GeminiGenerateContentResponse{
		Candidates: []*gemini.GeminiCandidate{
			{Content: &gemini.GeminiContent{Parts: []*gemini.GeminiPart{{Text: "hi"}}}},
		},
	}
	data, err := MarshalSSEFrame(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.HasPrefix(string(data), "data: ") {
		t.Errorf("missing data: prefix: %s", string(data))
	}
	if !strings.HasSuffix(string(data), "\n\n") {
		t.Errorf("missing trailing \\n\\n: %s", string(data))
	}
	// body 应是合法 JSON
	body := strings.TrimPrefix(string(data), "data: ")
	body = strings.TrimSuffix(body, "\n\n")
	var back gemini.GeminiGenerateContentResponse
	if err := json.Unmarshal([]byte(body), &back); err != nil {
		t.Errorf("body not valid JSON: %v\nbody=%s", err, body)
	}
}

func TestMarshalSSEFrame_WrongType(t *testing.T) {
	_, err := MarshalSSEFrame("not a response")
	if err == nil {
		t.Error("expected error for non-GeminiResponse input")
	}
}

func TestNormalizeSSEFrame_Empty(t *testing.T) {
	if NormalizeSSEFrame("") != "" {
		t.Errorf("empty should stay empty")
	}
}

func TestNormalizeSSEFrame_AddDataPrefix(t *testing.T) {
	got := NormalizeSSEFrame(`{"foo":1}`)
	if !strings.HasPrefix(got, "data: ") {
		t.Errorf("missing data: prefix: %s", got)
	}
	if !strings.HasSuffix(got, "\n\n") {
		t.Errorf("missing trailing \\n\\n: %s", got)
	}
}

func TestNormalizeSSEFrame_PreservesDataPrefix(t *testing.T) {
	got := NormalizeSSEFrame("data: {\"foo\":1}\n\n")
	if got != "data: {\"foo\":1}\n\n" {
		t.Errorf("unexpected transform: %s", got)
	}
}

// ============== RawData passthrough 测试（同协议流式透传） ==============

// TestMarshalSSEFrame_RawDataPassthrough 验证 RawData 非空时直接 emit 原始字节,
// 不走 json.Marshal 往返（保留 omitempty 零值字段与未知字段）。
//
// 镜像 OpenAI TestParseSSEEvent_PreservesRawBytes 与 Anthropic
// TestMarshalSSEFrameRaw_Verbatim 的模式 — 保证三协议对称。
func TestMarshalSSEFrame_RawDataPassthrough(t *testing.T) {
	// 含未知字段 groundingMetadata / urlContextMetadata / omitempty 零值 avgLogprobs
	body := `{"candidates":[{"content":{"parts":[{"text":"hi"}],"role":"model"},"finishReason":"STOP","index":0,"avgLogprobs":0.4}],"modelVersion":"gemini-2.5-pro","usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1,"totalTokenCount":2},"groundingMetadata":{"webSearchQueries":["q1","q2"]},"urlContextMetadata":{"urlMetadata":[{"url":"https://x.com"}]}}`
	resp := &gemini.GeminiGenerateContentResponse{
		RawData: []byte(body),
	}
	out, err := MarshalSSEFrame(resp)
	if err != nil {
		t.Fatalf("MarshalSSEFrame: %v", err)
	}
	want := "data: " + body + "\n\n"
	if string(out) != want {
		t.Errorf("passthrough output mismatch:\n got: %q\nwant: %q", out, want)
	}
}

// TestMarshalSSEFrame_BadJSONSurvives 验证非 JSON RawData 也能透传（不走 json.Marshal）。
func TestMarshalSSEFrame_BadJSONSurvives(t *testing.T) {
	body := `not-json-at-all-but-still-passthrough`
	resp := &gemini.GeminiGenerateContentResponse{
		RawData: []byte(body),
	}
	out, err := MarshalSSEFrame(resp)
	if err != nil {
		t.Fatalf("MarshalSSEFrame: %v", err)
	}
	want := "data: " + body + "\n\n"
	if string(out) != want {
		t.Errorf("passthrough should emit raw data frame, got %q want %q", out, want)
	}
}

// TestMarshalSSEFrame_TypedFallbackWhenRawDataEmpty 验证 RawData 为空时
// 走 typed marshal（非流式 / typed 转换路径）— 行为不变。
func TestMarshalSSEFrame_TypedFallbackWhenRawDataEmpty(t *testing.T) {
	resp := &gemini.GeminiGenerateContentResponse{
		ModelVersion: "gemini-2.5-pro",
	}
	out, err := MarshalSSEFrame(resp)
	if err != nil {
		t.Fatalf("MarshalSSEFrame: %v", err)
	}
	body := strings.TrimPrefix(string(out), "data: ")
	body = strings.TrimSuffix(body, "\n\n")
	if !strings.Contains(body, `"modelVersion":"gemini-2.5-pro"`) {
		t.Errorf("typed output should contain modelVersion, got: %q", body)
	}
}
