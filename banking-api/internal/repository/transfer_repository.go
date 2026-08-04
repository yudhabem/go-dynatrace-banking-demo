package repository

import (
	"context"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"gorm.io/gorm"
)

type TransferRepository struct {
	db *gorm.DB
}

func NewTransferRepository(db *gorm.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) FindAccount(ctx context.Context, account string) (*model.Account, error) {
	var acc model.Account

	err := r.db.WithContext(ctx).
		Where("account_number = ?", account).
		First(&acc).Error

	return &acc, err
}

func (r *TransferRepository) ExecuteTransfer(
	ctx context.Context,
	from *model.Account,
	to *model.Account,
	transaction *model.Transaction,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

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

func (r *TransferRepository) GetAccount(ctx context.Context, account string) (*model.Account, error) {
	var acc model.Account

	err := r.db.WithContext(ctx).
		Where("account_number = ?", account).
		First(&acc).Error

	return &acc, err
}

func (r *TransferRepository) GetTransactions(ctx context.Context) ([]model.Transaction, error) {
	var trx []model.Transaction

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&trx).Error

	return trx, err
}
