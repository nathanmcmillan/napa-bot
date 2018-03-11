package trader

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"../datastore"
	"../gdax"
	"../parse"
)

// Trade manages trading functions
type Trade struct {
	Datastore  *sql.DB
	Auth       *gdax.Authentication
	Settings   *gdax.Settings
	Signals    chan os.Signal
	Accounts   map[string]*datastore.Account
	Macd       map[string]*Macd
	Ticker     map[string]*MovingAverage
	Orders     map[string][]*gdax.Order
	OpenOrders map[string][]*gdax.Order
}

// NewTrader constructor
func NewTrader(db *sql.DB, auth *gdax.Authentication, settings *gdax.Settings, signals chan os.Signal) *Trade {
	trader := &Trade{}
	trader.Datastore = db
	trader.Auth = auth
	trader.Settings = settings
	trader.Signals = signals
	trader.Accounts = make(map[string]*datastore.Account)
	trader.Macd = make(map[string]*Macd)
	trader.Ticker = make(map[string]*MovingAverage)
	trader.Orders = make(map[string][]*gdax.Order, 0)
	trader.OpenOrders = make(map[string][]*gdax.Order, 0)
	return trader
}

// Run core loop
func (trader *Trade) Run() {

	messages := make(chan interface{})
	socketDone := make(chan bool)
	pollingDone := make(chan bool)
	go gdax.ExchangeSocket(trader.Settings, messages, socketDone)
	go gdax.Polling(trader.Auth, trader.Settings, messages, pollingDone)

	retryWait := int64(1)
	retryLimit := int64(3)
	var tries int64
	var err error

	for i := 0; i < len(trader.Settings.Products); i++ {
		product := trader.Settings.Products[i]
		trader.Ticker[product] = NewMovingAverage(10)
		trader.Orders[product] = make([]*gdax.Order, 0)
		trader.OpenOrders[product] = make([]*gdax.Order, 0)
	}

	var storedAccounts []*datastore.Account
	tries = 0
	for {
		storedAccounts, err = datastore.QueryAccounts(trader.Datastore)
		if err == nil {
			break
		}
		if tries >= retryLimit {
			log.Println(err)
			panic(err)
		}
		tries++
		fmt.Println("query accounts", tries, "/", retryLimit, "failed", err)
		time.Sleep(time.Second * time.Duration(retryWait))
	}
	fmt.Println("stored accounts :", storedAccounts)
	for i := 0; i < len(storedAccounts); i++ {
		account := storedAccounts[i]
		trader.Accounts[account.Product] = account
	}

	var storedOrders []*datastore.Order
	tries = 0
	for {
		storedOrders, err = datastore.QueryOrders(trader.Datastore)
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
	fmt.Println("stored orders :", storedOrders)
	for i := 0; i < len(storedOrders); i++ {
		order, err := gdax.GetOrder(trader.Auth, storedOrders[i].ExchangeID)
		if err != nil {
			log.Println(err)
			fmt.Println("error getting order", err)
			continue
		}
		trader.Orders[order.Product] = append(trader.Orders[order.Product], order)
	}
loop:
	for {
		select {
		case rawMessage := <-messages:
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
		case signalMessage := <-trader.Signals:
			fmt.Println(" signal", signalMessage, "closing trader")
			socketDone <- true
			pollingDone <- true
			break loop
		}
	}
}

func (trader *Trade) process(product string) {
	orders := trader.Orders[product]
	macd := trader.Macd[product]
	ticker := trader.Ticker[product].Current
	if orders == nil || macd == nil || macd.Signal == "wait" {
		return
	}
	if macd.Signal == "sell" {
		for i := 0; i < len(orders); i++ {
			order := orders[i]
			if gdax.ProfitPrice(product, order.ExecutedValue) >= ticker {
				trader.sell(product, order)
			}
		}
	} else if macd.Signal == "buy" {
		trader.buy(product)
	}
}

func (trader *Trade) sell(product string, order *gdax.Order) {
	var rawJs *bytes.Buffer
	parse.Begin(rawJs)
	parse.First(rawJs, "type", "market")
	parse.Append(rawJs, "side", "sell")
	parse.Append(rawJs, "product_id", product)
	parse.Append(rawJs, "size", strconv.FormatFloat(order.Size, 'f', -1, 64))
	parse.End(rawJs)
	exchangeOrder, err := gdax.PlaceOrder(trader.Auth, rawJs.String())
	if err != nil {
		log.Println(err)
		fmt.Println("error placing order", rawJs, err)
		return
	}
	fmt.Println(exchangeOrder)
	trader.OpenOrders[product] = append(trader.OpenOrders[product], exchangeOrder)
	go (func() {
		attempts := 0
		for {
			attempts++
			if attempts == 10 {
				fmt.Println("update order attempt limit")
				return
			}
			time.Sleep(time.Second)
			exchangeOrderUpdate, err := gdax.GetOrder(trader.Auth, exchangeOrder.ID)
			if err != nil {
				log.Println(err)
				fmt.Println("error getting order", err)
				continue
			}
			if exchangeOrderUpdate.Settled {
				datastore.RemoveOrder(trader.Datastore, order.ID)
				fmt.Println("original buy", order, "sold as", exchangeOrderUpdate)
				return
			}
		}
	})()
}

func (trader *Trade) buy(product string) {
	accounts, err := gdax.GetAccounts(trader.Auth)
	if err != nil {
		log.Println(err)
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
	if trader.Accounts[product].Funds < 2000.0 || usd.Available < 2000.0 {
		fmt.Println("not enough usd available $", usd.Available, "funds in account", product, "$", trader.Accounts[product].Funds)
		return
	}
	var rawJs *bytes.Buffer
	parse.Begin(rawJs)
	parse.First(rawJs, "type", "market")
	parse.Append(rawJs, "side", "buy")
	parse.Append(rawJs, "product_id", product)
	parse.Append(rawJs, "funds", "0.0")
	parse.End(rawJs)
	exchangeOrder, err := gdax.PlaceOrder(trader.Auth, rawJs.String())
	if err != nil {
		log.Println(err)
		fmt.Println("error placing order", rawJs, err)
		return
	}
	fmt.Println(exchangeOrder)
	datastore.ArchiveOrder(trader.Datastore, exchangeOrder.ID)
	trader.OpenOrders[product] = append(trader.OpenOrders[product], exchangeOrder)
	go (func() {
		attempts := 0
		for {
			attempts++
			if attempts == 10 {
				fmt.Println("update order attempt limit")
				return
			}
			time.Sleep(time.Second)
			exchangeOrderUpdate, err := gdax.GetOrder(trader.Auth, exchangeOrder.ID)
			if err != nil {
				log.Println(err)
				fmt.Println("error getting order", err)
				continue
			}
			if exchangeOrderUpdate.Settled {
				datastore.ArchiveOrder(trader.Datastore, exchangeOrder.ID)
				fmt.Println("buy", exchangeOrderUpdate)
				return
			}
		}
	})()
}
