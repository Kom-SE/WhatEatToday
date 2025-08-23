package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetUserRouter(routers *gin.Engine) {
	user := routers.Group("/user")
	user.Use(middlewares.CheckJWT())
	{
		// 获取用户信息面板
		user.GET("/info", controllers.GetUserInfo)
		// 修改用户信息
		user.PATCH("/update", controllers.UpdateUserInfo)
		// 修改用户头像
		user.PATCH("/avatar", controllers.UpdateUserAvatar)
		// 注销用户（软删除）
		user.DELETE("/delete", middlewares.CheckUserType(0, 1), controllers.DeleteUser)
		// 收藏食谱
		user.POST("/collect", middlewares.CheckUserType(1), controllers.CollectRecipe)
		// 获取收藏所有食谱
		user.GET("/collects", middlewares.CheckUserType(1), controllers.GetCollectedRecipes)
		// 取消收藏食谱
		user.DELETE("/collect", middlewares.CheckUserType(1), controllers.UncollectRecipe)
	}
}
