package routers

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	routers := gin.Default()

	routers.Use()
	{
		SetAuthRouter(routers)
	}
	return routers
}
