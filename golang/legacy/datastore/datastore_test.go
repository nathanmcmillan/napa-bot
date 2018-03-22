package datastore

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseDriver   = "sqlite3"
	databaseTestName = "../napa_test.db"
)

// TestOrders test orders
func TestOrders(t *testing.T) {
	db, err := sql.Open(databaseDriver, databaseTestName)
	if err != nil {
		t.Error("error", err)
	}
	defer db.Close()

	product := "LTC-USD"
	price := float64(10200)
	size := float64(0.05)

	err = ArchiveOrder(db, product, price, size)
	if err != nil {
		t.Error("error", err)
	}

	orderMap, err := QueryOrders(db)
	if orderMap == nil || err != nil {
		t.Error("error", err)
	}
	orders := orderMap[product]
	if orders == nil || len(orders) == 0 {
		t.Error("empty after insert")
	}
	order := orders[0]
	if order.Product != product || order.Price != price || order.Size != size {
		t.Error("invalid data")
	}

	err = RemoveOrder(db, order.ID)
	if err != nil {
		t.Error("error", err)
	}
	orderMap, err = QueryOrders(db)
	if orderMap == nil || err != nil {
		t.Error("error", err)
	}
	if orderMap[product] != nil {
		t.Error("not empty after delete")
	}
}
