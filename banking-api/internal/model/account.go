package model

import "time"

type Account struct {
	ID            uint   `gorm:"primaryKey"`
	AccountNumber string `gorm:"uniqueIndex"`
	CustomerID    string
	Balance       float64
	CreatedAt     time.Time
}
