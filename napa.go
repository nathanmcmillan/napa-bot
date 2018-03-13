package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"./datastore"
	"./gdax"
	"./parse"
	"./trader"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type listenLock struct {
	mu   *sync.Mutex
	todo bool
}

const (
	databaseDriver   = "sqlite3"
	databaseName     = "./napa.db"
	databaseTestName = "./napa_test.db"
	databaseSQL      = "./napa.sql"
)

var (
	indexFileHTML []byte
	indexFileJS   []byte
)

func input() {
	if os.Args[1] == "install" {
		install()
	} else if os.Args[1] == "server" {
		server()
	} else if os.Args[1] == "fund" {
		if len(os.Args) < 4 {
			fmt.Println("fund [product] [usd]")
			return
		}
		fund(os.Args[2], os.Args[3])
	} else if os.Args[1] == "list" {
		list()
	} else if os.Args[1] == "buy" {
		if len(os.Args) < 4 {
			fmt.Println("buy [product] [usd]")
			return
		}
		buy(os.Args[2], os.Args[3])
	} else if os.Args[1] == "sell" {
		if len(os.Args) < 4 {
			fmt.Println("sell [product] [amount]")
			return
		}
		sell(os.Args[2], os.Args[3])
	}
}

func buy(product, usd string) {
	fmt.Println("buy", product, "$", usd)
	fund, err := strconv.ParseFloat(usd, 64)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open(databaseDriver, databaseName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	settings := settings()
	auth := authentication()
	rest := gdax.NewRest(auth)
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	trade := trader.NewTrader(db, rest, settings, signals)
	order, err := trade.PlaceMarketBuy(product, fund)
	if err == nil {
		fmt.Println(order)
	}  else {
		fmt.Println(err)	
	}
}

func sell(product, amount string) {
	fmt.Println("sell", product, "amount", amount)
	_, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
}

func fund(product, usd string) {
	fmt.Println("funding", product, "$", usd)
	amount, err := strconv.ParseFloat(usd, 64)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open(databaseDriver, databaseName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	datastore.NewAccount(db, product, amount)
}

func list() {
	db, err := sql.Open(databaseDriver, databaseName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	accounts, err := datastore.QueryAccounts(db)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(accounts); i++ {
		account := accounts[i]
		fmt.Println("account", account.ID, account.Product, account.Funds)
	}
}

func install() {
	fmt.Println("deleting databases")
	err := os.Remove(databaseName)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	err = os.Remove(databaseTestName)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	fmt.Println("creating database")
	db, err := sql.Open(databaseDriver, databaseName)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = datastore.RunFile(db, databaseSQL)
	if err != nil {
		panic(err)
	}
	fmt.Println("creating test database")
	dbTest, err := sql.Open(databaseDriver, databaseTestName)
	if err != nil {
		panic(err)
	}
	defer dbTest.Close()
	err = datastore.RunFile(dbTest, databaseSQL)
	if err != nil {
		panic(err)
	}
}

func server() {
	fmt.Println("loading files")
	file, err := os.Open("napa.html")
	if err != nil {
		panic(err)
	}
	indexFileHTML, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file, err = os.Open("napa.js")
	if err != nil {
		panic(err)
	}
	indexFileJS, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening and serving")
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/napa.js", indexJS)
	http.HandleFunc("/websocket", clientSocket)
	http.ListenAndServe(":80", nil)
}

func indexHTML(writer http.ResponseWriter, request *http.Request) {
	writer.Write(indexFileHTML)
}

func indexJS(writer http.ResponseWriter, request *http.Request) {
	writer.Write(indexFileJS)
}

func exchangeSocket(clientSocket *websocket.Conn, lock *sync.Mutex, listen *listenLock) error {
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

func settings() *gdax.Settings {
	// load files
	public, err := getFile("./public.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	s := &gdax.Settings{}
	s.Products = parse.StringList(public, "products")
	s.Channels = parse.StringList(public, "channels")
	s.Seconds = parse.Integer(public, "seconds")
	s.EmaShort = parse.Integer(public, "ema-short")
	s.EmaLong = parse.Integer(public, "ema-long")
	s.RsiPeriods = parse.Integer(public, "rsi")
	return s
}

func authentication() *gdax.Authentication {
	private, err := getFile("../private.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	a := &gdax.Authentication{}
	a.Key = parse.Text(private, "key")
	a.Secret = parse.Text(private, "secret")
	a.Passphrase = parse.Text(private, "phrase")
	return a
}

func main() {
	fmt.Println("napa bot")

	// logging
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(logFile)

	if len(os.Args) > 1 {
		input()
		return
	}

	settings := settings()
	auth := authentication()
	rest := gdax.NewRest(auth)
	fmt.Println("settings", settings)

	// database
	db, err := sql.Open(databaseDriver, databaseName)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.Close()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// trade
	trade := trader.NewTrader(db, rest, settings, signals)
	trade.Run()

	time.Sleep(time.Second)
}
