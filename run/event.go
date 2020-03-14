package run

import "encoding/json"

type State string

const (
	QUEUED   = State("queued")
	RUNNING  = State("running")
	CANCELED = State("canceled")
	ERROR    = State("error")
	SUCCESS  = State("success")
)

type Event struct {
	Id    int         `json:"id"`
	State State       `json:"state"`
	Value interface{} `json:"value,omitempty"`
}

func (e *Event) Ended() bool {
	return e.State == CANCELED || e.State == SUCCESS || e.State == ERROR
}

func Parse(event string, data []byte) (*Event, error) {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		return nil, err
	}
	return &Event{
		State: State(event),
		Value: value,
	}, nil
}
