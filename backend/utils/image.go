package utils

import (
	"fmt"

	"os"
	"path/filepath"
	"strings"
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

// 上传图片函数
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

	// 获取上传的文件

	form, err := ctx.MultipartForm()
	if err != nil {
		return "", fmt.Errorf("failed to parse form data: %v", err)
	}

	image, exists := form.File[config.Field]

	// 如果没有上传图片，结束函数，允许不上传图片

	if !exists || len(image) == 0 {
		return "", nil
	}

	// 如果有传入旧路径，则考虑后续处理
	if len(oldpath) > 0 && oldpath[0] != "" {
		// 检查旧路径是否存在，如果存在则删除
		oldpathlist := strings.Split(oldpath[0], ",")
		if len(oldpath) > 0 {
			for _, oldPathStr := range oldpathlist {
				err := DeleteImage(oldPathStr)
				if err != nil {
					return "", fmt.Errorf("failed to delete old image: %v", err)
				}
			}
		}
		// 若旧图片路径的文件名与新图片路径的文件名不同，则重命名
		oldfolder := filepath.Join("./images", filepath.Dir(oldpathlist[0])[7:]) // 去掉前缀"/image/"
		if oldfolder != config.SavePath {
			if err := os.Rename(oldfolder, config.SavePath); err != nil {
				return "", fmt.Errorf("failed to rename old folder: %v", err)
			}
		}
	}
	// 确保保存目录存在
	if err := os.MkdirAll(config.SavePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// 依次处理image中图片文件
	var outurls []string

	for index, img := range image {
		// 检查文件类型
		ext := filepath.Ext(img.Filename)
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
		if img.Size > config.MaxSize {
			return "", fmt.Errorf("file size %d exceeds limit %d", img.Size, config.MaxSize)
		}

		// 生成文件名
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("%s_%v_%d_%s%s", config.Prefix, userid, index, timestamp, ext)

		// 保存文件
		fullPath := filepath.Join(config.SavePath, filename)
		if err := ctx.SaveUploadedFile(img, fullPath); err != nil {
			return "", fmt.Errorf("failed to save file: %v", err)
		}
		outurls = append(outurls, fmt.Sprintf("%s/%s", config.BasicURL, filename))
	}

	// 把outurls转换为字符串
	imageurl := strings.Join(outurls, ",")

	// 返回访问URL
	return imageurl, nil
}

// 删除图片函数
func DeleteImage(imageurl string) error {
	// 检查图片URL是否为空
	if imageurl == "" {
		return fmt.Errorf("image URL is empty")
	}

	// 更换拼接为相对路径
	imagepath := filepath.Join("./images", imageurl[7:]) // 去掉前缀"/image/"

	// 检查文件是否存在。若不存在，不做报错，终止函数
	if _, err := os.Stat(imagepath); os.IsNotExist(err) {
		// 文件不存在时，记录日志但不返回错误
		fmt.Printf("Warning: Image file does not exist, skipping deletion: %s\n", imagepath)
		return nil // 不返回错误，认为删除成功
	}

	// 删除文件
	if err := os.Remove(imagepath); err != nil {
		return fmt.Errorf("failed to delete image file: %v", err)
	}
	// 成功删除
	return nil
}

// 删除图片目录函数
func DeleteImageDir(dir string) error {
	// 拼贴相对路径
	dirPath := filepath.Join("./images", dir[7:]) // 去掉前缀"/image/"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %v", err)
	}

	// 删除目录
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("failed to delete directory: %v", err)
	}
	// 成功删除
	return nil
}
