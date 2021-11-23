package config

type WebConfig struct {
	Port            int `mapstructure:"port"`
	*UserSrvConfigs `mapstructure:"userserver"`
	*ServerConfig   `mapstructure:"serverinfo"`
	*LogConfig      `mapstructure:"log"`
	*MysqlConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*SmsConfig      `mapstructure:"sms"`
	*ConsulConfig   `mapstructure:"consul"`
}

type UserSrvConfigs struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name        string `mapstructure:"name"`
	UserSrvInfo string `mapstructure:"user_srv"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DbName       string `mapstructure:"dbname"`
	MaxConn      int    `mapstructure:"max_conn"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolsize"`
	Expire   int64  `mapstructure:"expire"`
}
type SmsConfig struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type NaconsConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	NameSpace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
