package main

import (
	"encoding/json"
	"fmt"
	"getLost/internals/comms"
	"log"
	"net"
)

var (
	serverIP = "localhost:3000"
)

type ServerConn struct {
	conn net.Conn
}

func GetServerConn() (*ServerConn, error) {
	connection, err := net.Dial("tcp", serverIP)
	if err != nil {
		return nil, err
	}

	return &ServerConn{
		conn: connection,
	}, nil

}

func (s *ServerConn) getMessageHello() error {
	r := comms.Request{
		Action: "hello",
	}
	j, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = s.conn.Write(j)
	return err
}
func (s *ServerConn) Start() {
	s.ReadingLoop()
}

func (s *ServerConn) ReadingLoop() {
	if err := s.getMessageHello(); err != nil {
		fmt.Println("error while sending a message : ", err)
	}
	buffer := make([]byte, 2000)
	for {
		n, _ := s.conn.Read(buffer)
		fmt.Println(string(buffer[:n]))
	}
}
func main() {
	server, err := GetServerConn()
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
