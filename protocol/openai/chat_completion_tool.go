package openai

import "encoding/json"

// OpenAIChatCompletionToolUnion 工具 union 接口。
//
// 7 种实现：Function / Custom / WebSearch / CodeInterpreter / FileSearch / ComputerUse / ImageGeneration。
type OpenAIChatCompletionToolUnion interface {
	isChatCompletionTool()
}

// OpenAIChatCompletionToolUnionFlat 工具 union 扁平化实现。
type OpenAIChatCompletionToolUnionFlat struct {
	// Type 标识当前激活的变体。
	Type string `json:"type,omitempty"`
	// Function 函数工具。
	Function *OpenAIChatCompletionToolFunctionParam `json:"function,omitempty"`
	// Custom 自定义工具。
	Custom *OpenAIChatCompletionToolCustomParam `json:"custom,omitempty"`
	// WebSearch Web 搜索工具。
	WebSearch *OpenAIChatCompletionToolWebSearchParam `json:"web_search,omitempty"`
	// CodeInterpreter 代码解释器工具。
	CodeInterpreter *OpenAIChatCompletionToolCodeInterpreterParam `json:"code_interpreter,omitempty"`
	// FileSearch 文件搜索工具。
	FileSearch *OpenAIChatCompletionToolFileSearchParam `json:"file_search,omitempty"`
	// ComputerUse 电脑使用工具。
	ComputerUse *OpenAIChatCompletionToolComputerUseParam `json:"computer_use_preview,omitempty"`
	// ImageGeneration 图像生成工具。
	ImageGeneration *OpenAIChatCompletionToolImageGenerationParam `json:"image_generation,omitempty"`
}

func (OpenAIChatCompletionToolUnionFlat) isChatCompletionTool() {}

func (u OpenAIChatCompletionToolUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "function":
		if u.Function == nil {
			return nil, nil
		}
		return json.Marshal(u.Function)
	case "custom":
		if u.Custom == nil {
			return nil, nil
		}
		return json.Marshal(u.Custom)
	case "web_search":
		if u.WebSearch == nil {
			return nil, nil
		}
		return json.Marshal(u.WebSearch)
	case "code_interpreter":
		if u.CodeInterpreter == nil {
			return nil, nil
		}
		return json.Marshal(u.CodeInterpreter)
	case "file_search":
		if u.FileSearch == nil {
			return nil, nil
		}
		return json.Marshal(u.FileSearch)
	case "computer_use_preview":
		if u.ComputerUse == nil {
			return nil, nil
		}
		return json.Marshal(u.ComputerUse)
	case "image_generation":
		if u.ImageGeneration == nil {
			return nil, nil
		}
		return json.Marshal(u.ImageGeneration)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionToolUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "function":
		var v OpenAIChatCompletionToolFunctionParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Function = &v
	case "custom":
		var v OpenAIChatCompletionToolCustomParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Custom = &v
	case "web_search":
		var v OpenAIChatCompletionToolWebSearchParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.WebSearch = &v
	case "code_interpreter":
		var v OpenAIChatCompletionToolCodeInterpreterParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.CodeInterpreter = &v
	case "file_search":
		var v OpenAIChatCompletionToolFileSearchParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.FileSearch = &v
	case "computer_use_preview":
		var v OpenAIChatCompletionToolComputerUseParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ComputerUse = &v
	case "image_generation":
		var v OpenAIChatCompletionToolImageGenerationParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ImageGeneration = &v
	}
	return nil
}

// OpenAIChatCompletionToolFunctionParam 函数工具。
type OpenAIChatCompletionToolFunctionParam struct {
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Function 函数定义。
	Function OpenAIChatCompletionToolFunctionDefinition `json:"function"`
}

func (OpenAIChatCompletionToolFunctionParam) isChatCompletionTool() {}

// OpenAIChatCompletionToolFunctionDefinition 函数定义。
type OpenAIChatCompletionToolFunctionDefinition struct {
	// Name 函数名。
	Name string `json:"name"`
	// Description 函数描述（用于模型理解何时调用）。
	Description string `json:"description,omitempty"`
	// Parameters 参数 JSON Schema。
	Parameters json.RawMessage `json:"parameters,omitempty"`
	// Strict 是否启用 strict 模式。
	Strict *bool `json:"strict,omitempty"`
}

