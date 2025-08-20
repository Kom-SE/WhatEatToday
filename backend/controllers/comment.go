package controllers

import (
	"errors"
	"main/global"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 用户发布评论
func AddComment(ctx *gin.Context) {
	// 获取当前用户ID
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 定义评论输入结构体
	type CommentInput struct {
		CommentText     string `json:"comment_text" binding:"required"` // 评论内容
		RecipeID        uint   `json:"recipe_id" binding:"required"`    // 评论对应的食谱ID
		ParentCommentID *uint  `json:"parent_comment_id"`               // 父评论ID，null表示没有父评论                            // 点赞数
	}

	// 绑定到json
	var input CommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid input",
			"test": "绑定json出错"})
		return
	}

	// 验证食谱是否存在
	var recipe models.Recipe
	if err := global.DB.First(&recipe, input.RecipeID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found",
			"test": "验证食谱是否存在出错"})
		return
	}

	// 如果是回复评论，验证父评论是否存在
	if input.ParentCommentID != nil {
		var parentcomment models.Comment
		if err := global.DB.First(&parentcomment, *input.ParentCommentID).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found",
				"test": "验证父评论是否存在出错"})
			return
		}

		// 检查父评论是否属于同一食谱
		if parentcomment.RecipeID != input.RecipeID {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment does not belong to this recipe",
				"test": "检查父评论是否属于同一食谱出错"})
			return
		}
	}

	// 创建评论
	comment := models.Comment{
		CommentText:     input.CommentText,
		CommentUserID:   uint(userid.(uint8)), // 将userid转换为uint
		RecipeID:        input.RecipeID,
		ParentCommentID: input.ParentCommentID,
		Likes:           0, // 初始化点赞数为0
	}

	// 保存评论到数据库
	if err := global.DB.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment added successfully",
		"comment": gin.H{
			"comment_id":        comment.ID,
			"comment_text":      comment.CommentText,
			"recipe_id":         comment.RecipeID,
			"parent_comment_id": comment.ParentCommentID,
			"likes":             comment.Likes,
			"comment_user_id":   comment.CommentUserID,
		},
	})
}

// 获取食谱评论
func GetRecipeComments(ctx *gin.Context) {
	// 获取食谱ID
	var input struct {
		RecipeID int `json:"recipe_id" binding:"required"` // 食谱ID
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 设置输出结构体
	type CommentOutput struct {
		ID              uint   `json:"comment_id"`        // 评论ID
		CommentText     string `json:"comment_text"`      // 评论内容
		RecipeID        uint   `json:"recipe_id"`         // 评论对应的食谱ID
		ParentCommentID *uint  `json:"parent_comment_id"` // 父评论ID，null表示没有父评论
		Likes           uint   `json:"likes"`             // 点赞数
		CommentUserID   uint   `json:"comment_user_id"`   // 评论用户ID
		UserName        string `json:"user_name"`         // 评论用户名称
		UserAvatar      string `json:"user_avatar"`       // 评论用户头像
	}
	// 查询食谱评论
	var comments []CommentOutput
	if err := global.DB.Table("comments").Select("comments.id, comment_text, parent_comment_id,recipe_id,likes,comment_user_id,users.name as user_name,users.avatar as user_avatar").
		Joins("LEFT JOIN users ON comments.comment_user_id = users.id").
		Where("comments.recipe_id = ?", input.RecipeID).Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	// 返回评论列表
	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

// 删除评论
func DeleteComment(ctx *gin.Context) {
	// 获取当前用户ID
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 获取用户类型
	usertype, exists := ctx.Get("usertype")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User type not found"})
		return
	}

	// 获取评论ID
	var input struct {
		CommentID uint `json:"comment_id" binding:"required"` // 评论ID
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 查询评论是否存在
	var comment models.Comment
	switch usertype.(uint8) {
	case 1: //普通用户
		if err := global.DB.Where("id = ? AND comment_user_id = ?", input.CommentID, userid).First(&comment).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or not authorized to delete"})
			return
		}
	case 0: // 管理员
		if err := global.DB.Where("id = ?", input.CommentID).First(&comment).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or not authorized to delete"})
			return
		}
	}

	// 删除评论，需删除子评论
	if err := global.DB.Unscoped().Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}

// (直接操作Mysql数据库)根据点赞状态进行点赞数增加与减少
func ToggleCommentLike(ctx *gin.Context) {
	// 获取当前用户ID
	userid, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 获取评论ID
	var input struct {
		CommentID uint `json:"comment_id" binding:"required"` // 评论ID
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 查询评论是否存在
	var comment models.Comment
	if err := global.DB.Where("id = ?", input.CommentID).First(&comment).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// 检查用户是否已经点赞
	var existingLike models.CommentLike
	if err := global.DB.Where("user_id = ? AND comment_id = ?", userid, input.CommentID).First(&existingLike).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户未点赞，添加点赞记录
			newLike := models.CommentLike{
				UserID:    uint(userid.(uint8)), // 将userid转换为uint
				CommentID: input.CommentID,
			}

			// 创建新的点赞记录
			if err := global.DB.Create(&newLike).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to add like"})
				return
			}

			// 增加点赞数
			comment.Likes++
			if err := global.DB.Save(&comment).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add like count"})
				return
			}
		}
		// 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Like added successfully",
		})
	} else {
		// 用户已点赞，删除点赞记录
		if err := global.DB.Unscoped().Delete(&existingLike).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
			return
		}

		// 减少点赞数，确保不小于0
		if comment.Likes > 0 {
			comment.Likes--
			if err := global.DB.Save(&comment).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like count"})
				return
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Like count is already zero"})
			return
		}
		// 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Like removed successfully",
		})
	}

// (操作redis数据库)记录用户评论的点赞状态改变
// (操作Mysql数据库)根据redis记录状态改变Mysql中实际数据，间隔时间：1h
}
