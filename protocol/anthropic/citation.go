package anthropic

// AnthropicCitationsParam
// ============== AnthropicCitationsParam union ==============
//
// AnthropicCitationsParam 是引文位置元数据的 union 接口。
// 对齐官方 anthropic-sdk-go TextCitationParamUnion（message.go:5784-5836）。
//
// 5 种实现按 type 字段字符串路由：
//   - AnthropicCitationCharLocationParam            （"char_location"）
//   - AnthropicCitationPageLocationParam            （"page_location"）
//   - AnthropicCitationContentBlockLocationParam    （"content_block_location"）
//   - AnthropicCitationSearchResultLocationParam    （"search_result_location"）
//   - AnthropicCitationWebSearchResultLocationParam （"web_search_result_location"）
//
// 通过 isCitationsParam() 私有方法封闭实现集；调用方通过 type switch 或各实现的
// Type 字段分发，详见 doc.go。
//
// 与 AnthropicCitationsConfig 的关系：Config（common.go）是块级开关
// （"是否为本 block 生成 citations"），本 union 是具体引文元数据。
// 两者作用层级不同，因此本 PR 不修改 message.go 现有的 Citations 字段。
type AnthropicCitationsParam interface {
	isCitationsParam()
}

// ============== 类型化 Type 字段 ==============

// AnthropicCitationCharLocationType char_location 的 type 字段类型。
type AnthropicCitationCharLocationType string

// AnthropicCitationCharLocationTypeCharLocation char_location 唯一合法取值。
const AnthropicCitationCharLocationTypeCharLocation AnthropicCitationCharLocationType = "char_location"

// AnthropicCitationPageLocationType page_location 的 type 字段类型。
type AnthropicCitationPageLocationType string

// AnthropicCitationPageLocationTypePageLocation page_location 唯一合法取值。
const AnthropicCitationPageLocationTypePageLocation AnthropicCitationPageLocationType = "page_location"

// AnthropicCitationContentBlockLocationType content_block_location 的 type 字段类型。
type AnthropicCitationContentBlockLocationType string

// AnthropicCitationContentBlockLocationTypeContentBlockLocation content_block_location 唯一合法取值。
const AnthropicCitationContentBlockLocationTypeContentBlockLocation AnthropicCitationContentBlockLocationType = "content_block_location"

// AnthropicCitationSearchResultLocationType search_result_location 的 type 字段类型。
type AnthropicCitationSearchResultLocationType string

// AnthropicCitationSearchResultLocationTypeSearchResultLocation search_result_location 唯一合法取值。
const AnthropicCitationSearchResultLocationTypeSearchResultLocation AnthropicCitationSearchResultLocationType = "search_result_location"

// AnthropicCitationWebSearchResultLocationType web_search_result_location 的 type 字段类型。
type AnthropicCitationWebSearchResultLocationType string

// AnthropicCitationWebSearchResultLocationTypeWebSearchResultLocation web_search_result_location 唯一合法取值。
const AnthropicCitationWebSearchResultLocationTypeWebSearchResultLocation AnthropicCitationWebSearchResultLocationType = "web_search_result_location"

// ============== 5 个实现 ==============

// AnthropicCitationCharLocationParam 字符级引文位置（普通文本）。
//
// 对齐官方 anthropic-sdk-go CitationCharLocationParam（message.go:582-599）。
// 必填：Type + CitedText + StartCharIndex + EndCharIndex；DocumentIndex/DocumentTitle 可选。
type AnthropicCitationCharLocationParam struct {
	Type           AnthropicCitationCharLocationType `json:"type"` // 固定 "char_location"
	CitedText      string                            `json:"cited_text"`
	DocumentIndex  *int64                            `json:"document_index,omitempty"`
	DocumentTitle  *string                           `json:"document_title,omitempty"`
	StartCharIndex int64                             `json:"start_char_index"`
	EndCharIndex   int64                             `json:"end_char_index"`
}

func (AnthropicCitationCharLocationParam) isCitationsParam() {}

