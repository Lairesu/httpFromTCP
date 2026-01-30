package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		// current line buffer
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(err)
			}

			// split 8-byte chunk
			parts := strings.Split(string(data[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				// fmt.Printf("read: %s\n",o str+parts[i]) // print complete line
				out <- str + parts[i]
				str = "" // reset buffer
			}
			str += parts[len(parts)-1] // last part , carried forward
		}
		if str != "" {
			// fmt.Printf("read: %s\n", str)
			out <- str
		}
	}()
	return out
}

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
		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}
