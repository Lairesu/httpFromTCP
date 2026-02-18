package server

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
	}

	go server.Listen()
	return server, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}
func (s *Server) handle(conn net.Conn) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!")
	conn.Write(out)
	conn.Close()
}
