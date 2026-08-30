package openai

import "encoding/json"

// OpenAIResponseOutputItemUnionFlat OutputItem union 扁平化实现（26 变体）。
type OpenAIResponseOutputItemUnionFlat struct {
	// Type 区分变体，核心类型：
	// "message" / "function_call" / "web_search_call" / "file_search_call" /
	// "computer_call" / "reasoning" / "mcp_call" / "mcp_list_tools" /
	// "code_interpreter_call" / "image_generation_call" /
	// "tool_search_call" / "tool_search_output" /
	// "local_shell_call" / "local_shell_call_output" /
	// "shell_call" / "shell_call_output" /
	// "apply_patch_call" / "apply_patch_call_output" /
	// "mcp_approval_request" / "mcp_approval_response" /
	// "custom_tool_call" / "custom_tool_call_output" /
	// "compaction" / "additional_tools" / "function_web_search_call"
	Type string `json:"type,omitempty"`

	// === message ===
	ID      string                               `json:"id,omitempty"`
	Role    string                               `json:"role,omitempty"`
	Status  string                               `json:"status,omitempty"`
	Content []OpenAIResponseOutputMessageContent `json:"content,omitempty"`

	// === function_call ===
	CallID    string `json:"call_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`

	// === web_search_call / function_web_search_call ===
	Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`

	// === file_search_call ===
	Queries []string                              `json:"queries,omitempty"`
	Results []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`

	// === computer_call ===
	CallID2  string                                   `json:"call_id,omitempty"`
	Action2  *OpenAIResponseInputComputerAction       `json:"action,omitempty"`
	Outputs2 []OpenAIResponseOutputComputerScreenshot `json:"outputs,omitempty"`

	// === reasoning ===
	Summary          []OpenAIResponseOutputReasoningSummary `json:"summary,omitempty"`
	ReasoningContent []string                               `json:"content,omitempty"`

	// === mcp_call ===
	ServerLabel string                       `json:"server_label,omitempty"`
	Output      string                       `json:"output,omitempty"`
	Error       *OpenAIResponseInputMcpError `json:"error,omitempty"`

	// === mcp_list_tools ===
	Tools []OpenAIResponseOutputMcpTool `json:"tools,omitempty"`

	// === code_interpreter_call ===
	ContainerID string                                     `json:"container_id,omitempty"`
	Code        string                                     `json:"code,omitempty"`
	Outputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`

	// === image_generation_call ===
	Result        string `json:"result,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`

	// === tool_search_call ===
	Queries2 []string `json:"queries,omitempty"`

	// === tool_search_output ===
	Tools2 []OpenAIResponseOutputMcpTool `json:"tools,omitempty"`

	// === local_shell_call ===
	Action3 *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`

	// === local_shell_call_output ===
	Output2 string `json:"output,omitempty"`

	// === shell_call ===
	Action4 *OpenAIResponseInputShellAction `json:"action,omitempty"`

	// === shell_call_output ===
	MaxOutputLength *int64 `json:"max_output_length,omitempty"`
	Output3         string `json:"output,omitempty"`

	// === apply_patch_call ===
	Operation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`

	// === mcp_approval_request ===
	ApprovalRequestID string `json:"approval_request_id,omitempty"`

	// === mcp_approval_response ===
	Approve bool   `json:"approve,omitempty"`
	Reason  string `json:"reason,omitempty"`

	// === custom_tool_call ===
	Input string `json:"input,omitempty"`

	// === custom_tool_call_output ===
	Output4 string `json:"output,omitempty"`

	// === compaction ===
	EncryptedContent string `json:"encrypted_content,omitempty"`
	CreatedBy        string `json:"created_by,omitempty"`

	// === additional_tools ===
	Role2  string                         `json:"role,omitempty"`
	Tools3 []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
}

func (OpenAIResponseOutputItemUnionFlat) isResponseOutputItem() {}

func (u OpenAIResponseOutputItemUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "message":
		return json.Marshal(struct {
			Type    string                               `json:"type"`
			ID      string                               `json:"id"`
			Role    string                               `json:"role"`
			Status  string                               `json:"status"`
			Content []OpenAIResponseOutputMessageContent `json:"content"`
		}{
			Type:    "message",
			ID:      u.ID,
			Role:    u.Role,
			Status:  u.Status,
			Content: u.Content,
		})
	case "function_call":
		return json.Marshal(struct {
			Type      string `json:"type"`
			ID        string `json:"id"`
			CallID    string `json:"call_id"`
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
			Status    string `json:"status,omitempty"`
		}{
			Type:      "function_call",
			ID:        u.ID,
			CallID:    u.CallID,
			Name:      u.Name,
			Arguments: u.Arguments,
			Status:    u.Status,
		})
	case "web_search_call":
		return json.Marshal(struct {
			Type   string                              `json:"type"`
			ID     string                              `json:"id"`
			Status string                              `json:"status"`
			Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`
		}{
			Type:   "web_search_call",
			ID:     u.ID,
			Status: u.Status,
			Action: u.Action,
		})
	case "file_search_call":
		return json.Marshal(struct {
			Type    string                                `json:"type"`
			ID      string                                `json:"id"`
			Status  string                                `json:"status"`
			Queries []string                              `json:"queries,omitempty"`
			Results []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`
		}{
			Type:    "file_search_call",
			ID:      u.ID,
			Status:  u.Status,
			Queries: u.Queries,
			Results: u.Results,
		})
	case "computer_call":
		return json.Marshal(struct {
			Type    string                                   `json:"type"`
			ID      string                                   `json:"id"`
			CallID  string                                   `json:"call_id,omitempty"`
			Status  string                                   `json:"status,omitempty"`
			Action  *OpenAIResponseInputComputerAction       `json:"action,omitempty"`
			Outputs []OpenAIResponseOutputComputerScreenshot `json:"outputs,omitempty"`
		}{
			Type:    "computer_call",
			ID:      u.ID,
			CallID:  u.CallID2,
			Status:  u.Status,
			Action:  u.Action2,
			Outputs: u.Outputs2,
		})
	case "reasoning":
		return json.Marshal(struct {
			Type    string                                 `json:"type"`
			ID      string                                 `json:"id"`
			Status  string                                 `json:"status,omitempty"`
			Summary []OpenAIResponseOutputReasoningSummary `json:"summary,omitempty"`
			Content []string                               `json:"content,omitempty"`
		}{
			Type:    "reasoning",
			ID:      u.ID,
			Status:  u.Status,
			Summary: u.Summary,
			Content: u.ReasoningContent,
		})
	case "mcp_call":
		return json.Marshal(struct {
			Type        string                       `json:"type"`
			ID          string                       `json:"id"`
			Name        string                       `json:"name,omitempty"`
			ServerLabel string                       `json:"server_label,omitempty"`
			Arguments   string                       `json:"arguments,omitempty"`
			Output      string                       `json:"output,omitempty"`
			Error       *OpenAIResponseInputMcpError `json:"error,omitempty"`
			Status      string                       `json:"status,omitempty"`
		}{
			Type:        "mcp_call",
			ID:          u.ID,
			Name:        u.Name,
			ServerLabel: u.ServerLabel,
			Arguments:   u.Output,
			Output:      u.Output,
			Error:       u.Error,
			Status:      u.Status,
		})
	case "mcp_list_tools":
		return json.Marshal(struct {
			Type        string                        `json:"type"`
			ID          string                        `json:"id"`
			ServerLabel string                        `json:"server_label"`
			Tools       []OpenAIResponseOutputMcpTool `json:"tools"`
			Error       *OpenAIResponseInputMcpError  `json:"error,omitempty"`
		}{
			Type:        "mcp_list_tools",
			ID:          u.ID,
			ServerLabel: u.ServerLabel,
			Tools:       u.Tools,
			Error:       u.Error,
		})
	case "code_interpreter_call":
		return json.Marshal(struct {
			Type        string                                     `json:"type"`
			ID          string                                     `json:"id"`
			Status      string                                     `json:"status"`
			ContainerID string                                     `json:"container_id,omitempty"`
			Code        string                                     `json:"code,omitempty"`
			Outputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`
		}{
			Type:        "code_interpreter_call",
			ID:          u.ID,
			Status:      u.Status,
			ContainerID: u.ContainerID,
			Code:        u.Code,
			Outputs:     u.Outputs,
		})
	case "image_generation_call":
		return json.Marshal(struct {
			Type          string `json:"type"`
			ID            string `json:"id"`
			Status        string `json:"status"`
			Result        string `json:"result,omitempty"`
			RevisedPrompt string `json:"revised_prompt,omitempty"`
		}{
			Type:          "image_generation_call",
			ID:            u.ID,
			Status:        u.Status,
			Result:        u.Result,
			RevisedPrompt: u.RevisedPrompt,
		})
	case "tool_search_call":
		return json.Marshal(struct {
			Type    string   `json:"type"`
			ID      string   `json:"id"`
			Status  string   `json:"status"`
			Queries []string `json:"queries,omitempty"`
		}{
			Type:    "tool_search_call",
			ID:      u.ID,
			Status:  u.Status,
			Queries: u.Queries2,
		})
	case "tool_search_output":
		return json.Marshal(struct {
			Type  string                        `json:"type"`
			ID    string                        `json:"id"`
			Tools []OpenAIResponseOutputMcpTool `json:"tools,omitempty"`
		}{
			Type:  "tool_search_output",
			ID:    u.ID,
			Tools: u.Tools2,
		})
	case "local_shell_call":
		return json.Marshal(struct {
			Type   string                               `json:"type"`
			ID     string                               `json:"id"`
			CallID string                               `json:"call_id,omitempty"`
			Status string                               `json:"status,omitempty"`
			Action *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`
		}{
			Type:   "local_shell_call",
			ID:     u.ID,
			CallID: u.CallID,
			Status: u.Status,
			Action: u.Action3,
		})
	case "local_shell_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id"`
			Output string `json:"output,omitempty"`
		}{
			Type:   "local_shell_call_output",
			ID:     u.ID,
			Output: u.Output2,
		})
	case "shell_call":
		return json.Marshal(struct {
			Type   string                          `json:"type"`
			ID     string                          `json:"id"`
			CallID string                          `json:"call_id,omitempty"`
			Status string                          `json:"status,omitempty"`
			Action *OpenAIResponseInputShellAction `json:"action,omitempty"`
		}{
			Type:   "shell_call",
			ID:     u.ID,
			CallID: u.CallID,
			Status: u.Status,
			Action: u.Action4,
		})
	case "shell_call_output":
		return json.Marshal(struct {
			Type            string `json:"type"`
			ID              string `json:"id"`
			MaxOutputLength *int64 `json:"max_output_length,omitempty"`
			Output          string `json:"output,omitempty"`
		}{
			Type:            "shell_call_output",
			ID:              u.ID,
			MaxOutputLength: u.MaxOutputLength,
			Output:          u.Output3,
		})
	case "apply_patch_call":
		return json.Marshal(struct {
			Type      string                                  `json:"type"`
			ID        string                                  `json:"id"`
			CallID    string                                  `json:"call_id,omitempty"`
			Status    string                                  `json:"status,omitempty"`
			Operation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`
		}{
			Type:      "apply_patch_call",
			ID:        u.ID,
			CallID:    u.CallID,
			Status:    u.Status,
			Operation: u.Operation,
		})
	case "apply_patch_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id"`
			CallID string `json:"call_id,omitempty"`
			Status string `json:"status,omitempty"`
		}{
			Type:   "apply_patch_call_output",
			ID:     u.ID,
			CallID: u.CallID,
			Status: u.Status,
		})
	case "mcp_approval_request":
		return json.Marshal(struct {
			Type              string `json:"type"`
			ID                string `json:"id"`
			ServerLabel       string `json:"server_label,omitempty"`
			Name              string `json:"name,omitempty"`
			Arguments         string `json:"arguments,omitempty"`
			ApprovalRequestID string `json:"approval_request_id,omitempty"`
		}{
			Type:              "mcp_approval_request",
			ID:                u.ID,
			ServerLabel:       u.ServerLabel,
			Name:              u.Name,
			Arguments:         u.Output,
			ApprovalRequestID: u.ApprovalRequestID,
		})
	case "mcp_approval_response":
		return json.Marshal(struct {
			Type              string `json:"type"`
			ID                string `json:"id,omitempty"`
			ApprovalRequestID string `json:"approval_request_id"`
			Approve           bool   `json:"approve"`
			Reason            string `json:"reason,omitempty"`
		}{
			Type:              "mcp_approval_response",
			ID:                u.ID,
			ApprovalRequestID: u.ApprovalRequestID,
			Approve:           u.Approve,
			Reason:            u.Reason,
		})
	case "custom_tool_call":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id"`
			CallID string `json:"call_id,omitempty"`
			Name   string `json:"name"`
			Input  string `json:"input"`
			Status string `json:"status,omitempty"`
		}{
			Type:   "custom_tool_call",
			ID:     u.ID,
			CallID: u.CallID,
			Name:   u.Name,
			Input:  u.Input,
			Status: u.Status,
		})
	case "custom_tool_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id"`
			CallID string `json:"call_id,omitempty"`
			Output string `json:"output"`
		}{
			Type:   "custom_tool_call_output",
			ID:     u.ID,
			CallID: u.CallID,
			Output: u.Output4,
		})
	case "compaction":
		return json.Marshal(struct {
			Type             string `json:"type"`
			ID               string `json:"id"`
			EncryptedContent string `json:"encrypted_content,omitempty"`
			CreatedBy        string `json:"created_by,omitempty"`
		}{
			Type:             "compaction",
			ID:               u.ID,
			EncryptedContent: u.EncryptedContent,
			CreatedBy:        u.CreatedBy,
		})
	case "additional_tools":
		return json.Marshal(struct {
			Type  string                         `json:"type"`
			ID    string                         `json:"id"`
			Role  string                         `json:"role,omitempty"`
			Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
		}{
			Type:  "additional_tools",
			ID:    u.ID,
			Role:  u.Role2,
			Tools: u.Tools3,
		})
	}
	return nil, nil
}

func (u *OpenAIResponseOutputItemUnionFlat) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var tmp struct {
		Type string `json:"type"`
		// message
		ID      string                               `json:"id"`
		Role    string                               `json:"role"`
		Status  string                               `json:"status"`
		Content []OpenAIResponseOutputMessageContent `json:"content"`
		// function_call
		CallID    string `json:"call_id"`
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
		// web_search_call
		Action *OpenAIResponseInputWebSearchAction `json:"action"`
		// file_search_call
		Queries []string                              `json:"queries"`
		Results []OpenAIResponseInputFileSearchResult `json:"results"`
		// computer_call
		ComputerCallID  string                                   `json:"call_id"`
		ComputerAction  *OpenAIResponseInputComputerAction       `json:"action"`
		ComputerOutputs []OpenAIResponseOutputComputerScreenshot `json:"outputs"`
		// reasoning
		Summary          []OpenAIResponseOutputReasoningSummary `json:"summary"`
		ReasoningContent []string                               `json:"content"`
		// mcp_call
		ServerLabel string                       `json:"server_label"`
		McpOutput   string                       `json:"output"`
		Error       *OpenAIResponseInputMcpError `json:"error"`
		// mcp_list_tools
		Tools []OpenAIResponseOutputMcpTool `json:"tools"`
		// code_interpreter_call
		ContainerID string                                     `json:"container_id"`
		Code        string                                     `json:"code"`
		CodeOutputs []OpenAIResponseInputCodeInterpreterOutput `json:"outputs"`
		// image_generation_call
		Result        string `json:"result"`
		RevisedPrompt string `json:"revised_prompt"`
		// tool_search_call
		ToolSearchQueries []string `json:"queries"`
		// tool_search_output
		ToolSearchTools []OpenAIResponseOutputMcpTool `json:"tools"`
		// local_shell_call
		LocalShellAction *OpenAIResponseInputLocalShellAction `json:"action"`
		// local_shell_call_output
		LocalShellOutput string `json:"output"`
		// shell_call
		ShellAction *OpenAIResponseInputShellAction `json:"action"`
		// shell_call_output
		MaxOutputLength *int64 `json:"max_output_length"`
		ShellOutput     string `json:"output"`
		// apply_patch_call
		Operation *OpenAIResponseInputApplyPatchOperation `json:"operation"`
		// mcp_approval_request
		ApprovalRequestID string `json:"approval_request_id"`
		// mcp_approval_response
		Approve bool   `json:"approve"`
		Reason  string `json:"reason"`
		// custom_tool_call
		Input string `json:"input"`
		// custom_tool_call_output
		CustomToolOutput string `json:"output"`
		// compaction
		EncryptedContent string `json:"encrypted_content"`
		CreatedBy        string `json:"created_by"`
		// additional_tools
		Role2  string                         `json:"role"`
		Tools3 []OpenAIResponsesToolUnionFlat `json:"tools"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "message":
		u.ID = tmp.ID
		u.Role = tmp.Role
		u.Status = tmp.Status
		u.Content = tmp.Content
	case "function_call":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Name = tmp.Name
		u.Arguments = tmp.Arguments
		u.Status = tmp.Status
	case "web_search_call":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.Action = tmp.Action
	case "function_web_search_call":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.Action = tmp.Action
	case "file_search_call":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.Queries = tmp.Queries
		u.Results = tmp.Results
	case "computer_call":
		u.ID = tmp.ID
		u.CallID2 = tmp.ComputerCallID
		u.Status = tmp.Status
		u.Action2 = tmp.ComputerAction
		u.Outputs2 = tmp.ComputerOutputs
	case "reasoning":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.Summary = tmp.Summary
		u.ReasoningContent = tmp.ReasoningContent
	case "mcp_call":
		u.ID = tmp.ID
		u.Name = tmp.Name
		u.ServerLabel = tmp.ServerLabel
		u.Output = tmp.McpOutput
		u.Error = tmp.Error
		u.Status = tmp.Status
	case "mcp_list_tools":
		u.ID = tmp.ID
		u.ServerLabel = tmp.ServerLabel
		u.Tools = tmp.Tools
		u.Error = tmp.Error
	case "code_interpreter_call":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.ContainerID = tmp.ContainerID
		u.Code = tmp.Code
		u.Outputs = tmp.CodeOutputs
	case "image_generation_call":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.Result = tmp.Result
		u.RevisedPrompt = tmp.RevisedPrompt
	case "tool_search_call":
		u.ID = tmp.ID
		u.Status = tmp.Status
		u.Queries2 = tmp.ToolSearchQueries
	case "tool_search_output":
		u.ID = tmp.ID
		u.Tools2 = tmp.ToolSearchTools
	case "local_shell_call":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Status = tmp.Status
		u.Action3 = tmp.LocalShellAction
	case "local_shell_call_output":
		u.ID = tmp.ID
		u.Output2 = tmp.LocalShellOutput
	case "shell_call":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Status = tmp.Status
		u.Action4 = tmp.ShellAction
	case "shell_call_output":
		u.ID = tmp.ID
		u.MaxOutputLength = tmp.MaxOutputLength
		u.Output3 = tmp.ShellOutput
	case "apply_patch_call":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Status = tmp.Status
		u.Operation = tmp.Operation
	case "apply_patch_call_output":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Status = tmp.Status
	case "mcp_approval_request":
		u.ID = tmp.ID
		u.ServerLabel = tmp.ServerLabel
		u.Name = tmp.Name
		u.Output = tmp.McpOutput
		u.ApprovalRequestID = tmp.ApprovalRequestID
	case "mcp_approval_response":
		u.ID = tmp.ID
		u.ApprovalRequestID = tmp.ApprovalRequestID
		u.Approve = tmp.Approve
		u.Reason = tmp.Reason
	case "custom_tool_call":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Name = tmp.Name
		u.Input = tmp.Input
		u.Status = tmp.Status
	case "custom_tool_call_output":
		u.ID = tmp.ID
		u.CallID = tmp.CallID
		u.Output4 = tmp.CustomToolOutput
	case "compaction":
		u.ID = tmp.ID
		u.EncryptedContent = tmp.EncryptedContent
		u.CreatedBy = tmp.CreatedBy
	case "additional_tools":
		u.ID = tmp.ID
		u.Role2 = tmp.Role2
		u.Tools3 = tmp.Tools3
	}
	return nil
}

