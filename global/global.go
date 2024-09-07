package global

const (
	CtlVersion    = "v1.0.1"
	ServerVersion = "v1.0.1"
	AgentBinName  = "node-env-check-agent"
	AgentVersion  = "v1.0.1"
)

var (
	ServerAddr string
	//server
	ConfigFile        string
	Database          string
	SubmitJobYamlPath string
	//ctl
	ResourceName     string
	ResourceFilePath string
	//agent
	ScriptPath         string
	ResultPath         string
	JobId              int64
	ExecTimeoutSeconds int
	NodeIp             string
)
