package router

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/handler"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/middleware"
)

func New(
	userHandler *handler.UserHandler,
	transferHandler *handler.TransferHandler,
	paymentHandler *handler.PaymentHandler,
) *gin.Engine {

	r := gin.Default()

	r.Use(otelgin.Middleware("banking-api"))
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	r.GET("/health", handler.Health)

	r.POST("/users/random", userHandler.Random)
	r.GET("/users", userHandler.GetAll)

	r.POST("/login/random", userHandler.Login)

	r.POST("/transfer", middleware.FailureInjection(), transferHandler.Transfer)

	// Tambahkan dua route ini
	r.GET("/accounts/:account", transferHandler.Inquiry)
	r.GET("/transactions", transferHandler.History)

	r.POST("/payment", middleware.FailureInjection(), paymentHandler.Payment)

	r.GET("/panic", func(c *gin.Context) {
		panic("this is panic test")
	})

	return r
}
