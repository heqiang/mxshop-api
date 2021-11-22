package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/utils"
	"mxshop-api/user-web/utils/jwt"
)

const (
	CTXUSERUSERNAMEKEY = "username"
	CTXUSERIDKEY       = "userid"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			utils.ResponseError(c, utils.CodeNeedAuth)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(authHeader)
		if err != nil {
			utils.ResponseError(c, utils.CodeInvaildAuth)
			c.Abort()
			return
		}
		// 将当前请求的UserId信息保存到请求的上下文c上
		c.Set(CTXUSERUSERNAMEKEY, mc.Username)
		c.Set(CTXUSERIDKEY, mc.UserId)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息

	}
}
