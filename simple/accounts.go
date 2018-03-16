package main

import (
	"encoding/json"
	"errors"
)

type account struct {
	id        string
	currency  string
	balance   *currency
	available *currency
	hold      *currency
	profileID string
}

func readAccounts(auth map[string]string) (map[string]*account, error) {
	status, body, err := privateRequest(auth, get, "/accounts", "")
	if err != nil {
		return nil, err
	}
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		var message interface{}
		errB := json.Unmarshal(body, &message)
		if errB != nil {
			return nil, errors.New(err.Error() + " -> " + errB.Error())
		}
		values, ok := message.(map[string]interface{})
		if !ok {
			return nil, errors.New(err.Error() + " -> parse error message")
		}
		str, _ := values["message"].(string)
		return nil, errors.New(err.Error() + " -> " + str)
	}
	accounts := make(map[string]*account, 0)
	for i := 0; i < len(decode); i++ {
		values, ok := decode[i].(map[string]interface{})
		if !ok {
			return nil, errors.New("parse error accounts")
		}
		a := &account{}
		a.id, _ = values["id"].(string)
		a.currency, _ = values["currency"].(string)
		temp, _ := values["balance"].(string)
		a.balance = newCurrency(temp)
		temp, _ = values["available"].(string)
		a.available = newCurrency(temp)
		temp, _ = values["hold"].(string)
		a.hold = newCurrency(temp)
		a.profileID, _ = values["profile_id"].(string)
		accounts[a.currency] = a
	}
	return accounts, nil
}
