package checkctl

import (
	"encoding/json"
	"fmt"
	"inspection/global"
	"inspection/pkg/utils"
	"k8s.io/klog/v2"
)

type Job struct {
	Name        string `json:"name"`
	ScriptName  string `json:"script_name"`
	DesiredName string `json:"desired_name"`
	IpString    string `json:"ip_string"`
}

func (this *Job) Get() error {
	klog.Info("get job")
	return nil
}

func (this *Job) Add() error {
	klog.Info("add job")
	url := "http://" + global.ServerAddr + "/api/v1/check-job"
	this.DesiredName = this.Name
	this.ScriptName = this.Name
	jsonDataString := fmt.Sprintf("{\"name\":\"%s\",\"script_name\":\"%s\",\"desired_name\":\"%s\",\"ip_string\":\"%s\"}", this.Name, this.ScriptName, this.DesiredName, this.IpString)
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

func (this *Job) Update() error {
	klog.Info("update job")
	return nil
}

func (this *Job) Delete() error {
	klog.Info("delete job")
	return nil
}
