package controllers

import (
	"main/global"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 增加基础标签种类
func AddLable(ctx *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Type        uint8  `json:"type" binding:"required"` // 1:口味 2:做法
		Description string `json:"description"`             // 标签具体内容：null则没有，1则口味，2则做法
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	// 检查是否已存在同名标签
	if err := global.DB.Where("name = ?", input.Name).First(&models.Label{}).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Label with this name already exists",
		})
		return
	}

	// 无同名标签，创建新标签
	label := models.Label{
		Name:        input.Name,
		Type:        input.Type,
		Description: input.Description,
	}

	if err := global.DB.Create(&label).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create label",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Label created successfully",
	})

}

// 删除已有标签-Name(连锁删除)
func DeleteLableByName(ctx *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	var label models.Label
	if err := global.DB.Where("name = ?", input.Name).First(&label).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Label not found",
		})
		return
	}

	if err := global.DB.Unscoped().Delete(&label).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete label",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Label deleted successfully",
	})
}

// 获取所有标签
func GetAllLable(ctx *gin.Context) {
	var labels []models.Label
	if err := global.DB.Find(&labels).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve labels",
		})
		return
	}
	type LabelResponse struct {
		Name        string `json:"name"`
		Type        uint8  `json:"type"`        // 1:口味 2:做法
		Description string `json:"description"` // 标签具体内容：null则没有，1则口味，2则做法
	}
	var response []LabelResponse
	for _, label := range labels {
		response = append(response, LabelResponse{
			Name:        label.Name,
			Type:        label.Type,
			Description: label.Description,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"labels": response,
	})
}

// 获取标签-Name-精准查询
func GetLableByName(ctx *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	var label models.Label
	if err := global.DB.Where("name = ?", input.Name).First(&label).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Label not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"label": gin.H{
			"name":        label.Name,
			"type":        label.Type,
			"description": label.Description,
		},
	})
}

// 获取标签-Name-模糊查询
func GetLableLikeName(ctx *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
		})
		return
	}

	// 使用模糊查询获取标签
	var labels []models.Label
	searchPattern := "%" + input.Name + "%"
	if err := global.DB.Where("name LIKE ?", searchPattern).Find(&labels).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve labels",
		})
		return
	}

	if len(labels) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"lables": "null",
			"count":  0,
		})
	} else {
		var response []gin.H
		for _, label := range labels {
			response = append(response, gin.H{
				"name":        label.Name,
				"type":        label.Type,
				"description": label.Description,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"labels": response,
			"count":  len(labels),
		})
	}
}
