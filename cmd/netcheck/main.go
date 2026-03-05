package main

import (
	"fmt"

	"github.com/Kiryue0/go-network-checker/internal/network"
)

func main() {

	fmt.Println("Hello World")
	fmt.Println(network.PingHost("google.com"))

}
