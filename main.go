package main

import (
	"fmt"
	"runtime"
	"time"
	// "flag"
	"net/http"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"cloud-config-client-go/conf"

	"oauth/routers"
	"oauth/controllers"
	"oauth/constant"
	"logs"
	"utils"
	"log"
)

func init() {
	// var PROFILE string
	// flag.StringVar(&PROFILE,"profile","test","execute environment")
	// flag.Parse()
	// fmt.Println(PROFILE)
	conf.Load(constant.EUREKA_SERVER, constant.CONFIG_SERVER, constant.APP_NAME, constant.PROFILE, constant.BRANCH, false)
	//conf.LoadLocalProperties("conf.properties")

	//database connection
	mysql_user := conf.GetString("mysql.username")
	mysql_pass := conf.GetString("mysql.password")
	mysql_urls := conf.GetString("mysql.url")
	mysql_port := conf.GetString("mysql.port")
	mysql_name := conf.GetString("mysql.conn.name")

	// timezone
	orm.DefaultTimeLoc = time.Local
	orm.RegisterDataBase("default", "mysql", mysql_user+":"+mysql_pass+"@tcp("+mysql_urls+":"+mysql_port+")/"+mysql_name+ "?charset=utf8&loc=Asia%2FShanghai")
	orm.Debug = true
	runtime.GOMAXPROCS(runtime.NumCPU())

	machines :=[]string{constant.EUREKA_SERVER}
	utils.StartEureka(constant.APP_NAME, conf.GetInteger("app.port"), nil, machines)

	//logs.InitSingleFile("log/log.log", orm.DebugLog.Logger)
	logs.InitTimeFile("log/saas/oauth/", time.Hour * 24, []*log.Logger{orm.DebugLog.Logger, utils.ULog.Logger}, logs.DEBUG)

	controllers.RootClient = conf.GetStringW("root.client", "adatafun")
	constant.TokenExpiry = time.Hour * time.Duration(conf.GetIntegerW("token.expiry", 8))
	constant.EmployeeApp = conf.GetString("employee.app")
	logs.Info("root client id:", controllers.RootClient)
	logs.Info("token expiry hours:", constant.TokenExpiry)

	indexSlice := controllers.OauthAccessTokensController{}.Init()
	controllers.InitSchedule(indexSlice)
}

// create by  chenwj on 2018-07-06
func main() {
	http.HandleFunc("/", routers.GetHandler())
	fmt.Printf("server at :%s\n", conf.GetString("app.port"))
	http.ListenAndServe(":" + conf.GetString("app.port"), nil)
}