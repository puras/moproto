package usage

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkganthropic "hegq.net/mogo/pkg/protocol/anthropic"
	pkgopenai "hegq.net/mogo/pkg/protocol/openai"
)

// ============== P0-1: Base 字段命名统一（Anthropic 风格为权威） ==============

func TestBase_NoAliasFields(t *testing.T) {
	// 验证 Base 只包含 InputTokens / OutputTokens / TotalTokens / FreshInputTokens，
	// 不存在 PromptTokens / CompletionTokens 别名（编译期+反射双保险）。
	//
	// 2026-08-25 PR 改造新增 FreshInputTokens（计费输入，去 cache 后口径）。
	baseType := reflect.TypeOf(Base{})
	want := []string{"InputTokens", "OutputTokens", "TotalTokens", "FreshInputTokens"}
	if got := baseType.NumField(); got != len(want) {
		t.Fatalf("Base 字段数=%d, want %d（发现意外字段）", got, len(want))
	}
	for i, name := range want {
		if got := baseType.Field(i).Name; got != name {
			t.Errorf("Base 第 %d 个字段=%q, want %q", i, got, name)
		}
	}
	// 同时确保 Usage 顶层也不暴露别名
	usageType := reflect.TypeOf(Usage{})
	for i := 0; i < usageType.NumField(); i++ {
		name := usageType.Field(i).Name
		if name == "PromptTokens" || name == "CompletionTokens" {
			t.Errorf("Usage 不应暴露别名字段 %s", name)
		}
	}
}

// ============== P0-2: TotalTokens 填充语义 ==============

func TestFromAnthropic_TotalTokens_AutoSum(t *testing.T) {
	// Anthropic 不返回 TotalTokens，转换函数自填为 InputTokens + OutputTokens
	src := &pkganthropic.AnthropicUsage{
		InputTokens:  200,
		OutputTokens: 80,
	}
	got := FromAnthropic(src)
	if got.TotalTokens != 280 {
		t.Errorf("Anthropic → TotalTokens 应为 InputTokens+OutputTokens=280，实际 %d", got.TotalTokens)
	}
	if got.InputTokens != 200 || got.OutputTokens != 80 {
		t.Errorf("Base 字段未正确映射")
	}
	if got.Protocol != "anthropic" {
		t.Errorf("Protocol 应为 anthropic，实际 %q", got.Protocol)
	}
}

func TestFromOpenAICompletion_TotalTokens_Passthrough(t *testing.T) {
	// OpenAI 直接给出 TotalTokens，转换函数原样透传（不自动重算）
	src := &pkgopenai.OpenAICompletionUsage{
		PromptTokens:     100,
		CompletionTokens: 50,
		// OpenAI 上游 TotalTokens 包含 cached_tokens，与 InputTokens+OutputTokens 不一定相等
		TotalTokens: 175,
	}
	got := FromOpenAICompletion(src)
	if got.InputTokens != 100 || got.OutputTokens != 50 {
		t.Errorf("OpenAI → Base 字段映射错误: Input=%d Output=%d", got.InputTokens, got.OutputTokens)
	}
	if got.TotalTokens != 175 {
		t.Errorf("OpenAI → TotalTokens 应原样透传 175，实际 %d", got.TotalTokens)
	}
	if got.Protocol != "openai" {
		t.Errorf("Protocol 应为 openai，实际 %q", got.Protocol)
	}
}

func TestFromAnthropic_Nil(t *testing.T) {
	if got := FromAnthropic(nil); got != nil {
		t.Errorf("FromAnthropic(nil) 应返回 nil")
	}
}

func TestFromOpenAICompletion_Nil(t *testing.T) {
	if got := FromOpenAICompletion(nil); got != nil {
		t.Errorf("FromOpenAICompletion(nil) 应返回 nil")
	}
}

// ============== IsZero 完整遍历分组 ==============

