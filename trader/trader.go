package trader

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"../datastore"
	"../gdax"
	"../analysis"
)

// Trade manages trading functions
type Trade struct {
	Datastore  *sql.DB
	Auth       *gdax.Authentication
	Settings   *gdax.Settings
	Signals    chan os.Signal
	Accounts   map[string]*datastore.Account
	Macd       map[string]*analysis.Macd
	Ticker     map[string]*analysis.MovingAverage
	Orders     map[string][]*gdax.Order
}

// NewTrader constructor
func NewTrader(db *sql.DB, auth *gdax.Authentication, settings *gdax.Settings, signals chan os.Signal) *Trade {
	trader := &Trade{}
	trader.Datastore = db
	trader.Auth = auth
	trader.Settings = settings
	trader.Signals = signals
	trader.Accounts = make(map[string]*datastore.Account)
	trader.Macd = make(map[string]*analysis.Macd)
	trader.Ticker = make(map[string]*analysis.MovingAverage)
	trader.Orders = make(map[string][]*gdax.Order, 0)
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
		trader.Ticker[product] = analysis.NewMovingAverage(10)
		trader.Orders[product] = make([]*gdax.Order, 0)
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
				mac := analysis.NewMacd(trader.Settings.EmaShort, trader.Settings.EmaLong, message.List[0].Closing)
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
