package main

import (
	"context"
	"log"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/config"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/database"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/handler"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/logger"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/observability"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/repository"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/router"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/service"
)

func main() {

	cfg := config.Load()

	// Initialize OpenTelemetry TracerProvider
	tp := observability.Init()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	// Initialize logger
	logg := logger.New()
	defer logg.Sync()

	// Connect database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate
	err = db.AutoMigrate(
		&model.Customer{},
		&model.Account{},
		&model.Transaction{},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Repository
	userRepo := repository.NewUserRepository(db)
	transferRepo := repository.NewTransferRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Service
	userService := service.NewUserService(userRepo)
	transferService := service.NewTransferService(transferRepo)
	paymentService := service.NewPaymentService(paymentRepo)

	// Handler
	userHandler := handler.NewUserHandler(userService)
	transferHandler := handler.NewTransferHandler(transferService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Router
	r := router.New(
		cfg,
		db,
		userHandler,
		transferHandler,
		paymentHandler,
	)

	logg.Info("Banking API started")

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
