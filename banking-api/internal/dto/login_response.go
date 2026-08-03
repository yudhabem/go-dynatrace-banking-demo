package dto

type LoginResponse struct {
	CustomerID    string  `json:"customerId"`
	Name          string  `json:"name"`
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
}
