package controller

import (
	"go-project/common"
	"go-project/consts"
	"go-project/service"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func UploadKnowledgeBaseFile(c *gin.Context) {
	knowledgeBaseId := c.PostForm("knowledgeBaseId")
	knowledgeBaseItemId := c.PostForm("knowledgeBaseItemId")
	filename := c.PostForm("filename")

	if knowledgeBaseId == "" || knowledgeBaseItemId == "" || filename == "" {
		common.ResponseErrorWithMsg(c, common.INVALID_PARAMS, "knowledgeBaseId or knowledgeBaseItemId or filename is empty")
		return
	}

	knowledgeBaseIndex := service.KnowledgeBaseIndex{KnowledgeBaseId: knowledgeBaseId, KnowledgeBaseItemId: knowledgeBaseItemId, Filename: filename}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		common.ResponseErrorWithMsg(c, common.INVALID_PARAMS, err)
		return
	}

	if file == nil {
		common.ResponseErrorWithMsg(c, common.INVALID_PARAMS, "file is empty")
		return
	}

	fullpath := path.Join(consts.KNOWLEDGE_BASE_TMP_PATH, file.Filename)
	err = c.SaveUploadedFile(file, fullpath)
	if err != nil {
		common.ResponseErrorWithMsg(c, common.INVALID_PARAMS, err)
		return
	}

	res, err := service.NewKnowledgeBase().UploadFile(c, fullpath, knowledgeBaseIndex)
	if err != nil {
		common.ResponseErrorWithMsg(c, common.SERVER_BUSY, err)
		return
	}

	_ = os.Remove(fullpath)

	common.ResponseSuccess(c, res)
}

func UploadKnowledgeBaseSite(c *gin.Context) {
	var params service.KnowledgeSiteUploadParams
	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}

	repo := service.NewGitRepo(params.GitUrl)
	err := repo.ProcessKnowledgeBase(c, params)
	if err != nil {
		common.ResponseErrorWithMsg(c, common.SERVER_BUSY, err)
		return
	}

	common.ResponseSuccess(c)
}

func DeleteKnowledgeBase(c *gin.Context) {
	err := service.NewKnowledgeBase().DeleteFileByKnowledgeBaseId(c, c.Param("knowledgeBaseId"))
	if err != nil {
		common.ResponseErrorWithMsg(c, common.SERVER_BUSY, err)
		return
	}
	common.ResponseSuccess(c)
}

func DeleteKnowledgeBaseItem(c *gin.Context) {
	err := service.NewKnowledgeBase().DeleteFileByKnowledgeBaseItemId(c, c.Param("knowledgeBaseItemId"))
	if err != nil {
		common.ResponseErrorWithMsg(c, common.SERVER_BUSY, err)
		return
	}
	common.ResponseSuccess(c)
}

func SearchKnowledgeBase(c *gin.Context) {
	var params service.KnowledgeSearchParams
	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}

	results, err := service.NewKnowledgeBase().Search(c, params)
	if err != nil {
		common.ResponseErrorWithMsg(c, common.SERVER_BUSY, err)
		return
	}

	common.ResponseSuccess(c, results)
}

func GetTextByDocumentSetId(c *gin.Context) {
	result, err := service.NewKnowledgeBase().Coll.GetDocumentSetById(c, c.Query("documentSetId"))
	if err != nil {
		common.ResponseErrorWithMsg(c, common.SERVER_BUSY, err)
		return
	}
	common.ResponseSuccess(c, result)
}
