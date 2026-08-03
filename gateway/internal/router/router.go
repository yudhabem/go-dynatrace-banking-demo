package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/handler"
	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/middleware"
)

func New() *gin.Engine {

	r := gin.Default()

	r.Use(middleware.RequestID())

	r.GET("/health", handler.Health)

	return r
}
