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
)

const (
	get     = "GET"
	post    = "POST"
	website = "https://api.gdax.com"
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

func publicRequest(method, path string) ([]byte, error) {
	client, request, err := request(method, website+path, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func privateRequest(a *auth, method, path, body string) ([]byte, error) {
	var data io.Reader
	if body != "" {
		message, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		data = bytes.NewReader(message)
	}
	client, request, err := request(method, website+path, data)
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	what := timestamp + method + path + body
	base64key, err := base64.StdEncoding.DecodeString(a.secret)
	if err != nil {
		return nil, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(what))
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))
	request.Header.Add("CB-ACCESS-KEY", a.key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", a.phrase)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
