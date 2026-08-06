package scenario

import (
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"
)

func Run(client *client.Client, users []model.User) {

	RunLogin(client, 10)

	RunInquiry(client, users, 15)

	RunTransfer(client, users, 5)

	RunPayment(client, users, 10)

	RunHistory(client, 8)
}
