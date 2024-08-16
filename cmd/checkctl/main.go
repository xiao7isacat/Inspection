package main

import (
	"flag"
	"fmt"
	"inspection/global"
	"k8s.io/klog/v2"
	"os"
)

func main() {
	klog.InitFlags(flag.CommandLine)
	flag.StringVar(&global.Add, "add", "", "add")
	flag.StringVar(&global.Delete, "delete", "", "delete")
	flag.StringVar(&global.Update, "update", "", "update")
	flag.StringVar(&global.Getone, "getone", "", "getone")
	flag.StringVar(&global.Getone, "getall", "", "getall")
	flag.StringVar(&global.Getone, "getall", "", "getall")
	flag.StringVar(&global.ServerAddr, "server_addr", "127.0.0.1", "server_addr")
	version := flag.Bool("v", false, "prints current roxy version")
	flag.Parse()
	if *version {
		fmt.Println(global.AgentVersion)
		os.Exit(0)
	}
}
