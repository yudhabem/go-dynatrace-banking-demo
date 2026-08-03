package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/response"
)

func Health(c *gin.Context) {

	response.Success(c, gin.H{
		"status":  "UP",
		"service": "api-gateway",
		"version": "1.0.0",
	})
}
