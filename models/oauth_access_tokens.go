package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"utils"
	"oauth/constant"
	"logs"
)

type OauthAccessTokens struct {
	Id         string             `json:"id" orm:"column(id);pk"`
	UserType   constant.USER_TYPE `json:"userType" orm:"column(user_type)"`
	UserId     int64              `json:"userId" orm:"column(user_id)"`
	ExpireTime time.Time          `json:"expireTime" orm:"column(expire_time)"`
	Status     int8               `json:"status" orm:"column(status)"`
	CreatedAt  time.Time          `json:"createdAt" orm:"column(created_at);type(timestamp);null"`
	UpdatedAt  time.Time          `json:"updatedAt" orm:"column(updated_at);type(timestamp);null"`
}

func (t *OauthAccessTokens) TableName() string {
	return "oauth_access_tokens"
}

func init() {
	orm.RegisterModel(new(OauthAccessTokens))
}

func (t *OauthAccessTokens) Add() (err error) {
	o := orm.NewOrm()
	_, err = o.Insert(t)
	return
}

func (t *OauthAccessTokens) Update(fields ...string) (err error) {
	o := orm.NewOrm()
	if _, err = o.Update(t, fields...); err != nil {
		logs.Error(err)
	}
	return
}

func ExpireTokenInit() {
	o := orm.NewOrm()
	_, err := o.Raw("update oauth_access_tokens set status = 1, updated_at = ? where status = 0 and expire_time < ?", time.Now(), time.Now()).Exec()
	logs.Log(logs.Error, err)
}

type TokenQuery struct {
	UserType int
	Status   int
}

func (t *TokenQuery) GetTokenList() (tokenList []OauthAccessTokens, err error) {
	o := orm.NewOrm()
	sqlCon := utils.SqlSelect{}
	sqlCon.Equal("user_type", t.UserType).
		Equal("status", t.Status).
		Greater("expire_time", time.Now()).
		OrderBy("user_id asc", "expire_time desc")
	_, err = o.Raw("select id, user_id, expire_time from oauth_access_tokens "+sqlCon.GetWhere().GetOrderBy().ToString(), sqlCon.GetParamWhere()...).QueryRows(&tokenList)
	return
}
