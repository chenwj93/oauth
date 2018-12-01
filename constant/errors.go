package constant

import "errors"

const (
	PASSWORD_MISTAKE = "密码错误"
	NEED_LOGIN = "请先登录"
)

var (
	GET_PWD_ERR = errors.New("miss param [clientId/username]")
	AUTH_DATA_ERR = errors.New("api access control list data is error")
)