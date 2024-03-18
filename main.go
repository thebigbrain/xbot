package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thebigbrain/xbot/internal"
	"github.com/thebigbrain/xbot/internal/chat"
)

func main() {
	// 初始化数据库和路由器
	defer internal.GetDB().Close()
	r := gin.Default()

	// 应用中间件
	internal.ApplyMiddlewares(r)

	// 设置路由
	// internal.SetupRoutes(r)

	chatService := chat.NewChatService(internal.GetDB())
	chat.SetupChatRoutes(r, chatService)

	// 启动服务器
	if err := r.Run(":5000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}

}
