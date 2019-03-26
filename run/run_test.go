package run

import (
	"encoding/json"
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
	j, err := json.Marshal(evts)
	assert.NoError(t, err)
	fmt.Println(string(j))
}
