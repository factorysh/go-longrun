package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/factorysh/go-longrun/rest"
	"github.com/factorysh/go-longrun/run"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	mux := http.NewServeMux()

	runs := run.New(30 * time.Second)
	h := rest.NewHandler(runs, "/user", func(r *run.Run, req *http.Request, arg map[string]interface{}) {
		r.Run(arg["name"])
		time.Sleep(100 * time.Millisecond)
		r.Success(nil)
	})
	mux.Handle("/user/", h)
	server := httptest.NewServer(mux)
	client := server.Client()

	req, err := http.NewRequest("POST", server.URL+"/user/", bytes.NewBuffer([]byte(`{"name": "Charly"}`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "text/event-stream")
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	fmt.Println(resp.Request.URL)
	eventsReader, err := Longrun(resp)
	assert.NoError(t, err)
	cpt := 0
	for {
		evt, err := eventsReader.Read()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		cpt++
		fmt.Println("Read event: ", evt)
	}
	assert.Equal(t, 3, cpt)
}
