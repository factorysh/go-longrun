package sse

import (
	"bufio"
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSSEReader(t *testing.T) {
	buff := &bytes.Buffer{}
	buff.WriteString("data: plop\n\n")
	client := NewSSEReader(bufio.NewReader(buff))
	evt, err := client.Read()
	assert.NoError(t, err)
	assert.Equal(t, "plop", evt.Data)
	evt, err = client.Read()
	assert.Nil(t, evt)
	assert.Equal(t, io.EOF, err)
}
func TestMoreSSEReader(t *testing.T) {
	buff := &bytes.Buffer{}
	buff.WriteString("id: 42\nevent: siesta\nretry: 60\ndata: beuha\ndata: aussi\n\n")
	client := NewSSEReader(bufio.NewReader(buff))
	evt, err := client.Read()
	assert.NoError(t, err)
	assert.Equal(t, "beuha\naussi", evt.Data)
	assert.Equal(t, "siesta", evt.Event)
	assert.Equal(t, "42", evt.Id)
	assert.Equal(t, 60*time.Millisecond, evt.Retry)
	evt, err = client.Read()
	assert.Nil(t, evt)
	assert.Equal(t, io.EOF, err)
}
