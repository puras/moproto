package openai

import "encoding/json"

// ============================================================================
// Flattened Union Structs
// ============================================================================

// OpenAIResponseInputMessageContentUnionFlat Input message content union 扁平化实现。
type OpenAIResponseInputMessageContentUnionFlat struct {
	Type  string                                           `json:"type,omitempty"` // "string" / "content_parts"
	Text  string                                           `json:"text,omitempty"`
	Parts []OpenAIResponseInputMessageContentPartUnionFlat `json:"content_parts,omitempty"`
}

func (OpenAIResponseInputMessageContentUnionFlat) isResponseInputMessageContent() {}

func (u OpenAIResponseInputMessageContentUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "string":
		if u.Text == "" {
			return nil, nil
		}
		// Return as raw JSON string value (with quotes)
		return []byte("\"" + u.Text + "\""), nil
	case "content_parts":
		if u.Parts == nil {
			return nil, nil
		}
		return json.Marshal(u.Parts)
	}
	return nil, nil
}

func (u *OpenAIResponseInputMessageContentUnionFlat) UnmarshalJSON(data []byte) error {
	// Check first char: '"' means raw string, '{' or '[' means object
	if len(data) > 0 && data[0] == '"' {
		// Raw string: "hi"
		u.Type = "string"
		return json.Unmarshal(data, &u.Text)
	}
	// Object with type field
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "string":
		// Should not reach here if first char check worked, but handle anyway
		if err := json.Unmarshal(data, &u.Text); err != nil {
			return err
		}
	case "content_parts":
		var parts []OpenAIResponseInputMessageContentPartUnionFlat
		if err := json.Unmarshal(data, &parts); err != nil {
			return err
		}
		u.Parts = parts
	}
	return nil
}

// OpenAIResponseInputMessageContentPartUnionFlat 内容块 union 扁平化实现。
type OpenAIResponseInputMessageContentPartUnionFlat struct {
	Type     string `json:"type,omitempty"` // "input_text" / "input_image" / "input_file"
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	Detail   string `json:"detail,omitempty"`
	FileID   string `json:"file_id,omitempty"`
}

func (OpenAIResponseInputMessageContentPartUnionFlat) isResponseInputMessageContentPart() {}

func (u OpenAIResponseInputMessageContentPartUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "input_text":
		return json.Marshal(struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{Type: "input_text", Text: u.Text})
	case "input_image":
		return json.Marshal(struct {
			Type     string `json:"type"`
			ImageURL string `json:"image_url,omitempty"`
			Detail   string `json:"detail,omitempty"`
		}{Type: "input_image", ImageURL: u.ImageURL, Detail: u.Detail})
	case "input_file":
		return json.Marshal(struct {
			Type   string `json:"type"`
			FileID string `json:"file_id"`
		}{Type: "input_file", FileID: u.FileID})
	}
	return nil, nil
}

func (u *OpenAIResponseInputMessageContentPartUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "input_text":
		var v struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Text = v.Text
	case "input_image":
		var v struct {
			ImageURL string `json:"image_url,omitempty"`
			Detail   string `json:"detail,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ImageURL = v.ImageURL
		u.Detail = v.Detail
	case "input_file":
		var v struct {
			FileID string `json:"file_id"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.FileID = v.FileID
	}
	return nil
}

// OpenAIResponseInputComputerCallOutputContentUnionFlat 电脑使用输出内容 union 扁平化实现。
type OpenAIResponseInputComputerCallOutputContentUnionFlat struct {
	Type     string `json:"type,omitempty"` // "computer_screenshot"
	FileID   string `json:"file_id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

func (OpenAIResponseInputComputerCallOutputContentUnionFlat) isResponseInputComputerCallOutputContent() {
}

func (u OpenAIResponseInputComputerCallOutputContentUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "computer_screenshot":
		return json.Marshal(struct {
			Type     string `json:"type"`
			FileID   string `json:"file_id,omitempty"`
			ImageURL string `json:"image_url,omitempty"`
		}{Type: "computer_screenshot", FileID: u.FileID, ImageURL: u.ImageURL})
	}
	return nil, nil
}

func (u *OpenAIResponseInputComputerCallOutputContentUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "computer_screenshot":
		var v struct {
			FileID   string `json:"file_id,omitempty"`
			ImageURL string `json:"image_url,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.FileID = v.FileID
		u.ImageURL = v.ImageURL
	}
	return nil
}

// OpenAIResponseInputFunctionCallOutputUnionFlat 函数调用输出内容 union 扁平化实现。
type OpenAIResponseInputFunctionCallOutputUnionFlat struct {
	Type string `json:"type,omitempty"` // "output_text"
	Text string `json:"text,omitempty"`
}

func (OpenAIResponseInputFunctionCallOutputUnionFlat) isResponseInputFunctionCallOutput() {}

func (u OpenAIResponseInputFunctionCallOutputUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "output_text":
		return json.Marshal(struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{Type: "output_text", Text: u.Text})
	}
	return nil, nil
}

