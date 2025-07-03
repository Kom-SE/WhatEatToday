package routers

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetCommentRouter(routers *gin.Engine) {
	comment := routers.Group("/comment")
	comment.Use(middlewares.CheckJWT())
	{
		// 增加评论
		comment.POST("/add", middlewares.CheckUserType(1), controllers.AddComment)
		// 删除评论
		comment.DELETE("/delete", middlewares.CheckUserType(0, 1), controllers.DeleteComment)
		// 获取食谱评论
		comment.GET("/get", middlewares.CheckUserType(0, 1), controllers.GetRecipeComments)
		// 根据点赞状态来决定点赞增加还是减少
		comment.PATCH("/like", middlewares.CheckUserType(1), controllers.ToggleCommentLike)
	}
}
