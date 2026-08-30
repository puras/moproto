package openai

// OpenAIServiceTier 服务等级（ChatCompletion / Responses 顶层 service_tier 字段）。
//
// 官方枚举："auto" | "default" | "flex" | "scale" | "priority"。
type OpenAIServiceTier string

const (
	OpenAIServiceTierAuto     OpenAIServiceTier = "auto"
	OpenAIServiceTierDefault  OpenAIServiceTier = "default"
	OpenAIServiceTierPriority OpenAIServiceTier = "priority"
	OpenAIServiceTierFlex     OpenAIServiceTier = "flex"
	OpenAIServiceTierScale    OpenAIServiceTier = "scale"
)

// OpenAIReasoningEffort 推理强度（ChatCompletion reasoning_effort / Responses reasoning.effort）。
//
// 官方枚举："none" | "minimal" | "low" | "medium" | "high" | "xhigh"。
//   - `gpt-5.1` 默认 `none`（不执行推理）。
//   - `xhigh` 仅 `gpt-5.1-codex-max` 及之后模型支持。
type OpenAIReasoningEffort string

const (
	OpenAIReasoningEffortNone    OpenAIReasoningEffort = "none"
	OpenAIReasoningEffortMinimal OpenAIReasoningEffort = "minimal"
	OpenAIReasoningEffortLow     OpenAIReasoningEffort = "low"
	OpenAIReasoningEffortMedium  OpenAIReasoningEffort = "medium"
	OpenAIReasoningEffortHigh    OpenAIReasoningEffort = "high"
	OpenAIReasoningEffortXHigh   OpenAIReasoningEffort = "xhigh"
)

// OpenAIVerbosity 输出详细度（ChatCompletion verbosity / Responses text.verbosity）。
type OpenAIVerbosity string

const (
	OpenAIVerbosityLow    OpenAIVerbosity = "low"
	OpenAIVerbosityMedium OpenAIVerbosity = "medium"
	OpenAIVerbosityHigh   OpenAIVerbosity = "high"
)

// OpenAIResponseTruncation 截断策略（Responses truncation 字段）。
type OpenAIResponseTruncation string

const (
	OpenAIResponseTruncationAuto     OpenAIResponseTruncation = "auto"
	OpenAIResponseTruncationDisabled OpenAIResponseTruncation = "disabled"
)

// OpenAIResponsePromptCacheRetention 缓存保留策略（Responses prompt_cache_retention）。
type OpenAIResponsePromptCacheRetention string

const (
	OpenAIResponsePromptCacheRetentionInMemory OpenAIResponsePromptCacheRetention = "in-memory"
	OpenAIResponsePromptCacheRetention24h      OpenAIResponsePromptCacheRetention = "24h"
)

// OpenAIAudioFormat 音频格式。
type OpenAIAudioFormat string

const (
	OpenAIAudioFormatWAV  OpenAIAudioFormat = "wav"
	OpenAIAudioFormatMP3  OpenAIAudioFormat = "mp3"
	OpenAIAudioFormatFLAC OpenAIAudioFormat = "flac"
	OpenAIAudioFormatOpus OpenAIAudioFormat = "opus"
	OpenAIAudioFormatPCM  OpenAIAudioFormat = "pcm16"
)

// OpenAIAudioVoice 语音类型（字符串形式）。
type OpenAIAudioVoice string

const (
	OpenAIAudioVoiceAlloy   OpenAIAudioVoice = "alloy"
	OpenAIAudioVoiceEcho    OpenAIAudioVoice = "echo"
	OpenAIAudioVoiceFable   OpenAIAudioVoice = "fable"
	OpenAIAudioVoiceOnyx    OpenAIAudioVoice = "onyx"
	OpenAIAudioVoiceNova    OpenAIAudioVoice = "nova"
	OpenAIAudioVoiceShimmer OpenAIAudioVoice = "shimmer"
	OpenAIAudioVoiceAsh     OpenAIAudioVoice = "ash"
	OpenAIAudioVoiceBallad  OpenAIAudioVoice = "ballad"
	OpenAIAudioVoiceCoral   OpenAIAudioVoice = "coral"
	OpenAIAudioVoiceSage    OpenAIAudioVoice = "sage"
	OpenAIAudioVoiceVerse   OpenAIAudioVoice = "verse"
	OpenAIAudioVoiceMarin   OpenAIAudioVoice = "marin"
	OpenAIAudioVoiceCedar   OpenAIAudioVoice = "cedar"
)

// OpenAIModalities 输出模态（ChatCompletion modalities 字段）。
type OpenAIModality string

const (
	OpenAIModalityText  OpenAIModality = "text"
	OpenAIModalityAudio OpenAIModality = "audio"
)

// OpenAIImageDetail 图片细节级别。
type OpenAIImageDetail string

const (
	OpenAIImageDetailAuto OpenAIImageDetail = "auto"
	OpenAIImageDetailLow  OpenAIImageDetail = "low"
	OpenAIImageDetailHigh OpenAIImageDetail = "high"
)

// OpenAIResponseFormatType 响应格式类型（ChatCompletion response_format / Responses text.format）。
type OpenAIResponseFormatType string

const (
	OpenAIResponseFormatTypeText       OpenAIResponseFormatType = "text"
	OpenAIResponseFormatTypeJSONObject OpenAIResponseFormatType = "json_object"
	OpenAIResponseFormatTypeJSONSchema OpenAIResponseFormatType = "json_schema"
)

// OpenAIEmbeddingModel 嵌入模型 ID。
type OpenAIEmbeddingModel string

const (
	OpenAIEmbeddingModelTextEmbeddingAda002 OpenAIEmbeddingModel = "text-embedding-ada-002"
	OpenAIEmbeddingModelTextEmbedding3Small OpenAIEmbeddingModel = "text-embedding-3-small"
	OpenAIEmbeddingModelTextEmbedding3Large OpenAIEmbeddingModel = "text-embedding-3-large"
)

// OpenAIEmbeddingEncodingFormat 嵌入编码格式。
type OpenAIEmbeddingEncodingFormat string

const (
	OpenAIEmbeddingEncodingFormatFloat  OpenAIEmbeddingEncodingFormat = "float"
	OpenAIEmbeddingEncodingFormatBase64 OpenAIEmbeddingEncodingFormat = "base64"
)

// OpenAIModerationModel 内容审核模型 ID。
type OpenAIModerationModel string

const (
	OpenAIModerationModelOmniModerationLatest     OpenAIModerationModel = "omni-moderation-latest"
	OpenAIModerationModelOmniModeration2024_09_26 OpenAIModerationModel = "omni-moderation-2024-09-26"
	OpenAIModerationModelTextModerationLatest     OpenAIModerationModel = "text-moderation-latest"
	OpenAIModerationModelTextModerationStable     OpenAIModerationModel = "text-moderation-stable"
)

// OpenAIModerationMultiModalInputType 多模态输入类型（image_url / text）。
type OpenAIModerationMultiModalInputType string

const (
	OpenAIModerationMultiModalInputTypeImageURL OpenAIModerationMultiModalInputType = "image_url"
	OpenAIModerationMultiModalInputTypeText     OpenAIModerationMultiModalInputType = "text"
)
