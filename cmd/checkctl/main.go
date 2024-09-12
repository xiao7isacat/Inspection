package main

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"inspection/global"
	"inspection/pkg/checkctl"
	"k8s.io/klog/v2"
	"os"
)

var RootCommand = &cobra.Command{
	Use: "checkctl",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCommand.AddCommand(checkctl.AddCmd)
	RootCommand.AddCommand(checkctl.GetCmd)
	RootCommand.AddCommand(checkctl.UpdateCmd)
	RootCommand.AddCommand(checkctl.DeleteCmd)
	RootCommand.AddCommand(checkctl.RunCmd)

	RootCommand.Flags().StringVar(&global.ServerAddr, "server_addr", "127.0.0.1:8092", "server_addr")
	RootCommand.PersistentFlags().StringVarP(&global.ResourceFilePath, "file", "f", "test", "指定传入的文件")
	RootCommand.PersistentFlags().StringVarP(&global.ResourceName, "name", "n", "test", "指定资源名称")
	RootCommand.PersistentFlags().StringVarP(&global.NodeAddrs, "node_addr", "", "127.0.0.1:8093", "node_addr")
}
func main() {
	klog.InitFlags(flag.CommandLine)

	if err := RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
