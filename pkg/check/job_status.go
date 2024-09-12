package check

type JobStatus struct {
	JobName     string `json:"name"`
	CheckName   string `json:"check_name"`
	ActualValue string `json:"actual_value"`
	ExpectValue string `json:"expect_value"`
	Node        string `json:"node"`
	Status      string `json:"status"`
	AllDone     bool   `json:"all_done"`
}