// AnthropicCitationPageLocationParam PDF 页面级引文位置。
//
// 对齐官方 anthropic-sdk-go CitationPageLocationParam（message.go:705-722）。
// 必填：Type + CitedText + StartPageNumber + EndPageNumber；FileID/DocumentIndex/DocumentTitle 可选。
type AnthropicCitationPageLocationParam struct {
	Type            AnthropicCitationPageLocationType `json:"type"` // 固定 "page_location"
	CitedText       string                            `json:"cited_text"`
	DocumentIndex   *int64                            `json:"document_index,omitempty"`
	DocumentTitle   *string                           `json:"document_title,omitempty"`
	FileID          *string                           `json:"file_id,omitempty"`
	StartPageNumber int64                             `json:"start_page_number"`
	EndPageNumber   int64                             `json:"end_page_number"`
}

func (AnthropicCitationPageLocationParam) isCitationsParam() {}

// AnthropicCitationContentBlockLocationParam 自定义内容块级引文位置。
//
// 对齐官方 anthropic-sdk-go CitationContentBlockLocationParam（message.go:643-672）。
// 必填：Type + CitedText + StartBlockIndex + EndBlockIndex；FileID/DocumentIndex/DocumentTitle 可选。
type AnthropicCitationContentBlockLocationParam struct {
	Type            AnthropicCitationContentBlockLocationType `json:"type"` // 固定 "content_block_location"
	CitedText       string                                    `json:"cited_text"`
	DocumentIndex   *int64                                    `json:"document_index,omitempty"`
	DocumentTitle   *string                                   `json:"document_title,omitempty"`
	FileID          *string                                   `json:"file_id,omitempty"`
	StartBlockIndex int64                                     `json:"start_block_index"`
	EndBlockIndex   int64                                     `json:"end_block_index"`
}

func (AnthropicCitationContentBlockLocationParam) isCitationsParam() {}

// AnthropicCitationSearchResultLocationParam 工具搜索结果级引文位置。
//
// 对齐官方 anthropic-sdk-go CitationSearchResultLocationParam（message.go:726-763）。
// 必填：Type + CitedText + StartBlockIndex + EndBlockIndex + SearchResultIndex；Title/Source 可选。
type AnthropicCitationSearchResultLocationParam struct {
	Type              AnthropicCitationSearchResultLocationType `json:"type"` // 固定 "search_result_location"
	CitedText         string                                    `json:"cited_text"`
	StartBlockIndex   int64                                     `json:"start_block_index"`
	EndBlockIndex     int64                                     `json:"end_block_index"`
	SearchResultIndex int64                                     `json:"search_result_index"`
	Title             *string                                   `json:"title,omitempty"`
	Source            *string                                   `json:"source,omitempty"`
}

func (AnthropicCitationSearchResultLocationParam) isCitationsParam() {}

// AnthropicCitationWebSearchResultLocationParam 网络搜索结果级引文位置。
//
// 对齐官方 anthropic-sdk-go CitationWebSearchResultLocationParam（message.go:766-783）。
// 必填：Type + CitedText + EncryptedIndex；Title/URL 可选。
type AnthropicCitationWebSearchResultLocationParam struct {
	Type           AnthropicCitationWebSearchResultLocationType `json:"type"` // 固定 "web_search_result_location"
	CitedText      string                                       `json:"cited_text"`
	EncryptedIndex string                                       `json:"encrypted_index"`
	Title          *string                                      `json:"title,omitempty"`
	URL            *string                                      `json:"url,omitempty"`
}

func (AnthropicCitationWebSearchResultLocationParam) isCitationsParam() {}

// ============== 构造器 ==============
//
// 构造器函数避免调用方手填 Type 字段导致的 wire-format 字符串错误。
// 对齐官方 anthropic-sdk-go 的 NewCitation*LocationParam 模式。

// NewAnthropicCitationCharLocationParam 构造 char_location 引文。
func NewAnthropicCitationCharLocationParam(citedText string, startCharIndex, endCharIndex int64) AnthropicCitationCharLocationParam {
	return AnthropicCitationCharLocationParam{
		Type:           AnthropicCitationCharLocationTypeCharLocation,
		CitedText:      citedText,
		StartCharIndex: startCharIndex,
		EndCharIndex:   endCharIndex,
	}
}

