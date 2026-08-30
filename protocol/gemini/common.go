package gemini

// GeminiRole 消息角色类型。
//
// 取值空间：user | model。system 角色在 Gemini 协议中走顶层 systemInstruction
// 字段，不在 Content.Role 中出现。
type GeminiRole string

const (
	GeminiRoleUser  GeminiRole = "user"
	GeminiRoleModel GeminiRole = "model"
)

// GeminiFinishReason candidate.finishReason 字段取值。
type GeminiFinishReason string

const (
	GeminiFinishReasonStop       GeminiFinishReason = "STOP"
	GeminiFinishReasonMaxTokens  GeminiFinishReason = "MAX_TOKENS"
	GeminiFinishReasonSafety     GeminiFinishReason = "SAFETY"
	GeminiFinishReasonRecitation GeminiFinishReason = "RECITATION"
	GeminiFinishReasonOther      GeminiFinishReason = "OTHER"
	GeminiFinishReasonBlocklist  GeminiFinishReason = "BLOCKLIST"
	GeminiFinishReasonProhibited GeminiFinishReason = "PROHIBITED_CONTENT"
	GeminiFinishReasonSpii       GeminiFinishReason = "SPII"
	GeminiFinishReasonMalformed  GeminiFinishReason = "MALFORMED_FUNCTION_CALL"
)

// GeminiHarmCategory SafetyCategory 字段枚举。
type GeminiHarmCategory string

const (
	GeminiHarmCategoryHateSpeech       GeminiHarmCategory = "HARM_CATEGORY_HATE_SPEECH"
	GeminiHarmCategoryHarassment       GeminiHarmCategory = "HARM_CATEGORY_HARASSMENT"
	GeminiHarmCategorySexuallyExplicit GeminiHarmCategory = "HARM_CATEGORY_SEXUALLY_EXPLICIT"
	GeminiHarmCategoryDangerousContent GeminiHarmCategory = "HARM_CATEGORY_DANGEROUS_CONTENT"
	GeminiHarmCategoryCivicIntegrity   GeminiHarmCategory = "HARM_CATEGORY_CIVIC_INTEGRITY"
)

// GeminiHarmBlockThreshold SafetySetting.threshold 字段取值。
type GeminiHarmBlockThreshold string

const (
	GeminiHarmBlockThresholdUnspecified GeminiHarmBlockThreshold = "HARM_BLOCK_THRESHOLD_UNSPECIFIED"
	GeminiHarmBlockThresholdLow         GeminiHarmBlockThreshold = "BLOCK_LOW_AND_ABOVE"
	GeminiHarmBlockThresholdMedium      GeminiHarmBlockThreshold = "BLOCK_MEDIUM_AND_ABOVE"
	GeminiHarmBlockThresholdHigh        GeminiHarmBlockThreshold = "BLOCK_ONLY_HIGH"
	GeminiHarmBlockThresholdNone        GeminiHarmBlockThreshold = "BLOCK_NONE"
)

// GeminiHarmProbability SafetyRating.probability 字段取值。
type GeminiHarmProbability string

const (
	GeminiHarmProbabilityUnspecified GeminiHarmProbability = "HARM_PROBABILITY_UNSPECIFIED"
	GeminiHarmProbabilityNegligible  GeminiHarmProbability = "NEGLIGIBLE"
	GeminiHarmProbabilityLow         GeminiHarmProbability = "LOW"
	GeminiHarmProbabilityMedium      GeminiHarmProbability = "MEDIUM"
	GeminiHarmProbabilityHigh        GeminiHarmProbability = "HIGH"
)

// GeminiFunctionCallingMode functionCallingConfig.mode 字段取值。
type GeminiFunctionCallingMode string

const (
	GeminiFunctionCallingModeAuto GeminiFunctionCallingMode = "AUTO"
	GeminiFunctionCallingModeAny  GeminiFunctionCallingMode = "ANY"
	GeminiFunctionCallingModeNone GeminiFunctionCallingMode = "NONE"
)
