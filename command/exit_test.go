package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExit_BuildMessage(t *testing.T) {
	assert.Equal(t, "exit", Exit{}.BuildMessage())
}
