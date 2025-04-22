package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许任何来源的请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 允许特定的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许特定的 HTTP 头部字段
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		// 如果有需要，可以设置其他跨域相关的头部信息

		// 放行请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
