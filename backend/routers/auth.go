package routers

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func SetAuthRouter(routers *gin.Engine) {
	auth := routers.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}
}
