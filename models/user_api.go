package models

import "github.com/astaxie/beego/orm"

type UserApi struct {
	UrlId   uint32 `orm:"column(url_id)"`
	UserId  int64 `orm:"column(user_id)"`
	Deleted int8  `orm:"column(deleted)"`
}

func (u *UserApi) Add(o orm.Ormer) (err error) {
	_, err = o.Raw("insert into user_api (url_id, user_id) values(?, ?)", u.UrlId, u.UserId).Exec()
	return
}

func (u *UserApi) Delete() (err error) {
	o := orm.NewOrm()
	sql := "delete from user_api where user_id = ? "
	param := []interface{}{u.UserId}
	if u.UrlId <= 0 {
		sql += " and url_id = ?"
		param = append(param, u.UrlId)
	}
	_, err = o.Raw(sql, param...).Exec()
	return
}

func (u *UserApi) GetApiList() (apilist []uint32, err error) {
	o := orm.NewOrm()
	var object []UserApi
	sql := "select url_id from user_api where deleted = 0 and user_id = ? order by url_id asc"
	_, err = o.Raw(sql, u.UserId).QueryRows(&object)
	if err == nil {
		apilist = make([]uint32, len(object))
		for i := 0; i < len(object); i++ {
			apilist[i] = object[i].UrlId
		}
	}
	return
}
