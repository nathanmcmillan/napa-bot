package main

import (
	"fmt"

	"./gdax"
)

func main() {
	fmt.Println("napa bot")
	gdax.ListenTicker()
	gdax.GetCurrencies()
}
