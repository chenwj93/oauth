package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"utils"
	"oauth/constant"
	"logs"
	"oauth/common"
)

type OauthUser struct {
	Id          int64               `json:"id" orm:"column(oauth_user_id);auto"`
	Name        string              `json:"loginName" orm:"column(name);size(255)"`
	ClientId    string              `json:"clientId" orm:"column(client_id);size(255)"`
	MobilePhone string              `json:"phone" orm:"column(mobile_phone);size(255)"`
	Password    string              `json:"password" orm:"column(password);size(60)"`
	Singleton   int8                `json:"singleton" orm:"column(singleton)" description:"是否禁止重复登录"`
	AuthStatus  constant.AUTHSTATUS `json:"authStatus" orm:"column(auth_status)" description:"0-未获取权限  1-已获取权限  2-权限已更新"`
	CreatedAt   time.Time           `json:"createdAt" orm:"column(created_at);type(timestamp);null"`
	UpdatedAt   time.Time           `json:"updatedAt" orm:"column(updated_at);type(timestamp);null"`
	Deleted     int8                `json:"-" orm:"column(deleted);null"`
	Status      int8                `json:"-" orm:"column(status);null"`
}

func (t *OauthUser) TableName() string {
	return "oauth_user"
}

func init() {
	orm.RegisterModel(new(OauthUser))
}

func (t *OauthUser) GetAbstract() string {
	//uis := common.LinkInterface(t.ClientId, t.Name)
	abs := utils.Md5(t.Name)
	return abs
}

func (t *OauthUser) Add() (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(t)
	return
}

func (t *OauthUser) GetById() (err error) {
	o := orm.NewOrm()
	if err = o.Read(t); err == nil {
		return nil
	}
	return err
}

func (t *OauthUser) GetPassword() (err error) {
	if t.Name == "" {
		return constant.GET_PWD_ERR
	}
	o := orm.NewOrm()
	err = o.Raw("select oauth_user_id, name, client_id, password, auth_status, singleton from oauth_user where deleted = 0 and status = 0 and name = ?", t.Name).QueryRow(t)
	return
}

func (t *OauthUser) Update(fields ...string) (v *OauthUser, err error) {
	o := orm.NewOrm()
	if _, err = o.Update(t, fields...); err != nil {
		logs.Error(err)
	}

	return t, err
}

func (t *TokenQuery) GetUserList() (userList []OauthUser, err error) {
	o := orm.NewOrm()
	sqlCon := utils.SqlSelect{}
	sqlCon.Equal("ot.user_type", t.UserType).
		Equal("ot.status", t.Status).
		Greater("ot.expire_time", time.Now())
	sql := `
			SELECT
				oauth_user_id ,
				name ,
				client_id ,
				password ,
				auth_status,
				singleton
			FROM
				oauth_user ou
			JOIN oauth_access_tokens ot ON ou.oauth_user_id = ot.user_id and ou.status = 0
			` + sqlCon.GetWhere().GetOrderBy().ToString() +
		` group by oauth_user_id ,
			name,
			client_id ,
			password ,
			singleton
		order by oauth_user_id asc`
	_, err = o.Raw(sql, sqlCon.GetParamWhere()...).QueryRows(&userList)
	return
}

func UpdateAuthStatus(userId int64, authStatus constant.AUTHSTATUS) error {
	user := OauthUser{Id: userId, AuthStatus: authStatus, UpdatedAt: time.Now()}
	_, err := user.Update("AuthStatus", "UpdatedAt")
	if err != nil {
		logs.Error(err)
	}
	return err
}

func ChangeUserStatus(clientId interface{}, userId []interface{}, status int) (nameList []StringValue, err error) {
	o := orm.NewOrm()
	var param []interface{}
	var sqlCon string
	param = append(param, status)
	if userId != nil {
		sqlCon = "and oauth_user_id in " + common.GenPlaceholder(len(userId))
		param = append(param, userId...)
	} else {
		sqlCon = "and (client_id = ? or client_id like ?)"
		param =append(param, clientId, utils.ParseString(clientId) + "#%")
	}

	_, err = o.Raw("update oauth_user  set status = ?, updated_at = now() where deleted = 0 " + sqlCon, param...).Exec()
	if err != nil {
		logs.Error(err)
		return
	}
	_, err = o.Raw("select name as value1 from oauth_user where deleted = 0 and status = ? " + sqlCon, param...).QueryRows(&nameList)
	return
}
