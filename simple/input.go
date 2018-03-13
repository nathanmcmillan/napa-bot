package main

import (
    "os"
    "fmt"   
)

func input() {
	if os.Args[1] == "install" {
        fmt.Println("installing")
	}
}