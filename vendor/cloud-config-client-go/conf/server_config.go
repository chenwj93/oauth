package conf

import (
	"encoding/json"
	"cloud-config-client-go/core"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type serverConfig struct {
	eurekaServer     string
	serviceId        string
	serverUri        string
	applicationName  string
	profile          string
	label            string
	enableRemoteConf bool
}

type configProperties struct {
	Name            string   `json:"name"`
	Profiles        []string `json:"profiles"`
	Label           string   `json:"label"`
	Version         string   `json:"version"`
	State           string   `json:"state"`
	PropertySources []source `json:"propertySources"`
}

type source struct {
	Name   string            `json:"name"`
	Source map[string]string `json:"source"`
}

type Instance struct {
	InstanceId  string `xml:"instanceId"`
	HostName    string `xml:"hostName"`
	HomePageUrl string `xml:"homePageUrl"`
	IpAddr      string `xml:"ipAddr"`
	Port        *Port  `xml:"port"`
}

type Port struct {
	Port    int  `xml:",chardata" json:"$"`
	Enabled bool `xml:"enabled,attr" json:"@enabled"`
}

type ServiceInfo struct {
	Application xml.Name   `xml:"application"`
	Name        string     `xml:"name"`
	Instance    []Instance `xml:"instance"`
}

var instance *serverConfig

func LoadLocalProperties(configs ... string) {
	var configMap map[string]string
	var err error
	for _, config := range configs {
		if strings.HasSuffix(config, ".properties") {
			configMap, err = core.ParsePropertiesFile(config)
		} else if strings.HasSuffix(config, "yml") || strings.HasSuffix(config, "yaml") {
			configMap, err = core.ParseYmlFile(config)
		}
		if err != nil {
			continue
		}
		properties = append(properties, configMap)
	}

}

func Load(eurekaServer, serviceId, applicationName, profile, label string, enableRemoteConfig bool) {
	if instance != nil {
		panic("配置中心只能配置一次")
	}
	instance = &serverConfig{
		eurekaServer, serviceId, "", applicationName, profile, label, enableRemoteConfig,
	}

	discoverConfigServerAndLoad(instance)
}

func discoverConfigServerAndLoad(instance *serverConfig) {
	instanceInfo, err := core.HttpRequest(instance.eurekaServer, "apps", instance.serviceId)
	if err != nil {
		panic("服务：" + instance.serviceId + " 不存在")
	}
	var instanceInfoMap = new(ServiceInfo)
	err = xml.Unmarshal(instanceInfo, instanceInfoMap)
	if err != nil || len(instanceInfoMap.Instance) == 0 {
		fmt.Println("err:", err)
		panic("服务：" + instance.serviceId + " 不可用")
	}
	instanceId := instanceInfoMap.Instance[0].IpAddr + ":" + strconv.Itoa(instanceInfoMap.Instance[0].Port.Port)
	//instanceId := instanceInfoMap.Instance[0].HomePageUrl
	instance.serverUri = instanceId
	loadProperties()
}

func LoadConf(eurekaServer, serviceId, serverUri, applicationName, profile, label string, enableRemoteConfig bool) {

	if instance != nil {
		panic("配置中心只能配置一次")
	}
	instance = &serverConfig{
		eurekaServer, serviceId, serverUri, applicationName, profile, label, enableRemoteConfig,
	}

	defer func() {
		if err := recover(); err != nil {
			loadProperties()
		}
	}()
	discoverConfigServerAndLoad(instance)
}

func LoadConfig(serverUri, applicationName, profile, label string, enableRemoteConfig bool) {
	//once.Do(func() {})
	if instance != nil {
		panic("配置中心只能配置一次")
	}
	instance = &serverConfig{
		"", "", serverUri, applicationName, profile, label, enableRemoteConfig,
	}
	loadProperties()
}

func loadProperties() {
	if instance.enableRemoteConf || properties == nil {
		configData, err := core.HttpRequest(instance.serverUri, instance.applicationName, instance.profile, instance.label)
		if err != nil {
			log.Println("加载配置失败:", err)
			return
		}
		properties = getPropertySource(configData)
	}
}

/*
{"name":"user","profiles":["test"],"label":"master","version":"c99730a081455abb136a845cf419c70d9f957a67","state":null,"propertySources":[{"name":"git@code.aliyun.com:fs-platform/cloud-config.git/config-files/user/user-test.properties","source":{"mysql.url":"localhost:7777"}},{"name":"git@code.aliyun.com:fs-platform/cloud-config.git/config-files/user/user-test.yaml","source":{"mysql.url":"localhost:8888"}}]}
 */
func getPropertySource(srcByte []byte) []map[string]string {
	var propertyInfo configProperties
	json.Unmarshal(srcByte, &propertyInfo)
	var sourceMaps []map[string]string
	var propertySourcesMap = propertyInfo.PropertySources
	if len(propertySourcesMap) > 0 {
		for _, item := range propertySourcesMap {
			sourceMaps = append(sourceMaps, item.Source)
		}
	}
	return sourceMaps
}
