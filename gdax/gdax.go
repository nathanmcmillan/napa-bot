package gdax

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

const (
	api       = "https://api.gdax.com"
	apiSocket = "wss://ws-feed.gdax.com"
	get       = "GET"
	post      = "POST"
)

func request(method, url string, body io.Reader) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "napa")
	return client, request, nil
}

func publicRequest(method, url string) ([]byte, error) {
	client, request, err := request(method, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func privateRequest(method, site, path, body string, auth *Authentication) ([]byte, error) {
	var data io.Reader
	if body != "" {
		message, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		data = bytes.NewReader(message)
	}

	client, request, err := request(method, site+path, data)
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	what := timestamp + method + path + body
	base64key, err := base64.StdEncoding.DecodeString(auth.Secret)
	if err != nil {
		return nil, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(what))
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))

	request.Header.Add("CB-ACCESS-KEY", auth.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", auth.Passphrase)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

// GetCurrencies list of currencies
func GetCurrencies() (interface{}, error) {
	body, err := publicRequest(get, api+"/currencies")
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
func GetBook(product string) (interface{}, error) {
	body, err := publicRequest(get, api+"/products/"+product+"/book?level=2")
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
func GetTicker(product string) (interface{}, error) {
	body, err := publicRequest(get, api+"/products/"+product+"/ticker")
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
func GetTrades(product string) (interface{}, error) {
	body, err := publicRequest(get, api+"/products/"+product+"/trades")
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
func GetHistory(product, start, end, granularity string) (*CandleList, error) {
	body, err := publicRequest(get, api+"/products/"+product+"/candles?start="+start+"&end="+end+"&granularity="+granularity)
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
func GetStats(product string) (interface{}, error) {
	body, err := publicRequest(get, api+"/products/"+product+"/stats")
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
func GetAccounts(private *Authentication) ([]Account, error) {
	body, err := privateRequest(get, api, "/accounts", "", private)
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
		account.Balance, _ = values["balance"].(string)
		account.Available, _ = values["available"].(string)
		account.Hold, _ = values["hold"].(string)
		account.ProfileID, _ = values["profile_id"].(string)

		accounts = append(accounts, account)
	}
	return accounts, nil
}

// PlaceOrder send a buy or sell order
func PlaceOrder(private *Authentication, js string) ([]Account, error) {
	return nil, errors.New("not implemented yet")
}

// ListOrders get open orders
func ListOrders(private *Authentication) (map[string][]*Order, error) {
	body, err := privateRequest(get, api, "/orders", "", private)
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
