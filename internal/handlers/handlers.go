package handlers

import (
	"database/sql" // Add this line
	"log"
	"os"
	"strings"
	"time"

	database "jobscheduler/internal/db"
	"jobscheduler/internal/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
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

func RunJob(c *gin.Context) {
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
		hosts := strings.Split(strings.Trim(job.Hosts, "[]"), ",")
		var host_list []string
		for _, host := range hosts {
			host_list = append(host_list, strings.TrimSpace(host))
		}
		var wg sync.WaitGroup
		wg.Add(len(host_list))
		start_time := time.Now()
		// err = database.InsertExecutionDetails(jobName, start_time, time.Time{}, "In Progress")
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }
		job_failed := false
		for _, host := range host_list {
			//hostStr := string(host)
			go func(host string) {
				defer wg.Done()

				key, err := os.ReadFile("path-to-pem-file") // replace with your private key path
				// log.Printf("Key is present for host %s", host)
				if err != nil {
					job_failed = true
					log.Printf("unable to read private key: %v", err)
					return
				}

				signer, err := ssh.ParsePrivateKey(key)
				if err != nil {
					job_failed = true
					log.Printf("unable to parse private key: %v", err)
					return
				}
				config := &ssh.ClientConfig{
					User: "user", // replace with your username
					Auth: []ssh.AuthMethod{
						ssh.PublicKeys(signer),
					},
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					Timeout:         30 * time.Second,
				}
				client, err := ssh.Dial("tcp", host+":22", config)
				if err != nil {
					// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					job_failed = true
					log.Printf("Error occurred while SSH to VM: %v", err)
					return
				}
				defer client.Close()
				session, err := client.NewSession()
				if err != nil {
					// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					job_failed = true
					log.Printf("Error occurred while creating new session: %v", err)
					return
				}
				defer session.Close()

				// replace 'command-to-run-job' with the command to run your job
				err = session.Run(job.JobCommand)
				if err != nil {
					// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					job_failed = true
					log.Printf("Error occurred while running job: %v", err)
					return
				}
			}(host)
		}
		wg.Wait()
		end_time := time.Now()
		if !job_failed {
			c.JSON(http.StatusOK, gin.H{"message": "Job run successfully"})
			err = database.InsertExecutionDetails(jobName, start_time, end_time, "Success")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Job failed"})
			err = database.InsertExecutionDetails(jobName, start_time, end_time, "Failed")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

	}
}

func GetExecutionHistory(c *gin.Context) {
	jobs := database.GetExecutionHistory()
	c.JSON(http.StatusOK, jobs)
}
