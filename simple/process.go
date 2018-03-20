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
		amount := "0.0" // amount := "10.0"
		if accounts["USD"].available.moreThan(newCurrency(amount)) {
			pending, status, err := buy(auth, product, amount)
			if err == nil && status == 200 {
				logger("buy", pending.id)
				orders.PushBack(pending)
				updates = true
			} else {
				if err == nil {
					err = errors.New("exchange response")
				}
				logger(err.Error(), "| status code", strconv.FormatInt(int64(status), 10))
			}
		}
	} else if algo.signal == "sell" {
		t, err := tick(product)
		if err != nil {
			logger(err.Error())
			return
		}
		var next *list.Element
		for e := orders.Front(); e != nil; e = next {
			next = e.Next()
			order := e.Value.(*order)
			min := profitPrice(order)
			fmt.Println("*", product, "|", min, ">", t.price, "*")
			if min.moreThan(t.price) {
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
		writeBytes(orderUpdateFile, []byte(buffer.String()))
		copyFile(orderFile, orderBackupFile)
		copyFile(orderUpdateFile, orderUpdateBackupFile)
		renameFile(orderUpdateFile, orderFile)
	}
}

func buy(a map[string]string, product, funds string) (*order, int, error) {
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

func sell(a map[string]string, o *order) (*order, int, error) {
	rawJs := &strings.Builder{}
	beginJs(rawJs)
	firstJs(rawJs, "type", "market")
	pushJs(rawJs, "side", "sell")
	pushJs(rawJs, "product_id", o.product)
	pushJs(rawJs, "size", o.size.str(precision[o.product]))
	endJs(rawJs)
	str := rawJs.String()
	logger(str)
	return postOrder(a, str)
}
