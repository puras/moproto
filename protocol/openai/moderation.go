package openai

import "encoding/json"

// === Moderation Input Union 类型 ===

// OpenAIModerationInputUnion 审核输入 union（string | []string | []MultiModalInput）。
// 对齐官方 SDK ModerationNewParamsInputUnion 的 3 种形态。
type OpenAIModerationInputUnion interface {
	isModerationInputUnion()
}

// OpenAIModerationInputUnionFlat 审核输入 union 扁平化实现。
type OpenAIModerationInputUnionFlat struct {
	// Type 标识当前激活的变体：string/string_array/multi_modal_array。
	Type string `json:"type,omitempty"`
	// String 单字符串变体。
	String *OpenAIModerationInputString `json:"string,omitempty"`
	// StringArray 多字符串变体。
	StringArray *OpenAIModerationInputStringArray `json:"string_array,omitempty"`
	// MultiModalArray 多模态输入项列表变体。
	MultiModalArray *OpenAIModerationInputMultiModalArray `json:"multi_modal_array,omitempty"`
}

func (OpenAIModerationInputUnionFlat) isModerationInputUnion() {}

func (u OpenAIModerationInputUnionFlat) MarshalJSON() ([]byte, error) {
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
	case "multi_modal_array":
		if u.MultiModalArray == nil {
			return nil, nil
		}
		return json.Marshal(u.MultiModalArray)
	}
	return nil, nil
}

func (u *OpenAIModerationInputUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		var v OpenAIModerationInputString
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.String = &v
	case "string_array":
		var v OpenAIModerationInputStringArray
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.StringArray = &v
	case "multi_modal_array":
		var v OpenAIModerationInputMultiModalArray
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.MultiModalArray = &v
	}
	return nil
}

// OpenAIModerationInputString 单字符串输入。
type OpenAIModerationInputString string

func (OpenAIModerationInputString) isModerationInputUnion() {}

// OpenAIModerationInputStringArray 多字符串输入。
type OpenAIModerationInputStringArray []string

func (OpenAIModerationInputStringArray) isModerationInputUnion() {}

// OpenAIModerationInputMultiModalArray 多模态输入项列表。
type OpenAIModerationInputMultiModalArray []OpenAIModerationMultiModalInputUnionParam

func (OpenAIModerationInputMultiModalArray) isModerationInputUnion() {}

// === Moderation MultiModal Input Union 嵌套类型 ===

// OpenAIModerationMultiModalInputUnionParam 多模态输入项 union（image_url | text）。
type OpenAIModerationMultiModalInputUnionParam interface {
	isModerationMultiModalInput()
}

// OpenAIModerationMultiModalInputUnionParamFlat 多模态输入项 union 扁平化实现。
type OpenAIModerationMultiModalInputUnionParamFlat struct {
	// Type 标识当前激活的变体：image_url/text。
	Type string `json:"type,omitempty"`
	// ImageURL 图片 URL 变体。
	ImageURL *OpenAIModerationMultiModalInputImageURLParam `json:"image_url,omitempty"`
	// Text text 变体。
	Text *OpenAIModerationMultiModalInputTextParam `json:"text,omitempty"`
}

func (OpenAIModerationMultiModalInputUnionParamFlat) isModerationMultiModalInput() {}

func (u OpenAIModerationMultiModalInputUnionParamFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "image_url":
		if u.ImageURL == nil {
			return nil, nil
		}
		return json.Marshal(u.ImageURL)
	case "text":
		if u.Text == nil {
			return nil, nil
		}
		return json.Marshal(u.Text)
	}
	return nil, nil
}

func (u *OpenAIModerationMultiModalInputUnionParamFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "image_url":
		var v OpenAIModerationMultiModalInputImageURLParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ImageURL = &v
	case "text":
		var v OpenAIModerationMultiModalInputTextParam
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Text = &v
	}
	return nil
}

// OpenAIModerationImageURL 图片 URL（嵌套结构）。
type OpenAIModerationImageURL struct {
	// URL 图片 URL 或 base64 data URL。
	URL string `json:"url"`
}

