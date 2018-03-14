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

func main() {
	fmt.Println("simple napa")
	signals := signals()
	logging()
	o := orders()
	fmt.Print("*orders*\n", o, "\n")
	a := authentication()
	acc, err := accounts(a)
	if err != nil {
		logger(err.Error())
	}
	fmt.Println(acc)
loop:
	for {
		wait := time.NewTimer(time.Second * time.Duration(5))
		select {
		case <-wait.C:
			fmt.Println("sleeping")
			continue
		case <-signals:
			wait.Stop()
			fmt.Println("signal interrupt")
			break loop
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

func signals() chan os.Signal {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	return s
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