// OpenAIResponseOutputItemUnion OutputItem union 接口。
//
// 23 种实现（4.6+）：ResponseOutputMessage / ResponseOutputFunctionCall /
// ResponseOutputWebSearchCall / ResponseOutputFileSearchCall / ResponseOutputComputerCall /
// ResponseOutputReasoning / ResponseOutputMcpCall / ResponseOutputMcpListTools /
// ResponseOutputCodeInterpreterCall / ResponseOutputImageGenerationCall /
// ResponseOutputToolSearchCall / ResponseOutputToolSearchOutput /
// ResponseOutputLocalShellCall(+Output) / ResponseOutputShellCall(+Output) /
// ResponseOutputApplyPatchToolCall(+Output) / ResponseOutputMcpApprovalRequest /
// ResponseOutputMcpApprovalResponse / ResponseOutputCustomToolCall(+Output) /
// ResponseOutputCompactionItem / ResponseOutputFunctionShellToolCall(+Output)。
type OpenAIResponseOutputItemUnion interface {
	isResponseOutputItem()
}

// OpenAIResponseOutputMessage 响应输出消息。
type OpenAIResponseOutputMessage struct {
	Type    string                               `json:"type"` // "message"
	ID      string                               `json:"id"`
	Role    string                               `json:"role"` // "assistant"
	Status  string                               `json:"status"`
	Content []OpenAIResponseOutputMessageContent `json:"content"`
}

