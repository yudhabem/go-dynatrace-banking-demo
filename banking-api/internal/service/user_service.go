package service

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"

	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/dto"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateRandomUser() (*model.Customer, *model.Account, error) {

	id := uuid.New().String()[:8]

	customer := &model.Customer{
		CustomerID: fmt.Sprintf("CUS-%s", id),
		Name:       fmt.Sprintf("User-%04d", rand.Intn(9999)),
		Email:      fmt.Sprintf("%s@mail.com", id),
		Phone:      fmt.Sprintf("0812%08d", rand.Intn(99999999)),
	}

	account := &model.Account{
		AccountNumber: fmt.Sprintf("1000%06d", rand.Intn(999999)),
		CustomerID:    customer.CustomerID,
		Balance:       float64(rand.Intn(100000000)),
	}

	if err := s.repo.CreateCustomer(customer); err != nil {
		return nil, nil, err
	}

	if err := s.repo.CreateAccount(account); err != nil {
		return nil, nil, err
	}

	return customer, account, nil
}

func (s *UserService) GetAllUsers() ([]model.Customer, error) {
	return s.repo.GetCustomers()
}

func (s *UserService) RandomLogin() (*dto.LoginResponse, error) {

	customer, account, err := s.repo.GetRandomUser()
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		CustomerID:    customer.CustomerID,
		Name:          customer.Name,
		AccountNumber: account.AccountNumber,
		Balance:       account.Balance,
	}, nil
}
