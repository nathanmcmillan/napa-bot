package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"io"
	"strconv"
	"time"
	"sync"
	"strings"
)

const (
	get     = "GET"
	post    = "POST"
	website = "https://api.gdax.com"
)

var (
	limit = &sync.Mutex{}
	rate = time.Millisecond * time.Duration(500)
)

func request(method, url string, body string) (*http.Client, *http.Request, error) {
	var post io.Reader
	if body != "" {
		post = strings.NewReader(body)
	}
	request, err := http.NewRequest(method, url, post)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "napa")
	return &http.Client{}, request, nil
}

func publicRequest(method, path string) ([]byte, int, error) {
	client, request, err := request(method, website+path, "")
	if err != nil {
		return nil, 0, err
	}
	limit.Lock()
	response, err := client.Do(request)
	time.Sleep(rate)
	limit.Unlock()
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	return read, response.StatusCode, err
}

func privateRequest(auth map[string]string, method, path, body string) ([]byte, int, error) {
	client, request, err := request(method, website+path, body)
	if err != nil {
		return nil, 0, err
	}
	limit.Lock()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	what := timestamp + method + path + body
	base64key, err := base64.StdEncoding.DecodeString(auth["secret"])
	if err != nil {
		return nil, 0, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(what))
	if err != nil {
		return nil, 0, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	request.Header.Add("CB-ACCESS-KEY", auth["key"])
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", auth["phrase"])
	response, err := client.Do(request)
	time.Sleep(rate)
	limit.Unlock()
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	return read, response.StatusCode, err
}
