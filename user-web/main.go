package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
)

func main() {
	initialize.InitConfig()
	fmt.Println(global.Conf)
	initialize.InitSrvConn()

	err := initialize.InitLogger(global.Conf.LogConfig)
	if err != nil {
		return
	}
	err = api.InitTrans("zh")
	if err != nil {
		zap.L().Error("初始化翻译器失败")
		return
	}
	//验证器注册
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("mobile", utils.ValidateMobie)
		if err != nil {
			zap.L().Error("验证器验证失败", zap.Error(err))
			return
		}
		_ = v.RegisterTranslation("mobile", api.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	conf := global.Conf
	viper.AutomaticEnv()
	//如果是本地开发环境端口号固定，线上环境启动获取端口号
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			conf.Port = port
		}
	}
	r := initialize.InitRouter()
	err = r.Run(fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		zap.L().Error("服务启动失败")
		return
	}

}
