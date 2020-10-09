package run

import (
	"encoding/json"

	"github.com/factorysh/go-longrun/sse"
	"github.com/google/uuid"
)

type Run struct {
	runs *Runs

	Events *sse.Events
	id     uuid.UUID
}

func (r *Run) append(evt *Event) {
	j, _ := json.Marshal(evt.Value)
	r.Events.Append(&sse.Event{
		Data:   string(j),
		Event:  string(evt.State),
		Ending: evt.Ended(),
	})
}

func (r *Run) Run(value interface{}) {
	r.append(&Event{
		State: RUNNING,
		Value: value,
	})
}

func (r *Run) Cancel() {
	r.append(&Event{
		State: CANCELED,
		Value: nil,
	})
	r.Events.Close()
}

func (r *Run) Error(err error) {
	r.append(&Event{
		State: ERROR,
		Value: err.Error(),
	})
	r.Events.Close()
}

func (r *Run) Success(value interface{}) {
	r.append(&Event{
		State: SUCCESS,
		Value: value,
	})
	r.Events.Close()
}

func (r *Run) Id() uuid.UUID {
	return r.id
}

func (r *Run) Since(since int) ([]*Event, error) {
	evts := r.Events.Since(since)
	resp := make([]*Event, len(evts))
	for i, e := range evts {
		var value interface{}
		err := json.Unmarshal([]byte(e.Data), &value)
		if err != nil {
			return nil, err
		}
		resp[i] = &Event{
			State: State(e.Event),
			Value: value,
		}
	}
	return resp, nil
}
