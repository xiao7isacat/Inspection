package web

import (
	"github.com/gin-gonic/gin"
	"inspection/models"
	"inspection/pkg/response"
	"k8s.io/klog/v2"
)

func NodeResultReport(c *gin.Context) {

	var input models.FailedNodeResult
	if err := c.BindJSON(&input); err != nil {
		response.JSONR(c, 400, err)
		return
	}
	// 成功的不记录，返回success

	if input.Succeed == true {
		input.FinalSucceed = 1
		klog.Infof("[NodeResultReport.node.success][ip:%v][job_id:%v]", input.NodeIp, input.JobId)
		//response.JSONR(c, 200, "success")
		//return
	} else {
		input.FinalSucceed = 2
		klog.Infof("[NodeResultReport.node.failed][ip:%v][job_id:%v]", input.NodeIp, input.JobId)
	}

	id, err := input.CreateOrUpdate()
	if err != nil {
		response.JSONR(c, 500, err)
		return
	}
	response.JSONR(c, 200, id)

	klog.Infof("[Failed.NodeResultAdd.success][ip:%v][job_id:%v]", input.NodeIp, input.JobId)
}
