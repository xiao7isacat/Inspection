package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"inspection/models"
	"inspection/pkg/common"
	"k8s.io/klog/v2"
)

func CheckJobAdd(c *gin.Context) {
	var (
		check_script   models.CheckScript
		desired_result models.DesiredResult
	)

	var input models.CheckJob
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}

	// 检查 script_name 和 desired_result_name
	check_script.Name = input.ScriptName
	desired_result.Name = input.DesiredName
	if err := check_script.GetOne(); err != nil {
		errInfo := fmt.Errorf("check_script %v", err)
		common.JSONR(c, 400, errInfo)
		return
	}

	if err := desired_result.GetOne(); err != nil {
		errInfo := fmt.Errorf("desired_result %v", err)
		common.JSONR(c, 400, errInfo)
		return
	}

	// TODO 校验cmdb配置 你自己去适配
	// 将ip list解析一下

	if input.IpJson != "" {
		var ips []string
		err := json.Unmarshal([]byte(input.IpJson), &ips)
		if err != nil {
			common.JSONR(c, 500, err)
			return
		}
		input.AllNum = int64(len(ips))
	}

	id, err := input.CreateOne()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}
	klog.Infof("[CheckJobAdd.success][CheckJob:%v]", input.Name)
	common.JSONR(c, 200, id)
}

/*func CronJobAdd(c *gin.Context) {

	var input models.CronJob
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	if input.RunAtHourMinute == "" {
		common.JSONR(c, 400, "RunAtHourMinute.empty")
		return
	}

	// 检查 script_name 和 desired_result_name
	script, err := models.CheckScriptGetByName(input.ScriptName)
	if script == nil {
		err = fmt.Errorf("script :%v not found in db: %w", input.ScriptName, err)

		common.JSONR(c, 400, err)
		return
	}

	result, err := models.DesiredResultGetByName(input.DesiredResultName)
	if result == nil {
		err = fmt.Errorf("result :%v not found in db: %w", input.DesiredResultName, err)

		common.JSONR(c, 400, err)
		return
	}
	cm := c.MustGet(common.CheckJobManager).(*check.CheckJobManager)
	// 校验cmdb配置
	// 将ip list解析一下

	if input.IpJson != "" {
		var ips []string
		err = json.Unmarshal([]byte(input.IpJson), &ips)
		if err != nil {
			common.JSONR(c, 500, err)
			return
		}
		input.AllNum = int64(len(ips))
	}

	if input.JobWaitCompleteMinutes == 0 {
		input.JobWaitCompleteMinutes = cm.Cg.JobCompleteMinutes
	}

	id, err := input.AddOne()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}
	klog.Infof("[CronJobAdd.success][CronJob:%v]", input.Name)
	common.JSONR(c, 200, id)
}*/

func CheckJobGets(c *gin.Context) {
	var checkJob models.CheckJob
	ljs, err := checkJob.GetList(false)
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, 200, ljs)
}
