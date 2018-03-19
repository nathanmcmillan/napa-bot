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
		"BTC-USD": 0.001,
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

func readOrder(auth map[string]string, orderID string) (*order, error) {
	body, _, err := privateRequest(auth, get, "/orders/"+orderID, "")
	if err != nil {
		return nil, err
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	return decodeOrder(decode), nil
}

func postOrder(auth map[string]string, rawJs string) (*order, error, int) {
	body, status, err := privateRequest(auth, post, "/orders", rawJs)
	if err != nil {
		return nil, err, status
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err, status
	}
	if status == 200 {
		return decodeOrder(decode), nil, status
	} else {
		return nil, errors.New(fmt.Sprint(decode)), status
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
