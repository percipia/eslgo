package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect_BuildMessage(t *testing.T) {
	assert.Equal(t, "connect", Connect{}.BuildMessage())
}
