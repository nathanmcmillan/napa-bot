package gdax

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
)

// GetCurrencies list of currencies
func (r *Rest) GetCurrencies() (interface{}, error) {
	body, err := r.publicRequest(get, api+"/currencies")
	if err != nil {
		return nil, err
	}
	var decode interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	return decode, nil
}

// GetBook map of level 2 product books
func (r *Rest) GetBook(product string) (interface{}, error) {
	body, err := r.publicRequest(get, api+"/products/"+product+"/book?level=2")
	if err != nil {
		return nil, err
	}
	var decode interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	return decode, nil
}

// GetTicker map of product ticker
func (r *Rest) GetTicker(product string) (interface{}, error) {
	body, err := r.publicRequest(get, api+"/products/"+product+"/ticker")
	if err != nil {
		return nil, err
	}
	var decode interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	return decode, nil
}

// GetTrades map of product trades
func (r *Rest) GetTrades(product string) (interface{}, error) {
	body, err := r.publicRequest(get, api+"/products/"+product+"/trades")
	if err != nil {
		return nil, err
	}
	var decode interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	return decode, nil
}

// GetHistory list of candles for product history
func (r *Rest) GetHistory(product, start, end, granularity string) (*CandleList, error) {
	body, err := r.publicRequest(get, api+"/products/"+product+"/candles?start="+start+"&end="+end+"&granularity="+granularity)
	if err != nil {
		return nil, err
	}
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	candles := make([]*Candle, 0)
	for i := 0; i < len(decode); i++ {
		values, ok := decode[i].([]interface{})
		if !ok {
			return nil, errors.New("not a list")
		}
		candle := &Candle{}
		floatTime, _ := values[0].(float64)
		candle.Time = int64(floatTime)
		candle.Low, _ = values[1].(float64)
		candle.High, _ = values[2].(float64)
		candle.Open, _ = values[3].(float64)
		candle.Closing, _ = values[4].(float64)
		candle.Volume, _ = values[5].(float64)
		candles = append(candles, candle)
	}
	sort.Sort(SortCandlesByTime(candles))
	return &CandleList{product, candles}, nil
}

// GetStats map of 24 hour product statistics
func (r *Rest) GetStats(product string) (interface{}, error) {
	body, err := r.publicRequest(get, api+"/products/"+product+"/stats")
	if err != nil {
		return nil, err
	}
	var decode interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	return decode, nil
}

