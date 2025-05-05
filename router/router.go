package router

import (
	"go-project/controller"
	"go-project/middleware"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init(mode string) *gin.Engine {
	if mode == gin.DebugMode || mode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.GinLogger())
	r.Use(middleware.GinRecovery(true))
	r.Use(middleware.CORS())

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

	api := r.Group("/v0")

	// 内部认证路由组
	internal := api.Group("")
	internal.Use(middleware.InternalAuthenticated())
	{
		internal.POST("/cache", controller.Set)
		internal.GET("/cache", controller.Get)
		internal.POST("/crawlers/fetchContent", controller.FetchContent)
		internal.POST("/crawlers/fetchLinks", controller.FetchLinks)
	}

	// 普通认证路由组
	auth := api.Group("")
	auth.Use(middleware.Authenticated())
	{
		auth.POST("/vector/uploadSite", controller.UploadSite)
		auth.POST("/vector/searchText", controller.SearchText)

		auth.POST("/knowledgeBase/uploadFile", controller.UploadKnowledgeBaseFile)
		auth.POST("/knowledgeBase/uploadSite", controller.UploadKnowledgeBaseSite)

		auth.DELETE("/knowledgeBase/:knowledgeBaseId", controller.DeleteKnowledgeBase)
		auth.DELETE("/knowledgeBaseItem/:knowledgeBaseItemId", controller.DeleteKnowledgeBaseItem)

		auth.POST("/knowledgeBase/searchText", controller.SearchKnowledgeBase)
		auth.GET("/knowledgeBase/text", controller.GetTextByDocumentSetId)
	}

	// /debug/pprof/ (性能分析首页)
	// /debug/pprof/heap (内存分析)
	// /debug/pprof/goroutine (协程分析)
	// /debug/pprof/profile (CPU分析)
	// /debug/pprof/trace (执行追踪)
	pprof.Register(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})

	return r
}
