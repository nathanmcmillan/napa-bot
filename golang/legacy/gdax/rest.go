package gdax

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

const (
	api       = "https://api.gdax.com"
	apiSocket = "wss://ws-feed.gdax.com"
	get       = "GET"
	post      = "POST"
	rateLimit = time.Millisecond * time.Duration(500)
)

type Rest struct {
	mutex *sync.Mutex
	auth  *Authentication
}

func NewRest(a *Authentication) *Rest {
	r := &Rest{}
	r.mutex = &sync.Mutex{}
	r.auth = a
	return r
}

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

func (r *Rest) publicRequest(method, url string) ([]byte, error) {
	client, request, err := request(method, url, nil)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	response, err := client.Do(request)
	time.Sleep(rateLimit)
	r.mutex.Unlock()

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}

func (r *Rest) privateRequest(method, site, path, body string) ([]byte, error) {
	var data io.Reader
	if body != "" {
		message, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		data = bytes.NewReader(message)
	}

	client, request, err := request(method, site+path, data)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	what := timestamp + method + path + body
	base64key, err := base64.StdEncoding.DecodeString(r.auth.Secret)
	if err != nil {
		return nil, err
	}
	hashMessage := hmac.New(sha256.New, base64key)
	_, err = hashMessage.Write([]byte(what))
	if err != nil {
		return nil, err
	}
	signature := base64.StdEncoding.EncodeToString(hashMessage.Sum(nil))

	request.Header.Add("CB-ACCESS-KEY", r.auth.Key)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", r.auth.Passphrase)

	response, err := client.Do(request)
	time.Sleep(rateLimit)
	r.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
