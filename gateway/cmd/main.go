package main

import (
	"log"

	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/config"
	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/database"
	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/logger"
	"github.com/yudhabem/go-dynatrace-banking-demo/gateway/internal/router"
)

func main() {

	cfg := config.Load()

	logg := logger.New()
	defer logg.Sync()

	_, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := router.New()

	logg.Info("Gateway started")

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
