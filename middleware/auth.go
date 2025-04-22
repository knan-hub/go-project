package middleware

import (
	"go-project/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求头中的x-api-key
		apiKey := ctx.GetHeader("x-api-key")

		// 校验x-api-key是否正确
		if apiKey != setting.Config.Self.INTERNAL_API_KEY {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 继续处理请求
		ctx.Next()
	}
}
