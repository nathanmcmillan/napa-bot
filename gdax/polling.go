package gdax

import (
	"time"
	"strconv"
	"fmt"
)

// Polling sends poll requests to exchange
func Polling(auth *Authentication, settings *Settings, messages chan interface{}) error {
	next := time.Now()
	interval := time.Second * time.Duration(settings.TimeInterval)
	for {
		time.Sleep(time.Second)
		if time.Now().Before(next) {
			continue
		}
		
		product := "BTC-USD"
		limit := int64(128)
		start := time.Now().Add(-time.Second * time.Duration(limit*settings.TimeInterval)).Format(time.RFC3339)
		end := time.Now().Format(time.RFC3339)
		history, err := GetHistory(product, start, end, strconv.FormatInt(settings.TimeInterval, 10))
		if err != nil {
			panic(err)
		}
		fmt.Println("History:", history)
		
		messages <- history
		next = time.Now().Add(interval)
	}
}
