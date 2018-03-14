package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
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
	o := orders()
	a := authentication()
	acc, err := accounts(a)
	if err != nil {
		logger(err.Error())
	}
	fmt.Println(acc)
	product := "BTC-USD"
	end := time.Now().UTC()
	start := time.Now().UTC().Add(beginning)
	fmt.Println("pollin", product, "from", start, "to", end)
	c, err := candles(a, product, start.Format(time.RFC3339), end.Format(time.RFC3339), granularity)
	for {
		fmt.Println("sleeping")
		wait := time.NewTimer(time.Second)
		<-wait.C
		if interrupt {
			break
		}
	}
}

func orders() []string {
	path := "orders.txt"
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

func authentication() *auth {
	path := "../../private.txt"
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			data := hashmap(contents)
			return &auth{data["key"], data["secret"], data["phrase"]}
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
