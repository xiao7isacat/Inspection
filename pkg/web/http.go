package web

import (
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"inspection/global"
	"inspection/pkg/check"
	"inspection/pkg/config"
	"k8s.io/klog/v2"
	"net/http"
	"time"
)

func StartServer(cf *config.Config, cm *check.CheckJobManager) error {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	//ginweb 使用prometheus sdk打点
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	// 外部指针变量传递给gin，在view中使用
	m := make(map[string]interface{})
	m[global.CheckJobManager] = cm
	r.Use(ConfigMiddleware(m))

	//路由
	configServerRouters(r)
	s := &http.Server{
		Addr:              cf.HttpAddr,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		MaxHeaderBytes:    1 << 2,
	}

	klog.Infof("[web.server.available.at:%v]", cf.HttpAddr)
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func StartAgent() error {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	//路由
	configAgentRouters(r)
	s := &http.Server{
		Addr:              ":" + global.AgentPort,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		MaxHeaderBytes:    1 << 2,
	}

	klog.Infof("[web.server.available.at:%v]", global.AgentPort)
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