// OpenAIChatCompletionToolCustomParam 自定义工具（4.6+）。
type OpenAIChatCompletionToolCustomParam struct {
	// Type 固定 "custom"。
	Type string `json:"type"` // "custom"
	// Custom 自定义工具定义。
	Custom OpenAIChatCompletionToolCustomDefinition `json:"custom"`
}

func (OpenAIChatCompletionToolCustomParam) isChatCompletionTool() {}

// OpenAIChatCompletionToolCustomDefinition 自定义工具定义。
type OpenAIChatCompletionToolCustomDefinition struct {
	// Name 工具名。
	Name string `json:"name"`
	// Description 工具描述。
	Description string `json:"description,omitempty"`
	// Format 工具输出格式（文本或 grammar）。
	Format OpenAIChatCompletionToolCustomFormatUnion `json:"format,omitempty"`
}

// OpenAIChatCompletionToolCustomFormatUnion 自定义工具格式 union。
type OpenAIChatCompletionToolCustomFormatUnion interface {
	isChatCompletionToolCustomFormat()
}

// OpenAIChatCompletionToolCustomFormatUnionFlat 自定义工具格式 union 扁平化实现。
type OpenAIChatCompletionToolCustomFormatUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"text"/"grammar"。
	Type string `json:"type,omitempty"`
	// Text 纯文本格式。
	Text *OpenAIChatCompletionToolCustomFormatText `json:"text,omitempty"`
	// Grammar 自定义语法格式。
	Grammar *OpenAIChatCompletionToolCustomFormatGrammar `json:"grammar,omitempty"`
}

func (OpenAIChatCompletionToolCustomFormatUnionFlat) isChatCompletionToolCustomFormat() {}

func (u OpenAIChatCompletionToolCustomFormatUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "text":
		if u.Text == nil {
			return nil, nil
		}
		return json.Marshal(u.Text)
	case "grammar":
		if u.Grammar == nil {
			return nil, nil
		}
		return json.Marshal(u.Grammar)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionToolCustomFormatUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "text":
		var v OpenAIChatCompletionToolCustomFormatText
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Text = &v
	case "grammar":
		var v OpenAIChatCompletionToolCustomFormatGrammar
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Grammar = &v
	}
	return nil
}

// OpenAIChatCompletionToolCustomFormatText 纯文本格式。
type OpenAIChatCompletionToolCustomFormatText struct {
	Type string `json:"type"` // "text"
}

func (OpenAIChatCompletionToolCustomFormatText) isChatCompletionToolCustomFormat() {}

// OpenAIChatCompletionToolCustomFormatGrammar 自定义语法格式。
type OpenAIChatCompletionToolCustomFormatGrammar struct {
	// Type 固定 "grammar"。
	Type string `json:"type"` // "grammar"
	// Syntax 语法类型（lark / regex）。
	Syntax string `json:"syntax"`
	// Definition 语法定义。
	Definition string `json:"definition"`
}

func (OpenAIChatCompletionToolCustomFormatGrammar) isChatCompletionToolCustomFormat() {}

// OpenAIChatCompletionToolWebSearchParam Web 搜索工具（hosted）。
type OpenAIChatCompletionToolWebSearchParam struct {
	// Type 固定 "web_search"。
	Type string `json:"type"` // "web_search" (hosted tool)
	// WebSearchOptions 搜索选项（可选）。
	WebSearchOptions *OpenAIChatCompletionWebSearchOptions `json:"web_search_options,omitempty"`
}

func (OpenAIChatCompletionToolWebSearchParam) isChatCompletionTool() {}

// OpenAIChatCompletionToolCodeInterpreterParam 代码解释器工具（hosted）。
type OpenAIChatCompletionToolCodeInterpreterParam struct {
	// Type 固定 "code_interpreter"。
	Type string `json:"type"` // "code_interpreter"
	// Container 容器 ID（可选）。
	Container string `json:"container,omitempty"`
}

func (OpenAIChatCompletionToolCodeInterpreterParam) isChatCompletionTool() {}

