package run

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	runs := New()
	run := runs.New()
	go func() {
		time.Sleep(5 * time.Second)
		run.Run("hop")
		run.Run("pim")
		run.Run("pam")
		run.Run("poum")
		time.Sleep(2 * time.Second)
		run.Cancel()
	}()
	i := 0
	for {
		evts, err := runs.Get(run.id, i)
		assert.NoError(t, err)
		i += len(evts)
		stop := false
		for _, evt := range evts {
			j, err := json.Marshal(evt)
			assert.NoError(t, err)
			fmt.Println(i, string(j))
			stop = stop || evt.Ended()
		}
		if stop {
			break
		}
	}
}
