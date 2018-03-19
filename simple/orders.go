package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	taker = map[string]*currency{
		"BTC-USD": newCurrency("1.0025"),
		"ETH-USD": newCurrency("1.003"),
		"LTC-USD": newCurrency("1.003"),
	}
	minimum = map[string]float64{
		"BTC-USD": 10.0,
	}
	precision = map[string]int{
		"BTC-USD": 3,
	}
)

type book []*order

func (b *book) push(o *order) {
	size := len(*b)
	for i := 0; i < size; i++ {
		if (*b)[i] == nil {
			(*b)[i] = o
			return
		}
	}
	*b = append(*b, o)
}

func (b *book) delete(i int) {
	size := len(*b)
	(*b)[i] = (*b)[size-1]
	*b = (*b)[:size-1]
}

type order struct {
	id            string
	price         *currency
	size          *currency
	product       string
	side          string
	stp           string
	typ           string
	timeInForce   string
	postOnly      bool
	createdAt     string
	fillFees      *currency
	filledSize    *currency
	executedValue *currency
	status        string
	settled       bool
}

func profitPrice(o *order) *currency {
	return o.executedValue.mul(taker[o.product])
}

func readOrder(auth map[string]string, orderID string) (*order, int, error) {
	body, status, err := privateRequest(auth, get, "/orders/"+orderID, "")
	return processOrder(body, status, err)
}

func postOrder(auth map[string]string, rawJs string) (*order, int, error) {
	body, status, err := privateRequest(auth, post, "/orders", rawJs)
	return processOrder(body, status, err)
}

func processOrder(body []byte, status int, err error) (*order, int, error) {
	if err != nil {
		return nil, 0, err
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, 0, err
	}
	if status == 200 {
		return decodeOrder(decode), status, nil
	} else {
		return nil, status, errors.New(fmt.Sprint(decode))
	}
}

func decodeOrder(d map[string]interface{}) *order {
	o := &order{}
	o.id, _ = d["id"].(string)
	temp, _ := d["price"].(string)
	o.price = newCurrency(temp)
	temp, _ = d["size"].(string)
	o.size = newCurrency(temp)
	o.product, _ = d["product_id"].(string)
	o.side, _ = d["side"].(string)
	o.stp, _ = d["stp"].(string)
	o.typ, _ = d["type"].(string)
	o.timeInForce, _ = d["time_in_force"].(string)
	o.postOnly, _ = d["post_only"].(bool)
	o.createdAt, _ = d["created_at"].(string)
	temp, _ = d["fill_fees"].(string)
	o.fillFees = newCurrency(temp)
	temp, _ = d["fillled_size"].(string)
	o.filledSize = newCurrency(temp)
	temp, _ = d["executed_value"].(string)
	o.executedValue = newCurrency(temp)
	o.status, _ = d["status"].(string)
	o.settled, _ = d["settled"].(bool)
	return o
}
