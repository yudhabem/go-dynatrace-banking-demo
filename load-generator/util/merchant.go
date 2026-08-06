package util

import "math/rand"

var merchants = []string{
	"PLN Prepaid",
	"Shopee",
	"Tokopedia",
	"Netflix",
	"Steam",
	"Indomaret",
	"Alfamart",
	"Gojek",
	"Grab",
	"Lazada",
	"PLN Postpaid",
	"Maxim",
	"Telkomsel",
	"XL Axiata",
	"Indosat",
	"By-U",
	"Apple",
	"Youtube",
	"Traveloka",
	"Agoda",
	"Uniqlo",
	"Zara",
	"Pertamina",
	"Spotify",
}

func RandomMerchant() string {
	return merchants[rand.Intn(len(merchants))]
}
