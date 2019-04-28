package run

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Run struct {
	runs      *Runs
	events    []*Event
	id        uuid.UUID
	bid       int64
	broadcast map[int64]func()
	block     sync.Mutex
	lock      sync.Mutex
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
	go func() {
		time.Sleep(r.runs.ttl)
		delete(r.runs.run, r.id)
	}()
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
