package openai_test

import (
	"encoding/json"
	"testing"

	openai "hegq.net/mogo/pkg/protocol/openai"
)

// === A. AudioTranslation 请求参数 ===

func TestWireFormatAudioTranslationParams(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Params_minimal",
			v: openai.OpenAITranslationParams{
				File:  "base64data",
				Model: "whisper-1",
			},
			want: `{"file":"base64data","model":"whisper-1"}`,
		},
		{
			name: "Params_with_response_format_verbose_json",
			v: func() openai.OpenAITranslationParams {
				temp := 0.0
				return openai.OpenAITranslationParams{
					File:           "base64data",
					Model:          "whisper-1",
					Prompt:         "Translate to English",
					Temperature:    &temp,
					ResponseFormat: openai.OpenAITranslationResponseFormatVerboseJSON,
				}
			}(),
			want: `{"file":"base64data","model":"whisper-1","prompt":"Translate to English","temperature":0,"response_format":"verbose_json"}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === B. AudioTranslation 响应 union ===

func TestWireFormatAudioTranslationResponse(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Response_basic",
			v: openai.OpenAITranslation{
				Text: "Hello, world!",
			},
			want: `{"text":"Hello, world!"}`,
		},
		{
			name: "Response_verbose",
			v: func() openai.OpenAITranslationVerbose {
				dur := 2.5
				return openai.OpenAITranslationVerbose{
					Text:     "Hello, world!",
					Language: "zh",
					Duration: &dur,
					Segments: []openai.OpenAITranslationSegment{
						{ID: 0, Start: 0.0, End: 2.5, Text: "你好世界"},
					},
				}
			}(),
			want: `{"text":"Hello, world!","language":"zh","duration":2.5,"segments":[{"id":0,"seek":0,"start":0,"end":2.5,"text":"你好世界","tokens":null,"temperature":0,"avg_logprob":0,"compression_ratio":0,"no_speech_prob":0}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %s\nwant %s", got, tc.want)
			}
		})
	}
}

// === C. AudioTranslation union 封闭性 ===

func TestAudioTranslationResponseUnionClosure(t *testing.T) {
	var _ openai.OpenAITranslationResponseUnion = openai.OpenAITranslation{}
	var _ openai.OpenAITranslationResponseUnion = openai.OpenAITranslationVerbose{}
}

// === D. AudioTranslation 5 种 ResponseFormat 命名常量 ===

func TestAudioTranslationResponseFormatNames(t *testing.T) {
	formats := []openai.OpenAITranslationResponseFormat{
		openai.OpenAITranslationResponseFormatJSON,
		openai.OpenAITranslationResponseFormatText,
		openai.OpenAITranslationResponseFormatSRT,
		openai.OpenAITranslationResponseFormatVerboseJSON,
		openai.OpenAITranslationResponseFormatVTT,
	}
	want := []string{"json", "text", "srt", "verbose_json", "vtt"}
	if len(formats) != len(want) {
		t.Fatalf("got %d formats, want %d", len(formats), len(want))
	}
	for i, f := range formats {
		if string(f) != want[i] {
			t.Errorf("formats[%d] = %q, want %q", i, string(f), want[i])
		}
	}
}
