package gdax

import (
	"encoding/json"
	"time"
)

// Poll settings for polling
type Poll struct {
	OrderTime   int64
	HistoryTime int64
}

// Polling sends poll requests to exchange
func Polling(auth *Authentication, settings *Poll, messages chan interface{}) error {
	orderTime := time.Second * time.Duration(settings.OrderTime)
	historyTime := time.Second * time.Duration(settings.HistoryTime)
	nextOrder := time.Now()
	nextHistory := time.Now()
	for {
		time.Sleep(time.Second)

		if time.Now().After(nextOrder) {
			time.Sleep(time.Second)
			orders, err := ListOrders(auth)
			if err != nil {
				panic(err)
			}
			out, err := json.Marshal(orders)
			if err != nil {
				panic(err)
			}
			messages <- string(out)
			nextOrder = time.Now().Add(orderTime)
		}

		if time.Now().After(nextHistory) {
			time.Sleep(time.Second)
			nextHistory = time.Now().Add(historyTime)
		}
	}
}
