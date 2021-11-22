package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(c *gin.Context) {
	smsForm := forms.SmsForm{}
	if err := c.ShouldBindJSON(&smsForm); err != nil {
		HandleValidtorError(c, err)
		return
	}
	confSms := global.Conf.SmsConfig
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", confSms.AccessKeyId, confSms.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = smsForm.Mobile                //手机号
	request.QueryParams["SignName"] = "mxshop"                          //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_228136691"               //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	if err != nil {
		fmt.Print(err.Error())
	}
	//将验证码保存起来 - redis
	var conf config.RedisConfig
	var sms forms.SmsForm
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
	})
	rdb.Set(context.Background(), sms.Mobile, smsCode, time.Duration(conf.Expire)*time.Second)

	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
