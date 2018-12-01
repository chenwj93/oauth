package controllers

import (
	"oauth/models"
	"utils"

	"time"
	"logs"
	"oauth/constant"
)

// @controller(oauthuser)
type OauthUserController struct {
}

func (c *OauthUserController) CheckPassword(userName, password interface{})(*models.OauthUser, bool){
	user := &models.OauthUser{Name:utils.ParseString(userName),
							}
	err := user.GetPassword()
	pwd := models.PasswordEncypt(password)
	if err != nil || user == nil || user.Password == utils.EMPTY_STRING || user.Password != pwd {
		return nil, false
	}
	return user, true
}

func (c *OauthUserController) StoreUserInfo(user *models.OauthUser, token, abs string) error{
	userInfo := &models.UserInfo{}
	err := utils.Assignment(*user, userInfo)
	if err != nil{
		logs.Error(err)
		return err
	}
	expireTime := time.Now().Add(constant.TokenExpiry)
	loginInfo := &models.LoginInfo{expireTime, userInfo}

	models.Uim.Put(abs, userInfo)
	models.Lim.Put(token, loginInfo)
	go ScheduleList.PushBack(&models.Index{abs, token, expireTime})
	userInfo.GetUserAuth()

	return nil
}