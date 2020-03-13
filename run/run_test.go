package run

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	runs := New(time.Hour)
	ctx := context.TODO()
	run := runs.NewRun(ctx)
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
		r, ok := runs.GetRun(run.id)
		assert.True(t, ok)
		i += r.Events.Size()
		stop := false
		evts := r.Events.Since(0)
		stop = stop || evts[len(evts)-1].Event == "canceled"
		if stop {
			break
		}
	}
}
