
package historian

import (
	"fmt"
	"../gdax"
    "database/sql"
)

// Account record
type Account struct {
	Id int64
	Funds float64
}

// Candle record
type Candle struct {
	Time int64
	Low float64
	High float64
	Open float64
	Closing float64
	Volume float64
}

func exec(db *sql.DB, query string) (error) {
    statement, err := db.Prepare(query)
    if err != nil {
		return err
	}
    _ , err = statement.Exec()
	return err
}

func getId(db *sql.DB, query string) (int64, error) {
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

func CreateDb(db *sql.DB) (error) {
	err := exec(db, "create table accounts (id integer primary key autoincrement, funds real);")
    if err != nil {
		return err
	}
	err = exec(db, "create table book (product text, unix integer, buy integer, price real, size real, complete integer);")
	if err != nil {
		return err
	}
	products := []string{"btc_usd, eth_usd, ltc_usd"}
	for i := 0; i < len(products); i++ {
		err = exec(db, fmt.Sprintf("create table %s (unix integer unique, low real, high real, open real, closing real, volume real);", products[i]))
		if err != nil {
			return err
		}
	}
	return nil
}

func NewAccount(db *sql.DB) (int64, error) {
    return getId(db, "insert into accounts default values")
}

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
		a.Id = id
		a.Funds = funds
		accounts = append(accounts, a)
	}
	rows.Close()
	return accounts, nil
}

func ArchiveBtcUsd(db *sql.DB, candle []gdax.Candle) (error) {
	for i := 0; i < len(candle); i++ {
		statement, err := db.Prepare("insert or ignore into btc_usd(unix, low, high, open, closing, volume) select ?, ?, ?, ?, ?, ?")
		if err != nil {
			return err
		}
		current := candle[i]
		_ , err = statement.Exec(current.Time, current.Low, current.High, current.Open, current.Closing, current.Volume)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBtcUsd(db *sql.DB) ([]Candle, error) {
	rows, err := db.Query("select * from btc_usd")
	if err != nil {
		return nil, err
	}
	candles := make([]Candle, 0)
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
		candle := Candle{unix, low, high, open, closing, volume}
		candles = append(candles, candle)
	}
	rows.Close()
	return candles, nil
}
