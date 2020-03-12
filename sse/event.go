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
}

func (e *Event) Write(w io.Writer) {
	if e.Id != "" {
		fmt.Fprintf(w, "id: %s\n", e.Id)
	}
	if e.Event != "" {
		fmt.Fprintf(w, "event: %s\n", e.Event)
	}
	if e.Retry != 0 {
		fmt.Fprintf(w, "retry: %d\n", e.Retry)
	}
	for _, data := range strings.Split(e.Data, "\n") {
		fmt.Fprintf(w, "data: %s\n", data)
	}
	w.Write([]byte("\n"))
}
