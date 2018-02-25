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
	"fmt"
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

func request(method, url string) (*http.Client, *http.Request) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	ok(err)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "napa")
	return client, request
}

func getRequest(url string) []byte {
	client, request := request(get, url)
	response, err := client.Do(request)
	ok(err)
	body, err := ioutil.ReadAll(response.Body)
	ok(err)
	response.Body.Close()
	return body
}

// GetCurrencies list of currencies
func GetCurrencies() interface{} {
	body := getRequest(api + "/currencies")
	var decode interface{}
	err := json.Unmarshal(body, &decode)
	ok(err)
	return decode
}

// GetBook map of level 2 product books
func GetBook(product string) interface{} {
	body := getRequest(api + "/products/" + product + "/book?level=2")
	var decode interface{}
	err := json.Unmarshal(body, &decode)
	ok(err)
	return decode
}

// GetTicker map of product ticker
func GetTicker(product string) interface{} {
	body := getRequest(api + "/products/" + product + "/ticker")
	var decode interface{}
	err := json.Unmarshal(body, &decode)
	ok(err)
	return decode
}

// GetTrades map of product trades
func GetTrades(product string) interface{} {
	body := getRequest(api + "/products/" + product + "/trades")
	var decode interface{}
	err := json.Unmarshal(body, &decode)
	ok(err)
	return decode
}

// GetHistory list of candles for product history
func GetHistory(product, start, end, granularity string) []Candle {
	body := getRequest(api + "/products/" + product + "/candles?start=" + start + "&end=" + end + "&granularity=" + granularity)
	var decode []interface{}
	err := json.Unmarshal(body, &decode)
	ok(err)
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
	return candles
}

// GetStats map of 24 hour product statistics
func GetStats(product string) interface{} {
	body := getRequest(api + "/products/" + product + "/stats")
	var decode interface{}
	err := json.Unmarshal(body, &decode)
	ok(err)
	return decode
}

// GetAccounts map of account balances
func GetAccounts(private *Authentication) []Profile {
	
	method := get
	url := api + "/accounts"
	
	client, request := request(method, url)
	
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	
	message := timestamp + method + url
	base64key, err := base64.StdEncoding.DecodeString(private.Secret)
	ok(err)
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(message))
	ok(err)
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	
	request.Header.Add("CB-ACCESS-KEY", private.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", private.Passphrase)
	
	response, err := client.Do(request)
	ok(err)
	body, err := ioutil.ReadAll(response.Body)
	ok(err)
	response.Body.Close()
	
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil
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
	return profiles
}

// PlaceOrder send a buy or sell order
func PlaceOrder(private *Authentication, js string) []Profile {
	
	method := post
	url := api + "/orders"
	data := json.Marshal(js)
	
	client, request := request(method, url, data)
	
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	
	message := timestamp + method + url + js
	base64key, err := base64.StdEncoding.DecodeString(private.Secret)
	ok(err)
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(message))
	ok(err)
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	
	request.Header.Add("CB-ACCESS-KEY", private.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", private.Passphrase)
	
	response, err := client.Do(request)
	ok(err)
	body, err := ioutil.ReadAll(response.Body)
	ok(err)
	response.Body.Close()
	
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var decodeError interface{}
		err = json.Unmarshal(body, &decodeError)
		fmt.Println(decodeError)
		return nil
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
	return profiles
}

func ok(err error) {
	if err != nil {
		panic(err)
	}
}
