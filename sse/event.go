package sse

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Event struct {
	dataExists bool
	Data       string
	Id         string
	Event      string
	Retry      time.Duration
	Ending     bool
}

func (e *Event) Write(w io.Writer) error {
	if e.Id != "" {
		_, err := fmt.Fprintf(w, "id: %s\n", e.Id)
		if err != nil {
			return err
		}
	}
	if e.Event != "" {
		_, err := fmt.Fprintf(w, "event: %s\n", e.Event)
		if err != nil {
			return err
		}
	}
	if e.Retry != 0 {
		_, err := fmt.Fprintf(w, "retry: %d\n", e.Retry)
		if err != nil {
			return err
		}
	}
	for _, data := range strings.Split(e.Data, "\n") {
		_, err := fmt.Fprintf(w, "data: %s\n", data)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\n"))
	return err
}
