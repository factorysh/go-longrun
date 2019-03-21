package run

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	runs := New()
	run := runs.New()
	run.Run("hop")
	run.Cancel()
	evts, err := runs.Get(run.id, 0)
	assert.NoError(t, err)
	assert.Len(t, evts, 3)
	fmt.Println(evts)
}
