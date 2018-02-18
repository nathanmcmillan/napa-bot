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
	history := gdax.GetHistory(product, "2018-02-16", "2018-02-18", "3600")
	fmt.Println("MACD", analyst.MovingAverageConvergenceDivergence(12, 26, history))
	fmt.Println("RSI", analyst.RelativeStrengthIndex(14, history))
	
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