package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUserType(allowedType ...int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userType, exists := ctx.Get("user_type")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "User type not found in context",
			})
			ctx.Abort()
			return
		}

		userType, _ = userType.(int)

		allowed := false
		for _, t := range allowedType {
			if userType == t {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied: insufficient permissions",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
