package longrun

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongrun(t *testing.T) {

	lr := New()
	run := lr.runs.New()

	args, err := json.Marshal(map[string]interface{}{
		"id": run.Id().String(),
	})
	assert.NoError(t, err)
	r, jerr := lr.Next(json.RawMessage(args))
	assert.Nil(t, jerr)
	jr, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`[{"state":"queued"}]`), jr)

}
