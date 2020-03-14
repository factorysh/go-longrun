package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/factorysh/go-longrun/run"
	"github.com/factorysh/go-longrun/sse"
)

type Reader struct {
	reader *sse.SSEReader
	ended  bool
}

func Longrun(resp *http.Response) (*Reader, error) {
	if resp.Header.Get("content-type") != "text/event-stream" {
		return nil, fmt.Errorf("Wrong content-type : %s", resp.Header.Get("content-type"))
	}
	return &Reader{
		reader: sse.NewSSEReader(resp.Body),
	}, nil
}

func (r *Reader) Read() (*run.Event, error) {
	if r.ended {
		return nil, io.EOF
	}
	raw, err := r.reader.Read()
	if err != nil {
		return nil, err
	}
	evt, err := run.Parse(raw.Event, []byte(raw.Data))
	if err != nil {
		return nil, err
	}
	r.ended = evt.Ended()
	return evt, nil
}
