package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"mxshop-api/user-web/utils"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 200, 5, 0.7, 200)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.L().Error("验证码生成错误：", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, "验证码生成错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
