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

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	connLock sync.Mutex
	dsn      string
)

// ConnectHandler handles database connection and saves the session in Redis
func ConnectHandler(c *gin.Context) {
	var conn models.DBConnection

	// Bind JSON to the structure
	if err := c.ShouldBindJSON(&conn); err != nil {
		log.Println("Error binding JSON:", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	// Parse interval as int
	interval, err := conn.Interval.Int64()
	if err != nil || interval <= 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid interval value")
		log.Println("Invalid interval:", conn.Interval)
		return
	}

	// Connection string
	newDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conn.Username, conn.Password, conn.Hostname, conn.Port, conn.Database)
	log.Println("Attempting to connect with DSN:", newDSN)

	// Check connect ot DB
	db, err := sql.Open("mysql", newDSN)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Unable to parse connection string")
		log.Println("Error parsing connection string:", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Unable to connect to database")
		log.Println("Database connection error:", err)
		return
	}

	// Save session into Redis
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())
	connData := map[string]interface{}{
		"hostname": conn.Hostname,
		"port":     conn.Port,
		"username": conn.Username,
		"password": conn.Password,
		"database": conn.Database,
		"interval": interval,
	}

	sessionData, err := json.Marshal(connData)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to serialize session data")
		log.Println("JSON marshal error:", err)
		return
	}

	ctx := context.Background()
	if err := utils.RedisClient.Set(ctx, sessionID, sessionData, 30*time.Minute).Err(); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save session")
		log.Println("Redis set error:", err)
		return
	}

	log.Println("Connected to database:", conn.Hostname)
	c.JSON(http.StatusOK, gin.H{
		"message":    "Connected successfully",
		"session_id": sessionID,
	})
}
