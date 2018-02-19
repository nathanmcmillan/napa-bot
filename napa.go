
package main

import (
	"fmt"
	"time"
	"./gdax"
	"os"
	"./analyst"
	"./historian"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseDriver = "sqlite3"
	databaseName = "./napa.db"
)

func remake() {
	fmt.Println("deleting database if exists")
	e := os.Remove(databaseName)
	if e != nil && !os.IsNotExist(e) {
		panic(e)
	}
	fmt.Println("creating database")
	db, e := sql.Open(databaseDriver, databaseName)
 	ok(e)
	historian.CreateDb(db)
	db.Close()
}

func app() {
	db, e := sql.Open(databaseDriver, databaseName)
 	ok(e)
	
	history := gdax.GetHistory("BTC-USD", "2018-02-17", "2018-02-18", "3600")
	historian.ArchiveBtcUsd(db, history)
	periods := historian.GetBtcUsd(db)
	fmt.Println("MACD", analyst.MovingAverageConvergenceDivergence(6, 12, periods))
	fmt.Println("RSI", analyst.RelativeStrengthIndex(7, periods))
	
	/* gdax.GetCurrencies()
	gdax.GetBook(product)
	gdax.GetTicker(product)
	gdax.GetTrades(product)
 	gdax.GetHistory(product, "2018-02-16", "2018-02-17", "3600")
	gdax.GetStats(product) */
	
	db.Close()
}

func main() {
	fmt.Println("napa bot")
	arguments := os.Args
	if len(arguments) > 1 {
		com := arguments[1]
		if com == "remake" {
			remake()
		}
	} else {
		app()
	}
}

func sleep(seconds int32) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func ok(e error) {
	if e != nil {
		panic(e)
	}
}