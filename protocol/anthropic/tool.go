package anthropic

import "encoding/json"

// ============== AnthropicTool union ==============

// AnthropicCustomToolType 自定义工具 type 字段类型。
//
// 取值："custom"（默认，可省略）。
type AnthropicCustomToolType string

const AnthropicCustomToolTypeCustom AnthropicCustomToolType = "custom"

// AnthropicWebSearchToolType web_search 工具 type 字段类型。
//
// 当前支持版本："web_search_20260209"。
type AnthropicWebSearchToolType string

const AnthropicWebSearchToolType20260209 AnthropicWebSearchToolType = "web_search_20260209"

// AnthropicWebFetchToolType web_fetch 工具 type 字段类型。
//
// 当前支持版本："web_fetch_20260209"。
type AnthropicWebFetchToolType string

const AnthropicWebFetchToolType20260209 AnthropicWebFetchToolType = "web_fetch_20260209"

// AnthropicCodeExecutionToolType code_execution 工具 type 字段类型。
//
// 当前支持版本："code_execution_20260120"。
type AnthropicCodeExecutionToolType string

const AnthropicCodeExecutionToolType20260120 AnthropicCodeExecutionToolType = "code_execution_20260120"

// AnthropicBashToolType bash 工具 type 字段类型。
//
// 当前支持版本："bash_20250124"。
type AnthropicBashToolType string

const AnthropicBashToolType20250124 AnthropicBashToolType = "bash_20250124"

// AnthropicTextEditorToolType text_editor 工具 type 字段类型（新版 / 4.6+ 默认）。
//
// 当前支持版本："text_editor_20250728"。
type AnthropicTextEditorToolType string

const AnthropicTextEditorToolType20250728 AnthropicTextEditorToolType = "text_editor_20250728"

// AnthropicTextEditorTool20250124Type text_editor 工具 type 字段类型（旧版 1）。
//
// 对齐官方 ToolTextEditor20250124Param（message.go:7578-7601）。
type AnthropicTextEditorTool20250124Type string

const AnthropicTextEditorToolType20250124 AnthropicTextEditorTool20250124Type = "text_editor_20250124"

// AnthropicTextEditorTool20250429Type text_editor 工具 type 字段类型（旧版 2）。
//
// 对齐官方 ToolTextEditor20250429Param（message.go:7612-7635）。
type AnthropicTextEditorTool20250429Type string

const AnthropicTextEditorToolType20250429 AnthropicTextEditorTool20250429Type = "text_editor_20250429"

// AnthropicMemoryToolType memory 工具 type 字段类型。
//
// 当前支持版本："memory_20250818"。
type AnthropicMemoryToolType string

const AnthropicMemoryToolType20250818 AnthropicMemoryToolType = "memory_20250818"

// ---- CodeExecutionTool 多版本号（官方 SDK L1109-1229 共 3 个旧版本） ----

// AnthropicCodeExecutionTool20250522Type code_execution_20250522 工具 type 字段类型。
type AnthropicCodeExecutionTool20250522Type string

const AnthropicCodeExecutionToolType20250522 AnthropicCodeExecutionTool20250522Type = "code_execution_20250522"

// AnthropicCodeExecutionTool20250825Type code_execution_20250825 工具 type 字段类型。
type AnthropicCodeExecutionTool20250825Type string

const AnthropicCodeExecutionToolType20250825 AnthropicCodeExecutionTool20250825Type = "code_execution_20250825"

// AnthropicCodeExecutionTool20260521Type code_execution_20260521 工具 type 字段类型。
type AnthropicCodeExecutionTool20260521Type string

const AnthropicCodeExecutionToolType20260521 AnthropicCodeExecutionTool20260521Type = "code_execution_20260521"

// ---- WebSearchTool 多版本号（官方 SDK 4 个版本） ----

// AnthropicWebSearchTool20250305Type web_search_20250305 工具 type 字段类型。
type AnthropicWebSearchTool20250305Type string

const AnthropicWebSearchToolType20250305 AnthropicWebSearchTool20250305Type = "web_search_20250305"