func TestUsage_IsZero(t *testing.T) {
	tests := []struct {
		name string
		u    *Usage
		want bool
	}{
		{"nil receiver", nil, true},
		{"zero value", &Usage{}, true},
		{"only base input", &Usage{Base: Base{InputTokens: 1}}, false},
		{"only cache", &Usage{Cache: &CacheMetrics{CachedTokens: 100}}, false},
		{"only reasoning", &Usage{Reasoning: &ReasoningTokens{ReasoningTokens: 5}}, false},
		{"only thinking", &Usage{Thinking: &ExtendedThinking{ThinkingTokens: 10}}, false},
		{"only server tools", &Usage{ServerTools: &ServerToolUsage{WebSearchRequests: 1}}, false},
		{"only audio", &Usage{Audio: &AudioTokens{InputAudioTokens: 1}}, false},
		{"only prediction", &Usage{Prediction: &PredictionTokens{AcceptedTokens: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============== P0-3: Add() 聚合方法 ==============

func TestAdd_NilReceiver(t *testing.T) {
	var u *Usage
	u.Add(&Usage{Base: Base{InputTokens: 100}})
	if u != nil {
		t.Errorf("nil receiver Add 应保持 nil")
	}
}

func TestAdd_NilOther(t *testing.T) {
	u := &Usage{Base: Base{InputTokens: 100}}
	u.Add(nil)
	if u.InputTokens != 100 {
		t.Errorf("nil other 不应修改 receiver")
	}
}

func TestAdd_BaseFields(t *testing.T) {
	u := &Usage{Base: Base{InputTokens: 100, OutputTokens: 50, TotalTokens: 175}}
	other := &Usage{Base: Base{InputTokens: 20, OutputTokens: 10, TotalTokens: 30}}
	u.Add(other)
	if u.InputTokens != 120 || u.OutputTokens != 60 || u.TotalTokens != 205 {
		t.Errorf("Base 累加错误: %+v", u.Base)
	}
}

func TestAdd_Cache_BothNil(t *testing.T) {
	u := &Usage{}
	other := &Usage{}
	u.Add(other)
	if u.Cache != nil {
		t.Errorf("双方 Cache 都为 nil 时结果应保持 nil，实际 %+v", u.Cache)
	}
}

func TestAdd_Cache_ReceiverNilOtherNot(t *testing.T) {
	u := &Usage{}
	other := &Usage{Cache: &CacheMetrics{
		CreationInputTokens:    100,
		ReadInputTokens:        50,
		Ephemeral5mInputTokens: 30,
	}}
	u.Add(other)
	if u.Cache == nil {
		t.Fatalf("receiver Cache 为 nil 时应自动创建")
	}
	if u.Cache.CreationInputTokens != 100 || u.Cache.ReadInputTokens != 50 || u.Cache.Ephemeral5mInputTokens != 30 {
		t.Errorf("Cache 字段累加错误: %+v", u.Cache)
	}
}

func TestAdd_Cache_BothNotNil(t *testing.T) {
	u := &Usage{Cache: &CacheMetrics{CreationInputTokens: 100, Ephemeral1hInputTokens: 50}}
	other := &Usage{Cache: &CacheMetrics{CreationInputTokens: 30, Ephemeral5mInputTokens: 20, CachedTokens: 5}}
	u.Add(other)
	if u.Cache.CreationInputTokens != 130 {
		t.Errorf("CreationInputTokens 累加错误: %d", u.Cache.CreationInputTokens)
	}
	if u.Cache.Ephemeral1hInputTokens != 50 {
		t.Errorf("Ephemeral1hInputTokens 应保持 50，实际 %d", u.Cache.Ephemeral1hInputTokens)
	}
	if u.Cache.Ephemeral5mInputTokens != 20 {
		t.Errorf("Ephemeral5mInputTokens 累加错误: %d", u.Cache.Ephemeral5mInputTokens)
	}
	if u.Cache.CachedTokens != 5 {
		t.Errorf("CachedTokens 应从 other 复制 5，实际 %d", u.Cache.CachedTokens)
	}
}

func TestAdd_Reasoning(t *testing.T) {
	u := &Usage{Reasoning: &ReasoningTokens{ReasoningTokens: 10}}
	other := &Usage{Reasoning: &ReasoningTokens{ReasoningTokens: 5}}
	u.Add(other)
	if u.Reasoning.ReasoningTokens != 15 {
		t.Errorf("ReasoningTokens 累加错误: %d", u.Reasoning.ReasoningTokens)
	}
}

func TestAdd_Thinking(t *testing.T) {
	u := &Usage{}
	other := &Usage{Thinking: &ExtendedThinking{ThinkingTokens: 200}}
	u.Add(other)
	if u.Thinking == nil || u.Thinking.ThinkingTokens != 200 {
		t.Errorf("Thinking 字段累加错误: %+v", u.Thinking)
	}
}

func TestAdd_ServerTools(t *testing.T) {
	u := &Usage{ServerTools: &ServerToolUsage{WebSearchRequests: 2, WebFetchRequests: 1}}
	other := &Usage{ServerTools: &ServerToolUsage{WebSearchRequests: 3}}
	u.Add(other)
	if u.ServerTools.WebSearchRequests != 5 || u.ServerTools.WebFetchRequests != 1 {
		t.Errorf("ServerTools 累加错误: %+v", u.ServerTools)
	}
}

func TestAdd_Audio(t *testing.T) {
	u := &Usage{Audio: &AudioTokens{InputAudioTokens: 10, OutputAudioTokens: 5}}
	other := &Usage{Audio: &AudioTokens{InputAudioTokens: 20}}
	u.Add(other)
	if u.Audio.InputAudioTokens != 30 || u.Audio.OutputAudioTokens != 5 {
		t.Errorf("Audio 累加错误: %+v", u.Audio)
	}
}

func TestAdd_Prediction(t *testing.T) {
	u := &Usage{Prediction: &PredictionTokens{AcceptedTokens: 10, RejectedTokens: 5}}
	other := &Usage{Prediction: &PredictionTokens{RejectedTokens: 3}}
	u.Add(other)
	if u.Prediction.AcceptedTokens != 10 || u.Prediction.RejectedTokens != 8 {
		t.Errorf("Prediction 累加错误: %+v", u.Prediction)
	}
}

func TestAdd_ServiceTierOverride(t *testing.T) {
	// ServiceTier 非累加语义：other 覆盖 receiver
	u := &Usage{ServiceTier: ServiceTierStandard}
	other := &Usage{ServiceTier: ServiceTierPriority}
	u.Add(other)
	if u.ServiceTier != ServiceTierPriority {
		t.Errorf("ServiceTier 应被 other 覆盖，实际 %q", u.ServiceTier)
	}
}

func TestAdd_ServiceTierEmptyKeeps(t *testing.T) {
	// other.ServiceTier 为空时保持 receiver 不变
	u := &Usage{ServiceTier: ServiceTierStandard}
	other := &Usage{}
	u.Add(other)
	if u.ServiceTier != ServiceTierStandard {
		t.Errorf("other.ServiceTier 为空时不应覆盖 receiver")
	}
}

func TestAdd_ProtocolOverride(t *testing.T) {
	u := &Usage{Protocol: "openai"}
	other := &Usage{Protocol: "anthropic"}
	u.Add(other)
	if u.Protocol != "anthropic" {
		t.Errorf("Protocol 应被 other 覆盖，实际 %q", u.Protocol)
	}
}

// ============== 流式场景仿真：模拟 message_delta 多次 Add ==============

func TestAdd_StreamSimulation(t *testing.T) {
	// 模拟 Anthropic 流式响应：
	//   message_start: input=100, output=0
	//   message_delta (3 次): output += 30 每次 + thinking_tokens += 10 每次
	total := &Usage{}
	total.Add(&Usage{
		Base:     Base{InputTokens: 100},
		Protocol: "anthropic",
	})
	for i := 0; i < 3; i++ {
		total.Add(&Usage{
			Base:     Base{OutputTokens: 30},
			Thinking: &ExtendedThinking{ThinkingTokens: 10},
		})
	}
	if total.InputTokens != 100 {
		t.Errorf("InputTokens=%d want 100", total.InputTokens)
	}
	if total.OutputTokens != 90 {
		t.Errorf("OutputTokens=%d want 90", total.OutputTokens)
	}
	if total.Thinking == nil || total.Thinking.ThinkingTokens != 30 {
		t.Errorf("ThinkingTokens 应为 30，实际 %+v", total.Thinking)
	}
	if total.Protocol != "anthropic" {
		t.Errorf("Protocol 应保持 anthropic，实际 %q", total.Protocol)
	}
}

// ============== FromAnthropicMessageDelta 流式 delta usage ==============

func TestFromAnthropicMessageDelta_Base(t *testing.T) {
	src := &pkganthropic.MessageDeltaUsage{
		InputTokens:              100,
		OutputTokens:             50,
		CacheCreationInputTokens: 10,
		CacheReadInputTokens:     5,
	}
	got := FromAnthropicMessageDelta(src)
	// 2026-08-25 PR 改造：TotalTokens 自算全上下文 = Input + Output + Creation + Read
	wantTotal := int64(100 + 50 + 10 + 5)
	if got.InputTokens != 100 || got.OutputTokens != 50 || got.TotalTokens != wantTotal {
		t.Errorf("Base 字段错误: %+v, want TotalTokens=%d", got.Base, wantTotal)
	}
	// FreshInputTokens 等于 InputTokens（Anthropic 协议 input_tokens 本就不含 cache）
	if got.FreshInputTokens != 100 {
		t.Errorf("FreshInputTokens=%d want 100（Anthropic 协议两字段相等）", got.FreshInputTokens)
	}
	if got.Cache == nil ||
		got.Cache.CreationInputTokens != 10 ||
		got.Cache.ReadInputTokens != 5 {
		t.Errorf("Cache 字段错误: %+v", got.Cache)
	}
	if got.Protocol != "anthropic" {
		t.Errorf("Protocol 应为 anthropic")
	}
}

func TestFromAnthropicMessageDelta_Thinking(t *testing.T) {
	src := &pkganthropic.MessageDeltaUsage{
		InputTokens:  100,
		OutputTokens: 80,
		OutputTokensDetails: &pkganthropic.AnthropicOutputTokensDetails{
			ThinkingTokens: 30,
		},
	}
	got := FromAnthropicMessageDelta(src)
	if got.Thinking == nil || got.Thinking.ThinkingTokens != 30 {
		t.Errorf("Thinking 字段错误: %+v", got.Thinking)
	}
}

func TestFromAnthropicMessageDelta_ServerTools(t *testing.T) {
	src := &pkganthropic.MessageDeltaUsage{
		InputTokens:  100,
		OutputTokens: 50,
		ServerToolUse: &pkganthropic.AnthropicServerToolUsage{
			WebSearchRequests: 2,
			WebFetchRequests:  1,
		},
	}
	got := FromAnthropicMessageDelta(src)
	if got.ServerTools == nil ||
		got.ServerTools.WebSearchRequests != 2 ||
		got.ServerTools.WebFetchRequests != 1 {
		t.Errorf("ServerTools 字段错误: %+v", got.ServerTools)
	}
}

func TestFromAnthropicMessageDelta_Nil(t *testing.T) {
	if got := FromAnthropicMessageDelta(nil); got != nil {
		t.Errorf("FromAnthropicMessageDelta(nil) 应返回 nil")
	}
}

// 流式 message_start + message_delta 累计仿真（验证 FromAnthropic + Add + FromAnthropicMessageDelta）
//
// 注意：Anthropic 流协议 message_delta 的 InputTokens 通常为 0（仅在最后一个
// delta 重发完整 input 上下文），OutputTokens 是本轮增量。message_start 的 input
// 已确定后，message_delta 只追加 output / thinking / server_tool_use。
func TestStream_MessageStartPlusDelta(t *testing.T) {
	start := FromAnthropic(&pkganthropic.AnthropicUsage{
		InputTokens:  200,
		OutputTokens: 0,
		CacheCreation: &pkganthropic.AnthropicCacheCreation{
			Ephemeral5mInputTokens: 50,
		},
	})
	delta := FromAnthropicMessageDelta(&pkganthropic.MessageDeltaUsage{
		InputTokens:          0, // 流式 delta 不重发 input
		OutputTokens:         80,
		CacheReadInputTokens: 10,
		OutputTokensDetails:  &pkganthropic.AnthropicOutputTokensDetails{ThinkingTokens: 20},
	})
	start.Add(delta)
	// 2026-08-25 PR 改造：TotalTokens 自算全上下文 = Input + Output + Creation + Read
	// start: 200+0+0+0=200; delta: 0+80+0+10=90; Add 后 290
	if start.InputTokens != 200 || start.OutputTokens != 80 || start.TotalTokens != 290 {
		t.Errorf("Base 累加错误: %+v", start.Base)
	}
	// FreshInputTokens 等于 InputTokens（Anthropic 协议两字段相等；delta 的 InputTokens=0 不影响）
	if start.FreshInputTokens != 200 {
		t.Errorf("FreshInputTokens=%d want 200", start.FreshInputTokens)
	}
	if start.Cache == nil ||
		start.Cache.CreationInputTokens != 0 || // message_start 的 CacheCreation 是 5m 嵌套而非顶层字段
		start.Cache.Ephemeral5mInputTokens != 50 ||
		start.Cache.ReadInputTokens != 10 {
		t.Errorf("Cache 累加错误: %+v", start.Cache)
	}
	if start.Thinking == nil || start.Thinking.ThinkingTokens != 20 {
		t.Errorf("Thinking 应来自 delta: %+v", start.Thinking)
	}
}

// ============== FromAnthropic 嵌套字段 ==============

func TestFromAnthropic_CacheNested(t *testing.T) {
	src := &pkganthropic.AnthropicUsage{
		InputTokens:              1000,
		OutputTokens:             200,
		CacheCreationInputTokens: 100,
		CacheReadInputTokens:     50,
		CacheCreation: &pkganthropic.AnthropicCacheCreation{
			Ephemeral5mInputTokens: 80,
			Ephemeral1hInputTokens: 20,
		},
	}
	got := FromAnthropic(src)
	if got.Cache == nil {
		t.Fatalf("Cache 应非 nil")
	}
	if got.Cache.CreationInputTokens != 100 ||
		got.Cache.ReadInputTokens != 50 ||
		got.Cache.Ephemeral5mInputTokens != 80 ||
		got.Cache.Ephemeral1hInputTokens != 20 {
		t.Errorf("Cache 嵌套字段映射错误: %+v", got.Cache)
	}
}

func TestFromAnthropic_ThinkingAndServerTools(t *testing.T) {
	src := &pkganthropic.AnthropicUsage{
		InputTokens:  100,
		OutputTokens: 50,
		OutputTokensDetails: &pkganthropic.AnthropicOutputTokensDetails{
			ThinkingTokens: 30,
		},
		ServerToolUse: &pkganthropic.AnthropicServerToolUsage{
			WebSearchRequests: 2,
			WebFetchRequests:  1,
		},
		ServiceTier:  pkganthropic.AnthropicUsageServiceTierPriority,
		InferenceGeo: "us-west",
	}
	got := FromAnthropic(src)
	if got.Thinking == nil || got.Thinking.ThinkingTokens != 30 {
		t.Errorf("Thinking 映射错误: %+v", got.Thinking)
	}
	if got.ServerTools == nil ||
		got.ServerTools.WebSearchRequests != 2 ||
		got.ServerTools.WebFetchRequests != 1 {
		t.Errorf("ServerTools 映射错误: %+v", got.ServerTools)
	}
	if got.ServiceTier != ServiceTierPriority {
		t.Errorf("ServiceTier=%q want priority", got.ServiceTier)
	}
	if got.InferenceGeo != "us-west" {
		t.Errorf("InferenceGeo=%q want us-west", got.InferenceGeo)
	}
}

// ============== FromOpenAICompletion 嵌套字段 ==============

func TestFromOpenAICompletion_CachedAndReasoning(t *testing.T) {
	src := &pkgopenai.OpenAICompletionUsage{
		PromptTokens:     200,
		CompletionTokens: 100,
		TotalTokens:      300,
		PromptTokensDetails: &pkgopenai.OpenAICompletionUsagePromptTokensDetails{
			CachedTokens: 80,
			AudioTokens:  20,
		},
		CompletionTokensDetails: &pkgopenai.OpenAICompletionUsageCompletionTokensDetails{
			ReasoningTokens:          25,
			AudioTokens:              5,
			AcceptedPredictionTokens: 10,
			RejectedPredictionTokens: 3,
		},
	}
	got := FromOpenAICompletion(src)
	if got.Cache == nil || got.Cache.CachedTokens != 80 {
		t.Errorf("Cache.CachedTokens 映射错误: %+v", got.Cache)
	}
	if got.Reasoning == nil || got.Reasoning.ReasoningTokens != 25 {
		t.Errorf("Reasoning 映射错误: %+v", got.Reasoning)
	}
	if got.Prediction == nil ||
		got.Prediction.AcceptedTokens != 10 ||
		got.Prediction.RejectedTokens != 3 {
		t.Errorf("Prediction 映射错误: %+v", got.Prediction)
	}
	if got.Audio == nil ||
		got.Audio.InputAudioTokens != 20 ||
		got.Audio.OutputAudioTokens != 5 {
		t.Errorf("Audio 映射错误: %+v", got.Audio)
	}
}

// ============== Clone 深拷贝 ==============

func TestClone_DeepCopy(t *testing.T) {
	src := &Usage{
		Base:         Base{InputTokens: 100},
		Cache:        &CacheMetrics{CachedTokens: 50},
		Protocol:     "openai",
		ServiceTier:  ServiceTierPriority,
		InferenceGeo: "us",
	}
	dst := src.Clone()
	if dst == src {
		t.Fatalf("Clone 应返回新指针")
	}
	if dst.Cache == src.Cache {
		t.Errorf("Cache 应深拷贝，不应共享指针")
	}
	// 修改 dst 不影响 src
	dst.Cache.CachedTokens = 999
	if src.Cache.CachedTokens != 50 {
		t.Errorf("修改 dst 后 src 不应受影响")
	}
}

func TestClone_Nil(t *testing.T) {
	if got := (*Usage)(nil).Clone(); got != nil {
		t.Errorf("nil.Clone() 应返回 nil")
	}
}

// ============== CacheMetrics.Total ==============

func TestCacheMetrics_Total(t *testing.T) {
	c := &CacheMetrics{
		CreationInputTokens:    100,
		ReadInputTokens:        50,
		CachedTokens:           30,
		Ephemeral5mInputTokens: 20,
		Ephemeral1hInputTokens: 10,
	}
	if got := c.Total(); got != 210 {
		t.Errorf("Cache.Total()=%d want 210", got)
	}
}

// ============== ToAnthropic 反向投影 ==============

func TestToAnthropic_FlatCacheAndNestedCache(t *testing.T) {
	src := &Usage{
		Base: Base{
			InputTokens:  100,
			OutputTokens: 50,
			TotalTokens:  150,
		},
		Cache: &CacheMetrics{
			CreationInputTokens:    20,
			ReadInputTokens:        10,
			Ephemeral5mInputTokens: 15,
			Ephemeral1hInputTokens: 5,
		},
	}
	got := ToAnthropic(src)

	// 顶层扁平字段
	if got.InputTokens != 100 || got.OutputTokens != 50 {
		t.Errorf("flat tokens: got (%d,%d), want (100,50)", got.InputTokens, got.OutputTokens)
	}
	if got.CacheCreationInputTokens != 20 || got.CacheReadInputTokens != 10 {
		t.Errorf("flat cache: got (%d,%d), want (20,10)", got.CacheCreationInputTokens, got.CacheReadInputTokens)
	}
	// 4.6+ 嵌套分桶
	if got.CacheCreation == nil {
		t.Fatal("CacheCreation 嵌套未生成")
	}
	if got.CacheCreation.Ephemeral5mInputTokens != 15 || got.CacheCreation.Ephemeral1hInputTokens != 5 {
		t.Errorf("nested cache: got (5m=%d,1h=%d), want (15,5)",
			got.CacheCreation.Ephemeral5mInputTokens, got.CacheCreation.Ephemeral1hInputTokens)
	}
}

func TestToAnthropic_ThinkingAndServerTools(t *testing.T) {
	src := &Usage{
		Base:         Base{InputTokens: 10, OutputTokens: 5},
		Thinking:     &ExtendedThinking{ThinkingTokens: 3},
		ServerTools:  &ServerToolUsage{WebFetchRequests: 2, WebSearchRequests: 1},
		ServiceTier:  ServiceTierPriority,
		InferenceGeo: "us-east-1",
	}
	got := ToAnthropic(src)

	if got.OutputTokensDetails == nil || got.OutputTokensDetails.ThinkingTokens != 3 {
		t.Errorf("OutputTokensDetails: got %+v", got.OutputTokensDetails)
	}
	if got.ServerToolUse == nil ||
		got.ServerToolUse.WebFetchRequests != 2 ||
		got.ServerToolUse.WebSearchRequests != 1 {
		t.Errorf("ServerToolUse: got %+v", got.ServerToolUse)
	}
	if got.ServiceTier != pkganthropic.AnthropicUsageServiceTierPriority {
		t.Errorf("ServiceTier: got %q", got.ServiceTier)
	}
	if got.InferenceGeo != "us-east-1" {
		t.Errorf("InferenceGeo: got %q", got.InferenceGeo)
	}
}

func TestToAnthropic_Nil(t *testing.T) {
	// nil 入参 → 返回 non-nil 空 AnthropicUsage（handler 期望 non-nil）。
	got := ToAnthropic(nil)
	if got == nil {
		t.Fatal("nil in should return non-nil empty *AnthropicUsage")
	}
	if got.InputTokens != 0 || got.OutputTokens != 0 {
		t.Errorf("empty usage has non-zero tokens: %+v", got)
	}
	// 嵌套字段全部 nil（omitempty 会在 marshal 时跳过）
	if got.CacheCreation != nil || got.OutputTokensDetails != nil || got.ServerToolUse != nil {
		t.Errorf("nested fields should be nil for empty Usage")
	}
}

// ============== P0-2: HasInput / BillableInput 应识别 cache_creation ==============

// HasInput 必须把 cache_creation_input_tokens 视为有效输入。
//
// 触发场景：上游 MiniMax 返回非标 message_start，input_tokens=0 但 cache_creation_input_tokens>0
// （MiniMax 把首次缓存写入的 token 全部计入 cache_creation 而非 input_tokens）。
// HasInput 返回 false 会导致 handler 走估算兜底，丢掉上游真实用量数据。
func TestHasInput_CacheCreationOnly(t *testing.T) {
	u := &Usage{
		Base: Base{
			InputTokens:  0,
			OutputTokens: 23,
		},
		Cache: &CacheMetrics{
			CreationInputTokens: 34464,
		},
	}
	if !u.HasInput() {
		t.Errorf("HasInput 应识别 cache_creation_input_tokens=34464 作为有效输入，但返回 false")
	}
}

func TestHasInput_Nil(t *testing.T) {
	if (&Usage{}).HasInput() {
		t.Errorf("空 Usage 的 HasInput 应返回 false")
	}
	if (*Usage)(nil).HasInput() {
		t.Errorf("nil Usage 的 HasInput 应返回 false")
	}
}

func TestHasInput_InputTokensOnly(t *testing.T) {
	// 纯 InputTokens > 0 应识别（回归保护）
	u := &Usage{Base: Base{InputTokens: 100}}
	if !u.HasInput() {
		t.Errorf("InputTokens=100 时 HasInput 应返回 true")
	}
}

func TestHasInput_CacheReadOnly(t *testing.T) {
	// 回归保护：CacheReadInputTokens > 0 时 HasInput 仍应 true
	u := &Usage{Cache: &CacheMetrics{ReadInputTokens: 100}}
	if !u.HasInput() {
		t.Errorf("CacheReadInputTokens=100 时 HasInput 应返回 true")
	}
}

// BillableInput 必须把 cache_creation 计入应计费输入。
//
// 现行实现当 InputTokens=0 时直接 return 0，会把 cache_creation 的 token 全部漏扣。
// 修复后：max(0, InputTokens - CacheReadInputTokens) + CacheCreationInputTokens。
func TestBillableInput_CacheCreationOnly(t *testing.T) {
	u := &Usage{
		Base: Base{InputTokens: 0},
		Cache: &CacheMetrics{
			CreationInputTokens: 34464,
		},
	}
	if got := u.BillableInput(); got != 34464 {
		t.Errorf("BillableInput = %d, want 34464（cache_creation 必须计入计费）", got)
	}
}

// TestBillableInput_InputMinusCacheRead 验证 FreshInputTokens + 0 = BillableInput。
//
// 2026-08-25 PR 改造：BillableInput 公式 = FreshInputTokens + CacheCreationInputTokens，
// 不再依赖 InputTokens - ReadInputTokens。Anthropic 协议下两字段相等，OpenAI/Gemini
// 协议下 FreshInputTokens 已扣 cache，调用方传 FreshInputTokens 直接得到计费值。
func TestBillableInput_InputMinusCacheRead(t *testing.T) {
	// Anthropic 协议口径：InputTokens=100，ReadInputTokens=30
	// FreshInputTokens 由 FromAnthropic 自动设为 InputTokens=100（Anthropic 不含 cache）
	// 公式 BillableInput = FreshInputTokens + 0 = 100
	u := &Usage{
		Base:  Base{InputTokens: 100, FreshInputTokens: 100},
		Cache: &CacheMetrics{ReadInputTokens: 30},
	}
	if got := u.BillableInput(); got != 100 {
		t.Errorf("BillableInput = %d, want 100（FreshInputTokens 已是 Anthropic 口径，无需再减 cache_read）", got)
	}
}

// TestBillableInput_OpenAI_FreshFromMapping 验证 OpenAI 协议下 FreshInputTokens=200
// （来自 FromOpenAICompletion/ResponsesUsage 的 mapping）直接作为 BillableInput。
// 公式：BillableInput = FreshInputTokens + CacheCreationInputTokens = 200 + 0 = 200。
func TestBillableInput_OpenAI_FreshFromMapping(t *testing.T) {
	u := &Usage{
		Base: Base{
			InputTokens:      1000, // 协议原始值（含 cached）
			FreshInputTokens: 200,  // mapping 后去 cache 口径
		},
		Cache: &CacheMetrics{CachedTokens: 800},
	}
	if got := u.BillableInput(); got != 200 {
		t.Errorf("BillableInput = %d, want 200（FreshInputTokens 已扣 cached）", got)
	}
}

func TestBillableInput_InputPlusCacheCreation(t *testing.T) {
	// FreshInputTokens=100, CacheCreation=20 → Billable=120
	u := &Usage{
		Base: Base{InputTokens: 100, FreshInputTokens: 100},
		Cache: &CacheMetrics{
			ReadInputTokens:     30,
			CreationInputTokens: 20,
		},
	}
	if got := u.BillableInput(); got != 120 {
		t.Errorf("BillableInput = %d, want 120（FreshInputTokens + CacheCreation）", got)
	}
}

func TestBillableInput_Nil(t *testing.T) {
	if got := (*Usage)(nil).BillableInput(); got != 0 {
		t.Errorf("nil Usage 的 BillableInput 应返回 0, got %d", got)
	}
}

// ============== P0-12: 多模态字段（ImageCount / Audio*Chars / VideoDurationSec） ==============

// TestUsage_MultimodalFields 验证 4 个多模态字段能正确参与 JSON 序列化/反序列化。
// 场景：Seedance 视频任务返回 VideoDurationSec=5 + ImageCount=2 + AudioInputChars=1500，
// 流式 message_start 通过 Usage 透传到计费层，必须保证 round-trip 字段名/值不丢。
func TestUsage_MultimodalFields(t *testing.T) {
	in := Usage{
		Base:             Base{InputTokens: 100, OutputTokens: 50},
		ImageCount:       2,
		AudioInputChars:  1500,
		AudioOutputChars: 800,
		VideoDurationSec: 5,
	}
	b, err := json.Marshal(in)
	require.NoError(t, err)
	var out Usage
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, int64(2), out.ImageCount)
	assert.Equal(t, int64(1500), out.AudioInputChars)
	assert.Equal(t, int64(800), out.AudioOutputChars)
	assert.Equal(t, int64(5), out.VideoDurationSec)
	// Base 字段也必须 round-trip 不变
	assert.Equal(t, int64(100), out.InputTokens)
	assert.Equal(t, int64(50), out.OutputTokens)
}

// TestUsage_Add_Multimodal 验证 Add() 正确累加 4 个多模态字段。
// 场景：流式响应中 message_start 带 ImageCount=1（封面预览），message_delta 追加 ImageCount=1 + VideoDurationSec=5。
func TestUsage_Add_Multimodal(t *testing.T) {
	a := &Usage{ImageCount: 1, VideoDurationSec: 3}
	b := &Usage{ImageCount: 2, VideoDurationSec: 5}
	a.Add(b)
	assert.Equal(t, int64(3), a.ImageCount)
	assert.Equal(t, int64(8), a.VideoDurationSec)
}

// TestUsage_Add_Multimodal_AudioChars 验证 Audio 字符字段累加（覆盖 Input/Output 双向）。
func TestUsage_Add_Multimodal_AudioChars(t *testing.T) {
	a := &Usage{AudioInputChars: 100, AudioOutputChars: 50}
	b := &Usage{AudioInputChars: 200, AudioOutputChars: 80}
	a.Add(b)
	assert.Equal(t, int64(300), a.AudioInputChars)
	assert.Equal(t, int64(130), a.AudioOutputChars)
}

// TestUsage_IsZero_Multimodal 验证多模态字段非零时 IsZero() 返回 false。
// 场景：纯多模态请求（Seedance 视频生成）无 Base token，IsZero 须识别 multimodal 信号。
func TestUsage_IsZero_Multimodal(t *testing.T) {
	tests := []struct {
		name string
		u    *Usage
		want bool
	}{
		{"only image", &Usage{ImageCount: 1}, false},
		{"only audio input", &Usage{AudioInputChars: 100}, false},
		{"only audio output", &Usage{AudioOutputChars: 50}, false},
		{"only video", &Usage{VideoDurationSec: 5}, false},
		{"all multimodal zero", &Usage{ImageCount: 0, AudioInputChars: 0, AudioOutputChars: 0, VideoDurationSec: 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============== OpenAI 三层口径（InputTokens / FreshInputTokens / TotalTokens）==============

// TestFromOpenAICompletion_CachedSubset 验证 2026-08-25 PR 改造后的三层口径：
//   - InputTokens = u.PromptTokens（协议原始值，含 cached）
//   - FreshInputTokens = max(PromptTokens - CachedTokens, 0)（去 cache 后的计费输入）
//   - TotalTokens = u.TotalTokens（透传上游）
//   - Cache.CachedTokens = u.PromptTokensDetails.CachedTokens
func TestFromOpenAICompletion_CachedSubset(t *testing.T) {
	src := &pkgopenai.OpenAICompletionUsage{
		PromptTokens:     1000,
		CompletionTokens: 50,
		TotalTokens:      1050,
		PromptTokensDetails: &pkgopenai.OpenAICompletionUsagePromptTokensDetails{
			CachedTokens: 800,
		},
	}
	got := FromOpenAICompletion(src)
	// InputTokens 是协议原始值（含 cached）
	if got.InputTokens != 1000 {
		t.Errorf("InputTokens 应为上游原始值 1000（含 cached），实际 %d", got.InputTokens)
	}
	// FreshInputTokens 是去 cache 后的计费输入
	if got.FreshInputTokens != 200 {
		t.Errorf("FreshInputTokens 应为 1000 - 800 = 200，实际 %d", got.FreshInputTokens)
	}
	// TotalTokens 透传上游
	if got.TotalTokens != 1050 {
		t.Errorf("TotalTokens 应透传上游 1050，实际 %d", got.TotalTokens)
	}
	if got.Cache == nil || got.Cache.CachedTokens != 800 {
		t.Errorf("Cache.CachedTokens 期望 800，实际 %+v", got.Cache)
	}
	// 桶互斥：FreshInput + CacheRead = 全上下文 1000（保证 CalcCostUsage 不重复计费）
	if got.FreshInputTokens+got.Cache.CachedTokens != 1000 {
		t.Errorf("桶互斥破坏：FreshInput %d + Cache.Cached %d != 1000",
			got.FreshInputTokens, got.Cache.CachedTokens)
	}
}

// TestFromOpenAICompletion_NoCacheDetails 验证无 PromptTokensDetails 时正常映射。
func TestFromOpenAICompletion_NoCacheDetails(t *testing.T) {
	src := &pkgopenai.OpenAICompletionUsage{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
	}
	got := FromOpenAICompletion(src)
	if got.InputTokens != 100 {
		t.Errorf("无 cached 时 InputTokens 应等于 promptTokens，实际 %d", got.InputTokens)
	}
	if got.FreshInputTokens != 100 {
		t.Errorf("无 cached 时 FreshInputTokens 应等于 InputTokens，实际 %d", got.FreshInputTokens)
	}
	if got.Cache != nil {
		t.Errorf("无 cached 时 Cache 应为 nil，实际 %+v", got.Cache)
	}
}

// TestFromOpenAICompletion_CachedExceedsPrompt 验证 cached 异常大于 prompt 时 FreshInputTokens 不出现负数。
func TestFromOpenAICompletion_CachedExceedsPrompt(t *testing.T) {
	src := &pkgopenai.OpenAICompletionUsage{
		PromptTokens:     10,
		CompletionTokens: 5,
		TotalTokens:      15,
		PromptTokensDetails: &pkgopenai.OpenAICompletionUsagePromptTokensDetails{
			CachedTokens: 100, // 异常：cached > prompt
		},
	}
	got := FromOpenAICompletion(src)
	// InputTokens 是协议原始值（异常保留，供审计）
	if got.InputTokens != 10 {
		t.Errorf("InputTokens 应保留上游原始值 10，实际 %d", got.InputTokens)
	}
	// FreshInputTokens max(10 - 100, 0) = 0
	if got.FreshInputTokens != 0 {
		t.Errorf("cached > prompt 时 FreshInputTokens 应 max(0, ...) = 0，实际 %d", got.FreshInputTokens)
	}
	if got.Cache == nil || got.Cache.CachedTokens != 100 {
		t.Errorf("Cache.CachedTokens 应保留原值 100 供审计，实际 %+v", got.Cache)
	}
}

// TestFromOpenAIResponsesUsage_CachedSubset 验证 OpenAI Responses API 的三层口径。
func TestFromOpenAIResponsesUsage_CachedSubset(t *testing.T) {
	src := &pkgopenai.OpenAIResponseUsage{
		InputTokens:  1000,
		OutputTokens: 50,
		TotalTokens:  1050,
		InputTokensDetails: pkgopenai.OpenAIResponseUsageInputTokensDetails{
			CachedTokens: 800,
		},
	}
	got := FromOpenAIResponsesUsage(src)
	if got.InputTokens != 1000 {
		t.Errorf("InputTokens 应为上游原始值 1000，实际 %d", got.InputTokens)
	}
	if got.FreshInputTokens != 200 {
		t.Errorf("FreshInputTokens 应为 1000 - 800 = 200，实际 %d", got.FreshInputTokens)
	}
	if got.TotalTokens != 1050 {
		t.Errorf("TotalTokens 应透传上游 1050，实际 %d", got.TotalTokens)
	}
	if got.Cache == nil || got.Cache.CachedTokens != 800 {
		t.Errorf("Cache.CachedTokens 期望 800，实际 %+v", got.Cache)
	}
	if got.FreshInputTokens+got.Cache.CachedTokens != 1000 {
		t.Errorf("桶互斥破坏：FreshInput %d + Cache.Cached %d != 1000",
			got.FreshInputTokens, got.Cache.CachedTokens)
	}
}

// TestFromOpenAIResponsesUsage_NoCache 验证无 cached 时两字段相等且都为原值。
func TestFromOpenAIResponsesUsage_NoCache(t *testing.T) {
	src := &pkgopenai.OpenAIResponseUsage{
		InputTokens:  100,
		OutputTokens: 50,
		TotalTokens:  150,
	}
	got := FromOpenAIResponsesUsage(src)
	if got.InputTokens != 100 {
		t.Errorf("无 cached 时 InputTokens 应等于 input_tokens，实际 %d", got.InputTokens)
	}
	if got.FreshInputTokens != 100 {
		t.Errorf("无 cached 时 FreshInputTokens 应等于 InputTokens，实际 %d", got.FreshInputTokens)
	}
	if got.Cache != nil {
		t.Errorf("无 cached 时 Cache 应为 nil，实际 %+v", got.Cache)
	}
}

// TestBillableInput_OpenAI_CachedSubset 验证 OpenAI 协议下 BillableInput() 等于 FreshInputTokens。
//
// 修复后：mapping 层已把 cached 扣除到 FreshInputTokens，调用方传 FreshInputTokens
// 直接得到计费值；CachedTokens 单独走 CacheRead 桶；桶互斥。
// TestBillableInput_OpenAI_CachedSubset 验证 OpenAI 协议下 BillableInput() 等于 FreshInputTokens。
//
// 修复后：mapping 层已把 cached 扣除到 FreshInputTokens，调用方传 FreshInputTokens
// 直接得到计费值；CachedTokens 单独走 CacheRead 桶；桶互斥。
// 总成本 = FreshInputTokens × InputPrice + CachedTokens × CacheReadPrice + ...
func TestBillableInput_OpenAI_CachedSubset(t *testing.T) {
	u := &Usage{
		Base: Base{
			InputTokens:      1000, // 协议原始值（含 cached）
			FreshInputTokens: 200,  // mapping 后 = prompt - cached
		},
		Cache: &CacheMetrics{CachedTokens: 800},
	}
	got := u.BillableInput()
	// 修复后：BillableInput = FreshInputTokens + CacheCreationInputTokens = 200 + 0 = 200
	if got != 200 {
		t.Errorf("BillableInput 应等于 FreshInputTokens=200，实际 %d", got)
	}
}
