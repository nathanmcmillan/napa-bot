package gdax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	api       = "https://api.gdax.com"
	apiSocket = "wss://ws-feed.gdax.com"
	get       = "GET"
	post      = "POST"
)

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

// ExchangeSocket dials websocket to exchange
func ExchangeSocket(products, channels []string) {
	fmt.Println("connecting to exchange")
	connection, _, err := websocket.DefaultDialer.Dial(apiSocket, nil)
	ok(err)

	var productList bytes.Buffer
	numProducts := len(products)
	for i := 0; i < numProducts; i++ {
		productList.WriteString(`"`)
		productList.WriteString(products[i])
		productList.WriteString(`"`)
		if i+1 < numProducts {
			productList.WriteString(`, `)
		}
	}

	var channelList bytes.Buffer
	numChannels := len(channels)
	for i := 0; i < numChannels; i++ {
		channelList.WriteString(`"`)
		channelList.WriteString(channels[i])
		channelList.WriteString(`"`)
		if i+1 < numChannels {
			channelList.WriteString(`, `)
		}
	}

	rawJs := fmt.Sprintf(`{"type":"subscribe", "product_ids":[%s], "channels":[%s]}`, productList.String(), channelList.String())
	js := json.RawMessage(rawJs)
	err = connection.WriteJSON(js)
	ok(err)
	fmt.Println("listening to exchange")
	for {

		var js interface{}
		err := connection.ReadJSON(&js)
		if err != nil {
			fmt.Println(err)
			break
		}
		message, ok := js.(map[string]interface{})
		if !ok {
			continue
		}
		messageType, ok := message["type"].(string)
		if !ok {
			continue
		}
		switch messageType {
		case "ticker":
			time, _ := message["time"].(string)
			productID, _ := message["product_id"].(string)
			price, _ := message["price"].(string)
			side, _ := message["side"].(string)
			fmt.Println(fmt.Sprintf(`{"uid":"ticker", "time":"%s", "product_id":"%s", "price":"%s", "side":"%s"}`, time, productID, price, side))
		case "snapshot":
			productID, _ := message["product_id"].(string)
			fmt.Println(fmt.Sprintf(`{"uid":"snapshot", "product_id":"%s"}`, productID))
		case "l2update":
			productID, _ := message["product_id"].(string)
			fmt.Println(fmt.Sprintf(`{"uid":"l2update", "product_id":"%s"}`, productID))
		}
	}
	connection.Close()
	fmt.Println("exchange connection closed")
}

func ok(err error) {
	if err != nil {
		panic(err)
	}
}
