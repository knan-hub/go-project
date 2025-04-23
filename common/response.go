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

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: SUCCESS,
		Msg:  SUCCESS.Msg(),
		Data: data,
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