// AnthropicWebSearchTool20260318Type web_search_20260318 工具 type 字段类型（含 response_inclusion 字段）。
type AnthropicWebSearchTool20260318Type string

const AnthropicWebSearchToolType20260318 AnthropicWebSearchTool20260318Type = "web_search_20260318"

// ---- WebFetchTool 多版本号（官方 SDK 4 个版本） ----

// AnthropicWebFetchTool20250910Type web_fetch_20250910 工具 type 字段类型。
type AnthropicWebFetchTool20250910Type string

const AnthropicWebFetchToolType20250910 AnthropicWebFetchTool20250910Type = "web_fetch_20250910"

// AnthropicWebFetchTool20260309Type web_fetch_20260309 工具 type 字段类型（新增 use_cache 字段）。
type AnthropicWebFetchTool20260309Type string

const AnthropicWebFetchToolType20260309 AnthropicWebFetchTool20260309Type = "web_fetch_20260309"

// AnthropicWebFetchTool20260318Type web_fetch_20260318 工具 type 字段类型。
type AnthropicWebFetchTool20260318Type string

const AnthropicWebFetchToolType20260318 AnthropicWebFetchTool20260318Type = "web_fetch_20260318"

// ---- ToolSearch 工具 ----

// AnthropicToolSearchToolBm25Type tool_search_tool_bm25_20251119 工具 type 字段类型。
type AnthropicToolSearchToolBm25Type string

const AnthropicToolSearchToolBm25Type20251119 AnthropicToolSearchToolBm25Type = "tool_search_tool_bm25_20251119"

// AnthropicToolSearchToolRegexType tool_search_tool_regex_20251119 工具 type 字段类型。
type AnthropicToolSearchToolRegexType string

const AnthropicToolSearchToolRegexType20251119 AnthropicToolSearchToolRegexType = "tool_search_tool_regex_20251119"

// AnthropicTool 工具定义 union（客户端传给服务端的工具声明）。
//
// 两种大类：
//   - 自定义工具：AnthropicCustomTool（需 input_schema）
//   - 服务端预定义工具：web_search / web_fetch / code_execution / bash / text_editor / memory
//
// 自定义工具的 InputSchema 是结构化 JSON Schema 对象（见 AnthropicToolInputSchema）；
// 预定义工具无此字段。
//
// 通过 isTool() 私有方法封闭实现集（外部包无法新增实现）。
// 调用方通过 type switch 或直接读取各实现的 Type 字段分发，详见 doc.go。
type AnthropicTool interface {
	isTool()
}

// AnthropicToolInputSchema JSON Schema for tool input.
//
// 对齐官方 anthropic-sdk-go ToolInputSchemaParam（message.go:6746-6761）：
//   - Type 固定为 "object"
//   - Properties 是嵌套 JSON Schema 节点（递归 any）；调用方按 map[string]any 类型断言
//   - Required 列出必填字段名
//   - ExtraFields 捕获未识别的 JSON Schema 元数据（如 $schema / additionalProperties / $ref 等）
type AnthropicToolInputSchema struct {
	Type        string         `json:"type"` // 固定 "object"
	Properties  any            `json:"properties,omitempty"`
	Required    []string       `json:"required,omitempty"`
	ExtraFields map[string]any `json:"-"`
}

// MarshalJSON 先序列化标准字段，再合并 ExtraFields。
//
// 对齐官方 param.MarshalWithExtras 的行为。
func (r AnthropicToolInputSchema) MarshalJSON() ([]byte, error) {
	// shadow 通过本方法内的 type alias 隐藏 ExtraFields。
	type shadow AnthropicToolInputSchema
	baseBytes, err := json.Marshal((*shadow)(&r))
	if err != nil {
		return nil, err
	}
	if len(r.ExtraFields) == 0 {
		return baseBytes, nil
	}
	var merged map[string]any
	if err := json.Unmarshal(baseBytes, &merged); err != nil {
		return nil, err
	}
	for k, v := range r.ExtraFields {
		merged[k] = v
	}
	return json.Marshal(merged)
}

