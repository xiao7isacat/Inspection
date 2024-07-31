package check

import (
	"context"
	"inspection/models"
	"inspection/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"time"
)

// 检测的管理器
type CheckJobManger struct {
	//配置文件
	Cg *config.CheckJobConf
}

// 初始化
func NewCheckJobManger(cg *config.Config) *CheckJobManger {
	return &CheckJobManger{Cg: cg.CheckJobC}
}

// 周期性检查数据库中是否有未下发的作业，有未下发的作业就下发
func (this *CheckJobManger) Run(ctx context.Context) error {
	//使用k8s中的wait库，周期性执行
	go wait.UntilWithContext(ctx, this.SpanCheckJob, time.Duration(this.Cg.CheckSubmitJobIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Info("RunCheckJobManger.exit.receive_quit_signal")
	return nil
}

// 下发任务
func (this *CheckJobManger) SpanCheckJob(ctx context.Context) {
	var checkJob models.CheckJob
	ljs, err := checkJob.GetList(true)
	if err != nil {
		klog.Errorf("models.CronJobGetUnDone.err:%v", err)
		return
	}
	if len(ljs) == 0 {
		klog.Warning("models.CronJobGetUnDone.zero")
		return
	}
	klog.Info("SpanCheckJob")

}
