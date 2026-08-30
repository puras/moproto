package openai

import "encoding/json"

// OpenAIChatCompletionParams 是 ChatCompletion API 的请求顶层类型（替代官方 SDK
// ChatCompletionNewParams 的 New 中缀，与 anthropic 包的简洁命名风格对齐）。
//
// 字段集与 OpenAI 官方 Go SDK ChatCompletionNewParams 严格对齐。
// 详见 https://platform.openai.com/docs/api-reference/chat/create。
//
// 实现 common.Request 接口（GetModel / GetFormat）。
func (r *OpenAIChatCompletionParams) GetModel() string  { return r.Model }
func (r *OpenAIChatCompletionParams) GetFormat() string { return "openai" }

type OpenAIChatCompletionParams struct {
	// === 必填字段 ===

	// Model 使用的模型 ID（如 "gpt-4o"）。
	Model string `json:"model"`
	// Messages 对话历史。
	Messages []OpenAIChatCompletionMessageParamUnion `json:"messages"`

	// === 流式开关 ===

	// Stream 是否流式输出（指针用于区分"未设置"）。
	Stream *bool `json:"stream,omitempty"`

	// === 采样参数 ===

	// Temperature 采样温度（0-2）。
	Temperature *float64 `json:"temperature,omitempty"`
	// TopP 核采样阈值。
	TopP *float64 `json:"top_p,omitempty"`
	// FrequencyPenalty 频率惩罚（-2 到 2）。
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	// PresencePenalty 存在惩罚（-2 到 2）。
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`
	// Seed 随机种子。
	Seed *int64 `json:"seed,omitempty"`

	// === 输出限制 ===

	// MaxTokens 最大输出 token 数（GPT-4 系列）。
	MaxTokens *int64 `json:"max_tokens,omitempty"`
	// MaxCompletionTokens 最大输出 token 数（o1 系列推荐）。
	MaxCompletionTokens *int64 `json:"max_completion_tokens,omitempty"`
	// N 每个 prompt 生成的完成数。
	N *int64 `json:"n,omitempty"`

	// === 停止条件 ===

	// Stop 停止序列（string 或 []string）。
	Stop OpenAIStopUnion `json:"stop,omitempty"`

	// === 工具调用 ===

	// Tools 工具定义列表。
	Tools []OpenAIChatCompletionToolUnion `json:"tools,omitempty"`
	// ToolChoice 工具选择策略。
	ToolChoice OpenAIChatCompletionToolChoiceOptionUnionParam `json:"tool_choice,omitempty"`
	// ParallelToolCalls 是否允许并行工具调用。
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`

	// === 响应格式 ===

	// ResponseFormat 响应格式（text / json_object / json_schema）。
	ResponseFormat OpenAIResponseFormatUnion `json:"response_format,omitempty"`

	// === 流式选项 ===

	// StreamOptions 流式选项（包含 usage 统计）。
	StreamOptions *OpenAIChatCompletionStreamOptionsParam `json:"stream_options,omitempty"`

	// === Logprobs ===

	// Logprobs 是否返回 logprobs。
	Logprobs *bool `json:"logprobs,omitempty"`
	// TopLogprobs 每个 token 返回的 top logprobs 数（0-20）。
	TopLogprobs *int64 `json:"top_logprobs,omitempty"`

	// === Logit bias ===

	// LogitBias token ID 到 bias 值的映射。
	LogitBias map[string]int64 `json:"logit_bias,omitempty"`

	// === 用户标识 ===

	// User 用户标识（用于滥用检测）。
	User string `json:"user,omitempty"`

	// === 元数据 ===

	// Metadata 自定义元数据（key-value 对，最多 16 对）。
	Metadata map[string]string `json:"metadata,omitempty"`

	// === 存储 ===

	// Store 是否存储对话用于评估/微调。
	Store *bool `json:"store,omitempty"`

	// === 服务等级 ===

	// ServiceTier 服务等级。
	ServiceTier OpenAIServiceTier `json:"service_tier,omitempty"`

	// === 音频输出（4.6+ GA）===

	// Audio 音频输出配置。
	Audio *OpenAIChatCompletionAudioParam `json:"audio,omitempty"`
	// Modalities 输出模态列表。
	Modalities []OpenAIModality `json:"modalities,omitempty"`

	// === 预测内容（4.6+）===

	// Prediction 预测内容（用于 speculative decoding）。
	Prediction *OpenAIChatCompletionPredictionContentParam `json:"prediction,omitempty"`

	// === 推理（o 系列）===

	// ReasoningEffort 推理强度。
	ReasoningEffort OpenAIReasoningEffort `json:"reasoning_effort,omitempty"`

	// === 详细度 ===

	// Verbosity 输出详细度。
	Verbosity OpenAIVerbosity `json:"verbosity,omitempty"`

	// === Web 搜索 ===

	// WebSearchOptions Web 搜索选项。
	WebSearchOptions *OpenAIChatCompletionWebSearchOptions `json:"web_search_options,omitempty"`

	// === 安全（4.6+）===

	// SafetyIdentifier 安全标识（用于检测滥用）。
	SafetyIdentifier string `json:"safety_identifier,omitempty"`

	// === 提示缓存 ===

	// PromptCacheKey 提示缓存键。
	PromptCacheKey string `json:"prompt_cache_key,omitempty"`

	// === 提示缓存保留策略（4.6+）===

	// PromptCacheRetention 缓存保留策略。
	PromptCacheRetention OpenAIResponsePromptCacheRetention `json:"prompt_cache_retention,omitempty"`

	// === 审核（4.6+）===

	// Moderation 审核配置（对输入与生成输出进行内容审核）。
	Moderation *OpenAIChatCompletionModerationParam `json:"moderation,omitempty"`

	// === 输出混淆（4.6+）===

	// IncludeObfuscation 是否在响应中对敏感字段（如 tool_calls.arguments）进行混淆。
	// 与 stream_options.include_obfuscation 不同：此为顶层全局开关，stream 内字段已被官方移除。
	IncludeObfuscation *bool `json:"include_obfuscation,omitempty"`
}

