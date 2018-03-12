package trader

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"

	"../datastore"
	"../gdax"
	"../parse"
)

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
	if trader.Accounts[product].Funds > usd.Available {
		msg := fmt.Sprintf("account funds greater than available USD %s > %s", trader.Accounts[product].Funds, usd.Available)
		log.Println(msg)
		fmt.Println(msg)
		return
	}
	amount := float64(5.0)
	exchangeOrder, err := trader.PlaceMarketBuy(product, amount)
	if err != nil {
		log.Println(err)
		fmt.Println("error placing order", err)
		return
	}
	fmt.Println(exchangeOrder)
	datastore.ArchiveOrder(trader.Datastore, exchangeOrder.ID)
	attempts := 0
	for {
		attempts++
		if attempts == 10 {
			fmt.Println("update order attempt limit")
			break
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
			msg := fmt.Sprintf("buy %s", exchangeOrderUpdate)
			log.Println(msg)
			fmt.Println(msg)
			break
		}
	}
}

func (trader *Trade) sell(product string, order *gdax.Order) {
	exchangeOrder, err := trader.PlaceMarketSell(product, order.Size)
	if err != nil {
		log.Println(err)
		fmt.Println("error placing order", err)
		return
	}
	fmt.Println(exchangeOrder)
	attempts := 0
	for {
		attempts++
		if attempts == 10 {
			fmt.Println("update order attempt limit")
			break
		}
		time.Sleep(time.Second)
		exchangeOrderUpdate, err := gdax.GetOrder(trader.Auth, exchangeOrder.ID)
		if err != nil {
			log.Println(err)
			fmt.Println("error getting order", err)
			continue
		}
		if exchangeOrderUpdate.Settled {
			fmt.Println("original buy", order, "sold as", exchangeOrderUpdate)
			datastore.RemoveOrder(trader.Datastore, order.ID)
			found := false
			for i := 0; i < len(trader.Orders[product]); i++ {
				if trader.Orders[product][i].ID == order.ID {
					found = true
					trader.Orders[product] = append(trader.Orders[product][:i], trader.Orders[product][i+1:]...)
					break
				}
			}
			if !found {
				msg := "original buy order not found in order list"
				log.Println(msg)
				fmt.Println(msg)	
			}
			break
		}
	}
}

func (trader *Trade) PlaceMarketBuy(product string, usd float64) (*gdax.Order, error) {
	var rawJs *bytes.Buffer
	parse.Begin(rawJs)
	parse.First(rawJs, "type", "market")
	parse.Append(rawJs, "side", "buy")
	parse.Append(rawJs, "product_id", product)
	parse.Append(rawJs, "funds", strconv.FormatFloat(usd, 'f', -1, 64))
	parse.End(rawJs)
	return gdax.PlaceOrder(trader.Auth, rawJs.String())
}

func (trader *Trade) PlaceMarketSell(product string, size float64) (*gdax.Order, error) {
	var rawJs *bytes.Buffer
	parse.Begin(rawJs)
	parse.First(rawJs, "type", "market")
	parse.Append(rawJs, "side", "sell")
	parse.Append(rawJs, "product_id", product)
	parse.Append(rawJs, "size", strconv.FormatFloat(size, 'f', -1, 64))
	parse.End(rawJs)
	return gdax.PlaceOrder(trader.Auth, rawJs.String())
}