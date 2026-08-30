package openai

import (
	"encoding/json"
	"strings"
)

// OpenAIResponseStreamEventUnion Responses 流事件 union 接口。
//
// 主线 53 变体，对齐官方 SDK ResponseStreamEventUnion 全部实现集。
type OpenAIResponseStreamEventUnion interface {
	isResponseStreamEvent()
}

// OpenAIResponseStreamEventType 流事件 type 字段类型。
type OpenAIResponseStreamEventType string

// OpenAIResponseStreamEventUnionFlat Responses 流事件 union 扁平化实现（53 变体）。
type OpenAIResponseStreamEventUnionFlat struct {
	// Type 区分变体。
	Type string `json:"type,omitempty"`

	// === 公共字段 ===
	Sequence     int64 `json:"sequence_number,omitempty"`
	OutputIndex  int64 `json:"output_index,omitempty"`
	ContentIndex int64 `json:"content_index,omitempty"`
	SummaryIndex int64 `json:"summary_index,omitempty"`

	// === response.created / response.queued / response.in_progress / response.completed / response.failed / response.incomplete ===
	Response *OpenAIResponse `json:"response,omitempty"`

	// === response.refusal.done ===
	Refusal       string `json:"refusal,omitempty"`
	RefusalItemID string `json:"item_id,omitempty"`

	// === error ===
	Error *OpenAIResponseError `json:"error,omitempty"`

	// === response.output_item.added / response.output_item.done ===
	ItemID2 string                        `json:"item_id,omitempty"`
	Item2   OpenAIResponseOutputItemUnion `json:"item,omitempty"`

	// === response.output_text.* / response.refusal.delta ===
	TextDelta string `json:"delta,omitempty"`

	// === response.output_text.done ===
	Text string `json:"text,omitempty"`

	// === response.content_part.* ===
	Part *OpenAIResponseOutputMessageContent `json:"part,omitempty"`

	// === response.reasoning_text.* ===
	ReasoningDelta string `json:"delta,omitempty"`

	// === response.reasoning_text.done ===
	ReasoningText string `json:"text,omitempty"`

	// === response.reasoning_summary_part.* ===
	Part2 *OpenAIResponseOutputReasoningSummary `json:"part,omitempty"`

	// === response.reasoning_summary_text.* ===
	SummaryDelta string `json:"delta,omitempty"`

	// === response.reasoning_summary_text.done ===
	SummaryText string `json:"text,omitempty"`

	// === response.function_call_arguments.* ===
	FuncDelta string `json:"delta,omitempty"`
	FuncArgs  string `json:"arguments,omitempty"`

	// === response.file_search_call.* / response.web_search_call.* ===
	SearchStatus string `json:"status,omitempty"`

	// === response.mcp_call.* ===
	McpDelta  string `json:"delta,omitempty"`
	McpArgs   string `json:"arguments,omitempty"`
	McpStatus string `json:"status,omitempty"`

	// === response.code_interpreter_call.* ===
	CodeDelta  string `json:"delta,omitempty"`
	Code       string `json:"code,omitempty"`
	CodeStatus string `json:"status,omitempty"`

	// === response.image_generation_call.* ===
	ImgStatus       string `json:"status,omitempty"`
	PartialImageIdx int64  `json:"partial_image_index,omitempty"`
	PartialImageB64 string `json:"b64_json,omitempty"`

	// === response.audio.delta / response.audio_transcript.delta ===
	AudioDelta string `json:"delta,omitempty"`

	// === response.audio.done ===
	Audio string `json:"audio,omitempty"`

	// === response.audio_transcript.done ===
	Transcript string `json:"transcript,omitempty"`

	// === response.custom_tool_call_input.* ===
	CustomDelta string `json:"delta,omitempty"`

	// === response.custom_tool_call_input.done ===
	CustomInput string `json:"input,omitempty"`

	// === response.output_text.annotation.added ===
	AnnotationIndex int64           `json:"annotation_index,omitempty"`
	Annotation      json.RawMessage `json:"annotation,omitempty"`
}

func (OpenAIResponseStreamEventUnionFlat) isResponseStreamEvent() {}

func (u OpenAIResponseStreamEventUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	// Lifecycle (6)
	case "response.created", "response.queued", "response.in_progress", "response.completed", "response.failed", "response.incomplete":
		return json.Marshal(struct {
			Type     string          `json:"type"`
			Sequence int64           `json:"sequence_number"`
			Response *OpenAIResponse `json:"response,omitempty"`
		}{
			Type:     u.Type,
			Sequence: u.Sequence,
			Response: u.Response,
		})
	// Error
	case "error":
		return json.Marshal(struct {
			Type     string               `json:"type"`
			Sequence int64                `json:"sequence_number"`
			Error    *OpenAIResponseError `json:"error,omitempty"`
		}{
			Type:     "error",
			Sequence: u.Sequence,
			Error:    u.Error,
		})
	// Output Text (4)
	case "response.output_text.delta":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Delta        string `json:"delta"`
		}{
			Type:         "response.output_text.delta",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Delta:        u.TextDelta,
		})
	case "response.output_text.done":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Text         string `json:"text"`
		}{
			Type:         "response.output_text.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Text:         u.Text,
		})
	// Refusal (2)
	case "response.refusal.delta":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Delta        string `json:"delta"`
		}{
			Type:         "response.refusal.delta",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Delta:        u.TextDelta,
		})
	case "response.refusal.done":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Refusal      string `json:"refusal"`
		}{
			Type:         "response.refusal.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Refusal:      u.Refusal,
		})
	// Audio (4)
	case "response.audio.delta":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Delta        string `json:"delta"`
		}{
			Type:         "response.audio.delta",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Delta:        u.AudioDelta,
		})
	case "response.audio.done":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Audio        string `json:"audio"`
		}{
			Type:         "response.audio.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Audio:        u.Audio,
		})
	case "response.audio_transcript.delta":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Delta        string `json:"delta"`
		}{
			Type:         "response.audio_transcript.delta",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Delta:        u.AudioDelta,
		})
	case "response.audio_transcript.done":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Transcript   string `json:"transcript"`
		}{
			Type:         "response.audio_transcript.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Transcript:   u.Transcript,
		})
	// Function Call (2)
	case "response.function_call_arguments.delta":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			ItemID      string `json:"item_id"`
			OutputIndex int64  `json:"output_index"`
			Delta       string `json:"delta"`
		}{
			Type:        "response.function_call_arguments.delta",
			Sequence:    u.Sequence,
			ItemID:      u.RefusalItemID,
			OutputIndex: u.OutputIndex,
			Delta:       u.FuncDelta,
		})
	case "response.function_call_arguments.done":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			ItemID      string `json:"item_id"`
			OutputIndex int64  `json:"output_index"`
			Arguments   string `json:"arguments"`
		}{
			Type:        "response.function_call_arguments.done",
			Sequence:    u.Sequence,
			ItemID:      u.RefusalItemID,
			OutputIndex: u.OutputIndex,
			Arguments:   u.FuncArgs,
		})
	// Code Interpreter (5)
	case "response.code_interpreter_call_code.delta":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
			Delta       string `json:"delta"`
		}{
			Type:        "response.code_interpreter_call_code.delta",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
			Delta:       u.CodeDelta,
		})
	case "response.code_interpreter_call_code.done":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
			Code        string `json:"code"`
		}{
			Type:        "response.code_interpreter_call_code.done",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
			Code:        u.Code,
		})
	case "response.code_interpreter_call.completed":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.code_interpreter_call.completed",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.code_interpreter_call.in_progress":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.code_interpreter_call.in_progress",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.code_interpreter_call.interpreting":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.code_interpreter_call.interpreting",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	// File Search (3)
	case "response.file_search_call.completed":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.file_search_call.completed",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.file_search_call.in_progress":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.file_search_call.in_progress",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.file_search_call.searching":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.file_search_call.searching",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	// Web Search (3)
	case "response.web_search_call.completed":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.web_search_call.completed",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.web_search_call.in_progress":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.web_search_call.in_progress",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.web_search_call.searching":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.web_search_call.searching",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	// Image Gen (4)
	case "response.image_generation_call.completed":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.image_generation_call.completed",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.image_generation_call.generating":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.image_generation_call.generating",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.image_generation_call.in_progress":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.image_generation_call.in_progress",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.image_generation_call.partial_image":
		return json.Marshal(struct {
			Type              string `json:"type"`
			Sequence          int64  `json:"sequence_number"`
			OutputIndex       int64  `json:"output_index"`
			ItemID            string `json:"item_id"`
			PartialImageIndex int64  `json:"partial_image_index"`
			B64JSON           string `json:"b64_json"`
		}{
			Type:              "response.image_generation_call.partial_image",
			Sequence:          u.Sequence,
			OutputIndex:       u.OutputIndex,
			ItemID:            u.RefusalItemID,
			PartialImageIndex: u.PartialImageIdx,
			B64JSON:           u.PartialImageB64,
		})
	// MCP (8)
	case "response.mcp_call_arguments.delta":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
			Delta       string `json:"delta"`
		}{
			Type:        "response.mcp_call_arguments.delta",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
			Delta:       u.McpDelta,
		})
	case "response.mcp_call_arguments.done":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
			Arguments   string `json:"arguments"`
		}{
			Type:        "response.mcp_call_arguments.done",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
			Arguments:   u.McpArgs,
		})
	case "response.mcp_call.completed":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.mcp_call.completed",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.mcp_call.failed":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.mcp_call.failed",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.mcp_call.in_progress":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			OutputIndex int64  `json:"output_index"`
			ItemID      string `json:"item_id"`
		}{
			Type:        "response.mcp_call.in_progress",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			ItemID:      u.RefusalItemID,
		})
	case "response.mcp_list_tools.completed":
		return json.Marshal(struct {
			Type     string `json:"type"`
			Sequence int64  `json:"sequence_number"`
			ItemID   string `json:"item_id"`
		}{
			Type:     "response.mcp_list_tools.completed",
			Sequence: u.Sequence,
			ItemID:   u.RefusalItemID,
		})
	case "response.mcp_list_tools.failed":
		return json.Marshal(struct {
			Type     string `json:"type"`
			Sequence int64  `json:"sequence_number"`
			ItemID   string `json:"item_id"`
		}{
			Type:     "response.mcp_list_tools.failed",
			Sequence: u.Sequence,
			ItemID:   u.RefusalItemID,
		})
	case "response.mcp_list_tools.in_progress":
		return json.Marshal(struct {
			Type     string `json:"type"`
			Sequence int64  `json:"sequence_number"`
			ItemID   string `json:"item_id"`
		}{
			Type:     "response.mcp_list_tools.in_progress",
			Sequence: u.Sequence,
			ItemID:   u.RefusalItemID,
		})
	// Output Item (2)
	case "response.output_item.added":
		return json.Marshal(struct {
			Type        string                        `json:"type"`
			Sequence    int64                         `json:"sequence_number"`
			OutputIndex int64                         `json:"output_index"`
			Item        OpenAIResponseOutputItemUnion `json:"item"`
		}{
			Type:        "response.output_item.added",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			Item:        u.Item2,
		})
	case "response.output_item.done":
		return json.Marshal(struct {
			Type        string                        `json:"type"`
			Sequence    int64                         `json:"sequence_number"`
			OutputIndex int64                         `json:"output_index"`
			Item        OpenAIResponseOutputItemUnion `json:"item"`
		}{
			Type:        "response.output_item.done",
			Sequence:    u.Sequence,
			OutputIndex: u.OutputIndex,
			Item:        u.Item2,
		})
	// Content Part (2)
	case "response.content_part.added":
		return json.Marshal(struct {
			Type         string                              `json:"type"`
			Sequence     int64                               `json:"sequence_number"`
			ItemID       string                              `json:"item_id"`
			OutputIndex  int64                               `json:"output_index"`
			ContentIndex int64                               `json:"content_index"`
			Part         *OpenAIResponseOutputMessageContent `json:"part"`
		}{
			Type:         "response.content_part.added",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Part:         u.Part,
		})
	case "response.content_part.done":
		return json.Marshal(struct {
			Type         string                              `json:"type"`
			Sequence     int64                               `json:"sequence_number"`
			ItemID       string                              `json:"item_id"`
			OutputIndex  int64                               `json:"output_index"`
			ContentIndex int64                               `json:"content_index"`
			Part         *OpenAIResponseOutputMessageContent `json:"part"`
		}{
			Type:         "response.content_part.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Part:         u.Part,
		})
	// Reasoning (6)
	case "response.reasoning_summary_part.added":
		return json.Marshal(struct {
			Type         string                                `json:"type"`
			Sequence     int64                                 `json:"sequence_number"`
			ItemID       string                                `json:"item_id"`
			OutputIndex  int64                                 `json:"output_index"`
			SummaryIndex int64                                 `json:"summary_index"`
			Part         *OpenAIResponseOutputReasoningSummary `json:"part"`
		}{
			Type:         "response.reasoning_summary_part.added",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			SummaryIndex: u.SummaryIndex,
			Part:         u.Part2,
		})
	case "response.reasoning_summary_part.done":
		return json.Marshal(struct {
			Type         string                                `json:"type"`
			Sequence     int64                                 `json:"sequence_number"`
			ItemID       string                                `json:"item_id"`
			OutputIndex  int64                                 `json:"output_index"`
			SummaryIndex int64                                 `json:"summary_index"`
			Part         *OpenAIResponseOutputReasoningSummary `json:"part"`
		}{
			Type:         "response.reasoning_summary_part.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			SummaryIndex: u.SummaryIndex,
			Part:         u.Part2,
		})
	case "response.reasoning_summary_text.delta":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			SummaryIndex int64  `json:"summary_index"`
			Delta        string `json:"delta"`
		}{
			Type:         "response.reasoning_summary_text.delta",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			SummaryIndex: u.SummaryIndex,
			Delta:        u.SummaryDelta,
		})
	case "response.reasoning_summary_text.done":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			SummaryIndex int64  `json:"summary_index"`
			Text         string `json:"text"`
		}{
			Type:         "response.reasoning_summary_text.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			SummaryIndex: u.SummaryIndex,
			Text:         u.SummaryText,
		})
	case "response.reasoning_text.delta":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Delta        string `json:"delta"`
		}{
			Type:         "response.reasoning_text.delta",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Delta:        u.ReasoningDelta,
		})
	case "response.reasoning_text.done":
		return json.Marshal(struct {
			Type         string `json:"type"`
			Sequence     int64  `json:"sequence_number"`
			ItemID       string `json:"item_id"`
			OutputIndex  int64  `json:"output_index"`
			ContentIndex int64  `json:"content_index"`
			Text         string `json:"text"`
		}{
			Type:         "response.reasoning_text.done",
			Sequence:     u.Sequence,
			ItemID:       u.RefusalItemID,
			OutputIndex:  u.OutputIndex,
			ContentIndex: u.ContentIndex,
			Text:         u.ReasoningText,
		})
	// Annotation (1)
	case "response.output_text.annotation.added":
		return json.Marshal(struct {
			Type            string          `json:"type"`
			Sequence        int64           `json:"sequence_number"`
			ItemID          string          `json:"item_id"`
			OutputIndex     int64           `json:"output_index"`
			ContentIndex    int64           `json:"content_index"`
			AnnotationIndex int64           `json:"annotation_index"`
			Annotation      json.RawMessage `json:"annotation"`
		}{
			Type:            "response.output_text.annotation.added",
			Sequence:        u.Sequence,
			ItemID:          u.RefusalItemID,
			OutputIndex:     u.OutputIndex,
			ContentIndex:    u.ContentIndex,
			AnnotationIndex: u.AnnotationIndex,
			Annotation:      u.Annotation,
		})
	// Custom Tool (2)
	case "response.custom_tool_call_input.delta":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			ItemID      string `json:"item_id"`
			OutputIndex int64  `json:"output_index"`
			Delta       string `json:"delta"`
		}{
			Type:        "response.custom_tool_call_input.delta",
			Sequence:    u.Sequence,
			ItemID:      u.RefusalItemID,
			OutputIndex: u.OutputIndex,
			Delta:       u.CustomDelta,
		})
	case "response.custom_tool_call_input.done":
		return json.Marshal(struct {
			Type        string `json:"type"`
			Sequence    int64  `json:"sequence_number"`
			ItemID      string `json:"item_id"`
			OutputIndex int64  `json:"output_index"`
			Input       string `json:"input"`
		}{
			Type:        "response.custom_tool_call_input.done",
			Sequence:    u.Sequence,
			ItemID:      u.RefusalItemID,
			OutputIndex: u.OutputIndex,
			Input:       u.CustomInput,
		})
	}
	return nil, nil
}

func (u *OpenAIResponseStreamEventUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var tmp struct {
		Type              string                                `json:"type"`
		Sequence          int64                                 `json:"sequence_number"`
		Response          *OpenAIResponse                       `json:"response"`
		Error             *OpenAIResponseError                  `json:"error"`
		ItemID            string                                `json:"item_id"`
		OutputIndex       int64                                 `json:"output_index"`
		ContentIndex      int64                                 `json:"content_index"`
		SummaryIndex      int64                                 `json:"summary_index"`
		Refusal           string                                `json:"refusal"`
		Delta             string                                `json:"delta"`
		Text              string                                `json:"text"`
		Audio             string                                `json:"audio"`
		Transcript        string                                `json:"transcript"`
		Arguments         string                                `json:"arguments"`
		Code              string                                `json:"code"`
		Input             string                                `json:"input"`
		PartialImageIndex int64                                 `json:"partial_image_index"`
		B64JSON           string                                `json:"b64_json"`
		Item              OpenAIResponseOutputItemUnion         `json:"item"`
		Part              *OpenAIResponseOutputMessageContent   `json:"part"`
		Part2             *OpenAIResponseOutputReasoningSummary `json:"part"`
		AnnotationIndex   int64                                 `json:"annotation_index"`
		Annotation        json.RawMessage                       `json:"annotation"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	u.Sequence = tmp.Sequence
	switch tmp.Type {
	case "response.created", "response.queued", "response.in_progress", "response.completed", "response.failed", "response.incomplete":
		u.Response = tmp.Response
	case "error":
		u.Error = tmp.Error
	case "response.output_text.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.TextDelta = tmp.Delta
	case "response.output_text.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.Text = tmp.Text
	case "response.refusal.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.TextDelta = tmp.Delta
	case "response.refusal.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.Refusal = tmp.Refusal
	case "response.audio.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.AudioDelta = tmp.Delta
	case "response.audio.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.Audio = tmp.Audio
	case "response.audio_transcript.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.AudioDelta = tmp.Delta
	case "response.audio_transcript.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.Transcript = tmp.Transcript
	case "response.function_call_arguments.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.FuncDelta = tmp.Delta
	case "response.function_call_arguments.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.FuncArgs = tmp.Arguments
	case "response.code_interpreter_call_code.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.CodeDelta = tmp.Delta
	case "response.code_interpreter_call_code.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.Code = tmp.Code
	case "response.code_interpreter_call.completed", "response.code_interpreter_call.in_progress", "response.code_interpreter_call.interpreting":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
	case "response.file_search_call.completed", "response.file_search_call.in_progress", "response.file_search_call.searching":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
	case "response.web_search_call.completed", "response.web_search_call.in_progress", "response.web_search_call.searching":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
	case "response.image_generation_call.completed", "response.image_generation_call.generating", "response.image_generation_call.in_progress":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
	case "response.image_generation_call.partial_image":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.PartialImageIdx = tmp.PartialImageIndex
		u.PartialImageB64 = tmp.B64JSON
	case "response.mcp_call_arguments.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.McpDelta = tmp.Delta
	case "response.mcp_call_arguments.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.McpArgs = tmp.Arguments
	case "response.mcp_call.completed", "response.mcp_call.failed", "response.mcp_call.in_progress":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
	case "response.mcp_list_tools.completed", "response.mcp_list_tools.failed", "response.mcp_list_tools.in_progress":
		u.RefusalItemID = tmp.ItemID
	case "response.output_item.added", "response.output_item.done":
		u.OutputIndex = tmp.OutputIndex
		u.Item2 = tmp.Item
	case "response.content_part.added", "response.content_part.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.Part = tmp.Part
	case "response.reasoning_summary_part.added", "response.reasoning_summary_part.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.SummaryIndex = tmp.SummaryIndex
		u.Part2 = tmp.Part2
	case "response.reasoning_summary_text.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.SummaryIndex = tmp.SummaryIndex
		u.SummaryDelta = tmp.Delta
	case "response.reasoning_summary_text.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.SummaryIndex = tmp.SummaryIndex
		u.SummaryText = tmp.Text
	case "response.reasoning_text.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.ReasoningDelta = tmp.Delta
	case "response.reasoning_text.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.ReasoningText = tmp.Text
	case "response.output_text.annotation.added":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.ContentIndex = tmp.ContentIndex
		u.AnnotationIndex = tmp.AnnotationIndex
		u.Annotation = tmp.Annotation
	case "response.custom_tool_call_input.delta":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.CustomDelta = tmp.Delta
	case "response.custom_tool_call_input.done":
		u.RefusalItemID = tmp.ItemID
		u.OutputIndex = tmp.OutputIndex
		u.CustomInput = tmp.Input
	}
	return nil
}

// === 状态变更事件（10 个）===

// OpenAIResponseCreatedEvent 创建事件。
type OpenAIResponseCreatedEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.created"
	Sequence int64                         `json:"sequence_number"`
	Response *OpenAIResponse               `json:"response"`
}

func (OpenAIResponseCreatedEvent) isResponseStreamEvent() {}

// OpenAIResponseQueuedEvent 入队事件。
type OpenAIResponseQueuedEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.queued"
	Sequence int64                         `json:"sequence_number"`
	Response *OpenAIResponse               `json:"response"`
}

func (OpenAIResponseQueuedEvent) isResponseStreamEvent() {}

// OpenAIResponseInProgressEvent 进行中事件。
type OpenAIResponseInProgressEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.in_progress"
	Sequence int64                         `json:"sequence_number"`
	Response *OpenAIResponse               `json:"response"`
}

func (OpenAIResponseInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseCompletedEvent 完成事件。
type OpenAIResponseCompletedEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.completed"
	Sequence int64                         `json:"sequence_number"`
	Response *OpenAIResponse               `json:"response"`
}

func (OpenAIResponseCompletedEvent) isResponseStreamEvent() {}

// OpenAIResponseFailedEvent 失败事件。
type OpenAIResponseFailedEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.failed"
	Sequence int64                         `json:"sequence_number"`
	Response *OpenAIResponse               `json:"response"`
}

func (OpenAIResponseFailedEvent) isResponseStreamEvent() {}

// OpenAIResponseIncompleteEvent 未完成事件。
type OpenAIResponseIncompleteEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.incomplete"
	Sequence int64                         `json:"sequence_number"`
	Response *OpenAIResponse               `json:"response"`
}

func (OpenAIResponseIncompleteEvent) isResponseStreamEvent() {}

// OpenAIResponseRefusalDoneEvent 拒绝完成事件。
type OpenAIResponseRefusalDoneEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.refusal.done"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Refusal      string                        `json:"refusal"`
}

func (OpenAIResponseRefusalDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseErrorEvent 错误事件。
type OpenAIResponseErrorEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "error"
	Sequence int64                         `json:"sequence_number"`
	Error    *OpenAIResponseError          `json:"error"`
}

func (OpenAIResponseErrorEvent) isResponseStreamEvent() {}

// OpenAIResponseOutputItemAddedEvent 输出项添加事件。
type OpenAIResponseOutputItemAddedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.output_item.added"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	Item        OpenAIResponseOutputItemUnion `json:"item"`
}

func (OpenAIResponseOutputItemAddedEvent) isResponseStreamEvent() {}

// OpenAIResponseOutputItemDoneEvent 输出项完成事件。
type OpenAIResponseOutputItemDoneEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.output_item.done"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	Item        OpenAIResponseOutputItemUnion `json:"item"`
}

func (OpenAIResponseOutputItemDoneEvent) isResponseStreamEvent() {}

// === 文本事件（4 个）===

// OpenAIResponseTextDeltaEvent 文本增量事件。
type OpenAIResponseTextDeltaEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.output_text.delta"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Delta        string                        `json:"delta"`
	Logprobs     []json.RawMessage             `json:"logprobs,omitempty"`
}

func (OpenAIResponseTextDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseTextDoneEvent 文本完成事件。
type OpenAIResponseTextDoneEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.output_text.done"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Text         string                        `json:"text"`
}

func (OpenAIResponseTextDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseContentPartAddedEvent 内容部分添加事件。
type OpenAIResponseContentPartAddedEvent struct {
	Type         OpenAIResponseStreamEventType       `json:"type"` // "response.content_part.added"
	Sequence     int64                               `json:"sequence_number"`
	ItemID       string                              `json:"item_id"`
	OutputIndex  int64                               `json:"output_index"`
	ContentIndex int64                               `json:"content_index"`
	Part         *OpenAIResponseOutputMessageContent `json:"part"`
}

func (OpenAIResponseContentPartAddedEvent) isResponseStreamEvent() {}

// OpenAIResponseContentPartDoneEvent 内容部分完成事件。
type OpenAIResponseContentPartDoneEvent struct {
	Type         OpenAIResponseStreamEventType       `json:"type"` // "response.content_part.done"
	Sequence     int64                               `json:"sequence_number"`
	ItemID       string                              `json:"item_id"`
	OutputIndex  int64                               `json:"output_index"`
	ContentIndex int64                               `json:"content_index"`
	Part         *OpenAIResponseOutputMessageContent `json:"part"`
}

func (OpenAIResponseContentPartDoneEvent) isResponseStreamEvent() {}

// === 推理事件（4 个）===

// OpenAIResponseReasoningTextDeltaEvent 推理文本增量事件。
type OpenAIResponseReasoningTextDeltaEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.reasoning_text.delta"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Delta        string                        `json:"delta"`
}

func (OpenAIResponseReasoningTextDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseReasoningTextDoneEvent 推理文本完成事件。
type OpenAIResponseReasoningTextDoneEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.reasoning_text.done"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Text         string                        `json:"text"`
}

func (OpenAIResponseReasoningTextDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseReasoningSummaryPartAddedEvent 推理摘要部分添加事件。
type OpenAIResponseReasoningSummaryPartAddedEvent struct {
	Type         OpenAIResponseStreamEventType         `json:"type"` // "response.reasoning_summary_part.added"
	Sequence     int64                                 `json:"sequence_number"`
	ItemID       string                                `json:"item_id"`
	OutputIndex  int64                                 `json:"output_index"`
	SummaryIndex int64                                 `json:"summary_index"`
	Part         *OpenAIResponseOutputReasoningSummary `json:"part"`
}

func (OpenAIResponseReasoningSummaryPartAddedEvent) isResponseStreamEvent() {}

// OpenAIResponseReasoningSummaryPartDoneEvent 推理摘要部分完成事件。
type OpenAIResponseReasoningSummaryPartDoneEvent struct {
	Type         OpenAIResponseStreamEventType         `json:"type"` // "response.reasoning_summary_part.done"
	Sequence     int64                                 `json:"sequence_number"`
	ItemID       string                                `json:"item_id"`
	OutputIndex  int64                                 `json:"output_index"`
	SummaryIndex int64                                 `json:"summary_index"`
	Part         *OpenAIResponseOutputReasoningSummary `json:"part"`
}

func (OpenAIResponseReasoningSummaryPartDoneEvent) isResponseStreamEvent() {}

// === 工具调用事件（6 个）===

// OpenAIResponseFunctionCallArgumentsDeltaEvent 函数调用参数增量事件。
type OpenAIResponseFunctionCallArgumentsDeltaEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.function_call_arguments.delta"
	Sequence    int64                         `json:"sequence_number"`
	ItemID      string                        `json:"item_id"`
	OutputIndex int64                         `json:"output_index"`
	Delta       string                        `json:"delta"`
}

func (OpenAIResponseFunctionCallArgumentsDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseFunctionCallArgumentsDoneEvent 函数调用参数完成事件。
type OpenAIResponseFunctionCallArgumentsDoneEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.function_call_arguments.done"
	Sequence    int64                         `json:"sequence_number"`
	ItemID      string                        `json:"item_id"`
	OutputIndex int64                         `json:"output_index"`
	Arguments   string                        `json:"arguments"`
}

func (OpenAIResponseFunctionCallArgumentsDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseFileSearchCallInProgressEvent 文件搜索进行中事件。
type OpenAIResponseFileSearchCallInProgressEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.file_search_call.in_progress"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseFileSearchCallInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseFileSearchCallSearchingEvent 文件搜索中事件。
type OpenAIResponseFileSearchCallSearchingEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.file_search_call.searching"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseFileSearchCallSearchingEvent) isResponseStreamEvent() {}

// OpenAIResponseFileSearchCallCompletedEvent 文件搜索完成事件。
type OpenAIResponseFileSearchCallCompletedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.file_search_call.completed"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseFileSearchCallCompletedEvent) isResponseStreamEvent() {}

// OpenAIResponseWebSearchCallInProgressEvent Web 搜索进行中事件。
type OpenAIResponseWebSearchCallInProgressEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.web_search_call.in_progress"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseWebSearchCallInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseWebSearchCallSearchingEvent Web 搜索中事件。
type OpenAIResponseWebSearchCallSearchingEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.web_search_call.searching"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseWebSearchCallSearchingEvent) isResponseStreamEvent() {}

// OpenAIResponseWebSearchCallCompletedEvent Web 搜索完成事件。
type OpenAIResponseWebSearchCallCompletedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.web_search_call.completed"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseWebSearchCallCompletedEvent) isResponseStreamEvent() {}

// === MCP 事件（9 个）===

// OpenAIResponseMcpCallInProgressEvent MCP 调用进行中事件。
type OpenAIResponseMcpCallInProgressEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.mcp_call.in_progress"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseMcpCallInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpCallArgumentsDeltaEvent MCP 调用参数增量事件。
type OpenAIResponseMcpCallArgumentsDeltaEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.mcp_call_arguments.delta"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
	Delta       string                        `json:"delta"`
}

func (OpenAIResponseMcpCallArgumentsDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpCallArgumentsDoneEvent MCP 调用参数完成事件。
type OpenAIResponseMcpCallArgumentsDoneEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.mcp_call_arguments.done"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
	Arguments   string                        `json:"arguments"`
}

func (OpenAIResponseMcpCallArgumentsDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpCallCompletedEvent MCP 调用完成事件。
type OpenAIResponseMcpCallCompletedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.mcp_call.completed"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseMcpCallCompletedEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpCallFailedEvent MCP 调用失败事件。
type OpenAIResponseMcpCallFailedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.mcp_call.failed"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseMcpCallFailedEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpListToolsInProgressEvent MCP 工具列表进行中事件。
type OpenAIResponseMcpListToolsInProgressEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.mcp_list_tools.in_progress"
	Sequence int64                         `json:"sequence_number"`
	ItemID   string                        `json:"item_id"`
}

func (OpenAIResponseMcpListToolsInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpListToolsCompletedEvent MCP 工具列表完成事件。
type OpenAIResponseMcpListToolsCompletedEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.mcp_list_tools.completed"
	Sequence int64                         `json:"sequence_number"`
	ItemID   string                        `json:"item_id"`
}

func (OpenAIResponseMcpListToolsCompletedEvent) isResponseStreamEvent() {}

// OpenAIResponseMcpListToolsFailedEvent MCP 工具列表失败事件。
type OpenAIResponseMcpListToolsFailedEvent struct {
	Type     OpenAIResponseStreamEventType `json:"type"` // "response.mcp_list_tools.failed"
	Sequence int64                         `json:"sequence_number"`
	ItemID   string                        `json:"item_id"`
}

func (OpenAIResponseMcpListToolsFailedEvent) isResponseStreamEvent() {}

// === Code Interpreter 事件（5 个）===

// OpenAIResponseCodeInterpreterCallInProgressEvent 代码解释器进行中事件。
type OpenAIResponseCodeInterpreterCallInProgressEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.code_interpreter_call.in_progress"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseCodeInterpreterCallInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseCodeInterpreterCallInterpretingEvent 代码解释器解释中事件。
type OpenAIResponseCodeInterpreterCallInterpretingEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.code_interpreter_call.interpreting"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseCodeInterpreterCallInterpretingEvent) isResponseStreamEvent() {}

// OpenAIResponseCodeInterpreterCallCodeDeltaEvent 代码增量事件。
type OpenAIResponseCodeInterpreterCallCodeDeltaEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.code_interpreter_call_code.delta"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
	Delta       string                        `json:"delta"`
}

func (OpenAIResponseCodeInterpreterCallCodeDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseCodeInterpreterCallCodeDoneEvent 代码完成事件。
type OpenAIResponseCodeInterpreterCallCodeDoneEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.code_interpreter_call_code.done"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
	Code        string                        `json:"code"`
}

func (OpenAIResponseCodeInterpreterCallCodeDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseCodeInterpreterCallCompletedEvent 代码解释器完成事件。
type OpenAIResponseCodeInterpreterCallCompletedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.code_interpreter_call.completed"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseCodeInterpreterCallCompletedEvent) isResponseStreamEvent() {}

// === Image Generation 事件（4 个）===

// OpenAIResponseImageGenCallInProgressEvent 图像生成进行中事件。
type OpenAIResponseImageGenCallInProgressEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.image_generation_call.in_progress"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseImageGenCallInProgressEvent) isResponseStreamEvent() {}

// OpenAIResponseImageGenCallGeneratingEvent 图像生成中事件。
type OpenAIResponseImageGenCallGeneratingEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.image_generation_call.generating"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseImageGenCallGeneratingEvent) isResponseStreamEvent() {}

// OpenAIResponseImageGenCallPartialImageEvent 图像部分生成事件。
type OpenAIResponseImageGenCallPartialImageEvent struct {
	Type              OpenAIResponseStreamEventType `json:"type"` // "response.image_generation_call.partial_image"
	Sequence          int64                         `json:"sequence_number"`
	OutputIndex       int64                         `json:"output_index"`
	ItemID            string                        `json:"item_id"`
	PartialImageIndex int64                         `json:"partial_image_index"`
	B64JSON           string                        `json:"b64_json"`
}

func (OpenAIResponseImageGenCallPartialImageEvent) isResponseStreamEvent() {}

// OpenAIResponseImageGenCallCompletedEvent 图像生成完成事件。
type OpenAIResponseImageGenCallCompletedEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.image_generation_call.completed"
	Sequence    int64                         `json:"sequence_number"`
	OutputIndex int64                         `json:"output_index"`
	ItemID      string                        `json:"item_id"`
}

func (OpenAIResponseImageGenCallCompletedEvent) isResponseStreamEvent() {}

// === Audio 事件（4 个）===

// OpenAIResponseAudioDeltaEvent 音频增量事件。
type OpenAIResponseAudioDeltaEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.audio.delta"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Delta        string                        `json:"delta"`
}

func (OpenAIResponseAudioDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseAudioDoneEvent 音频完成事件。
type OpenAIResponseAudioDoneEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.audio.done"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Audio        string                        `json:"audio"`
}

func (OpenAIResponseAudioDoneEvent) isResponseStreamEvent() {}

// OpenAIResponseAudioTranscriptDeltaEvent 音频转录增量事件。
type OpenAIResponseAudioTranscriptDeltaEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.audio_transcript.delta"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Delta        string                        `json:"delta"`
}

func (OpenAIResponseAudioTranscriptDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseAudioTranscriptDoneEvent 音频转录完成事件。
type OpenAIResponseAudioTranscriptDoneEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.audio_transcript.done"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Transcript   string                        `json:"transcript"`
}

func (OpenAIResponseAudioTranscriptDoneEvent) isResponseStreamEvent() {}

// === 推理摘要文本事件（2 个） ===

// OpenAIResponseReasoningSummaryTextDeltaEvent 推理摘要文本增量事件。
type OpenAIResponseReasoningSummaryTextDeltaEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.reasoning_summary_text.delta"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	SummaryIndex int64                         `json:"summary_index"`
	Delta        string                        `json:"delta"`
}

func (OpenAIResponseReasoningSummaryTextDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseReasoningSummaryTextDoneEvent 推理摘要文本完成事件。
type OpenAIResponseReasoningSummaryTextDoneEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.reasoning_summary_text.done"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	SummaryIndex int64                         `json:"summary_index"`
	Text         string                        `json:"text"`
}

func (OpenAIResponseReasoningSummaryTextDoneEvent) isResponseStreamEvent() {}

// === 拒绝文本事件（1 个 delta，与 OpenAIResponseRefusalDoneEvent 配对） ===

// OpenAIResponseRefusalDeltaEvent 拒绝文本增量事件。
type OpenAIResponseRefusalDeltaEvent struct {
	Type         OpenAIResponseStreamEventType `json:"type"` // "response.refusal.delta"
	Sequence     int64                         `json:"sequence_number"`
	ItemID       string                        `json:"item_id"`
	OutputIndex  int64                         `json:"output_index"`
	ContentIndex int64                         `json:"content_index"`
	Delta        string                        `json:"delta"`
}

func (OpenAIResponseRefusalDeltaEvent) isResponseStreamEvent() {}

// === 自定义工具调用输入事件（2 个） ===

// OpenAIResponseCustomToolCallInputDeltaEvent 自定义工具调用输入增量事件。
type OpenAIResponseCustomToolCallInputDeltaEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.custom_tool_call_input.delta"
	Sequence    int64                         `json:"sequence_number"`
	ItemID      string                        `json:"item_id"`
	OutputIndex int64                         `json:"output_index"`
	Delta       string                        `json:"delta"`
}

func (OpenAIResponseCustomToolCallInputDeltaEvent) isResponseStreamEvent() {}

// OpenAIResponseCustomToolCallInputDoneEvent 自定义工具调用输入完成事件。
type OpenAIResponseCustomToolCallInputDoneEvent struct {
	Type        OpenAIResponseStreamEventType `json:"type"` // "response.custom_tool_call_input.done"
	Sequence    int64                         `json:"sequence_number"`
	ItemID      string                        `json:"item_id"`
	OutputIndex int64                         `json:"output_index"`
	Input       string                        `json:"input"`
}

func (OpenAIResponseCustomToolCallInputDoneEvent) isResponseStreamEvent() {}

// === 输出文本注释事件（1 个） ===

// OpenAIResponseOutputTextAnnotationAddedEvent 输出文本注释添加事件。
type OpenAIResponseOutputTextAnnotationAddedEvent struct {
	Type            OpenAIResponseStreamEventType `json:"type"` // "response.output_text.annotation.added"
	Sequence        int64                         `json:"sequence_number"`
	ItemID          string                        `json:"item_id"`
	OutputIndex     int64                         `json:"output_index"`
	ContentIndex    int64                         `json:"content_index"`
	AnnotationIndex int64                         `json:"annotation_index"`
	// Annotation 字段官方 SDK 为 any；本包使用 json.RawMessage 以避免引入额外类型依赖。
	Annotation json.RawMessage `json:"annotation"`
}

func (OpenAIResponseOutputTextAnnotationAddedEvent) isResponseStreamEvent() {}

// === P2-1 已覆盖（2026-07-01） ===
//
// 官方 SDK ResponseStreamEventUnion 共 53 变体（不含 computer_call.* / queued_item.*），
// 全部已在上面实现。computer_call.* / queued_item.* 在官方 SDK 中不存在，注释移除。

// ParseResponseSSEEvent 解析 Responses SSE 行（"data: {...}"）为 OpenAIResponseStreamEventUnion。
func ParseResponseSSEEvent(line string) (OpenAIResponseStreamEventUnion, error) {
	data := strings.TrimPrefix(line, "data: ")
	data = strings.TrimSpace(data)
	if data == "" {
		return nil, nil
	}

	var typeOnly struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal([]byte(data), &typeOnly); err != nil {
		return nil, err
	}

	switch typeOnly.Type {
	// Lifecycle (6)
	case "response.created":
		var e OpenAIResponseCreatedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.in_progress":
		var e OpenAIResponseInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.completed":
		var e OpenAIResponseCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.incomplete":
		var e OpenAIResponseIncompleteEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.failed":
		var e OpenAIResponseFailedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.queued":
		var e OpenAIResponseQueuedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Error (1)
	case "error":
		var e OpenAIResponseErrorEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Output Text (4)
	case "response.output_text.delta":
		var e OpenAIResponseTextDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.output_text.done":
		var e OpenAIResponseTextDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.refusal.delta":
		var e OpenAIResponseRefusalDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.refusal.done":
		var e OpenAIResponseRefusalDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Audio (4)
	case "response.audio.delta":
		var e OpenAIResponseAudioDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.audio.done":
		var e OpenAIResponseAudioDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.audio_transcript.delta":
		var e OpenAIResponseAudioTranscriptDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.audio_transcript.done":
		var e OpenAIResponseAudioTranscriptDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Function Call (2)
	case "response.function_call_arguments.delta":
		var e OpenAIResponseFunctionCallArgumentsDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.function_call_arguments.done":
		var e OpenAIResponseFunctionCallArgumentsDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Code Interpreter (5)
	case "response.code_interpreter_call_code.delta":
		var e OpenAIResponseCodeInterpreterCallCodeDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.code_interpreter_call_code.done":
		var e OpenAIResponseCodeInterpreterCallCodeDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.code_interpreter_call.completed":
		var e OpenAIResponseCodeInterpreterCallCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.code_interpreter_call.in_progress":
		var e OpenAIResponseCodeInterpreterCallInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.code_interpreter_call.interpreting":
		var e OpenAIResponseCodeInterpreterCallInterpretingEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// File Search (3)
	case "response.file_search_call.completed":
		var e OpenAIResponseFileSearchCallCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.file_search_call.in_progress":
		var e OpenAIResponseFileSearchCallInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.file_search_call.searching":
		var e OpenAIResponseFileSearchCallSearchingEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Web Search (3)
	case "response.web_search_call.completed":
		var e OpenAIResponseWebSearchCallCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.web_search_call.in_progress":
		var e OpenAIResponseWebSearchCallInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.web_search_call.searching":
		var e OpenAIResponseWebSearchCallSearchingEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Image Gen (4)
	case "response.image_generation_call.completed":
		var e OpenAIResponseImageGenCallCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.image_generation_call.generating":
		var e OpenAIResponseImageGenCallGeneratingEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.image_generation_call.in_progress":
		var e OpenAIResponseImageGenCallInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.image_generation_call.partial_image":
		var e OpenAIResponseImageGenCallPartialImageEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// MCP (8)
	case "response.mcp_call_arguments.delta":
		var e OpenAIResponseMcpCallArgumentsDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_call_arguments.done":
		var e OpenAIResponseMcpCallArgumentsDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_call.completed":
		var e OpenAIResponseMcpCallCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_call.failed":
		var e OpenAIResponseMcpCallFailedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_call.in_progress":
		var e OpenAIResponseMcpCallInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_list_tools.completed":
		var e OpenAIResponseMcpListToolsCompletedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_list_tools.failed":
		var e OpenAIResponseMcpListToolsFailedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.mcp_list_tools.in_progress":
		var e OpenAIResponseMcpListToolsInProgressEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Output Item (2)
	case "response.output_item.added":
		var e OpenAIResponseOutputItemAddedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.output_item.done":
		var e OpenAIResponseOutputItemDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Content Part (2)
	case "response.content_part.added":
		var e OpenAIResponseContentPartAddedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.content_part.done":
		var e OpenAIResponseContentPartDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Reasoning (6)
	case "response.reasoning_summary_part.added":
		var e OpenAIResponseReasoningSummaryPartAddedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.reasoning_summary_part.done":
		var e OpenAIResponseReasoningSummaryPartDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.reasoning_summary_text.delta":
		var e OpenAIResponseReasoningSummaryTextDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.reasoning_summary_text.done":
		var e OpenAIResponseReasoningSummaryTextDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.reasoning_text.delta":
		var e OpenAIResponseReasoningTextDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.reasoning_text.done":
		var e OpenAIResponseReasoningTextDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Annotation (1)
	case "response.output_text.annotation.added":
		var e OpenAIResponseOutputTextAnnotationAddedEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	// Custom Tool (2)
	case "response.custom_tool_call_input.delta":
		var e OpenAIResponseCustomToolCallInputDeltaEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	case "response.custom_tool_call_input.done":
		var e OpenAIResponseCustomToolCallInputDoneEvent
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		return e, nil
	default:
		return nil, nil
	}
}
