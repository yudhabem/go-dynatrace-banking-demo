package scenario

import (
	"fmt"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/util"
)

func RunInquiry(client *client.Client, users []model.User, count int) {

	fmt.Println()
	fmt.Println("========== INQUIRY ==========")

	for i := 0; i < count; i++ {

		user := util.RandomUser(users)

		err := client.Inquiry(user.AccountNumber)

		if err != nil {
			fmt.Printf("[INQUIRY] %s FAILED\n", user.AccountNumber)
			continue
		}

		fmt.Printf("[INQUIRY] %s SUCCESS\n", user.AccountNumber)
	}
}
