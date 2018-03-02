package gdax

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// Ticker last match from exchange
type Ticker struct {
	Time      string
	ProductID string
	Price     string
	Side      string
}

// Snapshot snapshot of level 2 from exchange
type Snapshot struct {
	ProductID string
	Bids      []SnapshotTuple
	Asks      []SnapshotTuple
}

// SnapshotTuple tuple of price and size
type SnapshotTuple struct {
	Price string
	Size  string
}

// Update level 2 update from exchange
type Update struct {
	ProductID string
	Changes   []UpdateChange
}

// UpdateChange level 2 update details
type UpdateChange struct {
	Side  string
	Price string
	Size  string
}

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
			rawJs := Ticker{}
			rawJs.Time, _ = message["time"].(string)
			rawJs.ProductID, _ = message["product_id"].(string)
			rawJs.Price, _ = message["price"].(string)
			rawJs.Side, _ = message["side"].(string)
			messages <- rawJs
		case "snapshot":
			rawJs := Snapshot{}
			rawJs.ProductID, _ = message["product_id"].(string)
			rawJs.Bids = make([]SnapshotTuple, 0)
			rawJs.Asks = make([]SnapshotTuple, 0)
			list, ok := message["bids"].([]interface{})
			if !ok {
				continue
			}
			for i := 0; i < len(list); i++ {
				rawTuple, ok := list[i].([]interface{})
				if !ok {
					continue
				}
				tuple := SnapshotTuple{}
				tuple.Price, _ = rawTuple[0].(string)
				tuple.Size, _ = rawTuple[1].(string)
				rawJs.Bids = append(rawJs.Bids, tuple)
			}
			list, ok = message["asks"].([]interface{})
			if !ok {
				continue
			}
			for i := 0; i < len(list); i++ {
				rawTuple, ok := list[i].([]interface{})
				if !ok {
					continue
				}
				tuple := SnapshotTuple{}
				tuple.Price, _ = rawTuple[0].(string)
				tuple.Size, _ = rawTuple[1].(string)
				rawJs.Asks = append(rawJs.Asks, tuple)
			}
			messages <- rawJs
		case "l2update":
			rawJs := Update{}
			rawJs.ProductID, _ = message["product_id"].(string)
			rawJs.Changes = make([]UpdateChange, 0)
			list, ok := message["changes"].([]interface{})
			if !ok {
				continue
			}
			for i := 0; i < len(list); i++ {
				rawChange, ok := list[i].([]interface{})
				if !ok {
					continue
				}
				change := UpdateChange{}
				change.Side, _ = rawChange[0].(string)
				change.Price, _ = rawChange[1].(string)
				change.Size, _ = rawChange[2].(string)
				rawJs.Changes = append(rawJs.Changes, change)
			}
			messages <- rawJs
		}
	}
	return nil
}
