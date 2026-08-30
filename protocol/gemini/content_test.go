package gemini

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestGeminiPart_TextMarshal(t *testing.T) {
	p := GeminiPart{Text: "hello"}
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"text":"hello"}`
	if string(data) != want {
		t.Errorf("got %s, want %s", string(data), want)
	}
}

func TestGeminiPart_InlineData_Base64RoundTrip(t *testing.T) {
	orig := []byte("hello world")
	b := GeminiBlob{MIMEType: "image/png", Data: orig}
	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("marshal blob: %v", err)
	}
	if !containsBytes(data, []byte(base64.StdEncoding.EncodeToString(orig))) {
		t.Errorf("expected base64-encoded %q, got %s", orig, string(data))
	}
	var back GeminiBlob
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal blob: %v", err)
	}
	if string(back.Data) != string(orig) {
		t.Errorf("round trip: got %q, want %q", back.Data, orig)
	}
}

func TestGeminiContent_RoleDefault(t *testing.T) {
	c := GeminiContent{Parts: []*GeminiPart{{Text: "hi"}}}
	data, _ := json.Marshal(c)
	if string(data) != `{"parts":[{"text":"hi"}]}` {
		t.Errorf("role should be omitted, got %s", string(data))
	}
}

func TestGeminiContent_RoleModel(t *testing.T) {
	c := GeminiContent{Role: GeminiRoleModel, Parts: []*GeminiPart{{Text: "ok"}}}
	data, _ := json.Marshal(c)
	if string(data) != `{"parts":[{"text":"ok"}],"role":"model"}` {
		t.Errorf("role should be model, got %s", string(data))
	}
}

func containsBytes(haystack, needle []byte) bool {
	return len(haystack) >= len(needle) && string(haystack)[:len(needle)] == string(needle) ||
		len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(haystack, needle []byte) int {
outer:
	for i := 0; i+len(needle) <= len(haystack); i++ {
		for j := range needle {
			if haystack[i+j] != needle[j] {
				continue outer
			}
		}
		return i
	}
	return -1
}
