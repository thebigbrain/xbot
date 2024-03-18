package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 应用所有中间件到路由器
func ApplyMiddlewares(router *gin.Engine) {
	// 应用CORS中间件
	router.Use(corsMiddleware())

	// ...在这里添加其他中间件...
}

// CORS中间件的实现 (此代码为示例，并不完整)
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 在生产环境中应该是特定的域名而不是 '*'
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			// 如果是预检请求则直接返回200
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next() // 处理实际的请求
	}
}
