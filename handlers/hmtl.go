package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServeHTML serves the frontend HTML page using Gin
func ServeHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
