package routers

import (
	"context"
	"net/http"

	"utils"
	"logs"

	"oauth/constant"
	"oauth/controllers"
	"encoding/json"
	"oauth/common"
	"runtime/debug"
)

type BusinessServiceImpl struct {
}

//var route utils.Router
//
//func init(){
//	route.AddRouters(new(controllers.DemoController),
//		)
//	//route.AddRouter("test", new(controllers.TestController))
//	route.PrintURL()
//}

// create by chenwj on 2018-07-06 2018-03-07
func (msi *BusinessServiceImpl) Handle(operation string, paramInput map[string]interface{}) (r *utils.Response, err error) {
	defer func() {
		if re := recover(); re != nil {
			logs.Error(re, string(debug.Stack()))
			r, err = utils.ConcatFailed(nil)
		}
	}()

	logs.Info("-->param:", utils.SimplifyStringMap(common.HideAndCopy(paramInput), 300))

	//r, err = route.URLMapping(operation, paramInput)
	token := paramInput["access_token"]
	delete(paramInput, "access_token")
	ctx := context.WithValue(context.Background(), "token", token)
	if operation == "register" {
		ctx = context.WithValue(ctx, "operate", "register")
	}
	var con controllers.Accesser = &controllers.Access{ctx}

	switch operation {
	case "register":
		j, _ := json.Marshal(paramInput)
		r, err = con.Register(j)
	case "login":
		r, err = con.Login(paramInput)
	case "authenticate":
		r, err = con.Authenticate(paramInput)
	case "modifyPassword":
		r, err = con.ModifyPassword(paramInput)
	case "logout":
		r, err = con.Logout()
	case "authHadModified":
		r, err = con.AuthHadModified(paramInput)
	case "changeUserStatus":
		r, err = con.ChangeUserStatus(paramInput)
	default:
		r, err = utils.ConcatNotFound()
	}

	go controllers.SaveRequestLog(operation, paramInput, r, utils.ParseString(token))

	return
}

func GetHandler() http.HandlerFunc {
	f := func() utils.RouterInterface {
		return &BusinessServiceImpl{}
	}
	handler := utils.Handler{constant.ROOT_PATH, nil, f}
	return handler.Handle
}
