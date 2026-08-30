package openai_test

import (
	"encoding/json"
	"testing"

	openai "hegq.net/mogo/pkg/protocol/openai"
)

// === Audio Transcription P0 补齐测试（与 wire_format_responses_test.go 同包）===

func TestWireFormatAudioTranscriptionParamsAndStream(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "ParamsMinimal",
			v: openai.OpenAIAudioTranscriptionParams{
				File:  "file_xxx",
				Model: "whisper-1",
			},
			want: `{"file":"file_xxx","model":"whisper-1"}`,
		},
		{
			name: "ParamsVerboseJSONAutoChunkingLogprobs",
			v: openai.OpenAIAudioTranscriptionParams{
				File:             "audio.mp3",
				Model:            "gpt-4o-transcribe",
				Language:         "en",
				Prompt:           "hello",
				ResponseFormat:   openai.OpenAIAudioResponseFormatVerboseJSON,
				ChunkingStrategy: &openai.OpenAIAudioTranscriptionChunkingStrategyUnionImpl{Auto: "auto"},
				Include:          []openai.OpenAITranscriptionInclude{openai.OpenAITranscriptionIncludeLogprobs},
				Stream:           ptrBool(true),
			},
			want: `{"file":"audio.mp3","model":"gpt-4o-transcribe","language":"en","prompt":"hello","response_format":"verbose_json","chunking_strategy":"auto","include":["logprobs"],"stream":true}`,
		},
		{
			name: "ParamsVadConfig",
			v: openai.OpenAIAudioTranscriptionParams{
				File:  "audio.mp3",
				Model: "gpt-4o-transcribe",
				ChunkingStrategy: func() *openai.OpenAIAudioTranscriptionChunkingStrategyUnionImpl {
					vad, _ := json.Marshal(openai.OpenAIAudioTranscriptionChunkingStrategyVadConfig{
						Type:              "server_vad",
						SilenceDurationMs: ptrInt64(200),
						Threshold:         ptrFloat64(0.5),
					})
					return &openai.OpenAIAudioTranscriptionChunkingStrategyUnionImpl{
						Type:      "server_vad",
						VadConfig: vad,
					}
				}(),
			},
			want: `{"file":"audio.mp3","model":"gpt-4o-transcribe","chunking_strategy":{"type":"server_vad","silence_duration_ms":200,"threshold":0.5}}`,
		},
		{
			name: "ParamsKnownSpeakersDiarize",
			v: openai.OpenAIAudioTranscriptionParams{
				File:                   "meeting.mp3",
				Model:                  "gpt-4o-transcribe-diarize",
				KnownSpeakerNames:      []string{"Alice", "Bob"},
				KnownSpeakerReferences: []string{"data:audio/mp3;base64,AAA"},
			},
			want: `{"file":"meeting.mp3","model":"gpt-4o-transcribe-diarize","known_speaker_names":["Alice","Bob"],"known_speaker_references":["data:audio/mp3;base64,AAA"]}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("wire-format mismatch:\n  got:  %s\n  want: %s", got, c.want)
			}
		})
	}
}

func TestWireFormatAudioTranscriptionResponseUnion(t *testing.T) {
	t.Run("BaseWithTokensUsage", func(t *testing.T) {
		v := openai.OpenAIAudioTranscription{
			Text: "hello world",
			Usage: func() *openai.OpenAITranscriptionUsageUnionImpl {
				usage, _ := json.Marshal(openai.OpenAITranscriptionUsageTokens{
					Type:         "tokens",
					InputTokens:  10,
					OutputTokens: 5,
					TotalTokens:  15,
					InputTokenDetails: &openai.OpenAITranscriptionUsageTokensInputTokenDetails{
						AudioTokens: 8,
						TextTokens:  2,
					},
				})
				return &openai.OpenAITranscriptionUsageUnionImpl{Type: "tokens", Tokens: usage}
			}(),
		}
		got, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		want := `{"text":"hello world","usage":{"type":"tokens","input_tokens":10,"output_tokens":5,"total_tokens":15,"input_token_details":{"audio_tokens":8,"text_tokens":2}}}`
		if string(got) != want {
			t.Errorf("mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})

	t.Run("VerboseWithDurationUsage", func(t *testing.T) {
		v := openai.OpenAIAudioTranscriptionVerbose{
			Text:     "hi",
			Language: "en",
			Duration: ptrFloat64(3.5),
			Usage: func() *openai.OpenAITranscriptionUsageUnionImpl {
				usage, _ := json.Marshal(openai.OpenAITranscriptionUsageDuration{
					Type:    "duration",
					Seconds: 3.5,
				})
				return &openai.OpenAITranscriptionUsageUnionImpl{Type: "duration", Seconds: usage}
			}(),
		}
		got, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		want := `{"text":"hi","language":"en","duration":3.5,"usage":{"type":"duration","seconds":3.5}}`
		if string(got) != want {
			t.Errorf("mismatch:\n  got:  %s\n  want: %s", got, want)
		}
	})
}

func TestWireFormatAudioTranscriptionStreamEvents(t *testing.T) {
	cases := []struct {
		name string
		v    openai.OpenAITranscriptionStreamEvent
		want string
	}{
		{
			name: "TextSegmentDiarized",
			v: openai.OpenAITranscriptionTextSegmentEvent{
				Type:    openai.OpenAITranscriptionStreamEventTypeSegment,
				ID:      "seg_1",
				Start:   0.0,
				End:     1.2,
				Speaker: "Alice",
				Text:    "hello",
			},
			want: `{"type":"transcript.text.segment","id":"seg_1","end":1.2,"speaker":"Alice","start":0,"text":"hello"}`,
		},
		{
			name: "TextDelta",
			v: openai.OpenAITranscriptionTextDeltaEvent{
				Type:      openai.OpenAITranscriptionStreamEventTypeDelta,
				Delta:     " wor",
				SegmentID: "seg_1",
			},
			want: `{"type":"transcript.text.delta","delta":" wor","segment_id":"seg_1"}`,
		},
		{
			name: "TextDoneWithUsage",
			v: openai.OpenAITranscriptionTextDoneEvent{
				Type: openai.OpenAITranscriptionStreamEventTypeDone,
				Text: "hello world",
				Usage: &openai.OpenAITranscriptionUsageTokens{
					Type:         "tokens",
					InputTokens:  3,
					OutputTokens: 0,
					TotalTokens:  3,
				},
			},
			want: `{"type":"transcript.text.done","text":"hello world","usage":{"type":"tokens","input_tokens":3,"output_tokens":0,"total_tokens":3}}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := json.Marshal(c.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != c.want {
				t.Errorf("mismatch:\n  got:  %s\n  want: %s", got, c.want)
			}
		})
	}
}

