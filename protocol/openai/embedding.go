package openai

import "encoding/json"

// OpenAIEmbeddingInputUnion 输入 union（string | []string | []int64 | [][]int64）。
// 对齐官方 SDK EmbeddingNewParamsInputUnion 的 4 种形态。
type OpenAIEmbeddingInputUnion interface {
	isEmbeddingInputUnion()
}

// OpenAIEmbeddingInputUnionFlat 输入 union 扁平化实现。
type OpenAIEmbeddingInputUnionFlat struct {
	// Type 标识当前激活的变体：string/string_array/token_array/token_array_array。
	Type string `json:"type,omitempty"`
	// String 单字符串变体。
	String *OpenAIEmbeddingInputString `json:"string,omitempty"`
	// StringArray 多字符串变体。
	StringArray *OpenAIEmbeddingInputStringArray `json:"string_array,omitempty"`
	// TokenArray token id 数组变体。
	TokenArray *OpenAIEmbeddingInputTokenArray `json:"token_array,omitempty"`
	// TokenArrayArray token id 二维数组变体。
	TokenArrayArray *OpenAIEmbeddingInputTokenArrayArray `json:"token_array_array,omitempty"`
}

func (OpenAIEmbeddingInputUnionFlat) isEmbeddingInputUnion() {}

func (u OpenAIEmbeddingInputUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.String == nil {
			return nil, nil
		}
		return json.Marshal(u.String)
	case "string_array":
		if u.StringArray == nil {
			return nil, nil
		}
		return json.Marshal(u.StringArray)
	case "token_array":
		if u.TokenArray == nil {
			return nil, nil
		}
		return json.Marshal(u.TokenArray)
	case "token_array_array":
		if u.TokenArrayArray == nil {
			return nil, nil
		}
		return json.Marshal(u.TokenArrayArray)
	}
	return nil, nil
}

func (u *OpenAIEmbeddingInputUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIEmbeddingInputString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "string_array":
		var v OpenAIEmbeddingInputStringArray
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.StringArray = &v
	case "token_array":
		var v OpenAIEmbeddingInputTokenArray
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.TokenArray = &v
	case "token_array_array":
		var v OpenAIEmbeddingInputTokenArrayArray
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.TokenArrayArray = &v
	}
	return nil
}

// OpenAIEmbeddingInputString 单字符串输入。
type OpenAIEmbeddingInputString string

func (OpenAIEmbeddingInputString) isEmbeddingInputUnion() {}

// OpenAIEmbeddingInputStringArray 多字符串输入。
type OpenAIEmbeddingInputStringArray []string

func (OpenAIEmbeddingInputStringArray) isEmbeddingInputUnion() {}

// OpenAIEmbeddingInputTokenArray token id 数组输入。
type OpenAIEmbeddingInputTokenArray []int64

func (OpenAIEmbeddingInputTokenArray) isEmbeddingInputUnion() {}

// OpenAIEmbeddingInputTokenArrayArray token id 二维数组输入。
type OpenAIEmbeddingInputTokenArrayArray [][]int64

func (OpenAIEmbeddingInputTokenArrayArray) isEmbeddingInputUnion() {}

// OpenAIEmbeddingParams 嵌入请求参数。
type OpenAIEmbeddingParams struct {
	// Input 输入内容（string | []string | []int64 | [][]int64 四种形态之一）。
	Input *OpenAIEmbeddingInputUnionFlat `json:"input"`
	// Model 模型 ID（如 "text-embedding-3-small"）。
	Model string `json:"model"`
	// EncodingFormat 编码格式（float / base64）。
	EncodingFormat OpenAIEmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	// Dimensions 维度（可选，覆盖默认）。
	Dimensions *int64 `json:"dimensions,omitempty"`
	// User 用户标识。
	User string `json:"user,omitempty"`
}

// OpenAIEmbeddingResponse 嵌入响应。
type OpenAIEmbeddingResponse struct {
	// Object 固定 "list"。
	Object string `json:"object"`
	// Data 嵌入数据列表。
	Data []OpenAIEmbedding `json:"data"`
	// Model 模型。
	Model string `json:"model"`
	// Usage 用量。
	Usage OpenAIEmbeddingUsage `json:"usage"`
}

// OpenAIEmbedding 单个嵌入。
type OpenAIEmbedding struct {
	// Object 固定 "embedding"。
	Object string `json:"object"`
	// Index 索引。
	Index int64 `json:"index"`
	// Embedding 嵌入向量。
	Embedding []float64 `json:"embedding"`
}

// OpenAIEmbeddingUsage 用量。
type OpenAIEmbeddingUsage struct {
	// PromptTokens 输入 token 数。
	PromptTokens int64 `json:"prompt_tokens"`
	// TotalTokens 总 token 数。
	TotalTokens int64 `json:"total_tokens"`
}
