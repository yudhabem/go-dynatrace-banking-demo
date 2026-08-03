package repository

import (
	"github.com/yudhabem/go-dynatrace-banking-demo/banking-api/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateCustomer(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *UserRepository) CreateAccount(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *UserRepository) GetCustomers() ([]model.Customer, error) {

	var customers []model.Customer

	err := r.db.Find(&customers).Error

	return customers, err
}
