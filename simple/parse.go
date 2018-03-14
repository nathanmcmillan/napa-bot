package main

import (
	"strings"
)

func hashmap(src []byte) map[string]string {
	size := len(src)
	data := make(map[string]string)
	isKey := true
	var key strings.Builder
	var value strings.Builder
	for i := 0; i < size; i++ {
		c := src[i]
		if c == ' ' || c == '\t' {
			isKey = false
			continue
		}
		if c == '\n' || c == '\r' {
			if key.Len() > 0 && value.Len() > 0 {
				data[key.String()] = value.String()
			}
			key.Reset()
			value.Reset()
			isKey = true
			continue
		}
		if isKey {
			key.WriteByte(c)
		} else {
			value.WriteByte(c)
		}
	}
	return data
}

func list(src []byte) []string {
	size := len(src)
	data := make([]string, 0)
	var buffer strings.Builder
	for i := 0; i < size; i++ {
		c := src[i]
		if c == ' ' || c == '\t' {
			continue
		}
		if c == '\n' || c == '\r' {
			if buffer.Len() > 0 {
				data = append(data, buffer.String())
			}
			buffer.Reset()
			continue
		}
		buffer.WriteByte(c)
	}
	return data
}
