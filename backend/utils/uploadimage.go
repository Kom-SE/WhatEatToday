package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type ImageUploadConfig struct {
	Field       string   // 表单字段名
	MaxSize     int64    // 最大文件大小（字节）
	AllowedExts []string // 允许的文件扩展名
	SavePath    string   // 保存路径
	BasicURL    string   // 基础访问URL
	Prefix      string   // 文件名前缀
}

func UploadImage(ctx *gin.Context, config ImageUploadConfig, oldpath ...string) (string, error) {

	// 设置默认值
	if config.MaxSize == 0 {
		config.MaxSize = 5 * 1024 * 1024 // 默认5MB
	}
	if len(config.AllowedExts) == 0 {
		config.AllowedExts = []string{".jpg", ".jpeg", ".png", ".gif"}
	}
	if config.SavePath == "" {
		config.SavePath = "./images/"
	}
	// 获取用户ID
	userid, exists := ctx.Get("userid")
	if !exists {
		return "", fmt.Errorf("user not authenticated")
	}
	// 处理 oldpath 参数
	var oldPathStr string
	if len(oldpath) > 0 {
		oldPathStr = oldpath[0]
	}
	// 检查旧路径是否存在，如果存在则删除
	if oldPathStr != "" {
		removepath := filepath.Join("./images", oldPathStr[7:]) // 去掉前缀"/image/"
		if _, err := os.Stat(removepath); err == nil {
			fmt.Printf("Removing old image: %s\n", removepath)
			if err := os.Remove(removepath); err != nil {
				return "", fmt.Errorf("failed to remove old image: %v", err)
			}
			fmt.Printf("Old image removed: %s\n", removepath)
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("error checking old image path: %v", err)
		}
	}

	// 获取上传的文件
	image, err := ctx.FormFile(config.Field)
	if err != nil {
		return "", fmt.Errorf("failed to get image from form data: %v", err)
	}
	// 检查文件类型
	ext := filepath.Ext(image.Filename)
	allowed := false
	for _, allowedExt := range config.AllowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("unsupported file type %s. Allowed: %v", ext, config.AllowedExts)
	}

	// 检查文件大小
	if image.Size > config.MaxSize {
		return "", fmt.Errorf("file size %d exceeds limit %d", image.Size, config.MaxSize)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%v_%s%s", config.Prefix, userid, timestamp, ext)

	// 确保保存目录存在
	if err := os.MkdirAll(config.SavePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 保存文件
	fullPath := filepath.Join(config.SavePath, filename)
	if err := ctx.SaveUploadedFile(image, fullPath); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// 返回访问URL
	return fmt.Sprintf("%s/%s", config.BasicURL, filename), nil
}
