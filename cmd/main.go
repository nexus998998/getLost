package main

import (
	"fmt"
	"getLost/cmd/server"
)

func main() {
	server := server.NewServer(":3000")
	fmt.Println("localhost:3000")
	if err := server.Listen(); err != nil {
		panic(err)
	}
}
