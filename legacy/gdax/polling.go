package gdax

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	candleLimit = int64(350)
)

// Polling sends poll requests to exchange
func (r *Rest) Polling(settings *Settings, messages chan interface{}, done chan bool) {
	if settings.EmaLong > candleLimit {
		panic(errors.New("ema out of range"))
	}
	interval := time.Second * time.Duration(settings.Seconds)
	beginning := -interval * time.Duration(settings.EmaLong)
	granularity := strconv.FormatInt(settings.Seconds, 10)
	for {
		known := false
		var wait time.Duration

		for i := 0; i < len(settings.Products); i++ {
			now := time.Now().UTC()
			start := time.Now().UTC().Add(beginning)
			fmt.Println("polling", settings.Products[i], "from", start, "to", now)
			history, err := r.GetHistory(settings.Products[i], start.Format(time.RFC3339), now.Format(time.RFC3339), granularity)
			if err != nil {
				log.Println(err)
				messages <- err
				time.Sleep(time.Second)
				continue
			}
			if len(history.List) > 0 {
				messages <- history
				if !known {
					candle := history.List[len(history.List)-1]
					at := time.Unix(candle.Time, 0)
					wait = interval - time.Now().Sub(at)
					known = true
				}
			}
		}

		if known {
			if wait < 0 {
				wait = interval
			}
			fmt.Println("poll thread sleeping for", wait)

			waitTimer := time.NewTimer(wait)
			select {
			case <-waitTimer.C:
				continue
			case <-done:
				waitTimer.Stop()
				fmt.Println("closing poll thread")
				return
			}
		}
	}
}
