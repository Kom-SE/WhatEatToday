package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetLableRouter(routers *gin.Engine) {
	lable := routers.Group("/lable")
	lable.Use(middlewares.CheckJWT())
	{
		// 增加基础标签种类
		lable.POST("/add", middlewares.CheckUserType(0), controllers.AddLable)
		// 删除已有标签-Name(连锁删除)
		lable.DELETE("/delete", middlewares.CheckUserType(0), controllers.DeleteLableByName)
		// 获取所有标签
		lable.GET("/all", middlewares.CheckUserType(0, 1), controllers.GetAllLable)
		// 获取标签-Name-精准查询
		lable.GET("/get", middlewares.CheckUserType(0, 1), controllers.GetLableByName)
		// 获取标签-Name-模糊查询
		lable.GET("/likename", middlewares.CheckUserType(0, 1), controllers.GetLableLikeName)
	}
}
