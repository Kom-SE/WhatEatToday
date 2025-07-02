package routers

import (
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func SetCommentRouter(routers *gin.Engine) {
	comment := routers.Group("/comment")
	comment.Use(middlewares.CheckJWT())
	{
		// 增加评论
		// 删除评论
		// 获取食谱评论
		// 增加点赞
		// 减少点赞
		// 禁止该食谱评论
	}
}
