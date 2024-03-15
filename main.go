package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	sqlite "database/sql"

	"github.com/gin-gonic/gin"
	"github.com/tyloafer/langchaingo/llms"
	"github.com/tyloafer/langchaingo/llms/ollama"
	"github.com/tyloafer/langchaingo/schema"
	_ "modernc.org/sqlite"

	"github.com/google/uuid"
)

// corsMiddleware 创建一个CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 生产环境应指定具体域名
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// Message 定义了SSE中发送的消息的结构
type Message struct {
	SessionID string    `json:"sessionID"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

// 初始化数据库并创建所需的表
func initDB() *sqlite.DB {
	conn, err := sqlite.Open("sqlite", "chat_history.db")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	sqlStmt := `
    CREATE TABLE IF NOT EXISTS messages (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        session_id TEXT NOT NULL,
        sender TEXT NOT NULL,
        text TEXT NOT NULL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := conn.Exec(sqlStmt, nil); err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return conn
}

var db *sqlite.DB

func GetDB() *sqlite.DB {
	if db == nil {
		db = initDB()
	}

	return db
}

// saveMessage 保存消息到数据库
func saveMessage(msg Message) error {
	stmt, err := db.Prepare("INSERT INTO messages(session_id, sender, text) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(msg.SessionID, msg.Sender, msg.Text)
	return err
}

// retrieveHistory 获取聊天历史
func retrieveHistory() ([]Message, error) {
	rows, err := db.Query("SELECT session_id, sender, text, timestamp FROM messages ORDER BY timestamp ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.SessionID, &msg.Sender, &msg.Text, &msg.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, msg)
	}

	return history, nil
}

func getChatHistories() (r []schema.ChatMessage) {
	messages, _ := retrieveHistory()
	for _, msg := range messages {
		if msg.Sender == "human" {
			r = append(r, schema.HumanChatMessage{Content: msg.Text})
		} else {
			r = append(r, schema.AIChatMessage{Content: msg.Text})
		}
	}
	return r
}

func sendSse(c *gin.Context, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", data))
	// c.Writer.WriteString(fmt.Sprintf("%s", data))
	c.Writer.Flush()
}

func main() {
	defer GetDB().Close()

	router := gin.Default()

	// 应用CORS中间件
	router.Use(corsMiddleware())

	api := router.Group("/api")

	// 路由以获取聊天历史
	api.GET("/history", func(c *gin.Context) {
		// 这里应该是连接数据库或其他存储，获取聊天历史记录
		// 为了演示，我们只是返回一个空数组
		history, err := retrieveHistory()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve history"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"history": history})
	})

	// 使用HTTP流式响应发送消息的路由
	api.POST("/message", func(c *gin.Context) {
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
	})

	// SSE消息发送的路由
	api.POST("/send", func(c *gin.Context) {
		// 设置SSE头部信息
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		var input Message
		c.BindJSON(&input)

		if input.SessionID == "" {
			input.SessionID = uuid.NewString()
		}

		input.Sender = "human"
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
				Sender:    "bot",
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
			Sender:    "bot", // 指定系统作为消息发送者
			Text:      fullMessage,
			Timestamp: time.Now(),
		}
		err = saveMessage(responseMessage)
		if err != nil {
			log.Printf("Error saving message to database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving message"})
			return
		}

	})

	// 开始监听HTTP
	router.Run(":5000")
}
