package gemini

import (
	"encoding/base64"
	"encoding/json"
)

// GeminiBlob inline 二进制内容（base64 自动转换）。
//
// wire-format：Data 序列化为 base64 字符串（Google 协议约定）。
// 反序列化时自动 base64 解码回 []byte。
type GeminiBlob struct {
	MIMEType string `json:"mimeType"`
	Data     []byte `json:"data"`
}

// MarshalJSON 实现：Data → base64 字符串。
func (b GeminiBlob) MarshalJSON() ([]byte, error) {
	type blobAlias struct {
		MIMEType string `json:"mimeType"`
		Data     string `json:"data"`
	}
	return json.Marshal(blobAlias{MIMEType: b.MIMEType, Data: base64.StdEncoding.EncodeToString(b.Data)})
}

// UnmarshalJSON 实现：base64 字符串 → []byte。
func (b *GeminiBlob) UnmarshalJSON(data []byte) error {
	type blobAlias struct {
		MIMEType string `json:"mimeType"`
		Data     string `json:"data"`
	}
	var aux blobAlias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	raw, err := base64.StdEncoding.DecodeString(aux.Data)
	if err != nil {
		return err
	}
	b.MIMEType = aux.MIMEType
	b.Data = raw
	return nil
}
