package headers

import (
	"bytes"
	"fmt"
	"strconv"
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

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	// fmt.Printf("fieldLine: %q\n", string(fieldLine))
	// fmt.Printf("parts count: %d\n", len(parts))
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

	return string(name), string(value), nil
}

var rn = []byte("\r\n")

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func GetInt(h *Headers, name string, defaultValue int) int {
	valueStr, exists := h.Get(name)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func (h *Headers) Get(name string) (string, bool) {
	str, ok := h.headers[strings.ToLower(name)]
	return str, ok
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}

// This is old Set
func (h *Headers) Add(name, value string) {
	name = strings.ToLower(name)
	if existing, ok := h.headers[name]; ok && existing != "" {
		h.headers[name] = existing + "," + value
	} else {
		h.headers[name] = value
	}
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
		h.Add(name, value)
	}

	return read, done, nil
}
