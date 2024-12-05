# Database Job Monitor

A lightweight microservice designed to monitor and track database jobs with support for interval-based job checking and Redis-backed session management. The application features a simple web interface built with TailwindCSS and APIs built using Gin framework in Golang.

## Features

- **Database Monitoring:**
  - Connect to MySQL databases and fetch job statuses.
  - Monitor job status (e.g., Running, Completed, Failed).

- **Dynamic Intervals:**
  - Users can set custom intervals for job monitoring.

- **Redis-Backed Session Management:**
  - Store session data securely in Redis.

- **Prometheus Integration:**
  - Export job metrics to Prometheus for monitoring failed jobs.

- **Lightweight Frontend:**
  - Built using HTML and TailwindCSS for responsive UI.

## Prerequisites

- Docker
- Docker Compose
- Golang (if running locally)
- Redis and MySQL

## Setup Instructions

### Clone the Repository
```bash
# Clone the repository
git clone <repository-url>
cd database_job_monitor
```

### Environment Setup

- Create a `.env` file for environment-specific configurations.

```env
REDIS_HOST=redis
REDIS_PORT=6379
MYSQL_HOST=test-mysql
MYSQL_PORT=3306
MYSQL_USER=admin
MYSQL_PASSWORD=admin
MYSQL_DATABASE=test_db
```

### Run with Docker Compose
```bash
# Start services with Docker Compose
docker-compose up -d --build
```

### Folder Structure

```
├── Dockerfile
├── README.md
├── docker-compose.yaml
├── go.mod
├── go.sum
├── handlers
│   ├── connect.go
│   ├── hmtl.go
│   └── jobs.go
├── main.go
├── models
│   └── job.go
├── templates
│   └── index.html
└── utils
    ├── redis_client.go
    └── response.go
```

### Access the Application

1. Open a browser and navigate to `http://localhost:8080`
2. Connect to a MySQL database via the connection form.
3. View job statuses in the table.

## API Endpoints

### `/api/connect`
**POST**: Connect to the database.

- **Body:**
  ```json
  {
    "hostname": "127.0.0.1",
    "port": "3306",
    "username": "admin",
    "password": "admin",
    "database": "test_db",
    "interval": 30
  }
  ```

### `/api/jobs`
**GET**: Fetch job data for a connected session.

- **Headers:**
  - `Session-ID: <session_id>`

### `/metrics`
**GET**: Export metrics to Prometheus.

## Troubleshooting

- Ensure Docker and Docker Compose are installed and running.
- Verify that Redis and MySQL are reachable and properly configured.
- Check logs for error messages:
  ```bash
  docker logs database-monitor
  ```

## Technologies Used

- **Backend:** Golang, Gin Framework
- **Frontend:** HTML, TailwindCSS
- **Database:** MySQL
- **Cache/Session:** Redis
- **Monitoring:** Prometheus

## License
This project is licensed under the MIT License. See the LICENSE file for details.

---

#golang #redis #mysql #prometheus #gin #docker #tailwindcss
