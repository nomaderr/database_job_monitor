package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database_app/handlers"
	"database_app/utils"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	utils.InitializeRedis(redisHost, redisPort)

	http.HandleFunc("/connect", handlers.ConnectHandler) // API endpoint to connect to the database
	http.HandleFunc("/jobs", handlers.GetJobs)           // API endpoint to fetch jobs
	http.HandleFunc("/", handlers.ServeHTML)             // Serve the front-end

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
