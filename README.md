# 电商系统mxshop-web

> 以下api功能模块均采用基于gin框架的grpc远程服务调用，
>
> 此web部分只提供服务调用的相关接口
>
> 接口文档暂未更新，后续有需要在更新
>
> 服务接口:https://github.com/heqiang/mxshop_srvs
>
> 由于阿里云的短信服务暂时未开通成功，所以用户相关手机短信注册等功能暂不能实现，逻辑已完善，
> 后面项目完善后 采用go_zero框架进行重构

##### 一、api功能模块：  

+ 1 用户模块  

  + 用户登录
  + 用户注册
  + 用户信息修改 
  + 用户列表查看  

+ 用户权限验证

  + jwt

    

环境:

+ sdk:go1.17.3 windows/amd64  

+ 编辑器：goland 
+ 技术栈 go、gin、grpc、consul、nacos、docker
+ 数据库：mysql8、redis5.0

####  二、安装下载

2.1、下载:

api服务端

> git clone git@github.com:heqiang/mxshop_srvs.git   
>
> go mod tidy

api客户端

>  git@github.com:heqiang/mxshop-api.git
>
> go mod tidy

2.2 相关服务配置

grpc服务

> cd proto
>
> protoc --go_out=plugins=grpc:. *.proto

数据库服务

> ```console
> docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:8.0
> docker run --name some-redis -p 6379:6379 -d redis
> ```

服务发现与注册

> docker run -d --name=dev-consul  -p 8500:8500  -e CONSUL_BIND_INTERFACE=eth0 consul

配置中心nacos

> ```shell
> docker run --name nacos-quick -e MODE=standalone -p 8848:8848 -d nacos/nacos-server:2.0.2
> ```

进入nacos

step1 http://127.0.0.1:8848/nacos

step2 登录 账号密码都是nacos

step3 根据config.yaml对nacos进行相关配置，只需要配置命名空间及配置列表，配置信息如下

```json
{
  "port": 8081,
  "userserver": {
     # 更改为你本地的ipv4地址
    "host": "192.168.31.101",
    "port": 8081
  },
  "serverinfo": {
    "name": "hq",
    "user_srv": "mxshop"
  },
  "log": {
    "level": "debug",
    "filename": "web_app_log.log",
    "max_size": 200,
    "max_age": 30,
    "max_backups": 7
  },
  "redis": {
    "host": "redis507",
    "port": 6379,
    "password": "",
    "db": 0,
    "poolsize": 100,
    "expire": 300
  },
  "sms": {
    "accessKeyId": "xx",
    "accessKeySecret": "xx"
  },
  "consul": {
    "host": "127.0.0.1",
    "port": 8500,
     # 服务发现的名字 
    "name": "mxshop_srv"
  }
}
```

#### 三 运行：

<font style="color:red">注意</font>:运行之前需要先保证api服务的正常启动

> go  run main.go



