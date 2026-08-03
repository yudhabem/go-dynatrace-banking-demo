package router

import (
	"github.com/gin-gonic/gin"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/handler"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/middleware"
)

func New(
	userHandler *handler.UserHandler,
	transferHandler *handler.TransferHandler,
	paymentHandler *handler.PaymentHandler,
) *gin.Engine {

	r := gin.Default()

	r.Use(middleware.RequestID())

	r.GET("/health", handler.Health)

	r.POST("/users/random", userHandler.Random)
	r.GET("/users", userHandler.GetAll)

	r.POST("/login/random", userHandler.Login)

	r.POST("/transfer", transferHandler.Transfer)

	// Tambahkan dua route ini
	r.GET("/accounts/:account", transferHandler.Inquiry)
	r.GET("/transactions", transferHandler.History)

	r.POST("/payment", paymentHandler.Payment)

	return r
}
