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

// 用户获取信息面板
func GetUserInfo(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	fmt.Println("Username from context:", userid)

	var user models.User
	if err := global.DB.Where("id = ?", userid).First(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user information",
		})
		return
	}

	outfile := struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
		Avatar   string `json:"avatar"`
	}{
		Username: user.Username,
		Name:     user.Name,
		Gender:   user.Gender,
		Phone:    user.Phone,
		Address:  user.Address,
		Avatar:   user.Avatar,
	}

	if outfile.Avatar == "" {
		outfile.Avatar = "/image/avatar/default/rice.jpg"
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User information retrieved successfully",
		"data":    outfile,
	})
}

// 更新用户信息
func UpdateUserInfo(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not Found",
		})
		return
	}

	type UserUpdate struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Gender   string `json:"gender"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}

	var userupdate UserUpdate
	if err := ctx.ShouldBindJSON(&userupdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	result := global.DB.Model(&models.User{}).Where("id = ?", userid).Updates(userupdate)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user information",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User information updated successfully",
	})
}

// 注销用户
func DeleteUser(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	// 检查用户是否存在
	var user models.User
	if err := global.DB.Where("id = ?", userid).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.Avatar != "" {
		// 删除用户头像文件
		removepath := filepath.Join("./images", user.Avatar[7:]) // 去掉前缀"/image/"
		if _, err := os.Stat(removepath); err == nil {
			if err := os.Remove(removepath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to remove user avatar: " + err.Error(),
				})
				return
			}
		} else if !os.IsNotExist(err) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error checking user avatar path: " + err.Error(),
			})
			return
		}
	}

	result := global.DB.Where("id = ?", userid).Unscoped().Delete(&models.User{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found or already deleted",
		})
		return
	}
	// 清除用户的JWT令牌或相关会话信息（如果有的话）

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// 更新用户头像
func UpdateUserAvatar(ctx *gin.Context) {
	imageconfig := utils.ImageUploadConfig{
		Field:    "avatar",
		SavePath: "./images/avatar/user",
		BasicURL: "/image/avatar/user",
		Prefix:   "avatar",
	}

	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	// 获取旧头像路径,查询数据库获取
	var user models.User
	if err := global.DB.Where("id = ?", userid).First(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user information",
		})
		return
	}

	avatarurl, err := utils.UploadImage(ctx, imageconfig, user.Avatar)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to upload image: " + err.Error(),
		})
		return
	}

	if err := global.DB.Model(&models.User{}).Where("id = ?", userid).Update("avatar", avatarurl).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user avatar",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User avatar updated successfully",
		"avatar":  avatarurl,
	})
}

// 收藏食谱
func CollectRecipe(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var input struct {
		RecipeID uint `json:"recipe_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	// 检查是否已收藏该食谱
	var existingCollection models.CollectedRecipes
	if err := global.DB.Where("user_id = ? AND recipe_id = ?", userid, input.RecipeID).First(&existingCollection).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Recipe already collected",
		})
		return
	}

	// 创建收藏记录
	collection := models.NewCollectedRecipes()

	if err := global.DB.Create(&collection).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to collect recipe",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe collected successfully",
	})
}

// 获取所有收藏的食谱
func GetCollectedRecipes(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var collections []models.CollectedRecipes
	if err := global.DB.Where("user_id = ?", userid).Find(&collections).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve collected recipes",
		})
		return
	}

	var userfrecipe []models.Recipe
	for _, collection := range collections {
		var recipe models.Recipe
		if err := global.DB.Where("id = ?", collection.RecipeID).First(&recipe).Error; err == nil {
			userfrecipe = append(userfrecipe, recipe)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Collected recipes retrieved successfully",
		"data":    userfrecipe,
	})
}

// 取消收藏食谱
func UncollectRecipe(ctx *gin.Context) {
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var input struct {
		RecipeID uint `json:"recipe_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	result := global.DB.Where("user_id = ? AND recipe_id = ?", userid, input.RecipeID).Delete(&models.CollectedRecipes{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to uncollect recipe",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Collection not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe uncollected successfully",
	})
}
