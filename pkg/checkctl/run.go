package checkctl

import (
	"inspection/global"
	"inspection/pkg/utils"
	"k8s.io/klog/v2"
)

func Run() {
	url := "http://" + global.ServerAddr + "/api/v1/run-jobs"
	jsonData := []byte{}
	// 创建请求

	if err := utils.Post(url, jsonData); err != nil {
		klog.Info("请求:", url, "失败")
	}

	klog.Info("请求:", url, "成功")

}
