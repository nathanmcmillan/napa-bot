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

func writeBytes(path string, b []byte) error {
	return ioutil.WriteFile(path, b, 0644)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		logger("failed to open file", err.Error())
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		logger("failed to create file", err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		logger("failed to copy file", err.Error())
		return err
	}
	
	return nil
}

func renameFile(src, dest string) error {
	err := os.Rename(src, dest)
	if err != nil {
		logger("failed to rename file", err.Error())
		return err
	}
	return nil
}
