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
	}
	id, err := uuid.Parse(slugs[1])
	if err != nil {
		l.WithError(err).Error()
		w.WriteHeader(400)
		return
	}
	evts, err := s.runs.Subscribe(id, lei)
	if err != nil {
		// FIXME maybe some 404 if id doesn't exist
		l.WithError(err).Error()
		w.WriteHeader(400)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	l.Info("Starting SSE")
	cpt := 0
	for {
		evt := <-evts
		j, err := json.Marshal(evt)
		if err != nil {
			l.WithError(err).Error()
			return
		}
		w.Write([]byte(fmt.Sprintf("id: %d\n", cpt)))
		w.Write([]byte("data: "))
		w.Write(j)
		w.Write([]byte("\n\n"))
		if evt.Ended() {
			return
		}
		cpt++
	}
}
