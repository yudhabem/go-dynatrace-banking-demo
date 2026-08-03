package dto

type PaymentRequest struct {
	AccountNumber string  `json:"accountNumber" binding:"required"`
	Merchant      string  `json:"merchant" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}
