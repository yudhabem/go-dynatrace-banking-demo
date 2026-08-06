package scenario

import (
	"fmt"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/util"
)

func RunTransfer(client *client.Client, users []model.User, count int) {

	fmt.Println()
	fmt.Println("========== TRANSFER ==========")

	for i := 0; i < count; i++ {

		from, to := util.RandomTransfer(users)

		amount := util.RandomAmount()

		err := client.Transfer(
			from.AccountNumber,
			to.AccountNumber,
			amount,
		)

		if err != nil {

			fmt.Printf(
				"[TRANSFER] %s -> %s FAILED\n",
				from.AccountNumber,
				to.AccountNumber,
			)

			continue
		}

		fmt.Printf(
			"[TRANSFER] %s -> %s | %.0f SUCCESS\n",
			from.AccountNumber,
			to.AccountNumber,
			amount,
		)
	}
}
