package gdax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

// ExchangeSocket dials websocket to exchange
func ExchangeSocket(settings *Settings, messages chan interface{}) error {

	fmt.Println("dialing exchange")
	connection, _, err := websocket.DefaultDialer.Dial(apiSocket, nil)
	if err != nil {
		return err
	}
	defer connection.Close()

	var productList bytes.Buffer
	numProducts := len(settings.Products)
	for i := 0; i < numProducts; i++ {
		productList.WriteString(`"`)
		productList.WriteString(settings.Products[i])
		productList.WriteString(`"`)
		if i+1 < numProducts {
			productList.WriteString(`, `)
		}
	}

	var channelList bytes.Buffer
	numChannels := len(settings.Channels)
	for i := 0; i < numChannels; i++ {
		channelList.WriteString(`"`)
		channelList.WriteString(settings.Channels[i])
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
			ticker := Ticker{}
			temp, _ := message["time"].(string)
			ticker.Time, _ = strconv.ParseInt(temp, 10, 64)
			ticker.ProductID, _ = message["product_id"].(string)
			temp, _ = message["price"].(string)
			ticker.Price, _ = strconv.ParseFloat(temp, 64)
			ticker.Side, _ = message["side"].(string)
			messages <- ticker
		}
	}
	return nil
}
