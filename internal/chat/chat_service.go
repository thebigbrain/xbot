package chat

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/tyloafer/langchaingo/llms"
	"github.com/tyloafer/langchaingo/llms/ollama"
	"github.com/tyloafer/langchaingo/schema"
	// ... other imports...
)

type ChatService struct {
	db          *sql.DB
	historyLock sync.RWMutex // 保护以下的消息历史缓存
	historyMap  map[string][]Message
}

func NewChatService(db *sql.DB) *ChatService {
	return &ChatService{
		db:         db,
		historyMap: make(map[string][]Message),
	}
}

// 内存中缓存历史消息
func (cs *ChatService) cacheHistory(sessionID string, messages ...Message) {
	cs.historyLock.Lock()
	defer cs.historyLock.Unlock()
	if _, ok := cs.historyMap[sessionID]; !ok {
		cs.historyMap[sessionID] = []Message{}
	}
	cs.historyMap[sessionID] = append(cs.historyMap[sessionID], messages...)
}

// 从内存缓存中检索历史消息
func (cs *ChatService) getHistoryFromCache(sessionID string) ([]Message, bool) {
	cs.historyLock.RLock()
	defer cs.historyLock.RUnlock()
	history, ok := cs.historyMap[sessionID]
	return history, ok
}

// RetrieveHistory 用于检索聊天历史记录。
func (cs *ChatService) RetrieveHistory(ctx context.Context, sessionID string) ([]Message, error) {
	// 尝试从缓存中检索历史，如果找到了直接返回
	if history, ok := cs.getHistoryFromCache(sessionID); ok {
		return history, nil
	}

	messages := []Message{}

	// 假设有一个合适的查询来检索消息历史
	rows, err := cs.db.QueryContext(ctx, `
			SELECT session_id, sender, text, timestamp FROM messages WHERE session_id = ?
			ORDER BY timestamp ASC
	`, sessionID)

	if err != nil {
		// 实际中可能添加日志记录等
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.SessionID, &msg.Sender, &msg.Text, &msg.Timestamp); err != nil {
			// 实际中可能添加日志记录等
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		// 实际中可能添加日志记录等
		return nil, err
	}

	// 存储消息到缓存供下次使用
	cs.cacheHistory(sessionID, messages...)

	return messages, nil
}

// SaveMessage 用于保存消息到数据库。
func (cs *ChatService) SaveMessage(ctx context.Context, msg Message) (*Message, error) {
	_, err := cs.db.ExecContext(ctx, `
			INSERT INTO messages (session_id, sender, text, timestamp) VALUES (?, ?, ?, ?)
	`, msg.SessionID, msg.Sender, msg.Text, msg.Timestamp)

	if err != nil {
		// 实际中可能添加日志记录等
		return nil, err
	}

	// 写入数据库成功后添加到内存缓存
	cs.cacheHistory(msg.SessionID, msg)

	return &msg, nil
}

// ...其他服务方法...
func (cs *ChatService) getChatHistories(ctx context.Context, sessionID string) (r []schema.ChatMessage) {
	messages, _ := cs.RetrieveHistory(ctx, sessionID)
	for _, msg := range messages {
		if msg.Sender == "human" {
			r = append(r, schema.HumanChatMessage{Content: msg.Text})
		} else {
			r = append(r, schema.AIChatMessage{Content: msg.Text})
		}
	}
	return r
}

// ProcessMessageWithLLM 调用LLM以处理用户消息并获取回复
func (cs *ChatService) ProcessMessageWithLLM(ctx context.Context, msg Message, handleStream func(chunkStr string)) (*Message, error) {
	// 初始化LLM客户端。实际实现可能需要配置参数如API keys等
	// 以下为模拟逻辑
	llm, err := ollama.NewChat(ollama.WithLLMOptions(ollama.WithModel("codellama")))
	if err != nil {
		return nil, err // 实际中可能要记录错误或更详细处理
	}

	// 根据业务需求收集聊天历史，以供LLM生成回复
	// 此处代码略，但通常会包括从数据库中检索相关会话
	historyChats := cs.getChatHistories(ctx, msg.SessionID)

	historyChats = append(historyChats, schema.HumanChatMessage{Content: msg.Text})

	var responseBuffer strings.Builder
	var responseMsg Message

	// 调用LLM
	_, err = llm.Call(ctx, historyChats, llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
		chunkStr := string(chunk)
		handleStream(chunkStr)
		responseBuffer.WriteString(chunkStr)
		return nil
	}))

	if err != nil {
		return nil, err // 实际中可能要记录错误或更详细处理
	}

	responseMsg = Message{
		SessionID: msg.SessionID,
		Sender:    schema.ChatMessageTypeAI,
		Text:      responseBuffer.String(),
		Timestamp: time.Now(),
	}

	// 将LLM的回复保存至数据库，实际逻辑略
	_, saveErr := cs.SaveMessage(ctx, responseMsg)
	if saveErr != nil {
		return nil, saveErr // 实际中可能要记录错误或更详细处理
	}

	// 返回LLM的回复消息
	return &responseMsg, nil
}
