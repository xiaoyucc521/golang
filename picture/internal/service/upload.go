package service

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type upload struct{}

func NewUpload() *upload {
	return new(upload)
}

func (u *upload) Upload(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")

	// 获取文件后缀
	extName := path.Ext(file.Filename)
	// 支持上传的文件类型
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
	}

	if _, ok := allowExtMap[extName]; !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "文件类型不合法",
		})
		return
	}

	filename, _ := buildSaveName("image", extName)

	dir := "./" + "upload/" + filename

	_ = ctx.SaveUploadedFile(file, dir)

	ctx.JSON(http.StatusOK, gin.H{
		"filename": file.Filename,
		"size":     file.Size,
		"header":   file.Header.Get("Content-Type"),
	})
}

// 生成文件名
func buildSaveName(fileType, extName string) (string, error) {
	currentTime := time.Now()
	filePath := fmt.Sprintf("%v_%v%v", fileType, currentTime.Format("200601021504"), extName)
	return filePath, nil
}
