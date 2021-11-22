package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	userRouter.Use(middlewares.Cors())
	{
		userRouter.POST("/login", api.PasswordLogin)
		userRouter.POST("/create", api.RegisteForm)
		//userRouter.Use(middlewares.JWTAuthMiddleware())
		userRouter.GET("list", api.GetUserList)
		userRouter.POST("/updateUserInfo", api.UpdateUserInfo)
	}

}
