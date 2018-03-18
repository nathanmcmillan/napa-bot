package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	orderFile     = "orders.txt"
	orderSwapFile = "orders_swap.txt"
)

var (
	interrupt = false
	product string
	auth map[string]string
	algo *macd
	orders = make(book, 0)
)

func main() {
	fmt.Println("simple napa")
	signals()
	logging()
	var err error
	o := readList(orderFile)
	s := readMap("settings.txt")
	auth = readMap("../../private.txt")
	fmt.Println(o)
	fmt.Println(s)
	for i := 0; i < len(o); i++ {
		fmt.Println("reading order", o[i])
		order, err := readOrder(auth, o[i])
		if err != nil {
			logger(err.Error())
			return
		}
		fmt.Println(order.executedValue)
		orders.push(order)
		time.Sleep(time.Second)
	}
	product = s["product"]
	granularity := s["granularity"]
	ig, err := strconv.ParseInt(granularity, 10, 64)
	if err != nil {
		logger(err.Error())
		return
	}
	es, err := strconv.ParseInt(s["ema-short"], 10, 64)
	if err != nil {
		logger(err.Error())
		return
	}
	el, err := strconv.ParseInt(s["ema-long"], 10, 64)
	if err != nil {
		logger(err.Error())
		return
	}
	interval := time.Second * time.Duration(ig)
	offset := -interval * time.Duration(el)
	sleeping := time.Second * time.Duration(2)
	candleTime := int64(0)
	regulate := true
	for {
		if interrupt {
			break
		}
		end := time.Now().UTC()
		var start time.Time
		if candleTime > 0 {
			start = time.Unix(candleTime, 0)
		} else {
			start = end.Add(offset)
		}
		fmt.Println("*", product, "|", start.Format(time.Stamp), "|", end.Format(time.Stamp), "*")
		fmt.Println("*", product, "|", start.Format(time.RFC3339), "|", end.Format(time.RFC3339), "*")
		// ???
		c, err := candles(product, start.Format(time.RFC3339), end.Format(time.RFC3339), granularity)
		if err != nil {
			logger(err.Error())
			time.Sleep(time.Second)
			continue
		}
		var wait time.Duration
		if len(c) > 0 && c[len(c)-1].time > candleTime {
			var i int
			if algo == nil {
				algo = newMacd(es, el, c[0].closing)
				candleTime = c[0].time
				i = 1
			} else {
				i = 0
			}
			for i < len(c) {
				ctime := c[i].time
				if candleTime < ctime {
					algo.update(c[i].closing)
					candleTime = ctime
				}
				i++
			}
			fmt.Println("*", product, "| MACD", algo.current, "| SIGNAL", algo.signal, "*")
			process()
			if regulate {
				wait = interval - time.Now().Sub(time.Unix(c[len(c)-1].time, 0))
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
		fmt.Println("sleeping", wait)
		for wait > 0 {
			if interrupt {
				break
			}
			<-time.NewTimer(sleeping).C
			wait -= sleeping
		}
	}
}

func signals() {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go (func() {
		<-s
		interrupt = true
		fmt.Println("signal interrupt")
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
