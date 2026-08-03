package repository

import (
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"gorm.io/gorm"
)

type TransferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) FindAccount(account string) (*model.Account, error) {
	var acc model.Account

	err := r.db.
		Where("account_number = ?", account).
		First(&acc).Error

	return &acc, err
}

func (r *TransferRepository) ExecuteTransfer(
	from *model.Account,
	to *model.Account,
	transaction *model.Transaction,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(from).Error; err != nil {
			return err
		}

		if err := tx.Save(to).Error; err != nil {
			return err
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *TransferRepository) GetAccount(account string) (*model.Account, error) {
	var acc model.Account

	err := r.db.
		Where("account_number = ?", account).
		First(&acc).Error

	return &acc, err
}

func (r *TransferRepository) GetTransactions() ([]model.Transaction, error) {
	var trx []model.Transaction

	err := r.db.
		Order("created_at DESC").
		Find(&trx).Error

	return trx, err
}
