package openai

import "encoding/json"

// OpenAIResponsesToolUnion Responses Tool union 接口。
type OpenAIResponsesToolUnion interface {
	isResponsesTool()
}

// OpenAIResponsesToolUnionFlat Responses Tool union 扁平化实现（16 变体）。
type OpenAIResponsesToolUnionFlat struct {
	// Type 区分变体：
	// "function" / "web_search" / "code_interpreter" / "file_search" /
	// "image_generation" / "mcp" / "computer_use_preview" / "apply_patch" /
	// "code_execution" / "hosted" / "local_shell" / "shell" / "custom" /
	// "tool_search" / "web_search_preview" / "namespace"
	Type string `json:"type,omitempty"`

	// === function ===
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Strict      *bool           `json:"strict,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`

	// === web_search / web_search_preview ===
	Filters           *OpenAIResponsesWebSearchFilters      `json:"filters,omitempty"`
	SearchContextSize string                                `json:"search_context_size,omitempty"`
	UserLocation      *OpenAIResponsesWebSearchUserLocation `json:"user_location,omitempty"`

	// === code_interpreter ===
	Container OpenAIResponsesCodeInterpreterContainer `json:"container,omitempty"`

	// === file_search ===
	VectorStoreIDs []string                          `json:"vector_store_ids,omitempty"`
	MaxNumResults  *int64                            `json:"max_num_results,omitempty"`
	RankingOptions *OpenAIResponsesFileSearchRanking `json:"ranking_options,omitempty"`
	Filters2       json.RawMessage                   `json:"filters,omitempty"`

	// === image_generation ===
	Background    string `json:"background,omitempty"`
	InputFidelity string `json:"input_fidelity,omitempty"`
	Model         string `json:"model,omitempty"`
	Moderation    string `json:"moderation,omitempty"`
	Quality       string `json:"quality,omitempty"`
	OutputFormat  string `json:"output_format,omitempty"`
	Size          string `json:"size,omitempty"`

	// === mcp ===
	ServerLabel     string            `json:"server_label,omitempty"`
	ServerURL       string            `json:"server_url,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	AllowedTools    []string          `json:"allowed_tools,omitempty"`
	RequireApproval string            `json:"require_approval,omitempty"`

	// === computer_use_preview ===
	DisplayWidth  int64  `json:"display_width,omitempty"`
	DisplayHeight int64  `json:"display_height,omitempty"`
	Environment   string `json:"environment,omitempty"`

	// === local_shell / shell ===
	Environment2 *OpenAIResponsesShellEnvironment `json:"environment,omitempty"`

	// === custom ===
	Format *OpenAIResponsesCustomToolFormatGrammar `json:"format,omitempty"`

	// === tool_search ===
	Execution   string                            `json:"execution,omitempty"`
	Filters3    []string                          `json:"filters,omitempty"`
	MaxResults  *int64                            `json:"max_num_results,omitempty"`
	RankingOpts *OpenAIResponsesToolSearchRanking `json:"ranking_options,omitempty"`

	// === namespace ===
	NamespaceName string `json:"name,omitempty"`
}

func (OpenAIResponsesToolUnionFlat) isResponsesTool() {}

func (u OpenAIResponsesToolUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "function":
		return json.Marshal(struct {
			Type        string          `json:"type"`
			Name        string          `json:"name"`
			Description string          `json:"description,omitempty"`
			Strict      *bool           `json:"strict,omitempty"`
			Parameters  json.RawMessage `json:"parameters,omitempty"`
		}{
			Type:        "function",
			Name:        u.Name,
			Description: u.Description,
			Strict:      u.Strict,
			Parameters:  u.Parameters,
		})
	case "web_search":
		return json.Marshal(struct {
			Type              string                                `json:"type"`
			Filters           *OpenAIResponsesWebSearchFilters      `json:"filters,omitempty"`
			SearchContextSize string                                `json:"search_context_size,omitempty"`
			UserLocation      *OpenAIResponsesWebSearchUserLocation `json:"user_location,omitempty"`
		}{
			Type:              "web_search",
			Filters:           u.Filters,
			SearchContextSize: u.SearchContextSize,
			UserLocation:      u.UserLocation,
		})
	case "code_interpreter":
		return json.Marshal(struct {
			Type      string                                  `json:"type"`
			Container OpenAIResponsesCodeInterpreterContainer `json:"container"`
		}{
			Type:      "code_interpreter",
			Container: u.Container,
		})
	case "file_search":
		return json.Marshal(struct {
			Type           string                            `json:"type"`
			VectorStoreIDs []string                          `json:"vector_store_ids"`
			MaxNumResults  *int64                            `json:"max_num_results,omitempty"`
			RankingOptions *OpenAIResponsesFileSearchRanking `json:"ranking_options,omitempty"`
			Filters        json.RawMessage                   `json:"filters,omitempty"`
		}{
			Type:           "file_search",
			VectorStoreIDs: u.VectorStoreIDs,
			MaxNumResults:  u.MaxNumResults,
			RankingOptions: u.RankingOptions,
			Filters:        u.Filters2,
		})
	case "image_generation":
		return json.Marshal(struct {
			Type          string `json:"type"`
			Background    string `json:"background,omitempty"`
			InputFidelity string `json:"input_fidelity,omitempty"`
			Model         string `json:"model,omitempty"`
			Moderation    string `json:"moderation,omitempty"`
			Quality       string `json:"quality,omitempty"`
			OutputFormat  string `json:"output_format,omitempty"`
			Size          string `json:"size,omitempty"`
		}{
			Type:          "image_generation",
			Background:    u.Background,
			InputFidelity: u.InputFidelity,
			Model:         u.Model,
			Moderation:    u.Moderation,
			Quality:       u.Quality,
			OutputFormat:  u.OutputFormat,
			Size:          u.Size,
		})
	case "mcp":
		return json.Marshal(struct {
			Type            string            `json:"type"`
			ServerLabel     string            `json:"server_label"`
			ServerURL       string            `json:"server_url,omitempty"`
			Headers         map[string]string `json:"headers,omitempty"`
			AllowedTools    []string          `json:"allowed_tools,omitempty"`
			RequireApproval string            `json:"require_approval,omitempty"`
		}{
			Type:            "mcp",
			ServerLabel:     u.ServerLabel,
			ServerURL:       u.ServerURL,
			Headers:         u.Headers,
			AllowedTools:    u.AllowedTools,
			RequireApproval: u.RequireApproval,
		})
	case "computer_use_preview":
		return json.Marshal(struct {
			Type          string `json:"type"`
			DisplayWidth  int64  `json:"display_width"`
			DisplayHeight int64  `json:"display_height"`
			Environment   string `json:"environment"`
		}{
			Type:          "computer_use_preview",
			DisplayWidth:  u.DisplayWidth,
			DisplayHeight: u.DisplayHeight,
			Environment:   u.Environment,
		})
	case "apply_patch":
		return json.Marshal(struct {
			Type string `json:"type"`
		}{
			Type: "apply_patch",
		})
	case "code_execution":
		return json.Marshal(struct {
			Type string `json:"type"`
		}{
			Type: "code_execution",
		})
	case "hosted":
		return json.Marshal(struct {
			Type string `json:"type"`
		}{
			Type: u.Type,
		})
	case "local_shell":
		return json.Marshal(struct {
			Type string `json:"type"`
		}{
			Type: "local_shell",
		})
	case "shell":
		return json.Marshal(struct {
			Type        string                           `json:"type"`
			Environment *OpenAIResponsesShellEnvironment `json:"environment,omitempty"`
		}{
			Type:        "shell",
			Environment: u.Environment2,
		})
	case "custom":
		return json.Marshal(struct {
			Type        string                                  `json:"type"`
			Name        string                                  `json:"name"`
			Description string                                  `json:"description,omitempty"`
			Format      *OpenAIResponsesCustomToolFormatGrammar `json:"format,omitempty"`
		}{
			Type:        "custom",
			Name:        u.Name,
			Description: u.Description,
			Format:      u.Format,
		})
	case "tool_search":
		return json.Marshal(struct {
			Type           string                            `json:"type"`
			Execution      string                            `json:"execution,omitempty"`
			Filters        []string                          `json:"filters,omitempty"`
			MaxNumResults  *int64                            `json:"max_num_results,omitempty"`
			RankingOptions *OpenAIResponsesToolSearchRanking `json:"ranking_options,omitempty"`
		}{
			Type:           "tool_search",
			Execution:      u.Execution,
			Filters:        u.Filters3,
			MaxNumResults:  u.MaxResults,
			RankingOptions: u.RankingOpts,
		})
	case "web_search_preview":
		return json.Marshal(struct {
			Type              string                                `json:"type"`
			SearchContextSize string                                `json:"search_context_size,omitempty"`
			UserLocation      *OpenAIResponsesWebSearchUserLocation `json:"user_location,omitempty"`
		}{
			Type:              "web_search_preview",
			SearchContextSize: u.SearchContextSize,
			UserLocation:      u.UserLocation,
		})
	case "namespace":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
		}{
			Type:        "namespace",
			Name:        u.NamespaceName,
			Description: u.Description,
		})
	}
	return nil, nil
}

func (u *OpenAIResponsesToolUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var tmp struct {
		Type string `json:"type"`
		// function
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Strict      *bool           `json:"strict"`
		Parameters  json.RawMessage `json:"parameters"`
		// web_search / web_search_preview
		Filters           *OpenAIResponsesWebSearchFilters      `json:"filters"`
		SearchContextSize string                                `json:"search_context_size"`
		UserLocation      *OpenAIResponsesWebSearchUserLocation `json:"user_location"`
		// code_interpreter
		Container OpenAIResponsesCodeInterpreterContainer `json:"container"`
		// file_search
		VectorStoreIDs []string                          `json:"vector_store_ids"`
		MaxNumResults  *int64                            `json:"max_num_results"`
		RankingOptions *OpenAIResponsesFileSearchRanking `json:"ranking_options"`
		Filters2       json.RawMessage                   `json:"filters"`
		// image_generation
		Background    string `json:"background"`
		InputFidelity string `json:"input_fidelity"`
		Model         string `json:"model"`
		Moderation    string `json:"moderation"`
		Quality       string `json:"quality"`
		OutputFormat  string `json:"output_format"`
		Size          string `json:"size"`
		// mcp
		ServerLabel     string            `json:"server_label"`
		ServerURL       string            `json:"server_url"`
		Headers         map[string]string `json:"headers"`
		AllowedTools    []string          `json:"allowed_tools"`
		RequireApproval string            `json:"require_approval"`
		// computer_use_preview
		DisplayWidth  int64  `json:"display_width"`
		DisplayHeight int64  `json:"display_height"`
		Environment   string `json:"environment"`
		// local_shell / shell
		Environment2 *OpenAIResponsesShellEnvironment `json:"environment"`
		// custom
		Format *OpenAIResponsesCustomToolFormatGrammar `json:"format"`
		// tool_search
		Execution   string                            `json:"execution"`
		Filters3    []string                          `json:"filters"`
		MaxResults  *int64                            `json:"max_num_results"`
		RankingOpts *OpenAIResponsesToolSearchRanking `json:"ranking_options"`
		// namespace
		NamespaceName string `json:"name"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "function":
		u.Name = tmp.Name
		u.Description = tmp.Description
		u.Strict = tmp.Strict
		u.Parameters = tmp.Parameters
	case "web_search":
		u.Filters = tmp.Filters
		u.SearchContextSize = tmp.SearchContextSize
		u.UserLocation = tmp.UserLocation
	case "code_interpreter":
		u.Container = tmp.Container
	case "file_search":
		u.VectorStoreIDs = tmp.VectorStoreIDs
		u.MaxNumResults = tmp.MaxNumResults
		u.RankingOptions = tmp.RankingOptions
		u.Filters2 = tmp.Filters2
	case "image_generation":
		u.Background = tmp.Background
		u.InputFidelity = tmp.InputFidelity
		u.Model = tmp.Model
		u.Moderation = tmp.Moderation
		u.Quality = tmp.Quality
		u.OutputFormat = tmp.OutputFormat
		u.Size = tmp.Size
	case "mcp":
		u.ServerLabel = tmp.ServerLabel
		u.ServerURL = tmp.ServerURL
		u.Headers = tmp.Headers
		u.AllowedTools = tmp.AllowedTools
		u.RequireApproval = tmp.RequireApproval
	case "computer_use_preview":
		u.DisplayWidth = tmp.DisplayWidth
		u.DisplayHeight = tmp.DisplayHeight
		u.Environment = tmp.Environment
	case "apply_patch":
		// no additional fields
	case "code_execution":
		// no additional fields
	case "hosted":
		// no additional fields
	case "local_shell":
		// no additional fields
	case "shell":
		u.Environment2 = tmp.Environment2
	case "custom":
		u.Name = tmp.Name
		u.Description = tmp.Description
		u.Format = tmp.Format
	case "tool_search":
		u.Execution = tmp.Execution
		u.Filters3 = tmp.Filters3
		u.MaxResults = tmp.MaxResults
		u.RankingOpts = tmp.RankingOpts
	case "web_search_preview":
		u.SearchContextSize = tmp.SearchContextSize
		u.UserLocation = tmp.UserLocation
	case "namespace":
		u.NamespaceName = tmp.NamespaceName
		u.Description = tmp.Description
	}
	return nil
}

