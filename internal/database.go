package internal

import (
	sqlite "database/sql"
	"log"

	"github.com/tyloafer/langchaingo/schema"
	_ "modernc.org/sqlite"
)

var db *sqlite.DB
var historyCache []Message

func initDB() *sqlite.DB {
	conn, err := sqlite.Open("sqlite", "chatbot.db")
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

	db = conn

	// 加载现有历史至内存缓存
	historyCache, err = retrieveHistoryFromDB()
	if err != nil {
		log.Fatal("failed to load chat history from database:", err)
	}

	return conn
}

func GetDB() *sqlite.DB {
	if db == nil {
		db = initDB()
	}

	return db
}

// ...其他数据库相关的函数，比如 retrieveHistoryFromDB、saveMessage 等...

// retrieveHistory 获取聊天历史
func retrieveHistoryFromDB() ([]Message, error) {
	rows, err := GetDB().Query("SELECT session_id, sender, text, timestamp FROM messages ORDER BY timestamp ASC")
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

// retrieveHistory 返回内存中的聊天历史
func retrieveHistory() ([]Message, error) {
	// 直接返回内存缓存的历史记录，而不是查询数据库
	return historyCache, nil
}

// saveMessage 保存消息到数据库
func saveMessage(msg Message) error {
	stmt, err := GetDB().Prepare("INSERT INTO messages(session_id, sender, text) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(msg.SessionID, msg.Sender, msg.Text)

	historyCache = append(historyCache, msg)
	return err
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
