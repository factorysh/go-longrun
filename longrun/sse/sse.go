package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/factorysh/go-longrun/run"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type SSE struct {
	runs *run.Runs
}

func New(r *run.Runs) *SSE {
	return &SSE{
		runs: r,
	}
}

func (s *SSE) ServeRun(w http.ResponseWriter, l *log.Entry, _run *run.Run, lei int) {
	evts := _run.Subscribe(lei)
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	l.Info("Starting SSE")
	for {
		evt := <-evts
		j, err := json.Marshal(evt)
		if err != nil {
			l.WithError(err).Error()
			return
		}
		fmt.Fprintf(w, "id: %d\n", lei)
		fmt.Fprintf(w, "data: %s\n\n", j)
		if evt.Ended() {
			return
		}
		lei++
	}

}

func (s *SSE) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := log.WithField("url", r.URL.String())
	slugs := strings.Split(r.URL.Path, "/")
	if len(slugs) < 2 {
		w.WriteHeader(400)
		return
	}
	leiRaw := r.Header.Get("Last-Event-ID")
	lei := 0
	var err error
	if leiRaw != "" {
		lei, err = strconv.Atoi(leiRaw)
		if err != nil {
			l.WithError(err).Error()
			w.WriteHeader(400)
			return
		}
		lei++
	}
	id, err := uuid.Parse(slugs[1])
	if err != nil {
		l.WithError(err).Error()
		w.WriteHeader(400)
		return
	}
	_run, ok := s.runs.GetRun(id)
	if !ok {
		w.WriteHeader(404)
		return
	}
	s.ServeRun(w, l, _run, lei)
}
