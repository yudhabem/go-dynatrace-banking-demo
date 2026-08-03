package dto

type TransferRequest struct {
	FromAccount string  `json:"fromAccount" binding:"required"`
	ToAccount   string  `json:"toAccount" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}
