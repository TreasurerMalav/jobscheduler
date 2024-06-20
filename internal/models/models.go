package models

type Job struct {
	Job_name      string `json:"job_name"`
	Job_command   string `json:"job_command"`
	Is_scheduled  bool   `json:"is_scheduled"`
	Cron_schedule string `json:"cron_schedule"` // null if is_schedules is false
	Hosts         string `json:"hosts"`
}

type JobDetail struct {
	JobName                 string `json:"job_name"`
	JobCommand              string `json:"job_command"`
	IsScheduled             bool   `json:"is_scheduled"`
	CronSchedule            string `json:"cron_schedule"`
	Hosts                   string `json:"hosts"`
	User                    string `json:"user"`
	CreationTimestamp       string `json:"creation_timestamp"`
	LatestUpdationTimestamp string `json:"latest_updation_timestamp"`
}

type JobsExecutionHistory struct {
	ExecutionId string `json:"execution_id"`
	JobName     string `json:"job_name"`
	Status      string `json:"status"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}
