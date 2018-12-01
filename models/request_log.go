package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type RequestLog struct {
	Id       int64  `json:"id" orm:"column(id);pk"`
	Token    string `json:"token" orm:"column(token);size(40);null"`
	Request  string `json:"request" orm:"column(request);size(50);null"`
	Param    string `json:"param" orm:"column(param);size(255);null"`
	Response string `json:"response" orm:"column(response);size(255);null"`
}

func (t *RequestLog) TableName() string {
	return "request_log"
}

func init() {
	orm.RegisterModel(new(RequestLog))
}

func (t *RequestLog) Add() (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(t)
	return
}

func (t *RequestLog) GetById() (err error) {
	o := orm.NewOrm()
	if err = o.Read(t); err == nil {
		return nil
	}
	return err
}

func (t *RequestLog) Update(fields ...string) (err error) {
	o := orm.NewOrm()
	v := RequestLog{Id: t.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(t, fields...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}
