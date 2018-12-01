package constant


const(
	EUREKA_SERVER = "http://eureka-server.platform:3333/eureka"
	APP_NAME = "oauth"
	CONFIG_SERVER = "configServer"
	PROFILE = "hgh-test"
	BRANCH = "dev"
	ROOT_PATH = "/"
)

const PASSWORD_KEY = "fq398ptvyghq934y94tq34hgt62w45yh3w56qt34q534d4etgq234yt24554ytg3w456yhw4de5frtsrthws46j67k579i8l5yh465"

type LOGIN_STATUS int8

const(
	LOGIN_FIRST LOGIN_STATUS = iota
	LOGIN_FAILD
	LOGIN_SUCCESS
)

type USER_TYPE int8

const(
	USER_BACK USER_TYPE = iota + 1
	USER_WECHAT
)

type AUTHSTATUS int8

const(
	UNGET AUTHSTATUS = iota
	GETED
	UPDATED
)

const(
	STAR = "********"
	NULLNAME = "#T_T"
)