// OpenAIChatCompletionToolFileSearchParam 文件搜索工具（hosted）。
type OpenAIChatCompletionToolFileSearchParam struct {
	// Type 固定 "file_search"。
	Type string `json:"type"` // "file_search"
	// VectorStoreIDs 向量存储 ID 列表。
	VectorStoreIDs []string `json:"vector_store_ids,omitempty"`
	// MaxNumResults 最大结果数。
	MaxNumResults *int64 `json:"max_num_results,omitempty"`
	// RankingOptions 排序选项。
	RankingOptions *OpenAIChatCompletionToolFileSearchRanking `json:"ranking_options,omitempty"`
	// Filters 过滤器。
	Filters json.RawMessage `json:"filters,omitempty"`
}

func (OpenAIChatCompletionToolFileSearchParam) isChatCompletionTool() {}

// OpenAIChatCompletionToolFileSearchRanking 排序选项。
type OpenAIChatCompletionToolFileSearchRanking struct {
	// Ranker 排序算法。
	Ranker string `json:"ranker,omitempty"`
	// ScoreThreshold 分数阈值。
	ScoreThreshold float64 `json:"score_threshold,omitempty"`
}

// OpenAIChatCompletionToolComputerUseParam 电脑使用工具（hosted）。
type OpenAIChatCompletionToolComputerUseParam struct {
	// Type 固定 "computer_use_preview"。
	Type string `json:"type"` // "computer_use_preview"
	// DisplayWidth 显示器宽度。
	DisplayWidth int64 `json:"display_width"`
	// DisplayHeight 显示器高度。
	DisplayHeight int64 `json:"display_height"`
	// Environment 操作系统环境。
	Environment string `json:"environment"`
}

func (OpenAIChatCompletionToolComputerUseParam) isChatCompletionTool() {}

// OpenAIChatCompletionToolImageGenerationType 图像生成工具类型字符串。
type OpenAIChatCompletionToolImageGenerationType string

const (
	OpenAIChatCompletionToolImageGenerationTypeImageGeneration OpenAIChatCompletionToolImageGenerationType = "image_generation"
)

// OpenAIChatCompletionToolImageGenerationParam 图像生成工具（hosted，4.6+）。
type OpenAIChatCompletionToolImageGenerationParam struct {
	// Type 固定 "image_generation"。
	Type OpenAIChatCompletionToolImageGenerationType `json:"type"`
	// Background 背景类型（transparent / opaque / auto）。
	Background string `json:"background,omitempty"`
	// InputFidelity 输入保真度（high / low）。
	InputFidelity string `json:"input_fidelity,omitempty"`
	// Model 图像生成模型（如 "gpt-image-1"）。
	Model string `json:"model,omitempty"`
	// Moderation 审核级别（auto / low）。
	Moderation string `json:"moderation,omitempty"`
	// OutputCompression 输出压缩率（0-100）。
	OutputCompression *int64 `json:"output_compression,omitempty"`
	// PartialImages 部分图像模式（0=禁用，1-3=流式分片）。
	PartialImages *int64 `json:"partial_images,omitempty"`
	// Quality 图像质量（low / medium / high / auto）。
	Quality string `json:"quality,omitempty"`
	// Size 图像尺寸（如 "1024x1024" / "auto"）。
	Size string `json:"size,omitempty"`
	// OutputFormat 输出格式（png / jpeg / webp）。
	OutputFormat string `json:"output_format,omitempty"`
}

func (OpenAIChatCompletionToolImageGenerationParam) isChatCompletionTool() {}

// === ToolChoice union ===

// OpenAIChatCompletionToolChoiceOptionType 工具选择 union 中用于识别变体形态的 type 字段。
type OpenAIChatCompletionToolChoiceOptionType string

const (
	OpenAIChatCompletionToolChoiceOptionTypeAuto         OpenAIChatCompletionToolChoiceOptionType = "auto"
	OpenAIChatCompletionToolChoiceOptionTypeRequired     OpenAIChatCompletionToolChoiceOptionType = "required"
	OpenAIChatCompletionToolChoiceOptionTypeNone         OpenAIChatCompletionToolChoiceOptionType = "none"
	OpenAIChatCompletionToolChoiceOptionTypeAllowedTools OpenAIChatCompletionToolChoiceOptionType = "allowed_tools"
	OpenAIChatCompletionToolChoiceOptionTypeFunction     OpenAIChatCompletionToolChoiceOptionType = "function"
	OpenAIChatCompletionToolChoiceOptionTypeCustom       OpenAIChatCompletionToolChoiceOptionType = "custom"
)

