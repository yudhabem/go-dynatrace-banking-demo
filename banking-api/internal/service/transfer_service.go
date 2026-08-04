package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/dto"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/logger"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/repository"
)

type TransferService struct {
	repo *repository.TransferRepository
}

func NewTransferService(repo *repository.TransferRepository) *TransferService {
	return &TransferService{
		repo: repo,
	}
}

func (s *TransferService) Transfer(req dto.TransferRequest) (string, error) {
	logger.Log.Info(
		"transfer started",
		zap.String("fromAccount", req.FromAccount),
		zap.String("toAccount", req.ToAccount),
		zap.Float64("amount", req.Amount),
	)

	from, err := s.repo.FindAccount(req.FromAccount)
	if err != nil {
		logger.Log.Error(
			"source account not found",
			zap.String("account", req.FromAccount),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"source account loaded",
		zap.String("account", from.AccountNumber),
		zap.Float64("balance", from.Balance),
	)

	to, err := s.repo.FindAccount(req.ToAccount)
	if err != nil {
		logger.Log.Error(
			"destination account not found",
			zap.String("account", req.ToAccount),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"destination account loaded",
		zap.String("account", to.AccountNumber),
		zap.Float64("balance", to.Balance),
	)

	if from.Balance < req.Amount {
		logger.Log.Warn(
			"insufficient balance",
			zap.String("account", from.AccountNumber),
			zap.Float64("balance", from.Balance),
			zap.Float64("requested", req.Amount),
		)

		return "", errors.New("insufficient balance")
	}

	from.Balance -= req.Amount
	to.Balance += req.Amount

	trx := &model.Transaction{
		TransactionID: fmt.Sprintf("TRX-%s", uuid.New().String()[:8]),
		FromAccount:   from.AccountNumber,
		ToAccount:     to.AccountNumber,
		Amount:        req.Amount,
		Status:        "SUCCESS",
	}

	err = s.repo.ExecuteTransfer(from, to, trx)
	if err != nil {
		logger.Log.Error(
			"transfer failed",
			zap.String("transactionId", trx.TransactionID),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"transfer completed",
		zap.String("transactionId", trx.TransactionID),
		zap.String("fromAccount", trx.FromAccount),
		zap.String("toAccount", trx.ToAccount),
		zap.Float64("amount", trx.Amount),
		zap.String("status", trx.Status),
	)

	return trx.TransactionID, nil
}

func (s *TransferService) Inquiry(account string) (*model.Account, error) {
	return s.repo.GetAccount(account)
}

func (s *TransferService) History() ([]model.Transaction, error) {
	return s.repo.GetTransactions()
}
