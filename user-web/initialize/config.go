package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"mxshop-api/user-web/global"
)

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return
	}
	if err := viper.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed err%s\n", err))
	}
	//配置热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&global.Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed err:%s\n", err))
		}
	})

}
