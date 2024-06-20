package scheduledjob

import (
	database "jobscheduler/internal/db"
	"jobscheduler/internal/models"
	"net/http"

	"github.com/robfig/cron/v3"
)

var cGetJob *cron.Cron
var cRunScheduledJob *cron.Cron

func GetAndRunJobsOnLoop() {
	cGetJob = cron.New()
	cGetJob.AddFunc("@every 5m", func() {
		jobs := database.GetJobs()
		RunScheduledJob(jobs)
	})
	cGetJob.Start()
}

func RunScheduledJob(jobs []models.JobDetail) {
	if cRunScheduledJob != nil {
		cRunScheduledJob.Stop()
		cRunScheduledJob = nil
	}
	cRunScheduledJob = cron.New()
	for _, job := range jobs {
		job := job
		if job.IsScheduled {
			cRunScheduledJob.AddFunc(job.CronSchedule, func() {
				println(job.JobName)
				http.Post("http://localhost:8082/jobs/"+job.JobName+"/run", "application/json", nil)
			})
		}
	}
	cRunScheduledJob.Start()
}
