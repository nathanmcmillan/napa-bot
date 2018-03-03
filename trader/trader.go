package trader

import (
	"database/sql"
	"fmt"

	"../analyst"
	"../datastore"
	"../gdax"
)

// Run core loop
func Run(db *sql.DB, auth *gdax.Authentication, products []string, settings *analyst.Settings, messages chan interface{}) {

	//var rsi []float64
	//var macd []float64
	//var book []float64
	//var tickerAverage float64
	//var buyAverage float64
	//var sellAverage float64

	books := make(map[string]*Book)
	ticker := make(map[string]*MovingAverage)
	for i := 0; i < len(products); i++ {
		books[products[i]] = NewBook()
		ticker[products[i]] = NewMovingAverage(10)
	}

	accounts, err := gdax.GetAccounts(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Accounts:", accounts)

	openOrders, err := gdax.ListOrders(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Open Orders:", openOrders)

	orders, err := datastore.ListOrders(db)
	if err != nil {
		panic(err)
	}

	for {
		rawMessage := <-messages
		switch message := rawMessage.(type) {
		case gdax.Ticker:
			move := ticker[message.ProductID]
			move.Rolling(message.Price)
			review(orders[message.ProductID], move.Current)
		case gdax.Snapshot:
			book := books[message.ProductID]
			book.Snapshot(&message)
			fmt.Println("book init", book)
		case gdax.Update:
			book := books[message.ProductID]
			book.Update(&message)
			fmt.Println("book init", book)
		case string:
			fmt.Println("got a string ", message)
		}
	}
}

func review(orders []*gdax.Order, ticker float64) {
	fmt.Println("ticker", ticker)
	if orders == nil {
		return
	}
	for i := 0; i < len(orders); i++ {
		fmt.Println("orders")
	}
}

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
