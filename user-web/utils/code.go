package utils

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvaildParam
	CodeUserExist
	CodeUserNotExist
	CodeInvaildPassword
	CodeInvaildUserName
	CodeServerBusy
	CodeInvaildMobile
	CodeNeedAuth
	CodeInvaildAuth
	CodeNeedLogin
	CodeCreateSucess
	CodeBadRequest
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "请求成功",
	CodeInvaildParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvaildPassword: "密码错误",
	CodeInvaildUserName: "用户名错误",
	CodeInvaildMobile:   "手机号不存在",
	CodeServerBusy:      "服务繁忙",

	CodeNeedAuth:     "需要验证信息",
	CodeInvaildAuth:  "token验证失败",
	CodeNeedLogin:    "请登录",
	CodeCreateSucess: "创建成功",
	CodeBadRequest:   "请求错误",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
