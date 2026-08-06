package model

type User struct {
	CustomerID    string  `json:"customerId"`
	Name          string  `json:"name"`
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
}

type UserResponse struct {
	Success bool   `json:"success"`
	Data    []User `json:"data"`
}
