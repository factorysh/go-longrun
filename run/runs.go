package run

import (
	"context"
	"sync"
	"time"

	"github.com/factorysh/go-longrun/sse"
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

func (rr *Runs) NewRun(ctx context.Context) *Run {
	rr.sync.Lock()
	defer rr.sync.Unlock()
	r := &Run{
		runs:   rr,
		id:     uuid.New(),
		Events: sse.NewEvents(),
	}
	r.Events.Append(&sse.Event{
		Event: string(QUEUED),
		Data:  "null",
	})
	rr.run[r.id] = r
	return r
}

func (rr *Runs) GetRun(id uuid.UUID) (*Run, bool) {
	r, ok := rr.run[id]
	return r, ok
}
