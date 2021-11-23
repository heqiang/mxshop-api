package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"time"
)

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return
	}
	if err := viper.Unmarshal(&global.NacosConfig); err != nil {
		panic(fmt.Errorf("unmarshal conf failed err%s\n", err))
	}
	//配置热加载
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&global.NacosConfig); err != nil {
			panic(fmt.Errorf("unmarshal conf failed err:%s\n", err))
		}
	})

	//从nacos中
	sc := []constant.ServerConfig{{
		IpAddr: global.NacosConfig.Host,
		Port:   uint64(global.NacosConfig.Port),
	}}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NameSpace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "log",
		CacheDir:            "cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		zap.L().Error("nacos 客户端创建失败:", zap.Error(err))
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), &global.Conf)
	if err != nil {
		zap.L().Error("nacos反系列化失败", zap.Error(err))
		return
	}
	//fmt.Println(content) //字符串 - yaml
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生了变化...")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

	time.Sleep(300 * time.Second)

}
