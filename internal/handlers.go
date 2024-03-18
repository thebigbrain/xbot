package internal

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tyloafer/langchaingo/llms"
	"github.com/tyloafer/langchaingo/llms/ollama"
	"github.com/tyloafer/langchaingo/schema"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/history", getHistoryHandler)
	api.POST("/message", postMessageHandler)
	api.POST("/send", sendHandler)
}

func getHistoryHandler(c *gin.Context) {
	// 这里应该是连接数据库或其他存储，获取聊天历史记录
	// 为了演示，我们只是返回一个空数组
	history, err := retrieveHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

func postMessageHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 模拟实际消息，这里只是简单地每隔一秒发送一条消息
	for {
		time.Sleep(1 * time.Second)
		// 这将是实际的消息数据结构体
		message := map[string]string{"message": "New live message", "user": "Server"}
		c.SSEvent("message", message)
		c.Writer.Flush()
	}
}

func sendHandler(c *gin.Context) {
	// 设置SSE头部信息
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	var input Message
	c.BindJSON(&input)

	if input.SessionID == "" {
		input.SessionID = uuid.NewString()
	}

	input.Sender = schema.ChatMessageTypeHuman
	saveMessage(input)

	// sendSse(c, input)

	// 初始化机器人回复的字符串缓冲区
	var responseBuffer strings.Builder

	llm, err := ollama.NewChat(ollama.WithLLMOptions(ollama.WithModel("codellama")))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	historyChats := getChatHistories()
	historyChats = append(historyChats, schema.HumanChatMessage{Content: input.Text})

	// 使用流处理功能调用LLM
	_, err = llm.Call(ctx, historyChats, llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
		chunkStr := string(chunk)
		responseBuffer.WriteString(chunkStr)

		// 创建并发送流消息
		msg := Message{
			SessionID: input.SessionID,
			Sender:    schema.ChatMessageTypeAI,
			Text:      chunkStr,
			Timestamp: time.Now(),
		}
		sendSse(c, msg)
		return nil
	}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing your message"})
		return
	}

	// 流式回复结束，保存整个会话到数据库
	fullMessage := responseBuffer.String()
	responseMessage := Message{
		SessionID: input.SessionID,
		Sender:    schema.ChatMessageTypeAI,
		Text:      fullMessage,
		Timestamp: time.Now(),
	}
	err = saveMessage(responseMessage)
	if err != nil {
		log.Printf("Error saving message to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving message"})
		return
	}
}
