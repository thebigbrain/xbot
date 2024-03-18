package chat

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thebigbrain/xbot/internal/sse"
	"github.com/tyloafer/langchaingo/schema"
	// ... other imports ...
)

func SetupChatRoutes(router *gin.Engine, service *ChatService) {
	api := router.Group("/api")
	api.GET("/history", getHistoryHandler(service))
	api.POST("/message", postMessageHandler(service))
	api.POST("/send", sendHandler(service))
}

func getHistoryHandler(service *ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q Message
		c.ShouldBindQuery(&q)

		history, err := service.RetrieveHistory(c.Request.Context(), q.SessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve history"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"history": history})
	}
}

func postMessageHandler(_ *ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		clientGone := c.Writer.CloseNotify()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-clientGone:
				return
			case <-ticker.C:
				message := Message{
					// 这里应包含实际的消息数据
				}
				c.SSEvent("message", message)
				c.Writer.Flush()
			}
		}
	}
}

func sendHandler(service *ChatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input Message
		if err := c.BindQuery(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
			return
		}

		input.Sender = schema.ChatMessageTypeHuman
		dateNow := time.Now()

		ctx := c.Request.Context()
		// 保存消息
		savedMessage, err := service.SaveMessage(ctx, Message{
			Id:        uuid.NewString(),
			SessionID: input.SessionID,
			Sender:    input.Sender,
			Text:      input.Text,
			Timestamp: dateNow,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ts := time.Now()

		// 调用LLM处理用户的输入并获取回复
		responseMessage, err := service.ProcessMessageWithLLM(ctx, *savedMessage, func(chunkStr string) {
			sse.SendSse(c, Message{
				Id:        uuid.NewString(),
				SessionID: input.SessionID,
				Sender:    schema.ChatMessageTypeAI,
				Text:      chunkStr,
				Timestamp: ts,
			})
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// 发送LLM的回复作为SSE
		c.SSEvent("message", responseMessage)
		c.Writer.Flush()
	}
}

// ... 其他路由处理函数...
