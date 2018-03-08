package trader

import (
	"database/sql"
	"fmt"
	"time"

	"../datastore"
	"../gdax"
)

// Run core loop
func Run(db *sql.DB, auth *gdax.Authentication, settings *gdax.Settings) {

	messages := make(chan interface{})
	go gdax.ExchangeSocket(settings, messages)
	go gdax.Polling(auth, settings, messages)

	retryWait := int64(1)
	retryLimit := int64(3)
	var tries int64
	var err error

	macd := make(map[string]*Macd)
	ticker := make(map[string]*MovingAverage)

	for i := 0; i < len(settings.Products); i++ {
		macd[settings.Products[i]] = NewMacd(settings.EmaShort, settings.EmaLong)
		ticker[settings.Products[i]] = NewMovingAverage(10)
	}

	/* accounts, err := gdax.GetAccounts(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Accounts:", accounts)

	openOrders, err := gdax.ListOrders(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Open Orders:", openOrders) */

	var orders map[string][]*datastore.Order
	tries = 0
	for {
		orders, err = datastore.QueryOrders(db)
		if err == nil {
			break
		}
		if tries >= retryLimit {
			panic(err)
		}
		tries++
		fmt.Println("query order", tries, "/", retryLimit, "failed", err)
		time.Sleep(time.Second * time.Duration(retryWait))
	}
	fmt.Println("orders :", orders)

	for {
		rawMessage := <-messages
		switch message := rawMessage.(type) {
		case gdax.Ticker:
			move := ticker[message.ProductID]
			move.Rolling(message.Price)
			fmt.Println(message.ProductID, "| TICKER", move.Current, "|")
			tickerReview(orders[message.ProductID], move.Current)
		case *gdax.CandleList:
			mac := NewMacd(settings.EmaShort, settings.EmaLong)
			for j := 0; j < len(message.List); j++ {
				// fmt.Println("TIME", message.List[j].Time, "/", time.Unix(message.List[j].Time, 0), "CLOSING ", message.List[j].Closing)
				mac.Update(message.List[j].Closing)
			}
			fmt.Println(message.Product, "| MACD", mac.Current, "| SIGNAL", mac.Signal, "|")
			macd[message.Product] = mac
			macdReview(orders[message.Product], mac, ticker[message.Product].Current)
		case error:
			fmt.Println("error", message)
		}
	}
}

func tickerReview(orders []*datastore.Order, ticker float64) {
	if orders == nil {
		return
	}
	for i := 0; i < len(orders); i++ {
		order := orders[i]
		if order.ProfitPrice >= ticker {
			continueReview(order)
		}
	}
}

func continueReview(order *datastore.Order) {

}

func macdReview(orders []*datastore.Order, mac *Macd, ticker float64) {
	if orders == nil || mac.Signal == "wait" {
		return
	}
	if mac.Signal == "sell" {
		for i := 0; i < len(orders); i++ {
			order := orders[i]
			if order.ProfitPrice >= ticker {
				continueReview(order)
			}
		}
	} else if mac.Signal == "buy" {
		fmt.Println("buy code")
	}
}
