package repository

import (
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) GetAccount(account string) (*model.Account, error) {

	var acc model.Account

	err := r.db.
		Where("account_number = ?", account).
		First(&acc).Error

	return &acc, err
}

func (r *PaymentRepository) ExecutePayment(
	account *model.Account,
	transaction *model.Transaction,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(account).Error; err != nil {
			return err
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
}
