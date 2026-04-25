package main

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	quitCh     chan struct{}
	ln         net.Listener
}

func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddr: listenAddress,
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.ConnLoop()

	<-s.quitCh

	return nil
}

func (s *Server) ConnLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("connection accept error : ", err)
			continue
		}
		go ReadLoop(conn)
	}
}

// main handler for connections
func ReadLoop(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error while reading from a connection: ", err)
			continue
		}

		fmt.Println(string(buffer[:n]))
	}

}

func main() {
	server := NewServer(":3000")
	fmt.Println("localhost:3000")
	if err := server.Listen(); err != nil {
		panic(err)
	}
}
