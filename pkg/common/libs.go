package common

import (
	"context"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"k8s.io/klog/v2"
	"time"

	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
)

func GenEmptyString(input string) string {
	out := "empty_string"
	if input != "" {
		out = input
	}
	return out
}

func AnsiRunPlay(remoteHost string, extraVars map[string]interface{}, ansiYamlPath string) error {
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "smart",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: remoteHost + ",",
		ExtraVars: extraVars,
	}

	lplaybook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         []string{ansiYamlPath},
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithTransformers(
				results.Prepend("prome-shard"),
			),
		),
	}

	klog.Info(lplaybook.String())
	err := lplaybook.Run(ctx)
	return err
}