// OpenAIChatCompletionModerationParam 审核参数（对请求输入与生成输出运行审核）。
//
// 参考官方 ChatCompletionNewParamsModeration 字段集。
type OpenAIChatCompletionModerationParam struct {
	// Model 审核模型 ID（如 "omni-moderation-latest"）。
	Model string `json:"model"`
}

// OpenAIStopUnion 停止序列的 union（string | []string）。
type OpenAIStopUnion interface {
	isStopUnion()
}

// OpenAIStopUnionFlat 停止序列 union 扁平化实现。
type OpenAIStopUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"string"/"strings"。
	Type string `json:"type,omitempty"`
	// String 单个停止字符串。
	String *OpenAIStopUnionString `json:"string,omitempty"`
	// Strings 多个停止字符串。
	Strings *OpenAIStopUnionStrings `json:"strings,omitempty"`
}

func (OpenAIStopUnionFlat) isStopUnion() {}

func (u OpenAIStopUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "strings":
		if u.Strings == nil {
			return nil, nil
		}
		return json.Marshal(u.Strings)
	}
	return nil, nil
}

func (u *OpenAIStopUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIStopUnionString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "strings":
		var v OpenAIStopUnionStrings
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Strings = &v
	}
	return nil
}

// OpenAIStopUnionString 单个停止字符串。
type OpenAIStopUnionString string

func (OpenAIStopUnionString) isStopUnion() {}

// OpenAIStopUnionStrings 多个停止字符串。
type OpenAIStopUnionStrings []string

func (OpenAIStopUnionStrings) isStopUnion() {}

// OpenAIChatCompletionStreamOptionsParam 流式选项。
type OpenAIChatCompletionStreamOptionsParam struct {
	// IncludeUsage 是否在最后一个 chunk 包含 usage 统计。
	IncludeUsage *bool `json:"include_usage,omitempty"`
	// IncludeObfuscation 是否在流式响应中对敏感字段（如 tool_calls.arguments）进行混淆。
	IncludeObfuscation *bool `json:"include_obfuscation,omitempty"`
}

// OpenAIChatCompletionAudioParam 音频输出参数。
type OpenAIChatCompletionAudioParam struct {
	// Voice 语音类型（字符串或 VoiceID 对象）。
	Voice OpenAIAudioVoiceUnion `json:"voice"`
	// Format 音频输出格式。
	Format OpenAIAudioFormat `json:"format"`
}

// OpenAIAudioVoiceUnion 语音类型 union（string | VoiceID）。
type OpenAIAudioVoiceUnion interface {
	isAudioVoiceUnion()
}

// OpenAIAudioVoiceUnionFlat 语音类型 union 扁平化实现。
type OpenAIAudioVoiceUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"string"/"id"。
	Type string `json:"type,omitempty"`
	// String 字符串语音名。
	String *OpenAIAudioVoiceUnionString `json:"string,omitempty"`
	// ID 自定义语音 ID 对象。
	ID *OpenAIAudioVoiceUnionID `json:"id,omitempty"`
}

func (OpenAIAudioVoiceUnionFlat) isAudioVoiceUnion() {}

func (u OpenAIAudioVoiceUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "id":
		if u.ID == nil {
			return nil, nil
		}
		return json.Marshal(u.ID)
	}
	return nil, nil
}

func (u *OpenAIAudioVoiceUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIAudioVoiceUnionString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "id":
		var v OpenAIAudioVoiceUnionID
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ID = &v
	}
	return nil
}

// OpenAIAudioVoiceUnionString 字符串语音名。
type OpenAIAudioVoiceUnionString OpenAIAudioVoice

func (OpenAIAudioVoiceUnionString) isAudioVoiceUnion() {}

// OpenAIAudioVoiceUnionID 自定义语音 ID 对象。
type OpenAIAudioVoiceUnionID struct {
	ID string `json:"id"`
}

