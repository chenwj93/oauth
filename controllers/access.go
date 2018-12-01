package controllers

import (
	"utils"
	"oauth/models"
	"context"
	"oauth/constant"
	"time"
	"encoding/json"
	"logs"
	"strings"
	"fmt"
	"oauth/remote"
)

type Accesser interface {
	Register([]byte) (*utils.Response, error)
	Login(map[string]interface{}) (*utils.Response, error)
	Authenticate(map[string]interface{}) (*utils.Response, error)
	ModifyPassword(map[string]interface{}) (*utils.Response, error)
	Logout() (*utils.Response, error)
	AuthHadModified(map[string]interface{}) (*utils.Response, error)
	ChangeUserStatus(map[string]interface{}) (*utils.Response, error)
}

type Access struct {
	context.Context
}

var RootClient string

func (a *Access) Register(paramJson []byte) (*utils.Response, error) {
	token := utils.ParseString(a.Value("token"))
	loginInfo, ok := models.Lim.Get(token)
	if !ok {
		return utils.ConcatDeny()
	}
	var object models.OauthUser
	err := json.Unmarshal(paramJson, &object)
	if err != nil {
		logs.Error(err)
		return utils.ConcatFormat()
	}
	object.Password = models.PasswordEncypt(object.Password)
	object.CreatedAt = time.Now()
	createBranch := false
	if object.ClientId == ""{
		cid := loginInfo.User.ClientId
		if strings.Contains(cid, "#"){ //分公司创建员工账号
			object.ClientId = cid
		} else { //集团创建子公司账号
			object.ClientId = "##" + loginInfo.User.ClientId
			createBranch = true
		}
	}
	id, err := object.Add()
	var param interface{}
	if err == nil {
		if createBranch {
			object.Id = id
			err = object.GetById()
			param = map[string]interface{}{"loginId": id, "clientId": object.ClientId}
		} else {
			param = id
		}
	}
	return utils.PackageReturn(param, "用户名不可重复", err)
}

func (a *Access) Login(paramInput map[string]interface{}) (*utils.Response, error) {
	userName, okId := paramInput["userName"]
	password, okPs := paramInput["password"]
	logs.Info(userName)
	if !okId || !okPs {
		return utils.ConcatFormat()
	}
	var userctl OauthUserController
	abstract, token, uId, status, clientId := models.CheckDuplicateLogin(userName, password)
	if status == constant.LOGIN_FAILD { // 密码错误
		return utils.ConcatFailed(nil, constant.PASSWORD_MISTAKE)
	}
	if status == constant.LOGIN_FIRST { // 没有重复登录
		user, ok := userctl.CheckPassword(userName, password)
		if !ok { // 密码校验不通过
			return utils.ConcatFailed(nil, constant.PASSWORD_MISTAKE)
		}
		err := userctl.StoreUserInfo(user, token, abstract)
		if err != nil{
			return utils.ConcatFailed(nil)
		}
		uId = user.Id
		clientId = user.ClientId
	}
	go OauthAccessTokensController{}.Save(token, uId, constant.USER_BACK)
	var cType string = "company"
	if strings.Contains(clientId, RootClient) {
		cType = "fs"
	}else if strings.Contains(clientId, "#") {
		cType = "branchCom"
	}
	return utils.ConcatSuccess(map[string]interface{}{"access_token":token, "type":cType})
}

func (a *Access) Authenticate(paramInput map[string]interface{}) (*utils.Response, error) {
	token := utils.ParseString(a.Value("token"))
	loginInfo, ok := models.Lim.Get(token)
	if !ok || loginInfo.User.Status != 0{
		return utils.ConcatErr("please login ！", 401, 200, nil)
	}
	//if loginInfo.User.ClientId != RootClient { // 非风数账号，需要鉴权
	//	if loginInfo.User.AuthSet == nil{
	//		loginInfo.User.GetUserAuth()
	//	}
	//	if !loginInfo.User.CheckAuth(utils.ParseString(paramInput["operate"])) {
	//		return utils.ConcatDeny()
	//	}
	//}
	if loginInfo.User.UserName == "" && loginInfo.User.UserName != constant.NULLNAME{
		loginInfo.User.UserName = remote.GetUserName(loginInfo.User.Id, loginInfo.User.ClientId)
	}
	ret := map[string]interface{}{"userId":loginInfo.User.Id,
		"userName": utils.TEE(loginInfo.User.UserName == constant.NULLNAME, "", loginInfo.User.UserName),
		"clientId":loginInfo.User.ClientId}
	return utils.ConcatSuccess(ret)
}

