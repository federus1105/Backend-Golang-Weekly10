package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func SaveUploadedImage(ctx *gin.Context, file *multipart.FileHeader, prefix string, userID int) (string, error) {
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d_%d%s", prefix, time.Now().UnixNano(), userID, ext)
	location := filepath.Join("public", filename)

	if err := ctx.SaveUploadedFile(file, location); err != nil {
		return "", err
	}
	return filename, nil
}
