package request

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type parserState string

const (
	StateInit    parserState = "initialized"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       parserState
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_MALFORMED_REQUEST = fmt.Errorf("malformed request-line")
var UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var END_OF_LINE = []byte("\r\n")

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	i := bytes.Index(s, END_OF_LINE)

	if i == -1 {
		return nil, 0, nil
	}

	// get the  start line for parsing
	startLine := s[:i]
	RestOfMsg := i + len(END_OF_LINE)

	components := bytes.Split(startLine, []byte(" "))
	if len(components) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST
	}

	HTTPComponents := bytes.Split(components[2], []byte("/"))
	if len(HTTPComponents) != 2 || string(HTTPComponents[0]) != "HTTP" || string(HTTPComponents[1]) != "1.1" {
		return nil, 0, UNSUPPORTED_HTTP_VERSION
	}

	method := components[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, RestOfMsg, ERROR_MALFORMED_REQUEST
		}
	}
	rl := &RequestLine{
		Method:        string(components[0]),
		RequestTarget: string(components[1]),
		HttpVersion:   string(HTTPComponents[1]),
	}

	return rl, RestOfMsg, nil
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			read += n

			if done {
				length := headers.GetInt(r.Headers, "content-length", 0)
				if length > 0 {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}

		case StateBody:
			length := headers.GetInt(r.Headers, "content-length", 0)
			if length == 0 {
				r.state = StateDone
				break
			}
			//
			remainingData := min(length-len(r.Body), len(currentData))
			r.Body = append(r.Body, currentData[:remainingData]...)
			read += remainingData

			if len(r.Body) == length {
				r.state = StateDone
			}

		case StateDone:
			break outer
		default:
			panic("You are failure")
		}
	}
	return read, nil

}

func (r *Request) done() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    []byte{},
	}

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])

		if err != nil {
			return nil, err
		}
		bufLen += n

		// // debugging
		// slog.Info("Read from reader",
		// 	"n", n,
		// 	"bufLen", bufLen,
		// )

		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		// // debugging
		// if readN > bufLen {
		// 	slog.Info("Parse returned more than buffer length",
		// 		"readN", readN,
		// 		"bufLen", bufLen,
		// 	)
		// 	return nil, fmt.Errorf("parse returned readN > bufLen: %d > %d", readN, bufLen)
		// }

		if err == io.EOF && !request.done() {
			return nil, fmt.Errorf("unexpected EOF: request incomplete")
		}

		if readN == 0 && n == 0 {
			// Prevent infinite loop: no progress made
			return nil, fmt.Errorf("no progress parsing request, need more data")
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil

}
