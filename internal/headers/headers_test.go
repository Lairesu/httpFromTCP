package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFoo:    bar   \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, _ := headers.Get("HOST")
	assert.Equal(t, "localhost:42069", value)
	value, _ = headers.Get("Foo")
	assert.Equal(t, "bar", value)
	assert.Equal(t, 41, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Character
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple values
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:42068\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	value, _ = headers.Get("HOST")
	assert.Equal(t, "localhost:42069,localhost:42068", value)
	assert.Equal(t, 48, n)
	assert.True(t, done)

	// Test: Multiple different headers
	headers = NewHeaders()
	data = []byte("Host: localhost\r\nUser-Agent: Go\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)
	value, _ = headers.Get("host")
	assert.Equal(t, "localhost", value)
	value, _ = headers.Get("user-agent")
	assert.Equal(t, "Go", value)

}
