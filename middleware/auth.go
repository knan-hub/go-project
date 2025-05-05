package middleware

import (
	"go-project/api"
	"go-project/consts"
	"go-project/setting"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InternalAuthenticated() gin.HandlerFunc {
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

func Authenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, "Invalid authorization")
			ctx.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusUnauthorized, "Invalid authorization")
			ctx.Abort()
			return
		}

		valid, user := api.ValidateUserToken(parts[1])
		if !valid {
			ctx.JSON(http.StatusUnauthorized, "Unauthorization")
			ctx.Abort()
			return
		}

		// 将当前请求的userId信息保存到请求的上下文上
		// 后续的处理请求的函数中 可以用过c.Get(CtxUserIDKey) 来获取当前请求的用户信息
		ctx.Set(consts.CtxUserIDKey, user.UserId)

		ctx.Next()
	}
}
