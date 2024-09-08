package checkctl

import (
	"encoding/json"
	"fmt"
	"inspection/global"
	"inspection/pkg/utils"
	"k8s.io/klog/v2"
	"net/url"
	"os"
)

type Desired struct {
	Name             string `json:"name"`
	ResultJson       string `json:"result_json"`
	ResourceFilePath string `json:"resource_file_path"`
}

func (this *Desired) Get() error {
	klog.Info("get desired")
	baseURL := "http://" + global.ServerAddr + "/api/v1/one-desired-result"
	// 构建查询参数
	params := url.Values{}
	params.Add("result_name", this.Name)

	// 将查询参数添加到 URL
	urlWithParams := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	if _, err := utils.Get(urlWithParams); err != nil {
		return err
	}

	return nil
}

func (this *Desired) Add() error {
	klog.Info("add desired")
	url := "http://" + global.ServerAddr + "/api/v1/desired-result"

	resourceData, err := os.ReadFile(this.ResourceFilePath)
	if err != nil {
		return err
	}
	this.ResultJson = string(resourceData)
	jsonDataString := fmt.Sprintf("{\"name\":\"%s\",\"result_json\":\"%s\"}", this.Name, this.ResultJson)
	// 定义请求体
	jsonData, err := json.Marshal(&this)
	if err != nil {
		return err
	}
	// 创建请求

	if err := utils.Post(url, jsonData); err != nil {
		return err
	}

	klog.Info("添加内容:", jsonDataString, "成功")

	return nil
}

func (this *Desired) Update() error {
	klog.Info("update desired")
	url := "http://" + global.ServerAddr + "/api/v1/desired-result"
	//url := "http://127.0.0.1:8092/api/v1/check-script"

	resourceData, err := os.ReadFile(this.ResourceFilePath)
	if err != nil {
		return err
	}
	this.ResultJson = string(resourceData)
	jsonDataString := fmt.Sprintf("{\"name\":\"%s\",\"result_json\":\"%s\"}", this.Name, this.ResultJson)

	// 定义请求体
	jsonData, err := json.Marshal(&this)
	if err != nil {
		return err
	}
	// 创建请求

	if err := utils.Put(url, jsonData); err != nil {
		return err
	}

	klog.Info("添加内容:", jsonDataString, "成功")

	return nil
}

func (this *Desired) Delete() error {
	klog.Info("delete desired")
	return nil
}
