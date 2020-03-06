package rest

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetId(t *testing.T) {
	u, err := getId("/user", "/user/253ACCB1-4C4B-4F3A-8261-AB5CC8725EF8")
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("253ACCB1-4C4B-4F3A-8261-AB5CC8725EF8"), u)
}
