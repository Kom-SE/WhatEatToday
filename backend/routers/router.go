package routers

import (
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	routers := gin.Default()

	routers.Static("/image", "./images")

	routers.Use(middlewares.SetupCorsMiddleware())
	{
		SetAuthRouter(routers)
		SetUserRouter(routers)
		SetFoodRouter(routers)
		SetLableRouter(routers)
		SetRecipeRouter(routers)
		SetCommentRouter(routers)
	}
	return routers
}
