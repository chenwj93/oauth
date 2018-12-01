package common

import (
	"regexp"
	"oauth/constant"
)

var reg =`password|passwd|pwd|Password|Passwd|Pwd`

var R *regexp.Regexp

func init() {
	R = regexp.MustCompile(reg)
}

func HidePassword(param map[string]interface{})  {
	for k := range param {
		if R.MatchString(k) {
			param[k] = constant.STAR
		}
	}
}

func HideAndCopy(param map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for k := range param {
		if R.MatchString(k) {
			m[k] = constant.STAR
		}else {
			m[k] = param[k]
		}
	}
	return m
}