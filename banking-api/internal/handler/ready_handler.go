package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Ready(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		sqlDB, err := db.DB()
		if err != nil {

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "DOWN",
			})

			return
		}

		if err := sqlDB.Ping(); err != nil {

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "DOWN",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})

	}

}
