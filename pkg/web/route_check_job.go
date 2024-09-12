package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func WebStatusJobs(c *gin.Context) {
	var (
		allJobDone       bool
		jobStatus        check.JobStatus
		jobStatusList    []check.JobStatus
		checkJob         models.CheckJob
		checkJobs        []*models.CheckJob
		failedNodeResult models.FailedNodeResult
		desiredResult    models.DesiredResult
	)

	//获取任务列表
	checkJobs, err := checkJob.GetList()
	if err != nil {
		klog.Info("获取任务列表失败", err)
	}

	notSyncCheckJobs, err := checkJob.GetNotSyncList()
	if err != nil {
		klog.Info("获取未下发任务列表失败", err)
	}

	notComcheckJobs, err := checkJob.GetNotCompleteList()
	if err != nil {
		klog.Info("获取未完成任务列表失败", err)
	}

	if len(notComcheckJobs) == 0 && len(notSyncCheckJobs) == 0 {
		allJobDone = true
	}

	for _, job := range checkJobs {
		//获取任务和基线
		jobStatus.JobName = job.Name
		jobStatus.AllDone = allJobDone
		desiredResult.Name = job.DesiredName
		if err = desiredResult.GetOne(); err != nil {
			klog.Info("获取任务", desiredResult.Name, "基线失败", err)
		}
		desiredResultMap := map[string]string{}
		actualResultMap := map[string]string{}
		if err = json.Unmarshal([]byte(desiredResult.ResultJson), &desiredResultMap); err != nil {
			klog.Info("desiredResult.ResultJson.json.Unmarshal.err:%v", err)
		}
		for dk, dv := range desiredResultMap {
			jobStatus.ExpectValue = dv
			jobStatus.CheckName = dk
			//job未完成时获取ip，设置状态
			if job.JobHasSynced == 0 || job.JobHasComplete == 0 {
				jobStatus.Status = "Waiting"
				jobStatus.ActualValue = ""
				for _, ipList := range strings.Split(job.IpString, ",") {
					jobStatus.Node = ipList
					jobStatusList = append(jobStatusList, jobStatus)
				}
			} else {
				//job完成时获取结果
				failedNodeResult.JobId = int64(job.ID)
				JobNodeResults, err := failedNodeResult.GetList()
				if err != nil {
					//未获取到执行结果时，任务节点丢失
					if err == gorm.ErrRecordNotFound {
						jobStatus.Status = "Miss"
						jobStatus.ActualValue = ""
						for _, ipList := range strings.Split(job.IpString, ",") {
							jobStatus.Node = ipList
							jobStatusList = append(jobStatusList, jobStatus)
						}
					} else {
						klog.Info("获取任务", desiredResult.Name, "执行结果失败", err)
					}

				} else {
					//获取到执行结果时
					for _, JobNodeResult := range JobNodeResults {
						if JobNodeResult.FinalSucceed == 1 {
							jobStatus.Status = "Success"
						} else {
							jobStatus.Status = "Failed"
						}
						if err = json.Unmarshal([]byte(JobNodeResult.ResultJson), &actualResultMap); err != nil {
							klog.Info("failedNodeResult.ResultJson.json.Unmarshal.err:%v", err)
						}
						jobStatus.ActualValue = actualResultMap[dk]
						jobStatus.Node = JobNodeResult.NodeIp
						jobStatusList = append(jobStatusList, jobStatus)
					}
				}

			}

		}
	}
	response.JSONR(c, 200, jobStatusList)

}

