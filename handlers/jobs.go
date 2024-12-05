package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"database_app/models"
	"database_app/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics
var (
	failedJobsMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "failed_jobs_count",
			Help: "Number of failed jobs",
		},
		[]string{"database_name", "job_name"},
	)
)

func init() {
	prometheus.MustRegister(failedJobsMetric)
}

// GetJobs fetches jobs from the database for a specific session
func GetJobs(c *gin.Context) {
	sessionID := c.GetHeader("Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Session-ID header"})
		return
	}

	ctx := context.Background()
	sessionData, err := utils.RedisClient.Get(ctx, sessionID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		log.Println("Redis get error:", err)
		return
	}

	var conn models.DBConnection
	if err := json.Unmarshal([]byte(sessionData), &conn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse session data"})
		log.Println("JSON unmarshal error:", err)
		return
	}

	db, err := sql.Open("mysql", conn.Username+":"+conn.Password+"@tcp("+conn.Hostname+":"+conn.Port+")/"+conn.Database)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	var hostname string
	err = db.QueryRow("SELECT @@hostname").Scan(&hostname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hostname"})
		log.Println("Error fetching hostname:", err)
		return
	}

	rows, err := db.Query(`SELECT DATABASE() AS database_name, name AS job_name, status, last_run FROM jobs`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		job.Hostname = hostname // Hostname field
		if err := rows.Scan(&job.DatabaseName, &job.JobName, &job.Status, &job.LastRun); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process job data"})
			log.Println("Error scanning job:", err)
			return
		}

		// Update Prometheus metric for failed jobs
		if job.Status == "Failed" {
			failedJobsMetric.WithLabelValues(job.DatabaseName, job.JobName).Inc()
		}

		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, jobs)
}

// StartJobChecker periodically checks for jobs in all active sessions
// func StartJobChecker() {
// 	go func() {
// 		for {
// 			time.Sleep(10 * time.Second) // Frequency of Redis session checking

// 			ctx := context.Background()
// 			keys, err := utils.RedisClient.Keys(ctx, "session_*").Result()
// 			if err != nil {
// 				log.Println("Error fetching keys from Redis:", err)
// 				continue
// 			}

// 			for _, key := range keys {
// 				sessionData, err := utils.RedisClient.Get(ctx, key).Result()
// 				if err != nil {
// 					log.Println("Error fetching session data:", err)
// 					continue
// 				}

// 				var conn models.DBConnection
// 				if err := json.Unmarshal([]byte(sessionData), &conn); err != nil {
// 					log.Println("Error unmarshalling session data:", err)
// 					continue
// 				}

//					checkJobs(conn)
//				}
//			}
//		}()
//	}
func StartJobChecker() {
	go func() {
		for {
			ctx := context.Background()
			keys, err := utils.RedisClient.Keys(ctx, "session_*").Result()
			if err != nil {
				log.Println("Error fetching keys from Redis:", err)
				time.Sleep(10 * time.Second) // standart value
				continue
			}

			for _, key := range keys {
				sessionData, err := utils.RedisClient.Get(ctx, key).Result()
				if err != nil {
					log.Println("Error fetching session data:", err)
					continue
				}

				var conn models.DBConnection
				if err := json.Unmarshal([]byte(sessionData), &conn); err != nil {
					log.Println("Error unmarshalling session data:", err)
					continue
				}

				// interval from json.Number to int
				interval, err := conn.Interval.Int64()
				if err != nil || interval <= 0 {
					interval = 10 // Default value
				}

				checkJobs(conn)

				time.Sleep(time.Duration(interval) * time.Second) // Dynamic interval
			}
		}
	}()
}

// checkJobs fetches and logs jobs for a specific connection
func checkJobs(conn models.DBConnection) {
	db, err := sql.Open("mysql", conn.Username+":"+conn.Password+"@tcp("+conn.Hostname+":"+conn.Port+")/"+conn.Database)
	if err != nil {
		log.Println("Failed to connect to database for session:", conn.Hostname, err)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT name, status FROM jobs WHERE status = "Failed"`)
	if err != nil {
		log.Println("Error executing query for failed jobs:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var jobName, status string
		if err := rows.Scan(&jobName, &status); err != nil {
			log.Println("Error scanning job data:", err)
			continue
		}
		log.Printf("Job %s is in %s state\n", jobName, status)
	}
}
