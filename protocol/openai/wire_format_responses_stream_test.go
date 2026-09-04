package openai

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseResponseSSEEvent_Created(t *testing.T) {
	jsonStr := `{"type":"response.created","response":{"id":"resp_xxx"},"sequence_number":1}`
	event, err := ParseResponseSSEEvent("data: " + jsonStr)
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event == nil {
		t.Fatalf("event should not be nil")
	}
	e := event.(OpenAIResponseCreatedEvent)
	if got := string(e.Type); got != "response.created" {
		t.Errorf("Type = %q, want %q", got, "response.created")
	}
}

func TestParseResponseSSEEvent_OutputTextDelta(t *testing.T) {
	jsonStr := `{"type":"response.output_text.delta","delta":"hello","item_id":"item_xxx","output_index":0,"content_index":0}`
	event, err := ParseResponseSSEEvent("data: " + jsonStr)
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event == nil {
		t.Fatalf("event should not be nil")
	}
	e := event.(OpenAIResponseTextDeltaEvent)
	if got := string(e.Type); got != "response.output_text.delta" {
		t.Errorf("Type = %q, want %q", got, "response.output_text.delta")
	}
	if e.Delta != "hello" {
		t.Errorf("Delta = %q, want %q", e.Delta, "hello")
	}
}

func TestParseResponseSSEEvent_Error(t *testing.T) {
	jsonStr := `{"type":"error","error":{"code":"rate_limit_exceeded","message":"Rate limit exceeded"}}`
	event, err := ParseResponseSSEEvent("data: " + jsonStr)
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event == nil {
		t.Fatalf("event should not be nil")
	}
	e := event.(OpenAIResponseErrorEvent)
	if got := string(e.Type); got != "error" {
		t.Errorf("Type = %q, want %q", got, "error")
	}
	if e.Error == nil {
		t.Fatalf("Error should not be nil")
	}
	if e.Error.Code != "rate_limit_exceeded" {
		t.Errorf("Code = %q, want %q", e.Error.Code, "rate_limit_exceeded")
	}
}

