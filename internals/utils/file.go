package utils

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// maksimal ukuran file
const MaxFileSize = 500 * 1024

// tipe file yang di perbolehkan
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func UploadImageFile(ctx context.Context, file *multipart.FileHeader, uploadDir string, prefix string) (string, string, error) {
	if file == nil {
		return "", "", errors.New("file tidak ditemukan")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", "", errors.New("format file tidak didukung (hanya jpg, jpeg, png)")
	}

	if file.Size > MaxFileSize {
		return "", "", errors.New("ukuran file maksimal 500 kb")
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), prefix, ext)
	savePath := filepath.Join(uploadDir, filename)

	return savePath, filename, nil
}

