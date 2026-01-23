package controllers

import "github.com/gin-gonic/gin"

type ResponseData struct {
	// 统一响应结构：code/message/data
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseError 返回统一错误响应。
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(200, ResponseData{
		Code:    code,
		Message: code.Msg(),
	})
}

// ResponseErrorWithMsg 返回自定义错误文案。
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg string) {
	c.JSON(200, ResponseData{
		Code:    code,
		Message: msg,
	})
}

// ResponseSuccess 返回统一成功响应。
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	})
}
