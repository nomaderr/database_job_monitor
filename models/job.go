package models

type Job struct {
	Hostname     string `json:"hostname"`
	DatabaseName string `json:"database_name"`
	JobName      string `json:"job_name"`
	Status       string `json:"status"`
	LastRun      string `json:"last_run"`
}

type DBConnection struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}
