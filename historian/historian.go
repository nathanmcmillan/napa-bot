package historian

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"../gdax"
)

// Account record
type Account struct {
	ID    int64
	Funds float64
}

// Candle record
type Candle struct {
	Time    int64
	Low     float64
	High    float64
	Open    float64
	Closing float64
	Volume  float64
}

func exec(db *sql.DB, query string) error {
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}

func getID(db *sql.DB, query string) (int64, error) {
	statement, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	result, err := statement.Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
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
func NewAccount(db *sql.DB) (int64, error) {
	return getID(db, "insert into accounts default values")
}

// GetAccounts get list of accounts from database
func GetAccounts(db *sql.DB) ([]Account, error) {
	rows, err := db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}
	accounts := make([]Account, 0)
	var id int64
	var funds float64
	for rows.Next() {
		err = rows.Scan(&id, &funds)
		if err != nil {
			return nil, err
		}
		a := Account{}
		a.ID = id
		a.Funds = funds
		accounts = append(accounts, a)
	}
	rows.Close()
	return accounts, nil
}

// ArchiveBtcUsd create database records of btc usd
func ArchiveBtcUsd(db *sql.DB, candle []gdax.Candle) error {
	for i := 0; i < len(candle); i++ {
		statement, err := db.Prepare("insert or ignore into btc_usd(unix, low, high, open, closing, volume) select ?, ?, ?, ?, ?, ?")
		if err != nil {
			return err
		}
		current := candle[i]
		_, err = statement.Exec(current.Time, current.Low, current.High, current.Open, current.Closing, current.Volume)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetBtcUsd queries history of bitcoin
func GetBtcUsd(db *sql.DB, interval, from, to int64) ([]*Candle, error) {
	if to < from {
		return nil, errors.New("bad range")
	}
	rows, err := db.Query("select * from btc_usd where unix > ? and unix < ? order by unix", from, to)
	if err != nil {
		return nil, err
	}
	allCandles := make([]*Candle, 0)
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
		candle := &Candle{unix, low, high, open, closing, volume}
		allCandles = append(allCandles, candle)
	}
	rows.Close()

	indexOffset := from / interval
	numIndices := to/interval - indexOffset
	candles := make([]*Candle, numIndices)

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
