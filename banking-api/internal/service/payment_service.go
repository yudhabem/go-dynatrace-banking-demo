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

type PaymentService struct {
	repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) Payment(req dto.PaymentRequest) (string, error) {
	logger.Log.Info(
		"payment started",
		zap.String("account", req.AccountNumber),
		zap.String("merchant", req.Merchant),
		zap.Float64("amount", req.Amount),
	)

	account, err := s.repo.GetAccount(req.AccountNumber)
	if err != nil {
		logger.Log.Error(
			"payment account not found",
			zap.String("account", req.AccountNumber),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"payment account loaded",
		zap.String("account", account.AccountNumber),
		zap.Float64("balance", account.Balance),
	)

	if account.Balance < req.Amount {
		logger.Log.Warn(
			"insufficient balance",
			zap.String("account", account.AccountNumber),
			zap.Float64("balance", account.Balance),
			zap.Float64("requested", req.Amount),
		)

		return "", errors.New("insufficient balance")
	}

	account.Balance -= req.Amount

	logger.Log.Info(
		"payment balance deducted",
		zap.String("account", account.AccountNumber),
		zap.Float64("amount", req.Amount),
		zap.Float64("balance", account.Balance),
	)

	trx := &model.Transaction{
		TransactionID: fmt.Sprintf("PAY-%s", uuid.New().String()[:8]),
		FromAccount:   account.AccountNumber,
		Amount:        req.Amount,
		Type:          "PAYMENT",
		Merchant:      req.Merchant,
		Status:        "SUCCESS",
	}

	if err := s.repo.ExecutePayment(account, trx); err != nil {
		logger.Log.Error(
			"payment failed",
			zap.String("transactionId", trx.TransactionID),
			zap.Error(err),
		)

		return "", err
	}

	logger.Log.Info(
		"payment completed",
		zap.String("transactionId", trx.TransactionID),
		zap.String("account", trx.FromAccount),
		zap.String("merchant", trx.Merchant),
		zap.Float64("amount", trx.Amount),
		zap.String("status", trx.Status),
	)

	return trx.TransactionID, nil
}
