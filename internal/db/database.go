package db

import (
	sqlite "database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sqlite.DB

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

	return conn
}

func GetDB() *sqlite.DB {
	if db == nil {
		db = initDB()
	}

	return db
}