func (OpenAIResponseOutputMessage) isResponseOutputItem() {}

// OpenAIResponseOutputMessageContent 响应消息内容项。
type OpenAIResponseOutputMessageContent struct {
	Type        string                                      `json:"type"` // "output_text" / "refusal"
	Text        string                                      `json:"text,omitempty"`
	Annotations []OpenAIResponseOutputMessageAnnotation     `json:"annotations,omitempty"`
	Refusal     string                                      `json:"refusal,omitempty"`
	Logprobs    []OpenAIResponseOutputMessageContentLogprob `json:"logprobs,omitempty"`
}

// OpenAIResponseOutputMessageContentLogprob 输出文本 logprob 项。
type OpenAIResponseOutputMessageContentLogprob struct {
	Token       string                                                `json:"token"`
	Bytes       []int64                                               `json:"bytes,omitempty"`
	Logprob     float64                                               `json:"logprob"`
	TopLogprobs []OpenAIResponseOutputMessageContentLogprobTopLogprob `json:"top_logprobs,omitempty"`
}

// OpenAIResponseOutputMessageContentLogprobTopLogprob top logprob 项。
type OpenAIResponseOutputMessageContentLogprobTopLogprob struct {
	Token   string  `json:"token"`
	Bytes   []int64 `json:"bytes,omitempty"`
	Logprob float64 `json:"logprob"`
}

