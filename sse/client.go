package sse

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type Event struct {
	Data  string
	Id    string
	Event string
	Retry string
}

func (e *Event) JSON(v interface{}) error {
	return json.Unmarshal([]byte(e.Data), v)
}

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
		evt.Retry = value
	case "event":
		evt.Event = value
	case "data":
		evt.Data = evt.Data + value
	}
}
