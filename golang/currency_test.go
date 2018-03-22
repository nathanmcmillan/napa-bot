package main

import (
	"testing"
)

// TestCurrency test
func TestCurrency(t *testing.T) {
	var c *currency

	precision := 8
	str := "5957.11914015"
	c = newCurrency(str)
	expectString(t, c.str(precision), str)
}

func expectString(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Error("expected", expected, "got", actual)
	}
}
