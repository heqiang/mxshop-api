package initialize

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

// InitSrvConn 负载均衡
func InitSrvConn() {
	consulConfig := global.Conf.ConsulConfig
	dial, err := grpc.Dial(fmt.Sprintf("%s:%d/%s?wait=14s",
		consulConfig.Host, consulConfig.Port, consulConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.L().Error("用户服务连接失败：", zap.Error(err))
		return
	}
	client := proto.NewUserClient(dial)
	global.Client = client
}

func InitProtoClient() {
	// 从注册中心获取用户信息
	cfg := consulapi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.Conf.ConsulConfig.Host,
		global.Conf.ConsulConfig.Port)

	consulclient, err := consulapi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
	//筛选需要的Service
	data, err := consulclient.Agent().ServicesWithFilter(fmt.Sprintf("Service==\"%s\"",
		global.Conf.ConsulConfig.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	dial, err := grpc.Dial(fmt.Sprintf("%s:%d/mxshop?wait=14s&tag=hq", userSrvHost, userSrvPort), grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.L().Error("grpc连接失败", zap.Error(err))
		return
	}
	zap.L().Info(fmt.Sprintf("获取到的grpc服务地址为=>%s:%d", userSrvHost, userSrvPort))

	client := proto.NewUserClient(dial)
	global.Client = client
}
