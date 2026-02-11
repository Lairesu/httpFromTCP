package headers

import (
	"bytes"
	"fmt"
)

var rn = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	// whitespace
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Malformed field name")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		i := bytes.Index(data[read:], rn)
		if i == -1 {
			break // need more data
		}

		// Empty line means end of headers
		if i == 0 {
			done = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read:read+i])
		if err != nil {
			return 0, false, err
		}
		read += i + len(rn)
		h[name] = value
	}

	return read, done, nil
}
