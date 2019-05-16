package sse

import (
	"bufio"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"

	"github.com/factorysh/go-longrun/run"
)

func fixture(r *run.Run) {
	go func() {
		time.Sleep(2 * time.Second)
		r.Run("pim")
		r.Run("pam")
		time.Sleep(time.Second)
		r.Run("poum")
		time.Sleep(2 * time.Second)
		r.Success(new(interface{}))
	}()
}

func TestSSE(t *testing.T) {
	runs := run.New(5 * time.Minute)
	r := runs.New()
	fixture(r)
	s := New(runs)
	ts := httptest.NewServer(s)
	defer ts.Close()
	res, err := http.Get(fmt.Sprintf("%s/%s", ts.URL, r.Id().String()))
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	defer res.Body.Close()
	cpt := 0
	err = Reader(res.Body, func(evtRaw *Event) error {
		fmt.Println(evtRaw)
		var evt run.Event
		err := evtRaw.JSON(&evt)
		if err != nil {
			return err
		}
		assert.Equal(t, []run.State{run.QUEUED,
			run.RUNNING, run.RUNNING, run.RUNNING,
			run.SUCCESS}[cpt], evt.State)
		cpt++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, cpt)
}

func TestLastEventId(t *testing.T) {
	runs := run.New(5 * time.Minute)
	r := runs.New()
	fixture(r)
	s := New(runs)
	ts := httptest.NewServer(s)
	defer ts.Close()
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/%s", ts.URL, r.Id().String()), nil)
	req.Header.Set("Last-Event-Id", "2")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	reader := bufio.NewReader(res.Body)
	defer res.Body.Close()
	evts := make([]*run.Event, 0)
	err = Reader(reader, func(evtRaw *Event) error {
		var e run.Event
		err = evtRaw.JSON(&e)
		assert.NoError(t, err)
		evts = append(evts, &e)
		return nil
	})
	assert.NoError(t, err)
	assert.Len(t, evts, 2)
	assert.Equal(t, "poum", evts[0].Value)
	assert.Equal(t, run.SUCCESS, evts[1].State)

}
