package checkctl

import (
	"encoding/json"
	"fmt"
	"inspection/global"
	"k8s.io/klog/v2"
	"net/url"
	"os"
)

type Scrip struct {
	Name             string `json:"name"`
	ContentJson      string `json:"content_json"`
	ResourceFilePath string `json:"resource_file_path"`
}

func (this *Scrip) Get() error {
	klog.Info("get scrip")
	baseURL := "http://" + global.ServerAddr + "/api/v1/one-check-script"
	// 构建查询参数
	params := url.Values{}
	params.Add("script_name", this.Name)

	// 将查询参数添加到 URL
	urlWithParams := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	if err := Get(urlWithParams); err != nil {
		return err
	}
	return nil
}

func (this *Scrip) Add() error {
	klog.Info("add scrip")
	url := "http://" + global.ServerAddr + "/api/v1/check-script"
	//url := "http://127.0.0.1:8092/api/v1/check-script"

	resourceData, err := os.ReadFile(this.ResourceFilePath)
	if err != nil {
		return err
	}
	this.ContentJson = string(resourceData)
	jsonDataString := fmt.Sprintf("{\"name\":\"%s\",\"content_json\":\"%s\"}", this.Name, this.ContentJson)
	// 定义请求体
	jsonData, err := json.Marshal(&this)
	if err != nil {
		return err
	}
	// 创建请求

	if err := Post(url, jsonData); err != nil {
		return err
	}

	klog.Info("添加内容:", jsonDataString, "成功")

	return nil
}

func (this *Scrip) Update() error {
	klog.Info("update scrip")

	url := "http://" + global.ServerAddr + "/api/v1/check-script"
	//url := "http://127.0.0.1:8092/api/v1/check-script"

	resourceData, err := os.ReadFile(this.ResourceFilePath)
	if err != nil {
		return err
	}
	this.ContentJson = string(resourceData)
	jsonDataString := fmt.Sprintf("{\"name\":\"%s\",\"content_json\":\"%s\"}", this.Name, this.ContentJson)

	// 定义请求体
	jsonData, err := json.Marshal(&this)
	if err != nil {
		return err
	}
	// 创建请求

	if err := Put(url, jsonData); err != nil {
		return err
	}

	klog.Info("添加内容:", jsonDataString, "成功")

	return nil
}

func (this *Scrip) Delete() error {
	klog.Info("delete scrip")
	return nil
}
