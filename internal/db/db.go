package database

import (
	"database/sql"
	"jobscheduler/internal/models"
	"time"
)

type DatabaseCreds struct {
	User     string
	Password string
}

var db *sql.DB

func NewDatabaseCreds() *DatabaseCreds {
	return &DatabaseCreds{
		User:     "jobsadmin",
		Password: "******",
	}
}

func DBConnection(user, password, dbname string) *sql.DB {
	// connection to the database
	db, err := sql.Open("mysql", user+":"+password+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	// defer db.Close()

	return db
}

func InsertJob(job models.Job) error {
	// insert a new job a to the table jobs_details
	//var job models.Job
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()
	_, err := db.Query("INSERT INTO jobs_details VALUES(?, ?, ?, ?, ?, ?, now(), now())", job.Job_name, job.Job_command, job.Is_scheduled, job.Cron_schedule, job.Hosts, NewDatabaseCreds().User)
	defer db.Close()

	return err
}

func GetJobs() []models.JobDetail {
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM jobs_details")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var jobs []models.JobDetail
	for rows.Next() {
		var job models.JobDetail
		err := rows.Scan(&job.JobName, &job.JobCommand, &job.IsScheduled, &job.CronSchedule, &job.Hosts, &job.User, &job.CreationTimestamp, &job.LatestUpdationTimestamp)
		if err != nil {
			panic(err)
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func GetJobByName(jobName string) (models.JobDetail, error) {
	var job models.JobDetail
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()
	err := db.QueryRow("SELECT * FROM jobs_details WHERE Job_name = ?", jobName).Scan(&job.JobName, &job.JobCommand, &job.IsScheduled, &job.CronSchedule, &job.Hosts, &job.User, &job.CreationTimestamp, &job.LatestUpdationTimestamp)
	// if err != nil {
	// 	panic(err)
	// }

	return job, err
}

func UpdateJob(jobName string, job models.Job) (sql.Result, error) {
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()
	stmt, err := db.Prepare(`
	UPDATE jobs_details 
	SET job_command = ?, is_scheduled = ?, cron_schedule = ?, hosts = ?, latest_updation_timestamp = now()
	WHERE Job_name = ?
	`)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(job.Job_command, job.Is_scheduled, job.Cron_schedule, job.Hosts, jobName)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}

	return res, err
}

func DeleteJob(job string) (sql.Result, error) {
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()

	result, err := db.Exec("DELETE FROM jobs_details WHERE Job_name = ?", job)

	return result, err

}

func InsertExecutionDetails(job_name string, start_time time.Time, end_time time.Time, status string) error {
	// insert a new job a to the table jobs_details
	//var job models.Job
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()
	if end_time == (time.Time{}) {
		_, err := db.Query("INSERT INTO jobs_execution_history (job_name, start_time, status) VALUES(?, ?, ?)", job_name, start_time, status)
		defer db.Close()
		return err
	} else {
		_, err := db.Query("INSERT INTO jobs_execution_history (job_name, start_time, end_time, status) VALUES(?, ?, ?, ?)", job_name, start_time, end_time, status)
		defer db.Close()
		return err

	}

}

func GetExecutionHistory() []models.JobsExecutionHistory {
	DatabaseCreds := NewDatabaseCreds()
	db = DBConnection(DatabaseCreds.User, DatabaseCreds.Password, "jobs")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM jobs_execution_history")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var jobs_execution_history []models.JobsExecutionHistory
	for rows.Next() {
		var job models.JobsExecutionHistory
		err := rows.Scan(&job.ExecutionId, &job.JobName, &job.StartTime, &job.EndTime, &job.Status)
		if err != nil {
			panic(err)
		}
		jobs_execution_history = append(jobs_execution_history, job)
	}
	return jobs_execution_history
}
