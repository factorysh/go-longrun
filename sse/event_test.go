package sse

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	buffer := &bytes.Buffer{}
	evt := &Event{
		Event: "beuha",
		Data:  "aussi",
	}
	err := evt.Write(buffer)
	assert.NoError(t, err)
	line, err := buffer.ReadString('\n')
	assert.NoError(t, err)
	assert.Equal(t, "event: beuha\n", line)
}
