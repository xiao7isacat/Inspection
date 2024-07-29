package config

import (
	"fmt"
	pconfig "github.com/prometheus/common/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	LogLevel  string        `yaml:"log_level"`
	CheckJobC *CheckJobConf `yaml:"check_job_info"`
	CmdInfo   *CmdInfoConf  `yaml:"cmd_info"`
	HttpAddr  string        `yaml:"http_addr"`
}

type CmdInfoConf struct {
	AiK8sFirstPath            string                    `yaml:"ai_k8s_first_path"`
	AiK8sClusterExpose        bool                      `yaml:"ai_k8s_cluster_expose"`
	ApiAddr                   string                    `yaml:"api_addr"`
	HttpClientConfigBasicAuth *pconfig.HTTPClientConfig `yaml:"http_config"`
	HttpClientConfig          *pconfig.HTTPClientConfig `yaml:"-"`
	TimeOutSecond             int                       `yaml:"time_out_second"`
}

type CheckJobConf struct {
	RunCheckJobBatch             int     `yaml:"run_check_job_batch"`
	RunHostBatch                 int     `yaml:"run_host_batch"`
	JobComplateMinutes           float64 `yaml:"job_complate_minutes"`
	CheckServerAddr              string  `yaml:"check_server_addr"`
	NodeRunCheckdir              string  `yaml:"node_run_checkdir"`
	AgentBinDownloadAddr         string  `yaml:"agent_bin_download_addr"`
	CheckSubitJobIntervalSeconds int     `yaml:"check_subit_job_interval_seconds"`
	ComplateJobIntervalSeconds   int     `yaml:"complate_job_interval_seconds"`
	MetricsJobIntervalSeconds    int     `yaml:"metrics_job_interval_seconds"`
}

func LoadFile(filename string) (*Config, error) {
	cfg := &Config{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, cfg); err != nil {

		return nil, fmt.Errorf("解析yaml文件失败,%v", err)
	}
	return cfg, nil
}
