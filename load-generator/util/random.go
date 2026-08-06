package util

import (
	"math/rand"
	"time"

	"github.com/yudhabem/go-dynatrace-banking-demo/load-generator/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomUser(users []model.User) model.User {
	return users[rand.Intn(len(users))]
}

func RandomAmount() float64 {
	return float64(rand.Intn(500000) + 10000)
}