func (u *OpenAIResponseInputFunctionCallOutputUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "output_text":
		var v struct {
			Text string `json:"text"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.Text = v.Text
	}
	return nil
}

// OpenAIResponseInputItemUnionFlat InputItem union 扁平化实现（30 变体）。
type OpenAIResponseInputItemUnionFlat struct {
	Type string `json:"type,omitempty"`

	// message
	MessageRole    OpenAIResponseInputMessageItemRole          `json:"role,omitempty"`
	MessageContent *OpenAIResponseInputMessageContentUnionFlat `json:"content,omitempty"`
	MessageStatus  string                                      `json:"status,omitempty"`

	// function_call
	FunctionCallID        string `json:"id,omitempty"`
	FunctionCallCallID    string `json:"call_id,omitempty"`
	FunctionCallName      string `json:"name,omitempty"`
	FunctionCallArguments string `json:"arguments,omitempty"`
	FunctionCallStatus    string `json:"status,omitempty"`

	// function_call_output (P3)
	FunctionCallOutputID     string                                          `json:"id,omitempty"`
	FunctionCallOutputCallID string                                          `json:"call_id,omitempty"`
	FunctionCallOutputOutput *OpenAIResponseInputFunctionCallOutputUnionFlat `json:"output,omitempty"`
	FunctionCallOutputStatus string                                          `json:"status,omitempty"`

	// file_search_call
	FileSearchCallID      string                                `json:"id,omitempty"`
	FileSearchCallQueries []string                              `json:"queries,omitempty"`
	FileSearchCallStatus  string                                `json:"status,omitempty"`
	FileSearchCallResults []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`

	// web_search_call
	WebSearchCallID     string                              `json:"id,omitempty"`
	WebSearchCallStatus string                              `json:"status,omitempty"`
	WebSearchCallAction *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`

	// computer_call
	ComputerCallID     string                             `json:"id,omitempty"`
	ComputerCallCallID string                             `json:"call_id,omitempty"`
	ComputerCallAction *OpenAIResponseInputComputerAction `json:"action,omitempty"`
	ComputerCallStatus string                             `json:"status,omitempty"`

	// reasoning
	ReasoningID      string   `json:"id,omitempty"`
	ReasoningSummary []string `json:"summary,omitempty"`

	// item_reference
	ItemReferenceID string `json:"id,omitempty"`

	// apply_patch_call
	ApplyPatchCallID        string                                  `json:"id,omitempty"`
	ApplyPatchCallCallID    string                                  `json:"call_id,omitempty"`
	ApplyPatchCallStatus    string                                  `json:"status,omitempty"`
	ApplyPatchCallOperation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`

	// mcp_call
	McpCallID          string                       `json:"id,omitempty"`
	McpCallName        string                       `json:"name,omitempty"`
	McpCallServerLabel string                       `json:"server_label,omitempty"`
	McpCallArguments   string                       `json:"arguments,omitempty"`
	McpCallOutput      string                       `json:"output,omitempty"`
	McpCallError       *OpenAIResponseInputMcpError `json:"error,omitempty"`
	McpCallStatus      string                       `json:"status,omitempty"`

	// code_interpreter_call
	CodeInterpreterCallID          string                                     `json:"id,omitempty"`
	CodeInterpreterCallCallID      string                                     `json:"call_id,omitempty"`
	CodeInterpreterCallCode        string                                     `json:"code,omitempty"`
	CodeInterpreterCallOutputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`
	CodeInterpreterCallStatus      string                                     `json:"status,omitempty"`
	CodeInterpreterCallContainerID string                                     `json:"container_id,omitempty"`

	// image_generation_call
	ImageGenerationCallID     string `json:"id,omitempty"`
	ImageGenerationCallResult string `json:"result,omitempty"`
	ImageGenerationCallStatus string `json:"status,omitempty"`

	// tool_search_call
	ToolSearchCallID        string   `json:"id,omitempty"`
	ToolSearchCallQueries   []string `json:"queries,omitempty"`
	ToolSearchCallStatus    string   `json:"status,omitempty"`
	ToolSearchCallResults   []string `json:"results,omitempty"`
	ToolSearchCallExecution string   `json:"execution,omitempty"`

	// tool_search_output
	ToolSearchOutputID    string                                    `json:"id,omitempty"`
	ToolSearchOutputTools []OpenAIResponsesInputToolSearchToolEntry `json:"tools,omitempty"`

	// computer_call_output
	ComputerCallOutputID                       string                                                         `json:"id,omitempty"`
	ComputerCallOutputCallID                   string                                                         `json:"call_id,omitempty"`
	ComputerCallOutputOutput                   *OpenAIResponseInputComputerCallOutputContentUnionFlat         `json:"output,omitempty"`
	ComputerCallOutputAcknowledgedSafetyChecks []OpenAIResponseInputComputerCallOutputAcknowledgedSafetyCheck `json:"acknowledged_safety_checks,omitempty"`
	ComputerCallOutputStatus                   string                                                         `json:"status,omitempty"`

	// local_shell_call
	LocalShellCallID     string                               `json:"id,omitempty"`
	LocalShellCallCallID string                               `json:"call_id,omitempty"`
	LocalShellCallAction *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`
	LocalShellCallStatus string                               `json:"status,omitempty"`

	// local_shell_call_output
	LocalShellCallOutputID     string `json:"id,omitempty"`
	LocalShellCallOutputOutput string `json:"output,omitempty"`

	// shell_call
	ShellCallID     string                          `json:"id,omitempty"`
	ShellCallCallID string                          `json:"call_id,omitempty"`
	ShellCallAction *OpenAIResponseInputShellAction `json:"action,omitempty"`
	ShellCallStatus string                          `json:"status,omitempty"`

	// shell_call_output
	ShellCallOutputID              string `json:"id,omitempty"`
	ShellCallOutputMaxOutputLength *int64 `json:"max_output_length,omitempty"`
	ShellCallOutputOutput          string `json:"output,omitempty"`

	// apply_patch_call_output
	ApplyPatchCallOutputID     string `json:"id,omitempty"`
	ApplyPatchCallOutputCallID string `json:"call_id,omitempty"`
	ApplyPatchCallOutputStatus string `json:"status,omitempty"`

	// mcp_approval_request
	McpApprovalRequestID                string `json:"id,omitempty"`
	McpApprovalRequestServerLabel       string `json:"server_label,omitempty"`
	McpApprovalRequestName              string `json:"name,omitempty"`
	McpApprovalRequestArguments         string `json:"arguments,omitempty"`
	McpApprovalRequestApprovalRequestID string `json:"approval_request_id,omitempty"`

	// mcp_approval_response
	McpApprovalResponseApprovalRequestID string `json:"approval_request_id,omitempty"`
	McpApprovalResponseApprove           bool   `json:"approve"`
	McpApprovalResponseReason            string `json:"reason,omitempty"`

	// compaction_trigger
	CompactionTriggerID string `json:"id,omitempty"`

	// additional_tools
	AdditionalToolsID    string                         `json:"id,omitempty"`
	AdditionalToolsRole  string                         `json:"role,omitempty"`
	AdditionalToolsTools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`

	// custom_tool_call
	CustomToolCallID     string `json:"id,omitempty"`
	CustomToolCallCallID string `json:"call_id,omitempty"`
	CustomToolCallName   string `json:"name,omitempty"`
	CustomToolCallInput  string `json:"input,omitempty"`
	CustomToolCallStatus string `json:"status,omitempty"`

	// custom_tool_call_output
	CustomToolCallOutputID     string `json:"id,omitempty"`
	CustomToolCallOutputCallID string `json:"call_id,omitempty"`
	CustomToolCallOutputOutput string `json:"output,omitempty"`
	CustomToolCallOutputStatus string `json:"status,omitempty"`

	// easy_input_message (EasyInputMessage)
	EasyInputMessageRole    OpenAIResponseInputMessageItemRole          `json:"role,omitempty"`
	EasyInputMessageContent *OpenAIResponseInputMessageContentUnionFlat `json:"content,omitempty"`
}

func (OpenAIResponseInputItemUnionFlat) isResponseInputItem() {}

func (u OpenAIResponseInputItemUnionFlat) MarshalJSON() ([]byte, error) {
	if u.Type == "" {
		return nil, nil
	}
	switch u.Type {
	case "message":
		return json.Marshal(struct {
			Type    string                                      `json:"type"`
			Role    OpenAIResponseInputMessageItemRole          `json:"role"`
			Content *OpenAIResponseInputMessageContentUnionFlat `json:"content"`
			Status  string                                      `json:"status,omitempty"`
		}{Type: "message", Role: u.MessageRole, Content: u.MessageContent, Status: u.MessageStatus})
	case "function_call":
		return json.Marshal(struct {
			Type      string `json:"type"`
			ID        string `json:"id,omitempty"`
			CallID    string `json:"call_id,omitempty"`
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
			Status    string `json:"status,omitempty"`
		}{Type: "function_call", ID: u.FunctionCallID, CallID: u.FunctionCallCallID, Name: u.FunctionCallName, Arguments: u.FunctionCallArguments, Status: u.FunctionCallStatus})
	case "function_call_output":
		// P3 variant with OpenAIResponseInputFunctionCallOutputUnion output
		return json.Marshal(struct {
			Type   string                                          `json:"type"`
			ID     string                                          `json:"id,omitempty"`
			CallID string                                          `json:"call_id"`
			Output *OpenAIResponseInputFunctionCallOutputUnionFlat `json:"output"`
			Status string                                          `json:"status,omitempty"`
		}{Type: "function_call_output", ID: u.FunctionCallOutputID, CallID: u.FunctionCallOutputCallID, Output: u.FunctionCallOutputOutput, Status: u.FunctionCallOutputStatus})
	case "file_search_call":
		return json.Marshal(struct {
			Type    string                                `json:"type"`
			ID      string                                `json:"id,omitempty"`
			Queries []string                              `json:"queries,omitempty"`
			Status  string                                `json:"status,omitempty"`
			Results []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`
		}{Type: "file_search_call", ID: u.FileSearchCallID, Queries: u.FileSearchCallQueries, Status: u.FileSearchCallStatus, Results: u.FileSearchCallResults})
	case "web_search_call":
		return json.Marshal(struct {
			Type   string                              `json:"type"`
			ID     string                              `json:"id,omitempty"`
			Status string                              `json:"status,omitempty"`
			Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`
		}{Type: "web_search_call", ID: u.WebSearchCallID, Status: u.WebSearchCallStatus, Action: u.WebSearchCallAction})
	case "computer_call":
		return json.Marshal(struct {
			Type   string                             `json:"type"`
			ID     string                             `json:"id,omitempty"`
			CallID string                             `json:"call_id,omitempty"`
			Action *OpenAIResponseInputComputerAction `json:"action,omitempty"`
			Status string                             `json:"status,omitempty"`
		}{Type: "computer_call", ID: u.ComputerCallID, CallID: u.ComputerCallCallID, Action: u.ComputerCallAction, Status: u.ComputerCallStatus})
	case "reasoning":
		return json.Marshal(struct {
			Type    string   `json:"type"`
			ID      string   `json:"id,omitempty"`
			Summary []string `json:"summary,omitempty"`
		}{Type: "reasoning", ID: u.ReasoningID, Summary: u.ReasoningSummary})
	case "item_reference":
		return json.Marshal(struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		}{Type: "item_reference", ID: u.ItemReferenceID})
	case "apply_patch_call":
		return json.Marshal(struct {
			Type      string                                  `json:"type"`
			ID        string                                  `json:"id,omitempty"`
			CallID    string                                  `json:"call_id,omitempty"`
			Status    string                                  `json:"status,omitempty"`
			Operation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`
		}{Type: "apply_patch_call", ID: u.ApplyPatchCallID, CallID: u.ApplyPatchCallCallID, Status: u.ApplyPatchCallStatus, Operation: u.ApplyPatchCallOperation})
	case "mcp_call":
		return json.Marshal(struct {
			Type        string                       `json:"type"`
			ID          string                       `json:"id,omitempty"`
			Name        string                       `json:"name,omitempty"`
			ServerLabel string                       `json:"server_label,omitempty"`
			Arguments   string                       `json:"arguments,omitempty"`
			Output      string                       `json:"output,omitempty"`
			Error       *OpenAIResponseInputMcpError `json:"error,omitempty"`
			Status      string                       `json:"status,omitempty"`
		}{Type: "mcp_call", ID: u.McpCallID, Name: u.McpCallName, ServerLabel: u.McpCallServerLabel, Arguments: u.McpCallArguments, Output: u.McpCallOutput, Error: u.McpCallError, Status: u.McpCallStatus})
	case "code_interpreter_call":
		return json.Marshal(struct {
			Type        string                                     `json:"type"`
			ID          string                                     `json:"id,omitempty"`
			CallID      string                                     `json:"call_id,omitempty"`
			Code        string                                     `json:"code,omitempty"`
			Outputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`
			Status      string                                     `json:"status,omitempty"`
			ContainerID string                                     `json:"container_id,omitempty"`
		}{Type: "code_interpreter_call", ID: u.CodeInterpreterCallID, CallID: u.CodeInterpreterCallCallID, Code: u.CodeInterpreterCallCode, Outputs: u.CodeInterpreterCallOutputs, Status: u.CodeInterpreterCallStatus, ContainerID: u.CodeInterpreterCallContainerID})
	case "image_generation_call":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id,omitempty"`
			Result string `json:"result,omitempty"`
			Status string `json:"status,omitempty"`
		}{Type: "image_generation_call", ID: u.ImageGenerationCallID, Result: u.ImageGenerationCallResult, Status: u.ImageGenerationCallStatus})
	case "tool_search_call":
		return json.Marshal(struct {
			Type      string   `json:"type"`
			ID        string   `json:"id,omitempty"`
			Queries   []string `json:"queries,omitempty"`
			Status    string   `json:"status,omitempty"`
			Results   []string `json:"results,omitempty"`
			Execution string   `json:"execution,omitempty"`
		}{Type: "tool_search_call", ID: u.ToolSearchCallID, Queries: u.ToolSearchCallQueries, Status: u.ToolSearchCallStatus, Results: u.ToolSearchCallResults, Execution: u.ToolSearchCallExecution})
	case "tool_search_output":
		return json.Marshal(struct {
			Type  string                                    `json:"type"`
			ID    string                                    `json:"id,omitempty"`
			Tools []OpenAIResponsesInputToolSearchToolEntry `json:"tools,omitempty"`
		}{Type: "tool_search_output", ID: u.ToolSearchOutputID, Tools: u.ToolSearchOutputTools})
	case "computer_call_output":
		return json.Marshal(struct {
			Type                     string                                                         `json:"type"`
			ID                       string                                                         `json:"id,omitempty"`
			CallID                   string                                                         `json:"call_id,omitempty"`
			Output                   *OpenAIResponseInputComputerCallOutputContentUnionFlat         `json:"output,omitempty"`
			AcknowledgedSafetyChecks []OpenAIResponseInputComputerCallOutputAcknowledgedSafetyCheck `json:"acknowledged_safety_checks,omitempty"`
			Status                   string                                                         `json:"status,omitempty"`
		}{Type: "computer_call_output", ID: u.ComputerCallOutputID, CallID: u.ComputerCallOutputCallID, Output: u.ComputerCallOutputOutput, AcknowledgedSafetyChecks: u.ComputerCallOutputAcknowledgedSafetyChecks, Status: u.ComputerCallOutputStatus})
	case "local_shell_call":
		return json.Marshal(struct {
			Type   string                               `json:"type"`
			ID     string                               `json:"id,omitempty"`
			CallID string                               `json:"call_id,omitempty"`
			Action *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`
			Status string                               `json:"status,omitempty"`
		}{Type: "local_shell_call", ID: u.LocalShellCallID, CallID: u.LocalShellCallCallID, Action: u.LocalShellCallAction, Status: u.LocalShellCallStatus})
	case "local_shell_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id,omitempty"`
			Output string `json:"output,omitempty"`
		}{Type: "local_shell_call_output", ID: u.LocalShellCallOutputID, Output: u.LocalShellCallOutputOutput})
	case "shell_call":
		return json.Marshal(struct {
			Type   string                          `json:"type"`
			ID     string                          `json:"id,omitempty"`
			CallID string                          `json:"call_id,omitempty"`
			Action *OpenAIResponseInputShellAction `json:"action,omitempty"`
			Status string                          `json:"status,omitempty"`
		}{Type: "shell_call", ID: u.ShellCallID, CallID: u.ShellCallCallID, Action: u.ShellCallAction, Status: u.ShellCallStatus})
	case "shell_call_output":
		return json.Marshal(struct {
			Type            string `json:"type"`
			ID              string `json:"id,omitempty"`
			MaxOutputLength *int64 `json:"max_output_length,omitempty"`
			Output          string `json:"output,omitempty"`
		}{Type: "shell_call_output", ID: u.ShellCallOutputID, MaxOutputLength: u.ShellCallOutputMaxOutputLength, Output: u.ShellCallOutputOutput})
	case "apply_patch_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id,omitempty"`
			CallID string `json:"call_id,omitempty"`
			Status string `json:"status,omitempty"`
		}{Type: "apply_patch_call_output", ID: u.ApplyPatchCallOutputID, CallID: u.ApplyPatchCallOutputCallID, Status: u.ApplyPatchCallOutputStatus})
	case "mcp_approval_request":
		return json.Marshal(struct {
			Type              string `json:"type"`
			ID                string `json:"id,omitempty"`
			ServerLabel       string `json:"server_label,omitempty"`
			Name              string `json:"name,omitempty"`
			Arguments         string `json:"arguments,omitempty"`
			ApprovalRequestID string `json:"approval_request_id,omitempty"`
		}{Type: "mcp_approval_request", ID: u.McpApprovalRequestID, ServerLabel: u.McpApprovalRequestServerLabel, Name: u.McpApprovalRequestName, Arguments: u.McpApprovalRequestArguments, ApprovalRequestID: u.McpApprovalRequestApprovalRequestID})
	case "mcp_approval_response":
		return json.Marshal(struct {
			Type              string `json:"type"`
			ApprovalRequestID string `json:"approval_request_id"`
			Approve           bool   `json:"approve"`
			Reason            string `json:"reason,omitempty"`
		}{Type: "mcp_approval_response", ApprovalRequestID: u.McpApprovalResponseApprovalRequestID, Approve: u.McpApprovalResponseApprove, Reason: u.McpApprovalResponseReason})
	case "compaction_trigger":
		return json.Marshal(struct {
			Type string `json:"type"`
			ID   string `json:"id,omitempty"`
		}{Type: "compaction_trigger", ID: u.CompactionTriggerID})
	case "additional_tools":
		return json.Marshal(struct {
			Type  string                         `json:"type"`
			ID    string                         `json:"id,omitempty"`
			Role  string                         `json:"role,omitempty"`
			Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
		}{Type: "additional_tools", ID: u.AdditionalToolsID, Role: u.AdditionalToolsRole, Tools: u.AdditionalToolsTools})
	case "custom_tool_call":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id,omitempty"`
			CallID string `json:"call_id,omitempty"`
			Name   string `json:"name"`
			Input  string `json:"input"`
			Status string `json:"status,omitempty"`
		}{Type: "custom_tool_call", ID: u.CustomToolCallID, CallID: u.CustomToolCallCallID, Name: u.CustomToolCallName, Input: u.CustomToolCallInput, Status: u.CustomToolCallStatus})
	case "custom_tool_call_output":
		return json.Marshal(struct {
			Type   string `json:"type"`
			ID     string `json:"id,omitempty"`
			CallID string `json:"call_id,omitempty"`
			Output string `json:"output"`
			Status string `json:"status,omitempty"`
		}{Type: "custom_tool_call_output", ID: u.CustomToolCallOutputID, CallID: u.CustomToolCallOutputCallID, Output: u.CustomToolCallOutputOutput, Status: u.CustomToolCallOutputStatus})
	}
	return nil, nil
}