// UnmarshalJSON 先解析标准字段，再把未识别的顶层 key 收集到 ExtraFields。
func (r *AnthropicToolInputSchema) UnmarshalJSON(data []byte) error {
	// 第一步：用 shadow 解码标准字段
	type shadow AnthropicToolInputSchema
	var s shadow
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*r = AnthropicToolInputSchema(s)

	// 第二步：把 data unmarshal 成 map 收集未识别字段
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	known := map[string]struct{}{
		"type": {}, "properties": {}, "required": {},
	}
	for k, v := range raw {
		if _, ok := known[k]; ok {
			continue
		}
		if r.ExtraFields == nil {
			r.ExtraFields = map[string]any{}
		}
		var val any
		if err := json.Unmarshal(v, &val); err != nil {
			return err
		}
		r.ExtraFields[k] = val
	}
	return nil
}

// AnthropicCustomTool 自定义工具（用户定义的客户端工具）。
//
// 对齐官方 anthropic-sdk-go ToolParam（message.go:6693-6730）。
// 4 个扩展字段（EagerInputStreaming / DeferLoading / AllowedCallers / InputExamples）
// 全部为 omitempty，零值不写入 JSON，与现有 wire-format 完全向后兼容。
//
// AllowedCallers / InputExamples 与官方 SDK 字段类型完全一致：
//   - AllowedCallers: []string（官方即此类型；非 union）
//   - InputExamples:  []map[string]any（官方即此类型；允许 input 任意结构）
//
// 网关需校验 AllowedCallers 取值范围时可用 IsAllowedCallerValid() 助手函数。
type AnthropicCustomTool struct {
	Type         AnthropicCustomToolType         `json:"type,omitempty"` // 可省略（默认 "custom"）
	Name         string                          `json:"name"`
	Description  string                          `json:"description,omitempty"`
	InputSchema  AnthropicToolInputSchema        `json:"input_schema"`
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
	Strict       bool                            `json:"strict,omitempty"`

	// ---- 4 个扩展字段（官方 ToolParam 4.8+） ----

	// EagerInputStreaming 启用工具输入的提前流式。
	// true：增量流式生成工具输入；false：禁用；nil：使用默认值。
	EagerInputStreaming *bool `json:"eager_input_streaming,omitempty"`
	// DeferLoading 延迟加载工具定义。
	// true：工具不放入初始 system prompt，仅在被 tool_reference 引用时加载。
	DeferLoading *bool `json:"defer_loading,omitempty"`
	// AllowedCallers 控制谁能调用此工具。
	// 取值范围（与官方 AnthropicServerToolUseBlockParamCallerUnion 的所有 type 字符串一致）：
	//   - "direct"
	//   - "code_execution_20250825"
	//   - "code_execution_20260120"
	//   - "code_execution_20260521"
	// 类型与官方一致为 []string（非 union interface），因为跨网关注入场景下 union 增加的类型
	// 安全收益有限，反而增加调用方序列化负担。
	AllowedCallers []string `json:"allowed_callers,omitempty"`
	// InputExamples 工具调用示例（few-shot 增强）。
	// 类型与官方一致为 []map[string]any，允许 input 任意 JSON 结构。
	InputExamples []map[string]any `json:"input_examples,omitempty"`
}

// IsAllowedCallerValid 校验单个 AllowedCaller 取值是否在官方枚举范围内。
//
// 对应官方 AnthropicServerToolUseBlockParamCallerUnion 的所有 type 字符串
// （message.go:5526-5528）。
func IsAllowedCallerValid(caller string) bool {
	switch caller {
	case "direct",
		"code_execution_20250825",
		"code_execution_20260120",
		"code_execution_20260521":
		return true
	}
	return false
}

// ValidateAllowedCallers 返回首个不合法 caller 与 false；如全合法返回 "", true。
//
// 网关可在 audit / 反垃圾层调用；不影响 wire-format。
func (t *AnthropicCustomTool) ValidateAllowedCallers() (string, bool) {
	for _, c := range t.AllowedCallers {
		if !IsAllowedCallerValid(c) {
			return c, false
		}
	}
	return "", true
}

func (AnthropicCustomTool) isTool() {}

