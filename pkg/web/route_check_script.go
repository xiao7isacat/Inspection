package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inspection/models"
	"inspection/pkg/common"
	"k8s.io/klog/v2"
)

func CheckScriptAdd(c *gin.Context) {

	var input models.CheckScript
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	id, err := input.CreateOne()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}
	klog.Infof("[CheckScriptAdd.success][script:%v]", input.Name)
	common.JSONR(c, 200, id)
}

func CheckScriptGets(c *gin.Context) {
	var checkScript models.CheckScript
	ljs, err := checkScript.GetList()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, 200, ljs)
}

func CheckScriptPuts(c *gin.Context) {
	var input models.CheckScript
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

func CheckScriptGetByName(c *gin.Context) {
	scriptName := c.Query("script_name")
	if scriptName == "" {
		c.String(400, "empty name")
		return
	}
	scriptContent := ""

	var checkScript models.CheckScript
	checkScript.Name = scriptName
	err := checkScript.GetOne()
	if err != nil {
		c.String(500, fmt.Errorf("models.CheckScriptGetByName.err:%w", err).Error())
		return
	}

	scriptContent = checkScript.ContentJson

	c.String(200, scriptContent)
}