// NewAnthropicCitationPageLocationParam 构造 page_location 引文。
func NewAnthropicCitationPageLocationParam(citedText string, startPageNumber, endPageNumber int64) AnthropicCitationPageLocationParam {
	return AnthropicCitationPageLocationParam{
		Type:            AnthropicCitationPageLocationTypePageLocation,
		CitedText:       citedText,
		StartPageNumber: startPageNumber,
		EndPageNumber:   endPageNumber,
	}
}

// NewAnthropicCitationContentBlockLocationParam 构造 content_block_location 引文。
func NewAnthropicCitationContentBlockLocationParam(citedText string, startBlockIndex, endBlockIndex int64) AnthropicCitationContentBlockLocationParam {
	return AnthropicCitationContentBlockLocationParam{
		Type:            AnthropicCitationContentBlockLocationTypeContentBlockLocation,
		CitedText:       citedText,
		StartBlockIndex: startBlockIndex,
		EndBlockIndex:   endBlockIndex,
	}
}

// NewAnthropicCitationSearchResultLocationParam 构造 search_result_location 引文。
func NewAnthropicCitationSearchResultLocationParam(citedText string, startBlockIndex, endBlockIndex, searchResultIndex int64) AnthropicCitationSearchResultLocationParam {
	return AnthropicCitationSearchResultLocationParam{
		Type:              AnthropicCitationSearchResultLocationTypeSearchResultLocation,
		CitedText:         citedText,
		StartBlockIndex:   startBlockIndex,
		EndBlockIndex:     endBlockIndex,
		SearchResultIndex: searchResultIndex,
	}
}

// NewAnthropicCitationWebSearchResultLocationParam 构造 web_search_result_location 引文。
func NewAnthropicCitationWebSearchResultLocationParam(citedText, encryptedIndex string) AnthropicCitationWebSearchResultLocationParam {
	return AnthropicCitationWebSearchResultLocationParam{
		Type:           AnthropicCitationWebSearchResultLocationTypeWebSearchResultLocation,
		CitedText:      citedText,
		EncryptedIndex: encryptedIndex,
	}
}

// ============== AnthropicCitationsParam Getter Helpers ==============
//
// 对齐官方 anthropic-sdk-go TextCitationParamUnion 的 getter
//（message.go:5820-5967）。返回类型语义：
//   - *int64 / *string：nil 表示字段未设置（与官方 SDK 一致）
//   - string（值类型）：空字符串表示字段未设置（简化调用）
//
// 使用模式：
//
//	p := anthropic.AnthropicCitationCharLocationParam{...}
//	if cited := anthropic.GetCitationCitedText(p); cited != "" { ... }
//	if idx := anthropic.GetCitationStartCharIndex(p); idx != nil { ... }
//
// 所有函数对未匹配的 union 变体返回零值（nil / 空字符串），调用方
// 无需做 type switch 兜底。

// GetCitationCitedText 返回 cited_text 字段（5 个实现聚合）。
func GetCitationCitedText(p AnthropicCitationsParam) string {
	switch v := p.(type) {
	case AnthropicCitationCharLocationParam:
		return v.CitedText
	case AnthropicCitationPageLocationParam:
		return v.CitedText
	case AnthropicCitationContentBlockLocationParam:
		return v.CitedText
	case AnthropicCitationSearchResultLocationParam:
		return v.CitedText
	case AnthropicCitationWebSearchResultLocationParam:
		return v.CitedText
	}
	return ""
}

// GetCitationDocumentIndex 返回 document_index（3 个实现聚合：char / page / content_block）。
func GetCitationDocumentIndex(p AnthropicCitationsParam) *int64 {
	switch v := p.(type) {
	case AnthropicCitationCharLocationParam:
		return v.DocumentIndex
	case AnthropicCitationPageLocationParam:
		return v.DocumentIndex
	case AnthropicCitationContentBlockLocationParam:
		return v.DocumentIndex
	}
	return nil
}

// GetCitationDocumentTitle 返回 document_title（3 个实现聚合）。
func GetCitationDocumentTitle(p AnthropicCitationsParam) *string {
	switch v := p.(type) {
	case AnthropicCitationCharLocationParam:
		return v.DocumentTitle
	case AnthropicCitationPageLocationParam:
		return v.DocumentTitle
	case AnthropicCitationContentBlockLocationParam:
		return v.DocumentTitle
	}
	return nil
}

