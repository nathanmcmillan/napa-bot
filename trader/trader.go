package trader

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"../analyst"
	"../gdax"
)

// Run core loop
func Run(analysis *analyst.Analyst, db *sql.DB, auth *gdax.Authentication, messages chan interface{}) {

	// get account information
	accounts, err := gdax.GetAccounts(auth)
	if err != nil {
		panic(err)
	}
	fmt.Println("Accounts:", accounts)

	/*
		// get history
		product := "BTC-USD"
		limit := int64(128)
		start := time.Now().Add(-time.Second * time.Duration(limit*analysis.TimeInterval)).Format(time.RFC3339)
		end := time.Now().Format(time.RFC3339)
		history, err := gdax.GetHistory(product, start, end, strconv.FormatInt(analysis.TimeInterval, 10))
		if err != nil {
			panic(err)
		}
		fmt.Println("History:", history)

		// archive history
		historian.ArchiveBtcUsd(db, history)

		// analyze history
		from := time.Now().Add(-time.Second * time.Duration(limit*analysis.TimeInterval)).Unix()
		to := time.Now().Unix()
		candles, err := historian.GetBtcUsd(db, analysis.TimeInterval, from, to)
		if err != nil {
			panic(err)
		}
		fmt.Println("MACD", analyst.Macd(analysis.EmaShort, analysis.EmaLong, candles))
		fmt.Println("RSI", analyst.Rsi(analysis.RsiPeriods, candles)) */

	for {
		rawJs := <-messages
		var js interface{}
		err := json.Unmarshal([]byte(rawJs), &js)
		if err != nil {
			continue
		}
		message, ok := js.(map[string]interface{})
		if !ok {
			continue
		}
		messageType, ok := message["uid"].(string)
		if !ok {
			continue
		}
		switch messageType {
		case "ticker":
			fmt.Println("got ticker", rawJs)
		case "snapshot":
			fmt.Println("got snapshot", rawJs)
		case "l2update":
			fmt.Println("got l2 update", rawJs)
		}

	}
}
