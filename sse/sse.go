package sse

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func HandleSSE(ctx context.Context, e *Events, w http.ResponseWriter, l *log.Entry, lei int) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	ctx2, cancel := context.WithCancel(context.TODO())
	defer cancel()
	evts := e.Subscribe(ctx2, lei)
	h := w.Header()
	// https://html.spec.whatwg.org/multipage/server-sent-events.html
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	flusher.Flush()
	l.Info("Starting SSE")
	var evt *Event
	for {
		select {
		case evt = <-evts:
			evt.Id = fmt.Sprintf("%d", lei)
			evt.Write(w)
			flusher.Flush()
			lei++
			if evt.Ending {
				break
			}
		case <-ctx.Done():
			break
		}
	}
}
