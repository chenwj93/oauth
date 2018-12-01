package models

import (
	"errors"
	"utils"
	"logs"
	"github.com/astaxie/beego/orm"
	"strings"
	"sync"
	"oauth/constant"
	"net/http"
)

var URLMAP URLMap

func init()  {
	URLMAP.mm = make(map[string]uint32)
}

type URLMap struct {
	mm    map[string]uint32
	index uint32
	sync.RWMutex
}

func (u *URLMap) Store(url string, index uint32) {
	u.Lock()
	defer u.Unlock()
	u.mm[url] = index
	if u.index < index {
		u.index = index
	}
}

func (u *URLMap) StoreAuto(url string) uint32 {
	u.Lock()
	defer u.Unlock()
	u.index++
	u.mm[url] = u.index
	return u.index
}

func (u *URLMap) Get(url string) (index uint32, ok bool) {
	u.RLock()
	defer u.RUnlock()
	index, ok = u.mm[url]
	return
}

type Api struct {
	Id      uint32 `orm:"column(id);auto"`
	Url     string `orm:"column(url)"`
	Deleted int8   `orm:"column(deleted)"`
}

func (u *Api) TableName() string {
	return "api"
}

func init() {
	orm.RegisterModel(new(Api))
}

func (u *Api) Add(o orm.Ormer) (uint32, error) {
	id, err := o.Insert(u)
	if err != nil {
		logs.Error(err)
	}
	return uint32(id), err
}

func InsertBatch(authSet *[]interface{}, userId int64) {
	sql := strings.Builder{}
	sql.WriteString("insert into api (url) values ")
	for i := 0; i < len(*authSet); i++ {
		sql.WriteString("(?)")
		if i != len(*authSet)-1 {
			sql.WriteByte(',')
		}
	}
	o := orm.NewOrm()
	o.Raw(sql.String(), *authSet...).Exec()
}

func GetAuthFromRemote(userInfo *UserInfo) (error) {
	logs.Info("service call:")
	// todo 服务间调用
	menusMap, err := utils.ServiceCall(constant.EmployeeApp, "employee/getAccessList", http.MethodGet, map[string]interface{}{"userId": userInfo.Id}, map[string]interface{}{"clientId":"-1"})
	//menusMap, err := utils.GetDataByHttpGet("http://192.168.1.107:9092/employee/getAccessList?userId="+ utils.ParseString(userInfo.Id))
	if err != nil {
		logs.Error(err)
		return err
	}
	if utils.ParseInt(menusMap["status"]) != 200 {
		logs.Error(menusMap["msg"])
		return errors.New(utils.ParseString(menusMap["msg"]))
	}
	dataMap, ok := menusMap["data"].(map[string]interface{})
	if !ok {
		logs.Error(constant.AUTH_DATA_ERR, menusMap)
		return constant.AUTH_DATA_ERR
	}
	urlList, ok := dataMap["apiList"].([]interface{})
	if !ok {
		logs.Error(constant.AUTH_DATA_ERR, menusMap)
		return constant.AUTH_DATA_ERR
	}
	err = handleUrlList(urlList, userInfo)
	if err != nil {
		logs.Error(err)
		return err
	}
	userInfo.UserName = utils.ParseString(dataMap["userName"])
	return nil
}

func handleUrlList(urlList []interface{}, userInfo *UserInfo) (error) {
	o := orm.NewOrm()
	err := (&UserApi{UserId: userInfo.Id}).Delete()
	if err != nil {
		return err
	}
	ifUpdateStatus := true
	userApiList := make([]uint32, len(urlList))
	for i := 0; i < len(urlList); i ++ {
		url := strings.TrimSpace(utils.ParseString(urlList[i]))
		if url == ""{
			continue
		}
		if urlIndex, ok := URLMAP.Get(url); ok {
			userApiList[i] = urlIndex
			SaveUserApi(urlIndex, userInfo.Id, &ifUpdateStatus, o)
		} else {
			id, err := (&Api{Url: url}).Add(o)
			if err != nil || id <= 0 {
				id = URLMAP.StoreAuto(url)
			} else {
				URLMAP.Store(url, id)
				SaveUserApi(id, userInfo.Id, &ifUpdateStatus, o)
			}
			userApiList[i] = id
		}
	}
	if ifUpdateStatus {
		UpdateAuthStatus(userInfo.Id, constant.GETED)
	}
	userInfo.AuthStatus = constant.GETED
	SortUint32(userApiList)
	userInfo.AuthSet = userApiList
	return err
}

func SaveUserApi(urlId uint32, userId int64, ifUpdateStatus *bool, o orm.Ormer) {
	userApi := UserApi{UrlId: urlId, UserId: userId}
	err := userApi.Add(o)
	if err != nil {
		logs.Error(err)
		*ifUpdateStatus = false
	}
	return
}

func GetApiLocal() (apilist []Api, err error){
	o := orm.NewOrm()

	sql := "select url, id from api where deleted = 0 "
	_, err = o.Raw(sql).QueryRows(&apilist)
	return
}