// OpenAIResponsesFunctionToolParam 函数工具。
type OpenAIResponsesFunctionToolParam struct {
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Name 函数名。
	Name string `json:"name"`
	// Description 函数描述。
	Description string `json:"description,omitempty"`
	// Strict 是否启用 strict 模式。
	Strict *bool `json:"strict,omitempty"`
	// Parameters 参数 JSON Schema。
	Parameters json.RawMessage `json:"parameters,omitempty"`
}

func (OpenAIResponsesFunctionToolParam) isResponsesTool() {}

// OpenAIResponsesWebSearchToolParam Web 搜索工具。
type OpenAIResponsesWebSearchToolParam struct {
	// Type 固定 "web_search"。
	Type string `json:"type"` // "web_search"
	// Filters 过滤器（可选）。
	Filters *OpenAIResponsesWebSearchFilters `json:"filters,omitempty"`
	// SearchContextSize 搜索上下文大小。
	SearchContextSize string `json:"search_context_size,omitempty"`
	// UserLocation 用户位置。
	UserLocation *OpenAIResponsesWebSearchUserLocation `json:"user_location,omitempty"`
}

func (OpenAIResponsesWebSearchToolParam) isResponsesTool() {}

// OpenAIResponsesWebSearchFilters Web 搜索过滤器。
type OpenAIResponsesWebSearchFilters struct {
	AllowedDomains []string `json:"allowed_domains,omitempty"`
}

