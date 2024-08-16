package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	pconfig "github.com/prometheus/common/config"
	"inspection/global"
	"inspection/models"
	"inspection/pkg/utils"
	"io/ioutil"
	"k8s.io/klog/v2"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ()

func main() {

	flag.StringVar(&global.ServerAddr, "server_addr", "http://localhost:8087/api/v1/node-result-report", "The server addr ")
	flag.StringVar(&global.ScriptPath, "script_path", "./script_xxx", "local script path ")
	flag.StringVar(&global.ResultPath, "result_path", "./result", "local script path ")
	version := flag.Bool("v", false, "prints current roxy version")
	flag.IntVar(&global.ExecTimeoutSeconds, "exec_timeout_seconds", 10, "exec tw sec")
	flag.Int64Var(&global.JobId, "job_id", 0, "jobid")
	flag.Parse()
	nodeIp := GetLocalIp()

	if *version {
		fmt.Println(global.AgentVersion)
		os.Exit(0)
	}

	oneResult := models.FailedNodeResult{
		NodeIp: nodeIp,
		JobId:  global.JobId,
	}
	// 读取一下result json结果
	resultBytes, err := ioutil.ReadFile(global.ResultPath)
	if err != nil {
		klog.Errorf("ioutil.ReadFile(resultPath).err[file:%v][path:%v]", global.ResultPath, err)
		oneResult.ErrMsg = err.Error()
	}

	desiredResultMap := map[string]string{}
	actualResultMap := map[string]string{}
	err = json.Unmarshal(resultBytes, &desiredResultMap)
	if err != nil {
		klog.Errorf("ComputeOneJob.desiredResult.ResultJson.json.Unmarshal.err:%v", err)
		return
	}
	klog.Infof("[desiredResultMap.print][%v]", desiredResultMap)
	out, err := CommandWithTw("/bin/bash", global.ScriptPath)

	if err != nil {
		oneResult.ErrMsg = err.Error()
	}
	err = json.Unmarshal([]byte(out), &actualResultMap)

	if err != nil {
		klog.Errorf("ComputeOneJob.actualResultMap.json.Unmarshal.[err:%v][jsonStr:%v]", err, out)
		return
	}

	klog.Infof("[actualResultMap.print][%v]", actualResultMap)
	// 对比两边的结果
	same := true
	if len(desiredResultMap) != len(actualResultMap) {
		same = false
	}
	for dk, dv := range desiredResultMap {
		if dv != actualResultMap[dk] {
			same = false
		}
	}

	oneResult.Succeed = same
	oneResult.ResultJson = out
	// 写入本地结果
	actualResultPath := fmt.Sprintf("%s_%s", global.ResultPath, "actual")
	err = ioutil.WriteFile(actualResultPath, []byte(out), 0644)
	klog.Infof("[run.result.print][same:%v][WriteFile.err:%v]", same, err)
	// 最多尝试10次

	for i := 0; i < 10; i++ {
		hc := pconfig.HTTPClientConfig{}
		_, err := utils.PostWithBearerToken("report-result", hc, 10, global.ServerAddr, oneResult)
		if err == nil {
			klog.Info("report-result.success")
			break
		}
		klog.Errorf("report-result.err[addr:%v][err:%v]", global.ServerAddr, err)
		time.Sleep(5 * time.Second)
	}

}

func GetLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Printf("get local addr err:%v", err)
		return ""
	} else {
		localIp := strings.Split(conn.LocalAddr().String(), ":")[0]
		conn.Close()
		return localIp
	}
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
