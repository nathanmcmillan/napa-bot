package main

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	percentRange = 0.05	
)

func process() {
	if algo.signal == "wait" {
		return
	}
	updates := false
	ticker, err := tick(product)
	if err != nil {
		logger(err.Error())
		return
	}
	if algo.signal == "buy" {
		tickerFloat := ticker.price.float()
		for e := orders.Front(); e != nil; e = e.Next() {
			existingOrder := e.Value.(*order)
			price := priceOfCoin(existingOrder).float()
			if percentChange(tickerFloat, price) < percentRange {
				fmt.Println("* existing order", existingOrder.id, "bought at $", price, "withing", percentRange, "of last ticker $", tickerFloat, "*")
				return
			}
		}
		accounts, err := readAccounts(auth)
		if err != nil {
			logger(err.Error())
			return
		}
		fund := funds[product]
		availableUsd := accounts["USD"].available
		if fund.moreThan(availableUsd) && fund.moreThan(twenty) {
			amount := fund.div(two)
			pending, status, err := buy(auth, product, amount.str(2))
			if err == nil && status == 200 {
				logger("buy", pending.id)
				settledOrder := waitTilSettled(pending.id)
				funds[product] = funds[product].minus(settledOrder.executedValue).minus(settledOrder.fillFees)
				err = updateFundsFile()
				logger(amount.str(2), "->", settledOrder.id, "| funds $", funds[product].str(2))
				if err != nil {
					logger(err.Error())
					panic(err)
				}
				orders.PushBack(settledOrder)
				updates = true
			} else {
				if err == nil {
					err = errors.New("exchange response")
				}
				logger(err.Error(), "| status code", strconv.FormatInt(int64(status), 10))
			}
		}
	} else if algo.signal == "sell" {
		var next *list.Element
		for e := orders.Front(); e != nil; e = next {
			next = e.Next()
			orderToSell := e.Value.(*order)
			min, err := profitPrice(orderToSell)
			if err != nil {
				logger(err.Error())
				continue
			}
			fmt.Println("*", product, "|", ticker.price.str(10), ">", min.str(10), "? *")
			if ticker.price.moreThan(min) {
				pending, status, err := sell(auth, orderToSell)
				if err == nil && status == 200 {
					logger("sell", pending.id)
					settledOrder := waitTilSettled(pending.id)
					profits := settledOrder.executedValue.minus(orderToSell.executedValue).minus(settledOrder.fillFees)
					funds[product] = funds[product].plus(profits.mul(percent85))
					err = updateFundsFile()
					logger(orderToSell.id, "->", settledOrder.id, "| profit $", profits.str(2), "| funds $", funds[product].str(2))
					if err != nil {
						logger(err.Error())
						panic(err)
					}
					orders.Remove(e)
					updates = true
				} else {
					if err == nil {
						err = errors.New("exchange response")
					}
					logger(err.Error(), "| status code", strconv.FormatInt(int64(status), 10))
				}
			}
		}
	}
	if updates {
		var buffer strings.Builder
		for e := orders.Front(); e != nil; e = e.Next() {
			order := e.Value.(*order)
			buffer.WriteString(order.id)
			buffer.WriteString("\n")
		}
		err := updateOrderFile(buffer.String())
		if err != nil {
			logger(err.Error())
			panic(err)
		}
	}
}

func updateOrderFile(contents string) error {
	err := writeBytes(orderUpdateFile, []byte(contents))
	if err != nil {
		return err
	}
	err = copyFile(orderFile, orderBackupFile)
	if err != nil {
		return err
	}
	err = copyFile(orderUpdateFile, orderUpdateBackupFile)
	if err != nil {
		return err
	}
	return renameFile(orderUpdateFile, orderFile)
}

func updateFundsFile() error {
	var buffer strings.Builder
	for key, value := range funds {
		buffer.WriteString(key)
		buffer.WriteString(" ")
		buffer.WriteString(value.num.String())
		buffer.WriteString("\n")
	}
	err := writeBytes(fundsUpdateFile, []byte(buffer.String()))
	if err != nil {
		return err
	}
	err = copyFile(fundsFile, fundsBackupFile)
	if err != nil {
		return err
	}
	err = copyFile(fundsUpdateFile, fundsUpdateBackupFile)
	if err != nil {
		return err
	}
	return renameFile(fundsUpdateFile, fundsFile)
}

func waitTilSettled(orderID string) (*order) {
	for {
		time.Sleep(time.Second)
		orderUpdate, status, err := readOrder(auth, orderID)
		if orderUpdate == nil || err != nil {
			logger("could not update order " + orderID + " | status code " + fmt.Sprint(status))
			continue
		}
		if orderUpdate.settled {
			return orderUpdate	
		}
	}
}

func buy(a map[string]string, product, funds string) (*orderResponse, int, error) {
	rawJs := &strings.Builder{}
	beginJs(rawJs)
	firstJs(rawJs, "type", "market")
	pushJs(rawJs, "side", "buy")
	pushJs(rawJs, "product_id", product)
	pushJs(rawJs, "funds", funds)
	endJs(rawJs)
	str := rawJs.String()
	logger(str)
	return postOrder(a, str)
}

func sell(a map[string]string, o *order) (*orderResponse, int, error) {
	rawJs := &strings.Builder{}
	beginJs(rawJs)
	firstJs(rawJs, "type", "market")
	pushJs(rawJs, "side", "sell")
	pushJs(rawJs, "product_id", o.product)
	pushJs(rawJs, "size", o.filledSizeS)
	endJs(rawJs)
	str := rawJs.String()
	logger(str)
	return postOrder(a, str)
}
