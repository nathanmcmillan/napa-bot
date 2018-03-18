package main

import "strings"

func process() {

}

func buy(a map[string]string, product, funds string) (*order, error) {
	rawJs := &strings.Builder{}
	beginJs(rawJs)
	pushJs(rawJs, "type", "market")
	pushJs(rawJs, "side", "buy")
	pushJs(rawJs, "product_id", product)
	pushJs(rawJs, "funds", funds)
	endJs(rawJs)
	str := rawJs.String()
	logger(str)
	return postOrder(a, str)
}

func sell(a map[string]string, o *order) (*order, error) {
	rawJs := &strings.Builder{}
	beginJs(rawJs)
	firstJs(rawJs, "type", "market")
	pushJs(rawJs, "side", "sell")
	pushJs(rawJs, "product_id", o.product)
	pushJs(rawJs, "size", o.size.str(precision[o.product]))
	endJs(rawJs)
	str := rawJs.String()
	logger(str)
	return postOrder(a, str)
}
