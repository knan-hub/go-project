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
	api.Use(middleware.InternalAuthenticated())
	{
		api.POST("/cache", controller.Set)
		api.GET("/cache", controller.Get)
	}

	api.Use(middleware.Authenticated())
	{
		api.POST("/vector/uploadSite", controller.UploadSite)
		api.POST("/vector/searchText", controller.SearchText)

		api.POST("/knowledgeBase/uploadFile", controller.UploadKnowledgeBaseFile)
		api.POST("/knowledgeBase/uploadSite", controller.UploadKnowledgeBaseSite)

		api.DELETE("/knowledgeBase/:knowledgeBaseId", controller.DeleteKnowledgeBase)
		api.DELETE("/knowledgeBaseItem/:knowledgeBaseItemId", controller.DeleteKnowledgeBaseItem)

		api.POST("/knowledgeBase/searchText", controller.SearchKnowledgeBase)
		api.GET("/knowledgeBase/text", controller.GetTextByDocumentSetId)
	}

	api.Use(middleware.InternalAuthenticated())
	{
		// 从给定的URL中提取内容
		api.POST("/crawlers/fetchContent", controller.FetchContent)
		// 从给定的URL中提取链接
		api.POST("/crawlers/fetchLinks", controller.FetchLinks)
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
