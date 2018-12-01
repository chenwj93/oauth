package utils

import (
	"net/url"
	"io/ioutil"
	"bytes"
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"strings"
	"mime/multipart"
)

func GetDataByHttpGet(url string, param ...map[string]interface{}) (map[string]interface{}, error){
	var paramStr string
	if len(param) > 0 {
		for k, v := range param[0] {
			paramStr += "&" + k + "=" + ParseString(v)
		}
	}
	if len(paramStr) != 0{
		if strings.Contains(url, "?") {
			url += paramStr[1:]
		} else {
			url += "?" + paramStr[1:]
		}
	}
	ULog.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	if len(param) >1 {
		for k,v := range param[1] {
			fmt.Println(k, v)
			req.Header.Set(k, ParseString(v))
		}
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	m := make(map[string]interface{})
	json.Unmarshal(body, &m)
	return m, err
}

func GetDataByHttpPostForm(Url string, m map[string]interface{}) (ret map[string]interface{}, err error) {
	var param = make(url.Values)
	for key, ele := range m {
		//fmt.Println(utils.ParseString(ele))
		param[key] = []string{ParseString(ele)}
	}
	response, err := http.PostForm(Url, param)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &ret)
	return
}

func GetDataByHttpPost(Url string, m ...map[string]interface{}) (ret map[string]interface{}, err error) {
	var byteInfo []byte
	if len(m)>0{
		byteInfo, _ = json.Marshal(m[0])
	}
	fmt.Println(Url, string(byteInfo))
	req, _ := http.NewRequest("POST", Url, bytes.NewReader(byteInfo))
	req.Header.Set("Content-Type", "application/json")
	if len(m) >1 {
		for k,v := range m[1] {
			req.Header.Set(k, ParseString(v))
		}
	}
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &ret)
	return
}

func GetDataByHttpPostFormData(Url string, paramMap, headMap map[string]interface{}) (ret map[string]interface{}, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for k, v := range paramMap{
		w.WriteField(k, ParseString(v))
	}

	w.Close()
	fmt.Println(Url, paramMap)
	req, _ := http.NewRequest("POST", Url, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	for k, v := range headMap{
		req.Header.Set(k, ParseString(v))
	}

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &ret)
	return
}