// GetAccounts map of account balances
func (r *Rest) GetAccounts() ([]Account, error) {
	body, err := r.privateRequest(get, api, "/accounts", "")
	if err != nil {
		return nil, err
	}
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil, err
	}
	accounts := make([]Account, 0)
	for i := 0; i < len(decode); i++ {
		values, ok := decode[i].(map[string]interface{})
		if !ok {
			return nil, errors.New("parse error")
		}
		account := Account{}
		account.ID, _ = values["id"].(string)
		account.Currency, _ = values["currency"].(string)
		account.Balance, _ = values["balance"].(float64)
		account.Available, _ = values["available"].(float64)
		account.Hold, _ = values["hold"].(float64)
		account.ProfileID, _ = values["profile_id"].(string)
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// PlaceOrder send a buy or sell order
func (r *Rest) PlaceOrder(rawJs string) (*Order, error) {
	body, err := r.privateRequest(post, api, "/orders", rawJs)
	if err != nil {
		return nil, err
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	order := &Order{}
	order.ID, _ = decode["id"].(string)
	temp, _ := decode["price"].(string)
	order.Price, _ = strconv.ParseFloat(temp, 64)
	temp, _ = decode["size"].(string)
	order.Size, _ = strconv.ParseFloat(temp, 64)
	order.Product, _ = decode["product_id"].(string)
	order.Side, _ = decode["side"].(string)
	order.Stp, _ = decode["stp"].(string)
	order.Type, _ = decode["type"].(string)
	order.TimeInForce, _ = decode["time_in_force"].(string)
	order.PostOnly, _ = decode["post_only"].(bool)
	order.CreatedAt, _ = decode["created_at"].(string)
	temp, _ = decode["fill_fees"].(string)
	order.FillFees, _ = strconv.ParseFloat(temp, 64)
	temp, _ = decode["fillled_size"].(string)
	order.FilledSize, _ = strconv.ParseFloat(temp, 64)
	temp, _ = decode["executed_value"].(string)
	order.ExecutedValue, _ = strconv.ParseFloat(temp, 64)
	order.Status, _ = decode["status"].(string)
	order.Settled, _ = decode["settled"].(bool)
	return order, nil
}

// ListOrders get open orders
func (r *Rest) ListOrders() (map[string][]*Order, error) {
	body, err := r.privateRequest(get, api, "/orders", "")
	if err != nil {
		return nil, err
	}
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil, err
	}
	orders := make(map[string][]*Order)
	for i := 0; i < len(decode); i++ {
		values, ok := decode[i].(map[string]interface{})
		if !ok {
			return nil, errors.New("parse error")
		}
		order := &Order{}
		order.ID, _ = values["id"].(string)
		temp, _ := values["price"].(string)
		order.Price, _ = strconv.ParseFloat(temp, 64)
		temp, _ = values["size"].(string)
		order.Size, _ = strconv.ParseFloat(temp, 64)
		order.Product, _ = values["product_id"].(string)
		order.Side, _ = values["side"].(string)
		order.Stp, _ = values["stp"].(string)
		order.Type, _ = values["type"].(string)
		order.TimeInForce, _ = values["time_in_force"].(string)
		order.PostOnly, _ = values["post_only"].(bool)
		order.CreatedAt, _ = values["created_at"].(string)
		temp, _ = values["fill_fees"].(string)
		order.FillFees, _ = strconv.ParseFloat(temp, 64)
		temp, _ = values["fillled_size"].(string)
		order.FilledSize, _ = strconv.ParseFloat(temp, 64)
		temp, _ = values["executed_value"].(string)
		order.ExecutedValue, _ = strconv.ParseFloat(temp, 64)
		order.Status, _ = values["status"].(string)
		order.Settled, _ = values["settled"].(bool)

		if orders[order.Product] == nil {
			orders[order.Product] = make([]*Order, 0)
		}
		orders[order.Product] = append(orders[order.Product], order)
	}
	return orders, nil
}

// GetOrder get an order by id
func (r *Rest) GetOrder(orderID string) (*Order, error) {
	body, err := r.privateRequest(get, api, "/orders/"+orderID, "")
	if err != nil {
		return nil, err
	}
	var decode map[string]interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	order := &Order{}
	order.ID, _ = decode["id"].(string)
	temp, _ := decode["price"].(string)
	order.Price, _ = strconv.ParseFloat(temp, 64)
	temp, _ = decode["size"].(string)
	order.Size, _ = strconv.ParseFloat(temp, 64)
	order.Product, _ = decode["product_id"].(string)
	order.Side, _ = decode["side"].(string)
	order.Stp, _ = decode["stp"].(string)
	order.Type, _ = decode["type"].(string)
	order.TimeInForce, _ = decode["time_in_force"].(string)
	order.PostOnly, _ = decode["post_only"].(bool)
	order.CreatedAt, _ = decode["created_at"].(string)
	temp, _ = decode["fill_fees"].(string)
	order.FillFees, _ = strconv.ParseFloat(temp, 64)
	temp, _ = decode["fillled_size"].(string)
	order.FilledSize, _ = strconv.ParseFloat(temp, 64)
	temp, _ = decode["executed_value"].(string)
	order.ExecutedValue, _ = strconv.ParseFloat(temp, 64)
	order.Status, _ = decode["status"].(string)
	order.Settled, _ = decode["settled"].(bool)
	return order, nil
}