// OpenAIResponseOutputMessageAnnotation 响应消息注释。
type OpenAIResponseOutputMessageAnnotation struct {
	Type         string                                   `json:"type"` // "url_citation" / "file_citation" / "container_file_citation"
	URLCitation  *OpenAIResponseOutputMessageURLCitation  `json:"url_citation,omitempty"`
	FileCitation *OpenAIResponseOutputMessageFileCitation `json:"file_citation,omitempty"`
}

// OpenAIResponseOutputMessageURLCitation URL 引用注释。
type OpenAIResponseOutputMessageURLCitation struct {
	URL        string `json:"url"`
	Title      string `json:"title"`
	StartIndex int64  `json:"start_index"`
	EndIndex   int64  `json:"end_index"`
}

// OpenAIResponseOutputMessageFileCitation 文件引用注释。
type OpenAIResponseOutputMessageFileCitation struct {
	FileID   string `json:"file_id"`
	Filename string `json:"filename,omitempty"`
	Index    int64  `json:"index,omitempty"`
}

// OpenAIResponseOutputFunctionCall 响应输出函数调用。
type OpenAIResponseOutputFunctionCall struct {
	Type      string `json:"type"` // "function_call"
	ID        string `json:"id"`
	CallID    string `json:"call_id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
	Status    string `json:"status,omitempty"`
}

func (OpenAIResponseOutputFunctionCall) isResponseOutputItem() {}

// OpenAIResponseOutputWebSearchCall 响应输出 Web 搜索调用。
type OpenAIResponseOutputWebSearchCall struct {
	Type   string                              `json:"type"` // "web_search_call"
	ID     string                              `json:"id"`
	Status string                              `json:"status"`
	Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`
}

