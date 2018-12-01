package core

import (
	"strings"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrorNetwork   = fmt.Errorf("请求url失败")
	ErrorInternalServerError   = fmt.Errorf("网络错误")
	ErrorNotFound  = fmt.Errorf("url不存在")
	ErrorForbidden = fmt.Errorf("url不存在")
	ErrorUnknown   = fmt.Errorf("未知错误")
)
var lastPropertyVersion string

func HttpRequest(urlPath ...string) ([]byte, error) {
	url := urlPath[0]
	if strings.LastIndex(url,"/") == len(url) -1 {
		url = url[:len(url) - 1]
	}
	if strings.Index(url, "http://") == -1 {
		if strings.Index(url, "https://") == -1 {
			url = "http://" + url
		}
	}

	for i := 1; i < len(urlPath); i++ {
		url += "/" + urlPath[i]
	}

	req, err := http.NewRequest(http.MethodGet,url,nil)
	client := new(http.Client)
	respData,err:= client.Do(req)
	if err != nil {
		return nil, ErrorNetwork
	}


	if respData.StatusCode == http.StatusOK {
		propertyBranchInfo, _ := ioutil.ReadAll(respData.Body)
		return propertyBranchInfo, nil
	} else if respData.StatusCode == http.StatusNotFound {
		return nil, ErrorNotFound
	} else if respData.StatusCode == http.StatusForbidden {
		return nil, ErrorForbidden
	}else if respData.StatusCode == http.StatusInternalServerError {
		return nil,ErrorInternalServerError
	}
	return nil, nil
}
