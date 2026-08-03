package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/dto"
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

	from, err := s.repo.FindAccount(req.FromAccount)
	if err != nil {
		return "", err
	}

	to, err := s.repo.FindAccount(req.ToAccount)
	if err != nil {
		return "", err
	}

	if from.Balance < req.Amount {
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
		return "", err
	}

	return trx.TransactionID, nil
}

func (s *TransferService) Inquiry(account string) (*model.Account, error) {
	return s.repo.GetAccount(account)
}

func (s *TransferService) History() ([]model.Transaction, error) {
	return s.repo.GetTransactions()
}
