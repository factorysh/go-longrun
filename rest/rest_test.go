package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/factorysh/go-longrun/run"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetId(t *testing.T) {
	u, err := getId("/user", "/user/253ACCB1-4C4B-4F3A-8261-AB5CC8725EF8")
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("253ACCB1-4C4B-4F3A-8261-AB5CC8725EF8"), u)
}

func TestUrl(t *testing.T) {
	mux := http.NewServeMux()
	h := &Handler{
		runs: run.New(30 * time.Second),
		root: "/user",
	}
	r := h.runs.New()
	mux.Handle("/user/", h)
	server := httptest.NewServer(mux)
	client := server.Client()
	resp, err := client.Head(server.URL + "/user/" + r.Id().String())
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp, err = client.Head(server.URL + "/user/7D5E3597-9B8C-41A3-AB1C-DA6EAC94A7B8")
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	req, err := http.NewRequest("POST", server.URL+"/user/", bytes.NewBuffer([]byte(`{"name": "Bob"}`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	req, err = http.NewRequest("POST", server.URL+"/user/", bytes.NewBuffer([]byte(`{"name": "Charly"}`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "text/event-stream")
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()
	rjson := make(map[string]string)
	json.Unmarshal(body, &rjson)
	fmt.Println("Body: ", string(body))
	assert.Equal(t, 303, resp.StatusCode)
	assert.Equal(t, "/user/"+rjson["id"], resp.Header.Get("location"))
}
