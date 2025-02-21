package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func InitServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		s.msgch <- Message{
			from: conn.RemoteAddr().String(),
			payload: buff[:n],
		}
		conn.Write([]byte("thank you for your message!\n"))
	}
}

func (s *Server) acceptLoop() {
	for {
		con, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		fmt.Println("new connection to the server:", con.RemoteAddr())
		go s.readLoop(con)
	}
}

func main() {
	server := InitServer(":3000")
	go func() {
		for msg := range server.msgch {
			fmt.Printf("receive message from connection(%s):%s \n", msg.from ,string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
