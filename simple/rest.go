package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"sync"
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

func request(method, url string, body io.Reader) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "napa")
	return client, request, nil
}

func publicRequest(method, path string) ([]byte, error, int) {
	client, request, err := request(method, website+path, nil)
	if err != nil {
		return nil, err, 0
	}
	limit.Lock()
	response, err := client.Do(request)
	time.Sleep(rate)
	limit.Unlock()
	if err != nil {
		return nil, err, 0
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	return read, err, response.StatusCode
}

func privateRequest(auth map[string]string, method, path, body string) ([]byte, error, int) {
	var data io.Reader
	if body != "" {
		message, err := json.Marshal(body)
		if err != nil {
			return nil, err, 0
		}
		data = bytes.NewReader(message)
	}
	client, request, err := request(method, website+path, data)
	if err != nil {
		return nil, err, 0
	}
	limit.Lock()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	what := timestamp + method + path + body
	base64key, err := base64.StdEncoding.DecodeString(auth["secret"])
	if err != nil {
		return nil, err, 0
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(what))
	if err != nil {
		return nil, err, 0
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
		return nil, err, 0
	}
	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	return read, err, response.StatusCode
}
