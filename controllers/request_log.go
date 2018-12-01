package controllers

import (
	"oauth/models"
	"utils"
	"oauth/common"
	"strings"
	"strconv"
)

// @controller(requestlog)
type RequestLogController struct {
}

func SaveRequestLog(request string, param map[string]interface{}, response *utils.Response, token string){
	common.HidePassword(param)
	var res strings.Builder
	res.WriteString("{Code:")
	res.WriteString(strconv.Itoa(response.Code))
	res.WriteString(", Json:")
	res.WriteString(string(response.Json))
	res.WriteString("}")
	object := models.RequestLog{Request:request, Param: utils.ParseStringFromMap(param), Response:res.String(), Token:token}
	object.Add()
}
