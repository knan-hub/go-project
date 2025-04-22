package router

import (
	"go-project/logger"
	"go-project/middleware"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.DebugMode || mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery(true))
	r.Use(middleware.CORSMiddleware())

	// 返回纯文本格式的响应
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 返回JSON格式的响应
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})

	// /debug/pprof/ (性能分析首页)
	// /debug/pprof/heap (内存分析)
	// /debug/pprof/goroutine (协程分析)
	// /debug/pprof/profile (CPU分析)
	// /debug/pprof/trace (执行追踪)
	pprof.Register(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