func (u *OpenAIResponseInputItemUnionFlat) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.Type = tmp.Type
	switch tmp.Type {
	case "message":
		var v struct {
			Type    string                                      `json:"type"`
			Role    OpenAIResponseInputMessageItemRole          `json:"role"`
			Content *OpenAIResponseInputMessageContentUnionFlat `json:"content"`
			Status  string                                      `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.MessageRole = v.Role
		u.MessageContent = v.Content
		u.MessageStatus = v.Status
	case "function_call":
		var v struct {
			ID        string `json:"id,omitempty"`
			CallID    string `json:"call_id,omitempty"`
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
			Status    string `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.FunctionCallID = v.ID
		u.FunctionCallCallID = v.CallID
		u.FunctionCallName = v.Name
		u.FunctionCallArguments = v.Arguments
		u.FunctionCallStatus = v.Status
	case "function_call_output":
		var v struct {
			ID     string                                          `json:"id,omitempty"`
			CallID string                                          `json:"call_id"`
			Output *OpenAIResponseInputFunctionCallOutputUnionFlat `json:"output"`
			Status string                                          `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.FunctionCallOutputID = v.ID
		u.FunctionCallOutputCallID = v.CallID
		u.FunctionCallOutputOutput = v.Output
		u.FunctionCallOutputStatus = v.Status
	case "file_search_call":
		var v struct {
			ID      string                                `json:"id,omitempty"`
			Queries []string                              `json:"queries,omitempty"`
			Status  string                                `json:"status,omitempty"`
			Results []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.FileSearchCallID = v.ID
		u.FileSearchCallQueries = v.Queries
		u.FileSearchCallStatus = v.Status
		u.FileSearchCallResults = v.Results
	case "web_search_call":
		var v struct {
			ID     string                              `json:"id,omitempty"`
			Status string                              `json:"status,omitempty"`
			Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.WebSearchCallID = v.ID
		u.WebSearchCallStatus = v.Status
		u.WebSearchCallAction = v.Action
	case "computer_call":
		var v struct {
			ID     string                             `json:"id,omitempty"`
			CallID string                             `json:"call_id,omitempty"`
			Action *OpenAIResponseInputComputerAction `json:"action,omitempty"`
			Status string                             `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ComputerCallID = v.ID
		u.ComputerCallCallID = v.CallID
		u.ComputerCallAction = v.Action
		u.ComputerCallStatus = v.Status
	case "reasoning":
		var v struct {
			ID      string   `json:"id,omitempty"`
			Summary []string `json:"summary,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ReasoningID = v.ID
		u.ReasoningSummary = v.Summary
	case "item_reference":
		var v struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ItemReferenceID = v.ID
	case "apply_patch_call":
		var v struct {
			ID        string                                  `json:"id,omitempty"`
			CallID    string                                  `json:"call_id,omitempty"`
			Status    string                                  `json:"status,omitempty"`
			Operation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ApplyPatchCallID = v.ID
		u.ApplyPatchCallCallID = v.CallID
		u.ApplyPatchCallStatus = v.Status
		u.ApplyPatchCallOperation = v.Operation
	case "mcp_call":
		var v struct {
			ID          string                       `json:"id,omitempty"`
			Name        string                       `json:"name,omitempty"`
			ServerLabel string                       `json:"server_label,omitempty"`
			Arguments   string                       `json:"arguments,omitempty"`
			Output      string                       `json:"output,omitempty"`
			Error       *OpenAIResponseInputMcpError `json:"error,omitempty"`
			Status      string                       `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.McpCallID = v.ID
		u.McpCallName = v.Name
		u.McpCallServerLabel = v.ServerLabel
		u.McpCallArguments = v.Arguments
		u.McpCallOutput = v.Output
		u.McpCallError = v.Error
		u.McpCallStatus = v.Status
	case "code_interpreter_call":
		var v struct {
			ID          string                                     `json:"id,omitempty"`
			CallID      string                                     `json:"call_id,omitempty"`
			Code        string                                     `json:"code,omitempty"`
			Outputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`
			Status      string                                     `json:"status,omitempty"`
			ContainerID string                                     `json:"container_id,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.CodeInterpreterCallID = v.ID
		u.CodeInterpreterCallCallID = v.CallID
		u.CodeInterpreterCallCode = v.Code
		u.CodeInterpreterCallOutputs = v.Outputs
		u.CodeInterpreterCallStatus = v.Status
		u.CodeInterpreterCallContainerID = v.ContainerID
	case "image_generation_call":
		var v struct {
			ID     string `json:"id,omitempty"`
			Result string `json:"result,omitempty"`
			Status string `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ImageGenerationCallID = v.ID
		u.ImageGenerationCallResult = v.Result
		u.ImageGenerationCallStatus = v.Status
	case "tool_search_call":
		var v struct {
			ID        string   `json:"id,omitempty"`
			Queries   []string `json:"queries,omitempty"`
			Status    string   `json:"status,omitempty"`
			Results   []string `json:"results,omitempty"`
			Execution string   `json:"execution,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ToolSearchCallID = v.ID
		u.ToolSearchCallQueries = v.Queries
		u.ToolSearchCallStatus = v.Status
		u.ToolSearchCallResults = v.Results
		u.ToolSearchCallExecution = v.Execution
	case "tool_search_output":
		var v struct {
			ID    string                                    `json:"id,omitempty"`
			Tools []OpenAIResponsesInputToolSearchToolEntry `json:"tools,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ToolSearchOutputID = v.ID
		u.ToolSearchOutputTools = v.Tools
	case "computer_call_output":
		var v struct {
			ID                       string                                                         `json:"id,omitempty"`
			CallID                   string                                                         `json:"call_id,omitempty"`
			Output                   *OpenAIResponseInputComputerCallOutputContentUnionFlat         `json:"output,omitempty"`
			AcknowledgedSafetyChecks []OpenAIResponseInputComputerCallOutputAcknowledgedSafetyCheck `json:"acknowledged_safety_checks,omitempty"`
			Status                   string                                                         `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ComputerCallOutputID = v.ID
		u.ComputerCallOutputCallID = v.CallID
		u.ComputerCallOutputOutput = v.Output
		u.ComputerCallOutputAcknowledgedSafetyChecks = v.AcknowledgedSafetyChecks
		u.ComputerCallOutputStatus = v.Status
	case "local_shell_call":
		var v struct {
			ID     string                               `json:"id,omitempty"`
			CallID string                               `json:"call_id,omitempty"`
			Action *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`
			Status string                               `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.LocalShellCallID = v.ID
		u.LocalShellCallCallID = v.CallID
		u.LocalShellCallAction = v.Action
		u.LocalShellCallStatus = v.Status
	case "local_shell_call_output":
		var v struct {
			ID     string `json:"id,omitempty"`
			Output string `json:"output,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.LocalShellCallOutputID = v.ID
		u.LocalShellCallOutputOutput = v.Output
	case "shell_call":
		var v struct {
			ID     string                          `json:"id,omitempty"`
			CallID string                          `json:"call_id,omitempty"`
			Action *OpenAIResponseInputShellAction `json:"action,omitempty"`
			Status string                          `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ShellCallID = v.ID
		u.ShellCallCallID = v.CallID
		u.ShellCallAction = v.Action
		u.ShellCallStatus = v.Status
	case "shell_call_output":
		var v struct {
			ID              string `json:"id,omitempty"`
			MaxOutputLength *int64 `json:"max_output_length,omitempty"`
			Output          string `json:"output,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ShellCallOutputID = v.ID
		u.ShellCallOutputMaxOutputLength = v.MaxOutputLength
		u.ShellCallOutputOutput = v.Output
	case "apply_patch_call_output":
		var v struct {
			ID     string `json:"id,omitempty"`
			CallID string `json:"call_id,omitempty"`
			Status string `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.ApplyPatchCallOutputID = v.ID
		u.ApplyPatchCallOutputCallID = v.CallID
		u.ApplyPatchCallOutputStatus = v.Status
	case "mcp_approval_request":
		var v struct {
			ID                string `json:"id,omitempty"`
			ServerLabel       string `json:"server_label,omitempty"`
			Name              string `json:"name,omitempty"`
			Arguments         string `json:"arguments,omitempty"`
			ApprovalRequestID string `json:"approval_request_id,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.McpApprovalRequestID = v.ID
		u.McpApprovalRequestServerLabel = v.ServerLabel
		u.McpApprovalRequestName = v.Name
		u.McpApprovalRequestArguments = v.Arguments
		u.McpApprovalRequestApprovalRequestID = v.ApprovalRequestID
	case "mcp_approval_response":
		var v struct {
			ApprovalRequestID string `json:"approval_request_id"`
			Approve           bool   `json:"approve"`
			Reason            string `json:"reason,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.McpApprovalResponseApprovalRequestID = v.ApprovalRequestID
		u.McpApprovalResponseApprove = v.Approve
		u.McpApprovalResponseReason = v.Reason
	case "compaction_trigger":
		var v struct {
			ID string `json:"id,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.CompactionTriggerID = v.ID
	case "additional_tools":
		var v struct {
			ID    string                         `json:"id,omitempty"`
			Role  string                         `json:"role,omitempty"`
			Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.AdditionalToolsID = v.ID
		u.AdditionalToolsRole = v.Role
		u.AdditionalToolsTools = v.Tools
	case "custom_tool_call":
		var v struct {
			ID     string `json:"id,omitempty"`
			CallID string `json:"call_id,omitempty"`
			Name   string `json:"name"`
			Input  string `json:"input"`
			Status string `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.CustomToolCallID = v.ID
		u.CustomToolCallCallID = v.CallID
		u.CustomToolCallName = v.Name
		u.CustomToolCallInput = v.Input
		u.CustomToolCallStatus = v.Status
	case "custom_tool_call_output":
		var v struct {
			ID     string `json:"id,omitempty"`
			CallID string `json:"call_id,omitempty"`
			Output string `json:"output"`
			Status string `json:"status,omitempty"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		u.CustomToolCallOutputID = v.ID
		u.CustomToolCallOutputCallID = v.CallID
		u.CustomToolCallOutputOutput = v.Output
		u.CustomToolCallOutputStatus = v.Status
	}
	return nil
}

// ============================================================================
// Original Union Interfaces and Implementations
// ============================================================================

// OpenAIResponseInputItemUnion InputItem union 接口。
//
// 24 种实现（4.6+）：MessageItem / FunctionToolCallItem / FunctionToolCallOutputItem /
// FileSearchToolCallItem / WebSearchToolCallItem / ComputerToolCallItem /
// ReasoningItem / ItemReferenceItem / ApplyPatchToolCallItem / McpToolCallItem /
// CodeInterpreterToolCallItem / ImageGenerationCallItem /
// ToolSearchCallItem / ToolSearchOutputItem / ComputerCallOutputItem /
// LocalShellCallItem / LocalShellCallOutputItem / ShellCallItem / ShellCallOutputItem /
// ApplyPatchCallOutputItem / McpApprovalRequestItem / McpApprovalResponseItem /
// CompactionTriggerItem / EasyInputMessageItem。
type OpenAIResponseInputItemUnion interface {
	isResponseInputItem()
}

// OpenAIResponseInputMessageItemRole InputItem role 字段类型。
type OpenAIResponseInputMessageItemRole string

const (
	OpenAIResponseInputMessageItemRoleUser      OpenAIResponseInputMessageItemRole = "user"
	OpenAIResponseInputMessageItemRoleSystem    OpenAIResponseInputMessageItemRole = "system"
	OpenAIResponseInputMessageItemRoleAssistant OpenAIResponseInputMessageItemRole = "assistant"
)

// OpenAIResponseInputMessageItem 输入消息项。
type OpenAIResponseInputMessageItem struct {
	// Type 固定 "message"。
	Type string `json:"type"` // "message"
	// Role 消息角色。
	Role OpenAIResponseInputMessageItemRole `json:"role"`
	// Content 消息内容（string | []ContentPart）。
	Content *OpenAIResponseInputMessageContentUnionFlat `json:"content"`
	// Status 状态（in_progress / completed / incomplete）。
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputMessageItem) isResponseInputItem() {}

// OpenAIResponseInputMessageContentUnion Input message content union。
type OpenAIResponseInputMessageContentUnion interface {
	isResponseInputMessageContent()
}

type OpenAIResponseInputMessageContentString string

func (OpenAIResponseInputMessageContentString) isResponseInputMessageContent() {}

// OpenAIResponseInputMessageContentParts 内容块数组（input_text / input_image / input_file）。
type OpenAIResponseInputMessageContentParts []OpenAIResponseInputMessageContentPartUnion

func (OpenAIResponseInputMessageContentParts) isResponseInputMessageContent() {}

// OpenAIResponseInputMessageContentPartUnion 内容块 union（P2-4）。
type OpenAIResponseInputMessageContentPartUnion interface {
	isResponseInputMessageContentPart()
}

// OpenAIResponseInputMessageContentPartText 文本内容块。
type OpenAIResponseInputMessageContentPartText struct {
	Type string `json:"type"` // "input_text"
	Text string `json:"text"`
}

func (OpenAIResponseInputMessageContentPartText) isResponseInputMessageContentPart() {}

// OpenAIResponseInputMessageContentPartImage 图像内容块。
type OpenAIResponseInputMessageContentPartImage struct {
	Type     string `json:"type"` // "input_image"
	ImageURL string `json:"image_url,omitempty"`
	Detail   string `json:"detail,omitempty"` // "low" / "high" / "auto"
}

func (OpenAIResponseInputMessageContentPartImage) isResponseInputMessageContentPart() {}

// OpenAIResponseInputMessageContentPartFile 文件内容块。
type OpenAIResponseInputMessageContentPartFile struct {
	Type   string `json:"type"` // "input_file"
	FileID string `json:"file_id"`
}

func (OpenAIResponseInputMessageContentPartFile) isResponseInputMessageContentPart() {}

// OpenAIResponseInputFunctionToolCallItem 函数工具调用项。
type OpenAIResponseInputFunctionToolCallItem struct {
	// Type 固定 "function_call"。
	Type string `json:"type"` // "function_call"
	// ID 工具调用 ID。
	ID string `json:"id,omitempty"`
	// CallID 引用 ID。
	CallID string `json:"call_id,omitempty"`
	// Name 函数名。
	Name string `json:"name"`
	// Arguments 参数 JSON 字符串。
	Arguments string `json:"arguments"`
	// Status 状态。
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputFunctionToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputFunctionToolCallOutputItem 函数工具调用输出项。
type OpenAIResponseInputFunctionToolCallOutputItem struct {
	// Type 固定 "function_call_output"。
	Type string `json:"type"` // "function_call_output"
	// CallID 对应的工具调用 ID。
	CallID string `json:"call_id"`
	// Output 输出内容。
	Output string `json:"output"`
	// Status 状态。
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputFunctionToolCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputFileSearchToolCallItem 文件搜索工具调用项。
type OpenAIResponseInputFileSearchToolCallItem struct {
	Type    string                                `json:"type"` // "file_search_call"
	ID      string                                `json:"id,omitempty"`
	Queries []string                              `json:"queries,omitempty"`
	Status  string                                `json:"status,omitempty"`
	Results []OpenAIResponseInputFileSearchResult `json:"results,omitempty"`
}

func (OpenAIResponseInputFileSearchToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputFileSearchResult 文件搜索结果项。
type OpenAIResponseInputFileSearchResult struct {
	FileID   string  `json:"file_id"`
	Filename string  `json:"filename,omitempty"`
	Score    float64 `json:"score,omitempty"`
	Text     string  `json:"text,omitempty"`
}

// OpenAIResponseInputWebSearchToolCallItem Web 搜索工具调用项。
type OpenAIResponseInputWebSearchToolCallItem struct {
	Type   string                              `json:"type"` // "web_search_call"
	ID     string                              `json:"id,omitempty"`
	Status string                              `json:"status,omitempty"`
	Action *OpenAIResponseInputWebSearchAction `json:"action,omitempty"`
}

func (OpenAIResponseInputWebSearchToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputWebSearchAction Web 搜索动作。
type OpenAIResponseInputWebSearchAction struct {
	Type  string `json:"type"`            // "search" / "open_page" / "find_in_page"
	Query string `json:"query,omitempty"` // search 时
	URL   string `json:"url,omitempty"`   // open_page / find_in_page 时
}

// OpenAIResponseInputComputerToolCallItem 电脑使用工具调用项。
type OpenAIResponseInputComputerToolCallItem struct {
	Type   string                             `json:"type"` // "computer_call"
	ID     string                             `json:"id,omitempty"`
	CallID string                             `json:"call_id,omitempty"`
	Action *OpenAIResponseInputComputerAction `json:"action,omitempty"`
	Status string                             `json:"status,omitempty"`
}

func (OpenAIResponseInputComputerToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputComputerAction 电脑使用动作。
type OpenAIResponseInputComputerAction struct {
	Type    string   `json:"type"` // "click" / "type" / "screenshot" 等
	X       *int64   `json:"x,omitempty"`
	Y       *int64   `json:"y,omitempty"`
	Button  string   `json:"button,omitempty"`
	Text    string   `json:"text,omitempty"`
	Keys    []string `json:"keys,omitempty"`
	ScrollX *int64   `json:"scroll_x,omitempty"`
	ScrollY *int64   `json:"scroll_y,omitempty"`
}

// OpenAIResponseInputReasoningItem 推理项。
type OpenAIResponseInputReasoningItem struct {
	Type    string   `json:"type"` // "reasoning"
	ID      string   `json:"id,omitempty"`
	Summary []string `json:"summary,omitempty"`
}

func (OpenAIResponseInputReasoningItem) isResponseInputItem() {}

// OpenAIResponseInputItemReferenceItem 项目引用项。
type OpenAIResponseInputItemReferenceItem struct {
	Type string `json:"type"` // "item_reference"
	ID   string `json:"id"`
}

func (OpenAIResponseInputItemReferenceItem) isResponseInputItem() {}

// OpenAIResponseInputApplyPatchToolCallItem 应用补丁工具调用项。
type OpenAIResponseInputApplyPatchToolCallItem struct {
	Type      string                                  `json:"type"` // "apply_patch_call"
	ID        string                                  `json:"id,omitempty"`
	CallID    string                                  `json:"call_id,omitempty"`
	Status    string                                  `json:"status,omitempty"`
	Operation *OpenAIResponseInputApplyPatchOperation `json:"operation,omitempty"`
}

func (OpenAIResponseInputApplyPatchToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputApplyPatchOperation 应用补丁操作。
type OpenAIResponseInputApplyPatchOperation struct {
	Type       string `json:"type"` // "create_file" / "delete_file" / "update_file"
	Path       string `json:"path,omitempty"`
	Diff       string `json:"diff,omitempty"`
	OldText    string `json:"old_text,omitempty"`
	NewContent string `json:"new_content,omitempty"`
}

// OpenAIResponseInputMcpToolCallItem MCP 工具调用项。
type OpenAIResponseInputMcpToolCallItem struct {
	Type        string                       `json:"type"` // "mcp_call"
	ID          string                       `json:"id,omitempty"`
	Name        string                       `json:"name,omitempty"`
	ServerLabel string                       `json:"server_label,omitempty"`
	Arguments   string                       `json:"arguments,omitempty"`
	Output      string                       `json:"output,omitempty"`
	Error       *OpenAIResponseInputMcpError `json:"error,omitempty"`
	Status      string                       `json:"status,omitempty"`
}

func (OpenAIResponseInputMcpToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputMcpError MCP 错误。
type OpenAIResponseInputMcpError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// OpenAIResponseInputCodeInterpreterToolCallItem 代码解释器工具调用项。
type OpenAIResponseInputCodeInterpreterToolCallItem struct {
	Type        string                                     `json:"type"` // "code_interpreter_call"
	ID          string                                     `json:"id,omitempty"`
	CallID      string                                     `json:"call_id,omitempty"`
	Code        string                                     `json:"code,omitempty"`
	Outputs     []OpenAIResponseInputCodeInterpreterOutput `json:"outputs,omitempty"`
	Status      string                                     `json:"status,omitempty"`
	ContainerID string                                     `json:"container_id,omitempty"`
}

func (OpenAIResponseInputCodeInterpreterToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputCodeInterpreterOutput 代码解释器输出项。
type OpenAIResponseInputCodeInterpreterOutput struct {
	Type string `json:"type"` // "logs" / "image"
	Logs string `json:"logs,omitempty"`
	URL  string `json:"url,omitempty"`
}

// OpenAIResponseInputImageGenerationCallItem 图像生成调用项。
type OpenAIResponseInputImageGenerationCallItem struct {
	Type   string `json:"type"` // "image_generation_call"
	ID     string `json:"id,omitempty"`
	Result string `json:"result,omitempty"` // base64
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputImageGenerationCallItem) isResponseInputItem() {}

// === 4.6+ 扩展变体（12 种）===

// OpenAIResponseInputToolSearchCallItem 工具搜索调用项。
type OpenAIResponseInputToolSearchCallItem struct {
	Type      string   `json:"type"` // "tool_search_call"
	ID        string   `json:"id,omitempty"`
	Queries   []string `json:"queries,omitempty"`
	Status    string   `json:"status,omitempty"`
	Results   []string `json:"results,omitempty"`
	Execution string   `json:"execution,omitempty"`
}

func (OpenAIResponseInputToolSearchCallItem) isResponseInputItem() {}

// OpenAIResponseInputToolSearchOutputItem 工具搜索输出项（命名空间 / 工具列表）。
type OpenAIResponseInputToolSearchOutputItem struct {
	Type string `json:"type"` // "tool_search_output"
	ID   string `json:"id,omitempty"`
	// Tools 工具搜索返回的工具列表（OpenAIResponsesInputToolSearchToolEntry）。
	Tools []OpenAIResponsesInputToolSearchToolEntry `json:"tools,omitempty"`
}

func (OpenAIResponseInputToolSearchOutputItem) isResponseInputItem() {}

// OpenAIResponsesInputToolSearchToolEntry 工具搜索结果中的工具定义。
type OpenAIResponsesInputToolSearchToolEntry struct {
	Type        string `json:"type"` // "function"
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// OpenAIResponseInputComputerCallOutputItem 电脑使用输出项（含 acknowledged safety checks）。
type OpenAIResponseInputComputerCallOutputItem struct {
	Type                     string                                                         `json:"type"` // "computer_call_output"
	CallID                   string                                                         `json:"call_id"`
	Output                   *OpenAIResponseInputComputerCallOutputContentUnionFlat         `json:"output"`
	AcknowledgedSafetyChecks []OpenAIResponseInputComputerCallOutputAcknowledgedSafetyCheck `json:"acknowledged_safety_checks,omitempty"`
	Status                   string                                                         `json:"status,omitempty"`
}

func (OpenAIResponseInputComputerCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputComputerCallOutputContentUnion 电脑使用输出内容 union。
type OpenAIResponseInputComputerCallOutputContentUnion interface {
	isResponseInputComputerCallOutputContent()
}

// OpenAIResponseInputComputerCallOutputScreenshot 电脑使用 screenshot 输出。
type OpenAIResponseInputComputerCallOutputScreenshot struct {
	Type     string `json:"type"` // "computer_screenshot"
	FileID   string `json:"file_id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

func (OpenAIResponseInputComputerCallOutputScreenshot) isResponseInputComputerCallOutputContent() {}

// OpenAIResponseInputComputerCallOutputAcknowledgedSafetyCheck 已确认的安全检查项。
type OpenAIResponseInputComputerCallOutputAcknowledgedSafetyCheck struct {
	ID      string `json:"id,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// OpenAIResponseInputLocalShellCallItem 本地 shell 调用项（gpt-5.1 codex 用）。
type OpenAIResponseInputLocalShellCallItem struct {
	Type   string                               `json:"type"` // "local_shell_call"
	ID     string                               `json:"id,omitempty"`
	CallID string                               `json:"call_id,omitempty"`
	Action *OpenAIResponseInputLocalShellAction `json:"action,omitempty"`
	Status string                               `json:"status,omitempty"`
}

func (OpenAIResponseInputLocalShellCallItem) isResponseInputItem() {}

// OpenAIResponseInputLocalShellAction 本地 shell 动作。
type OpenAIResponseInputLocalShellAction struct {
	Type    string   `json:"type"` // "exec"
	Command string   `json:"command,omitempty"`
	Env     []string `json:"env,omitempty"`
	Timeout *int64   `json:"timeout_ms,omitempty"`
}

// OpenAIResponseInputLocalShellCallOutputItem 本地 shell 输出项。
type OpenAIResponseInputLocalShellCallOutputItem struct {
	Type string `json:"type"` // "local_shell_call_output"
	ID   string `json:"id,omitempty"`
	// Output 输出（stdout / stderr）。
	Output string `json:"output,omitempty"`
}

func (OpenAIResponseInputLocalShellCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputShellCallItem 远程 shell 调用项。
type OpenAIResponseInputShellCallItem struct {
	Type   string                          `json:"type"` // "shell_call"
	ID     string                          `json:"id,omitempty"`
	CallID string                          `json:"call_id,omitempty"`
	Action *OpenAIResponseInputShellAction `json:"action,omitempty"`
	Status string                          `json:"status,omitempty"`
}

func (OpenAIResponseInputShellCallItem) isResponseInputItem() {}

// OpenAIResponseInputShellAction 远程 shell 动作。
type OpenAIResponseInputShellAction struct {
	Type        string   `json:"type"` // "exec"
	Command     string   `json:"command,omitempty"`
	TimeLimitMs *int64   `json:"timelimitms,omitempty"`
	Env         []string `json:"env,omitempty"`
}

// OpenAIResponseInputShellCallOutputItem 远程 shell 输出项。
type OpenAIResponseInputShellCallOutputItem struct {
	Type            string `json:"type"` // "shell_call_output"
	ID              string `json:"id,omitempty"`
	MaxOutputLength *int64 `json:"max_output_length,omitempty"`
	Output          string `json:"output,omitempty"`
}

func (OpenAIResponseInputShellCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputApplyPatchCallOutputItem 应用补丁输出项。
type OpenAIResponseInputApplyPatchCallOutputItem struct {
	Type   string `json:"type"` // "apply_patch_call_output"
	CallID string `json:"call_id,omitempty"`
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputApplyPatchCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputMcpApprovalRequestItem MCP 工具审批请求项。
type OpenAIResponseInputMcpApprovalRequestItem struct {
	Type              string `json:"type"` // "mcp_approval_request"
	ID                string `json:"id,omitempty"`
	ServerLabel       string `json:"server_label,omitempty"`
	Name              string `json:"name,omitempty"`
	Arguments         string `json:"arguments,omitempty"`
	ApprovalRequestID string `json:"approval_request_id,omitempty"`
}

func (OpenAIResponseInputMcpApprovalRequestItem) isResponseInputItem() {}

// OpenAIResponseInputMcpApprovalResponseItem MCP 工具审批响应项。
type OpenAIResponseInputMcpApprovalResponseItem struct {
	Type              string `json:"type"` // "mcp_approval_response"
	ApprovalRequestID string `json:"approval_request_id"`
	Approve           bool   `json:"approve"`
	Reason            string `json:"reason,omitempty"`
}

func (OpenAIResponseInputMcpApprovalResponseItem) isResponseInputItem() {}

// OpenAIResponseInputCompactionTriggerItem 上下文压缩触发项。
type OpenAIResponseInputCompactionTriggerItem struct {
	Type string `json:"type"` // "compaction_trigger"
}

func (OpenAIResponseInputCompactionTriggerItem) isResponseInputItem() {}

// OpenAIResponseInputAdditionalToolsItem 附加工具集输入项（P3-2）。
type OpenAIResponseInputAdditionalToolsItem struct {
	// Type 固定 "additional_tools"。
	Type string `json:"type"` // "additional_tools"
	// ID 唯一 ID。
	ID string `json:"id,omitempty"`
	// Role 角色（仅 "developer" 支持）。
	Role string `json:"role,omitempty"`
	// Tools 附加工具列表。
	Tools []OpenAIResponsesToolUnionFlat `json:"tools,omitempty"`
}

func (OpenAIResponseInputAdditionalToolsItem) isResponseInputItem() {}

// OpenAIResponseInputCustomToolCallItem 自定义工具调用项（P3 补充）。
type OpenAIResponseInputCustomToolCallItem struct {
	// Type 固定 "custom_tool_call"。
	Type   string `json:"type"` // "custom_tool_call"
	ID     string `json:"id,omitempty"`
	CallID string `json:"call_id,omitempty"`
	Name   string `json:"name"`
	Input  string `json:"input"`
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputCustomToolCallItem) isResponseInputItem() {}

// OpenAIResponseInputCustomToolCallOutputItem 自定义工具输出项（P3 补充）。
type OpenAIResponseInputCustomToolCallOutputItem struct {
	// Type 固定 "custom_tool_call_output"。
	Type   string `json:"type"` // "custom_tool_call_output"
	ID     string `json:"id,omitempty"`
	CallID string `json:"call_id,omitempty"`
	Output string `json:"output"`
	Status string `json:"status,omitempty"`
}

func (OpenAIResponseInputCustomToolCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputFunctionCallOutputItem 函数调用输出项（P3 补充）。
type OpenAIResponseInputFunctionCallOutputItem struct {
	// Type 固定 "function_call_output"。
	Type   string                                          `json:"type"` // "function_call_output"
	ID     string                                          `json:"id,omitempty"`
	CallID string                                          `json:"call_id"`
	Output *OpenAIResponseInputFunctionCallOutputUnionFlat `json:"output"`
	Status string                                          `json:"status,omitempty"`
}

func (OpenAIResponseInputFunctionCallOutputItem) isResponseInputItem() {}

// OpenAIResponseInputFunctionCallOutputUnion 函数调用输出内容 union。
type OpenAIResponseInputFunctionCallOutputUnion interface {
	isResponseInputFunctionCallOutput()
}

// OpenAIResponseInputFunctionCallOutputText 文本输出。
type OpenAIResponseInputFunctionCallOutputText struct {
	Type string `json:"type"` // "output_text"
	Text string `json:"text"`
}

func (OpenAIResponseInputFunctionCallOutputText) isResponseInputFunctionCallOutput() {}

// OpenAIResponseInputEasyInputMessageItem 轻量输入消息（EasyInputMessage）。
//
// 用于在 Responses API 入口简化传入消息（无需 ID / Status）。
type OpenAIResponseInputEasyInputMessageItem struct {
	Type    string                                      `json:"type"` // "message"
	Role    OpenAIResponseInputMessageItemRole          `json:"role"`
	Content *OpenAIResponseInputMessageContentUnionFlat `json:"content"`
}

func (OpenAIResponseInputEasyInputMessageItem) isResponseInputItem() {}
