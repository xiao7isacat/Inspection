package check

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"inspection/global"
	"inspection/models"
	"inspection/pkg/common"
	"inspection/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

// 检测的管理器
type CheckJobManager struct {
	//配置文件
	Cg      *config.CheckJobConf
	Version string
}

const (
	AgentBinName = "node-env-check-agent"
)

// 初始化
func NewCheckJobManger(cg *config.Config) *CheckJobManager {
	return &CheckJobManager{Cg: cg.CheckJobC}
}

// 周期性检查数据库中是否有未下发的作业，有未下发的作业就下发
func (this *CheckJobManager) Run(ctx context.Context) error {
	//使用k8s中的wait库，周期性执行
	go wait.UntilWithContext(ctx, this.SpanCheckJob, time.Duration(this.Cg.CheckSubmitJobIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Info("RunCheckJobManger.exit.receive_quit_signal")
	return nil
}

// 下发任务
func (this *CheckJobManager) SpanCheckJob(ctx context.Context) {
	var checkJob models.CheckJob
	ljs, err := checkJob.GetNotSyncList()
	if err != nil {
		klog.Errorf("models.CronJobGetUnDone.err:%v", err)
		return
	}

	//所有任务已经下发
	if len(ljs) == 0 {
		klog.Warning("models.CronJobGetUnDone.zero")
		return
	}

	//对未下发的任务进行处理
	wp := workerpool.New(this.Cg.RunCheckJobBatch)
	for i := 0; i < len(ljs); i++ {
		job := ljs[i]
		if err := json.Unmarshal([]byte(job.IpJson), &job.IpList); err != nil {
			klog.Errorf("SpanCheckJob.ips.json.Unmarshal.err[job:%v][err:%v]", job, err)
			continue
		}
		wp.Submit(func() {
			this.SubmitJob(job)
		})
	}
	wp.StopWait()

	klog.Info("SpanCheckJob")

}

// 单个任务下发
func (this *CheckJobManager) SubmitJob(job *models.CheckJob) {
	remoteHost := strings.Join(job.IpList, ",")

	binFilePath := fmt.Sprintf("%s/%s",
		this.Cg.NodeRunCheckDir,
		global.AgentBinName,
	)
	thisJobDir := fmt.Sprintf("%s/%d_%s", this.Cg.NodeRunCheckDir, job.ID, job.Name)

	curlBinCmd := fmt.Sprintf("wget %s -O %s",
		this.Cg.AgentBinDownloadAddr,
		binFilePath,
	)
	checkBinVersionOrDownloadCmd := fmt.Sprintf(`%s -v| grep %s || %s  `,
		binFilePath,
		this.Version,
		curlBinCmd,
	)

	scriptFilePath := fmt.Sprintf("%s/%s.sh",
		thisJobDir,
		job.ScriptName,
	)
	resultFilePath := fmt.Sprintf("%s/%s.result",
		thisJobDir,
		job.DesiredName,
	)
	reportUrl := fmt.Sprintf("%s/api/v1/node-result-report", this.Cg.CheckServerAddr)
	/*

	   # 创建目录
	   [ ! -d "{{ NodeRunCheckDir }}" ] &&  mkdir {{ NodeRunCheckDir }}
	   # 下载agent的二进制
	   wget {{ AgentBinDownloadAddr  }} -O  {{ binFilePath }}
	   # curl 获取脚本
	   curl {{ CheckServerAddr  }}/api/v1/one-check-script?script_name={{ ScriptName }} > {{  scriptFilePath }}
	   # curl 获取基线
	   curl {{ CheckServerAddr  }}/api/v1/one-desired-result?result_name={{ DesiredResultName }} > {{  resultFilePath }}
	   # chmodCmd
	   chmod +x  {{ NodeRunCheckDir }}/*
	   # agent执行 ，并且给agent传参
	   {{ binFilePath }} -job_id={{ jobId }} -server_addr={{ reportUrl }} -result_path={{ resultFilePath }} -script_path={{ scriptFilePath }} &


	*/
	extraVars := map[string]interface{}{

		"NodeRunCheckDir":              this.Cg.NodeRunCheckDir,
		"thisJobDir":                   thisJobDir,
		"checkBinVersionOrDownloadCmd": checkBinVersionOrDownloadCmd, // 这里可以用nginx 或cdn
		"binFilePath":                  binFilePath,
		"CheckServerAddr":              this.Cg.CheckServerAddr,
		"ScriptName":                   job.ScriptName,
		"scriptFilePath":               scriptFilePath,
		"DesiredName":                  job.DesiredName,
		"resultFilePath":               resultFilePath,
		"jobId":                        job.ID,
		"reportUrl":                    reportUrl,
	}
	ansiYamlPath := global.SubmitJobYamlPath
	klog.V(2).Infof("SubmitJob.ips.common.AnsiRunPlay.print[job:%v] start", job.Name)
	if err := common.AnsiRunPlay(remoteHost, extraVars, ansiYamlPath); err != nil {
		klog.Errorf("SubmitJob.ips.common.AnsiRunPlay.print[job:%v][extraVars:%v][err:%v]", job.Name, extraVars, err)
		return
	}

	job.JobHasSynced = 1
	if err := job.Update(); err != nil {
		klog.Infof("SubmitJob.job.SetJobHasSynced.print[job:%v]update: false", job.Name)
		klog.V(2).Infof("SubmitJob.job.SetJobHasSynced.print[job:%v][update: false][err:%v]", job.Name, err)
		return
	}

	klog.Infof("SubmitJob.job.SetJobHasSynced.print[job:%v][updated:%v]", job.Name)
}

// 周期性的统计作业的成功 失败数量  标志位是 job_has_compute
func (cm *CheckJobManager) RunComputeJobManager(ctx context.Context) error {
	go wait.UntilWithContext(ctx, cm.ComputeCheckJob, time.Duration(cm.Cg.CompleteJobIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunComputeJobManager.exit.receive_quit_signal")
	return nil
}

func (cm *CheckJobManager) ComputeCheckJob(ctx context.Context) {
	// 获取还需要统计
	var checkJob models.CheckJob
	jobs, err := checkJob.GetNotCompleteList()
	if err != nil {
		klog.Errorf("models.CheckJobGetUnCompute.err:%v", err)
		return
	}
	if len(jobs) == 0 {
		klog.Warning("models.CheckJobGetUnCompute.zero")
		return
	}

	klog.Infof("models.CheckJobGetUnCompute.num:%v", len(jobs))
	wp := workerpool.New(cm.Cg.RunCheckJobBatch)
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		// 因为机器比较多，还没上报完，这时候统计的成功失败数量是 不准的
		if time.Now().Sub(job.UpdatedAt).Minutes() < float64(job.JobWaitCompleteMinutes) {
			klog.Infof("models.CheckJobGetUnCompute.still.in.wait:%v", job.Name)
			continue
		}

		wp.Submit(func() {
			// 启动任务
			cm.ComputeOneJob(job)
		})
	}
	wp.StopWait()
}

func (cm *CheckJobManager) ComputeOneJob(cj *models.CheckJob) {

}