// AnthropicCustomToolOption 函数式选项，用于构造 AnthropicCustomTool 的扩展字段。
type AnthropicCustomToolOption func(*AnthropicCustomTool)

// WithStrict 设置 strict 字段（保证 schema 验证）。
func WithStrict() AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.Strict = true }
}

// WithDescription 设置 description 字段。
func WithDescription(desc string) AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.Description = desc }
}

// WithCacheControl 设置缓存控制。
func WithCacheControl(cc *AnthropicCacheControlEphemeral) AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.CacheControl = cc }
}

// WithDeferLoading 设置 defer_loading 字段。
func WithDeferLoading(v bool) AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.DeferLoading = &v }
}

// WithEagerInputStreaming 设置 eager_input_streaming 字段。
func WithEagerInputStreaming(v bool) AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.EagerInputStreaming = &v }
}

// WithAllowedCallers 设置 allowed_callers 字段。
func WithAllowedCallers(callers ...string) AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.AllowedCallers = append(t.AllowedCallers, callers...) }
}

// WithInputExamples 设置 input_examples 字段。
func WithInputExamples(examples ...map[string]any) AnthropicCustomToolOption {
	return func(t *AnthropicCustomTool) { t.InputExamples = append(t.InputExamples, examples...) }
}

// NewAnthropicCustomTool 构造自定义工具（functional options 模式）。
//
// 必填参数：name、inputSchema。其它字段（description / cache_control / strict
// / 4 个扩展字段）通过 With* 选项设置。
//
// 示例：
//
//	t := anthropic.NewAnthropicCustomTool(
//	    "get_weather",
//	    anthropic.AnthropicToolInputSchema{Type: "object"},
//	    anthropic.WithDescription("获取天气信息"),
//	    anthropic.WithDeferLoading(true),
//	    anthropic.WithAllowedCallers("direct"),
//	)
func NewAnthropicCustomTool(name string, inputSchema AnthropicToolInputSchema, opts ...AnthropicCustomToolOption) AnthropicCustomTool {
	t := AnthropicCustomTool{
		Type:        AnthropicCustomToolTypeCustom,
		Name:        name,
		InputSchema: inputSchema,
	}
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

// ============== 预定义工具（type 字段固定，服务端识别） ==============

// AnthropicWebSearchTool web_search 工具（20260209 版）。
type AnthropicWebSearchTool struct {
	Type           AnthropicWebSearchToolType      `json:"type"` // 固定 "web_search_20260209"
	Name           string                          `json:"name"` // 通常 "web_search"
	MaxUses        int                             `json:"max_uses,omitempty"`
	AllowedDomains []string                        `json:"allowed_domains,omitempty"`
	BlockedDomains []string                        `json:"blocked_domains,omitempty"`
	UserLocation   *AnthropicUserLocation          `json:"user_location,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicWebSearchTool) isTool() {}

// AnthropicUserLocation 用户近似位置（web_search 工具上下文）。
//
// 对齐官方 anthropic-sdk-go UserLocationParam（message.go:8501-8515）。
// 当前由 AnthropicWebSearchTool 使用；保留为顶层共享类型以便未来 web_fetch 或
// 其他工具复用（官方 SDK WebFetchTool 当前也未使用此字段）。
//
// Type 固定为 "approximate"；其余字段全部可选（omitempty）。
type AnthropicUserLocation struct {
	Type     string `json:"type"` // 固定 "approximate"
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

// NewAnthropicUserLocation 构造用户近似位置。
//
// city / country / region / timezone 全部可选；调用方按需传 string 指针。
func NewAnthropicUserLocation(city, country, region, timezone string) *AnthropicUserLocation {
	return &AnthropicUserLocation{
		Type:     "approximate",
		City:     city,
		Country:  country,
		Region:   region,
		Timezone: timezone,
	}
}

// AnthropicWebFetchTool web_fetch 工具（20260209 版）。
type AnthropicWebFetchTool struct {
	Type           AnthropicWebFetchToolType       `json:"type"` // 固定 "web_fetch_20260209"
	Name           string                          `json:"name"` // 通常 "web_fetch"
	MaxUses        int                             `json:"max_uses,omitempty"`
	AllowedDomains []string                        `json:"allowed_domains,omitempty"`
	BlockedDomains []string                        `json:"blocked_domains,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicWebFetchTool) isTool() {}

// AnthropicCodeExecutionTool code_execution 工具（20260120 版）。
type AnthropicCodeExecutionTool struct {
	Type         AnthropicCodeExecutionToolType  `json:"type"` // 固定 "code_execution_20260120"
	Name         string                          `json:"name"` // 通常 "code_execution"
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicCodeExecutionTool) isTool() {}

// AnthropicBashTool bash 工具（20250124 版）。
type AnthropicBashTool struct {
	Type         AnthropicBashToolType           `json:"type"` // 固定 "bash_20250124"
	Name         string                          `json:"name"` // 通常 "bash"
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicBashTool) isTool() {}

// AnthropicTextEditorTool 文本编辑器工具（20250728 版，str_replace_based_edit_tool）。
type AnthropicTextEditorTool struct {
	Type         AnthropicTextEditorToolType     `json:"type"` // 固定 "text_editor_20250728"
	Name         string                          `json:"name"` // 通常 "str_replace_based_edit_tool"
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicTextEditorTool) isTool() {}

// AnthropicTextEditorTool20250124 文本编辑器工具（20250124 旧版，str_replace_editor）。
//
// 对齐官方 ToolTextEditor20250124Param（message.go:7578-7601）。
// 字段布局与官方一致：Type / Name / DeferLoading / Strict / AllowedCallers / CacheControl / InputExamples。
type AnthropicTextEditorTool20250124 struct {
	Type           AnthropicTextEditorTool20250124Type `json:"type"` // 固定 "text_editor_20250124"
	Name           string                              `json:"name"` // 固定 "str_replace_editor"
	DeferLoading   *bool                               `json:"defer_loading,omitempty"`
	Strict         *bool                               `json:"strict,omitempty"`
	AllowedCallers []string                            `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral     `json:"cache_control,omitempty"`
	InputExamples  []map[string]any                    `json:"input_examples,omitempty"`
}

func (AnthropicTextEditorTool20250124) isTool() {}

// AnthropicTextEditorTool20250429 文本编辑器工具（20250429 旧版，str_replace_based_edit_tool）。
//
// 对齐官方 ToolTextEditor20250429Param（message.go:7612-7635）。
type AnthropicTextEditorTool20250429 struct {
	Type           AnthropicTextEditorTool20250429Type `json:"type"` // 固定 "text_editor_20250429"
	Name           string                              `json:"name"` // 固定 "str_replace_based_edit_tool"
	DeferLoading   *bool                               `json:"defer_loading,omitempty"`
	Strict         *bool                               `json:"strict,omitempty"`
	AllowedCallers []string                            `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral     `json:"cache_control,omitempty"`
	InputExamples  []map[string]any                    `json:"input_examples,omitempty"`
}

func (AnthropicTextEditorTool20250429) isTool() {}

// AnthropicMemoryTool memory 工具（20250818 版）。
type AnthropicMemoryTool struct {
	Type         AnthropicMemoryToolType         `json:"type"` // 固定 "memory_20250818"
	Name         string                          `json:"name"` // 通常 "memory"
	CacheControl *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicMemoryTool) isTool() {}

// ============== 预定义工具：CodeExecutionTool 多版本（旧 3 版） ==============
//
// 3 个旧版 code_execution 工具结构相同（DeferLoading / Strict / AllowedCallers /
// CacheControl / Name / Type），区别仅在 Type 字符串与 `omitzero` 的 deprecation 语义。
// 对齐官方 anthropic-sdk-go CodeExecutionTool20250522/20250825/20260521Param
//（message.go:1109-1229）。

// AnthropicCodeExecutionTool20250522 code_execution 工具（20250522 旧版）。
type AnthropicCodeExecutionTool20250522 struct {
	Type           AnthropicCodeExecutionTool20250522Type `json:"type"` // 固定 "code_execution_20250522"
	Name           string                                 `json:"name"` // 通常 "code_execution"
	Strict         bool                                   `json:"strict,omitempty"`
	DeferLoading   bool                                   `json:"defer_loading,omitempty"`
	AllowedCallers []string                               `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral        `json:"cache_control,omitempty"`
}

func (AnthropicCodeExecutionTool20250522) isTool() {}

// AnthropicCodeExecutionTool20250825 code_execution 工具（20250825 旧版）。
type AnthropicCodeExecutionTool20250825 struct {
	Type           AnthropicCodeExecutionTool20250825Type `json:"type"` // 固定 "code_execution_20250825"
	Name           string                                 `json:"name"`
	Strict         bool                                   `json:"strict,omitempty"`
	DeferLoading   bool                                   `json:"defer_loading,omitempty"`
	AllowedCallers []string                               `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral        `json:"cache_control,omitempty"`
}

func (AnthropicCodeExecutionTool20250825) isTool() {}

// AnthropicCodeExecutionTool20260521 code_execution 工具（20260521 旧版）。
type AnthropicCodeExecutionTool20260521 struct {
	Type           AnthropicCodeExecutionTool20260521Type `json:"type"` // 固定 "code_execution_20260521"
	Name           string                                 `json:"name"`
	Strict         bool                                   `json:"strict,omitempty"`
	DeferLoading   bool                                   `json:"defer_loading,omitempty"`
	AllowedCallers []string                               `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral        `json:"cache_control,omitempty"`
}

func (AnthropicCodeExecutionTool20260521) isTool() {}

// ============== 预定义工具：WebSearchTool 多版本（其他 2 版） ==============
//
// 对齐官方 anthropic-sdk-go WebSearchTool20250305/20260318Param（message.go:9161-9263）。
// 20260318 新增 ResponseInclusion 字段（"full" | "excluded"），控制 code_execution
// 消费后的结果回传方式。

// AnthropicWebSearchTool20250305 web_search 工具（20250305 旧版）。
type AnthropicWebSearchTool20250305 struct {
	Type           AnthropicWebSearchTool20250305Type `json:"type"` // 固定 "web_search_20250305"
	Name           string                             `json:"name"` // 通常 "web_search"
	MaxUses        int                                `json:"max_uses,omitempty"`
	AllowedDomains []string                           `json:"allowed_domains,omitempty"`
	BlockedDomains []string                           `json:"blocked_domains,omitempty"`
	Strict         bool                               `json:"strict,omitempty"`
	DeferLoading   bool                               `json:"defer_loading,omitempty"`
	AllowedCallers []string                           `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral    `json:"cache_control,omitempty"`
	UserLocation   *AnthropicUserLocation             `json:"user_location,omitempty"`
}

func (AnthropicWebSearchTool20250305) isTool() {}

// AnthropicWebSearchToolResponseInclusion web_search 工具结果回传策略。
type AnthropicWebSearchToolResponseInclusion string

const (
	// AnthropicWebSearchToolResponseInclusionFull 完整回传结果块（默认）。
	AnthropicWebSearchToolResponseInclusionFull AnthropicWebSearchToolResponseInclusion = "full"
	// AnthropicWebSearchToolResponseInclusionExcluded 排除嵌套的 server_tool_use 与 result 块。
	AnthropicWebSearchToolResponseInclusionExcluded AnthropicWebSearchToolResponseInclusion = "excluded"
)

// AnthropicWebSearchTool20260318 web_search 工具（20260318 版，新增 response_inclusion）。
type AnthropicWebSearchTool20260318 struct {
	Type              AnthropicWebSearchTool20260318Type      `json:"type"` // 固定 "web_search_20260318"
	Name              string                                  `json:"name"`
	MaxUses           int                                     `json:"max_uses,omitempty"`
	AllowedDomains    []string                                `json:"allowed_domains,omitempty"`
	BlockedDomains    []string                                `json:"blocked_domains,omitempty"`
	Strict            bool                                    `json:"strict,omitempty"`
	DeferLoading      bool                                    `json:"defer_loading,omitempty"`
	AllowedCallers    []string                                `json:"allowed_callers,omitempty"`
	CacheControl      *AnthropicCacheControlEphemeral         `json:"cache_control,omitempty"`
	UserLocation      *AnthropicUserLocation                  `json:"user_location,omitempty"`
	ResponseInclusion AnthropicWebSearchToolResponseInclusion `json:"response_inclusion,omitempty"`
}

func (AnthropicWebSearchTool20260318) isTool() {}

// ============== 预定义工具：WebFetchTool 多版本（其他 3 版） ==============
//
// 对齐官方 anthropic-sdk-go WebFetchTool20250910/20260309/20260318Param
//（message.go:8570-8740）。20260309 起新增 UseCache 字段；20250910 还含 Citations 字段。

// AnthropicWebFetchTool20250910 web_fetch 工具（20250910 版）。
type AnthropicWebFetchTool20250910 struct {
	Type             AnthropicWebFetchTool20250910Type `json:"type"` // 固定 "web_fetch_20250910"
	Name             string                            `json:"name"`
	MaxContentTokens int64                             `json:"max_content_tokens,omitempty"`
	MaxUses          int                               `json:"max_uses,omitempty"`
	AllowedDomains   []string                          `json:"allowed_domains,omitempty"`
	BlockedDomains   []string                          `json:"blocked_domains,omitempty"`
	Strict           bool                              `json:"strict,omitempty"`
	DeferLoading     bool                              `json:"defer_loading,omitempty"`
	AllowedCallers   []string                          `json:"allowed_callers,omitempty"`
	CacheControl     *AnthropicCacheControlEphemeral   `json:"cache_control,omitempty"`
	Citations        *AnthropicCitationsConfig         `json:"citations,omitempty"`
}

func (AnthropicWebFetchTool20250910) isTool() {}

// AnthropicWebFetchTool20260309 web_fetch 工具（20260309 版，新增 use_cache 字段）。
type AnthropicWebFetchTool20260309 struct {
	Type             AnthropicWebFetchTool20260309Type `json:"type"` // 固定 "web_fetch_20260309"
	Name             string                            `json:"name"`
	MaxContentTokens int64                             `json:"max_content_tokens,omitempty"`
	MaxUses          int                               `json:"max_uses,omitempty"`
	UseCache         bool                              `json:"use_cache,omitempty"`
	AllowedDomains   []string                          `json:"allowed_domains,omitempty"`
	BlockedDomains   []string                          `json:"blocked_domains,omitempty"`
	Strict           bool                              `json:"strict,omitempty"`
	DeferLoading     bool                              `json:"defer_loading,omitempty"`
	AllowedCallers   []string                          `json:"allowed_callers,omitempty"`
	CacheControl     *AnthropicCacheControlEphemeral   `json:"cache_control,omitempty"`
	Citations        *AnthropicCitationsConfig         `json:"citations,omitempty"`
}

func (AnthropicWebFetchTool20260309) isTool() {}

// AnthropicWebFetchTool20260318 web_fetch 工具（20260318 版，简化字段集）。
type AnthropicWebFetchTool20260318 struct {
	Type             AnthropicWebFetchTool20260318Type `json:"type"` // 固定 "web_fetch_20260318"
	Name             string                            `json:"name"`
	MaxContentTokens int64                             `json:"max_content_tokens,omitempty"`
	MaxUses          int                               `json:"max_uses,omitempty"`
	UseCache         bool                              `json:"use_cache,omitempty"`
	AllowedDomains   []string                          `json:"allowed_domains,omitempty"`
	BlockedDomains   []string                          `json:"blocked_domains,omitempty"`
	Strict           bool                              `json:"strict,omitempty"`
	DeferLoading     bool                              `json:"defer_loading,omitempty"`
	AllowedCallers   []string                          `json:"allowed_callers,omitempty"`
	CacheControl     *AnthropicCacheControlEphemeral   `json:"cache_control,omitempty"`
}

func (AnthropicWebFetchTool20260318) isTool() {}

// ============== 预定义工具：ToolSearch（BM25 + Regex） ==============
//
// 对齐官方 anthropic-sdk-go ToolSearchToolBm25_20251119/Regex20251119Param
//（message.go:7277-7337）。结构相同：DeferLoading/Strict/AllowedCallers/
// CacheControl/Name/Type，Type 字符串区分 bm25 与 regex。

// AnthropicToolSearchToolBm25 工具搜索工具（BM25 算法，20251119 版）。
type AnthropicToolSearchToolBm25 struct {
	Type           AnthropicToolSearchToolBm25Type `json:"type"` // 固定 "tool_search_tool_bm25_20251119"
	Name           string                          `json:"name"` // 通常 "tool_search_tool_bm25"
	Strict         bool                            `json:"strict,omitempty"`
	DeferLoading   bool                            `json:"defer_loading,omitempty"`
	AllowedCallers []string                        `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral `json:"cache_control,omitempty"`
}

func (AnthropicToolSearchToolBm25) isTool() {}

// AnthropicToolSearchToolRegex 工具搜索工具（正则表达式算法，20251119 版）。
type AnthropicToolSearchToolRegex struct {
	Type           AnthropicToolSearchToolRegexType `json:"type"` // 固定 "tool_search_tool_regex_20251119"
	Name           string                           `json:"name"` // 通常 "tool_search_tool_regex"
	Strict         bool                             `json:"strict,omitempty"`
	DeferLoading   bool                             `json:"defer_loading,omitempty"`
	AllowedCallers []string                         `json:"allowed_callers,omitempty"`
	CacheControl   *AnthropicCacheControlEphemeral  `json:"cache_control,omitempty"`
}

func (AnthropicToolSearchToolRegex) isTool() {}

// ============== AnthropicToolChoice union ==============

// AnthropicToolChoiceAutoType tool_choice type 字段类型。
type AnthropicToolChoiceAutoType string

const AnthropicToolChoiceAutoTypeAuto AnthropicToolChoiceAutoType = "auto"

// AnthropicToolChoiceAnyType tool_choice type 字段类型。
type AnthropicToolChoiceAnyType string

const AnthropicToolChoiceAnyTypeAny AnthropicToolChoiceAnyType = "any"

// AnthropicToolChoiceToolType tool_choice type 字段类型。
type AnthropicToolChoiceToolType string

const AnthropicToolChoiceToolTypeTool AnthropicToolChoiceToolType = "tool"

// AnthropicToolChoiceNoneType tool_choice type 字段类型。
type AnthropicToolChoiceNoneType string

const AnthropicToolChoiceNoneTypeNone AnthropicToolChoiceNoneType = "none"

// AnthropicToolChoice 工具调用策略 union。
//
// 通过 isToolChoice() 私有方法封闭实现集（外部包无法新增实现）。
// 调用方通过 type switch 或直接读取各实现的 Type 字段分发，详见 doc.go。
type AnthropicToolChoice interface {
	isToolChoice()
}

// AnthropicToolChoiceAuto 模型自行决定（默认）。
type AnthropicToolChoiceAuto struct {
	Type                   AnthropicToolChoiceAutoType `json:"type"` // 固定 "auto"
	DisableParallelToolUse bool                        `json:"disable_parallel_tool_use,omitempty"`
}

func (AnthropicToolChoiceAuto) isToolChoice() {}

// AnthropicToolChoiceAny 必须调用某个工具。
type AnthropicToolChoiceAny struct {
	Type                   AnthropicToolChoiceAnyType `json:"type"` // 固定 "any"
	DisableParallelToolUse bool                       `json:"disable_parallel_tool_use,omitempty"`
}

func (AnthropicToolChoiceAny) isToolChoice() {}

// AnthropicToolChoiceTool 强制调用指定工具。
type AnthropicToolChoiceTool struct {
	Type                   AnthropicToolChoiceToolType `json:"type"` // 固定 "tool"
	Name                   string                      `json:"name"`
	DisableParallelToolUse bool                        `json:"disable_parallel_tool_use,omitempty"`
}

func (AnthropicToolChoiceTool) isToolChoice() {}

// AnthropicToolChoiceNone 禁用工具。
type AnthropicToolChoiceNone struct {
	Type AnthropicToolChoiceNoneType `json:"type"` // 固定 "none"
}

func (AnthropicToolChoiceNone) isToolChoice() {}