// OpenAIChatCompletionToolChoiceOptionUnionParam 工具选择策略 union。
//
// 6 种实现：Auto / Required / None / Allowed / HostedNamed / HostedNamedCustom。
type OpenAIChatCompletionToolChoiceOptionUnionParam interface {
	isChatCompletionToolChoiceOption()
}

// OpenAIChatCompletionToolChoiceOptionUnionParamFlat 工具选择策略 union 扁平化实现。
type OpenAIChatCompletionToolChoiceOptionUnionParamFlat struct {
	// Type 标识当前激活的变体。
	Type string `json:"type,omitempty"`
	// Auto 自动选择。
	Auto *OpenAIChatCompletionToolChoiceOptionAutoParam `json:"auto,omitempty"`
	// Required 强制调用工具。
	Required *OpenAIChatCompletionToolChoiceOptionRequiredParam `json:"required,omitempty"`
	// None 禁止调用工具。
	None *OpenAIChatCompletionToolChoiceOptionNoneParam `json:"none,omitempty"`
	// Allowed 允许的工具列表。
	Allowed *OpenAIChatCompletionToolChoiceOptionAllowedParam `json:"allowed_tools,omitempty"`
	// HostedNamed 指定 hosted 工具（function 命名）。
	HostedNamed *OpenAIChatCompletionToolChoiceOptionHostedNamedParam `json:"function,omitempty"`
	// HostedNamedCustom 指定 hosted 工具（custom 命名）。
	HostedNamedCustom *OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam `json:"custom,omitempty"`
}

func (OpenAIChatCompletionToolChoiceOptionUnionParamFlat) isChatCompletionToolChoiceOption() {}

func (u OpenAIChatCompletionToolChoiceOptionUnionParamFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "auto":
		if u.Auto == nil {
			return nil, nil
		}
		return json.Marshal(u.Auto)
	case "required":
		if u.Required == nil {
			return nil, nil
		}
		return json.Marshal(u.Required)
	case "none":
		if u.None == nil {
			return nil, nil
		}
		return json.Marshal(u.None)
	case "allowed_tools":
		if u.Allowed == nil {
			return nil, nil
		}
		return json.Marshal(u.Allowed)
	case "function":
		if u.HostedNamed == nil {
			return nil, nil
		}
		return json.Marshal(u.HostedNamed)
	case "custom":
		if u.HostedNamedCustom == nil {
			return nil, nil
		}
		return json.Marshal(u.HostedNamedCustom)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionToolChoiceOptionUnionParamFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "auto":
		var v OpenAIChatCompletionToolChoiceOptionAutoParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Auto = &v
	case "required":
		var v OpenAIChatCompletionToolChoiceOptionRequiredParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Required = &v
	case "none":
		var v OpenAIChatCompletionToolChoiceOptionNoneParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.None = &v
	case "allowed_tools":
		var v OpenAIChatCompletionToolChoiceOptionAllowedParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Allowed = &v
	case "function":
		var v OpenAIChatCompletionToolChoiceOptionHostedNamedParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.HostedNamed = &v
	case "custom":
		var v OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.HostedNamedCustom = &v
	}
	return nil
}

// OpenAIChatCompletionToolChoiceOptionAutoParam 自动选择。
type OpenAIChatCompletionToolChoiceOptionAutoParam struct {
	Type string `json:"type"` // "auto"
}

func (OpenAIChatCompletionToolChoiceOptionAutoParam) isChatCompletionToolChoiceOption() {}

// OpenAIChatCompletionToolChoiceOptionRequiredParam 强制调用工具。
type OpenAIChatCompletionToolChoiceOptionRequiredParam struct {
	Type string `json:"type"` // "required"
}

func (OpenAIChatCompletionToolChoiceOptionRequiredParam) isChatCompletionToolChoiceOption() {}

// OpenAIChatCompletionToolChoiceOptionNoneParam 禁止调用工具。
type OpenAIChatCompletionToolChoiceOptionNoneParam struct {
	Type string `json:"type"` // "none"
}