// OpenAIModerationMultiModalInputImageURLParam 多模态 image_url 输入项。
type OpenAIModerationMultiModalInputImageURLParam struct {
	// Type 固定 "image_url"。
	Type OpenAIModerationMultiModalInputType `json:"type"`
	// ImageURL 图片 URL 详情。
	ImageURL OpenAIModerationImageURL `json:"image_url"`
}

func (OpenAIModerationMultiModalInputImageURLParam) isModerationMultiModalInput() {}

// OpenAIModerationMultiModalInputTextParam 多模态 text 输入项。
type OpenAIModerationMultiModalInputTextParam struct {
	// Type 固定 "text"。
	Type OpenAIModerationMultiModalInputType `json:"type"`
	// Text 待审核文本。
	Text string `json:"text"`
}

func (OpenAIModerationMultiModalInputTextParam) isModerationMultiModalInput() {}

// === Moderation Params + Response ===

// OpenAIModerationParams 内容审核请求参数。
type OpenAIModerationParams struct {
	// Input 输入内容（string | []string | []MultiModalInput）。
	Input *OpenAIModerationInputUnionFlat `json:"input"`
	// Model 模型 ID（如 "omni-moderation-latest" / "text-moderation-latest"）。
	Model string `json:"model,omitempty"`
}

// OpenAIModerationResponse 审核响应。
type OpenAIModerationResponse struct {
	// ID 响应 ID。
	ID string `json:"id"`
	// Model 模型。
	Model string `json:"model"`
	// Results 审核结果列表。
	Results []OpenAIModeration `json:"results"`
}

// OpenAIModeration 单个审核结果。
type OpenAIModeration struct {
	// Flagged 是否被标记。
	Flagged bool `json:"flagged"`
	// Categories 各类别标记。
	Categories OpenAIModerationCategories `json:"categories"`
	// CategoryScores 各类别分数。
	CategoryScores OpenAIModerationCategoryScores `json:"category_scores"`
	// AppliedInputRiskLevel 应用的输入风险等级（omni-moderation）。
	AppliedInputRiskLevel string `json:"applied_input_risk_level,omitempty"`
}

// OpenAIModerationCategories 各类别标记。
type OpenAIModerationCategories struct {
	Hate                  bool `json:"hate,omitempty"`
	HateThreatening       bool `json:"hate/threatening,omitempty"`
	Harassment            bool `json:"harassment,omitempty"`
	HarassmentThreatening bool `json:"harassment/threatening,omitempty"`
	SelfHarm              bool `json:"self-harm,omitempty"`
	SelfHarmIntent        bool `json:"self-harm/intent,omitempty"`
	SelfHarmInstructions  bool `json:"self-harm/instructions,omitempty"`
	Sexual                bool `json:"sexual,omitempty"`
	SexualMinors          bool `json:"sexual/minors,omitempty"`
	Violence              bool `json:"violence,omitempty"`
	ViolenceGraphic       bool `json:"violence/graphic,omitempty"`
	Illicit               bool `json:"illicit,omitempty"`
	IllicitViolent        bool `json:"illicit/violent,omitempty"`
}

// OpenAIModerationCategoryScores 各类别分数。
type OpenAIModerationCategoryScores struct {
	Hate                  float64 `json:"hate,omitempty"`
	HateThreatening       float64 `json:"hate/threatening,omitempty"`
	Harassment            float64 `json:"harassment,omitempty"`
	HarassmentThreatening float64 `json:"harassment/threatening,omitempty"`
	SelfHarm              float64 `json:"self-harm,omitempty"`
	SelfHarmIntent        float64 `json:"self-harm/intent,omitempty"`
	SelfHarmInstructions  float64 `json:"self-harm/instructions,omitempty"`
	Sexual                float64 `json:"sexual,omitempty"`
	SexualMinors          float64 `json:"sexual/minors,omitempty"`
	Violence              float64 `json:"violence,omitempty"`
	ViolenceGraphic       float64 `json:"violence/graphic,omitempty"`
	Illicit               float64 `json:"illicit,omitempty"`
	IllicitViolent        float64 `json:"illicit/violent,omitempty"`
}
