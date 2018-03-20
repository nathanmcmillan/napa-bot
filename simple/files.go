package main

import (
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	retries = 10
)

func readList(path string) []string {
	try := 0
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return decodeList(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file", err.Error())
		if interrupt {
			return nil
		}
		if try == retries {
			panic(err)
		}
		time.Sleep(time.Second)
		try++
	}
}

func readMap(path string) map[string]string {
	try := 0
	for {
		contents, err := ioutil.ReadFile(path)
		if err == nil {
			return decodeHashmap(contents)
		}
		if os.IsNotExist(err) {
			logger("file not found:", path)
			os.Exit(0)
		}
		logger("failed to open file", err.Error())
		if interrupt {
			return nil
		}
		if try == retries {
			panic(err)
		}
		time.Sleep(time.Second)
		try++
	}
}

func writeBytes(path string, b []byte) {
	try := 0
	for {
		if interrupt {
			return
		}
		err := ioutil.WriteFile(path, b, 0644)
		if err == nil {
			return
		}
		logger("failed to write to file", err.Error())
		if try == retries {
			panic(err)
		}
		time.Sleep(time.Second)
		try++
	}
}

func copyFile(src, dest string) {
	try := 0
	for {
		if interrupt {
			return
		}
		try++

		in, err := os.Open(src)
		if err != nil {
			logger("failed to open file", err.Error())
			if try == retries {
				panic(err)
			}
			time.Sleep(time.Second)
			continue
		}

		out, err := os.Create(dest)
		if err != nil {
			in.Close()
			logger("failed to create file", err.Error())
			if try == retries {
				panic(err)
			}
			time.Sleep(time.Second)
			continue
		}

		_, err = io.Copy(out, in)
		if err != nil {
			in.Close()
			out.Close()
			logger("failed to copy file", err.Error())
			if try == retries {
				panic(err)
			}
			time.Sleep(time.Second)
			continue
		}

		in.Close()
		out.Close()
		return
	}
}

func renameFile(src, dest string) {
	try := 0
	for {
		if interrupt {
			return
		}
		err := os.Rename(src, dest)
		if err == nil {
			return
		}
		logger("failed to rename file", err.Error())
		if try == retries {
			panic(err)
		}
		time.Sleep(time.Second)
		try++
	}
}
