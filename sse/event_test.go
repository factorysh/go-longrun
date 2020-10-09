package sse

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleEvent(t *testing.T) {
	buffer := &bytes.Buffer{}
	evt := &Event{
		Event: "beuha",
		Data:  "aussi",
	}
	err := evt.Write(buffer)
	assert.NoError(t, err)
	for _, expected := range []string{"event: beuha\n", "data: aussi\n", "\n"} {
		line, err := buffer.ReadString('\n')
		assert.NoError(t, err)
		assert.Equal(t, expected, line)
	}
}

func TestFullEvent(t *testing.T) {
	buffer := &bytes.Buffer{}
	evt := &Event{
		Event: "beuha",
		Data:  "aussi",
		Id:    "42",
		Retry: time.Second / 10,
	}
	err := evt.Write(buffer)
	assert.NoError(t, err)
	datas := make(map[string]string)
	for {
		line, err := buffer.ReadString('\n')
		assert.NoError(t, err)
		if line == "\n" {
			break
		}
		blobs := strings.Split(line, ": ")
		datas[blobs[0]] = strings.Trim(blobs[1], " \n\t")
	}
	assert.Len(t, datas, 4)
	assert.Equal(t, map[string]string{
		"event": "beuha",
		"data":  "aussi",
		"id":    "42",
		"retry": "100",
	}, datas)
}
