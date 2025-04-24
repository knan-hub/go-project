package controller

import (
	"go-project/common"
	"go-project/logger"
	"go-project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FetchContent(c *gin.Context) {
	var params struct {
		URL  string `json:"url" binding:"required"`
		Type string `json:"type"` // "http"或"rod"或者不传默认为"rod"
	}
	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}

	url := params.URL
	_type := params.Type

	// service.ThroughCache()
	result, err := service.FetchContent(url, _type)
	if err != nil {
		logger.Logger.Error("error", zap.Error(err))
		common.ResponseError(c, common.SERVER_BUSY)
		return
	}

	common.ResponseSuccess(c, result)
}

func FetchLinks(c *gin.Context) {
	var params struct {
		URL  string `json:"url" binding:"required"`
		Type string `json:"type"` // "http"或"rod"或者不传默认为"rod"
	}
	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}

	url := params.URL

	var result interface{}
	var err error

	// 默认使用rod
	if params.Type == "http" {
		result, err = service.FetchLinksWithHttp(url)
	} else {
		result, err = service.FetchLinksWithRod(url)
	}

	if err != nil {
		logger.Logger.Error("error", zap.Error(err))
		common.ResponseError(c, common.SERVER_BUSY)
		return
	}

	common.ResponseSuccess(c, result)
}
