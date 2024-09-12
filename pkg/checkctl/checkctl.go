package checkctl

import (
	"github.com/spf13/cobra"
	"inspection/global"
	"k8s.io/klog/v2"
	"os"
)

var resourceMap = map[string]bool{
	"job":     true,
	"script":  true,
	"desired": true,
	"status":  true,
}

var AddCmd = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			klog.Error("请输入正确的资源类型")
			os.Exit(1)
		}
		if _, exists := resourceMap[args[0]]; !exists {
			klog.Error("请输入正确的资源类型")
			os.Exit(1)
		}
		if err := ResourceOperate(args[0], "add", global.ResourceName, global.ResourceFilePath); err != nil {
			klog.Error(err)
			os.Exit(1)
		}
	},
}

var UpdateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			klog.Error("请输入正确的资源")
			os.Exit(1)
		}
		if _, exists := resourceMap[args[0]]; !exists {
			klog.Error("请输入正确的资源")
			os.Exit(1)
		}
		if err := ResourceOperate(args[0], "update", global.ResourceName, global.ResourceFilePath); err != nil {
			klog.Error(err)
			os.Exit(1)
		}
	},
}

var GetCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			klog.Error("请输入正确的资源")
			os.Exit(1)
		}
		if _, exists := resourceMap[args[0]]; !exists {
			klog.Error("请输入正确的资源")
			os.Exit(1)
		}
		if err := ResourceOperate(args[0], "get", "", ""); err != nil {
			klog.Error(err)
			os.Exit(1)
		}
	},
}

var DeleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			klog.Error("请输入正确的资源")
			os.Exit(1)
		}
		if _, exists := resourceMap[args[0]]; !exists {
			klog.Error("请输入正确的资源")
			os.Exit(1)
		}
		/*if err := ResourceOperate(args[0], "delete", ""); err != nil {
			klog.Error(err)
			os.Exit(1)
		}*/
	},
}

var RunCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}
