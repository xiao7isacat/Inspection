package global

const (
	CtlVersion            = "v1.0.1"
	ServerVersion         = "v1.0.1"
	GLOBAL_VAR_JOB_PROBER = "job_prober"
	MasterPodName         = ""
	JobResourceResultUndo = ""
	GLOBAL_VAR_ALL_CONFIG = "ALL_CONFIG"
	AgentBinName          = "node-env-check-agent"
	AgentVersion          = "v1.0.1"
	CheckJobManager       = "CheckJobManager"
)

var (
	//server
	ServerAddr        string
	ConfigFile        string
	Database          string
	SubmitJobYamlPath string
	//ctl
	ResourceName     string
	ResourceFilePath string
	NodeAddrs        string
	//agent
	AgentPort          string
	ScriptPath         string
	ResultPath         string
	JobId              int64
	ExecTimeoutSeconds int
)
