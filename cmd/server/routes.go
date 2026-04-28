package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"getLost/internals/comms"
)

var routes = map[string]func(){
	"hello": printHello,
}

// takes any request message and directs it into the corresponding route function
func ServeRequest(msg []byte) error {
	var request comms.Request
	json.Unmarshal(msg, &request)
	f, ok := routes[request.Action]
	if !ok {
		return errors.New("error : ")
	}
	f()
	return nil
}

func printHello() {
	fmt.Println("hello world")
}
