package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"./decoding"
	"./analyst"
	"./gdax"
	"./historian"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type listenLock struct {
	mu   *sync.Mutex
	todo bool
}

const (
	databaseDriver = "sqlite3"
	databaseName   = "./napa.db"
)

var (
	contents []byte
)

func app() {
	// db, e := sql.Open(databaseDriver, databaseName)
	// ok(e)

	/* history := gdax.GetHistory("BTC-USD", "2018-02-17", "2018-02-18", "3600")
	historian.ArchiveBtcUsd(db, history)
	periods := historian.GetBtcUsd(db)
	fmt.Println("MACD", analyst.MovingAverageConvergenceDivergence(6, 12, periods))
	fmt.Println("RSI", analyst.RelativeStrengthIndex(7, periods)) */

	/* gdax.GetCurrencies()
		gdax.GetBook(product)
		gdax.GetTicker(product)
		gdax.GetTrades(product)
	 	gdax.GetHistory(product, "2018-02-16", "2018-02-17", "3600")
		gdax.GetStats(product) */

	// db.Close()
}

func install() {
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

func indexPage(writer http.ResponseWriter, request *http.Request) {
	writer.Write(contents)
}

func exchangeSocket(clientSocket *websocket.Conn, lock *sync.Mutex, listen *listenLock) {
	fmt.Println("connecting to exchange")
	url := "wss://ws-feed.gdax.com"
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	ok(err)
	js := json.RawMessage(`{"type":"subscribe", "product_ids":["BTC-USD"], "channels":["ticker"]}`)
	err = connection.WriteJSON(js)
	ok(err)
	fmt.Println("listening to exchange")
	for {
		var proceed bool
		listen.mu.Lock()
		proceed = listen.todo
		listen.mu.Unlock()
		if !proceed {
			break
		}
		var js interface{}
		err := connection.ReadJSON(&js)
		if err != nil {
			fmt.Println(err)
			break
		}
		message, ok := js.(map[string]interface{})
		if !ok {
			continue
		}
		messageType, ok := message["type"].(string)
		if !ok {
			continue
		}
		if messageType == "ticker" {
			time, _ := message["time"].(string)
			productID, _ := message["product_id"].(string)
			price, _ := message["price"].(string)
			side, _ := message["side"].(string)
			clientMessage := fmt.Sprintf(`{"uid":"ticker", "time":"%s", "product_id":"%s", "price":"%s", "side":"%s"}`, time, productID, price, side)
			go clientWrite(clientSocket, lock, clientMessage) // broadcast close exchange socket if fail ?
		}
	}
	connection.Close()
	fmt.Println("exchange connection closed")
}

func clientWrite(connection *websocket.Conn, lock *sync.Mutex, rawJs string) {
	js := json.RawMessage([]byte(rawJs))
	lock.Lock()
	err := connection.WriteJSON(js)
	lock.Unlock()
	fmt.Println("sent", rawJs)
	if err != nil {
		fmt.Println(err)
	}
}

func clientRead(connection *websocket.Conn) {
	var lock = &sync.Mutex{}
	var exchangeLock = &listenLock{}
	exchangeLock.mu = &sync.Mutex{}
	exchangeLock.todo = false
	for {
		var js interface{}
		err := connection.ReadJSON(&js)
		if err != nil {
			fmt.Println(err)
			connection.Close()
			break
		}
		message, ok := js.(map[string]interface{})
		if !ok {
			continue
		}
		uid, ok := message["uid"].(string)
		if !ok {
			continue
		}
		switch uid {
		case "sub-exchange":
			var current bool
			exchangeLock.mu.Lock()
			current = exchangeLock.todo
			exchangeLock.todo = true
			exchangeLock.mu.Unlock()
			if current {
				continue
			}
			go clientWrite(connection, lock, `{"uid":"log", "message":"subbing to exchange"}`)
			go exchangeSocket(connection, lock, exchangeLock)
		case "unsub-exchange":
			go clientWrite(connection, lock, `{"uid":"log", "message":"unsubbing from exchange"}`)
			exchangeLock.mu.Lock()
			exchangeLock.todo = false
			exchangeLock.mu.Unlock()
		}
	}
}

func clientSocket(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Origin") != "http://"+request.Host {
		http.Error(writer, "origin not allowed", 403)
		return
	}
	upgrader := websocket.Upgrader{}
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, "could not open websocket", 400)
		return
	}
	go clientRead(connection)
}

func main() {
	fmt.Println("napa bot")

	/*fmt.Println("loading files")
	file, err := os.Open("napa.html")
	ok(err)
	contents, err = ioutil.ReadAll(file)
	ok(err)

	fmt.Println("listening and serving")
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/websocket", clientSocket)
	http.ListenAndServe(":80", nil)*/
	
	public := getFile("./public.json")
	fmt.Println(public)
	
	var analysis = analyst.Analyst{}
	analysis.TimeInterval = decoding.GetInteger(public, "interval")
	analysis.EmaShort = decoding.GetInteger(public, "ema-short")
	analysis.EmaLong = decoding.GetInteger(public, "ema-long")
	analysis.RsiPeriods = decoding.GetInteger(public, "rsi")
	
	fmt.Println(analysis)
	
	gdax.GetAccounts()
}

func sleep(seconds int32) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func ok(e error) {
	if e != nil {
		panic(e)
	}
}

func getFile(path string) map[string]interface{} {
	file, err := os.Open(path)
	ok(err)
	contents, err := ioutil.ReadAll(file)
	ok(err)
	var decode interface{}
	err = json.Unmarshal(contents, &decode)
	ok(err)
	js, _ := decode.(map[string]interface{})
	return js
}
