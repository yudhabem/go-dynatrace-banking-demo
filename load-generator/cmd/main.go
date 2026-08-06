package main

import (
	"fmt"
	"log"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/client"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/config"
	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/scenario"
)

func main() {

	fmt.Println("======================================")
	fmt.Println(" Dynatrace Banking Load Generator")
	fmt.Println("======================================")

	cfg := config.Load()

	c := client.New(cfg)

	users, err := c.LoadUsers()
	if err != nil {
		log.Fatalf("failed to load users: %v", err)
	}

	fmt.Printf("Loaded %d users\n", len(users))

	scenario.Run(c, users)

	fmt.Println()
	fmt.Println("======================================")
	fmt.Println(" Load Generator Finished")
	fmt.Println("======================================")
}
