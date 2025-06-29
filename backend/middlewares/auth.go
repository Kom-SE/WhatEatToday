package middlewares

import (
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 验证JWT令牌
func CheckJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			ctx.Abort()
			return
		}
		username, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}
		user_type, ok := ctx.Get("user_type")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context"})
			return
		}
		ctx.Set("user_type", user_type)
		ctx.Set("username", username)
		ctx.Next()
	}
}
