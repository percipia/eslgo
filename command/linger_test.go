package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoLinger_BuildMessage(t *testing.T) {
	assert.Equal(t, "nolinger", Linger{}.BuildMessage())
}

func TestLinger_BuildMessage(t *testing.T) {
	assert.Equal(t, "linger", Linger{true}.BuildMessage())
}
