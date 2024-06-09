package v1

import (
	"log/slog"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	slog.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{statusCode, message})
}