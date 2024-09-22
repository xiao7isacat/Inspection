package checkctl

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"inspection/global"
	"inspection/pkg/utils"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type JobStatus struct {
	JobName     string `json:"name"`
	CheckName   string `json:"check_name"`
	ActualValue string `json:"actual_value"`
	ExpectValue string `json:"expect_value"`
	Node        string `json:"node"`
	Status      string `json:"status"`
	AllDone     bool   `json:"all_done"`
}

func (this *JobStatus) Get() error {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)
	fmt.Print("\033[2J\033[H")
	go func() {
		for {

			// 发送 HTTP GET 请求
			var (
				jobInfos []JobStatus
				done     bool
			)

			baseURL := "http://" + global.ServerAddr + "/api/v1/ctl-status-job"
			// 构建查询参数
			jobResourceDate, err := utils.Get(baseURL)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// clearScreen 使用终端控制序列清除屏幕

			//fmt.Print("\033[2J\033[H")
			cmd := exec.Command("tput", "reset")
			cmd.Stdout = os.Stdout
			cmd.Run()

			if err = json.Unmarshal(jobResourceDate, &jobInfos); err != nil {
				fmt.Printf("结果转换格式失败：%v", err)
				os.Exit(1)
			}
			originalHeader := []string{"任务名称", "检测名称", "节点", "状态", "期望值", "真实值"}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(originalHeader)
			// 不显示表格线
			table.SetBorder(false)
			// 不显示水平线条
			table.SetAutoWrapText(false)
			// 不显示竖直线条
			table.SetColumnSeparator(" ")
			for _, jobInfo := range jobInfos {
				table.Append([]string{jobInfo.JobName, jobInfo.CheckName, jobInfo.Node, jobInfo.Status, jobInfo.ExpectValue, jobInfo.ActualValue})
				done = jobInfo.AllDone
			}
			table.Render()

			if done {
				err = ioutil.WriteFile("check_result.txt", jobResourceDate, 0755)
				if err != nil {
					fmt.Printf("Error writing file check_result: %v\n", err)
					os.Exit(1)
				}
				os.Exit(0)
			}
			select {
			case <-exitSignal:
				// 接收到退出信号
				fmt.Print("\033[2J\033[H")
				os.Exit(0)
			default:
				// 等待一段时间再次发送请求
				time.Sleep(1 * time.Second)
			}

		}
	}()
	<-exitSignal
	return nil
}

func (this *JobStatus) Add() error {
	klog.Info("add job")

	return nil
}

func (this *JobStatus) Update() error {
	klog.Info("update job")
	return nil
}

func (this *JobStatus) Delete() error {
	klog.Info("delete job")
	return nil
}
