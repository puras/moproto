package anthropic

// 代码执行 (code_execution) + bash_code_execution + text_editor_code_execution
// 这三类工具响应块内部的子结构类型。
//
// 对应官方 SDK：
//   - CodeExecutionToolResultBlockContentUnion  （message.go:1267-1288）
//   - BashCodeExecutionToolResultBlockContentUnion （message.go:297-318）
//   - TextEditorCodeExecutionToolResultBlockContentUnion （message.go:6107-6135）
//
// 这两类 Content 在外层 struct 中仍是 json.RawMessage（content_block.go 第 122/132/141 行），
// 以便服务端按需做 union 分发。本文件提供 union 中的具体 struct 形态，供调用方在 type switch
// json.RawMessage 后做精确解码/编码使用。
//
// 设计原则：仅当网关做"二次 marshal"（如把响应转发给客户端并需要 revalidate 类型）时，
// 这些子结构才被实际用到；纯透传场景下 RawMessage 已足够。

// ============== Type 别名 + 常量 ==============

// CodeExecutionOutputType code_execution_output 子块的 type 字段类型。
type CodeExecutionOutputType string

const CodeExecutionOutputTypeCodeExecutionOutput CodeExecutionOutputType = "code_execution_output"

// CodeExecutionResultType code_execution_result 子块的 type 字段类型。
type CodeExecutionResultType string

const CodeExecutionResultTypeCodeExecutionResult CodeExecutionResultType = "code_execution_result"

// EncryptedCodeExecutionResultType encrypted_code_execution_result 子块的 type 字段类型。
type EncryptedCodeExecutionResultType string

const EncryptedCodeExecutionResultTypeEncryptedCodeExecutionResult EncryptedCodeExecutionResultType = "encrypted_code_execution_result"

// CodeExecutionToolResultErrorType code_execution_tool_result_error 子块的 type 字段类型。
type CodeExecutionToolResultErrorType string

const CodeExecutionToolResultErrorTypeCodeExecutionToolResultError CodeExecutionToolResultErrorType = "code_execution_tool_result_error"

// BashCodeExecutionOutputType bash_code_execution_output 子块的 type 字段类型。
type BashCodeExecutionOutputType string

const BashCodeExecutionOutputTypeBashCodeExecutionOutput BashCodeExecutionOutputType = "bash_code_execution_output"

// BashCodeExecutionResultType bash_code_execution_result 子块的 type 字段类型。
type BashCodeExecutionResultType string

const BashCodeExecutionResultTypeBashCodeExecutionResult BashCodeExecutionResultType = "bash_code_execution_result"

// BashCodeExecutionToolResultErrorType bash_code_execution_tool_result_error 子块的 type 字段类型。
type BashCodeExecutionToolResultErrorType string

const BashCodeExecutionToolResultErrorTypeBashCodeExecutionToolResultError BashCodeExecutionToolResultErrorType = "bash_code_execution_tool_result_error"

// TextEditorCodeExecutionViewResultType text_editor_code_execution_view_result 子块的 type 字段类型。
type TextEditorCodeExecutionViewResultType string

const TextEditorCodeExecutionViewResultTypeTextEditorCodeExecutionViewResult TextEditorCodeExecutionViewResultType = "text_editor_code_execution_view_result"

// TextEditorCodeExecutionCreateResultType text_editor_code_execution_create_result 子块的 type 字段类型。
type TextEditorCodeExecutionCreateResultType string

const TextEditorCodeExecutionCreateResultTypeTextEditorCodeExecutionCreateResult TextEditorCodeExecutionCreateResultType = "text_editor_code_execution_create_result"

// TextEditorCodeExecutionStrReplaceResultType text_editor_code_execution_str_replace_result 子块的 type 字段类型。
type TextEditorCodeExecutionStrReplaceResultType string

const TextEditorCodeExecutionStrReplaceResultTypeTextEditorCodeExecutionStrReplaceResult TextEditorCodeExecutionStrReplaceResultType = "text_editor_code_execution_str_replace_result"

// TextEditorCodeExecutionToolResultErrorType text_editor_code_execution_tool_result_error 子块的 type 字段类型。
type TextEditorCodeExecutionToolResultErrorType string

const TextEditorCodeExecutionToolResultErrorTypeTextEditorCodeExecutionToolResultError TextEditorCodeExecutionToolResultErrorType = "text_editor_code_execution_tool_result_error"

// CodeExecutionToolResultErrorCode code_execution 错误码枚举（已存在；此处保留以明确其在子结构中的位置）。
//
// 对齐官方 CodeExecutionToolResultErrorCode 常量集合（message.go:1442-1458）。
type AnthropicCodeExecutionToolResultErrorCode string

