package core

import (
	"os"
	"bufio"
	"io"
	"strings"
	"cloud-config-client-go/commom"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ParsePropertiesFile(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var configMap = make(map[string]string)
	reader := bufio.NewReader(file)
	for {
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		line := string(lineBytes)
		if len(lineBytes) == 0 || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		if !strings.Contains(line, "=") {
			panic("解析配置文件失败；配置文件:" + fileName + "；配置项:" + line)
		}
		//line = strings.TrimSpace(line[:strings.Index(line, "//")])
		lineConfig := strings.Split(line, "=")
		configMap[strings.TrimSpace(lineConfig[0])] = strings.TrimSpace(lineConfig[1])
	}
	return configMap, err
}

func ParseYmlFile(fileName string) (map[string]string, error) {
	fileByte, err := ioutil.ReadFile(fileName)
	var configMap = make(map[string]interface{})
	yaml.Unmarshal(fileByte, &configMap)
	if err != nil {
		panic(err)
	}
	var configs = make(map[string]string)
	for key, value := range configMap {
		resMap := recursion(value)
		for k, v := range resMap {
			configs[key+"."+k] = common.ParseString(v)
		}
	}
	return configs, err
}

func recursion(item interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	interfaceMap, ok := item.(map[interface{}]interface{})
	if !ok {
		res[""] = item
		return res
	}
	var key string
	for k, v := range interfaceMap {
		recRes := recursion(v)
		for nk, nv := range recRes {
			key = common.ParseString(k)
			if len(nk) > 0 {
				key += "." + nk
			}
			res[key] = nv
		}
	}

	return res

}
