package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type Response struct {
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.Write(b)
	return err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case StatusInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		return fmt.Errorf("Great, you found new Status. Unrecognized error code")
	}
}
