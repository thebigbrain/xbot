package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thebigbrain/xbot/internal"
)

func main() {
	// 初始化数据库和路由器
	defer internal.GetDB().Close()
	router := gin.Default()

	// 应用中间件
	internal.ApplyMiddlewares(router)

	// 设置路由
	internal.SetupRoutes(router)

	// 开始监听HTTP
	router.Run(":5000")
}
