package scenario

import (
	"fmt"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
)

func RunHistory(client *client.Client, count int) {

	fmt.Println()
	fmt.Println("========== HISTORY ==========")

	for i := 0; i < count; i++ {

		err := client.History()

		if err != nil {

			fmt.Println("[HISTORY] FAILED")
			continue
		}

		fmt.Println("[HISTORY] SUCCESS")
	}
}
