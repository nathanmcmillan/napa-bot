package gdax

import (
	"time"
	"strconv"
	"fmt"
)

const (
	candleLimit = int64(350)
)

// Polling sends poll requests to exchange
func Polling(auth *Authentication, settings *Settings, messages chan interface{}) error {
	interval := time.Second * time.Duration(settings.TimeInterval)
	intervalString := strconv.FormatInt(settings.TimeInterval, 10)
	next := time.Now()
	for {
		difference := time.Now().Sub(next)
		fmt.Println("time difference", difference)
		if difference > 0 {
			time.Sleep(difference)
			continue
		}
		
		for i := 0; i < len(settings.Products); i++ {
			product := settings.Products[i]
			fmt.Println("history product", product)
			start := time.Now().Add(-time.Second * time.Duration(candleLimit*settings.TimeInterval)).Format(time.RFC3339)
			end := time.Now().Format(time.RFC3339)
			history, err := GetHistory(product, start, end, intervalString)
			if err != nil {
				panic(err)
			}
			fmt.Println("History:", history)
			messages <- history
			
			time.Sleep(time.Second)
		}
		
		next = time.Now().Add(interval)
	}
}
