package trader

import (
	"database/sql"
	"fmt"

	"../analyst"
	"../gdax"
)

// Run core loop
func Run(analyst *analyst.Analysis, db *sql.DB, auth *gdax.Authentication, messages chan interface{}) {

	//var rsi []float64
	//var macd []float64
	//var ticker []float64
	//var book []float64
	//var tickerAverage float64
	//var buyAverage float64
	//var sellAverage float64

	accounts, err := gdax.GetAccounts(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Accounts:", accounts)

	orders, err := gdax.ListOrders(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Orders:", orders)

	/* // get history
	product := "BTC-USD"
	product_table := "btc_usd"
	limit := int64(128)
	start := time.Now().Add(-time.Second * time.Duration(limit*analysis.TimeInterval)).Format(time.RFC3339)
	end := time.Now().Format(time.RFC3339)
	history, err := gdax.GetHistory(product, start, end, strconv.FormatInt(analysis.TimeInterval, 10))
	if err != nil {
		panic(err)
	}
	fmt.Println("History:", history)

	// archive history
	historian.ArchiveBtcUsd(db, history)

	// analyze history
	from := time.Now().Add(-time.Second * time.Duration(limit*analysis.TimeInterval)).Unix()
	to := time.Now().Unix()
	candles, err := historian.GetBtcUsd(db, analysis.TimeInterval, from, to)
	if err != nil {
		panic(err)
	}
	fmt.Println("MACD", analyst.Macd(analysis.EmaShort, analysis.EmaLong, candles))
	fmt.Println("RSI", analyst.Rsi(analysis.RsiPeriods, candles)) */

	for {
		rawJs := <-messages
		switch rawJs.(type) {
		case gdax.Ticker:
			fmt.Println("got ticker", rawJs)
		case gdax.Snapshot:
			fmt.Println("got a snapshot", rawJs)
		case gdax.Update:
			fmt.Println("got an update", rawJs)
		case string:
			fmt.Println("got a string ", rawJs)
		}

	}
}
