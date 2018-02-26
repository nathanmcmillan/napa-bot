package gdax

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
	"io"
	"fmt"
	"bytes"
)

const (
	api  = "https://api.gdax.com"
	get  = "GET"
	post = "POST"
)

// Authentication private data
type Authentication struct {
	Key string
	Secret string
	Passphrase string
}

// Profile exchange account product data
type Profile struct {
	ID    string
	Currency     string
	Balance    string
	Available    string
	Hold string
	ProfileID  string
}

// Candle product data
type Candle struct {
	Time    float64
	Low     float64
	High    float64
	Open    float64
	Closing float64
	Volume  float64
}

// Order an order placed on the exchange
type Order struct {
	ID string
	Price string
	Size string
	Product string
	Side string
	Stp string
	Type string
	TimeInForce string
	PostOnly bool
	CreatedAt string
	FillFees string
	FilledSize string
	ExecutedValue string
	Status string
	Settled bool
}

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

func getRequest(url string) ([]byte, error) {
	client, request, err := request(get, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	return body, nil
}

// GetCurrencies list of currencies
func GetCurrencies() (interface{}, error) {
	body, err := getRequest(api + "/currencies")
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
	body, err := getRequest(api + "/products/" + product + "/book?level=2")
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
	body, err := getRequest(api + "/products/" + product + "/ticker")
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
	body, err := getRequest(api + "/products/" + product + "/trades")
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
func GetHistory(product, start, end, granularity string) ([]Candle, error) {
	body, err := getRequest(api + "/products/" + product + "/candles?start=" + start + "&end=" + end + "&granularity=" + granularity)
	if err != nil {
		return nil, err
	}
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	candles := make([]Candle, 0)
	for _, list := range decode {
		values := list.([]interface{})
		candle := Candle{}
		candle.Time = values[0].(float64)
		candle.Low = values[1].(float64)
		candle.High = values[2].(float64)
		candle.Open = values[3].(float64)
		candle.Closing = values[4].(float64)
		candle.Volume = values[5].(float64)
		candles = append(candles, candle)
	}
	return candles, nil
}

// GetStats map of 24 hour product statistics
func GetStats(product string) (interface{}, error) {
	body, err := getRequest(api + "/products/" + product + "/stats")
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
func GetAccounts(private *Authentication) ([]Profile, error) {
	
	method := get
	url := api + "/accounts"
	
	client, request, err := request(method, url, nil)
	if err != nil {
		return nil, err
	}
	
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + url
	base64key, err := base64.StdEncoding.DecodeString(private.Secret)
	if err != nil {
		return nil, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	
	request.Header.Add("CB-ACCESS-KEY", private.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", private.Passphrase)
	
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil, err
	}
	profiles := make([]Profile, 0)
	for _, list := range decode {
		values := list.(map[string]interface{})
		profile := Profile{}
		profile.ID, _ = values["id"].(string)
		profile.Currency, _ = values["currency"].(string)
		profile.Balance, _ = values["balance"].(string)
		profile.Available, _ = values["available"].(string)
		profile.Hold, _ = values["hold"].(string)
		profile.ProfileID, _ = values["profile_id"].(string)
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// PlaceOrder send a buy or sell order
func PlaceOrder(private *Authentication, js string) ([]Profile, error) {	
	method := post
	url := api + "/orders"
	data, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}
	
	client, request, err := request(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + url + js
	base64key, err := base64.StdEncoding.DecodeString(private.Secret)
	if err != nil {
		return nil, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	
	request.Header.Add("CB-ACCESS-KEY", private.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", private.Passphrase)
	
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil, nil
	}
	profiles := make([]Profile, 0)
	for _, list := range decode {
		values := list.(map[string]interface{})
		profile := Profile{}
		profile.ID, _ = values["id"].(string)
		profile.Currency, _ = values["currency"].(string)
		profile.Balance, _ = values["balance"].(string)
		profile.Available, _ = values["available"].(string)
		profile.Hold, _ = values["hold"].(string)
		profile.ProfileID, _ = values["profile_id"].(string)
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// ListOrders get open orders
func ListOrders(private *Authentication) ([]Profile, error) {
	method := get
	url := api + "/orders"
	
	client, request, err := request(method, url, nil)
	if err != nil {
		return nil, err
	}
	
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + url
	base64key, err := base64.StdEncoding.DecodeString(private.Secret)
	if err != nil {
		return nil, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	
	request.Header.Add("CB-ACCESS-KEY", private.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", private.Passphrase)
	
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil, err
	}
	profiles := make([]Profile, 0)
	for _, list := range decode {
		values := list.(map[string]interface{})
		profile := Profile{}
		profile.ID, _ = values["id"].(string)
		profile.Currency, _ = values["currency"].(string)
		profile.Balance, _ = values["balance"].(string)
		profile.Available, _ = values["available"].(string)
		profile.Hold, _ = values["hold"].(string)
		profile.ProfileID, _ = values["profile_id"].(string)
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