func (OpenAIChatCompletionToolChoiceOptionNoneParam) isChatCompletionToolChoiceOption() {}

// OpenAIChatCompletionToolChoiceOptionAllowedParam 允许的工具列表（4.6+）。
type OpenAIChatCompletionToolChoiceOptionAllowedParam struct {
	// Type 固定 "allowed_tools"。
	Type string `json:"type"` // "allowed_tools"
	// AllowedTools 允许的工具定义列表。
	AllowedTools []OpenAIChatCompletionToolChoiceOptionAllowedTool `json:"allowed_tools"`
	// Mode 选择模式（auto / required）。
	Mode string `json:"mode,omitempty"`
}

func (OpenAIChatCompletionToolChoiceOptionAllowedParam) isChatCompletionToolChoiceOption() {}

// OpenAIChatCompletionToolChoiceOptionAllowedTool 允许的工具项。
type OpenAIChatCompletionToolChoiceOptionAllowedTool struct {
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Function 函数定义。
	Function OpenAIChatCompletionToolChoiceOptionAllowedFunction `json:"function"`
}

// OpenAIChatCompletionToolChoiceOptionAllowedFunction 允许的函数定义。
type OpenAIChatCompletionToolChoiceOptionAllowedFunction struct {
	Name string `json:"name"`
}

// OpenAIChatCompletionToolChoiceOptionHostedNamedParam 指定 hosted 工具（function 命名）。
type OpenAIChatCompletionToolChoiceOptionHostedNamedParam struct {
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Function 函数引用。
	Function OpenAIChatCompletionToolChoiceOptionHostedNamedFunction `json:"function"`
}

func (OpenAIChatCompletionToolChoiceOptionHostedNamedParam) isChatCompletionToolChoiceOption() {}

// OpenAIChatCompletionToolChoiceOptionHostedNamedFunction hosted 工具引用。
type OpenAIChatCompletionToolChoiceOptionHostedNamedFunction struct {
	Name string `json:"name"`
}

// OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam 指定 hosted 工具（custom 命名）。
type OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam struct {
	// Type 固定 "custom"。
	Type string `json:"type"` // "custom"
	// Custom 自定义工具引用。
	Custom OpenAIChatCompletionToolChoiceOptionHostedNamedCustom `json:"custom"`
}

func (OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam) isChatCompletionToolChoiceOption() {
}

// OpenAIChatCompletionToolChoiceOptionHostedNamedCustom hosted 自定义工具引用。
type OpenAIChatCompletionToolChoiceOptionHostedNamedCustom struct {
	Name string `json:"name"`
}

// === 消息中的 ToolCall union ===

// OpenAIChatCompletionMessageToolCallUnion 消息中工具调用 union（Request 侧）。
//
// 2 种实现：Function / Custom。
type OpenAIChatCompletionMessageToolCallUnion interface {
	isChatCompletionMessageToolCall()
}

// OpenAIChatCompletionMessageToolCallUnionFlat 消息中工具调用 union 扁平化实现。
type OpenAIChatCompletionMessageToolCallUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"function"/"custom"。
	Type string `json:"type,omitempty"`
	// Function 函数工具调用。
	Function *OpenAIChatCompletionMessageToolCallFunctionParam `json:"function,omitempty"`
	// Custom 自定义工具调用。
	Custom *OpenAIChatCompletionMessageToolCallCustomParam `json:"custom,omitempty"`
}

func (OpenAIChatCompletionMessageToolCallUnionFlat) isChatCompletionMessageToolCall() {}

func (u OpenAIChatCompletionMessageToolCallUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "function":
		if u.Function == nil {
			return nil, nil
		}
		return json.Marshal(u.Function)
	case "custom":
		if u.Custom == nil {
			return nil, nil
		}
		return json.Marshal(u.Custom)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionMessageToolCallUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "function":
		var v OpenAIChatCompletionMessageToolCallFunctionParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Function = &v
	case "custom":
		var v OpenAIChatCompletionMessageToolCallCustomParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Custom = &v
	}
	return nil
}

