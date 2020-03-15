package sse

import (
	"context"
	"sync"
)

type Events struct {
	events    []*Event
	bid       int64
	broadcast map[int64]*subscriber
	block     sync.RWMutex
	lock      sync.RWMutex
}

type subscriber struct {
	events chan *Event
	ctx    context.Context
}

func NewEvents() *Events {
	return &Events{
		events:    make([]*Event, 0),
		broadcast: make(map[int64]*subscriber),
	}
}

func (e *Events) nextBid() int64 {
	e.block.Lock()
	defer e.block.Unlock()
	e.bid++
	return e.bid
}

func (e *Events) Since(since int) []*Event {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.events[since:]
}

func (e *Events) Subscribe(ctx context.Context, since int) chan *Event {
	s := &subscriber{
		events: make(chan *Event, 100),
		ctx:    ctx,
	}
	bid := e.nextBid()
	e.block.Lock()
	e.broadcast[bid] = s
	e.block.Unlock()
	e.lock.RLock()
	for i := since; i < len(e.events); i++ {
		s.events <- e.events[i]
	}
	e.lock.RUnlock()
	go func() {
		select {
		case <-ctx.Done():
		}
		e.block.Lock()
		delete(e.broadcast, bid)
		e.block.Unlock()
	}()
	return s.events
}

func (e *Events) Append(evt *Event) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.events = append(e.events, evt)
	e.block.RLock()
	defer e.block.RUnlock()
	for _, s := range e.broadcast {
		s.events <- evt
	}
}

func (e *Events) Close() {
	e.block.Lock()
	for k, s := range e.broadcast {
		close(s.events)
		delete(e.broadcast, k)
	}
	e.block.Unlock()
}

func (e *Events) Size() int {
	e.lock.RLock()
	defer e.lock.RUnlock()
	return len(e.events)
}
