package main

import (
	"encoding/json"
)

type ticker struct {
	tradeID int64
	price   *currency
	size    *currency
	bid     *currency
	ask     *currency
	volume  *currency
	time    string
}

func tick(product string) (*ticker, error) {
	body, _, err := publicRequest(get, "/products/"+product+"/ticker")
	if err != nil {
		return nil, err
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	t := &ticker{}
	t.tradeID, _ = decode["trade_id"].(int64)
	temp, _ := decode["price"].(string)
	t.price = newCurrency(temp)
	temp, _ = decode["size"].(string)
	t.size = newCurrency(temp)
	temp, _ = decode["bid"].(string)
	t.bid = newCurrency(temp)
	temp, _ = decode["ask"].(string)
	t.ask = newCurrency(temp)
	temp, _ = decode["volume"].(string)
	t.volume = newCurrency(temp)
	t.time, _ = decode["time"].(string)
	return t, nil
}
