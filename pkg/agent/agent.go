package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pconfig "github.com/prometheus/common/config"
	"inspection/global"
	"inspection/models"
	"inspection/pkg/check"
	"inspection/pkg/utils"
	"io/ioutil"
	"k8s.io/klog/v2"
	"log"
	"os/exec"
	"time"
)

func ExecJobs(cj *check.CheckJobManager) error {

	if cj == nil {
		log.Println("checkJobManger is nil")
		return fmt.Errorf("checkJobManger is nil")
	}
	if cj.Cg == nil {
		log.Println("checkJobManger.Cg is nil")
		return fmt.Errorf("checkJobManger.Cg is nil")
	}
	if cj.AgentParameters == nil {
		log.Println("checkJobManger.AgentParameters is nil")
		return fmt.Errorf("checkJobManger.AgentParameters is nil")
	}

	//检测工作目录是否存在
	if cj.AgentParameters.JobDir == "" {
		return fmt.Errorf("job %v work directory is nil", cj.AgentParameters.ScriptName)
	}

	if err := cj.Download(); err != nil {
		return fmt.Errorf("job %v download err:%v", cj.AgentParameters.ScriptName, err)
	}

	oneResult := models.FailedNodeResult{
		NodeIp: cj.AgentParameters.NodeIP,
		JobId:  cj.AgentParameters.JobId,
	}

	desiredFileName := cj.AgentParameters.JobDir + "/" + cj.AgentParameters.ResultName + ".result"
	scriptFileName := cj.AgentParameters.JobDir + "/" + cj.AgentParameters.ResultName + ".sh"

	resultBytes, err := ioutil.ReadFile(desiredFileName)
	if err != nil {
		klog.Errorf("ioutil.ReadFile(resultPath).err[file:%v][path:%v]", cj.AgentParameters.ResultName, err)
		oneResult.ErrMsg = err.Error()
	}
	klog.Infof("desired is %v", string(resultBytes))

	desiredResultMap := map[string]string{}
	actualResultMap := map[string]string{}
	if err = json.Unmarshal(resultBytes, &desiredResultMap); err != nil {
		klog.Errorf("ComputeOneJob.desiredResult.ResultJson.json.Unmarshal.err:%v", err)
		return fmt.Errorf("Unmarshal resultBytes to desiredResultMap failed %v", err)
	}
	klog.Infof("[desiredResultMap.print][%v]", desiredResultMap)
	out, err := CommandWithTw("/bin/bash", scriptFileName)
	if err != nil {
		oneResult.ErrMsg = err.Error()
	}
	err = json.Unmarshal([]byte(out), &actualResultMap)
	klog.Infof("[actualResultMap.print][%v]", actualResultMap)
	same := true
	if len(desiredResultMap) != len(actualResultMap) {
		same = false
	} else {
		for dk, dv := range desiredResultMap {
			if !utils.Duibi(dv, actualResultMap[dk]) {
				break
			}

		}
	}

	oneResult.Succeed = same
	oneResult.ResultJson = out
	klog.V(2).Info(oneResult)
	actualResultPath := fmt.Sprintf("%s_%s", desiredFileName, "actual")
	err = ioutil.WriteFile(actualResultPath, []byte(out), 0644)
	klog.Infof("[run.result.print][same:%v][WriteFile.err:%v]", same, err)

	reportUrl := cj.Cg.CheckServerAddr + "/api/v1/node-result-report"
	for i := 0; i < 10; i++ {
		hc := pconfig.HTTPClientConfig{}
		_, err := utils.PostWithBearerToken("report-result", hc, 10, reportUrl, oneResult)
		if err == nil {
			klog.Info("report-result.success")
			break
		}
		klog.Errorf("report-result.err[addr:%v][err:%v]", reportUrl, err)
		time.Sleep(3 * time.Second)
	}
	return nil
}

func CommandWithTw(name string, arg ...string) (string, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Duration(global.ExecTimeoutSeconds)*time.Second)
	defer cancel() //releases resources if slowOperation completes before timeout elapses
	cmd := exec.CommandContext(ctxt, name, arg...)
	//当经过Timeout时间后，程序依然没有运行完，则会杀掉进程，ctxt也会有err信息
	if out, err := cmd.Output(); err != nil {
		//检测报错是否是因为超时引起的
		if ctxt.Err() != nil && ctxt.Err() == context.DeadlineExceeded {
			return "", errors.New("command timeout")

		}
		return string(out), err
	} else {
		return string(out), nil
	}
}
