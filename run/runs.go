package run

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Runs struct {
	sync sync.Mutex
	ttl  time.Duration
	run  map[uuid.UUID]*Run
}

func New(ttl time.Duration) *Runs {
	return &Runs{
		ttl: ttl,
		run: make(map[uuid.UUID]*Run),
	}
}

func (rr *Runs) New() *Run {
	rr.sync.Lock()
	defer rr.sync.Unlock()
	r := &Run{
		runs:      rr,
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
