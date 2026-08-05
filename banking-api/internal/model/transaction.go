package model

import "time"

type Transaction struct {
	ID            uint   `gorm:"primaryKey"`
	TransactionID string `gorm:"uniqueIndex"`

	FromAccount string
	ToAccount   string
	Amount      float64

	Type     string
	Merchant string

	Status    string
	CreatedAt time.Time
}
