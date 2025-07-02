package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetRecipeRouter(routers *gin.Engine) {
	recipe := routers.Group("/recipe")
	recipe.Use(middlewares.CheckJWT())
	{
		// 用户添加食谱
		recipe.POST("/add", middlewares.CheckUserType(0, 1), controllers.AddRecipe)
		// 用户更新自己的食谱
		recipe.PATCH("/update/:id", middlewares.CheckUserType(0, 1), controllers.UpdateRecipe)
		// 用户删除自己的食谱
		recipe.DELETE("/delete", middlewares.CheckUserType(0, 1), controllers.DeleteRecipe)

	}
}
