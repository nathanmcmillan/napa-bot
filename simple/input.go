package main

import (
	"fmt"
	"os"
)

func input() {
	if os.Args[1] == "install" {
		fmt.Println("installing")
	}
}