func CtlStatusJobs(c *gin.Context) {
	var (
		allJobDone    bool
		jobStatus     check.JobStatus
		jobStatusList []check.JobStatus
		checkJob      models.CheckJob
		checkJobs     []*models.CheckJob
	)

	//获取任务列表
	checkJobs, err := checkJob.GetList()
	if err != nil {
		klog.Info("获取任务列表失败", err)
	}
	notSyncCheckJobs, err := checkJob.GetNotSyncList()
	if err != nil {
		klog.Info("获取未下发任务列表失败", err)
	}

	notComcheckJobs, err := checkJob.GetNotCompleteList()
	if err != nil {
		klog.Info("获取未完成任务列表失败", err)
	}

	if len(notComcheckJobs) == 0 && len(notSyncCheckJobs) == 0 {
		allJobDone = true
	}

	for _, job := range checkJobs {
		var (
			failedNodeResult models.FailedNodeResult
			desiredResult    models.DesiredResult
		)
		//获取任务和基线
		jobStatus.JobName = job.Name
		jobStatus.AllDone = allJobDone
		desiredResult.Name = job.DesiredName
		if err = desiredResult.GetOne(); err != nil {
			klog.Info("命令行获取任务", desiredResult.Name, "基线失败", err)
		}
		desiredResultMap := map[string]string{}
		actualResultMap := map[string]string{}
		if err = json.Unmarshal([]byte(desiredResult.ResultJson), &desiredResultMap); err != nil {
			klog.Info("desiredResult.ResultJson.json.Unmarshal.err:%v", err)
		}
		for dk, dv := range desiredResultMap {
			jobStatus.ExpectValue = dv
			jobStatus.CheckName = dk
			//job未下发时获取ip，设置状态
			if job.JobHasSynced == 0 {
				jobStatus.Status = "Waiting"
				jobStatus.ActualValue = ""
				for _, ipList := range strings.Split(job.IpString, ",") {
					jobStatus.Node = ipList
					jobStatusList = append(jobStatusList, jobStatus)
				}
				continue
			}
			//job未执行时获取状态
			if job.JobHasComplete == 0 {
				jobStatus.Status = "Running"
				jobStatus.ActualValue = ""
				for _, ipList := range strings.Split(job.IpString, ",") {
					jobStatus.Node = ipList
					jobStatusList = append(jobStatusList, jobStatus)
				}
				continue
			}

			//job完成时获取结果
			failedNodeResult.JobId = int64(job.ID)
			JobNodeResults, err := failedNodeResult.GetList()
			if err != nil {
				klog.Info("获取任务", desiredResult.Name, "执行结果失败", err)
			}
			//未获取到执行结果时，任务节点丢失
			if len(JobNodeResults) == 0 {
				jobStatus.Status = "Miss"
				jobStatus.ActualValue = ""
				for _, ipList := range strings.Split(job.IpString, ",") {
					jobStatus.Node = ipList
					jobStatusList = append(jobStatusList, jobStatus)
				}
				continue
			}

			//获取到执行结果时
			for _, JobNodeResult := range JobNodeResults {
				if JobNodeResult.FinalSucceed == 1 {
					jobStatus.Status = "Success"
				} else {
					jobStatus.Status = "Failed"
				}
				if err = json.Unmarshal([]byte(JobNodeResult.ResultJson), &actualResultMap); err != nil {
					klog.Info("failedNodeResult.ResultJson.json.Unmarshal.err:%v", err)
				}
				jobStatus.ActualValue = actualResultMap[dk]
				jobStatus.Node = JobNodeResult.NodeIp
				jobStatusList = append(jobStatusList, jobStatus)
			}

		}

	}

	response.JSONR(c, 200, jobStatusList)

}

func RunJobs(c *gin.Context) {
	var checkJob models.CheckJob
	checkJobs, err := checkJob.GetSyncList()
	if err != nil {
		response.JSONR(c, 500, err)
	}

	for _, job := range checkJobs {
		var failedNodeResult models.FailedNodeResult
		failedNodeResult.JobId = int64(job.ID)
		if err := failedNodeResult.Delete(); err != nil {
			response.JSONR(c, 500, err)
		}

		job.JobHasSynced = 0
		job.JobHasComplete = 0
		if err := job.UpdateStatus(); err != nil {
			response.JSONR(c, 500, err)
		}
	}

	response.JSONR(c, 200, "ok")
}
