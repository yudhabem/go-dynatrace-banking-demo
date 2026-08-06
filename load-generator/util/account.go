package util

import "github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"

func RandomTransfer(users []model.User) (model.User, model.User) {

	from := RandomUser(users)

	to := RandomUser(users)

	for from.AccountNumber == to.AccountNumber {
		to = RandomUser(users)
	}

	return from, to
}
