package chat

import (
	"time"

	"github.com/tyloafer/langchaingo/schema"
)

// Message 定义了 API 消息的结构。
type Message struct {
	ID        string                 `json:"id"`
	SessionID string                 `json:"sessionId"`
	Sender    schema.ChatMessageType `json:"sender"`
	Text      string                 `json:"text"`
	Timestamp time.Time              `json:"timestamp"`
}

// 为其他模型定义留下扩展空间...