func (OpenAIResponseOutputWebSearchCall) isResponseOutputItem() {}

// OpenAIResponseOutputFileSearchCall 响应输出文件搜索调用。
type OpenAIResponseOutputFileSearchCall struct {
	Type    string                                `json:"type"` // "file_search_call"
	ID      string                                `json:"id"`
	Status  string                                `json:"status"`
	Queries []string                              `json:"queries,omitempty"`
	Results []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`
}

func (OpenAIResponseOutputFileSearchCall) isResponseOutputItem() {}

// OpenAIResponseOutputComputerCall 响应输出电脑使用调用。
type OpenAIResponseOutputComputerCall struct {
	Type    string                                   `json:"type"` // "computer_call"
	ID      string                                   `json:"id"`
	CallID  string                                   `json:"call_id,omitempty"`
	Status  string                                   `json:"status,omitempty"`
	Action  *OpenAIResponseInputComputerAction       `json:"action,omitempty"`
	Outputs []OpenAIResponseOutputComputerScreenshot `json:"outputs,omitempty"`
}

func (OpenAIResponseOutputComputerCall) isResponseOutputItem() {}

// OpenAIResponseOutputComputerScreenshot 电脑使用截图。
type OpenAIResponseOutputComputerScreenshot struct {
	Type     string `json:"type"` // "computer_screenshot"
	FileID   string `json:"file_id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// OpenAIResponseOutputReasoning 响应输出推理。
type OpenAIResponseOutputReasoning struct {
	Type    string                                 `json:"type"` // "reasoning"
	ID      string                                 `json:"id"`
	Status  string                                 `json:"status,omitempty"`
	Summary []OpenAIResponseOutputReasoningSummary `json:"summary,omitempty"`
	Content []string                               `json:"content,omitempty"` // 推理文本
}

func (OpenAIResponseOutputReasoning) isResponseOutputItem() {}

// OpenAIResponseOutputReasoningSummary 推理摘要项。
type OpenAIResponseOutputReasoningSummary struct {
	Type string `json:"type"` // "summary_text"
	Text string `json:"text"`
}

// OpenAIResponseOutputMcpCall 响应输出 MCP 调用。
type OpenAIResponseOutputMcpCall struct {
	Type        string                       `json:"type"` // "mcp_call"
	ID          string                       `json:"id"`
	Name        string                       `json:"name,omitempty"`
	ServerLabel string                       `json:"server_label,omitempty"`
	Arguments   string                       `json:"arguments,omitempty"`
	Output      string                       `json:"output,omitempty"`
	Error       *OpenAIResponseInputMcpError `json:"error,omitempty"`
	Status      string                       `json:"status,omitempty"`
}

func (OpenAIResponseOutputMcpCall) isResponseOutputItem() {}

// OpenAIResponseOutputMcpListTools 响应输出 MCP 工具列表。
type OpenAIResponseOutputMcpListTools struct {
	Type        string                        `json:"type"` // "mcp_list_tools"
	ID          string                        `json:"id"`
	ServerLabel string                        `json:"server_label"`
	Tools       []OpenAIResponseOutputMcpTool `json:"tools"`
	Error       *OpenAIResponseInputMcpError  `json:"error,omitempty"`
}

func (OpenAIResponseOutputMcpListTools) isResponseOutputItem() {}

// OpenAIResponseOutputMcpTool MCP 工具定义。
type OpenAIResponseOutputMcpTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"input_schema,omitempty"`
	Annotations json.RawMessage `json:"annotations,omitempty"`
}

// OpenAIResponseOutputCodeInterpreterCall 响应输出代码解释器调用。
type OpenAIResponseOutputCodeInterpreterCall struct {
	Type        string                                     `json:"type"` // "code_interpreter_call"
	ID          string                                     `json:"id"`
	Status      string                                     `json:"status"`
	ContainerID string                                     `json:"container_id,omitempty"`
	Code        string                                     `json:"code,omitempty"`
	Outputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`
}

