package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"net"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
type Server struct {
	listener net.Listener
	closed   bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		closed:   false,
		handler:  handler,
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
	defer conn.Close()

	headers := response.GetDefaultHeaders(0)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, *headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, r)
	if handlerError != nil {
		// updating the body and Content-Length
		body := []byte(handlerError.Message)
		headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))

		response.WriteStatusLine(conn, handlerError.StatusCode)
		response.WriteHeaders(conn, *headers)
		conn.Write(body)
		return
	}

	// 200 OK PATH
	body := writer.Bytes()
	headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))

	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, *headers)
	conn.Write(body)
}
