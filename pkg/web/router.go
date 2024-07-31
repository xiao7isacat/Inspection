package web

import (
	"github.com/gin-gonic/gin"
	"time"
)

func configRouters(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("healthy", func(c *gin.Context) {
		c.String(200, "ok")
	})
	api.GET("now-ts", GetNowTs)

	api.POST("/check-script", CheckScriptAdd)
	api.PUT("/check-script", CheckScriptPuts)
	api.GET("/check-script", CheckScriptGets)
	api.GET("/one-check-script", CheckScriptGetByName)

	api.POST("/desired-result", DesiredResultAdd)
	api.GET("/desired-result", DesiredResultGets)
	api.PUT("/desired-result", DesiredResultPuts)
	api.GET("/one-desired-result", DesiredResultByName)

	api.POST("/check-job", CheckJobAdd)
	//api.POST("/cron-job", CronJobAdd)
	api.GET("/check-job", CheckJobGets)

	//api.POST("/node-result-report", NodeResultReport)

}

func GetNowTs(c *gin.Context) {
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
