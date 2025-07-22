package controllers

import (
	"fmt"
	"main/global"
	"main/models"
	"main/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 添加基础食材种类
func AddFood(ctx *gin.Context) {
	foodname := ctx.PostForm("name")
	if foodname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Food name is required"},
		)
		return
	}

	if err := global.DB.Where("name = ?", foodname).First(&models.Food{}).Error; err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Food already exists",
		})
		return
	}

	description := ctx.PostForm("description")

	imageconfig := utils.ImageUploadConfig{
		Field:    "image",
		SavePath: "./images/food",
		BasicURL: "/image/food",
		Prefix:   "food",
	}

	imageurl, err := utils.UploadImage(ctx, imageconfig)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid image upload",
		})
		return
	}
	food := models.Food{
		Name:        foodname,
		Description: description,
		Image:       imageurl,
	}

	if err := global.DB.Create(&food).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add food",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Food added successfully",
	})
}

// 删除已有食材
func DeleteFoodByName(ctx *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var food models.Food
	if input.Name != "" {
		if err := global.DB.Where("name = ?", input.Name).First(&food).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Food not found",
			})
			return
		}
	}

	if food.Image != "" {
		removepath := filepath.Join("./images", food.Image[7:])
		fmt.Printf("Removing old avatar: %s\n", removepath)
		if err := os.Remove(removepath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to remove old avatar: %v", err),
			})
			return
		}
		fmt.Printf("Old avatar removed: %s\n", removepath)
	}

	if err := utils.DeleteFoodByID(int(food.ID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete food",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Food deleted successfully",
	})
}

// 获取所有食材
func GetAllFood(ctx *gin.Context) {
	var foods []models.Food
	if err := global.DB.Select("name,description,image").Find(&foods).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve foods",
		})
		return
	}

	if len(foods) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No foods found",
		})
		return
	}

	var foodResponses []gin.H
	for _, food := range foods {
		foodResponses = append(foodResponses, gin.H{
			"name":        food.Name,
			"description": food.Description,
			"image":       food.Image,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"foods": foodResponses,
	})
}

// 获取食材-Name
func GetFoodByName(ctx *gin.Context) {

	inputname := ctx.Param("name")

	if inputname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Food name is required",
		})
		return
	}

	var food models.Food
	if err := global.DB.Where("name LIKE ?", "%"+inputname+"%").First(&food).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Food not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"food": gin.H{
			"name":        food.Name,
			"description": food.Description,
			"image":       food.Image,
		}})

}

// 更新食材信息
func UpdateFoodByName(ctx *gin.Context) {
	foodname := ctx.PostForm("name")
	newfoodname := ctx.PostForm("newname")
	description := ctx.PostForm("description")

	// 检查食材是否存在
	var food models.Food
	if err := global.DB.Where("name = ?", foodname).First(&food).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Food not found",
		})
		return
	}

	foodupdate := make(map[string]interface{})

	if foodname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Food name is required",
		})
		return
	}

	if newfoodname != "" {
		foodupdate["name"] = newfoodname // 用于更新食材名称
	}
	if description != "" {
		foodupdate["description"] = description // 更新描述
	}

	imagecoinfig := utils.ImageUploadConfig{
		Field:    "image",
		SavePath: "./images/food",
		BasicURL: "/image/food",
		Prefix:   "food",
	}
	imageurl, err := utils.UploadImage(ctx, imagecoinfig)
	if err == nil {
		foodupdate["image"] = imageurl
	}

	if len(foodupdate) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields to update",
		})
		return
	}

	if err := global.DB.Model(&food).Updates(foodupdate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update food",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Food updated successfully",
	})
}
