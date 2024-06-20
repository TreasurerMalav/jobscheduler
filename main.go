package main

import (
	"jobscheduler/internal/handlers"
	"jobscheduler/pkg/scheduledjob"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.GET("/jobs", handlers.GetJobs)
	router.POST("/jobs", handlers.CreateJob)
	router.GET("/jobs/:job_name", handlers.GetJobByName)
	router.PATCH("/jobs/:job_name", handlers.UpdateJob)
	router.DELETE("/jobs/:job_name", handlers.DeleteJob)
	router.POST("/jobs/:job_name/run", handlers.RunJob)
	router.GET("/jobs/execution/history", handlers.GetExecutionHistory)
	scheduledjob.GetAndRunJobsOnLoop()
	router.Run("localhost:8082")
}
