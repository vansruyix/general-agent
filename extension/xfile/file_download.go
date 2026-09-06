package xfile

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
)

func Download(ctx *gin.Context, filePath, fileName string) {
	if fileName == "" {
		fileName = filepath.Base(filePath)
	}
	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/octet-stream")

	// 发送文件
	ctx.File(filePath)
}

func DownloadZIP(ctx *gin.Context, filePath, fileName string) {
	if fileName == "" {
		fileName = filepath.Base(filePath)
	}
	// 设置响应头
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/zip")

	// 发送文件
	ctx.File(filePath)
}
