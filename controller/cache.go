package controller

import (
	"fmt"
	"go-project/common"
	"go-project/redis"

	"github.com/gin-gonic/gin"
)

func Set(c *gin.Context) {
	var params struct {
		Key   string `json:"key" validate:"required"`
		Value string `json:"value"`
	}

	if err := common.BindAndValidate(c, &params); err != nil {
		return
	}
	fmt.Println("123")
	err := redis.Set(c, params.Key, params.Value)
	if err != nil {
		common.ResponseError(c, common.REDIS_SET_ERROR)
		return
	}

	common.ResponseSuccess(c)
}

func Get(c *gin.Context) {
	var key = c.Query("key")
	result, err := redis.Get(c, key)
	if err != nil {
		common.ResponseError(c, common.REDIS_GET_ERROR)
		return
	}
	common.ResponseSuccess(c, result)
}
