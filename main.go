package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thebigbrain/xbot/internal/chat"
	"github.com/thebigbrain/xbot/internal/db"
	"github.com/thebigbrain/xbot/internal/middlewares"
)

func main() {
	// 初始化数据库和路由器
	defer db.GetDB().Close()
	r := gin.Default()

	// 应用中间件
	middlewares.ApplyMiddlewares(r)

	// 设置路由
	// internal.SetupRoutes(r)

	chatService := chat.NewChatService(db.GetDB())
	chat.SetupChatRoutes(r, chatService)

	// 启动服务器
	if err := r.Run(":5000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
