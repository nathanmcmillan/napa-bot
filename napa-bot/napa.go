package main

import (
	"fmt"
	"time"
	"./gdax"
	"./analyst"
)

func main() {
	fmt.Println("napa bot")
	product := "LTC-USD"
	history := gdax.GetHistory(product, "2018-02-16", "2018-02-17", "3600")
	rsi := analyst.RelativeStrengthIndex(history)
	fmt.Println("RSI ", rsi)
	
	/* gdax.GetCurrencies()
	gdax.GetBook(product)
	gdax.GetTicker(product)
	gdax.GetTrades(product)
 	gdax.GetHistory(product, "2018-02-16", "2018-02-17", "3600")
	gdax.GetStats(product) */
}

func sleep(seconds int32) {
	time.Sleep(time.Second * time.Duration(seconds))
}