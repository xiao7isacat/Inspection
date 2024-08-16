package web

import (
	"github.com/gin-gonic/gin"
	"inspection/models"
	"inspection/pkg/common"
	"k8s.io/klog/v2"
)

func NodeResultReport(c *gin.Context) {

	var input models.FailedNodeResult
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	// 成功的不记录，返回success
	if input.Succeed == true {
		input.FinalSuccess = 1
		klog.Infof("[NodeResultReport.node.success][ip:%v][job_id:%v]", input.NodeIp, input.JobId)
		//common.JSONR(c, 200, "success")
		//return
	}

	id, err := input.CreateOrUpdate()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}
	common.JSONR(c, 200, id)

	klog.Infof("[Failed.NodeResultAdd.success][ip:%v][job_id:%v]", input.NodeIp, input.JobId)
}
