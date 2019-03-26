package longrun

import (
	"encoding/json"
	"errors"

	"github.com/factorysh/go-longrun/run"

	"github.com/google/uuid"

	"github.com/bitwurx/jrpc2"
)

type NextParams struct {
	Id  string `json:"id"`
	uid uuid.UUID
	N   int `json:"n"`
}

func (np *NextParams) FromPositional(params []interface{}) error {
	if len(params) == 0 {
		return errors.New("At least one argument")
	}

	np.Id = params[0].(string)
	if len(params) == 2 {
		np.N = params[1].(int)
	}
	return nil
}

type LongRun struct {
	Runs *run.Runs
}

func New() *LongRun {
	return &LongRun{run.New()}
}

func (lr *LongRun) Next(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(NextParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(p.Id)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.ErrorMsg("Not an UUID " + err.Error()),
		}
	}
	events, err := lr.Runs.Get(uid, p.N)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InternalErrorCode,
			Message: jrpc2.ErrorMsg(err.Error()),
		}
	}
	return events, nil
}
