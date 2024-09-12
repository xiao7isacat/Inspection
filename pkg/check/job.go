package check

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"inspection/models"
	"inspection/pkg/config"
	"inspection/pkg/utils"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// 检测的管理器
type CheckJobManager struct {
	//配置文件
	Cg              *config.CheckJobConf
	Version         string
	AgentParameters *AgentArgs
	rw              sync.RWMutex
}

// agent执行参数
type AgentArgs struct {
	JobDir     string `json:"job_dir"`
	ScriptName string `json:"script_name"`
	ResultName string `json:"result_file_name"`
	JobId      int64  `json:"job_id"`
	NodeIP     string `json:"node_ip"`
}

// server初始化
func ServerCheckJobManger(cg *config.Config) *CheckJobManager {
	return &CheckJobManager{Cg: cg.CheckJobC}
}

// 初始化空的任务控制器
func NilCheckJobManger() *CheckJobManager {
	return &CheckJobManager{}
}

// 周期性检查数据库中是否有未下发的作业，有未下发的作业就下发
func (this *CheckJobManager) RunCheckJobManger(ctx context.Context) error {
	//使用k8s中的wait库，周期性执行
	go wait.UntilWithContext(ctx, this.SpanCheckJob, time.Duration(this.Cg.CheckSubmitJobIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Info("RunCheckJobManger.exit.receive_quit_signal")
	return nil
}

// 下发任务
func (this *CheckJobManager) SpanCheckJob(ctx context.Context) {
	var (
		checkJob models.CheckJob
	)

	ljs, err := checkJob.GetNotSyncList()
	if err != nil {
		klog.Errorf("models.SpanCheckJobGetUnDone.err:%v", err)
		return
	}

	//所有任务已经下发
	if len(ljs) == 0 {
		klog.Warning("models.SpanCheckJobGetUnDone.zero")
		return
	}

	//对未下发的任务进行处理
	wp := workerpool.New(this.Cg.RunCheckJobBatch)
	for _, job := range ljs {
		wp.Submit(func() {
			this.SubmitJob(job)
		})
	}
	wp.StopWait()
	//this.SubmitJob(job)

	klog.Info("SpanCheckJob")

}

func (this *CheckJobManager) SubmitJob(job *models.CheckJob) {
	agentParameters := &AgentArgs{}
	agentParameters.JobId = int64(job.ID)
	agentParameters.JobDir = fmt.Sprintf("workjobdir"+"/%d_%s", job.ID, job.Name)
	agentParameters.ScriptName = job.ScriptName
	agentParameters.ResultName = job.DesiredName
	checkJobManger := NilCheckJobManger()
	checkJobManger.Cg = this.Cg
	checkJobManger.AgentParameters = agentParameters
	if checkJobManger.AgentParameters == nil {
		klog.V(2).Infof("SubmitJob.job.Post.print[nil]", job.Name, agentParameters.JobDir)
	}

	job.IpList = strings.Split(job.IpString, ",")
	klog.V(2).Infof("SubmitJob.job.Post.print[job:%v][date: %v]", job.Name, checkJobManger.AgentParameters.JobDir)

	wp := workerpool.New(len(job.IpList))
	for _, host := range job.IpList {
		wp.Submit(func() {
			url := "http://" + host + "/api/v1/run-check-script"
			checkJobManger.rw.Lock()
			checkJobManger.AgentParameters.NodeIP = host
			jsonData, err := json.Marshal(checkJobManger)
			checkJobManger.rw.Unlock()
			if err != nil {
				klog.Infof("SubmitJob.job.print[job:%v][Marshal: false][err:%v]", job.Name, err)
			}
			if err = utils.Post(url, jsonData); err != nil {
				klog.Infof("SubmitJob.job.Post.print[job:%v][err:%v]", job.Name, err)
			}
			klog.Infof("SubmitJob.job.Post.print[job:%v][url: %v]", job.Name, url)
		})
	}
	wp.StopWait()

	job.JobHasSynced = 1
	if err := job.Update(); err != nil {
		klog.Infof("SubmitJob.job.SetJobHasSynced.print[job:%v]update: false", job.Name)
		klog.V(2).Infof("SubmitJob.job.SetJobHasSynced.print[job:%v][update: false][err:%v]", job.Name, err)
		return
	}

	klog.Infof("SubmitJob.job.SetJobHasSynced.print[job:%v][updated]", job.Name)
}

// 下载脚本和基线
func (this *CheckJobManager) Download() error {
	var (
		desiredData []byte
		scriptData  []byte
	)

	if this.Cg == nil {
		log.Println("Download CheckJobManger.Cg is nil")
		return fmt.Errorf("Download CheckJobManger.Cg is nil")
	}
	if this.AgentParameters == nil {
		log.Println("Download CheckJobManger.AgentParameters is nil")
		return fmt.Errorf("Download CheckJobManger.AgentParameters is nil")
	}

	//检测工作目录是否存在
	if this.AgentParameters.JobDir == "" {
		return fmt.Errorf("job %v work directory is nil", this.AgentParameters.ScriptName)
	}
	workPath := this.AgentParameters.JobDir
	_, err := os.Stat(workPath)
	if err != nil {
		// 如果返回错误，则目录不存在
		if os.IsNotExist(err) {
			// 创建目录
			err = os.MkdirAll(workPath, 0755) // 设置目录权限
			if err != nil {
				log.Printf("Error creating directory: %v\n", err)
				return fmt.Errorf("Error creating directory: %v\n", err)
			}
			log.Println("Directory created:", workPath)
		} else {
			log.Printf("Error checking directory: %v\n", err)
			return fmt.Errorf("Error checking directory: %v\n", err)

		}
	}

	//下载基线
	desiredUrl := this.Cg.CheckServerAddr + "/api/v1/one-desired-result"
	desiredparams := url.Values{}
	desiredparams.Add("result_name", this.AgentParameters.ResultName)
	dersiredUrlWithParams := fmt.Sprintf("%s?%s", desiredUrl, desiredparams.Encode())
	if desiredData, err = utils.Get(dersiredUrlWithParams); err != nil {
		log.Printf("checkJobManger.AgentParameters download desired from %v,failed: %v", dersiredUrlWithParams, err)
		return fmt.Errorf("checkJobManger.AgentParameters download desired from %v,failed: %v", dersiredUrlWithParams, err)
	}
	desiredFileName := this.AgentParameters.JobDir + "/" + this.AgentParameters.ResultName + ".result"
	err = ioutil.WriteFile(desiredFileName, desiredData, 0755)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return fmt.Errorf("Error writing file: %v\n", err)
	}

	//下载脚本
	scriptUrl := this.Cg.CheckServerAddr + "/api/v1/one-check-script"
	scriptparams := url.Values{}
	scriptparams.Add("script_name", this.AgentParameters.ResultName)
	scriptUrlWithParams := fmt.Sprintf("%s?%s", scriptUrl, scriptparams.Encode())
	if scriptData, err = utils.Get(scriptUrlWithParams); err != nil {
		log.Printf("checkJobManger.AgentParameters download script from %v,failed: %v", scriptUrlWithParams, err)
		return fmt.Errorf("checkJobManger.AgentParameters download script from %v,failed: %v", scriptUrlWithParams, err)
	}
	scriptFileName := this.AgentParameters.JobDir + "/" + this.AgentParameters.ResultName + ".sh"
	err = ioutil.WriteFile(scriptFileName, scriptData, 0755)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return fmt.Errorf("Error writing file: %v\n", err)
	}

	return nil
}

