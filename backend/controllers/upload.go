package controllers

import (
	"aithink/pkg/snowflake"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// uploadDir 上传文件存储目录
const uploadDir = "./uploads"

// UploadImageHandler 处理图片上传：保存到本地并返回访问路径。
func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, "file is required")
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".png"
	}
	filename := fmt.Sprintf("%d%s", snowflake.GenIDByInt(), ext)
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	c.JSON(http.StatusOK, ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data: gin.H{
			"url": "/uploads/" + filename,
		},
	})
}
