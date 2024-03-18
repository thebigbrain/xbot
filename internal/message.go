package internal

import (
	"time"

	"github.com/tyloafer/langchaingo/schema"
)

// Message 定义了SSE中发送的消息的结构
type Message struct {
	SessionID string                 `json:"sessionID"`
	MessageID string                 `json:"id"`
	Sender    schema.ChatMessageType `json:"sender"`
	Text      string                 `json:"text"`
	Timestamp time.Time              `json:"timestamp"`
}
