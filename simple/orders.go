package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type orderResponse struct {
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
	filledSizeS   string
	executedValue *currency
	status        string
	settled       bool
}

type order struct {
	id            string
	size          *currency
	product       string
	side          string
	stp           string
	funds *currency
	specifiedFunds *currency
	typ           string
	postOnly      bool
	createdAt     string
	doneAt	 	  string
	doneReason    string
	fillFees      *currency
	filledSize    *currency
	filledSizeS   string
	executedValue *currency
	status        string
	settled       bool
}

func profitPrice(o *order) (*currency, error) {
	// executed value / filled_size = price of coin
	// specified_funds / funds = required margin
	funds := o.funds
	specifiedFunds := o.specifiedFunds
	executedValue := o.executedValue
	filledSize := o.filledSize
	margin := specifiedFunds.div(funds).minus(one).mul(two).plus(one)
	return executedValue.div(filledSize).mul(margin), nil
}

func readOrder(auth map[string]string, orderID string) (*order, int, error) {
	body, status, err := privateRequest(auth, get, "/orders/"+orderID, "")
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
	}
	return nil, status, errors.New(fmt.Sprint(decode))
}

func postOrder(auth map[string]string, rawJs string) (*orderResponse, int, error) {
	body, status, err := privateRequest(auth, post, "/orders", rawJs)
	if err != nil {
		return nil, 0, err
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, 0, err
	}
	if status == 200 {
		return decodeOrderResponse(decode), status, nil
	}
	return nil, status, errors.New(fmt.Sprint(decode))
}

func decodeOrderResponse(d map[string]interface{}) *orderResponse {
	o := &orderResponse{}
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
	o.filledSizeS, _ = d["filled_size"].(string)
	o.filledSize = newCurrency(o.filledSizeS)
	temp, _ = d["executed_value"].(string)
	o.executedValue = newCurrency(temp)
	o.status, _ = d["status"].(string)
	o.settled, _ = d["settled"].(bool)
	return o
}

func decodeOrder(d map[string]interface{}) *order {
	o := &order{}
	o.id, _ = d["id"].(string)
	temp, _ := d["size"].(string)
	o.size = newCurrency(temp)
	o.product, _ = d["product_id"].(string)
	o.side, _ = d["side"].(string)
	o.stp, _ = d["stp"].(string)
	temp, _ = d["funds"].(string)
	o.funds = newCurrency(temp)
	temp, _ = d["specified_funds"].(string)
	o.specifiedFunds = newCurrency(temp)
	o.typ, _ = d["type"].(string)
	o.postOnly, _ = d["post_only"].(bool)
	o.createdAt, _ = d["created_at"].(string)
	o.doneAt, _ = d["done_at"].(string)
	o.doneReason, _ = d["done_reason"].(string)
	temp, _ = d["fill_fees"].(string)
	o.fillFees = newCurrency(temp)
	o.filledSizeS, _ = d["filled_size"].(string)
	o.filledSize = newCurrency(o.filledSizeS)
	temp, _ = d["executed_value"].(string)
	o.executedValue = newCurrency(temp)
	o.status, _ = d["status"].(string)
	o.settled, _ = d["settled"].(bool)
	return o
}
