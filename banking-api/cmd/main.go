package main

import (
	"log"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/config"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/database"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/handler"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/logger"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/repository"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/router"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/service"
)

func main() {

	cfg := config.Load()

	logg := logger.New()
	defer logg.Sync()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.Customer{},
		&model.Account{},
		&model.Transaction{},
	)

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepository(db)

	svc := service.NewUserService(repo)

	userHandler := handler.NewUserHandler(svc)

	transferRepo := repository.NewTransferRepository(db)

	transferService := service.NewTransferService(transferRepo)

	transferHandler := handler.NewTransferHandler(transferService)

	r := router.New(
		userHandler,
		transferHandler,
	)

	logg.Info("Gateway started")

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
