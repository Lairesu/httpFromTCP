package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)

func main() {
	Listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer Listener.Close()

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// calling RequestFromReader
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Request Line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Targe: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
