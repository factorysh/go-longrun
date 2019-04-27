package run

import (
	"errors"
	"fmt"
	"sync"
	"time"

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
	events    []*Event
	id        uuid.UUID
	bid       int64
	broadcast map[int64]func()
	block     sync.Mutex
	lock      sync.Mutex
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
		id:        uuid.New(),
		events:    []*Event{&Event{QUEUED, nil}},
		broadcast: make(map[int64]func()),
		lock:      sync.Mutex{},
	}
	rr.run[r.id] = r
	return r
}

func (rr *Runs) Get(id uuid.UUID, since int) ([]*Event, error) {
	r, ok := rr.run[id]
	if !ok {
		return nil, errors.New("Unknown run")
	}
	if len(r.events) <= since {
		wait := make(chan interface{})
		// FIXME forever growing array
		bid := r.nextBid()
		r.broadcast[bid] = func() {
			wait <- new(interface{})
		}
		defer delete(r.broadcast, bid)
		select {
		case <-wait:
			fmt.Println("Wait")
		case <-time.After(10 * time.Second):
			fmt.Println("oups, timeout")
		}

	}
	return r.events[since:], nil
}

func (rr *Runs) Subscribe(id uuid.UUID, since int) (chan *Event, error) {
	r, ok := rr.run[id]
	if !ok {
		return nil, errors.New("Unknown run")
	}
	return r.Subscribe(since), nil
}

func (r *Run) Subscribe(since int) chan *Event {
	c := make(chan *Event)
	go func() {
		cpt := 0
		bid := r.nextBid()
		for {
			if since+cpt >= len(r.events) {
				c := make(chan interface{})
				r.broadcast[bid] = func() {
					c <- new(interface{})
				}
				<-c
				delete(r.broadcast, bid)
			}
			for _, evt := range r.events[since+cpt:] {
				c <- evt
				if evt.Ended() {
					return
				}
				cpt++
			}
		}
	}()
	return c
}

func (r *Run) nextBid() int64 {
	r.block.Lock()
	defer r.block.Unlock()
	r.bid++
	return r.bid
}

func (r *Run) append(evt *Event) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.events = append(r.events, evt)
	for _, f := range r.broadcast {
		go f()
	}
}

func (r *Run) Run(value interface{}) {
	r.append(&Event{RUNNING, value})
}

func (r *Run) Cancel() {
	r.append(&Event{CANCELED, nil})
}

func (r *Run) Error(err error) {
	r.append(&Event{ERROR, err.Error()})
}

func (r *Run) Success(value interface{}) {
	r.append(&Event{SUCCESS, value})
}

func (r *Run) Id() uuid.UUID {
	return r.id
}

func (e *Event) Ended() bool {
	return e.State == CANCELED || e.State == SUCCESS || e.State == ERROR
}
