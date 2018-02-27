package historian

import (
	"database/sql"

	"../gdax"
)

// Account account
type Account struct {
	ID    int64
	Funds float64
}

// Period time and closing price
type Period struct {
	Unix    int64
	Closing float64
}

func exec(db *sql.DB, query string) {
	statement, e1 := db.Prepare(query)
	ok(e1)
	_, e2 := statement.Exec()
	ok(e2)
}

func getID(db *sql.DB, query string) int64 {
	statement, e1 := db.Prepare(query)
	ok(e1)
	result, e2 := statement.Exec()
	ok(e2)
	id, e3 := result.LastInsertId()
	ok(e3)
	return id
}

// CreateDb creates the sqlite database
func CreateDb(db *sql.DB) {
	exec(db, "create table accounts (id integer primary key autoincrement, funds real);")
	exec(db, "create table book (product text, unix integer, buy integer, price real, size real, complete integer);")
	exec(db, "create table btc_usd (unix integer unique, closing real);")
	exec(db, "create table eth_usd (unix integer unique, closing real);")
	exec(db, "create table ltc_usd (unix integer unique, closing real);")
}

// NewAccount inserts a new account
func NewAccount(db *sql.DB) int64 {
	return getID(db, "insert into accounts default values")
}

// GetAccounts gets all accounts
func GetAccounts(db *sql.DB) []Account {
	rows, e := db.Query("select * from accounts")
	ok(e)
	accounts := make([]Account, 0)
	var id int64
	var funds float64
	for rows.Next() {
		e = rows.Scan(&id, &funds)
		ok(e)
		a := Account{}
		a.ID = id
		a.Funds = funds
		accounts = append(accounts, a)
	}
	rows.Close()
	return accounts
}

// ArchiveBtcUsd inserts history of bitcoin
func ArchiveBtcUsd(db *sql.DB, history []gdax.Candle) {
	for i := 0; i < len(history); i++ {
		statement, e1 := db.Prepare("insert or ignore into btc_usd(unix, closing) select ?, ?")
		ok(e1)
		_, e2 := statement.Exec(history[i].Time, history[i].Closing)
		ok(e2)
	}
}

// GetBtcUsd queries history of bitcoin
func GetBtcUsd(db *sql.DB, interval, from, to int64) []Period {
	rows, e := db.Query("select * from btc_usd where unix > ? and unix < ? order by unix", from, to)
	ok(e)
	timeRecords := make([]Period, 0)
	var unix int64
	var closing float64
	for rows.Next() {
		e = rows.Scan(&unix, &closing)
		ok(e)
		p := Period{unix, closing}
		timeRecords = append(timeRecords, p)
	}
	rows.Close()

	currentUnix := from
	lastRecordUnix := interval
	periods := make([]Period, 0)
	for i := 0; i < len(timeRecords); i++ {
		recordUnix := timeRecords[i].Unix
		if recordUnix%interval < lastRecordUnix {
			p := Period{unix, closing}
			periods = append(periods, p)

			currentUnix += interval
			if currentUnix >= to {
				break
			}
		}
		lastRecordUnix = recordUnix
	}
	return periods
}

func ok(e error) {
	if e != nil {
		panic(e)
	}
}
