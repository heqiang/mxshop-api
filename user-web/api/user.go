package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
	"mxshop-api/user-web/utils"
	"mxshop-api/user-web/utils/jwt"
	"net/http"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Internal:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "网络错误",
				})
			default:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "其他错误" + e.Message(),
				})

			}
		}
	}
}

func HandleValidtorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		zap.L().Error("invaild param", zap.Error(err))
		utils.ResponseError(ctx, utils.CodeInvaildParam)
		return
	}
	utils.ResponseErrorWithMsg(ctx, utils.CodeInvaildParam, RemoveTopStruct(errs.Translate(Trans)))
	return
}

func GetUserList(ctx *gin.Context) {
	cfg := consulapi.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.Conf.ConsulConfig.Host,
		global.Conf.ConsulConfig.Port)

	consulclient, err := consulapi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	userSrvHost := ""
	userSrvPort := 0
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
	dial, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("grpc连接失败", zap.Error(err))
		return
	}
	client := proto.NewUserClient(dial)
	zap.L().Info("用户列表查询")
	resp, err := client.GetUserList(context.Background(), &proto.PageInfo{})
	if err != nil {
		zap.L().Error("用户列表获取失败：", zap.Error(err))
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range resp.Data {
		data := make(map[string]interface{})
		data["id"] = value.Id
		data["name"] = value.NickNmae
		data["birthday"] = value.Birthday
		data["gender"] = value.Gender
		result = append(result, data)
	}
	utils.RsponseSuccess(ctx, utils.CodeSuccess, gin.H{
		"total": resp.Total,
		"data":  result,
	})
}

func PasswordLogin(c *gin.Context) {
	loginForm := forms.PasswordLoginForm{}
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		HandleValidtorError(c, err)
		return
	}
	if !store.Verify(loginForm.CaptchaId, loginForm.Captcha, true) {
		utils.RsponseSuccess(c, utils.CodeBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}
	res, err := global.Client.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: loginForm.Mobile,
	})
	if err != nil {
		zap.L().Error("GetUserByMobile error:", zap.Error(err))
		utils.ResponseError(c, utils.CodeInvaildMobile)
		return
	}
	checkBool, err := global.Client.CheckPassword(context.Background(), &proto.CheckPwd{
		EncryptedPassword: res.Password,
		Password:          loginForm.Password,
	})
	if !checkBool.Success {
		utils.ResponseError(c, utils.CodeInvaildPassword)
		return
	}
	token, err := jwt.GenToken(res.NickNmae, int64(res.Id))
	if err != nil {
		zap.L().Error("token生成有误", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	utils.RsponseSuccess(c, utils.CodeSuccess, gin.H{
		"token": token,
	})
}

func CreateUser(c *gin.Context) {
	var userinfo forms.UserInfo
	if err := c.ShouldBindJSON(&userinfo); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("参数有误", zap.Error(err))
			utils.ResponseError(c, utils.CodeInvaildParam)
			return
		}
		zap.L().Error("参数有误", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeInvaildParam, RemoveTopStruct(errs.Translate(Trans)))
		return
	}
	_, err := global.Client.CreateUser(c, &proto.CreateUserInfo{
		NickName: userinfo.NickName,
		PassWord: userinfo.PassWord,
		Mobile:   userinfo.Mobile,
	})
	if err != nil {
		zap.L().Error("用户创建失败", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, err.Error())
		return
	}
	utils.RsponseSuccess(c, utils.CodeCreateSucess, nil)

}

type UserInfo struct {
	UserId   int64
	NickName string `json:"nickname"`
	Gender   string `json:"gender"`
	Birthday uint64 `json:"birthday"`
}

func UpdateUserInfo(c *gin.Context) {
	var userinfo UserInfo
	_ = c.ShouldBindJSON(&userinfo)
	UserId, err := utils.GetCurrentUser(c)
	if err != nil {
		zap.L().Error("token验证失败:", zap.Error(err))
		utils.ResponseError(c, utils.CodeNeedAuth)
		return
	}
	_, err = global.Client.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(UserId),
		NickName: userinfo.NickName,
		Gender:   userinfo.Gender,
		Birthday: userinfo.Birthday,
	})
	if err != nil {
		zap.L().Error("用户更新失败", zap.Error(err))
		utils.ResponseError(c, utils.CodeUserNotExist)
		return
	}
	utils.RsponseSuccess(c, utils.CodeSuccess, nil)

}

func RegisteForm(c *gin.Context) {
	var register forms.RegisterForm
	if err := c.ShouldBindJSON(&register); err != nil {
		HandleValidtorError(c, err)
		return
	}
	conf := global.Conf.RedisConfig
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
	})
	value, err := rdb.Get(context.Background(), register.Mobile).Result()
	if err == redis.Nil {
		utils.ResponseErrorWithMsg(c, utils.CodeBadRequest, gin.H{
			"code": "验证码有误",
		})
		return
	}
	if value != register.Code {
		utils.ResponseErrorWithMsg(c, utils.CodeBadRequest, gin.H{
			"code": "验证码有误",
		})
		return
	}
	_, err = global.Client.CreateUser(c, &proto.CreateUserInfo{
		NickName: register.Mobile,
		PassWord: register.Password,
		Mobile:   register.Mobile,
	})
	if err != nil {
		zap.L().Error("用户创建失败", zap.Error(err))
		utils.ResponseErrorWithMsg(c, utils.CodeServerBusy, err.Error())
		return
	}
	utils.RsponseSuccess(c, utils.CodeCreateSucess, nil)
}
func GetUserDetial(c *gin.Context) {

}
