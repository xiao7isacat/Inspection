package main

import (
	"context"
	"flag"
	"fmt"
	"inspection/database"
	"inspection/global"
	"inspection/models"
	"inspection/pkg/config"
	esl "inspection/pkg/signal"
	"inspection/pkg/web"
	"k8s.io/klog/v2"
)

func main() {
	//klog本身也有自己的命令行参数，让其可以使用klog的命令行参数
	klog.InitFlags(flag.CommandLine)
	flag.StringVar(&global.ConfigFile, "config", "./node-env-check.yaml", "config file")
	flag.StringVar(&global.Database, "database", "sqlite", "dbname")
	flag.Parse()

	sConfig, err := config.LoadFile(global.ConfigFile)
	if err != nil {
		klog.Errorln(err)
		return
	}
	klog.V(2).Infof("config.LoadFile.success.print:%+v", sConfig)

	//初始db连接
	if err := database.ConnectDb(global.Database); err != nil {
		klog.Errorln("database:", global.Database, "error: ", err)
		return
	}

	models.AutoMigrat()

	//接受信号，开始编排
	group, stopChan := esl.SetupStopSignalContext()
	ctlAll, cancelAll := context.WithCancel(context.Background())

	//接收退出信号的ctx
	group.Go(func() error {
		klog.Infof("[stop chan watch start backend]")
		for {
			select {
			case <-stopChan:
				klog.Infof("[stop chan receive quite signal exit]")
				cancelAll()
				return nil

			}
		}
	})

	group.Go(func() error {
		klog.Infof("[metrics web start backend]")
		errChan := make(chan error)
		go func() {
			errChan <- web.StartServer(sConfig)
		}()

		select {
		case err := <-errChan:
			klog.Errorf("[web.server.error][err:%v]", err)
			return err
		case <-ctlAll.Done():
			klog.Info("receive.quit.singal.web.server.exit")
			return nil
		}
	})

	fmt.Println(ctlAll)
	if err := group.Wait(); err != nil {
		//pianc
		klog.Fatal(err)
	}

}
