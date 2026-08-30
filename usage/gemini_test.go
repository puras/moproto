package usage

import (
	"testing"
)

func TestFromGeminiBytes_ValidResponse(t *testing.T) {
	raw := []byte(`{"candidates":[{"content":{"parts":[{"text":"hi"}]}}],"usageMetadata":{"promptTokenCount":10,"candidatesTokenCount":5,"totalTokenCount":15,"cachedContentTokenCount":3,"thoughtsTokenCount":2}}`)
	u := FromGeminiBytes(raw)
	if u == nil {
		t.Fatal("nil usage")
	}
	// 2026-08-25 PR 改造：三层口径
	//   - InputTokens = promptTokenCount（协议原始值，含 cached）
	//   - FreshInputTokens = max(prompt - cached, 0) = 7
	//   - TotalTokens = totalTokenCount（透传上游）
	if u.InputTokens != 10 || u.OutputTokens != 5 || u.TotalTokens != 15 || u.FreshInputTokens != 7 {
		t.Errorf("base: got %+v want {InputTokens:10 OutputTokens:5 TotalTokens:15 FreshInputTokens:7}", u.Base)
	}
	if u.Cache == nil || u.Cache.CachedTokens != 3 {
		t.Errorf("cache: got %+v", u.Cache)
	}
	if u.Thinking == nil || u.Thinking.ThinkingTokens != 2 {
		t.Errorf("thinking: got %+v", u.Thinking)
	}
	if u.Protocol != "gemini" {
		t.Errorf("protocol: got %q", u.Protocol)
	}
}

func TestFromGeminiBytes_NoUsageMetadata(t *testing.T) {
	raw := []byte(`{"candidates":[{"content":{"parts":[{"text":"hi"}]}}]}`)
	if u := FromGeminiBytes(raw); u != nil {
		t.Errorf("expected nil when no usageMetadata, got %+v", u)
	}
}

func TestFromGeminiBytes_Empty(t *testing.T) {
	if u := FromGeminiBytes(nil); u != nil {
		t.Errorf("nil bytes should return nil usage")
	}
	if u := FromGeminiBytes([]byte{}); u != nil {
		t.Errorf("empty bytes should return nil usage")
	}
}

func TestFromGeminiBytes_MalformedJSON(t *testing.T) {
	if u := FromGeminiBytes([]byte("not json")); u != nil {
		t.Errorf("malformed JSON should return nil, got %+v", u)
	}
}

// TestFromGeminiBytes_CachedSubset 验证 2026-08-25 PR 改造后的三层口径：
//   - InputTokens = promptTokenCount（协议原始值，含 cached）
//   - FreshInputTokens = max(prompt - cached, 0)（去 cache 后的计费输入）
//   - TotalTokens = totalTokenCount（透传上游）
//   - Cache.CachedTokens = cachedContentTokenCount
func TestFromGeminiBytes_CachedSubset(t *testing.T) {
	raw := []byte(`{"usageMetadata":{"promptTokenCount":1000,"candidatesTokenCount":50,"totalTokenCount":1050,"cachedContentTokenCount":800}}`)
	u := FromGeminiBytes(raw)
	if u == nil {
		t.Fatal("nil usage")
	}
	// InputTokens 是协议原始值（含 cached）
	if u.InputTokens != 1000 {
		t.Errorf("InputTokens 应为上游原始值 1000（含 cached），实际 %d", u.InputTokens)
	}
	// FreshInputTokens 是去 cache 后的计费输入
	if u.FreshInputTokens != 200 {
		t.Errorf("FreshInputTokens 应为 1000 - 800 = 200，实际 %d", u.FreshInputTokens)
	}
	// TotalTokens 透传上游
	if u.TotalTokens != 1050 {
		t.Errorf("TotalTokens 应透传上游 1050，实际 %d", u.TotalTokens)
	}
	if u.Cache == nil || u.Cache.CachedTokens != 800 {
		t.Errorf("Cache.CachedTokens 期望 800，实际 %+v", u.Cache)
	}
	// 桶互斥：FreshInput + CacheRead = 全上下文 1000
	if u.FreshInputTokens+u.Cache.CachedTokens != 1000 {
		t.Errorf("桶互斥破坏：FreshInput %d + Cache.Cached %d != 1000",
			u.FreshInputTokens, u.Cache.CachedTokens)
	}
}

// TestFromGeminiBytes_NoCache 验证无 cached 时两字段相等。
func TestFromGeminiBytes_NoCache(t *testing.T) {
	raw := []byte(`{"usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":50,"totalTokenCount":150}}`)
	u := FromGeminiBytes(raw)
	if u == nil {
		t.Fatal("nil usage")
	}
	if u.InputTokens != 100 {
		t.Errorf("无 cached 时 InputTokens 应等于 promptTokenCount，实际 %d", u.InputTokens)
	}
	if u.FreshInputTokens != 100 {
		t.Errorf("无 cached 时 FreshInputTokens 应等于 InputTokens，实际 %d", u.FreshInputTokens)
	}
	if u.Cache != nil {
		t.Errorf("无 cached 时 Cache 应为 nil，实际 %+v", u.Cache)
	}
}
