package conf

import (
	"strconv"
	"fmt"
	"encoding/json"
)

var properties []map[string]string

func checkPropertyIsEmpty() {
	if len(properties) == 0 {
		panic("配置信息加载失败")
	}
}

func PrintAll()  {
	j,_ := json.Marshal(properties)
	fmt.Println(string(j))
}

func Get(key string) string {
	return GetString(key)
}

func GetString(key string) string {
	checkPropertyIsEmpty()
	for _, temp := range properties {
		if val, ok := temp[key]; ok {
			return val
		}
	}
	panic("配置项 " + key + "  不存在")
}

func GetStringW(key, defaultValue string) string {
	if len(properties) == 0 {
		return defaultValue
	}
	for _, temp := range properties {
		if val, ok := temp[key]; ok {
			return val
		}
	}
	return defaultValue
}

func GetInteger(key string) int {
	checkPropertyIsEmpty()
	var strI = GetStringW(key, "0")
	i, err := strconv.Atoi(strI)
	if err != nil {
		panic("不能将 string 转为 int； key is " + key + "; value is " + strI)
	}
	return i
}

func GetIntegerW(key string, defaultValue int) int {
	var strI = GetStringW(key,strconv.Itoa(defaultValue))
	i, err := strconv.Atoi(strI)
	if err != nil {
		panic("不能将 string 转为 int； key is " + key + "; value is " + strI)
	}
	return i
}

func GetFloat(key string) float64 {
	strF := GetString(key)
	f, err := strconv.ParseFloat(strF, 64)
	if err != nil {
		panic("不能将 string 转为 float； key is " + key + "; value is " + strF)
	}
	return f
}

func GetFloatW(key string, defaultValue float64) float64 {
	strF := GetStringW(key, strconv.FormatFloat(defaultValue, 'f', -1, 64))
	f, err := strconv.ParseFloat(strF, 64)
	if err != nil {
		panic("不能将 string 转为 float； key is " + key + "; value is " + strF)
	}
	return f
}

func GetBool(key string) bool {
	strB := GetString(key)
	b, err := strconv.ParseBool(strB)
	if err != nil {
		panic("不能将 string 转为 bool； key is " + key + "; value is " + strB)
	}
	return b
}

func GetBoolW(key string, defaultValue bool) bool {
	strB := GetStringW(key, strconv.FormatBool(defaultValue))
	b, err := strconv.ParseBool(strB)
	if err != nil {
		panic("不能将 string 转为 bool； key is " + key + "; value is " + strB)
	}
	return b
}
