package controller

import (
	"go-project/common"
	"go-project/model"
	"go-project/mysql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateServerNodeHandler(c *gin.Context) {
	params := new(model.ParamPostServerNode)
	if err := common.BindAndValidate(c, params); err != nil {
		return
	}

	var node = model.ServerNode{
		InternalIp: params.Url,
	}

	if err := mysql.DB.Create(&node).Error; err != nil {
		common.ResponseError(c, common.DATABASE_ERROR)
		return
	}

	common.ResponseSuccess(c, gin.H{
		"id": node.ID,
	})
}

func GetServerNodeHandler(c *gin.Context) {
	var nodes []model.ServerNode

	if err := mysql.DB.Limit(20).Find(&nodes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			common.ResponseSuccess(c, []model.ServerNode{})
			return
		}
		common.ResponseError(c, common.DATABASE_ERROR)
		return
	}

	common.ResponseSuccess(c, nodes)
}
