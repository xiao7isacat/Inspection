package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inspection/models"
	"inspection/pkg/common"
	"k8s.io/klog/v2"
)

func DesiredResultAdd(c *gin.Context) {

	var input models.DesiredResult
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	fmt.Println(input)
	id, err := input.CreateOne()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}
	klog.Infof("[DesiredResult.success][DesiredResult:%v]", input.Name)
	common.JSONR(c, 200, id)
}

func DesiredResultGets(c *gin.Context) {
	var desiredResult models.DesiredResult
	ljs, err := desiredResult.GetList()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, 200, ljs)
}

func DesiredResultPuts(c *gin.Context) {
	var input models.DesiredResult
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}

	err := input.Update()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	klog.Infof("[CheckScript.update][CheckScript:%v]", input.Name)
	common.JSONR(c, 200)
}

// 给agent用的 根据基线的名称获取 基线json内容
func DesiredResultByName(c *gin.Context) {
	resultName := c.Query("result_name")
	if resultName == "" {
		c.String(400, "empty name")
		return
	}
	ResultJson := ""
	var desiredResult models.DesiredResult
	desiredResult.Name = resultName
	err := desiredResult.GetOne()
	if err != nil {
		c.String(500, fmt.Errorf("models.DesiredResultByName.err:%w", err).Error())
		return
	}
	ResultJson = desiredResult.ResultJson
	c.String(200, ResultJson)
}
