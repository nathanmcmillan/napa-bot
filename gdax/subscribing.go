package gdax

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// ExchangeSocket dials websocket to exchange
func ExchangeSocket(products, channels []string, messages chan interface{}) error {

	fmt.Println("dialing exchange")
	connection, _, err := websocket.DefaultDialer.Dial(apiSocket, nil)
	if err != nil {
		return err
	}
	defer connection.Close()

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
	if err != nil {
		return err
	}

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
			rawJs := fmt.Sprintf(`{"uid":"ticker", "time":"%s", "product_id":"%s", "price":"%s", "side":"%s"}`, time, productID, price, side)
			messages <- rawJs
		case "snapshot":
			productID, _ := message["product_id"].(string)
			rawJs := fmt.Sprintf(`{"uid":"snapshot", "product_id":"%s"}`, productID)
			messages <- rawJs
		case "l2update":
			productID, _ := message["product_id"].(string)
			rawJs := fmt.Sprintf(`{"uid":"l2update", "product_id":"%s"}`, productID)
			messages <- rawJs
		}
	}
	return nil
}
