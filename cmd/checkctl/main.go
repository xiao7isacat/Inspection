package main

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"inspection/cmd"
	"inspection/global"
	"k8s.io/klog/v2"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "ctl",
	Short: "ctl is control for Inspection",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCommand.AddCommand(cmd.AddCommand)
	RootCommand.AddCommand(cmd.GetCommand)
	RootCommand.AddCommand(cmd.UpdateCommand)
	RootCommand.AddCommand(cmd.DeleteCommand)

	RootCommand.Flags().StringVar(&global.ServerAddr, "server_addr", "127.0.0.1", "server_addr")
}
func main() {
	klog.InitFlags(flag.CommandLine)
	flag.StringVar(&global.Add, "add", "", "add")
	flag.StringVar(&global.Delete, "delete", "", "delete")
	flag.StringVar(&global.Update, "update", "", "update")
	flag.StringVar(&global.Getone, "getone", "", "getone")
	flag.StringVar(&global.Getone, "getall", "", "getall")
	version := flag.Bool("v", false, "prints current roxy version")
	flag.Parse()
	if *version {
		fmt.Println(global.AgentVersion)
		os.Exit(0)
	}

	if err := RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
