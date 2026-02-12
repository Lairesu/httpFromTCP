package headers

import (
	"bytes"
	"fmt"
	"strings"
)

// tokens
func isToken(str []byte) bool {
	for _, ch := range str {
		found := false
		if ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9' {
			found = true
		}

		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}

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

	// at least 1 character
	if len(name) == 0 {
		return "", "", fmt.Errorf("Malformed Header field-name (token should have at least one character)")
	}

	return strings.ToLower(string(name)), string(value), nil
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

		name, value, err := parseHeader(data[read : read+i])
		if err != nil {
			return 0, false, err
		}

		// checking if token contains the must things
		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("Malformed Header field-name")
		}
		read += i + len(rn)

		// field-name are case insensitivity
		h[name] = value
	}

	return read, done, nil
}
