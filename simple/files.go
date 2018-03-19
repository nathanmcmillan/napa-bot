package main

import (
	"io/ioutil"
	"os"
	"time"
)

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

func writeList(path string, list []byte) {
	fmt.Println("not working...")
	for {
		err := ioutil.WriteFile(path, list, 0644)
		if err == nil {
			return
		}
		logger("failed to write to file", path)
		time.Sleep(time.Second)
	}
}
