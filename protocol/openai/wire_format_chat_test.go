package openai_test

import (
	"encoding/json"
	"testing"

	openai "hegq.net/mogo/pkg/protocol/openai"
)

// === A. ChatCompletionParams 顶层字段 marshal 测试（8 case） ===

func TestWireFormatChatCompletionParams(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Params_minimal",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{
						Role:    openai.OpenAIChatCompletionMessageRoleUser,
						Content: openai.OpenAIChatCompletionUserMessageContentString("hello"),
					},
				},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}]}`,
		},
		{
			name: "Params_with_sampling",
			v: func() openai.OpenAIChatCompletionParams {
				temp := 0.7
				topP := 0.9
				freq := 0.5
				pres := 0.5
				return openai.OpenAIChatCompletionParams{
					Model:            "gpt-4o",
					Messages:         []openai.OpenAIChatCompletionMessageParamUnion{openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")}},
					Temperature:      &temp,
					TopP:             &topP,
					FrequencyPenalty: &freq,
					PresencePenalty:  &pres,
				}
			}(),
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"temperature":0.7,"top_p":0.9,"frequency_penalty":0.5,"presence_penalty":0.5}`,
		},
		{
			name: "Params_with_tools",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				Tools: []openai.OpenAIChatCompletionToolUnion{
					openai.OpenAIChatCompletionToolFunctionParam{
						Type: "function",
						Function: openai.OpenAIChatCompletionToolFunctionDefinition{
							Name:        "get_weather",
							Description: "Get the weather",
							Parameters:  json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"}}}`),
						},
					},
				},
				ToolChoice: openai.OpenAIChatCompletionToolChoiceOptionAutoParam{Type: "auto"},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"tools":[{"type":"function","function":{"name":"get_weather","description":"Get the weather","parameters":{"type":"object","properties":{"city":{"type":"string"}}}}}],"tool_choice":{"type":"auto"}}`,
		},
		{
			name: "Params_with_response_format_json_schema",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				ResponseFormat: openai.OpenAIResponseFormatJSONSchemaParam{
					Type: openai.OpenAIResponseFormatTypeJSONSchema,
					JSONSchema: openai.OpenAIJSONSchema{
						Name:   "weather",
						Schema: json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"}}}`),
					},
				},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"response_format":{"type":"json_schema","json_schema":{"name":"weather","schema":{"type":"object","properties":{"city":{"type":"string"}}}}}}`,
		},
		{
			name: "Params_with_audio",
			v: func() openai.OpenAIChatCompletionParams {
				return openai.OpenAIChatCompletionParams{
					Model: "gpt-4o-audio",
					Messages: []openai.OpenAIChatCompletionMessageParamUnion{
						openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
					},
					Audio: &openai.OpenAIChatCompletionAudioParam{
						Voice:  openai.OpenAIAudioVoiceUnionString(openai.OpenAIAudioVoiceAlloy),
						Format: openai.OpenAIAudioFormatWAV,
					},
					Modalities: []openai.OpenAIModality{openai.OpenAIModalityText, openai.OpenAIModalityAudio},
				}
			}(),
			want: `{"model":"gpt-4o-audio","messages":[{"role":"user","content":"hi"}],"audio":{"voice":"alloy","format":"wav"},"modalities":["text","audio"]}`,
		},
		{
			name: "Params_with_streaming",
			v: func() openai.OpenAIChatCompletionParams {
				s := true
				inc := true
				return openai.OpenAIChatCompletionParams{
					Model: "gpt-4o",
					Messages: []openai.OpenAIChatCompletionMessageParamUnion{
						openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
					},
					Stream:        &s,
					StreamOptions: &openai.OpenAIChatCompletionStreamOptionsParam{IncludeUsage: &inc},
				}
			}(),
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"stream":true,"stream_options":{"include_usage":true}}`,
		},
		{
			name: "Params_with_reasoning",
			v: openai.OpenAIChatCompletionParams{
				Model: "o1",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				ReasoningEffort: openai.OpenAIReasoningEffortHigh,
				Verbosity:       openai.OpenAIVerbosityMedium,
			},
			want: `{"model":"o1","messages":[{"role":"user","content":"hi"}],"reasoning_effort":"high","verbosity":"medium"}`,
		},
		{
			name: "Params_with_web_search_options",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				WebSearchOptions: &openai.OpenAIChatCompletionWebSearchOptions{
					SearchContextSize: "high",
				},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"web_search_options":{"search_context_size":"high"}}`,
		},
		{
			name: "Params_with_web_search_filters",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				WebSearchOptions: &openai.OpenAIChatCompletionWebSearchOptions{
					SearchContextSize: "high",
					Filters:           &openai.OpenAIChatCompletionWebSearchFilters{AllowDomains: []string{"example.com", "anthropic.com"}},
				},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"web_search_options":{"search_context_size":"high","filters":{"allow_domains":["example.com","anthropic.com"]}}}`,
		},
		{
			name: "Params_with_prompt_cache_retention",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				PromptCacheRetention: openai.OpenAIResponsePromptCacheRetention24h,
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"prompt_cache_retention":"24h"}`,
		},
		{
			name: "Params_with_moderation",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				Moderation: &openai.OpenAIChatCompletionModerationParam{Model: "omni-moderation-latest"},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"moderation":{"model":"omni-moderation-latest"}}`,
		},
		{
			name: "Params_with_include_obfuscation",
			v: func() openai.OpenAIChatCompletionParams {
				s := true
				inc := true
				return openai.OpenAIChatCompletionParams{
					Model: "gpt-4o",
					Messages: []openai.OpenAIChatCompletionMessageParamUnion{
						openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
					},
					Stream:             &s,
					StreamOptions:      &openai.OpenAIChatCompletionStreamOptionsParam{IncludeObfuscation: &inc},
					IncludeObfuscation: &inc,
				}
			}(),
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"stream":true,"stream_options":{"include_obfuscation":true},"include_obfuscation":true}`,
		},
		{
			name: "Params_with_image_generation_tool",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("Draw a cat")},
				},
				Tools: []openai.OpenAIChatCompletionToolUnion{
					openai.OpenAIChatCompletionToolImageGenerationParam{
						Type:         openai.OpenAIChatCompletionToolImageGenerationTypeImageGeneration,
						Quality:      "high",
						Size:         "1024x1024",
						OutputFormat: "png",
					},
				},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"Draw a cat"}],"tools":[{"type":"image_generation","quality":"high","size":"1024x1024","output_format":"png"}]}`,
		},
		{
			name: "Params_with_tool_choice_custom",
			v: openai.OpenAIChatCompletionParams{
				Model: "gpt-4o",
				Messages: []openai.OpenAIChatCompletionMessageParamUnion{
					openai.OpenAIChatCompletionUserMessageParam{Role: openai.OpenAIChatCompletionMessageRoleUser, Content: openai.OpenAIChatCompletionUserMessageContentString("hi")},
				},
				ToolChoice: openai.OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam{
					Type:   "custom",
					Custom: openai.OpenAIChatCompletionToolChoiceOptionHostedNamedCustom{Name: "my_tool"},
				},
			},
			want: `{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}],"tool_choice":{"type":"custom","custom":{"name":"my_tool"}}}`,
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

// === B. MessageParam union marshal 测试（6 case） ===

func TestWireFormatChatCompletionMessageParamUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "UserMessage_text",
			v: openai.OpenAIChatCompletionUserMessageParam{
				Role:    openai.OpenAIChatCompletionMessageRoleUser,
				Content: openai.OpenAIChatCompletionUserMessageContentString("hi"),
			},
			want: `{"role":"user","content":"hi"}`,
		},
		{
			name: "UserMessage_multimodal",
			v: openai.OpenAIChatCompletionUserMessageParam{
				Role: openai.OpenAIChatCompletionMessageRoleUser,
				Content: openai.OpenAIChatCompletionUserMessageContentParts{
					openai.OpenAIChatCompletionContentPartTextParam{Type: openai.OpenAIChatCompletionContentPartTypeText, Text: "what's in this image?"},
					openai.OpenAIChatCompletionContentPartImageParam{
						Type:     openai.OpenAIChatCompletionContentPartTypeImageURL,
						ImageURL: openai.OpenAIChatCompletionContentPartImageURL{URL: "https://example.com/cat.jpg"},
					},
				},
			},
			want: `{"role":"user","content":[{"type":"text","text":"what's in this image?"},{"type":"image_url","image_url":{"url":"https://example.com/cat.jpg"}}]}`,
		},
		{
			name: "AssistantMessage_with_tool_calls",
			v: openai.OpenAIChatCompletionAssistantMessageParam{
				Role:    openai.OpenAIChatCompletionMessageRoleAssistant,
				Content: ptrString("Let me check the weather."),
				ToolCalls: []openai.OpenAIChatCompletionMessageToolCallUnion{
					openai.OpenAIChatCompletionMessageToolCallFunctionParam{
						ID:   "call_123",
						Type: "function",
						Function: openai.OpenAIChatCompletionMessageToolCallFunction{
							Name:      "get_weather",
							Arguments: `{"city":"SF"}`,
						},
					},
				},
			},
			want: `{"role":"assistant","content":"Let me check the weather.","tool_calls":[{"id":"call_123","type":"function","function":{"name":"get_weather","arguments":"{\"city\":\"SF\"}"}}]}`,
		},
		{
			name: "AssistantMessage_with_audio",
			v: openai.OpenAIChatCompletionAssistantMessageParam{
				Role:    openai.OpenAIChatCompletionMessageRoleAssistant,
				Content: ptrString(""),
				Audio: &openai.OpenAIChatCompletionMessageAudioParam{
					ID:         "audio_123",
					Data:       "base64data",
					ExpiresAt:  1234567890,
					Transcript: "Hello there.",
				},
			},
			want: `{"role":"assistant","content":"","audio":{"id":"audio_123","data":"base64data","expires_at":1234567890,"transcript":"Hello there."}}`,
		},
		{
			name: "SystemMessage_string",
			v: openai.OpenAIChatCompletionSystemMessageParam{
				Role:    openai.OpenAIChatCompletionMessageRoleSystem,
				Content: openai.OpenAIChatCompletionSystemMessageContentString("You are a helpful assistant."),
			},
			want: `{"role":"system","content":"You are a helpful assistant."}`,
		},
		{
			name: "ToolMessage_with_error",
			v: openai.OpenAIChatCompletionToolMessageParam{
				Role:       openai.OpenAIChatCompletionMessageRoleTool,
				Content:    openai.OpenAIChatCompletionToolMessageContentString("Tool not found"),
				ToolCallID: "call_123",
			},
			want: `{"role":"tool","content":"Tool not found","tool_call_id":"call_123"}`,
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

// === C. ContentPart union marshal 测试（7 case，Annotation 在响应侧） ===

func TestWireFormatChatCompletionContentPartUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Text",
			v: openai.OpenAIChatCompletionContentPartTextParam{
				Type: openai.OpenAIChatCompletionContentPartTypeText,
				Text: "hello",
			},
			want: `{"type":"text","text":"hello"}`,
		},
		{
			name: "Image",
			v: openai.OpenAIChatCompletionContentPartImageParam{
				Type:     openai.OpenAIChatCompletionContentPartTypeImageURL,
				ImageURL: openai.OpenAIChatCompletionContentPartImageURL{URL: "https://example.com/cat.jpg", Detail: openai.OpenAIImageDetailAuto},
			},
			want: `{"type":"image_url","image_url":{"url":"https://example.com/cat.jpg","detail":"auto"}}`,
		},
		{
			name: "InputAudio",
			v: openai.OpenAIChatCompletionContentPartInputAudioParam{
				Type: openai.OpenAIChatCompletionContentPartTypeInputAudio,
				InputAudio: openai.OpenAIChatCompletionContentPartInputAudio{
					Data:   "base64data",
					Format: openai.OpenAIAudioFormatWAV,
				},
			},
			want: `{"type":"input_audio","input_audio":{"data":"base64data","format":"wav"}}`,
		},
		{
			name: "File",
			v: openai.OpenAIChatCompletionContentPartFileParam{
				Type: openai.OpenAIChatCompletionContentPartTypeFile,
				File: openai.OpenAIChatCompletionContentPartFile{
					FileData: "base64data",
					Filename: "doc.pdf",
				},
			},
			want: `{"type":"file","file":{"file_data":"base64data","filename":"doc.pdf"}}`,
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

// === D. Tool union marshal 测试（6 case） ===

func TestWireFormatChatCompletionToolUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Function",
			v: openai.OpenAIChatCompletionToolFunctionParam{
				Type: "function",
				Function: openai.OpenAIChatCompletionToolFunctionDefinition{
					Name:        "get_weather",
					Description: "Get the weather",
					Parameters:  json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"}},"required":["city"]}`),
				},
			},
			want: `{"type":"function","function":{"name":"get_weather","description":"Get the weather","parameters":{"type":"object","properties":{"city":{"type":"string"}},"required":["city"]}}}`,
		},
		{
			name: "Custom",
			v: openai.OpenAIChatCompletionToolCustomParam{
				Type: "custom",
				Custom: openai.OpenAIChatCompletionToolCustomDefinition{
					Name:        "my_tool",
					Description: "Custom tool",
				},
			},
			want: `{"type":"custom","custom":{"name":"my_tool","description":"Custom tool"}}`,
		},
		{
			name: "WebSearch",
			v: openai.OpenAIChatCompletionToolWebSearchParam{
				Type: "web_search",
			},
			want: `{"type":"web_search"}`,
		},
		{
			name: "CodeInterpreter",
			v: openai.OpenAIChatCompletionToolCodeInterpreterParam{
				Type:      "code_interpreter",
				Container: "container_123",
			},
			want: `{"type":"code_interpreter","container":"container_123"}`,
		},
		{
			name: "FileSearch",
			v: func() openai.OpenAIChatCompletionToolFileSearchParam {
				max := int64(5)
				return openai.OpenAIChatCompletionToolFileSearchParam{
					Type:           "file_search",
					VectorStoreIDs: []string{"vs_1"},
					MaxNumResults:  &max,
					RankingOptions: &openai.OpenAIChatCompletionToolFileSearchRanking{Ranker: "auto", ScoreThreshold: 0.5},
				}
			}(),
			want: `{"type":"file_search","vector_store_ids":["vs_1"],"max_num_results":5,"ranking_options":{"ranker":"auto","score_threshold":0.5}}`,
		},
		{
			name: "ComputerUse",
			v: openai.OpenAIChatCompletionToolComputerUseParam{
				Type:          "computer_use_preview",
				DisplayWidth:  1024,
				DisplayHeight: 768,
				Environment:   "browser",
			},
			want: `{"type":"computer_use_preview","display_width":1024,"display_height":768,"environment":"browser"}`,
		},
		{
			name: "ImageGeneration",
			v: func() openai.OpenAIChatCompletionToolImageGenerationParam {
				comp := int64(85)
				parts := int64(3)
				return openai.OpenAIChatCompletionToolImageGenerationParam{
					Type:              openai.OpenAIChatCompletionToolImageGenerationTypeImageGeneration,
					Background:        "auto",
					InputFidelity:     "high",
					Model:             "gpt-image-1",
					Moderation:        "auto",
					OutputCompression: &comp,
					OutputFormat:      "png",
					PartialImages:     &parts,
					Quality:           "high",
					Size:              "1024x1024",
				}
			}(),
			want: `{"type":"image_generation","background":"auto","input_fidelity":"high","model":"gpt-image-1","moderation":"auto","output_compression":85,"partial_images":3,"quality":"high","size":"1024x1024","output_format":"png"}`,
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

// === E. ToolChoice union marshal 测试（4 case） ===

func TestWireFormatChatCompletionToolChoiceUnion(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Auto",
			v:    openai.OpenAIChatCompletionToolChoiceOptionAutoParam{Type: "auto"},
			want: `{"type":"auto"}`,
		},
		{
			name: "Required",
			v:    openai.OpenAIChatCompletionToolChoiceOptionRequiredParam{Type: "required"},
			want: `{"type":"required"}`,
		},
		{
			name: "None",
			v:    openai.OpenAIChatCompletionToolChoiceOptionNoneParam{Type: "none"},
			want: `{"type":"none"}`,
		},
		{
			name: "Allowed",
			v: openai.OpenAIChatCompletionToolChoiceOptionAllowedParam{
				Type: "allowed_tools",
				AllowedTools: []openai.OpenAIChatCompletionToolChoiceOptionAllowedTool{
					{
						Type:     "function",
						Function: openai.OpenAIChatCompletionToolChoiceOptionAllowedFunction{Name: "get_weather"},
					},
				},
				Mode: "auto",
			},
			want: `{"type":"allowed_tools","allowed_tools":[{"type":"function","function":{"name":"get_weather"}}],"mode":"auto"}`,
		},
		{
			name: "HostedNamedCustom",
			v: openai.OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam{
				Type:   "custom",
				Custom: openai.OpenAIChatCompletionToolChoiceOptionHostedNamedCustom{Name: "my_custom_tool"},
			},
			want: `{"type":"custom","custom":{"name":"my_custom_tool"}}`,
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

// === F. ChatCompletion 响应 marshal 测试（5 case） ===

func TestWireFormatChatCompletionResponse(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Response_minimal",
			v: openai.OpenAIChatCompletion{
				ID:      "chatcmpl-123",
				Object:  "chat.completion",
				Created: 1234567890,
				Model:   "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChoice{
					{
						Index: 0,
						Message: openai.OpenAIChatCompletionChoiceMessage{
							Role:    "assistant",
							Content: "Hello there.",
						},
						FinishReason: "stop",
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"Hello there."},"finish_reason":"stop"}]}`,
		},
		{
			name: "Response_with_audio_output",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "gpt-4o-audio",
				Choices: []openai.OpenAIChatCompletionChoice{
					{
						Index: 0,
						Message: openai.OpenAIChatCompletionChoiceMessage{
							Role: "assistant", Content: "",
							Audio: &openai.OpenAIChatCompletionChoiceMessageAudio{
								ID: "audio_1", Data: "data", ExpiresAt: 1234567890, Transcript: "hi",
							},
						},
						FinishReason: "stop",
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o-audio","choices":[{"index":0,"message":{"role":"assistant","content":"","audio":{"id":"audio_1","data":"data","expires_at":1234567890,"transcript":"hi"}},"finish_reason":"stop"}]}`,
		},
		{
			name: "Response_with_url_citation",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "gpt-4o-search",
				Choices: []openai.OpenAIChatCompletionChoice{
					{
						Index: 0,
						Message: openai.OpenAIChatCompletionChoiceMessage{
							Role: "assistant", Content: "According to the source.",
							Annotations: []openai.OpenAIChatCompletionChoiceMessageAnnotation{
								{
									Type: "url_citation",
									URLCitation: &openai.OpenAIChatCompletionMessageAnnotationURLCitation{
										URL: "https://example.com", Title: "Example", StartIndex: 0, EndIndex: 10,
									},
								},
							},
						},
						FinishReason: "stop",
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o-search","choices":[{"index":0,"message":{"role":"assistant","content":"According to the source.","annotations":[{"type":"url_citation","url_citation":{"url":"https://example.com","title":"Example","start_index":0,"end_index":10}}]},"finish_reason":"stop"}]}`,
		},
		{
			name: "Response_with_usage_breakdown",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChoice{
					{Index: 0, Message: openai.OpenAIChatCompletionChoiceMessage{Role: "assistant", Content: "hi"}, FinishReason: "stop"},
				},
				Usage: &openai.OpenAICompletionUsage{
					PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15,
					PromptTokensDetails: &openai.OpenAICompletionUsagePromptTokensDetails{CachedTokens: 8},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"hi"},"finish_reason":"stop"}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15,"prompt_tokens_details":{"cached_tokens":8}}}`,
		},
		{
			name: "Response_with_service_tier",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "gpt-4o",
				Choices:     []openai.OpenAIChatCompletionChoice{{Index: 0, Message: openai.OpenAIChatCompletionChoiceMessage{Role: "assistant", Content: "hi"}, FinishReason: "stop"}},
				ServiceTier: openai.OpenAIServiceTierPriority,
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"hi"},"finish_reason":"stop"}],"service_tier":"priority"}`,
		},
		{
			name: "Response_with_reasoning",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "deepseek-reasoner",
				Choices: []openai.OpenAIChatCompletionChoice{
					{
						Index: 0,
						Message: openai.OpenAIChatCompletionChoiceMessage{
							Role:             "assistant",
							Content:          "The answer is 42.",
							Reasoning:        "Let me think about this...",
							ReasoningContent: "Let me think about this...",
						},
						FinishReason: "stop",
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"deepseek-reasoner","choices":[{"index":0,"message":{"role":"assistant","content":"The answer is 42.","reasoning":"Let me think about this...","reasoning_content":"Let me think about this..."},"finish_reason":"stop"}]}`,
		},
		{
			name: "Response_with_content_filter",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChoice{
					{
						Index: 0,
						Message: openai.OpenAIChatCompletionChoiceMessage{
							Role: "assistant", Content: "hi",
						},
						FinishReason: "stop",
						ContentFilterResults: &openai.OpenAIChatCompletionContentFilterResults{
							Hate:            &openai.OpenAIChatCompletionContentFilterResult{Filtered: false, Severity: "safe"},
							HateThreatening: &openai.OpenAIChatCompletionContentFilterResult{Filtered: false, Severity: "safe"},
							SelfHarm:        &openai.OpenAIChatCompletionContentFilterResult{Filtered: false, Severity: "safe"},
							Sexual:          &openai.OpenAIChatCompletionContentFilterResult{Filtered: false, Severity: "safe"},
							Violence:        &openai.OpenAIChatCompletionContentFilterResult{Filtered: false, Severity: "safe"},
						},
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"hi"},"finish_reason":"stop","content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"hate_threatening":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}]}`,
		},
		{
			name: "Response_with_moderation",
			v: openai.OpenAIChatCompletion{
				ID: "chatcmpl-123", Object: "chat.completion", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChoice{
					{Index: 0, Message: openai.OpenAIChatCompletionChoiceMessage{Role: "assistant", Content: "hi"}, FinishReason: "stop"},
				},
				Moderation: &openai.OpenAIChatCompletionResponseModeration{
					Input:  &openai.OpenAIChatCompletionModerationInput{Flagged: false, Model: "omni-moderation-latest"},
					Output: &openai.OpenAIChatCompletionModerationOutput{Flagged: true, Model: "omni-moderation-latest"},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"message":{"role":"assistant","content":"hi"},"finish_reason":"stop"}],"moderation":{"input":{"flagged":false,"model":"omni-moderation-latest"},"output":{"flagged":true,"model":"omni-moderation-latest"}}}`,
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

// === G. Chunk 流式 marshal 测试（4 case） ===

func TestWireFormatChatCompletionChunk(t *testing.T) {
	stop := "stop"
	cases := []struct {
		name string
		v    any
		want string
	}{
		{
			name: "Chunk_role_chunk",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{Index: 0, Delta: openai.OpenAIChatCompletionChunkChoiceDelta{Role: "assistant"}},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant"}}]}`,
		},
		{
			name: "Chunk_content_delta",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{Index: 0, Delta: openai.OpenAIChatCompletionChunkChoiceDelta{Content: "hel"}},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"hel"}}]}`,
		},
		{
			name: "Chunk_tool_call_delta",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: openai.OpenAIChatCompletionChunkChoiceDelta{
							ToolCalls: []openai.OpenAIChatCompletionChunkChoiceDeltaToolCall{
								{Index: 0, ID: "call_123", Type: "function", Function: &openai.OpenAIChatCompletionMessageToolCallFunction{Name: "get_weather", Arguments: `{"city"`}},
							},
						},
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_123","type":"function","function":{"name":"get_weather","arguments":"{\"city\""}}]}}]}`,
		},
		{
			name: "Chunk_finish_reason",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{Index: 0, Delta: openai.OpenAIChatCompletionChunkChoiceDelta{}, FinishReason: &stop},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
		},
		{
			name: "Chunk_reasoning_delta",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "deepseek-reasoner",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{Index: 0, Delta: openai.OpenAIChatCompletionChunkChoiceDelta{
						Role:             "assistant",
						ReasoningContent: "Let me think...",
					}},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"deepseek-reasoner","choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"Let me think..."}}]}`,
		},
		{
			name: "Chunk_content_filter_results",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{
						Index: 0,
						Delta: openai.OpenAIChatCompletionChunkChoiceDelta{Content: "hi"},
						ContentFilterResults: &openai.OpenAIChatCompletionContentFilterResults{
							Hate: &openai.OpenAIChatCompletionContentFilterResult{Filtered: false, Severity: "safe"},
						},
					},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"hi"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"}}}]}`,
		},
		{
			name: "Chunk_moderation",
			v: openai.OpenAIChatCompletionChunk{
				ID: "chatcmpl-123", Object: "chat.completion.chunk", Created: 1234567890, Model: "gpt-4o",
				Choices: []openai.OpenAIChatCompletionChunkChoice{
					{Index: 0, Delta: openai.OpenAIChatCompletionChunkChoiceDelta{Content: "hi"}, FinishReason: &stop},
				},
				Moderation: &openai.OpenAIChatCompletionResponseModeration{
					Input: &openai.OpenAIChatCompletionModerationInput{Flagged: false, Model: "omni-moderation-latest"},
				},
			},
			want: `{"id":"chatcmpl-123","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":"hi"},"finish_reason":"stop"}],"moderation":{"input":{"flagged":false,"model":"omni-moderation-latest"}}}`,
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

// === H. Union 封闭性测试 ===

func TestChatCompletionMessageParamUnionClosure(t *testing.T) {
	var _ openai.OpenAIChatCompletionMessageParamUnion = openai.OpenAIChatCompletionUserMessageParam{}
	var _ openai.OpenAIChatCompletionMessageParamUnion = openai.OpenAIChatCompletionAssistantMessageParam{}
	var _ openai.OpenAIChatCompletionMessageParamUnion = openai.OpenAIChatCompletionSystemMessageParam{}
	var _ openai.OpenAIChatCompletionMessageParamUnion = openai.OpenAIChatCompletionDeveloperMessageParam{}
	var _ openai.OpenAIChatCompletionMessageParamUnion = openai.OpenAIChatCompletionToolMessageParam{}
	var _ openai.OpenAIChatCompletionMessageParamUnion = openai.OpenAIChatCompletionFunctionMessageParam{}
}

func TestChatCompletionContentPartUnionClosure(t *testing.T) {
	var _ openai.OpenAIChatCompletionContentPartUnion = openai.OpenAIChatCompletionContentPartTextParam{}
	var _ openai.OpenAIChatCompletionContentPartUnion = openai.OpenAIChatCompletionContentPartImageParam{}
	var _ openai.OpenAIChatCompletionContentPartUnion = openai.OpenAIChatCompletionContentPartInputAudioParam{}
	var _ openai.OpenAIChatCompletionContentPartUnion = openai.OpenAIChatCompletionContentPartFileParam{}
}

func TestChatCompletionToolUnionClosure(t *testing.T) {
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolFunctionParam{}
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolCustomParam{}
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolWebSearchParam{}
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolCodeInterpreterParam{}
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolFileSearchParam{}
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolComputerUseParam{}
	var _ openai.OpenAIChatCompletionToolUnion = openai.OpenAIChatCompletionToolImageGenerationParam{}
}

func TestChatCompletionToolChoiceUnionClosure(t *testing.T) {
	var _ openai.OpenAIChatCompletionToolChoiceOptionUnionParam = openai.OpenAIChatCompletionToolChoiceOptionAutoParam{}
	var _ openai.OpenAIChatCompletionToolChoiceOptionUnionParam = openai.OpenAIChatCompletionToolChoiceOptionRequiredParam{}
	var _ openai.OpenAIChatCompletionToolChoiceOptionUnionParam = openai.OpenAIChatCompletionToolChoiceOptionNoneParam{}
	var _ openai.OpenAIChatCompletionToolChoiceOptionUnionParam = openai.OpenAIChatCompletionToolChoiceOptionAllowedParam{}
	var _ openai.OpenAIChatCompletionToolChoiceOptionUnionParam = openai.OpenAIChatCompletionToolChoiceOptionHostedNamedParam{}
	var _ openai.OpenAIChatCompletionToolChoiceOptionUnionParam = openai.OpenAIChatCompletionToolChoiceOptionHostedNamedCustomParam{}
}

func TestChatCompletionStreamChunkWithModeration(t *testing.T) {
	jsonStr := `{"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"hello"},"finish_reason":null}],"moderation":{"id":"mod_xxx","output":"safe","version":"2024-05"}}`
	var chunk openai.OpenAIChatCompletionStreamChunkEvent
	err := json.Unmarshal([]byte(jsonStr), &chunk)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if chunk.Moderation == nil {
		t.Fatalf("expected Moderation not nil")
	}
	if chunk.Moderation.ID != "mod_xxx" {
		t.Errorf("got Moderation.ID %q, want %q", chunk.Moderation.ID, "mod_xxx")
	}
}

func TestChatCompletionStreamChunkWithServiceTier(t *testing.T) {
	jsonStr := `{"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"content":""}}],"service_tier":"auto"}`
	var chunk openai.OpenAIChatCompletionStreamChunkEvent
	err := json.Unmarshal([]byte(jsonStr), &chunk)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if chunk.ServiceTier != "auto" {
		t.Errorf("got ServiceTier %q, want %q", chunk.ServiceTier, "auto")
	}
}

func TestChatCompletionStreamChunkWithLogprobs(t *testing.T) {
	jsonStr := `{"id":"chatcmpl-xxx","object":"chat.completion.chunk","created":1234567890,"model":"gpt-4","choices":[{"index":0,"delta":{"content":"hi"},"logprobs":{"content":[{"token":"hi","logprob":-0.5,"prob":0.6}]}}]}`
	var chunk openai.OpenAIChatCompletionStreamChunkEvent
	err := json.Unmarshal([]byte(jsonStr), &chunk)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if chunk.Choices[0].Logprobs == nil {
		t.Fatalf("expected Logprobs not nil")
	}
	if chunk.Choices[0].Logprobs.Content[0].Token != "hi" {
		t.Errorf("got Content[0].Token %q, want %q", chunk.Choices[0].Logprobs.Content[0].Token, "hi")
	}
}

// 辅助函数：返回字符串指针
func ptrString(s string) *string {
	return &s
}

// 辅助函数：返回 bool 指针
func ptrBool(b bool) *bool { return &b }

// 辅助函数：返回 int64 指针
func ptrInt64(i int64) *int64 { return &i }

// 辅助函数：返回 float64 指针
func ptrFloat64(f float64) *float64 { return &f }

// === Z. 枚举全量命名常量测试（4.6+ 扩展） ===

func TestWireFormatEnumExtensions(t *testing.T) {
	t.Run("ReasoningEffort_full_set", func(t *testing.T) {
		efforts := []openai.OpenAIReasoningEffort{
			openai.OpenAIReasoningEffortNone,
			openai.OpenAIReasoningEffortMinimal,
			openai.OpenAIReasoningEffortLow,
			openai.OpenAIReasoningEffortMedium,
			openai.OpenAIReasoningEffortHigh,
			openai.OpenAIReasoningEffortXHigh,
		}
		want := []string{"none", "minimal", "low", "medium", "high", "xhigh"}
		if len(efforts) != len(want) {
			t.Fatalf("got %d, want %d", len(efforts), len(want))
		}
		for i, e := range efforts {
			if string(e) != want[i] {
				t.Errorf("[%d] got %q want %q", i, string(e), want[i])
			}
		}
	})

	t.Run("ServiceTier_full_set", func(t *testing.T) {
		tiers := []openai.OpenAIServiceTier{
			openai.OpenAIServiceTierAuto,
			openai.OpenAIServiceTierDefault,
			openai.OpenAIServiceTierPriority,
			openai.OpenAIServiceTierFlex,
			openai.OpenAIServiceTierScale,
		}
		want := []string{"auto", "default", "priority", "flex", "scale"}
		if len(tiers) != len(want) {
			t.Fatalf("got %d, want %d", len(tiers), len(want))
		}
		for i, v := range tiers {
			if string(v) != want[i] {
				t.Errorf("[%d] got %q want %q", i, string(v), want[i])
			}
		}
	})

	t.Run("AudioVoice_includes_marin_cedar", func(t *testing.T) {
		if string(openai.OpenAIAudioVoiceMarin) != "marin" {
			t.Errorf("marin = %q", string(openai.OpenAIAudioVoiceMarin))
		}
		if string(openai.OpenAIAudioVoiceCedar) != "cedar" {
			t.Errorf("cedar = %q", string(openai.OpenAIAudioVoiceCedar))
		}
	})

	t.Run("ResponseToolChoiceMode_full_set", func(t *testing.T) {
		modes := []openai.OpenAIResponseToolChoiceOptionMode{
			openai.OpenAIResponseToolChoiceOptionModeAuto,
			openai.OpenAIResponseToolChoiceOptionModeNone,
			openai.OpenAIResponseToolChoiceOptionModeRequired,
		}
		want := []string{"auto", "none", "required"}
		for i, m := range modes {
			if string(m) != want[i] {
				t.Errorf("[%d] got %q want %q", i, string(m), want[i])
			}
		}
	})
}