const (
	AnthropicCodeExecutionToolResultErrorCodeInvalidToolInput AnthropicCodeExecutionToolResultErrorCode = "invalid_tool_input"
	AnthropicCodeExecutionToolResultErrorCodeUnavailable      AnthropicCodeExecutionToolResultErrorCode = "unavailable"
	AnthropicCodeExecutionToolResultErrorCodeContainerExpired AnthropicCodeExecutionToolResultErrorCode = "container_expired"
	AnthropicCodeExecutionToolResultErrorCodeTimeLimit        AnthropicCodeExecutionToolResultErrorCode = "time_limit_exceeded"
	AnthropicCodeExecutionToolResultErrorCodePermissionDenied AnthropicCodeExecutionToolResultErrorCode = "permission_denied"
)

// BashCodeExecutionToolResultErrorCode bash_code_execution 错误码枚举。
type AnthropicBashCodeExecutionToolResultErrorCode string

const (
	AnthropicBashCodeExecutionToolResultErrorCodeInvalidToolInput AnthropicBashCodeExecutionToolResultErrorCode = "invalid_tool_input"
	AnthropicBashCodeExecutionToolResultErrorCodeUnavailable      AnthropicBashCodeExecutionToolResultErrorCode = "unavailable"
	AnthropicBashCodeExecutionToolResultErrorCodeTimeLimit        AnthropicBashCodeExecutionToolResultErrorCode = "time_limit_exceeded"
)

// TextEditorCodeExecutionToolResultErrorCode text_editor_code_execution 错误码枚举。
type AnthropicTextEditorCodeExecutionToolResultErrorCode string

const (
	AnthropicTextEditorCodeExecutionToolResultErrorCodeInvalidToolInput AnthropicTextEditorCodeExecutionToolResultErrorCode = "invalid_tool_input"
	AnthropicTextEditorCodeExecutionToolResultErrorCodeUnavailable      AnthropicTextEditorCodeExecutionToolResultErrorCode = "unavailable"
	AnthropicTextEditorCodeExecutionToolResultErrorCodeTimeLimit        AnthropicTextEditorCodeExecutionToolResultErrorCode = "time_limit_exceeded"
	AnthropicTextEditorCodeExecutionToolResultErrorCodeFileNotFound     AnthropicTextEditorCodeExecutionToolResultErrorCode = "file_not_found"
)

// TextEditorCodeExecutionViewResultFileType text_editor 视图结果的文件类型。
type AnthropicTextEditorCodeExecutionViewResultFileType string

const (
	AnthropicTextEditorCodeExecutionViewResultFileTypeText  AnthropicTextEditorCodeExecutionViewResultFileType = "text"
	AnthropicTextEditorCodeExecutionViewResultFileTypeImage AnthropicTextEditorCodeExecutionViewResultFileType = "image"
	AnthropicTextEditorCodeExecutionViewResultFileTypePDF   AnthropicTextEditorCodeExecutionViewResultFileType = "pdf"
)

// ============== code_execution 子结构（3 个 + Error） ==============

// AnthropicCodeExecutionOutputBlock code_execution 单个输出文件引用。
//
// 对应官方 CodeExecutionOutputBlockParam（message.go:1048-1054）。
type AnthropicCodeExecutionOutputBlock struct {
	Type   CodeExecutionOutputType `json:"type"` // 固定 "code_execution_output"
	FileID string                  `json:"file_id"`
}

// AnthropicCodeExecutionResultBlock code_execution 成功结果（含 stdout/stderr/return_code/output 文件）。
//
// 对应官方 CodeExecutionResultBlockParam（message.go:1089-1098）。
type AnthropicCodeExecutionResultBlock struct {
	Type       CodeExecutionResultType             `json:"type"` // 固定 "code_execution_result"
	Content    []AnthropicCodeExecutionOutputBlock `json:"content,omitempty"`
	ReturnCode int64                               `json:"return_code"`
	Stderr     string                              `json:"stderr"`
	Stdout     string                              `json:"stdout"`
}

// AnthropicEncryptedCodeExecutionResultBlock code_execution 加密结果（stdout 被服务端加密）。
//
// 对应官方 EncryptedCodeExecutionResultBlockParam（message.go:3299-3308）。
type AnthropicEncryptedCodeExecutionResultBlock struct {
	Type            EncryptedCodeExecutionResultType    `json:"type"` // 固定 "encrypted_code_execution_result"
	Content         []AnthropicCodeExecutionOutputBlock `json:"content,omitempty"`
	EncryptedStdout string                              `json:"encrypted_stdout"`
	ReturnCode      int64                               `json:"return_code"`
	Stderr          string                              `json:"stderr"`
}

