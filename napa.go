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

	"./parse"
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

func install() (error) {
	fmt.Println("deleting database if exists")
	err := os.Remove(databaseName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Println("creating database")
	db, err := sql.Open(databaseDriver, databaseName)
	if err != nil {
		return err
	}
	historian.CreateDb(db)
	db.Close()
	return nil
}

func indexPage(writer http.ResponseWriter, request *http.Request) {
	writer.Write(contents)
}

func exchangeSocket(clientSocket *websocket.Conn, lock *sync.Mutex, listen *listenLock) (error) {
	fmt.Println("connecting to exchange")
	url := "wss://ws-feed.gdax.com"
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	js := json.RawMessage(`{"type":"subscribe", "product_ids":["BTC-USD"], "channels":["ticker"]}`)
	err = connection.WriteJSON(js)
	if err != nil {
		return err
	}
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
	return nil
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
	
	public, err := getFile("./public.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(public)
	
	private, err := getFile("../private.json")
	if err != nil {
		panic(err)
	}
	fmt.Println(private)
	
	analysis := analyst.Analyst{}
	analysis.TimeInterval = parse.Integer(public, "interval")
	analysis.EmaShort = parse.Integer(public, "ema-short")
	analysis.EmaLong = parse.Integer(public, "ema-long")
	analysis.RsiPeriods = parse.Integer(public, "rsi")
	fmt.Println(analysis)
	
	auth := gdax.Authentication{}
	auth.Key = parse.Text(private, "key")
	auth.Secret = parse.Text(private, "secret")
	auth.Passphrase = parse.Text(private, "passphrase")
	fmt.Println(auth)
	
	gdax.GetAccounts(&auth)
}

func sleep(seconds int32) {
	time.Sleep(time.Second * time.Duration(seconds))
}

func getFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var decode interface{}
	err = json.Unmarshal(contents, &decode)
	if err != nil {
		return nil, err
	}
	js, _ := decode.(map[string]interface{})
	return js, nil
}
