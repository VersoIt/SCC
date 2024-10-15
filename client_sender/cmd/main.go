package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello from client_sender")
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("client_sender")
	}
}