func (a *Access) ModifyPassword(paramInput map[string]interface{}) (*utils.Response, error) {
	token := a.Context.Value("token")
	passwordOld, passwordNew := paramInput["passwordOld"], paramInput["passwordNew"]

	if utils.IsEmptyBat(token, passwordOld, passwordNew){
		return utils.ConcatFormat()
	}
	loginInfo, ok := models.Lim.Get(utils.ParseString(token))
	if !ok {
		return utils.ConcatFailed(nil, constant.NEED_LOGIN)
	}
	userInfo := loginInfo.User
	pwdOld := models.PasswordEncypt(passwordOld)
	if pwdOld != userInfo.Password{
		return utils.ConcatFailed(nil, constant.PASSWORD_MISTAKE)
	}
	pwdNew := models.PasswordEncypt(passwordNew)
	_, err := (&models.OauthUser{Id:userInfo.Id, Password:pwdNew, UpdatedAt:time.Now()}).Update("Password", "UpdatedAt")
	if err == nil{
		userInfo.Password = pwdNew
	}
	return utils.PackageReturn(nil, nil, err)
}

func (a *Access) Logout() (*utils.Response, error) {
	token := utils.ParseString(a.Value("token"))
	loginInfo, ok := models.Lim.Get(token)
	if ok && loginInfo.User.LastLoginToken == token{
		abs := loginInfo.User.GetAbstract()
		models.Uim.Delete(abs)
	}
	models.Lim.Delete(token)
	OauthAccessTokensController{}.Expire(token)
	return utils.ConcatSuccess()
}

func (a *Access) AuthHadModified(paramInput map[string]interface{}) (*utils.Response, error) {
	userId, ok := paramInput["data"]
	if !ok {
		return utils.ConcatFormat()
	}
	userList, ok := userId.([]interface{})
	if !ok {
		return utils.ConcatFormat()
	}
	user := models.OauthUser{UpdatedAt: time.Now()}
	for _, uid := range userList {
		user.Id = utils.ParseInt64(uid)
		authUser, err := user.Update("AuthStatus", "UpdatedAt")
		if err == nil {
			userInfo, ok := models.Uim.Get(authUser.GetAbstract())
			if ok {
				userInfo.AuthStatus = constant.UPDATED
			}
		} else if !strings.HasSuffix(err.Error(), "no row found") {
			logs.Error(err)
			return utils.ConcatFailed(nil, err.Error())
		}
	}
	return utils.ConcatSuccess()
}

func (a *Access) ChangeUserStatus(paramInput map[string]interface{}) (*utils.Response, error) {
	clientId := utils.ParseString(paramInput["clientId"])
	userId, ok := paramInput["userId"].([]interface{})
	status := utils.ParseInt(paramInput["status"])
	var err error
	var nameList []models.StringValue
	if clientId != ""{
		nameList, err = models.ChangeUserStatus(clientId, nil, status)
	} else if ok && len(userId) != 0 {
		nameList, err = models.ChangeUserStatus(nil, userId, status)
	}
	if err == nil {
		go changeLoginUserStatus(nameList, status)
	}

	return utils.PackageReturn(nil, nil, err)
}

func changeLoginUserStatus(nameList []models.StringValue, status int)  {
	for i := range nameList {
		name := nameList[i].Value1
		abs := utils.Md5(name)
		userInfo, ok := models.Uim.Get(abs)
		if ok {
			logs.Info(fmt.Sprintf("change login user [%s] status to [%d]", name, status))
			userInfo.Status = int8(status)
			models.Lim.Delete(userInfo.LastLoginToken)
		}
	}
}
//
////
//func (a *Access) auth(operate string) bool {
//	token := utils.ParseString(a.Value("token"))
//	loginInfo, ok := models.Lim.Get(token)
//	if !ok {
//		return false
//	}
//	if loginInfo.User.AuthSet == nil{
//		loginInfo.User.GetUserAuth()
//	}
//	return loginInfo.User.CheckAuth(operate)
//}

