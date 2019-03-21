package run

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type State string

const (
	QUEUED   = State("queued")
	RUNNING  = State("running")
	CANCELED = State("canceled")
	ERROR    = State("error")
	SUCCESS  = State("success")
)

type Run struct {
	events []*Event
	id     uuid.UUID
}

type Event struct {
	State State       `json:"state"`
	Value interface{} `json:"value,omitempty"`
}

type Runs struct {
	sync sync.Mutex
	run  map[uuid.UUID]*Run
}

func New() *Runs {
	return &Runs{
		run: make(map[uuid.UUID]*Run),
	}
}

func (rr *Runs) New() *Run {
	rr.sync.Lock()
	defer rr.sync.Unlock()
	r := &Run{
		id:     uuid.New(),
		events: []*Event{&Event{QUEUED, nil}},
	}
	rr.run[r.id] = r
	return r
}

func (rr *Runs) Get(id uuid.UUID, since int) ([]*Event, error) {
	r, ok := rr.run[id]
	if !ok {
		return nil, errors.New("Unknown run")
	}
	return r.events[since:], nil
}

func (r *Run) Run(value interface{}) {
	r.events = append(r.events, &Event{RUNNING, value})
}

func (r *Run) Cancel() {
	r.events = append(r.events, &Event{CANCELED, nil})
}

func (r *Run) Error(err error) {
	r.events = append(r.events, &Event{ERROR, err.Error()})
}

func (r *Run) Success(value interface{}) {
	r.events = append(r.events, &Event{SUCCESS, value})
}

func (r *Run) Id() uuid.UUID {
	return r.id
}
