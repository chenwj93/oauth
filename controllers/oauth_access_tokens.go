package controllers

import (
	"time"

	"oauth/models"
	"oauth/constant"
	"fmt"
	"logs"
	"strings"
)

// @controller(oauthAccessTokens)
type OauthAccessTokensController struct {
}

//add created by  chenwj on 2018-07-06
// @router(method=save)
func (OauthAccessTokensController) Save(token string, userId int64, userType constant.USER_TYPE) (err error) {
	object := models.OauthAccessTokens{Id: token,
	 	UserType: userType,
		UserId: userId,
		ExpireTime: time.Now().Add(time.Hour * 8),
		CreatedAt: time.Now()}
	err = object.Add()
	return
}

// 初始化已登录用户信息
func (OauthAccessTokensController) Init() *models.IndexSlice{
	models.ExpireTokenInit()
	tokenQ := models.TokenQuery{1, 0}
	tL, e := tokenQ.GetTokenList()
	if e != nil {
		logs.Error("initial login user faild", e)
		return nil
	}
	logs.Info(tL)
	uL, e := tokenQ.GetUserList()
	if e != nil {
		logs.Error("initial login user faild", e)
		return nil
	}
	logs.Info(uL)
	apiList, err := models.GetApiLocal()
	if err != nil{
		logs.Fatal(err)
		panic("initial api list failed")
	}
	for i := 0; i < len(apiList); i++{
		api := strings.TrimSpace(apiList[i].Url)
		models.URLMAP.Store(api, apiList[i].Id)
		logs.Info(api)
	}
	return models.StoreLoginUser(tL, uL)
}


func (OauthAccessTokensController) Expire(token string) {
	object := models.OauthAccessTokens{Id: token,
		Status: 1,
		UpdatedAt: time.Now()}
	err := object.Update("Status", "UpdatedAt")
	if err != nil{
		logs.Error(fmt.Sprintf("token expire status update failed[token:%s]-%s", token, err.Error()))
	}
	return
}