package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/config"
)

func Version(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"version": cfg.Version,
		})

	}

}
