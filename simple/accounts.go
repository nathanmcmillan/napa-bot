package main

import (
	"encoding/json"
	"errors"
)

type account struct {
	id        string
	currency  string
	balance   float64
	available float64
	hold      float64
	profileID string
}

func accounts(a *auth) ([]*account, error) {
	body, err := privateRequest(a, get, "/accounts", "")
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
	accounts := make([]*account, 0)
	for i := 0; i < len(decode); i++ {
		values, ok := decode[i].(map[string]interface{})
		if !ok {
			return nil, errors.New("parse error accounts")
		}
		account := &account{}
		account.id, _ = values["id"].(string)
		account.currency, _ = values["currency"].(string)
		account.balance, _ = values["balance"].(float64)
		account.available, _ = values["available"].(float64)
		account.hold, _ = values["hold"].(float64)
		account.profileID, _ = values["profile_id"].(string)
		accounts = append(accounts, account)
	}
	return accounts, nil
}