func TestParseResponseSSEEvent_AllTypes(t *testing.T) {
	testCases := []struct {
		typeVal string
		json    string
	}{
		{"response.created", `{"type":"response.created","response":{"id":"r1"},"sequence_number":1}`},
		{"response.in_progress", `{"type":"response.in_progress","response":{"id":"r2"},"sequence_number":2}`},
		{"response.completed", `{"type":"response.completed","response":{"id":"r3"},"sequence_number":3}`},
		{"response.incomplete", `{"type":"response.incomplete","response":{"id":"r4"},"sequence_number":4}`},
		{"response.failed", `{"type":"response.failed","response":{"id":"r5"},"sequence_number":5}`},
		{"response.queued", `{"type":"response.queued","response":{"id":"r6"},"sequence_number":6}`},
		{"error", `{"type":"error","error":{"code":"rate_limit","message":"rate limit"}}`},
		{"response.output_text.delta", `{"type":"response.output_text.delta","delta":"hi","item_id":"i1","output_index":0}`},
		{"response.output_text.done", `{"type":"response.output_text.done","text":"hi","item_id":"i1"}`},
		{"response.refusal.delta", `{"type":"response.refusal.delta","delta":"no","item_id":"i1"}`},
		{"response.refusal.done", `{"type":"response.refusal.done","refusal":"no","item_id":"i1"}`},
		{"response.audio.delta", `{"type":"response.audio.delta","delta":"YWJj"}`},
		{"response.audio.done", `{"type":"response.audio.done"}`},
		{"response.audio_transcript.delta", `{"type":"response.audio_transcript.delta","delta":"transcript"}`},
		{"response.audio_transcript.done", `{"type":"response.audio_transcript.done","transcript":"full"}`},
		{"response.function_call_arguments.delta", `{"type":"response.function_call_arguments.delta","arguments":"{\"name\":\"","item_id":"fc1"}`},
		{"response.function_call_arguments.done", `{"type":"response.function_call_arguments.done","arguments":"{}","item_id":"fc1"}`},
		{"response.code_interpreter_call_code.delta", `{"type":"response.code_interpreter_call_code.delta","delta":"print(","item_id":"ci1"}`},
		{"response.code_interpreter_call_code.done", `{"type":"response.code_interpreter_call_code.done","item_id":"ci1"}`},
		{"response.code_interpreter_call.completed", `{"type":"response.code_interpreter_call.completed","item_id":"ci1"}`},
		{"response.code_interpreter_call.in_progress", `{"type":"response.code_interpreter_call.in_progress","item_id":"ci1"}`},
		{"response.code_interpreter_call.interpreting", `{"type":"response.code_interpreter_call.interpreting","item_id":"ci1"}`},
		{"response.file_search_call.completed", `{"type":"response.file_search_call.completed","item_id":"fs1"}`},
		{"response.file_search_call.in_progress", `{"type":"response.file_search_call.in_progress","item_id":"fs1"}`},
		{"response.file_search_call.searching", `{"type":"response.file_search_call.searching","item_id":"fs1"}`},
		{"response.web_search_call.completed", `{"type":"response.web_search_call.completed","item_id":"ws1"}`},
		{"response.web_search_call.in_progress", `{"type":"response.web_search_call.in_progress","item_id":"ws1"}`},
		{"response.web_search_call.searching", `{"type":"response.web_search_call.searching","item_id":"ws1"}`},
		{"response.image_generation_call.completed", `{"type":"response.image_generation_call.completed","item_id":"ig1"}`},
		{"response.image_generation_call.generating", `{"type":"response.image_generation_call.generating","item_id":"ig1"}`},
		{"response.image_generation_call.in_progress", `{"type":"response.image_generation_call.in_progress","item_id":"ig1"}`},
		{"response.image_generation_call.partial_image", `{"type":"response.image_generation_call.partial_image","b64_json":"abc"}`},
		{"response.mcp_call_arguments.delta", `{"type":"response.mcp_call_arguments.delta","arguments":"{}","item_id":"mcp1"}`},
		{"response.mcp_call_arguments.done", `{"type":"response.mcp_call_arguments.done","item_id":"mcp1"}`},
		{"response.mcp_call.completed", `{"type":"response.mcp_call.completed","item_id":"mcp1"}`},
		{"response.mcp_call.failed", `{"type":"response.mcp_call.failed","item_id":"mcp1"}`},
		{"response.mcp_call.in_progress", `{"type":"response.mcp_call.in_progress","item_id":"mcp1"}`},
		{"response.mcp_list_tools.completed", `{"type":"response.mcp_list_tools.completed"}`},
		{"response.mcp_list_tools.failed", `{"type":"response.mcp_list_tools.failed"}`},
		{"response.mcp_list_tools.in_progress", `{"type":"response.mcp_list_tools.in_progress"}`},
		{"response.content_part.added", `{"type":"response.content_part.added","item_id":"item_x","content_index":0}`},
		{"response.content_part.done", `{"type":"response.content_part.done","item_id":"item_x","content_index":0}`},
		{"response.reasoning_summary_part.added", `{"type":"response.reasoning_summary_part.added","item_id":"rs1","summary_index":0}`},
		{"response.reasoning_summary_part.done", `{"type":"response.reasoning_summary_part.done","item_id":"rs1","summary_index":0}`},
		{"response.reasoning_summary_text.delta", `{"type":"response.reasoning_summary_text.delta","delta":"summary","item_id":"rs1"}`},
		{"response.reasoning_summary_text.done", `{"type":"response.reasoning_summary_text.done","text":"summary text","item_id":"rs1"}`},
		{"response.reasoning_text.delta", `{"type":"response.reasoning_text.delta","delta":"thinking","item_id":"rt1"}`},
		{"response.reasoning_text.done", `{"type":"response.reasoning_text.done","item_id":"rt1"}`},
		{"response.output_text.annotation.added", `{"type":"response.output_text.annotation.added","item_id":"item_x","annotation_index":0}`},
		{"response.custom_tool_call_input.delta", `{"type":"response.custom_tool_call_input.delta","delta":"input","item_id":"ctc1"}`},
		{"response.custom_tool_call_input.done", `{"type":"response.custom_tool_call_input.done","input":"done","item_id":"ctc1"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.typeVal, func(t *testing.T) {
			event, err := ParseResponseSSEEvent("data: " + tc.json)
			if err != nil {
				t.Fatalf("ParseResponseSSEEvent (type %s): %v", tc.typeVal, err)
			}
			if event == nil {
				t.Fatalf("event should not be nil for type %s", tc.typeVal)
			}

			var parsed struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal([]byte(tc.json), &parsed); err != nil {
				t.Fatalf("Unmarshal tc.json: %v", err)
			}
			if parsed.Type != tc.typeVal {
				t.Errorf("parsed.Type = %q, want %q", parsed.Type, tc.typeVal)
			}
		})
	}
}

func TestParseResponseSSEEvent_UnknownType(t *testing.T) {
	event, err := ParseResponseSSEEvent(`data: {"type":"unknown.event"}`)
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event != nil {
		t.Errorf("event should be nil, got %v", event)
	}
}

func TestParseResponseSSEEvent_Empty(t *testing.T) {
	event, err := ParseResponseSSEEvent("")
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event != nil {
		t.Errorf("event should be nil, got %v", event)
	}
}

func TestParseResponseSSEEvent_NoPrefix(t *testing.T) {
	jsonStr := `{"type":"response.created","response":{"id":"r1"},"sequence_number":1}`
	event, err := ParseResponseSSEEvent(jsonStr)
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event == nil {
		t.Fatalf("event should not be nil")
	}
	if got := string(event.(OpenAIResponseCreatedEvent).Type); got != "response.created" {
		t.Errorf("Type = %q, want %q", got, "response.created")
	}
}

func TestParseResponseSSEEvent_WithWhitespace(t *testing.T) {
	jsonStr := `{"type":"response.created","response":{"id":"r1"},"sequence_number":1}`
	event, err := ParseResponseSSEEvent("data:   " + jsonStr + "   ")
	if err != nil {
		t.Fatalf("ParseResponseSSEEvent: %v", err)
	}
	if event == nil {
		t.Fatalf("event should not be nil")
	}
	if got := string(event.(OpenAIResponseCreatedEvent).Type); got != "response.created" {
		t.Errorf("Type = %q, want %q", got, "response.created")
	}
}

func TestParseResponseSSEEvent_InvalidJSON(t *testing.T) {
	if _, err := ParseResponseSSEEvent("data: {invalid json}"); err == nil {
		t.Fatalf("expected error on invalid JSON")
	}
}

func TestParseResponseSSEEvent_LineByLineSSE(t *testing.T) {
	lines := []string{
		`data: {"type":"response.created","response":{"id":"r1"},"sequence_number":1}`,
		`data: {"type":"response.in_progress","response":{"id":"r1"},"sequence_number":2}`,
		`data: {"type":"response.output_text.delta","delta":"Hello","item_id":"i1","output_index":0}`,
		`data: {"type":"response.output_text.done","text":"Hello","item_id":"i1","output_index":0}`,
		`data: {"type":"response.completed","response":{"id":"r1"},"sequence_number":5}`,
		``,
		`data: [DONE]`,
	}

	var events []OpenAIResponseStreamEventUnion
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "data: [DONE]") {
			continue
		}
		event, err := ParseResponseSSEEvent(line)
		if err != nil {
			t.Fatalf("ParseResponseSSEEvent (line %q): %v", line, err)
		}
		if event != nil {
			events = append(events, event)
		}
	}

	if got := len(events); got != 5 {
		t.Fatalf("len(events) = %d, want 5", got)
	}
	if got := string(events[0].(OpenAIResponseCreatedEvent).Type); got != "response.created" {
		t.Errorf("events[0].Type = %q, want %q", got, "response.created")
	}
	if got := string(events[1].(OpenAIResponseInProgressEvent).Type); got != "response.in_progress" {
		t.Errorf("events[1].Type = %q, want %q", got, "response.in_progress")
	}
	if got := string(events[2].(OpenAIResponseTextDeltaEvent).Type); got != "response.output_text.delta" {
		t.Errorf("events[2].Type = %q, want %q", got, "response.output_text.delta")
	}
	if got := string(events[3].(OpenAIResponseTextDoneEvent).Type); got != "response.output_text.done" {
		t.Errorf("events[3].Type = %q, want %q", got, "response.output_text.done")
	}
	if got := string(events[4].(OpenAIResponseCompletedEvent).Type); got != "response.completed" {
		t.Errorf("events[4].Type = %q, want %q", got, "response.completed")
	}
}
