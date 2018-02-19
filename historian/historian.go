
package historian

import (
	"../gdax"
    "database/sql"
)

type Account struct {
	Id int64
	Funds float64
}

type Period struct {
	Time int64
	Closing float64
}

func exec(db *sql.DB, query string) {
    statement, e1 := db.Prepare(query)
    ok(e1)
    _ , e2 := statement.Exec()
    ok(e2)
}

func getId(db *sql.DB, query string) (int64) {
    statement, e1 := db.Prepare(query)
    ok(e1)
    result, e2 := statement.Exec()
    ok(e2)
    id, e3 := result.LastInsertId()
    ok(e3)
    return id
}

func CreateDb(db *sql.DB) {
    exec(db, "create table accounts (id integer primary key autoincrement, funds real);")
    exec(db, "create table book (product text, unix integer, buy integer, price real, size real, complete integer);")
    exec(db, "create table btc_usd (unix integer unique, closing real);")
    exec(db, "create table eth_usd (unix integer unique, closing real);")
    exec(db, "create table ltc_usd (unix integer unique, closing real);")
}

func NewAccount(db *sql.DB) (int64) {
    return getId(db, "insert into accounts default values")
}

func GetAccounts(db *sql.DB) ([]Account) {
	rows, e := db.Query("select * from accounts")
	ok(e)
	accounts := make([]Account, 0)
	var id int64
	var funds float64
	for rows.Next() {
		e = rows.Scan(&id, &funds)
 		ok(e)
		a := Account{}
		a.Id = id
		a.Funds = funds
		accounts = append(accounts, a)
	}
	rows.Close()
	return accounts
}

func ArchiveBtcUsd(db *sql.DB, history []gdax.Candle) {
	for i := 0; i < len(history); i++ {
		statement, e1 := db.Prepare("insert or ignore into btc_usd(unix, closing) select ?, ?")
		ok(e1)
		_ , e2 := statement.Exec(history[i].Time, history[i].Closing)
		ok(e2)
	}
}

func GetBtcUsd(db *sql.DB) ([]Period) {
	rows, e := db.Query("select * from btc_usd")
	ok(e)
	periods := make([]Period, 0)
	var unix int64
	var closing float64
	for rows.Next() {
		e = rows.Scan(&unix, &closing)
 		ok(e)
		p := Period{unix, closing}
		periods = append(periods, p)
	}
	rows.Close()
	return periods
}

func ok(e error) {
	if e != nil {
		panic(e)
	}
}
