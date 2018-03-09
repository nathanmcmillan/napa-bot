package trader

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"../datastore"
	"../gdax"
	"../parse"
)

// Trade manages trading functions
type Trade struct {
	Datastore        *sql.DB
	Auth             *gdax.Authentication
	Settings         *gdax.Settings
	Macd             map[string]*Macd
	Ticker           map[string]*MovingAverage
	StoredOrders     map[string][]*datastore.Order
	OpenTransactions []*gdax.Order
}

// NewTrader constructor
func NewTrader(db *sql.DB, auth *gdax.Authentication, settings *gdax.Settings) *Trade {
	trader := &Trade{}
	trader.Datastore = db
	trader.Auth = auth
	trader.Settings = settings
	trader.Macd = make(map[string]*Macd)
	trader.Ticker = make(map[string]*MovingAverage)
	return trader
}

// Run core loop
func (trader *Trade) Run() {

	messages := make(chan interface{})
	go gdax.ExchangeSocket(trader.Settings, messages)
	go gdax.Polling(trader.Auth, trader.Settings, messages)

	retryWait := int64(1)
	retryLimit := int64(3)
	var tries int64
	var err error

	for i := 0; i < len(trader.Settings.Products); i++ {
		trader.Ticker[trader.Settings.Products[i]] = NewMovingAverage(10)
	}

	tries = 0
	for {
		trader.StoredOrders, err = datastore.QueryOrders(trader.Datastore)
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
	fmt.Println("stored orders :", trader.StoredOrders)

	for {
		rawMessage := <-messages
		switch message := rawMessage.(type) {
		case gdax.Ticker:
			move := trader.Ticker[message.ProductID]
			move.Rolling(message.Price)
			fmt.Println(message.ProductID, "| TICKER", move.Current, "|")
			trader.process(message.ProductID)
		case *gdax.CandleList:
			mac := NewMacd(trader.Settings.EmaShort, trader.Settings.EmaLong, message.List[0].Closing)
			for j := 1; j < len(message.List); j++ {
				mac.Update(message.List[j].Closing)
			}
			fmt.Println(message.Product, "| MACD", mac.Current, "| SIGNAL", mac.Signal, "|")
			trader.Macd[message.Product] = mac
			trader.process(message.Product)
		case error:
			fmt.Println("error", message)
		}
	}
}

func (trader *Trade) process(product string) {
	orders := trader.StoredOrders[product]
	macd := trader.Macd[product]
	ticker := trader.Ticker[product].Current
	if orders == nil || macd.Signal == "wait" {
		return
	}
	if macd.Signal == "sell" {
		for i := 0; i < len(orders); i++ {
			order := orders[i]
			if order.ProfitPrice >= ticker {
				trader.sell(product, order)
			}
		}
	} else if macd.Signal == "buy" {
		trader.buy(product)
	}
}

func (trader *Trade) sell(product string, order *datastore.Order) {
	var rawJs *bytes.Buffer
	parse.Begin(rawJs)
	parse.First(rawJs, "type", "market")
	parse.Append(rawJs, "side", "sell")
	parse.Append(rawJs, "product_id", product)
	parse.Append(rawJs, "size", strconv.FormatFloat(order.Size, 'f', -1, 64))
	parse.End(rawJs)
	transaction, err := gdax.PlaceOrder(trader.Auth, rawJs.String())
	if err != nil {
		fmt.Println("error placing order", rawJs, err)
		return
	}
	fmt.Println(transaction)
	trader.OpenTransactions = append(trader.OpenTransactions, transaction)
	go (func() {
		attempts := 0
		for {
			attempts++
			if attempts == 10 {
				fmt.Println("update order attempt limit")
				return
			}
			time.Sleep(time.Second)
			updated, err := gdax.GetOrder(trader.Auth, transaction.ID)
			if err != nil {
				fmt.Println("error getting order", err)
				continue
			}
			if updated.Settled {
				datastore.RemoveOrder(trader.Datastore, order.ID)
				fmt.Println(updated)
				return
			}
		}
	})()
}

func (trader *Trade) buy(product string) {
	accounts, err := gdax.GetAccounts(trader.Auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Accounts:", accounts)
	var usd *gdax.Account
	for i := 0; i < len(accounts); i++ {
		account := accounts[i]
		if account.Currency == "USD" {
			usd = &account
			break
		}
	}
	if usd == nil {
		fmt.Println("USD wallet not found")
		return
	}
	if usd.Available < 2000.0 {
		fmt.Println("not enough usd available $", usd.Available)
		return
	}
	var rawJs *bytes.Buffer
	parse.Begin(rawJs)
	parse.First(rawJs, "type", "market")
	parse.Append(rawJs, "side", "buy")
	parse.Append(rawJs, "product_id", product)
	parse.Append(rawJs, "funds", "1.0")
	parse.End(rawJs)
	transaction, err := gdax.PlaceOrder(trader.Auth, rawJs.String())
	if err != nil {
		fmt.Println("error placing order", rawJs, err)
		return
	}
	fmt.Println(transaction)
	trader.OpenTransactions = append(trader.OpenTransactions, transaction)
	go (func() {
		attempts := 0
		for {
			attempts++
			if attempts == 10 {
				fmt.Println("update order attempt limit")
				return
			}
			time.Sleep(time.Second)
			updated, err := gdax.GetOrder(trader.Auth, transaction.ID)
			if err != nil {
				fmt.Println("error getting order", err)
				continue
			}
			if updated.Settled {
				datastore.ArchiveOrder(trader.Datastore, updated.Product, updated.ExecutedValue, updated.Size)
				fmt.Println(updated)
				return
			}
		}
	})()
}
