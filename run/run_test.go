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
		time.Sleep(500 * time.Millisecond)
		run.Run("hop")
		run.Run("pim")
		run.Run("pam")
		run.Run("poum")
		time.Sleep(200 * time.Millisecond)
		run.Cancel()
	}()
	i := 0
	for {
		r, ok := runs.GetRun(run.id)
		assert.True(t, ok)
		evts, err := r.Since(i)
		if len(evts) == 0 {
			continue
		}
		i += len(evts)
		assert.NoError(t, err)
		if evts[len(evts)-1].Ended() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
