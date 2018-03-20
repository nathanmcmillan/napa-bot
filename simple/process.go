package main

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func process() {
	updates := false
	if algo.signal == "buy" {
		accounts, err := readAccounts(auth)
		if err != nil {
			logger(err.Error())
			return
		}
		amount := newCurrency("0.0") // amount := "10.0"
		if accounts["USD"].available.moreThan(amount) {
			pending, status, err := buy(auth, product, amount.str(2))
			if err == nil && status == 200 {
				logger("buy", pending.id)
				hollow := &order{}
				hollow.id = pending.id
				hollow.settled = false
				orders.PushBack(hollow)
				updates = true
			} else {
				if err == nil {
					err = errors.New("exchange response")
				}
				logger(err.Error(), "| status code", strconv.FormatInt(int64(status), 10))
			}
		}
	} else if algo.signal == "sell" {
		ticker, err := tick(product)
		if err != nil {
			logger(err.Error())
			return
		}
		var next *list.Element
		for e := orders.Front(); e != nil; e = next {
			next = e.Next()
			order := e.Value.(*order)
			if !order.settled {
				orderUpdate, status, err := readOrder(auth, order.id)
				if orderUpdate == nil || err != nil {
					logger("could not update order | status code " + fmt.Sprint(status) + " | " + fmt.Sprint(order))
					continue
				}
				order = orderUpdate
				e.Value = orderUpdate
			}
			min, err := profitPrice(order)
			if err != nil {
				logger(err.Error())
				continue
			}
			fmt.Println("*", product, "|", ticker.price.str(10), ">", min.str(10), "? *")
			if ticker.price.moreThan(min) {
				pending, status, err := sell(auth, order)
				if err == nil && status == 200 {
					fmt.Println("sell", order, pending)
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
