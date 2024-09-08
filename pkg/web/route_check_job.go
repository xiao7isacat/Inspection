package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inspection/global"
	"inspection/models"
	"inspection/pkg/check"
	"inspection/pkg/response"
	"k8s.io/klog/v2"
	"strings"
)

func CheckJobAdd(c *gin.Context) {
	var (
		check_script   models.CheckScript
		desired_result models.DesiredResult
		//ip_addr_info   models.IpAddrInfo
	)

	var checkJob models.CheckJob
	if err := c.BindJSON(&checkJob); err != nil {
		response.JSONR(c, 400, err)
		return
	}

	// 检查 script_name 和 desired_result_name
	check_script.Name = checkJob.ScriptName
	desired_result.Name = checkJob.DesiredName
	if err := check_script.GetOne(); err != nil {
		errInfo := fmt.Errorf("check_script %v", err)
		response.JSONR(c, 400, errInfo)
		return
	}

	if err := desired_result.GetOne(); err != nil {
		errInfo := fmt.Errorf("desired_result %v", err)
		response.JSONR(c, 400, errInfo)
		return
	}

	// TODO 校验cmdb配置 你自己去适配
	// 将ip list解析一下

	if checkJob.IpString == "" {
		klog.Infof("[CheckJob.IpString][is nill]", checkJob.Name)
		err := fmt.Errorf("job %v ip is nil", checkJob.Name)
		response.JSONR(c, 500, err)
		return
	}

	ips := strings.Split(checkJob.IpString, ",")
	checkJob.AllNum = int64(len(ips))

	cm, _ := c.MustGet(global.CheckJobManager).(*check.CheckJobManager)

	if checkJob.JobWaitCompleteMinutes == 0 {
		checkJob.JobWaitCompleteMinutes = cm.Cg.JobCompleteMinutes
	}

	/*for _, ip := range ips {
		ip_addr_info.Ip = ip
		ipExist, err := ip_addr_info.CheckExist()
		if err != nil {
			if err != nil {
				response.JSONR(c, 500, err)
				return
			}
		}
		if !ipExist {
			ip_addr_info.Ip = ip
			_, err = ip_addr_info.CreateOne()
			if err != nil {
				response.JSONR(c, 500, err)
				return
			}
		}
	}*/

	id, err := checkJob.CreateOne()
	if err != nil {
		response.JSONR(c, 500, err)
		return
	}
	klog.Infof("[CheckJobAdd.success][CheckJob:%v]", checkJob.Name)
	response.JSONR(c, 200, id)
}

func CheckJobGets(c *gin.Context) {
	var checkJob models.CheckJob
	ljs, err := checkJob.GetList()
	if err != nil {
		response.JSONR(c, 500, err)
		return
	}

	response.JSONR(c, 200, ljs)
}
