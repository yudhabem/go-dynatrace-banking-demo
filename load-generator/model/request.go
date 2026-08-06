package model

type TransferRequest struct {
	FromAccount string  `json:"fromAccount"`
	ToAccount   string  `json:"toAccount"`
	Amount      float64 `json:"amount"`
}

type PaymentRequest struct {
	AccountNumber string  `json:"accountNumber"`
	Merchant      string  `json:"merchant"`
	Amount        float64 `json:"amount"`
}
