# Use Go 1.23 for compatibility
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Ensure the file exists
RUN ls -l /app

# Expose the application on port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
