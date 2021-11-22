package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("base")
	{
		userRouter.GET("captcha", api.GetCaptcha)
		userRouter.GET("sen_sms", api.SendSms)
	}

}
