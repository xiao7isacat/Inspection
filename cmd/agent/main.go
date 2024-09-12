package main

import (
	"flag"
	"fmt"
	"inspection/global"
	"inspection/pkg/web"
	"k8s.io/klog/v2"
	"os"
)

var ()

func main() {
	klog.InitFlags(flag.CommandLine)
	version := flag.Bool("version", false, "prints current roxy version")
	flag.IntVar(&global.ExecTimeoutSeconds, "exec_timeout_seconds", 10, "exec tw sec")
	flag.StringVar(&global.AgentPort, "agentport", "8093", "agent port")
	flag.Parse()

	if *version {
		fmt.Println(global.AgentVersion)
		os.Exit(0)
	}

	if err := web.StartAgent(); err != nil {
		klog.Errorf("[start.web.error][err:%v]", err)
		os.Exit(1)
	}
}
