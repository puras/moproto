package gemini

// GeminiSafetySetting 顶层 safetySettings 数组元素。
type GeminiSafetySetting struct {
	Category  GeminiHarmCategory       `json:"category,omitempty"`
	Threshold GeminiHarmBlockThreshold `json:"threshold,omitempty"`
}

// GeminiSafetyRating 安全评级（candidate.safetyRatings / promptFeedback.safetyRatings）。
type GeminiSafetyRating struct {
	Category    GeminiHarmCategory    `json:"category,omitempty"`
	Probability GeminiHarmProbability `json:"probability,omitempty"`
	Blocked     bool                  `json:"blocked,omitempty"`
}
