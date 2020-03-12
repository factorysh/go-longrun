package sse

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	dataExists bool
	Data       string
	Id         string
	Event      string
	Retry      time.Duration
}

type SSEReader struct {
	scanner *bufio.Scanner
}

func NewSSEReader(r io.Reader) *SSEReader {
	return &SSEReader{bufio.NewScanner(r)}
}

func (s *SSEReader) Read() (*Event, error) {
	evt := &Event{}
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if line == "" {
			return evt, nil
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		switch len(parts) {
		case 1:
			event(evt, parts[0], "")
		case 2:
			event(evt, parts[0], parts[1][:len(parts[1])])
		}
	}
	return nil, io.EOF
}

// DEPRECATED
func Reader(r io.Reader, visitor func(*Event) error) error {
	scanner := bufio.NewScanner(r)
	evt := &Event{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			err := visitor(evt)
			if err != nil {
				return err
			}
			evt = &Event{}
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		switch len(parts) {
		case 1:
			event(evt, parts[0], "")
		case 2:
			event(evt, parts[0], parts[1][:len(parts[1])])
		}
	}
	return scanner.Err()
}

func event(evt *Event, key, value string) {
	if strings.HasPrefix(value, " ") {
		value = value[1:]
	}
	if strings.HasSuffix(value, "\r") {
		value = value[:len(value)-1]
	}
	switch key {
	case "id":
		evt.Id = value
	case "retry":
		retry, err := strconv.Atoi(value)
		if err == nil {
			// like Mozilla, we doesn't throw an error
			evt.Retry = time.Duration(retry) * time.Millisecond
		}
	case "event":
		evt.Event = value
	case "data":
		if evt.dataExists {
			evt.Data = fmt.Sprintf("%s\n%s", evt.Data, value)
		} else {
			evt.Data = value
			evt.dataExists = true
		}
	}
}