// 周期性的统计作业的成功 失败数量  标志位是 job_has_compute
func (this *CheckJobManager) RunComputeJobManager(ctx context.Context) error {
	go wait.UntilWithContext(ctx, this.ComputeCheckJob, time.Duration(this.Cg.CompleteJobIntervalSeconds)*time.Second)
	<-ctx.Done()
	klog.Infof("RunComputeJobManager.exit.receive_quit_signal")
	return nil
}

// 统计任务
func (this *CheckJobManager) ComputeCheckJob(ctx context.Context) {
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
	wp := workerpool.New(this.Cg.RunCheckJobBatch)
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		// 因为机器比较多，还没上报完，这时候统计的成功失败数量是 不准的
		date := time.Now().Sub(job.UpdatedAt).Minutes()
		if time.Now().Sub(job.UpdatedAt).Minutes() < float64(job.JobWaitCompleteMinutes) {
			klog.Info(date, float64(job.JobWaitCompleteMinutes))
			klog.Infof("models.CheckJobGetUnCompute.still.in.wait:%v", job.Name)
			continue
		}

		wp.Submit(func() {
			// 启动任务
			this.ComputeOneJob(job)
		})
	}
	wp.StopWait()
}

func (this *CheckJobManager) ComputeOneJob(cj *models.CheckJob) {
	var successNodeResult models.FailedNodeResult
	successNodeResult.JobId = int64(cj.ID)
	successNodeResult.FinalSucceed = 1
	thisJobSuccessNodes, err := successNodeResult.GetList()
	if err != nil {
		klog.Errorf("ComputeOneJob.models.SuccessNodeResultByJobId.err[err:%v][jobName:%v]", err, cj.Name)
		return
	}

	var failedNodeResult models.FailedNodeResult
	failedNodeResult.JobId = int64(cj.ID)
	failedNodeResult.FinalSucceed = 2
	thisJobFailedNodes, err := failedNodeResult.GetList()
	if err != nil {
		klog.Errorf("ComputeOneJob.models.FailedNodeResultByJobId.err[err:%v][jobName:%v]", err, cj.Name)
		return
	}

	cj.SuccessNum = int64(len(thisJobSuccessNodes))
	cj.FailedNum = int64(len(thisJobFailedNodes))
	klog.Info("failed num ", cj.FailedNum)
	cj.FailedNum = int64(0)
	cj.MissNum = cj.AllNum - cj.SuccessNum - cj.FailedNum
	cj.JobHasComplete = 1
	if err = cj.UpdateNodeStatus(); err != nil {
		klog.Errorf("ComputeOneJob.models.cj.Update.err[err:%v][jobName:%v]", err, cj.Name)
		return
	}
	klog.Infof("ComputeOneJob.models.cj.Update[updated:%v][jobName:%v]", cj.Name)
	// 标记为结束
}

