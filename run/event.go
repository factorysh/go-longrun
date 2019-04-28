package run

type State string

const (
	QUEUED   = State("queued")
	RUNNING  = State("running")
	CANCELED = State("canceled")
	ERROR    = State("error")
	SUCCESS  = State("success")
)

type Event struct {
	State State       `json:"state"`
	Value interface{} `json:"value,omitempty"`
}

func (e *Event) Ended() bool {
	return e.State == CANCELED || e.State == SUCCESS || e.State == ERROR
}
