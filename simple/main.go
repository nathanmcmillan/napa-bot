package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	orderFile             = "orders.txt"
	orderBackupFile       = "orders_backup.txt"
	orderUpdateFile       = "orders_update.txt"
	orderUpdateBackupFile = "orders_update_backup.txt"
)

var (
	interrupt = false
	product   string
	auth      map[string]string
	algo      *macd
	orders    = list.New()
)

func main() {
	fmt.Println("simple napa")
	signals()
	logging()
	auth = readMap("../../private.txt")
	initOrders()
	p, granularity, granularityInt, emaShort, emaLong := initSettings()
	product = p
	interval := time.Second * time.Duration(granularityInt)
	offset := -interval * time.Duration(emaLong)
	sleeping := time.Second * time.Duration(2)
	candleTime := int64(0)
	regulate := true
	for {
		if interrupt {
			break
		}
		end := time.Now().UTC()
		start := end.Add(offset)
		fmt.Println("*", product, "|", start.Format(time.Stamp), "|", end.Format(time.Stamp), "*")
		candles, err := candles(product, start.Format(time.RFC3339), end.Format(time.RFC3339), granularity)
		if err != nil {
			logger(err.Error())
			time.Sleep(time.Second)
			continue
		}
		size := len(candles)
		var wait time.Duration
		if size > 0 && candles[size-1].time > candleTime {
			algo = newMacd(emaShort, emaLong, candles[0].closing)
			candleTime = candles[size-1].time
			for i := 1; i < size; i++ {
				algo.update(candles[i].closing)
			}
			fmt.Println("*", product, "| MACD", algo.current, "| SIGNAL", algo.signal, "*")
			process()
			if regulate {
				wait = interval - time.Now().Sub(time.Unix(candles[size-1].time, 0))
				if wait < 0 {
					wait = interval
				}
				regulate = false
			} else {
				wait = interval
			}
		} else {
			wait = time.Second * time.Duration(6)
		}
		fmt.Println("* sleeping", wait, "*")
		for wait > 0 {
			if interrupt {
				break
			}
			<-time.NewTimer(sleeping).C
			wait -= sleeping
		}
	}
}

func initOrders() {
	o := readList(orderFile)
	for i := 0; i < len(o); i++ {
		fmt.Println("fetching order:", o[i])
		order, status, err := readOrder(auth, o[i])
		if order == nil && err == nil {
			err = errors.New("order is null | status code " + strconv.FormatInt(int64(status), 10))
		}
		if err != nil {
			logger(err.Error())
			panic(err)
		}
		fmt.Println(order)
		orders.PushBack(order)
	}
}

func initSettings() (string, string, int64, int64, int64) {
	s := readMap("settings.txt")
	granularity := s["granularity"]
	granularityInt, err := strconv.ParseInt(granularity, 10, 64)
	if err != nil {
		logger(err.Error())
		panic(err)
	}
	emaShort, err := strconv.ParseInt(s["ema-short"], 10, 64)
	if err != nil {
		logger(err.Error())
		panic(err)
	}
	emaLong, err := strconv.ParseInt(s["ema-long"], 10, 64)
	if err != nil {
		logger(err.Error())
		panic(err)
	}
	return s["product"], granularity, granularityInt, emaShort, emaLong
}

func signals() {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go (func() {
		<-s
		interrupt = true
		fmt.Println("\nsignal interrupt")
	})()
}

func logging() {
	path := "log.txt"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger("failed to open log file:", path)
		os.Exit(0)
	}
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(file)
}

func logger(s ...string) {
	log.Println(s)
	fmt.Println(s)
}
