package controllers

import (
	"encoding/json"

	"main/global"
	"main/models"
	"main/utils"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 用户增加食谱,使用form-data格式来传递数据，一个键值为json-json，另一个是来传送图片信息
func AddRecipe(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	bindjson := form.Value["json"][0]
	recipeInput := models.NewRecipe()
	if err := json.Unmarshal([]byte(bindjson), &recipeInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// 检查是否存在该食谱
	var existingRecipe models.Recipe
	if err := global.DB.Select("id").Where("name = ?", recipeInput.Name).First(&existingRecipe).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Recipe with this name already exists"})
		return
	}
	// 接受一张或者多张图片
	imageconfig := utils.ImageUploadConfig{
		Field:    "image",
		SavePath: "./images/recipe" + "/" + recipeInput.Name,
		BasicURL: "/image/recipe" + "/" + recipeInput.Name,
		Prefix:   "recipe",
	}

	imageurls, err := utils.UploadImage(ctx, imageconfig)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image upload: " + err.Error()})
		return
	}

	if imageurls != "" {
		recipeInput.Images = imageurls
	}

	// 使用recipeinput在数据库中创建食谱
	recipeInput.AuthorID = uint(userid.(uint8)) // 将userid转换为uint类型

	if err := global.DB.Create(&recipeInput).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe created successfully",
		"recipe": gin.H{
			"id":              recipeInput.ID,
			"name":            recipeInput.Name,
			"author_id":       recipeInput.AuthorID,
			"description":     recipeInput.Description,
			"images":          recipeInput.Images,
			"food_id":         recipeInput.FoodID,
			"cook_time":       recipeInput.CookTime,
			"process":         recipeInput.Process,
			"likes":           recipeInput.Likes,
			"comment_allowed": recipeInput.CommentAllowed,
		}})
}

// 用户更新自己的食谱
func UpdateRecipe(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	// 从 URL 参数获取食谱 ID
	recipeID := ctx.Param("id")
	if recipeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Recipe ID is required"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	bindjson := form.Value["json"][0]
	recipeInput := models.NewRecipe() // 传入的json数据
	if err := json.Unmarshal([]byte(bindjson), &recipeInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// 检查食谱是否存在且属于当前用户
	var existingRecipe models.Recipe
	usertype, exist := ctx.Get("usertype")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User type not found"})
		return
	}

	switch usertype.(uint8) {
	case 1: // 普通用户
		if err := global.DB.Where("id = ? AND author_id = ?", recipeID, userid).First(&existingRecipe).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found or does not belong to the user"})
			return
		}
	case 0: // 管理员
		if err := global.DB.Where("id = ?", recipeID).First(&existingRecipe).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
			return
		}
	}

	// 检查图片更新
	imageurls := existingRecipe.Images
	var newimageurls string
	// 如果表单中有图片字段，则进行图片上传
	if newimage, exists := form.File["image"]; exists && len(newimage) > 0 {
		imageconfig := utils.ImageUploadConfig{
			Field:    "image",
			SavePath: "./images/recipe" + "/" + recipeInput.Name,
			BasicURL: "/image/recipe" + "/" + recipeInput.Name,
			Prefix:   "recipe",
		}
		var err error
		newimageurls, err = utils.UploadImage(ctx, imageconfig, imageurls)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image upload: " + err.Error()})
			return
		}
	}

	recipeInput.Images = newimageurls

	// 更新食谱信息
	if err := global.DB.Model(&existingRecipe).Updates(recipeInput).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe updated successfully",
	})
}

// 用户删除自己的食谱
func DeleteRecipe(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// 从json中获取食谱ID
	recipeID := ctx.Query("id")
	if recipeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Recipe ID is required"})
		return
	}

	// 检查食谱是否存在且属于当前用户
	var existingRecipe models.Recipe
	usertype, exist := ctx.Get("usertype")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User type not found"})
		return
	}

	switch usertype.(uint8) {
	case 1: // 普通用户
		if err := global.DB.Where("id = ? AND author_id = ?", recipeID, userid).First(&existingRecipe).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found or does not belong to the user"})
			return
		}
	case 0: // 管理员
		if err := global.DB.Where("id = ?", recipeID).First(&existingRecipe).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
			return
		}
	}

	// 检查是否有图片需要删除
	if existingRecipe.Images != "" {
		imagelist := strings.Split(existingRecipe.Images, ",")
		for _, image := range imagelist {
			if err := utils.DeleteImage(image); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image: " + err.Error()})
				return
			}
		}
		// 删除图片目录
		err := utils.DeleteImageDir(filepath.Dir(imagelist[0]))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image directory: " + err.Error()})
			return
		}
	}

	// 删除食谱
	if err := global.DB.Unscoped().Delete(&existingRecipe).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe: " + err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe deleted successfully",
	})
}
