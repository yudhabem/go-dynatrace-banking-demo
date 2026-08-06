package scenario

import (
	"fmt"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
)

func RunLogin(client *client.Client, count int) {

	fmt.Println()
	fmt.Println("========== LOGIN ==========")

	for i := 0; i < count; i++ {

		err := client.Login()

		if err != nil {
			fmt.Printf("[LOGIN %d] FAILED : %v\n", i+1, err)
			continue
		}

		fmt.Printf("[LOGIN %d] SUCCESS\n", i+1)
	}
}
