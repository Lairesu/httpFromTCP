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

// creating our own writer
type writerState int

const (
	stateStatusLine writerState = iota
	stateHeaders
	stateBody
	stateDone
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer, state: stateStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateStatusLine {
		return fmt.Errorf("cannot write status line in current state")
	}
	switch statusCode {
	case StatusOK:
		_, err := w.writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		w.state = stateHeaders
		return err
	case StatusBadRequest:
		_, err := w.writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		w.state = stateHeaders
		return err
	case StatusInternalServerError:
		_, err := w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		w.state = stateHeaders
		return err
	default:
		return fmt.Errorf("Great, you found new Status. Unrecognized error code")
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateHeaders {
		return fmt.Errorf("cannot write headers in current state")
	}
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Append(b, "\r\n")
	_, err := w.writer.Write(b)
	w.state = stateBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot write body in current state")
	}
	n, err := w.writer.Write(p)
	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot write chunked body in current state")
	}

	// writing chunk size in hex
	n := len(p)
	_, err := w.writer.Write([]byte(fmt.Sprintf("%x\r\n", n)))
	if err != nil {
		return 0, err
	}

	// write the chunk itself
	_, err = w.writer.Write(p)
	if err != nil {
		return 0, nil
	}

	// CRLF after chunk
	_, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}

	// return length of data
	return n, nil
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("cannot finish chunked body in current state")
	}

	// write final zero-length chunk
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return n, err
	}

	// mark state done
	w.state = stateDone
	return n, nil
}