func (OpenAIAudioVoiceUnionID) isAudioVoiceUnion() {}

// OpenAIChatCompletionPredictionContentType 预测内容类型。
type OpenAIChatCompletionPredictionContentType string

const (
	OpenAIChatCompletionPredictionContentTypeContent OpenAIChatCompletionPredictionContentType = "content"
	OpenAIChatCompletionPredictionContentTypeText    OpenAIChatCompletionPredictionContentType = "text"
)

// OpenAIChatCompletionPredictionContentParam 预测内容参数。
type OpenAIChatCompletionPredictionContentParam struct {
	// Type 固定为 "content"。
	Type OpenAIChatCompletionPredictionContentType `json:"type"`
	// Content 预测的内容（text 数组）。
	Content []OpenAIChatCompletionPredictionContentText `json:"content"`
}

// OpenAIChatCompletionPredictionContentText 预测文本块。
type OpenAIChatCompletionPredictionContentText struct {
	Type OpenAIChatCompletionPredictionContentType `json:"type"`
	Text string                                    `json:"text"`
}

// OpenAIChatCompletionWebSearchOptions Web 搜索选项。
type OpenAIChatCompletionWebSearchOptions struct {
	// SearchContextSize 搜索上下文大小。
	SearchContextSize string `json:"search_context_size,omitempty"`
	// UserLocation 用户位置（用于本地化搜索）。
	UserLocation *OpenAIChatCompletionWebSearchUserLocation `json:"user_location,omitempty"`
	// Filters 搜索结果过滤（域名白名单等）。
	Filters *OpenAIChatCompletionWebSearchFilters `json:"filters,omitempty"`
}

// OpenAIChatCompletionWebSearchFilters Web 搜索过滤。
type OpenAIChatCompletionWebSearchFilters struct {
	// AllowDomains 允许的域名白名单（仅返回来自这些域名的结果）。
	AllowDomains []string `json:"allow_domains,omitempty"`
}

// OpenAIChatCompletionWebSearchUserLocation 用户位置。
type OpenAIChatCompletionWebSearchUserLocation struct {
	Type     string `json:"type"` // 固定 "approximate"
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

// OpenAIResponseFormatUnion 响应格式 union。
type OpenAIResponseFormatUnion interface {
	isResponseFormatUnion()
}

// OpenAIResponseFormatUnionFlat 响应格式 union 扁平化实现。
type OpenAIResponseFormatUnionFlat struct {
	// Type 标识当前激活的变体：""(空)/"text"/"json_object"/"json_schema"。
	Type string `json:"type,omitempty"`
	// Text 纯文本格式。
	Text *OpenAIResponseFormatTextParam `json:"text,omitempty"`
	// JSONObject JSON 对象格式。
	JSONObject *OpenAIResponseFormatJSONObjectParam `json:"json_object,omitempty"`
	// JSONSchema JSON Schema 格式。
	JSONSchema *OpenAIResponseFormatJSONSchemaParam `json:"json_schema,omitempty"`
}

func (OpenAIResponseFormatUnionFlat) isResponseFormatUnion() {}

func (u OpenAIResponseFormatUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "text":
		if u.Text == nil {
			return nil, nil
		}
		return json.Marshal(u.Text)
	case "json_object":
		if u.JSONObject == nil {
			return nil, nil
		}
		return json.Marshal(u.JSONObject)
	case "json_schema":
		if u.JSONSchema == nil {
			return nil, nil
		}
		return json.Marshal(u.JSONSchema)
	}
	return nil, nil
}

func (u *OpenAIResponseFormatUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "text":
		var v OpenAIResponseFormatTextParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Text = &v
	case "json_object":
		var v OpenAIResponseFormatJSONObjectParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.JSONObject = &v
	case "json_schema":
		var v OpenAIResponseFormatJSONSchemaParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.JSONSchema = &v
	}
	return nil
}

// OpenAIResponseFormatTextParam 纯文本格式。
type OpenAIResponseFormatTextParam struct {
	Type OpenAIResponseFormatType `json:"type"` // "text"
}

func (OpenAIResponseFormatTextParam) isResponseFormatUnion() {}

// OpenAIResponseFormatJSONObjectParam JSON 对象格式。
type OpenAIResponseFormatJSONObjectParam struct {
	Type OpenAIResponseFormatType `json:"type"` // "json_object"
}

func (OpenAIResponseFormatJSONObjectParam) isResponseFormatUnion() {}

// OpenAIResponseFormatJSONSchemaParam JSON Schema 格式。
type OpenAIResponseFormatJSONSchemaParam struct {
	Type       OpenAIResponseFormatType `json:"type"` // "json_schema"
	JSONSchema OpenAIJSONSchema         `json:"json_schema"`
}

func (OpenAIResponseFormatJSONSchemaParam) isResponseFormatUnion() {}

// OpenAIJSONSchema JSON Schema 定义。
type OpenAIJSONSchema struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Schema      json.RawMessage `json:"schema"`
	Strict      *bool           `json:"strict,omitempty"`
}
