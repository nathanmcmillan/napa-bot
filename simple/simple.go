package main

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"time"
	"log"
	"io/ioutil"
)

func main() {
	fmt.Println("simple napa")
	
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	
	f := funds()
	fmt.Println(f)
}

func funds() string {
	for {
		contents, err := ioutil.ReadFile("funds.txt")
		if err == nil {
			return string(contents)
		}
		log.Println("failed to open file")
		time.Sleep(time.Second)
	}
}