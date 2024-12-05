package utils

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
	c.Abort()
}
