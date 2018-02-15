package gdax

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ListenTicker() {
	fmt.Println("listening")
}

func GetCurrencies() {
	client := &http.Client{}
	request, e := http.NewRequest("GET", "https://api.gdax.com/currencies", nil)
	if e != nil {
		panic(e)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "napa")

	response, e := client.Do(request)
	if e != nil {
		panic(e)
	}
	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(body))
}
