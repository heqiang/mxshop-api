package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 错误相应
func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseErrorWithMsg 带信息的错误响应
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// RsponseSuccess 成功响应
func RsponseSuccess(c *gin.Context, code ResCode, data interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
