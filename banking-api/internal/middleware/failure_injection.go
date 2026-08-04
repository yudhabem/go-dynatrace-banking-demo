package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/logger"
)

const (
	simulateHeader = "X-Simulate"
	slowDelay      = 5 * time.Second
)

// FailureInjection simulates latency or an internal-server error for demo traffic.
// Supported values for X-Simulate are "slow" and "error".
func FailureInjection() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get(RequestIDKey)

		switch c.GetHeader(simulateHeader) {
		case "slow":
			logger.Log.Warn(
				"simulated slow request",
				zap.Any("requestId", requestID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Duration("delay", slowDelay),
			)
			time.Sleep(slowDelay)
		case "error":
			logger.Log.Warn(
				"simulated request failure",
				zap.Any("requestId", requestID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "simulated internal server error",
			})
			return
		}

		c.Next()
	}
}
