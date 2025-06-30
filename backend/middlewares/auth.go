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
		userid, usertype, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("usertype", usertype)
		ctx.Set("userid", userid)
		ctx.Next()
	}
}
