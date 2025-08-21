package middlewares

import (
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 验证JWT令牌
func CheckJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userid, usertype uint8
		var rtoken string
		atoken, err := ctx.Cookie("atoken")
		if err == nil {
			// 如果存在atoken，则验证其有效性
			userid, usertype, err = utils.ParseJWT(atoken)
		}

		// 如果没有atoken或验证失败，则尝试使用rtoken刷新atoken
		if err != nil {
			rtoken, err = ctx.Cookie("rtoken")
			if err != nil {
				// rtoenken不存在，返回401状态码，要求用户重新登录
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
				ctx.Abort()
				return
			}

			// 刷新atoken和rtoken
			atoken, rtoken, err = utils.RefreshAToken(rtoken)
			if err != nil {
				ctx.AbortWithStatusJSON(401, gin.H{"error": "请重新登录"})
				return
			}
			// 设置新的atoken和rtoken到Cookie
			ctx.SetCookie("atoken", atoken, 3600*2, "/", "", false, true)
			ctx.SetCookie("rtoken", rtoken, 3600*24*30, "/", "", false, true)
			// 重新解析atoken
			userid, usertype, err = utils.ParseJWT(atoken)
			if err != nil {
				ctx.AbortWithStatusJSON(401, gin.H{"error": "请重新登录"})
				return
			}
		}
		// 将用户ID和用户类型存储在上下文中
		ctx.Set("usertype", usertype)
		ctx.Set("userid", userid)
		ctx.Next()
	}
}
