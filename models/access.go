package models

import (
	"time"
	"sync"
	"utils"
	"oauth/constant"
	"strings"
	"logs"
)

var (
	Lim LoginInfoMap
	Uim UserInfoMap
)

type UserInfo struct {
	OauthUser
	AuthSet        []uint32
	LastLoginToken string
	UserName 		string
}

func (u *UserInfo) GetUserAuth(){
	var err error
	if u.AuthStatus != constant.GETED { // 从角色模块获取权限
		err = GetAuthFromRemote(u)
	} else {
		u.AuthSet, err = (&UserApi{UserId:u.Id}).GetApiList()
	}
	if err != nil{
		logs.Error(err)
	}
	return
}

func (u *UserInfo) CheckAuth(operate string) bool {
	logs.Info("操作", operate)
	index, ok := URLMAP.Get(operate)
	if !ok {
		logs.Warn("权限不足")
		return false
	}
	for i, j := 0, len(u.AuthSet) -1; i <= j; {
		m := (i + j)/ 2
		if u.AuthSet[m] == index {
			return true
		}
		if u.AuthSet[m] > index {
			j = m - 1
		}else {
			i = m + 1
		}
	}
	return false
}

type UserInfoMap struct {
	sync.Map
}

func (u *UserInfoMap) Get(abstract string) (userInfo *UserInfo, ok bool) {
	val, ok := u.Load(abstract)
	if ok {
		userInfo, ok = val.(*UserInfo)
	}
	return
}

func (u *UserInfoMap) Put(abstract string, val *UserInfo) {
	u.Store(abstract, val)
}

type LoginInfo struct {
	ExpireTime time.Time `json:"-"`
	User       *UserInfo
}

type LoginInfoMap struct {
	sync.Map
}

func (u *LoginInfoMap) Get(token string) (loginInfo *LoginInfo, ok bool) {
	val, ok := u.Load(token)
	if ok {
		loginInfo, ok = val.(*LoginInfo)
	}
	return
}

func (u *LoginInfoMap) Put(token string, val *LoginInfo) {
	u.Store(token, val)
}

// return param status: 1-未找到 2- 密码错误 3-成功
func CheckDuplicateLogin(userName, password interface{}) (abs, token string, userId int64, status constant.LOGIN_STATUS, clientId string) {
	uis := utils.ParseString(userName)
	token = GenerateToken(uis)
	abs = utils.Md5(uis)
	userInfoBack, ok := Uim.Get(abs)
	if !ok {
		return abs, token, -1, constant.LOGIN_FIRST, ""
	}
	password = PasswordEncypt(password)
	if userInfoBack.Status != 0 || password != userInfoBack.Password {
		return utils.EMPTY_STRING, utils.EMPTY_STRING, -1, constant.LOGIN_FAILD, ""
	}

	if userInfoBack.Singleton == 1 { // 禁止重复登录
		oat := OauthAccessTokens{Id: userInfoBack.LastLoginToken, Status: 1, UpdatedAt: time.Now()}
		oat.Update("Status", "UpdatedAt")
		Lim.Delete(userInfoBack.LastLoginToken) // 删除前一个登录信息
	}
	loginInfo := &LoginInfo{time.Now().Add(constant.TokenExpiry), userInfoBack}
	userInfoBack.LastLoginToken = token
	Lim.Put(token, loginInfo)
	return abs, token, userInfoBack.Id, constant.LOGIN_SUCCESS, userInfoBack.ClientId
}

func GenerateToken(uis string) string {
	var sb strings.Builder
	t1 := utils.Md5(uis + utils.ParseString(time.Now().UnixNano()))
	t2 := utils.GenerateCahrCode8()
	for i, v := range t2 {
		sb.WriteString(t1[i*4: i*4+4])
		sb.WriteByte(byte(v))
	}
	return sb.String()
}

func PasswordEncypt(password interface{}) string {
	return utils.Md5(utils.ParseString(password) + constant.PASSWORD_KEY)
}

func StoreLoginUser(tokenList []OauthAccessTokens, userList []OauthUser) ( *IndexSlice){
	var indexSlice IndexSlice
	ui := 0
	ifCanStore := false //登录信息是否存在匹配用户信息，是则可存储
	var backUser *UserInfo
	var backAbs string
	for i := 0; i < len(tokenList); i++ {
		if backUser != nil &&backUser.Id == tokenList[i].UserId && backUser.Singleton == 1{ // 禁止多人登录
			oat := OauthAccessTokens{Id:tokenList[i].Id, Status: 1, UpdatedAt: time.Now()}
			go oat.Update("Status", "UpdatedAt")
			continue
		}
		if backUser != nil &&backUser.Id != tokenList[i].UserId{ // token用户切换, 尚未匹配user，所以不可存储
			ifCanStore = false
		}
		for ui < len(userList) && userList[ui].Id < tokenList[i].UserId  { // 调整userlist的下标，至user与token的user一致
			ui++
		}
		if !ifCanStore && ui < len(userList) && tokenList[i].UserId == userList[ui].Id { // user匹配一致，存储userinfo,并备份其指针
			userInfo := &UserInfo{}
			utils.Assignment(userList[ui], userInfo)
			userInfo.LastLoginToken = tokenList[i].Id
			Uim.Put(userInfo.GetAbstract(), userInfo)
			ifCanStore = true
			backUser = userInfo
			backAbs = backUser.GetAbstract()
		}
		if ifCanStore{ // 存储登录信息，摘要为前面备份的userInfo指针
			loginInfo := &LoginInfo{tokenList[i].ExpireTime, backUser}
			Lim.Put(tokenList[i].Id, loginInfo)
			indexSlice = append(indexSlice, &Index{backAbs, tokenList[i].Id, tokenList[i].ExpireTime})
		}
	}
	return &indexSlice
}

