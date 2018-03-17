package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	interrupt = false
)

func main() {
	fmt.Println("simple napa")
	signals()
	logging()
	o := readList("orders.txt")
	s := readMap("settings.txt")
	a := readMap("../../private.txt")
	fmt.Println(o)
	fmt.Println(s)
	acc, err := readAccounts(a)
	if err != nil {
		logger(err.Error())
		return
	}
	orders := make(book, 0)
	for i := 0; i < len(o); i++ {
		fmt.Println("reading order", o[i])
		order, err := readOrder(a, o[i])
		if err != nil {
			logger(err.Error())
			return
		}
		fmt.Println(order.executedValue)
		orders.push(order)
		time.Sleep(time.Second)
	}
	product := s["product"]
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
	for {
		if interrupt {
			break
		}
		end := time.Now().UTC()
		start := time.Now().UTC().Add(offset)
		fmt.Println("polling", product, "from", start, "to", end)
		c, err := candles(product, start.Format(time.RFC3339), end.Format(time.RFC3339), granularity)
		if err != nil {
			logger(err.Error())
			time.Sleep(time.Second)
			continue
		}
		var wait time.Duration
		if len(c) > 0 {
			m := newMacd(es, el, c[0].closing)
			for i := 1; i < len(c); i++ {
				m.update(c[i].closing)
			}
			fmt.Println("*", product, "| MACD", m.current, "| SIGNAL", m.signal, "*")
			// start process
			if m.signal == "buy" {
				zero := newCurrency("0.0")
				if acc["USD"].available.cmp(zero) > 0 {
					buy(a, "5.0")
					orders.push(nil)
					// orders = append(orders, order)
					// write to orders.txt
				}
			} else if m.signal == "sell" {
				t, err := tick(product)
				if err != nil {
					logger(err.Error())
					time.Sleep(time.Second)
					continue
				}
				size := len(o)
				for i := 0; i < size; i++ {
					order := orders[i]
					min := profitPrice(order)
					fmt.Println("*", product, "|", min, ">", t.price, "*")
					if min.cmp(t.price) > 0 {
						sell(a, order)
						orders.delete(i)
						// orders remove slice at index of order
						// write to orders.txt
					}
				}
			}
			// end process
			wait = interval - time.Now().Sub(time.Unix(c[len(c)-1].time, 0))
			if wait < 0 {
				wait = interval
			}
		} else {
			wait = time.Second
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

func readList(path string) []string {
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return list(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file")
		time.Sleep(time.Second)
	}
}

func readMap(path string) map[string]string {
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return hashmap(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file")
		time.Sleep(time.Second)
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(file)
}

func logger(s ...string) {
	log.Println(s)
	fmt.Println(s)
}
