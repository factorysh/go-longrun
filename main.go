package main

import (
	"encoding/json"
	"errors"

	"github.com/bitwurx/jrpc2"
	"github.com/factorysh/go-longrun/longrun"
)

// This struct is used for unmarshaling the method params
type AddParams struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
}

// Each params struct must implement the FromPositional method.
// This method will be passed an array of interfaces if positional parameters
// are passed in the rpc call
func (ap *AddParams) FromPositional(params []interface{}) error {
	if len(params) != 2 {
		return errors.New("exactly two integers are required")
	}

	x := params[0].(float64)
	y := params[1].(float64)
	ap.X = &x
	ap.Y = &y

	return nil
}

// Each method should match the prototype <fn(json.RawMessage) (inteface{}, *ErrorObject)>
func Add(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(AddParams)

	// ParseParams is a helper function that automatically invokes the FromPositional
	// method on the params instance if required
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}

	if p.X == nil || p.Y == nil {
		return nil, &jrpc2.ErrorObject{
			Code:    jrpc2.InvalidParamsCode,
			Message: jrpc2.InvalidParamsMsg,
			Data:    "exactly two integers are required",
		}
	}

	return *p.X + *p.Y, nil
}

func main() {
	// create a new server instance
	s := jrpc2.NewServer(":8888", "/api/v1/rpc", map[string]string{})

	// register the add method
	s.Register("add", jrpc2.Method{Method: Add})

	l := longrun.New()
	s.Register("longrun.next", jrpc2.Method{Method: l.Next})
	// register the subtract method to proxy another rpc server
	// s.Register("add", jrpc2.Method{Url: "http://localhost:9999/api/v1/rpc"})

	// start the server instance
	s.Start()
}
