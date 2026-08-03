package model

import "time"

type Customer struct {
	ID         uint   `gorm:"primaryKey"`
	CustomerID string `gorm:"uniqueIndex"`
	Name       string
	Email      string
	Phone      string
	CreatedAt  time.Time
}
