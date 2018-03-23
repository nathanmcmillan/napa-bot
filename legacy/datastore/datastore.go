package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"../gdax"
)

func exec(db *sql.DB, query string) error {
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}

// RunFile executes all sql statements
func RunFile(db *sql.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	statements := strings.Split(string(contents), ";")
	for i := 0; i < len(statements); i++ {
		query := strings.TrimSpace(statements[i])
		if query == "" {
			continue
		}
		err = exec(db, query)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAccount create database account
func NewAccount(db *sql.DB, product string, funds float64) (int64, error) {
	statement, err := db.Prepare("insert into accounts(product, funds) select ?, ?")
	if err != nil {
		return 0, err
	}
	result, err := statement.Exec(product, funds)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// QueryAccounts get list of accounts from database
func QueryAccounts(db *sql.DB) ([]*Account, error) {
	rows, err := db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, 0)
	var id int64
	var product string
	var funds float64
	for rows.Next() {
		err = rows.Scan(&id, &product, &funds)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &Account{id, product, funds})
	}
	rows.Close()
	return accounts, nil
}

// ArchiveCoin inserts historical records of coin product
func ArchiveCoin(product string, db *sql.DB, candle []gdax.Candle) error {
	statement, err := db.Prepare("insert or ignore into history(unix, product, low, high, open, closing, volume) select ?, ?, ?, ?, ?, ?, ?")
	if err != nil {
		return err
	}
	for i := 0; i < len(candle); i++ {
		current := candle[i]
		_, err = statement.Exec(current.Time, product, current.Low, current.High, current.Open, current.Closing, current.Volume)
		if err != nil {
			return err
		}
	}
	return nil
}

// QueryCoin queries history of coin product
func QueryCoin(product string, db *sql.DB, interval, from, to int64) ([]*gdax.Candle, error) {
	if to < from {
		return nil, errors.New("bad range")
	}
	rows, err := db.Query("select * from history where product = ? and unix > ? and unix < ? order by unix", product, from, to)
	if err != nil {
		return nil, err
	}
	allCandles := make([]*gdax.Candle, 0)
	var unix int64
	var low float64
	var high float64
	var open float64
	var closing float64
	var volume float64
	for rows.Next() {
		err = rows.Scan(&unix, &low, &high, &open, &closing, &volume)
		if err != nil {
			return nil, err
		}
		candle := &gdax.Candle{unix, low, high, open, closing, volume}
		allCandles = append(allCandles, candle)
	}
	rows.Close()

	indexOffset := from / interval
	numIndices := to/interval - indexOffset
	candles := make([]*gdax.Candle, numIndices)

	fmt.Println(indexOffset, numIndices)
	for i := 0; i < len(allCandles); i++ {
		current := allCandles[i]
		currentIndex := current.Time/interval - indexOffset
		fmt.Println(currentIndex)
		if currentIndex < 0 {
			continue
		}
		if currentIndex > numIndices {
			break
		}
		candles[currentIndex] = current
		for j := i - 1; j > 0; j-- {
			if candles[j] == nil {
				candles[j] = current
			}
		}
	}
	return candles, nil
}

// ArchiveOrder inserts new order
func ArchiveOrder(db *sql.DB, exchange_id string) error {
	statement, err := db.Prepare("insert into orders(exchange_id) select ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(exchange_id)
	if err != nil {
		return err
	}
	return nil
}

// RemoveOrder remove existing order
func RemoveOrder(db *sql.DB, exchange_id string) error {
	statement, err := db.Prepare("delete from orders where exchange_id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(exchange_id)
	if err != nil {
		return err
	}
	return nil
}

// QueryOrders lists stored orders
func QueryOrders(db *sql.DB) ([]*Order, error) {
	rows, err := db.Query("select * from orders")
	if err != nil {
		return nil, err
	}
	orders := make([]*Order, 0)
	var id int64
	var exchange_id string
	for rows.Next() {
		err = rows.Scan(&id, &exchange_id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &Order{id, exchange_id})
	}
	rows.Close()
	return orders, nil
}
