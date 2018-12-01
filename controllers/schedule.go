package controllers

import (
	"time"
	"oauth/models"
	"container/list"
	"sync"
	"sort"
	"logs"
	"oauth/constant"
)

type ListSync struct {
	ll    list.List
	mutex sync.Mutex
}

func (l *ListSync) PushBack(index *models.Index) {
	l.mutex.Lock()
	l.ll.PushBack(index)
	l.mutex.Unlock()
}

func (l *ListSync) PopHead() (element *list.Element, index *models.Index) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	element = l.ll.Front()
	if element != nil {
		index, _ = element.Value.(*models.Index)
	}
	return
}

func (l *ListSync) RemoveHead(element *list.Element) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.ll.Remove(element)
	return
}

var ScheduleList ListSync

func InitSchedule(indexList *models.IndexSlice) {
	sort.Sort(*indexList)
	for i := 0; i < indexList.Len();  i++ {
		ScheduleList.ll.PushBack((*indexList)[i])
	}

	expireMax := constant.TokenExpiry
	go ScheduleLim(expireMax)
}

func ScheduleLim(expireMax time.Duration) {
	defer protect(expireMax)
	tick := make(<-chan time.Time)
	var index *models.Index
	var ele *list.Element
	for {
		if index == nil{
			ele, index = ScheduleList.PopHead()
		}
		for index != nil && index.Expire.Before(time.Now()) {
			ScheduleList.RemoveHead(ele)
			go ClearLoginInfo(index.Abs, index.Token)
			index = nil
			ele, index = ScheduleList.PopHead()
		}
		if index != nil{
			tick = time.Tick(index.Expire.Sub(time.Now()))
			logs.Info("next: ", index.Expire.Sub(time.Now()))
		} else {
			tick = time.Tick(expireMax)
			logs.Info("无登录信息...")
		}
		<-tick
	}
}


func protect(expireMax time.Duration){
	if r := recover(); r != nil{
		logs.Error(r)
	}
	go ScheduleLim(expireMax)
	return
}

func ClearLoginInfo(abs, token string) {
	logs.Info("删除：", abs, token)
	go OauthAccessTokensController{}.Expire(token)
	go models.Lim.Delete(token)
	if userInfo, ok := models.Uim.Get(abs); ok && userInfo.LastLoginToken == token{
		go models.Uim.Delete(abs)
	}
	return
}
