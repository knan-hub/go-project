package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code Code        `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseSuccess(c *gin.Context, data ...interface{}) {
	var respData interface{}
	if len(data) == 0 {
		respData = nil
	} else if len(data) == 1 {
		respData = data[0]
	} else {
		respData = data
	}
	c.JSON(http.StatusOK, &Response{
		Code: SUCCESS,
		Msg:  SUCCESS.Msg(),
		Data: respData,
	})
	c.Abort()
}

func ResponseError(c *gin.Context, code Code) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
	c.Abort()
}

func ResponseErrorWithMsg(c *gin.Context, code Code, msg interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
	c.Abort()
}