// OpenAIResponsesWebSearchUserLocation Web 搜索用户位置。
type OpenAIResponsesWebSearchUserLocation struct {
	Type     string `json:"type"` // "approximate"
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

// OpenAIResponsesCodeInterpreterToolParam 代码解释器工具。
type OpenAIResponsesCodeInterpreterToolParam struct {
	// Type 固定 "code_interpreter"。
	Type string `json:"type"` // "code_interpreter"
	// Container 容器配置。
	Container OpenAIResponsesCodeInterpreterContainer `json:"container"`
}

func (OpenAIResponsesCodeInterpreterToolParam) isResponsesTool() {}

// OpenAIResponsesCodeInterpreterContainer 容器配置。
type OpenAIResponsesCodeInterpreterContainer struct {
	// Type 固定 "auto"。
	Type string `json:"type"` // "auto"
}

// OpenAIResponsesFileSearchToolParam 文件搜索工具。
type OpenAIResponsesFileSearchToolParam struct {
	// Type 固定 "file_search"。
	Type string `json:"type"` // "file_search"
	// VectorStoreIDs 向量存储 ID 列表。
	VectorStoreIDs []string `json:"vector_store_ids"`
	// MaxNumResults 最大结果数。
	MaxNumResults *int64 `json:"max_num_results,omitempty"`
	// RankingOptions 排序选项。
	RankingOptions *OpenAIResponsesFileSearchRanking `json:"ranking_options,omitempty"`
	// Filters 过滤器。
	Filters json.RawMessage `json:"filters,omitempty"`
}

func (OpenAIResponsesFileSearchToolParam) isResponsesTool() {}

// OpenAIResponsesFileSearchRanking 排序选项。
type OpenAIResponsesFileSearchRanking struct {
	Ranker         string  `json:"ranker,omitempty"`
	ScoreThreshold float64 `json:"score_threshold,omitempty"`
}

// OpenAIResponsesImageGenerationToolParam 图像生成工具。
type OpenAIResponsesImageGenerationToolParam struct {
	// Type 固定 "image_generation"。
	Type string `json:"type"` // "image_generation"
	// Background 背景（transparent / opaque / auto）。
	Background string `json:"background,omitempty"`
	// InputFidelity 输入保真度（high / low）。
	InputFidelity string `json:"input_fidelity,omitempty"`
	// Model 模型名（如 "gpt-image-1"）。
	Model string `json:"model,omitempty"`
	// Moderation 内容审核（auto / low）。
	Moderation string `json:"moderation,omitempty"`
	// Quality 质量（low / medium / high / auto）。
	Quality string `json:"quality,omitempty"`
	// OutputFormat 输出格式（png / webp / jpeg）。
	OutputFormat string `json:"output_format,omitempty"`
	// Size 尺寸（1024x1024 / 1024x1536 / 1536x1024 / auto）。
	Size string `json:"size,omitempty"`
}

func (OpenAIResponsesImageGenerationToolParam) isResponsesTool() {}

// OpenAIResponsesMcpToolParam MCP 工具。
type OpenAIResponsesMcpToolParam struct {
	// Type 固定 "mcp"。
	Type string `json:"type"` // "mcp"
	// ServerLabel 服务器标签。
	ServerLabel string `json:"server_label"`
	// ServerURL 服务器 URL。
	ServerURL string `json:"server_url,omitempty"`
	// Headers HTTP 头。
	Headers map[string]string `json:"headers,omitempty"`
	// AllowedTools 允许的工具列表。
	AllowedTools []string `json:"allowed_tools,omitempty"`
	// RequireApproval 是否需要批准（always / never）。
	RequireApproval string `json:"require_approval,omitempty"`
}

func (OpenAIResponsesMcpToolParam) isResponsesTool() {}

// OpenAIResponsesComputerToolParam 电脑使用工具。
type OpenAIResponsesComputerToolParam struct {
	// Type 固定 "computer_use_preview"。
	Type string `json:"type"` // "computer_use_preview"
	// DisplayWidth 显示器宽度。
	DisplayWidth int64 `json:"display_width"`
	// DisplayHeight 显示器高度。
	DisplayHeight int64 `json:"display_height"`
	// Environment 操作系统环境。
	Environment string `json:"environment"`
}

func (OpenAIResponsesComputerToolParam) isResponsesTool() {}

// OpenAIResponsesApplyPatchToolParam 应用补丁工具。
type OpenAIResponsesApplyPatchToolParam struct {
	// Type 固定 "apply_patch"。
	Type string `json:"type"` // "apply_patch"
}

func (OpenAIResponsesApplyPatchToolParam) isResponsesTool() {}

// OpenAIResponsesCodeExecutionToolParam 代码执行工具。
type OpenAIResponsesCodeExecutionToolParam struct {
	// Type 固定 "code_execution"。
	Type string `json:"type"` // "code_execution"
}

func (OpenAIResponsesCodeExecutionToolParam) isResponsesTool() {}

// OpenAIResponsesHostedToolParam hosted 工具（通用兜底）。
type OpenAIResponsesHostedToolParam struct {
	// Type hosted 工具类型（如 "web_search_preview"）。
	Type string `json:"type"`
}

func (OpenAIResponsesHostedToolParam) isResponsesTool() {}

// === ToolChoice union 实现 ===

// OpenAIResponseToolChoiceOptionAuto 自动选择。
type OpenAIResponseToolChoiceOptionAuto struct {
	Type string `json:"type"` // "auto"
}

func (OpenAIResponseToolChoiceOptionAuto) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionRequired 强制调用。
type OpenAIResponseToolChoiceOptionRequired struct {
	Type string `json:"type"` // "required"
}

func (OpenAIResponseToolChoiceOptionRequired) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionNone 禁止调用。
type OpenAIResponseToolChoiceOptionNone struct {
	Type string `json:"type"` // "none"
}

func (OpenAIResponseToolChoiceOptionNone) isResponseToolChoiceUnion() {}

// === ToolChoice 命名常量基线 ===
//
// 工具选择策略中 auto / none / required 三种字符串形态通过 string alias 承载，
// 编译期防止 wire-format 字符串拼写错误。
type OpenAIResponseToolChoiceOptionMode string

const (
	OpenAIResponseToolChoiceOptionModeAuto     OpenAIResponseToolChoiceOptionMode = "auto"
	OpenAIResponseToolChoiceOptionModeNone     OpenAIResponseToolChoiceOptionMode = "none"
	OpenAIResponseToolChoiceOptionModeRequired OpenAIResponseToolChoiceOptionMode = "required"
)

// === ToolChoice 高级变体（4.6+）===

// OpenAIResponseToolChoiceOptionAllowedParam 允许的工具列表（带 mode 模式）。
type OpenAIResponseToolChoiceOptionAllowedParam struct {
	// Type 固定 "allowed_tools"。
	Type string `json:"type"` // "allowed_tools"
	// Mode 选择模式（auto / required）。
	Mode OpenAIResponseToolChoiceOptionMode `json:"mode"`
	// Tools 允许的工具列表。
	Tools []OpenAIResponseToolChoiceOptionAllowedToolEntry `json:"tools"`
}

func (OpenAIResponseToolChoiceOptionAllowedParam) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionAllowedToolEntry 允许的工具项。
type OpenAIResponseToolChoiceOptionAllowedToolEntry struct {
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Name 函数名。
	Name string `json:"name"`
}

// OpenAIResponseToolChoiceOptionApplyPatchParam 指定 apply_patch 工具。
type OpenAIResponseToolChoiceOptionApplyPatchParam struct {
	// Type 固定 "apply_patch"。
	Type string `json:"type"` // "apply_patch"
}

func (OpenAIResponseToolChoiceOptionApplyPatchParam) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionCustomParam 指定 custom 工具（按名称）。
type OpenAIResponseToolChoiceOptionCustomParam struct {
	// Type 固定 "custom"。
	Type string `json:"type"` // "custom"
	// Name 工具名。
	Name string `json:"name"`
}

func (OpenAIResponseToolChoiceOptionCustomParam) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionFunctionParam 指定 function 工具（按名称）。
type OpenAIResponseToolChoiceOptionFunctionParam struct {
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Name 函数名。
	Name string `json:"name"`
}

func (OpenAIResponseToolChoiceOptionFunctionParam) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionMcpParam 指定 MCP 工具。
type OpenAIResponseToolChoiceOptionMcpParam struct {
	// Type 固定 "mcp"。
	Type string `json:"type"` // "mcp"
	// ServerLabel MCP 服务器标签。
	ServerLabel string `json:"server_label"`
	// Name 工具名。
	Name string `json:"name"`
}

func (OpenAIResponseToolChoiceOptionMcpParam) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionShellParam 指定 shell 工具。
type OpenAIResponseToolChoiceOptionShellParam struct {
	// Type 固定 "shell"。
	Type string `json:"type"` // "shell"
}

func (OpenAIResponseToolChoiceOptionShellParam) isResponseToolChoiceUnion() {}

// OpenAIResponseToolChoiceOptionTypesParam 复合工具选择（type 为字符串数组）。
//
// 用于指定"可调用工具类型集合"（如 ["function", "web_search"]）。
type OpenAIResponseToolChoiceOptionTypesParam struct {
	// Type 工具类型列表（如 ["function", "web_search"]）。
	Type []string `json:"type"`
}

func (OpenAIResponseToolChoiceOptionTypesParam) isResponseToolChoiceUnion() {}

// === Responses Tool 扩展变体（4.6+）===

// OpenAIResponsesLocalShellToolParam 本地 shell 工具（gpt-5.1 codex 用）。
type OpenAIResponsesLocalShellToolParam struct {
	// Type 固定 "local_shell"。
	Type string `json:"type"` // "local_shell"
}

func (OpenAIResponsesLocalShellToolParam) isResponsesTool() {}

// OpenAIResponsesShellToolParam 远程 shell 工具（gpt-5.1 codex 用）。
type OpenAIResponsesShellToolParam struct {
	// Type 固定 "shell"。
	Type string `json:"type"` // "shell"
	// Environment 执行环境配置。
	Environment *OpenAIResponsesShellEnvironment `json:"environment,omitempty"`
}

func (OpenAIResponsesShellToolParam) isResponsesTool() {}

// OpenAIResponsesShellEnvironment shell 执行环境配置。
type OpenAIResponsesShellEnvironment struct {
	// Type 固定 "container_auto"。
	Type string `json:"type"` // "container_auto"
}

// OpenAIResponsesCustomToolParam 自定义 grammar 工具（4.6+）。
type OpenAIResponsesCustomToolParam struct {
	// Type 固定 "custom"。
	Type string `json:"type"` // "custom"
	// Name 工具名。
	Name string `json:"name"`
	// Description 工具描述。
	Description string `json:"description,omitempty"`
	// Format 输出格式约束（grammar）。
	Format *OpenAIResponsesCustomToolFormatGrammar `json:"format,omitempty"`
}

func (OpenAIResponsesCustomToolParam) isResponsesTool() {}

// OpenAIResponsesCustomToolFormatGrammar 自定义语法格式约束。
type OpenAIResponsesCustomToolFormatGrammar struct {
	// Type 固定 "grammar"。
	Type string `json:"type"` // "grammar"
	// Syntax 语法类型（lark / regex）。
	Syntax string `json:"syntax"`
	// Definition 语法定义。
	Definition string `json:"definition"`
}

// OpenAIResponsesToolSearchToolParam 工具搜索工具（按名称/语义检索工具）。
type OpenAIResponsesToolSearchToolParam struct {
	// Type 固定 "tool_search"。
	Type string `json:"type"` // "tool_search"
	// Execution 执行模式（auto / bm25 / regex）。
	Execution string `json:"execution,omitempty"`
	// Filters 过滤规则。
	Filters []string `json:"filters,omitempty"`
	// MaxResults 最大结果数。
	MaxResults *int64 `json:"max_num_results,omitempty"`
	// RankingOpts 排序选项。
	RankingOpts *OpenAIResponsesToolSearchRanking `json:"ranking_options,omitempty"`
}

func (OpenAIResponsesToolSearchToolParam) isResponsesTool() {}

// OpenAIResponsesToolSearchRanking 工具搜索排序选项。
type OpenAIResponsesToolSearchRanking struct {
	// Ranker 排序算法。
	Ranker string `json:"ranker,omitempty"`
	// ScoreThreshold 分数阈值。
	ScoreThreshold float64 `json:"score_threshold,omitempty"`
}

// OpenAIResponsesWebSearchPreviewToolParam Web 搜索（preview 版本，旧版 hosted 工具）。
type OpenAIResponsesWebSearchPreviewToolParam struct {
	// Type 固定 "web_search_preview"。
	Type string `json:"type"` // "web_search_preview"
	// SearchContextSize 搜索上下文大小。
	SearchContextSize string `json:"search_context_size,omitempty"`
	// UserLocation 用户位置。
	UserLocation *OpenAIResponsesWebSearchUserLocation `json:"user_location,omitempty"`
}

func (OpenAIResponsesWebSearchPreviewToolParam) isResponsesTool() {}

// OpenAIResponsesNamespaceToolParam 命名空间工具（批量挂载工具集合）。
type OpenAIResponsesNamespaceToolParam struct {
	// Type 固定 "namespace"。
	Type string `json:"type"` // "namespace"
	// Name 命名空间名。
	Name string `json:"name"`
	// Description 命名空间描述。
	Description string `json:"description,omitempty"`
}

func (OpenAIResponsesNamespaceToolParam) isResponsesTool() {}
