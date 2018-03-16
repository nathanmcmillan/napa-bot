package main

import (
	"math/big"
)

type currency struct {
	num *big.Rat
}

func newCurrency(num string) *currency {
	c := &currency{}
	c.num, _ = new(big.Rat).SetString(num)
	return c
}

func (c *currency) add(o *currency) *currency {
	n := &currency{}
	n.num.Add(c.num, o.num)
	return n
}

func (c *currency) mul(o *currency) *currency {
	n := &currency{}
	n.num.Mul(c.num, o.num)
	return n
}

func (c *currency) cmp(o *currency) int {
	return c.num.Cmp(o.num)
}
