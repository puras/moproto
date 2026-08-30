package gemini

// GeminiSchema OpenAPI 3.0 子集（Gemini responseSchema + FunctionDeclaration.Parameters 用）。
//
// 字段集对齐官方 SDK schema.Schema struct（types.go:2339），但简化为 nil-安全。
// 任意字段省略时由 omitempty 跳过；递归结构支持嵌套 object。
type GeminiSchema struct {
	Type        string                   `json:"type,omitempty"`
	Format      string                   `json:"format,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Enum        []string                 `json:"enum,omitempty"`
	Properties  map[string]*GeminiSchema `json:"properties,omitempty"`
	Required    []string                 `json:"required,omitempty"`
	Items       *GeminiSchema            `json:"items,omitempty"`
	MinItems    *int32                   `json:"minItems,omitempty"`
	MaxItems    *int32                   `json:"maxItems,omitempty"`
	Minimum     *float64                 `json:"minimum,omitempty"`
	Maximum     *float64                 `json:"maximum,omitempty"`
	MinLength   *int32                   `json:"minLength,omitempty"`
	MaxLength   *int32                   `json:"maxLength,omitempty"`
	Pattern     string                   `json:"pattern,omitempty"`
	Nullable    bool                     `json:"nullable,omitempty"`
	AnyOf       []*GeminiSchema          `json:"anyOf,omitempty"`
}
