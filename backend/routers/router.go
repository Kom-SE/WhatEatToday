package routers

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	routers := gin.Default()

	routers.Static("/image", "./images")

	routers.Use()
	{
		SetAuthRouter(routers)
		SetUserRouter(routers)
		SetFoodRouter(routers)
		SetLableRouter(routers)
		SetRecipeRouter(routers)
	}
	return routers
}
