package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetFoodRouter(routers *gin.Engine) {
	food := routers.Group("/food")
	food.Use(middlewares.CheckJWT())
	{
		// 增加基础食材种类(管理员可以添加，普通用户只能调用另外的接口函数)
		food.POST("/add", middlewares.CheckUserType(0, 1), controllers.AddFood)
		// 删除已有食材-Name(连锁删除)
		food.DELETE("/delete", middlewares.CheckUserType(0), controllers.DeleteFoodByName)
		// 获取所有食材
		food.GET("/all", middlewares.CheckUserType(0, 1), controllers.GetAllFood)
		// 获取食材-Name
		food.GET("/get", middlewares.CheckUserType(0, 1), controllers.GetFoodByName)
	}
}