// OpenAIChatCompletionMessageToolCallFunctionParam 函数工具调用。
type OpenAIChatCompletionMessageToolCallFunctionParam struct {
	// ID 工具调用 ID。
	ID string `json:"id"`
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Function 函数调用详情。
	Function OpenAIChatCompletionMessageToolCallFunction `json:"function"`
}

func (OpenAIChatCompletionMessageToolCallFunctionParam) isChatCompletionMessageToolCall() {}

// OpenAIChatCompletionMessageToolCallFunction 函数调用详情。
type OpenAIChatCompletionMessageToolCallFunction struct {
	// Name 函数名。
	Name string `json:"name"`
	// Arguments 参数 JSON 字符串。
	Arguments string `json:"arguments"`
}

// OpenAIChatCompletionMessageToolCallCustomParam 自定义工具调用。
type OpenAIChatCompletionMessageToolCallCustomParam struct {
	// ID 工具调用 ID。
	ID string `json:"id"`
	// Type 固定 "custom"。
	Type string `json:"type"` // "custom"
	// Custom 自定义工具调用详情。
	Custom OpenAIChatCompletionMessageToolCallCustom `json:"custom"`
}

func (OpenAIChatCompletionMessageToolCallCustomParam) isChatCompletionMessageToolCall() {}

// OpenAIChatCompletionMessageToolCallCustom 自定义工具调用详情。
type OpenAIChatCompletionMessageToolCallCustom struct {
	// Name 工具名。
	Name string `json:"name"`
	// Input 工具输入。
	Input string `json:"input"`
}

// OpenAIChatCompletionChoiceMessageToolCallUnion 响应侧工具调用 union。
//
// 2 种实现：Function / Custom。
type OpenAIChatCompletionChoiceMessageToolCallUnion interface {
	isChatCompletionChoiceMessageToolCall()
}

// OpenAIChatCompletionChoiceMessageToolCallUnionFlat 响应侧工具调用 union 扁平化实现。
type OpenAIChatCompletionChoiceMessageToolCallUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"function"/"custom"。
	Type string `json:"type,omitempty"`
	// Function 函数工具调用。
	Function *OpenAIChatCompletionChoiceMessageToolCallFunctionParam `json:"function,omitempty"`
	// Custom 自定义工具调用。
	Custom *OpenAIChatCompletionChoiceMessageToolCallCustomParam `json:"custom,omitempty"`
}

func (OpenAIChatCompletionChoiceMessageToolCallUnionFlat) isChatCompletionChoiceMessageToolCall() {}

func (u OpenAIChatCompletionChoiceMessageToolCallUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "function":
		if u.Function == nil {
			return nil, nil
		}
		return json.Marshal(u.Function)
	case "custom":
		if u.Custom == nil {
			return nil, nil
		}
		return json.Marshal(u.Custom)
	}
	return nil, nil
}

func (u *OpenAIChatCompletionChoiceMessageToolCallUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "function":
		var v OpenAIChatCompletionChoiceMessageToolCallFunctionParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Function = &v
	case "custom":
		var v OpenAIChatCompletionChoiceMessageToolCallCustomParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Custom = &v
	}
	return nil
}

// OpenAIChatCompletionChoiceMessageToolCallFunctionParam 响应侧函数工具调用。
type OpenAIChatCompletionChoiceMessageToolCallFunctionParam struct {
	// ID 工具调用 ID。
	ID string `json:"id"`
	// Type 固定 "function"。
	Type string `json:"type"` // "function"
	// Function 函数调用详情。
	Function OpenAIChatCompletionMessageToolCallFunction `json:"function"`
}

func (OpenAIChatCompletionChoiceMessageToolCallFunctionParam) isChatCompletionChoiceMessageToolCall() {
}

// OpenAIChatCompletionChoiceMessageToolCallCustomParam 响应侧自定义工具调用。
type OpenAIChatCompletionChoiceMessageToolCallCustomParam struct {
	// ID 工具调用 ID。
	ID string `json:"id"`
	// Type 固定 "custom"。
	Type string `json:"type"` // "custom"
	// Custom 自定义工具调用详情。
	Custom OpenAIChatCompletionMessageToolCallCustom `json:"custom"`
}

func (OpenAIChatCompletionChoiceMessageToolCallCustomParam) isChatCompletionChoiceMessageToolCall() {}
