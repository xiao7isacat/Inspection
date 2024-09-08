package config

import (
	"fmt"
	pconfig "github.com/prometheus/common/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	LogLevel  string        `yaml:"log_level" json:"log_level"`
	CheckJobC *CheckJobConf `yaml:"check_job_info" json:"check_job_c"`
	CmdInfo   *CmdInfoConf  `yaml:"cmd_info" json:"cmd_info"`
	HttpAddr  string        `yaml:"http_addr" json:"http_addr"`
}

type CmdInfoConf struct {
	AiK8sFirstPath            string                    `yaml:"ai_k8s_first_path" json:"ai_k_8_s_first_path"`
	AiK8sClusterExpose        bool                      `yaml:"ai_k8s_cluster_expose" json:"ai_k_8_s_cluster_expose"`
	ApiAddr                   string                    `yaml:"api_addr" json:"api_addr"`
	HttpClientConfigBasicAuth *pconfig.HTTPClientConfig `yaml:"http_config" json:"http_config"`
	HttpClientConfig          *pconfig.HTTPClientConfig `yaml:"-"`
	TimeOutSecond             int                       `yaml:"time_out_second"`
}

type CheckJobConf struct {
	RunCheckJobBatch              int    `yaml:"run_check_job_batch" json:"run_check_job_batch"`
	RunHostBatch                  int64  `yaml:"run_host_batch" json:"run_host_batch"`
	JobCompleteMinutes            int64  `yaml:"job_complete_minutes" json:"job_complete_minutes"`
	CheckServerAddr               string `yaml:"check_server_addr" json:"check_server_addr"`
	NodeRunCheckDir               string `yaml:"node_run_check_dir" json:"node_run_check_dir"`
	AgentBinDownloadAddr          string `yaml:"agent_bin_download_addr" json:"agent_bin_download_addr"`
	CheckSubmitJobIntervalSeconds int64  `yaml:"check_submit_job_interval_seconds" json:"check_submit_job_interval_seconds"`
	CompleteJobIntervalSeconds    int64  `yaml:"complete_job_interval_seconds" json:"complete_job_interval_seconds"`
	MetricsJobIntervalSeconds     int64  `yaml:"metrics_job_interval_seconds" json:"metrics_job_interval_seconds"`
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
