package gdax

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
)

const (
	api = "https://api.gdax.com"
	get = "GET"
	post = "POST"
)

type Candle struct {
	Time float64
	Low float64
	High float64
	Open float64
	Closing float64
	Volume float64
}

func ok(e error) {
	if e != nil {
		panic(e)
	}
}

func request(method, url string) (*http.Client, *http.Request) {
	client := &http.Client{}
	request, e := http.NewRequest(method, url, nil)
	ok(e)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "napa")
	return client, request
}

func getRequest(url string) ([]byte) {
	client, request := request(get, url)
	response, e := client.Do(request)
	ok(e)
	body, e := ioutil.ReadAll(response.Body)
	ok(e)
	response.Body.Close()
	return body
}

func GetCurrencies() (interface{}) {
	body := getRequest(api + "/currencies")
	var decode interface{}
	e := json.Unmarshal(body, &decode)
	ok(e)
	return decode
}

func GetBook(product string) (interface{}) {
	body := getRequest(api + "/products/" + product + "/book?level=2")
	var decode interface{}
	e := json.Unmarshal(body, &decode)
	ok(e)
	return decode
}

func GetTicker(product string) (interface{}) {
	body := getRequest(api + "/products/" + product + "/ticker")
	var decode interface{}
	e := json.Unmarshal(body, &decode)
	ok(e)
	return decode
}

func GetTrades(product string) (interface{}) {
	body := getRequest(api + "/products/" + product + "/trades")
	var decode interface{}
	e := json.Unmarshal(body, &decode)
	ok(e)
	return decode
}

func GetHistory(product, start, end, granularity string) ([]Candle) {
	body := getRequest(api + "/products/" + product + "/candles?start=" + start + "&end=" + end + "&granularity=" + granularity)
	var decode []interface{}
	e := json.Unmarshal(body, &decode)
	ok(e)
	candles := make([]Candle, 0)
	for _ , list := range decode {
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

func GetStats(product string) (interface{}) {
	body := getRequest(api + "/products/" + product + "/stats")
	var decode interface{}
	e := json.Unmarshal(body, &decode)
	ok(e)
	return decode
}
