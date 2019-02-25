package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

//redis 实例
var Client *redis.Client

type Lock struct {
	key      string
	overTime int
	lockV    int64
}

func NewRedisLock(key string, overTime int) *Lock {
	return &Lock{key: key, overTime: overTime}
}

//create by cwj on 2017-11-11
//get lock from redis
func (l *Lock) Get() {
	if l.key == "" || Client == nil {
		panic("lock undefined")
	}
	for {
		lockV := time.Now().UnixNano()

		script := "if redis.call('get', KEYS[1])  == false then redis.call('set', KEYS[1], ARGV[1]); redis.call('expire', KEYS[1], 5) return 'OK' else return redis.call('get', KEYS[1]) end"
		ret := Client.Eval(script, []string{l.key}, lockV)
		ULog.Println("redis setnx ret:", ret.String())
		if ret != nil {
			if ParseString(ret.Val()) == "OK" { //获取锁成功，直接返回
				l.lockV = lockV
				return
			} else if preLock := ParseInt64(ret.Val()); preLock > 0 {
				if (lockV - preLock) > int64(l.overTime)*int64(time.Second) { //前一个获取锁进程处理时间已经超时
					script = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('set', KEYS[1], ARGV[2]) else return -1 end"
					ret := Client.Eval(script, []string{l.key}, preLock, lockV)

					if ret.Val() == "OK" { //说明获取锁成功
						//*wait <- true
						l.lockV = lockV
						return
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}

//create by cwj on 2017-11-11
//release lock from redis
func (l *Lock) Release() {
	if l.key == "" || Client == nil {
		panic("lock undefined")
	}
	//Client.Del(key)
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	Client.Eval(script, []string{l.key}, l.lockV)
	//fmt.Println(cmd.Result())
	//<-*wait
}

func CheckBusy(key string, timeLine int64) bool {
	interval := time.Second * time.Duration(timeLine)
	ULog.Println(key, interval)
	ok := Client.SetNX(key, 1, interval)
	ULog.Println(ok)
	if ok.Err() != nil {
		return false
	}
	return !ok.Val()
}

func GetSerialNo(key string, expiry time.Time, length int) string {
	ULog.Println("key:", key)
	no := Client.Incr(key)
	if no.Val()%10 == 1 {
		Client.ExpireAt(key, expiry)
	}
	return fmt.Sprintf("%."+strconv.Itoa(length)+"d", no.Val())
}

func GetDateSerialNo(key string, length int) string {
	no := GetSerialNo(TodayWithout()+key, GetToday24(), length)
	ULog.Println("redis:", no)
	return no
}
