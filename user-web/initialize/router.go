package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/router"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	ApiGroup := r.Group("/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return r

}
