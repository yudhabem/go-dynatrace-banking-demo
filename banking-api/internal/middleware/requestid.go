package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {

	return func(c *gin.Context) {

		requestID := uuid.New().String()

		c.Header("X-Request-ID", requestID)

		c.Set("requestID", requestID)

		c.Next()
	}
}
