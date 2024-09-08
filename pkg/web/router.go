package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func configServerRouters(r *gin.Engine) {

	r.StaticFS("static-file", http.Dir("static-file"))
	log.Println(os.Getwd())

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

	api.POST("/node-result-report", NodeResultReport)

}

func configAgentRouters(r *gin.Engine) {

	r.StaticFS("static-file", http.Dir("static-file"))
	log.Println(os.Getwd())

	api := r.Group("/api/v1")
	api.GET("healthy", func(c *gin.Context) {
		c.String(200, "ok")
	})
	api.GET("now-ts", GetNowTs)

	api.POST("/run-check-script", CheckScriptRun)
	//api.GET("/check-script", CheckScriptGets)
}

func GetNowTs(c *gin.Context) {
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
