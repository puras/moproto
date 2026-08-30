package anthropic_test

import (
	"testing"

	anthropic "hegq.net/moproto/protocol/anthropic"
)

// TestGetCitationCitedText 验证 5 个实现都能取到 cited_text。
func TestGetCitationCitedText(t *testing.T) {
	cases := []struct {
		name string
		p    anthropic.AnthropicCitationsParam
		want string
	}{
		{
			name: "char_location",
			p:    anthropic.NewAnthropicCitationCharLocationParam("hello", 0, 5),
			want: "hello",
		},
		{
			name: "page_location",
			p:    anthropic.NewAnthropicCitationPageLocationParam("section 1", 2, 3),
			want: "section 1",
		},
		{
			name: "content_block_location",
			p:    anthropic.NewAnthropicCitationContentBlockLocationParam("block content", 0, 10),
			want: "block content",
		},
		{
			name: "search_result_location",
			p:    anthropic.NewAnthropicCitationSearchResultLocationParam("matched", 0, 5, 5),
			want: "matched",
		},
		{
			name: "web_search_result_location",
			p:    anthropic.NewAnthropicCitationWebSearchResultLocationParam("matched text", "enc_blob"),
			want: "matched text",
		},
	}
	for _, c := range cases {
		if got := anthropic.GetCitationCitedText(c.p); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

// TestGetCitationDocumentIndex 验证 3 个实现（char/page/content_block）能取到 document_index。
func TestGetCitationDocumentIndex(t *testing.T) {
	idx := int64(7)
	cases := []struct {
		name string
		p    anthropic.AnthropicCitationsParam
		want *int64
	}{
		{
			name: "char_location_with_doc",
			p: anthropic.AnthropicCitationCharLocationParam{
				Type:           anthropic.AnthropicCitationCharLocationTypeCharLocation,
				CitedText:      "x",
				DocumentIndex:  &idx,
				StartCharIndex: 0,
				EndCharIndex:   1,
			},
			want: &idx,
		},
		{
			name: "page_location_with_doc",
			p: anthropic.AnthropicCitationPageLocationParam{
				Type:            anthropic.AnthropicCitationPageLocationTypePageLocation,
				CitedText:       "x",
				DocumentIndex:   &idx,
				StartPageNumber: 1,
				EndPageNumber:   2,
			},
			want: &idx,
		},
		{
			name: "content_block_location_with_doc",
			p: anthropic.AnthropicCitationContentBlockLocationParam{
				Type:            anthropic.AnthropicCitationContentBlockLocationTypeContentBlockLocation,
				CitedText:       "x",
				DocumentIndex:   &idx,
				StartBlockIndex: 0,
				EndBlockIndex:   1,
			},
			want: &idx,
		},
		{
			name: "search_result_location_no_doc",
			p:    anthropic.NewAnthropicCitationSearchResultLocationParam("x", 0, 1, 1),
			want: nil,
		},
		{
			name: "web_search_result_location_no_doc",
			p:    anthropic.NewAnthropicCitationWebSearchResultLocationParam("x", "enc"),
			want: nil,
		},
	}
	for _, c := range cases {
		got := anthropic.GetCitationDocumentIndex(c.p)
		if !int64PtrEq(got, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, ptrToStr(got), ptrToStr(c.want))
		}
	}
}

// TestGetCitationCharLocationGetters 验证 char_location 专属 getter。
func TestGetCitationCharLocationGetters(t *testing.T) {
	p := anthropic.AnthropicCitationCharLocationParam{
		Type:           anthropic.AnthropicCitationCharLocationTypeCharLocation,
		CitedText:      "x",
		StartCharIndex: 10,
		EndCharIndex:   20,
	}
	if got := anthropic.GetCitationStartCharIndex(p); got == nil || *got != 10 {
		t.Errorf("StartCharIndex: got %v, want 10", got)
	}
	if got := anthropic.GetCitationEndCharIndex(p); got == nil || *got != 20 {
		t.Errorf("EndCharIndex: got %v, want 20", got)
	}
	// 不匹配其它变体时返回 nil
	if got := anthropic.GetCitationStartPageNumber(p); got != nil {
		t.Errorf("StartPageNumber: got %v, want nil", got)
	}
}

// TestGetCitationPageLocationGetters 验证 page_location 专属 getter。
func TestGetCitationPageLocationGetters(t *testing.T) {
	p := anthropic.AnthropicCitationPageLocationParam{
		Type:            anthropic.AnthropicCitationPageLocationTypePageLocation,
		CitedText:       "x",
		StartPageNumber: 5,
		EndPageNumber:   10,
	}
	if got := anthropic.GetCitationStartPageNumber(p); got == nil || *got != 5 {
		t.Errorf("StartPageNumber: got %v, want 5", got)
	}
	if got := anthropic.GetCitationEndPageNumber(p); got == nil || *got != 10 {
		t.Errorf("EndPageNumber: got %v, want 10", got)
	}
	if got := anthropic.GetCitationStartCharIndex(p); got != nil {
		t.Errorf("StartCharIndex: got %v, want nil", got)
	}
}

// TestGetCitationBlockIndexGetters 验证 block_index getter 聚合 2 实现。
func TestGetCitationBlockIndexGetters(t *testing.T) {
	p1 := anthropic.NewAnthropicCitationContentBlockLocationParam("x", 3, 7)
	if got := anthropic.GetCitationStartBlockIndex(p1); got == nil || *got != 3 {
		t.Errorf("ContentBlock StartBlockIndex: got %v, want 3", got)
	}
	if got := anthropic.GetCitationEndBlockIndex(p1); got == nil || *got != 7 {
		t.Errorf("ContentBlock EndBlockIndex: got %v, want 7", got)
	}

	p2 := anthropic.NewAnthropicCitationSearchResultLocationParam("x", 8, 12, 1)
	if got := anthropic.GetCitationStartBlockIndex(p2); got == nil || *got != 8 {
		t.Errorf("SearchResult StartBlockIndex: got %v, want 8", got)
	}
	if got := anthropic.GetCitationEndBlockIndex(p2); got == nil || *got != 12 {
		t.Errorf("SearchResult EndBlockIndex: got %v, want 12", got)
	}
}

// TestGetCitationSearchResultIndex 验证 search_result_index getter。
func TestGetCitationSearchResultIndex(t *testing.T) {
	p := anthropic.NewAnthropicCitationSearchResultLocationParam("x", 0, 1, 99)
	if got := anthropic.GetCitationSearchResultIndex(p); got == nil || *got != 99 {
		t.Errorf("SearchResultIndex: got %v, want 99", got)
	}
	// 不匹配
	p2 := anthropic.NewAnthropicCitationCharLocationParam("x", 0, 1)
	if got := anthropic.GetCitationSearchResultIndex(p2); got != nil {
		t.Errorf("SearchResultIndex (wrong type): got %v, want nil", got)
	}
}

// TestGetCitationEncryptedIndex 验证 encrypted_index getter。
func TestGetCitationEncryptedIndex(t *testing.T) {
	p := anthropic.NewAnthropicCitationWebSearchResultLocationParam("x", "enc_blob_xyz")
	if got := anthropic.GetCitationEncryptedIndex(p); got != "enc_blob_xyz" {
		t.Errorf("EncryptedIndex: got %q, want enc_blob_xyz", got)
	}
	p2 := anthropic.NewAnthropicCitationCharLocationParam("x", 0, 1)
	if got := anthropic.GetCitationEncryptedIndex(p2); got != "" {
		t.Errorf("EncryptedIndex (wrong type): got %q, want empty", got)
	}
}

// TestGetCitationTitle 验证 title getter 聚合 2 实现。
func TestGetCitationTitle(t *testing.T) {
	title := "result title"
	p1 := anthropic.AnthropicCitationSearchResultLocationParam{
		Type:              anthropic.AnthropicCitationSearchResultLocationTypeSearchResultLocation,
		CitedText:         "x",
		StartBlockIndex:   0,
		EndBlockIndex:     1,
		SearchResultIndex: 1,
		Title:             &title,
	}
	if got := anthropic.GetCitationTitle(p1); got == nil || *got != "result title" {
		t.Errorf("SearchResult Title: got %v, want result title", got)
	}
	p2 := anthropic.AnthropicCitationWebSearchResultLocationParam{
		Type:           anthropic.AnthropicCitationWebSearchResultLocationTypeWebSearchResultLocation,
		CitedText:      "x",
		EncryptedIndex: "enc",
		Title:          &title,
	}
	if got := anthropic.GetCitationTitle(p2); got == nil || *got != "result title" {
		t.Errorf("WebSearch Title: got %v, want result title", got)
	}
}

// TestGetCitationURL 验证 url getter。
func TestGetCitationURL(t *testing.T) {
	url := "https://example.com"
	p := anthropic.AnthropicCitationWebSearchResultLocationParam{
		Type:           anthropic.AnthropicCitationWebSearchResultLocationTypeWebSearchResultLocation,
		CitedText:      "x",
		EncryptedIndex: "enc",
		URL:            &url,
	}
	if got := anthropic.GetCitationURL(p); got != "https://example.com" {
		t.Errorf("URL: got %q, want https://example.com", got)
	}
	// nil URL
	p2 := anthropic.NewAnthropicCitationWebSearchResultLocationParam("x", "enc")
	if got := anthropic.GetCitationURL(p2); got != "" {
		t.Errorf("URL (nil): got %q, want empty", got)
	}
}

// TestGetCitationSource 验证 source getter。
func TestGetCitationSource(t *testing.T) {
	src := "src"
	p := anthropic.AnthropicCitationSearchResultLocationParam{
		Type:              anthropic.AnthropicCitationSearchResultLocationTypeSearchResultLocation,
		CitedText:         "x",
		StartBlockIndex:   0,
		EndBlockIndex:     1,
		SearchResultIndex: 1,
		Source:            &src,
	}
	if got := anthropic.GetCitationSource(p); got != "src" {
		t.Errorf("Source: got %q, want src", got)
	}
	p2 := anthropic.NewAnthropicCitationSearchResultLocationParam("x", 0, 1, 1)
	if got := anthropic.GetCitationSource(p2); got != "" {
		t.Errorf("Source (nil): got %q, want empty", got)
	}
}

// TestGetCitationType 验证 type 字段 getter。
func TestGetCitationType(t *testing.T) {
	cases := []struct {
		name string
		p    anthropic.AnthropicCitationsParam
		want string
	}{
		{
			name: "char_location",
			p:    anthropic.NewAnthropicCitationCharLocationParam("x", 0, 1),
			want: "char_location",
		},
		{
			name: "page_location",
			p:    anthropic.NewAnthropicCitationPageLocationParam("x", 1, 2),
			want: "page_location",
		},
		{
			name: "content_block_location",
			p:    anthropic.NewAnthropicCitationContentBlockLocationParam("x", 0, 1),
			want: "content_block_location",
		},
		{
			name: "search_result_location",
			p:    anthropic.NewAnthropicCitationSearchResultLocationParam("x", 0, 1, 1),
			want: "search_result_location",
		},
		{
			name: "web_search_result_location",
			p:    anthropic.NewAnthropicCitationWebSearchResultLocationParam("x", "enc"),
			want: "web_search_result_location",
		},
	}
	for _, c := range cases {
		if got := anthropic.GetCitationType(c.p); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

// TestGetCitationDocumentTitle 验证 document_title getter 聚合 3 实现。
func TestGetCitationDocumentTitle(t *testing.T) {
	title := "doc title"
	p := anthropic.AnthropicCitationCharLocationParam{
		Type:           anthropic.AnthropicCitationCharLocationTypeCharLocation,
		CitedText:      "x",
		DocumentTitle:  &title,
		StartCharIndex: 0,
		EndCharIndex:   1,
	}
	if got := anthropic.GetCitationDocumentTitle(p); got == nil || *got != "doc title" {
		t.Errorf("DocumentTitle: got %v, want doc title", got)
	}
	// 不匹配
	p2 := anthropic.NewAnthropicCitationSearchResultLocationParam("x", 0, 1, 1)
	if got := anthropic.GetCitationDocumentTitle(p2); got != nil {
		t.Errorf("DocumentTitle (wrong type): got %v, want nil", got)
	}
}

// int64PtrEq 比较两个 *int64 是否相等。
func int64PtrEq(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// ptrToStr 格式化 *int64 用于错误输出。
func ptrToStr(p *int64) string {
	if p == nil {
		return "<nil>"
	}
	return formatInt(*p)
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
