package remote

import (
	"utils"
	"oauth/constant"
	"logs"
)

func GetUserName(userId int64, clientId string) string {
	ret, err := utils.ServiceCall(constant.EmployeeApp, "employee/getName", "get", map[string]interface{}{"userId": userId}, map[string]interface{}{"clientId":clientId})
	if err != nil {
		logs.Error(err)
		return ""
	}
	if utils.ParseInt(ret["status"]) != 200 {
		return ""
	}
	uname := utils.ParseString(ret["data"])
	return utils.TEE(uname == "", constant.NULLNAME, uname).(string)
}