package frames

import (
	"crypto/rand"
	"encoding/hex"
)

// fallbackMessageID crypto/rand 失败时的兜底 ID。
//
// 实际触发概率极低（仅在系统熵源失效时发生），但下游 SDK 对 message_id 格式
// （"msg_<24位hex>"）强校验，所以兜底也必须保持合法。
const fallbackMessageID = "msg_000000000000000000000000"

// GenerateMessageID 生成 Anthropic 风格的 message id（"msg_<24 位 hex>"）。
//
// 用途：网关在 Anthropic 协议代理路径下需要对外写出 message_start 事件，
// 其中的 message.id 字段必须是 "msg_" 前缀 + 24 位 hex 的格式，
// 因为我们做的是协议转发（ID 由 Gateway 生成而非上游分配）。
//
// 实现：crypto/rand 生成 12 字节 → hex 编码为 24 字符串。失败时退回固定值。
func GenerateMessageID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return fallbackMessageID
	}
	return "msg_" + hex.EncodeToString(b)
}
