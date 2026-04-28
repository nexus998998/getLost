package server

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	quitCh     chan struct{}
	ln         net.Listener
}

var (
	queueingPlayers = make(chan Player)
)

type Player struct {
	conn net.Conn
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
		go s.matchPlayersLoop()
	}
}

// main handler for connections
func ReadLoop(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == net.ErrClosed {
				fmt.Println("connection closed")
				break
			}
			fmt.Println("error while reading from a connection: ", err)
			continue
		}
		msg := buffer[:n]
		fmt.Println(string(msg))
		// map of funcs to messages
		// possible improvement : use json to send the action : and other headers
		ServeRequest(msg)
	}
}

func (s *Server) matchPlayersLoop() {
	for {
		player1, player2 := <-queueingPlayers, <-queueingPlayers
		_, err := player1.conn.Write([]byte("now you got queued with someone !"))
		if err != nil {
			fmt.Println("error while sending ", err)
			continue
		}
		_, err = player2.conn.Write([]byte("now you got queued with someone !"))
		if err != nil {
			fmt.Println("error while sending ", err)
			continue
		}
		// gameLoop lives here
	}
}
