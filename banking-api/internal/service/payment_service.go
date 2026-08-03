package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/dto"
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

	account, err := s.repo.GetAccount(req.AccountNumber)
	if err != nil {
		return "", err
	}

	if account.Balance < req.Amount {
		return "", errors.New("insufficient balance")
	}

	account.Balance -= req.Amount

	trx := &model.Transaction{
		TransactionID: fmt.Sprintf("PAY-%s", uuid.New().String()[:8]),
		FromAccount:   account.AccountNumber,
		Amount:        req.Amount,
		Type:          "PAYMENT",
		Merchant:      req.Merchant,
		Status:        "SUCCESS",
	}

	if err := s.repo.ExecutePayment(account, trx); err != nil {
		return "", err
	}

	return trx.TransactionID, nil
}
