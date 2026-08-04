package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/logger"
)

func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {

		defer func() {

			if err := recover(); err != nil {

				requestID, _ := c.Get(RequestIDKey)

				logger.Log.Error(
					"panic recovered",
					zap.Any("requestId", requestID),
					zap.Any("panic", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.ByteString("stack", debug.Stack()),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "internal server error",
				})
			}

		}()

		c.Next()
	}
}
