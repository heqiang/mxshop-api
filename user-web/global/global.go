package global

import (
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

var (
	Client      proto.UserClient
	Conf        *config.WebConfig
	NacosConfig *config.NaconsConfig
)
