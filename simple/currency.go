package main

import (
	"math/big"
)

var (
	zero = newCurrency("0.0")
	one = newCurrency("1.0")
	two = newCurrency("2.0")
)

type currency struct {
	num *big.Rat
}

func newCurrency(num string) *currency {
	c := &currency{}
	c.num, _ = new(big.Rat).SetString(num)
	return c
}

func (c *currency) plus(o *currency) *currency {
	n := &currency{}
	n.num = new(big.Rat)
	n.num.Add(c.num, o.num)
	return n
}

func (c *currency) minus(o *currency) *currency {
	n := &currency{}
	n.num = new(big.Rat)
	n.num.Sub(c.num, o.num)
	return n
}

func (c *currency) mul(o *currency) *currency {
	n := &currency{}
	n.num = new(big.Rat)
	n.num.Mul(c.num, o.num)
	return n
}

func (c *currency) div(o *currency) *currency {
	n := &currency{}
	n.num = new(big.Rat)
	if o.num.Cmp(zero.num) != 0 {
		n.num.Quo(c.num, o.num)
	}
	return n
}

func (c *currency) moreThan(o *currency) bool {
	return c.num.Cmp(o.num) > 0
}

func (c *currency) str(precision int) string {
	return c.num.FloatString(precision)
}
