package sse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"

	"github.com/factorysh/go-longrun/run"
)

func TestSSE(t *testing.T) {
	runs := run.New(5 * time.Minute)
	r := runs.New()
	go func() {
		time.Sleep(2 * time.Second)
		r.Run("pim")
		r.Run("pam")
		time.Sleep(time.Second)
		r.Run("poum")
		time.Sleep(2 * time.Second)
		r.Success(new(interface{}))
	}()
	s := New(runs)
	ts := httptest.NewServer(s)
	defer ts.Close()
	fmt.Println(r.Id())
	res, err := http.Get(fmt.Sprintf("%s/%s", ts.URL, r.Id().String()))
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	reader := bufio.NewReader(res.Body)
	defer res.Body.Close()
	cpt := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		fmt.Print(line)
		if strings.HasPrefix(line, "data: ") {
			var evt run.Event
			err = json.Unmarshal([]byte(line[6:]), &evt)
			assert.NoError(t, err)
			fmt.Println(evt)
			assert.Equal(t, []run.State{run.QUEUED,
				run.RUNNING, run.RUNNING, run.RUNNING,
				run.SUCCESS}[cpt], evt.State)
			cpt++
		}
	}

}
