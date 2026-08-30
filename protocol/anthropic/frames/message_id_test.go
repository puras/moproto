package frames

import (
	"regexp"
	"strings"
	"testing"
)

var messageIDPattern = regexp.MustCompile(`^msg_[0-9a-f]{24}$`)

// TestGenerateMessageID_Format 验证生成的 id 严格符合 "msg_<24位hex>" 格式。
func TestGenerateMessageID_Format(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := GenerateMessageID()
		if !messageIDPattern.MatchString(id) {
			t.Fatalf("invalid message id format: %q", id)
		}
	}
}

// TestGenerateMessageID_Unique 验证连续调用不会产生重复 id（统计独立性）。
func TestGenerateMessageID_Unique(t *testing.T) {
	seen := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		id := GenerateMessageID()
		if seen[id] {
			t.Fatalf("duplicate id after %d calls: %q", i, id)
		}
		seen[id] = true
	}
}

// TestFallbackMessageID_LegalFormat 兜底常量本身也必须通过格式校验。
func TestFallbackMessageID_LegalFormat(t *testing.T) {
	if !messageIDPattern.MatchString(fallbackMessageID) {
		t.Fatalf("fallback id violates the same format: %q", fallbackMessageID)
	}
	if !strings.HasPrefix(fallbackMessageID, "msg_") {
		t.Fatalf("fallback id missing prefix: %q", fallbackMessageID)
	}
}
