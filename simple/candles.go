package main

import (
	"encoding/json"
	"errors"
	"sort"
)

type candle struct {
	time    int64
	low     float64
	high    float64
	open    float64
	closing float64
	volume  float64
}

type sortCandles []*candle

func (ls sortCandles) Len() int {
	return len(ls)
}
func (ls sortCandles) Swap(left, right int) {
	ls[left], ls[right] = ls[right], ls[left]
}
func (ls sortCandles) Less(left, right int) bool {
	return ls[left].time < ls[right].time
}

func candles(product, start, end, granularity string) ([]*candle, error) {
	body, err := publicRequest(get, "/products/"+product+"/candles?start="+start+"&end="+end+"&granularity="+granularity)
	if err != nil {
		return nil, err
	}
	var decode []interface{}
	err = json.Unmarshal(body, &decode)
	if err != nil {
		return nil, err
	}
	candles := make([]*candle, 0)
	for i := 0; i < len(decode); i++ {
		values, ok := decode[i].([]interface{})
		if !ok {
			return nil, errors.New("not a list")
		}
		c := &candle{}
		floatTime, _ := values[0].(float64)
		c.time = int64(floatTime)
		c.low, _ = values[1].(float64)
		c.high, _ = values[2].(float64)
		c.open, _ = values[3].(float64)
		c.closing, _ = values[4].(float64)
		c.volume, _ = values[5].(float64)
		candles = append(candles, c)
	}
	sort.Sort(sortCandles(candles))
	return candles, nil
}
