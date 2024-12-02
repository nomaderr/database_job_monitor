package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"database_app/models"
	"database_app/utils"

	_ "github.com/go-sql-driver/mysql"
)

var (
	connLock sync.Mutex
	dsn      string
	dbConn   models.DBConnection
)

// ConnectHandler handles database connection based on user input
func ConnectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var conn models.DBConnection
	if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON input")
		log.Println("Error decoding JSON:", err)
		return
	}

	if conn.Hostname == "" || conn.Port == "" || conn.Username == "" || conn.Password == "" || conn.Database == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "All fields are required")
		log.Println("Missing fields in connection data")
		return
	}

	// Формируем строку подключения
	newDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conn.Username, conn.Password, conn.Hostname, conn.Port, conn.Database)
	log.Println("Attempting to connect with DSN:", newDSN)

	// Проверяем подключение к базе данных
	db, err := sql.Open("mysql", newDSN)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to parse connection string")
		log.Println("Error parsing connection string:", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to connect to database")
		log.Println("Database connection error:", err)
		return
	}

	// Сохраняем сессию в Redis
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
	sessionData, err := json.Marshal(conn) // Сериализация структуры в JSON
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to serialize session data")
		log.Println("JSON marshal error:", err)
		return
	}

	ctx := context.Background()
	if err := utils.RedisClient.Set(ctx, sessionID, sessionData, 30*time.Minute).Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save session")
		log.Println("Redis set error:", err)
		return
	}

	log.Println("Connected to database:", conn.Hostname)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Connected successfully", "session_id": sessionID})
}