func TestAudioTranscriptionUnionClosure(t *testing.T) {
	// ChunkingStrategy union
	var _ openai.OpenAIAudioTranscriptionChunkingStrategyUnion = openai.OpenAIAudioTranscriptionChunkingStrategyAutoValue
	var _ openai.OpenAIAudioTranscriptionChunkingStrategyUnion = openai.OpenAIAudioTranscriptionChunkingStrategyVadConfig{}

	// Response union
	var _ openai.OpenAIAudioTranscriptionResponseUnion = openai.OpenAIAudioTranscription{}
	var _ openai.OpenAIAudioTranscriptionResponseUnion = openai.OpenAIAudioTranscriptionVerbose{}

	// Usage union
	var _ openai.OpenAITranscriptionUsageUnion = openai.OpenAITranscriptionUsageTokens{}
	var _ openai.OpenAITranscriptionUsageUnion = openai.OpenAITranscriptionUsageDuration{}

	// StreamEvent union
	var _ openai.OpenAITranscriptionStreamEvent = openai.OpenAITranscriptionTextSegmentEvent{}
	var _ openai.OpenAITranscriptionStreamEvent = openai.OpenAITranscriptionTextDeltaEvent{}
	var _ openai.OpenAITranscriptionStreamEvent = openai.OpenAITranscriptionTextDoneEvent{}
}