func (OpenAIResponseOutputCodeInterpreterCall) isResponseOutputItem() {}

// OpenAIResponseOutputImageGenerationCall 响应输出图像生成调用。
type OpenAIResponseOutputImageGenerationCall struct {
	Type          string `json:"type"` // "image_generation_call"
	ID            string `json:"id"`
	Status        string `json:"status"`
	Result        string `json:"result,omitempty"` // base64
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

func (OpenAIResponseOutputImageGenerationCall) isResponseOutputItem() {}

// === 4.6+ 扩展变体（13 种）===

// OpenAIResponseOutputToolSearchCall 响应输出工具搜索调用。
type OpenAIResponseOutputToolSearchCall struct {
	Type    string   `json:"type"` // "tool_search_call"
	ID      string   `json:"id"`
	Status  string   `json:"status"`
	Queries []string `json:"queries,omitempty"`
}

func (OpenAIResponseOutputToolSearchCall) isResponseOutputItem() {}

// OpenAIResponseOutputToolSearchOutput 响应输出工具搜索结果。
type OpenAIResponseOutputToolSearchOutput struct {
	Type string `json:"type"` // "tool_search_output"
	ID   string `json:"id"`
	// Tools 工具列表。
	Tools []OpenAIResponseOutputMcpTool `json:"tools,omitempty"`
}

func (OpenAIResponseOutputToolSearchOutput) isResponseOutputItem() {}

// OpenAIResponseOutputLocalShellCall 响应输出本地 shell 调用。
type OpenAIResponseOutputLocalShellCall struct {
	Type   string                               `json:"type"` // "local_shell_call"
	ID     string                               `json:"id"`
	CallID string                               `json:"call_id,omitempty"`
	Status string                               `json:"status,omitempty"`
	Action *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`
}

func (OpenAIResponseOutputLocalShellCall) isResponseOutputItem() {}

// OpenAIResponseOutputLocalShellCallOutput 响应输出本地 shell 输出。
type OpenAIResponseOutputLocalShellCallOutput struct {
	Type   string `json:"type"` // "local_shell_call_output"
	ID     string `json:"id"`
	Output string `json:"output,omitempty"`
}

func (OpenAIResponseOutputLocalShellCallOutput) isResponseOutputItem() {}

// OpenAIResponseOutputShellCall 响应输出远程 shell 调用。
type OpenAIResponseOutputShellCall struct {
	Type   string                          `json:"type"` // "shell_call"
	ID     string                          `json:"id"`
	CallID string                          `json:"call_id,omitempty"`
	Status string                          `json:"status,omitempty"`
	Action *OpenAIResponseInputShellAction `json:"action,omitempty"`
}

func (OpenAIResponseOutputShellCall) isResponseOutputItem() {}

// OpenAIResponseOutputShellCallOutput 响应输出远程 shell 输出。
type OpenAIResponseOutputShellCallOutput struct {
	Type            string `json:"type"` // "shell_call_output"
	ID              string `json:"id"`
	MaxOutputLength *int64 `json:"max_output_length,omitempty"`
	Output          string `json:"output,omitempty"`
}

func (OpenAIResponseOutputShellCallOutput) isResponseOutputItem() {}

// OpenAIResponseOutputApplyPatchToolCall 响应输出应用补丁调用。
//
// 注：与 Responses Tool 侧的 OpenAIResponsesApplyPatchToolParam（type=apply_patch）
// 是不同概念；本结构是 OutputItem 输出项（type=apply_patch_call）。
type OpenAIResponseOutputApplyPatchToolCall struct {
	Type      string                                  `json:"type"` // "apply_patch_call"
	ID        string                                  `json:"id"`
	CallID    string                                  `json:"call_id,omitempty"`
	Status    string                                  `json:"status,omitempty"`
	Operation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`
}

func (OpenAIResponseOutputApplyPatchToolCall) isResponseOutputItem() {}

// OpenAIResponseOutputApplyPatchToolCallOutput 响应输出应用补丁结果。
type OpenAIResponseOutputApplyPatchToolCallOutput struct {
	Type   string `json:"type"` // "apply_patch_call_output"
	ID     string `json:"id"`
	CallID string `json:"call_id,omitempty"`
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseOutputApplyPatchToolCallOutput) isResponseOutputItem() {}

// OpenAIResponseOutputMcpApprovalRequest 响应输出 MCP 工具审批请求。
type OpenAIResponseOutputMcpApprovalRequest struct {
	Type              string `json:"type"` // "mcp_approval_request"
	ID                string `json:"id"`
	ServerLabel       string `json:"server_label,omitempty"`
	Name              string `json:"name,omitempty"`
	Arguments         string `json:"arguments,omitempty"`
	ApprovalRequestID string `json:"approval_request_id,omitempty"`
}

func (OpenAIResponseOutputMcpApprovalRequest) isResponseOutputItem() {}

// OpenAIResponseOutputMcpApprovalResponse 响应输出 MCP 工具审批响应。
type OpenAIResponseOutputMcpApprovalResponse struct {
	Type              string `json:"type"` // "mcp_approval_response"
	ID                string `json:"id,omitempty"`
	ApprovalRequestID string `json:"approval_request_id"`
	Approve           bool   `json:"approve"`
	Reason            string `json:"reason,omitempty"`
}

func (OpenAIResponseOutputMcpApprovalResponse) isResponseOutputItem() {}

// OpenAIResponseOutputCustomToolCall 响应输出自定义工具调用。
type OpenAIResponseOutputCustomToolCall struct {
	Type   string `json:"type"` // "custom_tool_call"
	ID     string `json:"id"`
	CallID string `json:"call_id,omitempty"`
	Name   string `json:"name"`
	Input  string `json:"input"`
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseOutputCustomToolCall) isResponseOutputItem() {}

// OpenAIResponseOutputCustomToolCallOutput 响应输出自定义工具结果。
type OpenAIResponseOutputCustomToolCallOutput struct {
	Type   string `json:"type"` // "custom_tool_call_output"
	ID     string `json:"id"`
	CallID string `json:"call_id,omitempty"`
	Output string `json:"output"`
}

func (OpenAIResponseOutputCustomToolCallOutput) isResponseOutputItem() {}

// OpenAIResponseOutputCompactionItem 响应输出上下文压缩项。
type OpenAIResponseOutputCompactionItem struct {
	Type             string `json:"type"` // "compaction"
	ID               string `json:"id"`
	EncryptedContent string `json:"encrypted_content,omitempty"`
	CreatedBy        string `json:"created_by,omitempty"`
}

func (OpenAIResponseOutputCompactionItem) isResponseOutputItem() {}

// OpenAIResponseOutputFunctionShellToolCall 响应输出 shell 工具调用（FunctionShellTool）。
type OpenAIResponseOutputFunctionShellToolCall struct {
	Type   string                          `json:"type"` // "shell_call"
	ID     string                          `json:"id"`
	CallID string                          `json:"call_id,omitempty"`
	Status string                          `json:"status,omitempty"`
	Action *OpenAIResponseInputShellAction `json:"action,omitempty"`
}

func (OpenAIResponseOutputFunctionShellToolCall) isResponseOutputItem() {}

// OpenAIResponseOutputFunctionShellToolCallOutput 响应输出 shell 工具结果。
type OpenAIResponseOutputFunctionShellToolCallOutput struct {
	Type            string `json:"type"` // "shell_call_output"
	ID              string `json:"id"`
	MaxOutputLength *int64 `json:"max_output_length,omitempty"`
	Output          string `json:"output,omitempty"`
}

func (OpenAIResponseOutputFunctionShellToolCallOutput) isResponseOutputItem() {}

// OpenAIResponseOutputAdditionalTools 响应输出附加工具集（P3-2）。
type OpenAIResponseOutputAdditionalTools struct {
	// Type 固定 "additional_tools"。
	Type string `json:"type"` // "additional_tools"
	// ID 唯一 ID。
	ID string `json:"id"`
	// Role 角色。
	Role string `json:"role,omitempty"`
	// Tools 工具定义列表。
	Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
}

func (OpenAIResponseOutputAdditionalTools) isResponseOutputItem() {}

// OpenAIResponseOutputFunctionWebSearch 响应输出 web_search 工具调用（P3 补充）。
type OpenAIResponseOutputFunctionWebSearch struct {
	// Type 固定 "web_search_call"。
	Type   string                              `json:"type"` // "web_search_call"
	ID     string                              `json:"id"`
	Status string                              `json:"status,omitempty"`
	Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`
}

func (OpenAIResponseOutputFunctionWebSearch) isResponseOutputItem() {}
