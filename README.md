Database Job Monitor

Project Overview

The Database Job Monitor is a lightweight microservice application designed to connect to a MySQL database, monitor scheduled jobs, and provide insights through a user-friendly web interface. It also integrates with Redis for session management and Prometheus for job alerting. The backend is implemented in Go using the Gin web framework.

Features

Database Connection: Connect to a MySQL database with user-provided credentials and interval.

Job Monitoring: Fetch and display job statuses from the database.

Session Management: Store connection details and user-defined intervals in Redis.

Dynamic Intervals: User-defined job monitoring intervals.

Prometheus Integration: Alerts for failed jobs.

Frontend: Simple web interface built with HTML, CSS (TailwindCSS), and JavaScript.

Technologies Used

Backend: Go, Gin Framework

Database: MySQL

Session Storage: Redis

Monitoring: Prometheus

Frontend: TailwindCSS, HTML, JavaScript

Containerization: Docker, Docker Compose

Project Structure

├── Dockerfile
├── README.md
├── docker-compose.yaml
├── go.mod
├── go.sum
├── handlers
│   ├── connect.go
│   ├── html.go
│   └── jobs.go
├── main.go
├── models
│   └── job.go
├── templates
│   └── index.html
└── utils
    ├── redis_client.go
    └── response.go

Setup Instructions

Prerequisites

Docker and Docker Compose installed.

MySQL and Redis are included in the Docker Compose setup.

Build and Run

Clone the repository:

git clone <repository_url>
cd database_job_monitor

Build and start the application:

docker-compose up -d --build

Access the web interface:
Navigate to http://localhost:8080 in your browser.

Configuration

The application uses environment variables for configuration. These are defined in the docker-compose.yaml file:

MYSQL_HOST: MySQL server hostname (default: test-mysql)

MYSQL_PORT: MySQL server port (default: 3306)

MYSQL_USER: MySQL username

MYSQL_PASSWORD: MySQL password

MYSQL_DATABASE: MySQL database name

REDIS_HOST: Redis server hostname (default: redis)

REDIS_PORT: Redis server port (default: 6379)

Endpoints

API Endpoints

POST /api/connect

Connect to the MySQL database with user-provided credentials.

Request Body:

{
  "hostname": "127.0.0.1",
  "port": "3306",
  "username": "root",
  "password": "password",
  "database": "test_db",
  "interval": 30
}

GET /api/jobs

Fetch job statuses from the connected database.

Requires Session-ID header.

Frontend

Static HTML served at /.

Prometheus Integration

The application exposes a /metrics endpoint for Prometheus to scrape. Failed jobs are monitored and exposed as metrics.

Known Issues

Redis and MySQL should be correctly networked and accessible by the database-monitor container.

Ensure that job table structure matches the expected format.

Future Enhancements

Extend support for additional database types.

Implement advanced alerting mechanisms.

Add authentication and authorization for API endpoints.

License

This project is licensed under the MIT License. See the LICENSE file for details.

