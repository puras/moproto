package frames

import "testing"

// ============== OpenAI SSE 帧规整测试 ==============

func TestNormalizeSSEFrame_Empty(t *testing.T) {
	if got := NormalizeSSEFrame(""); got != "" {
		t.Errorf("空字符串应保持空，实际 %q", got)
	}
}

func TestNormalizeSSEFrame_BlankLine(t *testing.T) {
	// SSE 用空行做事件分隔符，必须原样保留
	tests := []string{"", "\n", "   ", "\t\n", "\n\n"}
	for _, in := range tests {
		got := NormalizeSSEFrame(in)
		// TrimRight 空白后是空 → 返回 "\n"
		if got != "\n" && got != "" {
			t.Errorf("空行 %q 应原样保留，实际 %q", in, got)
		}
	}
}

func TestNormalizeSSEFrame_CompleteFrame(t *testing.T) {
	// 完整帧原样保留
	in := "data: {\"id\":\"chatcmpl-xxx\"}\n\n"
	got := NormalizeSSEFrame(in)
	if got != in {
		t.Errorf("完整帧应原样保留，got %q want %q", got, in)
	}
}

func TestNormalizeSSEFrame_DataLineOnly(t *testing.T) {
	// 缺尾部空行 → 补一个 \n
	in := "data: {\"id\":\"chatcmpl-xxx\"}\n"
	want := "data: {\"id\":\"chatcmpl-xxx\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("data 行应补 \\n\\n，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_BareJSON(t *testing.T) {
	// 裸 JSON → 包成 data: ...\n\n
	in := `{"id":"chatcmpl-xxx","object":"chat.completion.chunk"}`
	want := "data: {\"id\":\"chatcmpl-xxx\",\"object\":\"chat.completion.chunk\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("裸 JSON 应包成 data 帧，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_BareJSONWithNewline(t *testing.T) {
	// 裸 JSON + 单换行 → 包成 data: ...\n\n
	in := `{"id":"chatcmpl-xxx"}` + "\n"
	want := "data: {\"id\":\"chatcmpl-xxx\"}\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("裸 JSON+换行 应包成 data 帧，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_PreservesDataPrefix(t *testing.T) {
	// 已有 data: 前缀时不重复添加
	in := "data: hello"
	want := "data: hello\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("已有前缀不应重复添加，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_TrailingSpaces(t *testing.T) {
	// 尾部空格（不影响 wire-format）应被 TrimRight
	in := "data: hello   "
	want := "data: hello\n\n"
	if got := NormalizeSSEFrame(in); got != want {
		t.Errorf("尾部空格应 trim，got %q want %q", got, want)
	}
}

func TestNormalizeSSEFrame_Done(t *testing.T) {
	// [DONE] 哨兵帧也走规整（虽然 caller 通常不会把 [DONE] 透传）
	in := "data: [DONE]\n\n"
	got := NormalizeSSEFrame(in)
	if got != in {
		t.Errorf("[DONE] 帧应原样保留，got %q want %q", got, in)
	}
}
