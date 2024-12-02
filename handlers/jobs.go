package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database_app/models"
	"database_app/utils"
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Получаем session_id из заголовка
	sessionID := r.Header.Get("Session-ID")
	if sessionID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Session-ID header")
		return
	}

	// Извлекаем данные подключения из Redis
	ctx := context.Background()
	sessionData, err := utils.RedisClient.Get(ctx, sessionID).Result()
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
		log.Println("Redis get error:", err)
		return
	}

	// Десериализуем данные подключения из JSON
	var conn models.DBConnection
	if err := json.Unmarshal([]byte(sessionData), &conn); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to parse session data")
		log.Println("JSON unmarshal error:", err)
		return
	}

	// Формируем строку DSN
	currentDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conn.Username, conn.Password, conn.Hostname, conn.Port, conn.Database)

	// Подключаемся к базе данных
	db, err := sql.Open("mysql", currentDSN)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to connect to database")
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	// Выполняем SQL-запрос
	var hostname string
	err = db.QueryRow("SELECT @@hostname").Scan(&hostname)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch hostname")
		log.Println("Error fetching hostname:", err)
		return
	}

	rows, err := db.Query(`SELECT DATABASE() AS database_name, name AS job_name, status, last_run FROM jobs`)
	if err != nil {
		log.Println("Error executing query:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch jobs")
		return
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		job.Hostname = hostname
		if err := rows.Scan(&job.DatabaseName, &job.JobName, &job.Status, &job.LastRun); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process job data")
			log.Println("Error scanning job:", err)
			return
		}
		jobs = append(jobs, job)
	}

	if len(jobs) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No jobs found in the database")
		return
	}

	json.NewEncoder(w).Encode(jobs)
}
