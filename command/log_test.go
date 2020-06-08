package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLog_BuildMessage(t *testing.T) {
	assert.Equal(t, "log 9", Log{Enabled: true, Level: 9}.BuildMessage())
}

func TestNoLog_BuildMessage(t *testing.T) {
	assert.Equal(t, "nolog", Log{}.BuildMessage())
}
