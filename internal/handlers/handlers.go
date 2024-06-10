package handlers

import (
	"database/sql" // Add this line

	database "jobscheduler/internal/db"
	"jobscheduler/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// var user = "jobsadmin"

func CreateJob(c *gin.Context) {
	var newJob models.Job

	if err := c.BindJSON(&newJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.InsertJob(newJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Job created successfully"})
	}

}

func GetJobs(c *gin.Context) {
	jobs := database.GetJobs()
	c.JSON(http.StatusOK, jobs)

}

func GetJobByName(c *gin.Context) {
	jobName := c.Param("job_name")
	job, err := database.GetJobByName(jobName)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	} else {
		c.JSON(http.StatusOK, job)
	}

}

func UpdateJob(c *gin.Context) {
	jobName := c.Param("job_name")

	var updateJob models.Job
	if err := c.BindJSON(&updateJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.UpdateJob(jobName, updateJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Job updated successfully"})

}

func DeleteJob(c *gin.Context) {
	deleteJobName := c.Param("job_name")
	result, err := database.DeleteJob(deleteJobName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})

}
