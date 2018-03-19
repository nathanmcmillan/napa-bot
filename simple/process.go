package main

import (
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
		amount := "10.0"
		if accounts["USD"].available.moreThan(newCurrency(amount)) {
			pending, status, err := buy(auth, product, amount)
			if err == nil && status == 200 {
				fmt.Println(pending.id)
				orders.push(pending)
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
		for i := 0; i < len(orders); i++ {
			order := orders[i]
			min := profitPrice(order)
			fmt.Println("*", product, "|", min, ">", t.price, "*")
			if min.moreThan(t.price) {
				pending, status, err := sell(auth, order)
				if err == nil && status == 200 {
					fmt.Println(pending.id)
					orders.delete(i)
					i--
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
		for i := 0; i < len(orders); i++ {
			buffer.WriteString(orders[i].id)
			buffer.WriteByte('\n')
		}
		writeList(orderSwapFile, []byte(buffer.String()))
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