// 下发任务第一版
/*func (this *CheckJobManager) SpanCheckJob(ctx context.Context) {
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

//   # 创建目录
//   [ ! -d "{{ NodeRunCheckDir }}" ] &&  mkdir {{ NodeRunCheckDir }}
//   # 下载agent的二进制
//   wget {{ AgentBinDownloadAddr  }} -O  {{ binFilePath }}
//   # curl 获取脚本
//   curl {{ CheckServerAddr  }}/api/v1/one-check-script?script_name={{ ScriptName }} > {{  scriptFilePath }}
//   # curl 获取基线
//   curl {{ CheckServerAddr  }}/api/v1/one-desired-result?result_name={{ DesiredResultName }} > {{  resultFilePath }}
//   # chmodCmd
//   chmod +x  {{ NodeRunCheckDir }}/*
//   # agent执行 ，并且给agent传参
//   {{ binFilePath }} -job_id={{ jobId }} -server_addr={{ reportUrl }} -result_path={{ resultFilePath }} -script_path={{ scriptFilePath }} &

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
}*/

// 下发任务第二版
/*func (this *CheckJobManager) SpanCheckJob(ctx context.Context) {
	var (
		checkJob   models.CheckJob
		ipAddrInfo models.IpAddrInfo
		HostsJob   = make(map[string]*AgentCheckJobManagers)
	)

	ljs, err := checkJob.GetNotSyncList()
	if err != nil {
		klog.Errorf("models.SpanCheckJobGetUnDone.err:%v", err)
		return
	}

	//所有任务已经下发
	if len(ljs) == 0 {
		klog.Warning("models.SpanCheckJobGetUnDone.zero")
		return
	}

	ipAddrInfos, err := ipAddrInfo.GetList()
	if err != nil {
		klog.Errorf("models.ipAddrInfoGet.err:%v", err)
		return
	}

	//遍历ip
	for _, ipInfo := range ipAddrInfos {
		//遍历job
		for _, job := range ljs {
			//如果job中有该ip，将job添加进hostjob的结构中。这样在发送请求时。一个host发送一次就行
			if strings.HasPrefix(job.IpString, ipInfo.Ip) {
				agentCheckJobManger := NilCheckJobManger()
				agentCheckJobManger.AgentJobDir = fmt.Sprintf("%d_%s", job.ID, job.Name)
				agentCheckJobManger.AgentScriptName = fmt.Sprintf("%s.sh", job.ScriptName)
				agentCheckJobManger.AgentResultFileName = fmt.Sprintf("%s.result", job.DesiredName)
				HostsJob[ipInfo.Ip].checkJobMangerList = append(HostsJob[ipInfo.Ip].checkJobMangerList, agentCheckJobManger)
				HostsJob[ipInfo.Ip].Cg = this.Cg
				HostsJob[ipInfo.Ip].Version = this.Version
				HostsJob[ipInfo.Ip].Name = this.Version
			}
		}
	}

	//对未下发的任务进行处理
	wp := workerpool.New(len(ipAddrInfos))
	for _, ipInfo := range ipAddrInfos {
		HostsJob[ipInfo.Ip].SubmitJob(ipInfo.Ip)
	}
	wp.StopWait()
	//this.SubmitJob(job)

	klog.Info("SpanCheckJob")
}

func (this *AgentCheckJobManagers) SubmitJob(nodeIp string) {

	klog.V(2).Infof("SubmitJob.job.Post.print[job:%v][date: %v]", this.Name, &this)
	jsonData, err := json.Marshal(&this)
	if err != nil {
		klog.Infof("SubmitJob.job.print[job:%v][Marshal: false][err:%v]", this.Name, err)
	}

	url := "http://" + nodeIp + "/api/v1/run-check-script"

	if err := utils.Post(url, jsonData); err != nil {
		klog.Infof("SubmitJob.job.Post.print[job:%v][url: %v][err:%v]", this.Name, url, err)
	}
	klog.Infof("SubmitJob.job.Post.print[job:%v][url: %v]", this.Name, url)

	//job.JobHasSynced = 1
	//if err := job.Update(); err != nil {
	//	klog.Infof("SubmitJob.job.SetJobHasSynced.print[job:%v]update: false", this.Name)
	//	klog.V(2).Infof("SubmitJob.job.SetJobHasSynced.print[job:%v][update: false][err:%v]", job.Name, err)
	//	return
	//}

	klog.Infof("SubmitJob.job.SetJobHasSynced.print[job:%v][updated]", job.Name)
}*/