// AnthropicCodeExecutionToolResultError code_execution 错误。
//
// 对应官方 CodeExecutionToolResultErrorParam（message.go:1463-1471）。
type AnthropicCodeExecutionToolResultError struct {
	Type      CodeExecutionToolResultErrorType          `json:"type"` // 固定 "code_execution_tool_result_error"
	ErrorCode AnthropicCodeExecutionToolResultErrorCode `json:"error_code,omitempty"`
}

// ============== bash_code_execution 子结构（2 个 + Error） ==============

// AnthropicBashCodeExecutionOutputBlock bash 工具单个输出条目。
//
// 对应官方 BashCodeExecutionOutputBlockParam（message.go:212-228）。
type AnthropicBashCodeExecutionOutputBlock struct {
	Type   BashCodeExecutionOutputType `json:"type"` // 固定 "bash_code_execution_output"
	Output string                      `json:"output"`
}

// AnthropicBashCodeExecutionResultBlock bash 工具成功结果。
//
// 对应官方 BashCodeExecutionResultBlockParam（message.go:253-268）。
type AnthropicBashCodeExecutionResultBlock struct {
	Type       BashCodeExecutionResultType             `json:"type"` // 固定 "bash_code_execution_result"
	Content    []AnthropicBashCodeExecutionOutputBlock `json:"content,omitempty"`
	ReturnCode int64                                   `json:"return_code"`
	Stderr     string                                  `json:"stderr"`
	Stdout     string                                  `json:"stdout"`
}

// AnthropicBashCodeExecutionToolResultError bash 工具错误。
//
// 对应官方 BashCodeExecutionToolResultErrorParam（message.go:463-471）。
type AnthropicBashCodeExecutionToolResultError struct {
	Type      BashCodeExecutionToolResultErrorType          `json:"type"` // 固定 "bash_code_execution_tool_result_error"
	ErrorCode AnthropicBashCodeExecutionToolResultErrorCode `json:"error_code,omitempty"`
}

// ============== text_editor_code_execution 子结构（4 个 + Error） ==============

// AnthropicTextEditorCodeExecutionViewResultBlock text_editor view 命令结果。
//
// 对应官方 TextEditorCodeExecutionViewResultBlockParam（message.go:6438-6449）。
type AnthropicTextEditorCodeExecutionViewResultBlock struct {
	Type       TextEditorCodeExecutionViewResultType              `json:"type"` // 固定 "text_editor_code_execution_view_result"
	Content    string                                             `json:"content"`
	FileType   AnthropicTextEditorCodeExecutionViewResultFileType `json:"file_type,omitempty"`
	NumLines   *int64                                             `json:"num_lines,omitempty"`
	StartLine  *int64                                             `json:"start_line,omitempty"`
	TotalLines *int64                                             `json:"total_lines,omitempty"`
}

// AnthropicTextEditorCodeExecutionCreateResultBlock text_editor create 命令结果。
//
// 对应官方 TextEditorCodeExecutionCreateResultBlockParam（message.go:6017-6023）。
type AnthropicTextEditorCodeExecutionCreateResultBlock struct {
	Type         TextEditorCodeExecutionCreateResultType `json:"type"` // 固定 "text_editor_code_execution_create_result"
	IsFileUpdate bool                                    `json:"is_file_update"`
}

// AnthropicTextEditorCodeExecutionStrReplaceResultBlock text_editor str_replace 命令结果。
//
// 对应官方 TextEditorCodeExecutionStrReplaceResultBlockParam（message.go:6060-6070）。
type AnthropicTextEditorCodeExecutionStrReplaceResultBlock struct {
	Type     TextEditorCodeExecutionStrReplaceResultType `json:"type"` // 固定 "text_editor_code_execution_str_replace_result"
	NewLines *int64                                      `json:"new_lines,omitempty"`
	NewStart *int64                                      `json:"new_start,omitempty"`
	OldLines *int64                                      `json:"old_lines,omitempty"`
	OldStart *int64                                      `json:"old_start,omitempty"`
	Lines    []string                                    `json:"lines,omitempty"`
}

// AnthropicTextEditorCodeExecutionToolResultError text_editor 工具错误。
//
// 对应官方 TextEditorCodeExecutionToolResultErrorParam（message.go:6383-6392）。
type AnthropicTextEditorCodeExecutionToolResultError struct {
	Type         TextEditorCodeExecutionToolResultErrorType          `json:"type"` // 固定 "text_editor_code_execution_tool_result_error"
	ErrorCode    AnthropicTextEditorCodeExecutionToolResultErrorCode `json:"error_code,omitempty"`
	ErrorMessage *string                                             `json:"error_message,omitempty"`
}
