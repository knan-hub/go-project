package controller

import (
	"fmt"
	"go-project/common"
	"go-project/logger"
	"go-project/service"

	"github.com/gin-gonic/gin"
)

func UploadSite(c *gin.Context) {
	var params service.GitRepoParams
	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}
	if err := service.NewGitRepo(params.GitUrl).ProcessSite(c, params.LastCommitId, params.NewCommitId); err != nil {
		logger.Logger.Error(fmt.Sprintf("upload site error: %v", err))
		common.ResponseError(c, common.SERVER_BUSY)
		return
	}
	common.ResponseSuccess(c)
}

func SearchText(c *gin.Context) {
	// 从body里面提取参数，允许为空
	var params service.SiteSearchParams
	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}

	sitePage := service.NewSitePage()
	results, err := sitePage.Search(c, params)
	if err != nil {
		common.ResponseError(c, common.SERVER_BUSY)
		return
	}
	common.ResponseSuccess(c, results)

}
