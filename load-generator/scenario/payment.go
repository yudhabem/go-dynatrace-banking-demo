package scenario

import (
	"fmt"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/util"
)

func RunPayment(client *client.Client, users []model.User, count int) {

	fmt.Println()
	fmt.Println("========== PAYMENT ==========")

	for i := 0; i < count; i++ {

		user := util.RandomUser(users)

		merchant := util.RandomMerchant()

		amount := util.RandomAmount()

		err := client.Payment(
			user.AccountNumber,
			merchant,
			amount,
		)

		if err != nil {

			fmt.Printf(
				"[PAYMENT] %s FAILED\n",
				user.AccountNumber,
			)

			continue
		}

		fmt.Printf(
			"[PAYMENT] %s -> %s | %.0f SUCCESS\n",
			user.AccountNumber,
			merchant,
			amount,
		)
	}
}
