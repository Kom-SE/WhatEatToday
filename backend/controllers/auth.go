package controllers

import (
	"main/global"
	"main/models"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册账户
func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(
			http.StatusBadRequest, gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	user.UserType = 1 // 默认用户类型为1（普通用户）

	var flag bool
	flag, err := utils.IsExists(&models.User{}, "username", user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if flag {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "用户已经存在",
		})
		return
	}

	// 生成bcrypt密码
	hashpw, err := utils.BcryptPW(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.Password = hashpw
	token, err := utils.GenerateJWT(user.ID, user.UserType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 设置默认用户类型为1（普通用户）,其他权限序需要后台修改

	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"token":   token})
}

// 登录账户
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	if err := global.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !utils.CheckPW(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "wong password",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.UserType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "登录成功",
		"uid":      user.ID,
		"token":    token,
		"usertype": user.UserType,
	})
}