// GetCitationStartCharIndex 返回 char_location 的 start_char_index。
func GetCitationStartCharIndex(p AnthropicCitationsParam) *int64 {
	if v, ok := p.(AnthropicCitationCharLocationParam); ok {
		return &v.StartCharIndex
	}
	return nil
}

// GetCitationEndCharIndex 返回 char_location 的 end_char_index。
func GetCitationEndCharIndex(p AnthropicCitationsParam) *int64 {
	if v, ok := p.(AnthropicCitationCharLocationParam); ok {
		return &v.EndCharIndex
	}
	return nil
}

// GetCitationStartPageNumber 返回 page_location 的 start_page_number。
func GetCitationStartPageNumber(p AnthropicCitationsParam) *int64 {
	if v, ok := p.(AnthropicCitationPageLocationParam); ok {
		return &v.StartPageNumber
	}
	return nil
}

// GetCitationEndPageNumber 返回 page_location 的 end_page_number。
func GetCitationEndPageNumber(p AnthropicCitationsParam) *int64 {
	if v, ok := p.(AnthropicCitationPageLocationParam); ok {
		return &v.EndPageNumber
	}
	return nil
}

// GetCitationStartBlockIndex 返回 content_block_location 或 search_result_location
// 的 start_block_index。
func GetCitationStartBlockIndex(p AnthropicCitationsParam) *int64 {
	switch v := p.(type) {
	case AnthropicCitationContentBlockLocationParam:
		return &v.StartBlockIndex
	case AnthropicCitationSearchResultLocationParam:
		return &v.StartBlockIndex
	}
	return nil
}

// GetCitationEndBlockIndex 返回 content_block_location 或 search_result_location
// 的 end_block_index。
func GetCitationEndBlockIndex(p AnthropicCitationsParam) *int64 {
	switch v := p.(type) {
	case AnthropicCitationContentBlockLocationParam:
		return &v.EndBlockIndex
	case AnthropicCitationSearchResultLocationParam:
		return &v.EndBlockIndex
	}
	return nil
}

// GetCitationSearchResultIndex 返回 search_result_location 的 search_result_index。
func GetCitationSearchResultIndex(p AnthropicCitationsParam) *int64 {
	if v, ok := p.(AnthropicCitationSearchResultLocationParam); ok {
		return &v.SearchResultIndex
	}
	return nil
}

// GetCitationEncryptedIndex 返回 web_search_result_location 的 encrypted_index。
func GetCitationEncryptedIndex(p AnthropicCitationsParam) string {
	if v, ok := p.(AnthropicCitationWebSearchResultLocationParam); ok {
		return v.EncryptedIndex
	}
	return ""
}

// GetCitationTitle 返回 search_result_location 或 web_search_result_location
// 的 title 字段（均为 *string）。
func GetCitationTitle(p AnthropicCitationsParam) *string {
	switch v := p.(type) {
	case AnthropicCitationSearchResultLocationParam:
		return v.Title
	case AnthropicCitationWebSearchResultLocationParam:
		return v.Title
	}
	return nil
}

// GetCitationURL 返回 web_search_result_location 的 url 字段。
func GetCitationURL(p AnthropicCitationsParam) string {
	if v, ok := p.(AnthropicCitationWebSearchResultLocationParam); ok {
		if v.URL == nil {
			return ""
		}
		return *v.URL
	}
	return ""
}

// GetCitationSource 返回 search_result_location 的 source 字段。
func GetCitationSource(p AnthropicCitationsParam) string {
	if v, ok := p.(AnthropicCitationSearchResultLocationParam); ok {
		if v.Source == nil {
			return ""
		}
		return *v.Source
	}
	return ""
}

// GetCitationType 返回 union 的 type 字段字符串（5 实现聚合）。
func GetCitationType(p AnthropicCitationsParam) string {
	switch v := p.(type) {
	case AnthropicCitationCharLocationParam:
		return string(v.Type)
	case AnthropicCitationPageLocationParam:
		return string(v.Type)
	case AnthropicCitationContentBlockLocationParam:
		return string(v.Type)
	case AnthropicCitationSearchResultLocationParam:
		return string(v.Type)
	case AnthropicCitationWebSearchResultLocationParam:
		return string(v.Type)
	}
	return ""
}
